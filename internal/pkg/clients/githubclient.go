package clients

import (
	"io"
	"net/http"
	"strings"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
)

type GithubSolutionsResponse struct {
	ProblemName      string
	ProblemOriginUrl string
	IsSolvedProblem  bool
}

const (
	ProblemNameCSVIndex      = 0
	ProblemOriginUrlCSVIndex = 1
	IsSolvedProblemCSVIndex  = 2
	githubSolvedProblemsURL  = "https://raw.githubusercontent.com/DushnoAndTochka/solutions_algorithmic_problems/main/solved_problems.csv"
)

func (cli *Client) GetSolutionsList() ([]*GithubSolutionsResponse, error) {
	logger := log.GetLogger()

	req, err := http.NewRequest("GET", githubSolvedProblemsURL, nil)
	if err != nil {
		logger.Error("GithubClient NewRequest err: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		logger.Error("GithubClient DoRequest err: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	return parseProblemSolutionResponseBody(resp)
}

func parseProblemSolutionResponseBody(resp *http.Response) ([]*GithubSolutionsResponse, error) {
	logger := log.GetLogger()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("GithubClient ReadBody err: %s", err)
		return nil, err
	}

	responseRows := strings.Split(string(respBody[:]), "\r\n")

	var solutionsResponse []*GithubSolutionsResponse

	for _, row := range responseRows {
		rowElements := strings.Split(row, ",")

		if len(rowElements) != 3 || rowElements[0] == "problem_name" {
			continue
		}

		solutionResponse := &GithubSolutionsResponse{
			ProblemName:      rowElements[ProblemNameCSVIndex],
			ProblemOriginUrl: rowElements[ProblemOriginUrlCSVIndex],
			IsSolvedProblem:  rowElements[IsSolvedProblemCSVIndex] == "True",
		}
		solutionsResponse = append(solutionsResponse, solutionResponse)
	}

	return solutionsResponse, nil
}

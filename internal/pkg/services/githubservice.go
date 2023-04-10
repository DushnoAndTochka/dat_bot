package services

import (
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/clients"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/customerrors"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/storages"
)

func SyncGithubSolutions() {
	storage := storages.GetStorage()
	logger := log.GetLogger()
	cli := clients.GetClient()

	githubSolutions, err := cli.GetSolutionsList()
	if err != nil {
		logger.Error("GithubService GetSolutionsList failed: %w", err)
		return
	}

	solutions, err := storage.SolutionsGetAll()
	if err != nil {
		logger.Error("GithubService SolutionsGetAll failed: %w", err)
		return
	}

	processGithubSolutions(githubSolutions, solutions)
}

func processGithubSolutions(
	githubSolutions []*clients.GithubSolutionsResponse, solutions map[string]*models.Solution) {
	storage := storages.GetStorage()
	logger := log.GetLogger()

	var knownSolution *models.Solution
	var ok bool

	for _, githuSolution := range githubSolutions {
		if githuSolution == nil {
			continue
		}
		knownSolution, ok = solutions[githuSolution.ProblemName]
		if ok {
			if knownSolution.IsSolved == githuSolution.IsSolvedProblem {
				continue
			} else {
				knownSolution.IsSolved = githuSolution.IsSolvedProblem
				err := storage.SolutionUpdateOrCreate(knownSolution)
				if err != nil {
					logger.Error("GithubService SyncGithub: Fail to SolutionUpdateOrCreate: %w", err)
				}
			}
		} else {
			problem, err := models.NewProblemFromUrl(githuSolution.ProblemOriginUrl)
			if err != nil {
				logger.Error("GithubService SyncGithub: Fail to init problem: %s", customerrors.CustomErrors[err])
				continue
			}
			err = storage.ProblemGetOrCreate(problem)
			if err != nil {
				logger.Error("GithubService SyncGithub: Fail to GetOrCreate problem: %w", err)
				continue
			}
			knownSolution = &models.Solution{
				Name:      githuSolution.ProblemName,
				IsSolved:  githuSolution.IsSolvedProblem,
				ProblemID: problem.ID,
			}
			err = storage.SolutionUpdateOrCreate(knownSolution)
			if err != nil {
				logger.Error("GithubService SyncGithub: Fail to SolutionUpdateOrCreate: %w", err)
			}
		}
	}
}

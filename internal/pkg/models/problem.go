package models

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RegexForSourceProblem string

const (
	LeetCodeRegexProblemUrl RegexForSourceProblem = `^https:\/\/leetcode.com\/problems\/(?P<problem_name>[a-z\-]*)\/?$`
)

type ProblemSourceFromSource string

const (
	LeetCodeUrl ProblemSourceFromSource = "https://leetcode.com/problems/$1/"
)

type ProblemSource string

const (
	LeetCodeSource ProblemSource = "leetcode"
)

var ProblemSources = map[ProblemSource]RegexForSourceProblem{
	LeetCodeSource: LeetCodeRegexProblemUrl,
}

var ProblemUrlFormSources = map[ProblemSource]ProblemSourceFromSource{
	LeetCodeSource: LeetCodeUrl,
}

type Problem struct {
	ID     uuid.UUID
	Name   string
	Source ProblemSource
}

func NewProblemFromUrl(url string) (*Problem, error) {
	var err error
	var problemName string

	for problemSource, regex := range ProblemSources {
		re := regexp.MustCompile(string(regex))
		matches := re.FindStringSubmatch(url)
		if len(matches) == 0 {
			err = fmt.Errorf("ProblemModels: Cant parse url.")
			continue
		}
		problemNameIndex := re.SubexpIndex("problem_name")
		problemName = matches[problemNameIndex]
		problem := &Problem{
			Name:   problemName,
			Source: problemSource,
		}
		return problem, nil
	}
	return nil, err
}

func (u *Problem) ScanProblemRow(rows pgx.Row) error {
	return rows.Scan(
		&u.ID,
		&u.Name,
		&u.Source,
	)
}

func (u *Problem) ScanProblemRows(rows pgx.Rows) error {
	return rows.Scan(
		&u.ID,
		&u.Name,
		&u.Source,
	)
}

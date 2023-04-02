package models

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	LeetCodeRegexProblemUrl = regexp.MustCompile(`^https:\/\/leetcode.com\/problems\/(?P<problem_name>[a-z\-]*)\/?$`)
)

type problemSourceFromSource string

const (
	LeetCodeUrl problemSourceFromSource = "https://leetcode.com/problems/$1/"
)

type problemSource string

const (
	LeetCodeSource problemSource = "leetcode"
)

var ProblemSources = map[problemSource]*regexp.Regexp{
	LeetCodeSource: LeetCodeRegexProblemUrl,
}

var ProblemUrlFormSources = map[problemSource]problemSourceFromSource{
	LeetCodeSource: LeetCodeUrl,
}

type ErrorProblem struct {
	Err        error
	ErrMessage string
}

var (
	NotSupportedURL *ErrorProblem = &ErrorProblem{
		Err:        fmt.Errorf("notSupported"),
		ErrMessage: "Данный URL не поддерживается",
	}
)

type problemStatus string

const (
	CloseStatus = "Close"
	OpenStatus  = "Open"
)

type Problem struct {
	ID     uuid.UUID
	Name   string
	Source problemSource
	Status problemStatus
}

func NewProblemFromUrl(url string) (*Problem, error) {
	var err error
	var problemName string

	for problemSource, regex := range ProblemSources {
		matches := regex.FindStringSubmatch(url)
		if len(matches) == 0 {
			err = NotSupportedURL.Err
			continue
		}
		problemNameIndex := regex.SubexpIndex("problem_name")
		problemName = matches[problemNameIndex]
		problem := &Problem{
			Name:   problemName,
			Source: problemSource,
			Status: OpenStatus,
		}
		return problem, nil
	}
	return nil, err
}

func (p *Problem) ScanProblemRow(rows pgx.Row) error {
	return rows.Scan(
		&p.ID,
		&p.Name,
		&p.Source,
		&p.Status,
	)
}

func (p *Problem) ScanProblemRows(rows pgx.Rows) error {
	return rows.Scan(
		&p.ID,
		&p.Name,
		&p.Source,
		&p.Status,
	)
}

package models

import (
	"fmt"
	"regexp"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/customerrors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// регексы для матчинга урлов.
var (
	LeetCodeRegexProblemUrl = regexp.MustCompile(`^https:\/\/leetcode.com\/problems\/(?P<problem_name>[a-z\-0-9]*)\/?$`)
)

// типы урлов, которые могут быть обработанны и из них будут получены названия проблем
type problemSourceFromSource string

const (
	LeetCodeUrl problemSourceFromSource = "https://leetcode.com/problems/%s/"
)

// типа ресурсов проблем
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

// допустимые статусы для проблем
type problemStatus string

const (
	CloseStatus problemStatus = "Close"
	OpenStatus  problemStatus = "Open"
)

var EnableStatuses = []problemStatus{OpenStatus, CloseStatus}

type ProblemName string

type Problem struct {
	ID               uuid.UUID
	Name             ProblemName
	Source           problemSource
	Status           problemStatus
	CountSuggestions CountSuggestions
}

func NewProblem(id uuid.UUID, name string, source string, status string) (*Problem, error) {
	pSource := problemSource(source)
	_, ok := ProblemSources[pSource]
	if !ok {
		return nil, fmt.Errorf("NotCorrectSource")
	}

	pStatus := problemStatus(status)
	ok = false
	for _, enambleStatus := range EnableStatuses {
		if pStatus == enambleStatus {
			ok = true
			break
		}
	}

	if !ok {
		return nil, fmt.Errorf("NotCorrectStatus")
	}

	return &Problem{
		ID:     id,
		Name:   ProblemName(name),
		Source: pSource,
		Status: pStatus,
	}, nil
}


// инициализаци проблемы из представленного URL. Проверяет валидность урла.
func NewProblemFromUrl(url string) (*Problem, error) {
	var err error
	var problemName string

	for problemSource, regex := range ProblemSources {
		matches := regex.FindStringSubmatch(url)
		if len(matches) == 0 {
			err = customerrors.ErrNotSupportedURL
			continue
		}
		problemNameIndex := regex.SubexpIndex("problem_name")
		problemName = matches[problemNameIndex]
		problem := &Problem{
			Name:   ProblemName(problemName),
			Source: problemSource,
			Status: OpenStatus,
		}
		return problem, nil
	}
	return nil, err
}

func (p *Problem) GetOriginalUrl() string {
	source, ok := ProblemUrlFormSources[problemSource(p.Source)]
	if ok {
		return fmt.Sprintf(string(source), string(p.Name))
	} else {
		return ""
	}
}

func (p *Problem) ScanProblemRow(row pgx.Row) error {
	return row.Scan(
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

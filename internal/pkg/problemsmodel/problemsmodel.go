package problemsmodel

import (
	"context"
	"fmt"
	"regexp"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
)

func NewProblem(url string) (*ProblemModel, error) {
	re := regexp.MustCompile(`^https:\/\/leetcode.com\/problems\/(?P<problem_name>[a-z\-]*)\/?$`)
	matches := re.FindStringSubmatch(url)
	problemNameIndex := re.SubexpIndex("problem_name")
	if len(matches) == 0 {
		return nil, fmt.Errorf("ProblemModels: Cant parse url.")
	}

	name := matches[problemNameIndex]
	problem := &ProblemModel{
		Name: name,
		url:  url,
	}
	return problem, nil
}

func GetByName(name string) (*ProblemModel, error) {
	conn := dbconnection.GetPoolConnections()
	var id int
	var url string
	err := conn.QueryRow(context.Background(), "select id, url from widgets where id=$1", 42).Scan(&id, &url)
	if err != nil {
		return nil, err
	}

	problem := &ProblemModel{
		ID:   id,
		Name: name,
		url:  url,
	}

	return problem, nil

}

func GetAll() ([]*ProblemModel, error) {

}

func GetByUrl(url string) (*ProblemModel, error) {
	conn := dbconnection.GetPoolConnections()
	var id int
	var name string
	err := conn.QueryRow(context.Background(), "select id, url from problems where id=$1", 42).Scan(&id, &name)
	if err != nil {
		return nil, err
	}

	problem := &ProblemModel{
		ID:   id,
		Name: name,
		url:  url,
	}

	return problem, nil

}

type ProblemModel struct {
	ID   int
	Name string
	url  string
}

func (p *ProblemModel) Create() error {
	conn := dbconnection.GetPoolConnections()
	var id int
	err := conn.QueryRow(context.Background(), "INSERT INTO problems (name, url) VALUES ($1, $2);", p.Name, &p.url).Scan(&id)
	p.ID = id
	return err
}

package proposesmodel

import (
	"context"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/problemsmodel"
)

func NewPropose(problem *problemsmodel.ProblemModel, tgUuid int) *ProposeModel {
	// problem, err := problemsmodel.NewProblem(url)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = problemsmodel.GetByName(problem.Name)

	// if err != nil {
	// 	err = problem.Create()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	propose := ProposeModel{
		TgUuid:    tgUuid,
		proposeId: problem.ID,
	}
	return &propose
}

func GetByProblem(problem *problemsmodel.ProblemModel) (*ProposeModel, error) {
	conn := dbconnection.GetPoolConnections()
	var id int
	var problemId int
	var tgUuid int
	err := conn.QueryRow(context.Background(), "select id, problem_id as problemId, tg_uuid as tgUuid from widgets where id=$1", 42).Scan(&id, &tgUuid, &problemId)
	if err != nil {
		return nil, err
	}

	propose := ProposeModel{
		TgUuid:    tgUuid,
		proposeId: problem.ID,
	}

	return &propose, nil
}

type ProposeModel struct {
	ID        int
	TgUuid    int
	proposeId int
}

func (p *ProposeModel) Create() error {
	conn := dbconnection.GetPoolConnections()
	var id int
	err := conn.QueryRow(context.Background(), "INSERT INTO problems (name, url) VALUES ($1, $2);", p.TgUuid, &p.proposeId).Scan(&id)
	p.ID = id
	return err
}

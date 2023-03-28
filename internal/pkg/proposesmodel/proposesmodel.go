package proposesmodel

import (
	"context"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/problemsmodel"
)

func NewPropose(problem *problemsmodel.ProblemModel, tgUuid int64) *ProposeModel {
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
		problemID: problem.ID,
	}
	return &propose
}

func GetByProblem(problem *problemsmodel.ProblemModel) (*ProposeModel, error) {
	conn := dbconnection.GetPoolConnections()
	var id int
	var problemId int
	var tgUuid int64
	err := conn.QueryRow(context.Background(), "select id, problem_id as problemId, tg_uuid as tgUuid from widgets proposes id=$1", 42).Scan(&id, &tgUuid, &problemId)
	if err != nil {
		return nil, err
	}

	propose := ProposeModel{
		TgUuid:    tgUuid,
		problemID: problem.ID,
	}

	return &propose, nil
}

func GetByUuid(tgUuid int64) (*ProposeModel, error) {
	conn := dbconnection.GetPoolConnections()
	var id int
	var problemId int

	err := conn.QueryRow(
		context.Background(),
		`SELECT id, problem_id as problemId, tg_uuid as tgUuid 
		FROM proposes WHERE tg_uuid=$1`, tgUuid).Scan(&id, &tgUuid, &problemId)

	if err != nil {
		return nil, err
	}

	if problemId == 0 {
		return nil, nil
	}
	problem, err := problemsmodel.GetByID(problemId)

	if err != nil {
		return nil, err
	}

	propose := ProposeModel{
		TgUuid:    tgUuid,
		problemID: problemId,
		ID:        id,
		Problem:   problem,
	}

	return &propose, nil
}

type ProposeModel struct {
	ID        int
	TgUuid    int64
	problemID int
	Problem   *problemsmodel.ProblemModel
}

func (p *ProposeModel) Create() error {
	conn := dbconnection.GetPoolConnections()
	var id int
	err := conn.QueryRow(context.Background(), "INSERT INTO problems (name, problem_id) VALUES ($1, $2);", p.TgUuid, &p.problemID).Scan(&id)
	p.ID = id
	return err
}

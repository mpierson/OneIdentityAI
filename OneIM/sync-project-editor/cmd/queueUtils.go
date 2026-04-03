package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type JobQueue struct {
	oneim.Specials
	UID_Job        string
	JobChainName   string
	Ready2EXE      string
	BasisObjectKey string
	Priority       int
	WasError       bool
	StartAt        sql.NullTime
	ErrorMessages  *string
	ParamIN        *string
	Parameters     []string `mapstructure:",omitzero"`
	FilePath       *string  `mapstructure:"-"`
}

func FireDBEvent(db *sqlx.DB, objectType string, whereClause string, eventName string, priority int) error {

	genProcId, err := dbx.GetNewId(db)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`exec QBM_PJobCreate_HOFireEvent
        @objecttype = '%s',
        @whereclause = '%s' , 
		@EventName = '%s',
        @priority = %v,
        @GenProcID = '%s'`, objectType, whereClause, eventName, priority, genProcId)
	_, err = db.Exec(q)

	return err
}

func WaitForTaskStart(db *sqlx.DB, ctx context.Context, whereClause string, interval time.Duration) error {

	jobQueueCheck := func(db *sqlx.DB) (bool, error) {
		nTasks, err := dbx.GetTableCount(db, "JobQueue", whereClause)
		if err != nil {
			return false, err
		}
		return (nTasks > 0), nil
	}
	return dbx.WaitForDBResult(db, ctx, jobQueueCheck, interval)
}

func WaitForTaskFinish(db *sqlx.DB, ctx context.Context, whereClause string, interval time.Duration) error {
	jobQueueCheck := func(db *sqlx.DB) (bool, error) {
		nTasks, err := dbx.GetTableCount(db, "JobQueue", whereClause)
		if err != nil {
			return false, err
		}
		return (nTasks == 0), nil
	}
	return dbx.WaitForDBResult(db, ctx, jobQueueCheck, interval)
}

func GetTask(db *sqlx.DB, whereClause string) (*JobQueue, error) {
	t, err := dbx.GetStructSingletonByWC[JobQueue](db, whereClause)
	return &t, err
}

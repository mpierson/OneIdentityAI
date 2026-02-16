package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DialogSchedule struct {
	oneim.Specials
	Displayable
	UID_DialogSchedule       string
	UID_DialogTableBelongsTo string
	UID_DialogTimeZone       string
	Frequency                int
	FrequencyType            *string
	FrequencySubType         *string
	StartTime                string
	StartDate                *time.Time
	Table                    string `mapstructure:",omitzero"`
	TimeZone                 string `mapstructure:",omitzero"`
}

type JobAutoStart struct {
	oneim.Specials
	Displayable
	UID_JobAutoStart   string
	UID_DialogSchedule string
	UID_QBMEvent       string
	ObjectKeyTarget    string
	WhereClause        *string
	Schedule           string `mapstructure:",omitzero"`
}

var ScheduleCmd = CreateBaseCommand(
	"schedule",
	"synchronization schedule commands",
	`View and update synchronization schedule (DialogSchedule).`,
	showSchedules,
)

var ShowScheduleCmd = CreateShowCommand(
	"show details of one synchronization schedule",
	`View sync project schedule (DialogSchedule).`,
	nil,
	showSchedules,
)

func showSchedules(c *cobra.Command, db *sqlx.DB) error {

	id, _ := c.Flags().GetString("id")

	wc := "UID_DialogTableBelongsTo = (select UID_DialogTable from DialogTable where TableName='JobAutoStart')"
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DialogSchedule='%s'`, id)
	}

	return ShowDPRObjects[DialogSchedule](db, wc, fillScheduleData)
}

func fillScheduleData(db *sqlx.DB, t *DialogSchedule) error {

	t.Table, _ = GetForeignDisplay(db, "DialogTable", t.UID_DialogTableBelongsTo)
	tz, _ := dbx.GetForeignSingleton(db, "DialogTimeZone", "UID_DialogTimeZone", t.UID_DialogTimeZone)
	t.TimeZone = tz["ShortName"].(string)

	return nil
}

var InsertScheduleCmd = CreateInsertCommand(
	"create a new synchronization Schedule",
	`Create a new sync Schedule (DialogSchedule).`,
	[]string{"name"},
	insertSchedule,
)

func insertSchedule(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DialogSchedule](c, db, "", newSchedule_cmd)
}

func newSchedule_cmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DialogSchedule, error) {

	tzId, err := c.Flags().GetString("time-zone-id")
	if len(tzId) == 0 {
		tzId, err = dbx.GetTableValue(db, "DialogTimeZone", "UID_DialogTimeZone", "ShortName='UTC'")
		if err != nil {
			return nil, err
		}
	}

	return NewSchedule(db, id, objectKey, name, clrId, tzId)
}

func NewSchedule(db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
	timeZoneId string,
) (*DialogSchedule, error) {

	tableId, err := dbx.GetTableValue(db, "DialogTable", "UID_DialogTable", "TableName='JobAutoStart'")
	if err != nil {
		return nil, err
	}

	description := "Created by SPEd"
	startDate := time.Now().UTC()

	t := DialogSchedule{
		UID_DialogSchedule:       id,
		UID_DialogTimeZone:       timeZoneId,
		UID_DialogTableBelongsTo: tableId,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		Displayable: Displayable{
			Name:        &name,
			Description: &description,
		},
		StartDate: &startDate,
	}

	return &t, nil
}

func InsertNewSchedule(db *sqlx.DB, name string, timeZoneId string) (*DialogSchedule, error) {

	fNew := func(db *sqlx.DB, id string, objectKey string, name string, clrId string) (*DialogSchedule, error) {
		return NewSchedule(db, id, objectKey, name, clrId, timeZoneId)
	}

	t, _, err := InsertNewDPRObject[DialogSchedule](db, name, "", fNew)
	return t, err

}

func GetTimeZoneIdByName(db *sqlx.DB, name string) (string, error) {
	if !IsValidIdOrName(name) {
		return "", errors.New("invalid time zone: " + name)
	}
	wc := fmt.Sprintf(`ShortName='%s'`, name)
	return dbx.GetTableValue(db, "DialogTimeZone", "UID_DialogTimeZone", wc)
}

func InsertNewJobAutoStart(db *sqlx.DB,
	startInfoName string, startInfoObjectKey string,
	scheduleId string,
) (*JobAutoStart, error) {

	// event is same for all sync projects
	wc := "DisplayName = 'RUN - DPRProjectionStartInfo'"
	eventId, err := dbx.GetTableValue(db, "QBMEvent", "UID_QBMEvent", wc)

	// trigger only for given start info
	jobStartWC := fmt.Sprintf(`XObjectKey='%s'`, startInfoObjectKey)

	description := "Created by SPEd"

	fNew := func(db *sqlx.DB, id string, objectKey string, name string, clrId string) (*JobAutoStart, error) {
		return NewJobAutoStart(db, id, objectKey, name, clrId, scheduleId, eventId, startInfoObjectKey, jobStartWC, description)
	}

	autoStartName := fmt.Sprintf(`Synchronization of %s`, startInfoName)

	t, _, err := InsertNewDPRObject[JobAutoStart](db, autoStartName, "", fNew)
	return t, err
}

func NewJobAutoStart(db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
	scheduleId string,
	eventId string,
	startInfoObjectKey string,
	whereClause string,
	description string,
) (*JobAutoStart, error) {

	t := JobAutoStart{
		UID_JobAutoStart:   id,
		UID_DialogSchedule: scheduleId,
		UID_QBMEvent:       eventId,
		ObjectKeyTarget:    startInfoObjectKey,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		Displayable: Displayable{
			Name:        &name,
			Description: &description,
		},
	}

	if len(whereClause) > 0 {
		t.WhereClause = &whereClause
	}

	return &t, nil
}

var UpdateScheduleCmd = CreateUpdateCommand(
	"update an existing synchronization Schedule",
	`Update attributes of a sync project Schedule (DialogSchedule).`,
	updateSchedule,
)

func updateSchedule(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DialogSchedule](c, db)
}

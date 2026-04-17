package cmd

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRProjectionStartInfo struct {
	oneim.Specials
	Displayable
	UID_DPRProjectionStartInfo    string
	UID_DPRShell                  string
	UID_QBMClrType                string
	UID_DPRProjectionConfig       *string
	UID_DPRRootObjConnectionInfo  *string
	UID_DPRSystemVariableSet      *string
	BulkLevel                     int
	DebugMode                     bool
	FailureHandlingMode           *string
	FailureHandlingRetryCycles    int
	LastStart                     *time.Time
	LoadPartitionedThreshold      int
	MaintenanceRetryCycles        int
	MaintenanceType               *string
	PartitionSize                 int
	ProjectionDirection           *string
	RevisionHandling              *string
	StartGroupConcurrenceBehavior *string
	StartGroupName                *string
	SysConcurrenceCacheLifetime   int
	SysConcurrenceCheckMode       *string
	UseSingleProcessContextExec   bool
	Workflow                      string `mapstructure:",omitzero"`
	RootObj                       string `mapstructure:",omitzero"`
	VariableSet                   string `mapstructure:",omitzero"`
}

var STARTINFO_MaintenanceType_Invalidate = "Invalidate"

var StartInfoCmd = CreateBaseCommand(
	"start-info",
	"sync project start info commands",
	`View and update synchronization start info (DPRProjectionStartInfo).`,
	showStartInfos,
)

var ShowStartInfoCmd = CreateShowCommand(
	"show details of one synchronization start info",
	`View sync project start info (DPRProjectionStartInfo).`,
	[]string{"shell"},
	showStartInfos,
)

func showStartInfos(c *cobra.Command, db *sqlx.DB) error {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return err
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRShell='%s'`, shellId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRProjectionStartInfo='%s'`, id)
	}

	return ShowDPRObjects[DPRProjectionStartInfo](db, wc, fillStartInfoData)
}

func fillStartInfoData(db *sqlx.DB, t *DPRProjectionStartInfo) error {

	if t.UID_DPRProjectionConfig != nil {
		t.Workflow, _ = GetForeignDisplay(db, "DPRProjectionConfig", *t.UID_DPRProjectionConfig)
	}
	if t.UID_DPRSystemVariableSet != nil {
		t.VariableSet, _ = GetForeignDisplay(db, "DPRSystemVariableSet", *t.UID_DPRSystemVariableSet)
	}
	if t.UID_DPRRootObjConnectionInfo != nil {
		t.RootObj, _ = GetRootObjDisplay(db, *t.UID_DPRRootObjConnectionInfo)
	}

	return nil
}

var InsertStartInfoCmd = CreateInsertCommand(
	"create a new synchronization start info",
	`Create a new sync Connection (DPRProjectionStartInfo).`,
	[]string{"shell", "name", "direction"},
	insertStartInfo,
)

func insertStartInfo(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRProjectionStartInfo](c, db, "VI.Projector.Projection.ProjectionStartInfo", newStartInfo)
}

func newStartInfo(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRProjectionStartInfo, error) {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	variableSetId := ""
	useDefaultVars, _ := c.Flags().GetBool("use-default-variables")
	if useDefaultVars {
		variableSet, err := GetDefaultVariableSet(db, shellId)
		if err != nil {
			return nil, err
		}
		variableSetId = variableSet.UID_DPRSystemVariableSet
	} else {
		variableSetId, _ = c.Flags().GetString("variable-set-id")
	}

	wfId, _ := c.Flags().GetString("workflow-id")
	if len(wfId) > 0 {
		// validate that workflow exists
		_, err = GetStructId_MustExist[DPRProjectionConfig](c, "workflow-id", db)
		if err != nil {
			return nil, err
		}
	} else {
		// caller may have provided wf name
		wfName, _ := c.Flags().GetString("workflow-name")
		if len(wfName) > 0 {
			wf, err := GetWorkflowByName(db, shellId, wfName)
			if err != nil {
				return nil, err
			}
			if wf != nil {
				wfId = wf.UID_DPRProjectionConfig
			}
		}
	}

	direction, err := c.Flags().GetString("direction")
	if err != nil {
		return nil, err
	}

	t := DPRProjectionStartInfo{
		UID_DPRProjectionStartInfo: id,
		UID_DPRShell:               shellId,
		UID_QBMClrType:             clrId,
		Specials:                   oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &id,
			DisplayName: &name,
		},
		PartitionSize:            1024,
		BulkLevel:                1024,
		LoadPartitionedThreshold: 8,
		MaintenanceRetryCycles:   2,
		MaintenanceType:          &STARTINFO_MaintenanceType_Invalidate,
	}
	if len(variableSetId) > 0 {
		t.UID_DPRSystemVariableSet = &variableSetId
	}
	if len(wfId) > 0 {
		t.UID_DPRProjectionConfig = &wfId
	}
	if len(direction) > 0 {
		t.ProjectionDirection = &direction
	}

	return &t, nil
}

var UpdateStartInfoCmd = CreateUpdateCommand(
	"update an existing synchronization startup record",
	`Update attributes of a sync project startup (DPRProjectionStartInfo).`,
	updateStartInfo,
)

func updateStartInfo(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRProjectionStartInfo](c, db)
}

var AddScheduleToStartInfoCmd = createDPRCommand(
	"add-schedule",
	"add a new schedule and automation to a start info",
	`Add a schedule and related automation objects to an existing start info object.`,
	[]string{"shell", "id"},
	addAutomationToStartInfo,
)

func addAutomationToStartInfo(c *cobra.Command, db *sqlx.DB) error {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return err
	}

	startInfoId, err := GetStructId_MustExist[DPRProjectionStartInfo](c, "id", db)
	if err != nil {
		return err
	}
	startInfo, err := dbx.GetStructSingleton[DPRProjectionStartInfo](db, startInfoId)
	if err != nil {
		return err
	} else if startInfo.UID_DPRShell != shellId {
		return errors.New("Shell and StartInfo do not match")
	}

	// TODO: check for existing schedule / JobAutoStart

	schedule, err := addScheduleToStartInfo(c, db, &startInfo)
	if err != nil {
		return err
	}

	_, err = InsertNewJobAutoStart(db, *startInfo.DisplayName, startInfo.Specials.XObjectKey, schedule.UID_DialogSchedule)

	return nil
}

func addScheduleToStartInfo(c *cobra.Command, db *sqlx.DB, startInfo *DPRProjectionStartInfo) (*DialogSchedule, error) {

	timeZoneName, _ := c.Flags().GetString("time-zone")
	if len(timeZoneName) == 0 {
		return nil, errors.New("missing time zone")
	}
	timeZoneId, err := GetTimeZoneIdByName(db, timeZoneName)
	if err != nil {
		return nil, err
	} else if len(timeZoneId) == 0 {
		return nil, errors.New("invalid time zone: " + timeZoneName)
	}

	scheduleType, _ := c.Flags().GetString("type")
	if len(scheduleType) == 0 {
		return nil, errors.New("missing schedule type")
	}
	freq, _ := c.Flags().GetInt("frequency")
	if freq <= 0 {
		return nil, errors.New("missing schedule frequency")
	}
	startTime, _ := c.Flags().GetString("start-time")
	if len(startTime) == 0 {
		return nil, errors.New("missing start time")
	} else {
		// validate time string
		re := regexp.MustCompile("[01][0-9]:[0-5][0-9]")
		if !re.MatchString(startTime) {
			return nil, errors.New("Invalid start time")
		}
	}

	scheduleName := fmt.Sprintf(`Run of %s - %s`, *startInfo.DisplayName, scheduleType)

	fNewSched := func(db *sqlx.DB, id string, objectKey string, name string, clrId string) (*DialogSchedule, error) {
		t, err := NewSchedule(db, id, objectKey, name, clrId, timeZoneId)
		if err != nil {
			return nil, err
		}

		t.FrequencyType = &scheduleType
		t.Frequency = freq
		t.StartTime = startTime

		return t, nil
	}

	t, _, err := InsertNewDPRObject[DialogSchedule](db, scheduleName, "", fNewSched)
	return t, err
}

var AddRootObjectToStartInfoCmd = createDPRCommand(
	"add-root-object",
	"add root object to a start info",
	`Add a root object and related information to an existing start info object. (DPRRootObjConnectionInfo)`,
	[]string{"shell", "id"},
	addRootObjToStartInfo,
)

func addRootObjToStartInfo(c *cobra.Command, db *sqlx.DB) error {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return err
	}

	startInfoId, err := GetStructId_MustExist[DPRProjectionStartInfo](c, "id", db)
	if err != nil {
		return err
	}
	startInfo, err := dbx.GetStructSingleton[DPRProjectionStartInfo](db, startInfoId)
	if err != nil {
		return err
	} else if startInfo.UID_DPRShell != shellId {
		return errors.New("Shell / start info mismatch")
	}
	// TODO: check for existing root obj?

	connectionId := ""
	useDefaultConnection, _ := c.Flags().GetBool("use-default-connection")
	if useDefaultConnection {
		conn, err := GetDefaultTargetSystemConnection(db, shellId)
		if err != nil {
			return err
		}
		if conn != nil {
			connectionId = conn.UID_DPRSystemConnection
		} else {
			return errors.New("target system connection not found")
		}
	} else {
		connectionId, err = GetStructId_MustExist[DPRSystemConnection](c, "connection-id", db)
		if err != nil {
			return err
		}
	}

	variableSetId := ""
	useDefaultVariables, _ := c.Flags().GetBool("use-default-variables")
	if useDefaultVariables {
		vars, err := GetDefaultVariableSet(db, shellId)
		if err != nil {
			return err
		}
		if vars != nil {
			variableSetId = vars.UID_DPRSystemVariableSet
		} else {
			return errors.New("variable set not found")
		}
	} else {
		variableSetId, err = GetStructId_MustExist[DPRSystemVariableSet](c, "variable-set-id", db)
		if err != nil {
			return err
		}
	}

	serverId, _ := c.Flags().GetString("server-id")
	if len(serverId) > 0 {
		// verify that server exists
		serverId, err = GetStructId_MustExist[QBMServer](c, "server-id", db)
		if err != nil {
			return err
		}
	} else {
		// server name provided by caller?
		serverName, _ := c.Flags().GetString("server-name")
		if len(serverName) > 0 {
			server, err := GetServerByName(db, serverName)
			if err != nil {
				return err
			}
			serverId = server.UID_QBMServer
		}
	}
	if len(serverId) == 0 {
		return errors.New("missing server id")
	}

	rootObjKey, _ := c.Flags().GetString("root-object-key")
	if len(rootObjKey) == 0 {
		// table name provided by caller?
		tableName, _ := c.Flags().GetString("table-name")
		if len(tableName) > 0 {
			if !IsValidIdOrName(tableName) {
				return errors.New("invalid table name: " + tableName)
			}
			// lookup object key of given table
			rootObjKey, err = dbx.GetTableValue(db, "DialogTable", "XObjectKey", fmt.Sprintf("TableName='%s'", tableName))
			if err != nil {
				return err
			}
		}
	}
	if len(rootObjKey) == 0 {
		return errors.New("root object key is required")
	}
	// verify that target object exists
	_, _, err = dbx.GetSingletonTableDataByKey(db, rootObjKey)
	if err != nil {
		return err
	}

	// create the root obj
	t, err := InsertNewRootObj(db, connectionId, variableSetId, serverId, rootObjKey)
	if err != nil {
		return err
	}

	// update start info with root obj
	startInfo.UID_DPRRootObjConnectionInfo = &t.UID_DPRRootObjConnectionInfo
	err = UpdateStruct[DPRProjectionStartInfo](db, &startInfo, startInfoId)
	if err != nil {
		return err
	}

	return nil
}

var RunStartInfoCmd = createDPRCommand(
	"run",
	"run synchronization",
	`Start a synchronization as configured in given start info.`,
	[]string{"shell", "id"},
	runStartInfo,
)

func runStartInfo(c *cobra.Command, db *sqlx.DB) error {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return err
	}

	startInfoId, err := GetStructId_MustExist[DPRProjectionStartInfo](c, "id", db)
	if err != nil {
		return err
	}
	startInfo, err := dbx.GetStructSingleton[DPRProjectionStartInfo](db, startInfoId)
	if err != nil {
		return err
	} else if startInfo.UID_DPRShell != shellId {
		return errors.New("Shell / start info mismatch")
	}

	wc := fmt.Sprintf(`XObjectKey=''%s''`, startInfo.Specials.XObjectKey)
	err = FireDBEvent(db, "DPRProjectionStartInfo", wc, "RUN", 5)
	if err != nil {
		return err
	}

	// check for completed event handler (task should be present after running stored proc above)
	wc = fmt.Sprintf(`JobChainName like 'Created by QBMDBQueueProcess: fire event RUN for object type DPRProjectionStartInfo'
					  AND TaskName = 'FireGenEvent'
					  AND ParamIN like '%%%s%%'`,
		startInfo.Specials.XObjectKey)
	task, err := GetTask(db, wc)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*90)
	defer cancel()
	wc = fmt.Sprintf(`UID_Job='%s' AND Ready2EXE<>'DELETE'`, task.UID_Job)
	WaitForTaskFinish(db, ctx, wc, 2*time.Second)

	// task was successful?
	wc = fmt.Sprintf(`UID_Job='%s'`, task.UID_Job)
	task, err = GetTask(db, wc)
	if err != nil {
		return err
	} else if task.Ready2EXE == "DELETE" {
		// task is complete, check for errors
		if task.ErrorMessages != nil {
			fmt.Println("Failed to start sync: " + *task.ErrorMessages)
			return nil
		}
		// otherwise, good to carry on
	} else {
		fmt.Println("Failed to start sync: \n" + "  task state: " + task.Ready2EXE)
		return nil
	}

	// fetch sync task
	wc = fmt.Sprintf(`TaskName='FullProjection' AND BasisObjectKey='%s'`, startInfo.Specials.XObjectKey)
	task, err = GetTask(db, wc)
	if err != nil {
		return err
	}
	// wait for task to start
	wc = fmt.Sprintf(`UID_Job='%s' AND Ready2EXE<>'TRUE'`, task.UID_Job)
	err = WaitForTaskStart(db, ctx, wc, 2*time.Second)
	if err != nil {
		return err
	}

	// sync started successfully?
	wc = fmt.Sprintf(`UID_Job='%s'`, task.UID_Job)
	task, err = GetTask(db, wc)
	if err != nil {
		return err
	} else if task.Ready2EXE == "PROCESSING" {
		fmt.Println(task.UID_Job)
	} else {
		fmt.Println(fmt.Sprintf(`Synchronization failed: %v`, *task.ErrorMessages))
		return nil
	}

	return nil
}

var GetRunStatusCmd = createDPRCommand(
	"sync-status",
	"fetch status of a running sync",
	`Fetch status of a running sync for the given start info.`,
	[]string{"shell", "id", "job-id"},
	runStatus,
)

func runStatus(c *cobra.Command, db *sqlx.DB) error {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return err
	}

	startInfoId, err := GetStructId_MustExist[DPRProjectionStartInfo](c, "id", db)
	if err != nil {
		return err
	}
	startInfo, err := dbx.GetStructSingleton[DPRProjectionStartInfo](db, startInfoId)
	if err != nil {
		return err
	} else if startInfo.UID_DPRShell != shellId {
		return errors.New("Shell / start info mismatch")
	}

	jobId, _ := c.Flags().GetString("job-id")
	if !IsValidIdOrName(jobId) {
		return errors.New("invalid job id: " + jobId)
	}

	// get journal
	j, err := GetJournalMessages(db, startInfoId, jobId)
	if err != nil {
		fmt.Println("no status found")
		return err
	}
	fmt.Println(*j.ProjectionState)

	return nil
}

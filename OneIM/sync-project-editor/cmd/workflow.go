package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRProjectionConfig struct {
	oneim.Specials
	Displayable
	UID_DPRProjectionConfig string
	UID_DPRShell            string
	UID_QBMClrType          string
	ProjectionDirection     *string
	Steps                   []string `mapstructure:",omitzero"`
	Connections             []string `mapstructure:",omitzero"`
}

type DPRProjectionConfigHasConnect struct {
	oneim.Specials
	UID_DPRProjectionConfig string
	UID_DPRSystemConnection string
}

var WorkflowCmd = CreateBaseCommand(
	"workflow",
	"sync workflow commands",
	`View and update synchronization workflow records (DPRProjectionConfig).`,
	showWorkflows,
)

var ShowWorkflowCmd = CreateShowCommand(
	"show details of one synchronization workflow",
	`View sync workflow detail (DPRProjectionConfig).`,
	[]string{"shell"},
	showWorkflows,
)

func showWorkflows(c *cobra.Command, db *sqlx.DB) error {

	shellId, _ := c.Flags().GetString("shell")
	if len(shellId) == 0 {
		return errors.New("shell id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRShell='%s'`, shellId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRProjectionConfig='%s'`, id)
	}

	return ShowDPRObjects[DPRProjectionConfig](db, wc, fillWorkflowData)
}

func fillWorkflowData(db *sqlx.DB, t *DPRProjectionConfig) error {

	t.Steps, _ = GetChildDisplays(db, "DPRProjectionConfigStep", "UID_DPRProjectionConfig", t.UID_DPRProjectionConfig)
	t.Connections, _ = GetStructDisplays(
		db,
		"DPRSystemConnection",
		fmt.Sprintf(
			`UID_DPRSystemConnection in (
						select UID_DPRSystemConnection from DPRProjectionConfigHasConnect phc
						 where UID_DPRProjectionConfig = '%s')`,
			t.UID_DPRProjectionConfig,
		),
	)

	return nil
}

var InsertWorkflowCmd = CreateInsertCommand(
	"create a new synchronization workflow",
	`Create a new sync workflow (DPRProjectionConfig).`,
	[]string{"shell", "name", "direction"},
	insertWorkflow,
)

func insertWorkflow(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRProjectionConfig](c, db, "VI.Projector.Projection.ProjectionConfiguration", newWorkflow)
}

func newWorkflow(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRProjectionConfig, error) {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}
	shellName, err := dbx.GetTableValue(db, "DPRShell", "DisplayName", fmt.Sprintf(`UID_DPRShell='%s'`, shellId))
	if err != nil {
		return nil, err
	}

	displayName := fmt.Sprintf(`%s - %s`, shellName, name)

	direction, err := c.Flags().GetString("direction")
	if err != nil {
		return nil, err
	}

	t := DPRProjectionConfig{
		UID_DPRProjectionConfig: id,
		UID_DPRShell:            shellId,
		UID_QBMClrType:          clrId,
		Specials:                oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:                 &name,
			DisplayName:          &name,
			DisplayNameQualified: &displayName,
		},
	}

	if len(direction) > 0 {
		t.ProjectionDirection = &direction
	}

	return &t, nil
}

var AddConnectionToWorkflowCmd = createDPRCommand(
	"add-connection",
	"add a system connection to a workflow",
	`Add system connection to synchronization workflow (DPRProjectionConfigHasConnect).`,
	[]string{"id", "connection-id"},
	addConnectionToWorkflow_cmd,
)

func addConnectionToWorkflow_cmd(c *cobra.Command, db *sqlx.DB) error {

	id, err := GetStructId_MustExist[DPRProjectionConfig](c, "id", db)
	if err != nil {
		return err
	}

	connectionId, err := GetStructId_MustExist[DPRSystemConnection](c, "connection-id", db)
	if err != nil {
		return err
	}

	return addOneConnectionToWorkflow(db, id, connectionId)
}

func addOneConnectionToWorkflow(db *sqlx.DB, workflowId string, connectionId string) error {

	// TODO: check for existing record

	objectKey := fmt.Sprintf(`<Key><T>DPRProjectionConfigHasConnect</T><P>%s</P><P>%s</P></Key>`, workflowId, connectionId)

	t := DPRProjectionConfigHasConnect{
		UID_DPRProjectionConfig: workflowId,
		UID_DPRSystemConnection: connectionId,
		Specials:                oneim.Specials{XObjectKey: objectKey},
	}

	err := InsertDPRObject[DPRProjectionConfigHasConnect](db, &t)
	if err != nil {
		return err
	}
	fmt.Println(objectKey)

	return nil
}

var AddAllConnectionsToWorkflowCmd = createDPRCommand(
	"add-all-connections",
	"add all system connections associated with the project to a workflow",
	`Add all system connection in the given project to a synchronization workflow (DPRProjectionConfigHasConnect).`,
	[]string{"id"},
	addAllConnectionsToWorkflow,
)

func addAllConnectionsToWorkflow(c *cobra.Command, db *sqlx.DB) error {

	id, err := GetStructId_MustExist[DPRProjectionConfig](c, "id", db)
	if err != nil {
		return err
	}

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return err
	}

	connections, err := GetAllConnections(db, shellId)
	if err != nil {
		return err
	}

	for _, v := range connections {
		err = addOneConnectionToWorkflow(db, id, v.UID_DPRSystemConnection)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetWorkflowByName(db *sqlx.DB, shellId string, wfName string) (*DPRProjectionConfig, error) {
	if !IsValidIdOrName(shellId) {
		return nil, errors.New("invalid shell id")
	}
	if !IsValidIdOrName(wfName) {
		return nil, errors.New("invalid wf id")
	}

	wc := fmt.Sprintf(`UID_DPRShell='%s' and Name like '%s'`, shellId, wfName)
	ts, err := dbx.GetStructData[DPRProjectionConfig](db, "DPRProjectionConfig", wc)
	if err != nil {
		return nil, err
	}

	if len(ts) == 0 {
		return nil, nil
	} else if len(ts) > 1 {
		return nil, errors.New("workflow name is not unique")
	}

	return &ts[0], nil
}

func GetAllWorkflowConnections(db *sqlx.DB, wfId string) ([]DPRSystemConnection, error) {
	if !IsValidIdOrName(wfId) {
		return nil, errors.New("invalid wf id")
	}
	wc := fmt.Sprintf(`
		UID_DPRSystemConnection in (
			select UID_DPRSystemConnection from DPRProjectionConfigHasConnect where UID_DPRProjectionConfig ='%s'
		)
	`, wfId)
	return dbx.GetStructData[DPRSystemConnection](db, "DPRSystemConnection", wc)
}

func GetWFDefaultOneIMConnection(db *sqlx.DB, wfId string) (*DPRSystemConnection, error) {
	return getWFConnectionByName(db, wfId, "MainConnection")
}

func GetWFDefaultTargetSystemConnection(db *sqlx.DB, wfId string) (*DPRSystemConnection, error) {
	return getWFConnectionByName(db, wfId, "ConnectedSystemConnection")
}

func getWFConnectionByName(db *sqlx.DB, wfId string, name string) (*DPRSystemConnection, error) {
	ts, err := GetAllWorkflowConnections(db, wfId)
	if err != nil {
		return nil, err
	}

	// find first MainConnection
	for _, v := range ts {
		if *v.Name == name {
			return &v, nil
		}
	}

	return nil, nil
}

var UpdateWorkflowCmd = CreateUpdateCommand(
	"update an existing synchronization worklfow",
	`Update attributes of a sync project workflow (DPRProjectionConfig).`,
	updateWorkflow,
)

func updateWorkflow(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRProjectionConfig](c, db)
}

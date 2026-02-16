package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRProjectionConfigStep struct {
	oneim.Specials
	Displayable
	UID_DPRProjectionConfigStep  string
	UID_DPRProjectionConfig      string
	UID_QBMClrType               string
	UID_DPRSystemMap             string
	UID_DPRSystemObjectMatchSets string
	UID_LeftDPRSystemConnection  string
	UID_RightDPRSystemConnection string
	UID_LeftDPRProjectionQuota   *string
	UID_RightDPRProjectionQuota  *string
	Workflow                     string `mapstructure:",omitzero"`
	Map                          string `mapstructure:",omitzero"`
	MatchSets                    string `mapstructure:",omitzero"`
	LeftConnection               string `mapstructure:",omitzero"`
	RightConnection              string `mapstructure:",omitzero"`
	LeftQuota                    string `mapstructure:",omitzero"`
	RightQuota                   string `mapstructure:",omitzero"`
}

var WorkflowStepCmd = CreateBaseCommand(
	"workflow-step",
	"workflow step commands",
	`View and update synchronization steps in a workflow (DPRProjectionConfigStep).`,
	showWorkflowSteps,
)

var ShowWorkflowStepCmd = CreateShowCommand(
	"show synchronization workflow step",
	`View sync project workflow step (DPRProjectionConfgStep).`,
	[]string{"workflow-id"},
	showWorkflowSteps,
)

func showWorkflowSteps(c *cobra.Command, db *sqlx.DB) error {

	wfId, err := GetStructId_MustExist[DPRProjectionConfig](c, "workflow-id", db)
	if err != nil {
		return nil
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRProjectionConfig='%s'`, wfId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRProjectionConfigStep='%s'`, id)
	}

	return ShowDPRObjects[DPRProjectionConfigStep](db, wc, fillWorkflowStepData)
}

func fillWorkflowStepData(db *sqlx.DB, t *DPRProjectionConfigStep) error {

	t.Map, _ = GetForeignDisplay(db, "DPRSystemMap", t.UID_DPRSystemMap)
	t.MatchSets, _ = GetForeignDisplay(db, "DPRSystemObjectMatchSets", t.UID_DPRSystemObjectMatchSets)
	t.LeftConnection, _ = GetForeignDisplay(db, "DPRSystemConnection", t.UID_LeftDPRSystemConnection)
	t.RightConnection, _ = GetForeignDisplay(db, "DPRSystemConnection", t.UID_RightDPRSystemConnection)
	if t.UID_LeftDPRProjectionQuota != nil {
		t.LeftQuota, _ = GetForeignDisplay(db, "DPRProjectionQuota", *t.UID_LeftDPRProjectionQuota)
	}
	if t.UID_RightDPRProjectionQuota != nil {
		t.RightQuota, _ = GetForeignDisplay(db, "DPRProjectionQuota", *t.UID_RightDPRProjectionQuota)
	}

	return nil
}

var InsertWorkflowStepCmd = CreateInsertCommand(
	"create a new synchronization workflow step",
	`Create a new sync workflow step (DPRProjectionConfigStep).`,
	[]string{"workflow-id", "name", "map-id"},
	insertWorkflowStep,
)

func insertWorkflowStep(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRProjectionConfigStep](c, db, "VI.Projector.Projection.ProjectionStep", newWorkflowStep)
}

func newWorkflowStep(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRProjectionConfigStep, error) {

	wfId, err := GetStructId_MustExist[DPRProjectionConfig](c, "workflow-id", db)
	if err != nil {
		return nil, err
	}

	mapId, err := GetStructId_MustExist[DPRSystemMap](c, "map-id", db)
	if err != nil {
		return nil, err
	}

	// has user requested default connections?
	defaultConnections, err := c.Flags().GetBool("use-default-connections")
	if err != nil {
		return nil, err
	}
	leftConnectionId := ""
	rightConnectionId := ""
	if defaultConnections {
		leftConnection, err := GetWFDefaultOneIMConnection(db, wfId)
		if err != nil || leftConnection == nil {
			return nil, errors.New("unable to find default Identity Manager connection")
		}
		leftConnectionId = leftConnection.UID_DPRSystemConnection

		rightConnection, err := GetWFDefaultTargetSystemConnection(db, wfId)
		if err != nil || rightConnection == nil {
			return nil, errors.New("unable to find default target systtem connection")
		}
		rightConnectionId = rightConnection.UID_DPRSystemConnection
	} else {

		leftConnectionId, err = GetStructId_MustExist[DPRSystemConnection](c, "left-connection-id", db)
		if err != nil {
			return nil, err
		}
		rightConnectionId, err = GetStructId_MustExist[DPRSystemConnection](c, "right-connection-id", db)
		if err != nil {
			return nil, err
		}

	}

	// has user requested default match sets?
	matchSetsId := ""
	defaultMatchSets, err := c.Flags().GetBool("include-default-match-sets")
	if err != nil {
		return nil, err
	}
	if defaultMatchSets {
		matchSets, err := CreateNewMatchSetsWithDefaults(db, id)
		if err != nil {
			return nil, err
		}
		matchSetsId = matchSets.UID_DPRSystemObjectMatchSets
	} else {
		matchSetsId, err = GetStructId_MustExist[DPRSystemObjectMatchSets](c, "match-sets-id", db)
		if err != nil {
			return nil, errors.New("must include flag match-sets-id or include-default-match-sets")
		}
	}

	t := DPRProjectionConfigStep{
		UID_DPRProjectionConfigStep:  id,
		UID_QBMClrType:               clrId,
		UID_DPRProjectionConfig:      wfId,
		UID_DPRSystemMap:             mapId,
		UID_LeftDPRSystemConnection:  leftConnectionId,
		UID_RightDPRSystemConnection: rightConnectionId,
		UID_DPRSystemObjectMatchSets: matchSetsId,
		Specials:                     oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
	}

	return &t, nil
}

var UpdateWorkflowStepCmd = CreateUpdateCommand(
	"update an existing synchronization workflow step",
	`Update attributes of a sync project workflow step (DPRProjectionConfigStep).`,
	updateWorkflowStep,
)

func updateWorkflowStep(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRProjectionConfigStep](c, db)
}

var AddSchemaMethodCmd = createDPRCommand(
	"add-schema-method",
	"add a schema method to a workflow step",
	"Add a schema method (insert, update, etc.) to a workflow step. (DPRSystemObjectMatchSetMA)",
	[]string{"id", "side", "method", "match-set"},
	addSchemaMethod,
)

func addSchemaMethod(c *cobra.Command, db *sqlx.DB) error {

	stepId, err := GetStructId_MustExist[DPRProjectionConfigStep](c, "id", db)
	step, err := dbx.GetStructSingleton[DPRProjectionConfigStep](db, stepId)
	if err != nil {
		return err
	}

	// fetch map for this step
	systemMap, err := dbx.GetStructSingleton[DPRSystemMap](db, step.UID_DPRSystemMap)
	if err != nil {
		return err
	}

	// which side of map?
	side, _ := c.Flags().GetString("side")
	if len(side) == 0 {
		return errors.New("side (left or right) is required")
	}
	// get schema class associated with selected side
	schemaClassId := ""
	if strings.ToUpper(side) == "LEFT" {
		schemaClassId = systemMap.UID_LeftDPRSchemaClass
	} else if strings.ToUpper(side) == "RIGHT" {
		schemaClassId = systemMap.UID_RightDPRSchemaClass
	} else {
		return errors.New("invalid schema side: " + side)
	}
	if len(schemaClassId) == 0 {
		return errors.New("failed to lookup schema class")
	}
	// fetch schema type of given class
	stWC := fmt.Sprintf(`
			UID_DPRSchemaType = (
				select UID_DPRSchemaType from DPRSchemaClass where UID_DPRSchemaClass='%s'
			)`, schemaClassId)
	st, err := dbx.GetStructSingletonByWC[DPRSchemaClass](db, stWC)
	if err != nil {
		return nil
	}

	methodName, _ := c.Flags().GetString("method")
	if len(methodName) == 0 {
		return errors.New("method is required")
	}
	if !IsValidIdOrName(methodName) {
		return errors.New("invalid method: " + methodName)
	}
	// fetch schema method by name / class
	smWC := fmt.Sprintf(`UID_DPRSchemaType='%s' and Name='%s'`, st.UID_DPRSchemaType, methodName)
	sm, err := dbx.GetStructSingletonByWC[DPRSchemaMethod](db, smWC)
	if err != nil {
		return nil
	}

	// fetch match set collection for this wf step
	matchSetCol, err := dbx.GetStructSingleton[DPRSystemObjectMatchSets](db, step.UID_DPRSystemObjectMatchSets)
	if err != nil {
		return err
	}
	// which match set?
	matchSetName, _ := c.Flags().GetString("match-set")
	if len(matchSetName) == 0 {
		return errors.New("match set is required")
	}
	matchSet, err := GetMatchSetByName(db, matchSetCol.UID_DPRSystemObjectMatchSets, matchSetName)

	// insert new match set assignment
	msaName := methodName + step.UID_DPRProjectionConfigStep
	msa, _, err := InsertNewMatchSetMethodAssignment(db, msaName,
		matchSet.UID_DPRSystemObjectMatchSet, sm.UID_DPRSchemaMethod,
		side, fmt.Sprintf("To%s", side))

	fmt.Println(msa.UID_DPRSystemObjectMatchSetMA)
	return err
}

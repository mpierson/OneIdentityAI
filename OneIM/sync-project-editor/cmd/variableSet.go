package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

type DPRSystemVariableSet struct {
	oneim.Specials
	Displayable
	UID_DPRSystemVariableSet       string
	UID_DPRShell                   string
	UID_QBMClrType                 string
	UID_DPRSystemVariableSetParent *string
	Parent                         string   `mapstructure:",omitzero"`
	Variables                      []string `mapstructure:",omitzero"`
}

var VariableSetCmd = CreateBaseCommand(
	"variable-set",
	"sync project variable set commands",
	`View and update synchronization variable set (DPRSystemVariableSet).`,
	showVariableSets,
)

var ShowVariableSetCmd = CreateShowCommand(
	"show details of one synchronization variable set",
	`View sync project variable set (DPRSystemVariableSet).`,
	[]string{"shell"},
	showVariableSets,
)

func showVariableSets(c *cobra.Command, db *sqlx.DB) error {

	shellId, _ := c.Flags().GetString("shell")
	if len(shellId) == 0 {
		return errors.New("shell id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRShell='%s'`, shellId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSystemVariableSet='%s'`, id)
	}

	return ShowDPRObjects[DPRSystemVariableSet](db, wc, fillVariableSetData)
}

func fillVariableSetData(db *sqlx.DB, t *DPRSystemVariableSet) error {

	if t.UID_DPRSystemVariableSetParent != nil {
		t.Parent, _ = GetForeignDisplay(db, "DPRSystemVariableSet", *t.UID_DPRSystemVariableSetParent)
	}
	t.Variables, _ = GetChildDisplays(db, "DPRSystemVariable", "UID_DPRSystemVariableSet", t.UID_DPRSystemVariableSet)

	return nil
}

var InsertVariableSetCmd = CreateInsertCommand(
	"create a new synchronization variable set",
	`Create a new sync Connection (DPRSystemVariableSet).`,
	[]string{"shell", "name"},
	insertVariableSet,
)

func insertVariableSet(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemVariableSet](c, db, "VI.Projector.Variables.SystemVariableSet", newVariableSet)
}

func newVariableSet(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemVariableSet, error) {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	parentId, _ := c.Flags().GetString("parent-id")

	t := DPRSystemVariableSet{
		UID_DPRSystemVariableSet:       id,
		UID_DPRShell:                   shellId,
		UID_QBMClrType:                 clrId,
		UID_DPRSystemVariableSetParent: &parentId,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
	}

	return &t, nil
}

var UpdateVariableSetCmd = CreateUpdateCommand(
	"update an existing variable set",
	`Update attributes of a sync project variable set (DPRSystemVariableSet).`,
	updateVariableSet,
)

func updateVariableSet(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSystemVariableSet](c, db)
}

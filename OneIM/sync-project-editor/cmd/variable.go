package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

type DPRSystemVariable struct {
	oneim.Specials
	Displayable
	UID_DPRSystemVariable    string
	UID_DPRSystemVariableSet string
	UID_QBMClrType           string
	GenerateValueScript      *string
	ScriptLanguage           *string
	IsSecret                 bool
	IsSystemVariable         bool
	Value                    *string
}

var VariableCmd = CreateBaseCommand(
	"variable",
	"sync project variable commands",
	`View and update synchronization variable (DPRSystemVariable).`,
	showVariables,
)

var ShowVariableCmd = CreateShowCommand(
	"show details of one synchronization start info",
	`View sync project start info (DPRProjectionStartInfo).`,
	[]string{"variable-set"},
	showVariables,
)

func showVariables(c *cobra.Command, db *sqlx.DB) error {

	vsId, _ := c.Flags().GetString("variable-set")
	if len(vsId) == 0 {
		return errors.New("variable set id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRSystemVariableSet='%s'`, vsId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSystemVariable='%s'`, id)
	}

	return ShowDPRObjects[DPRSystemVariable](db, wc, fillVariableData)
}

func fillVariableData(db *sqlx.DB, t *DPRSystemVariable) error {

	return nil
}

var InsertVariableCmd = CreateInsertCommand(
	"create a new synchronization variable",
	`Create a new synchronization variable (DPRSystemVariable) and return the UID_SystemVariable of the new variable.`,
	[]string{"variable-set", "name"},
	insertVariable,
)

func insertVariable(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemVariable](c, db, "VI.Projector.Variables.SystemVariable", newVariable)
}

func newVariable(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemVariable, error) {

	vsId, err := GetStructId_MustExist[DPRSystemVariableSet](c, "variable-set", db)
	if err != nil {
		return nil, err
	}

	val, _ := c.Flags().GetString("value")
	isSecret, _ := c.Flags().GetBool("secret")

	t := DPRSystemVariable{
		UID_DPRSystemVariable:    id,
		UID_QBMClrType:           clrId,
		UID_DPRSystemVariableSet: vsId,
		Value:                    &val,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
		IsSecret: isSecret,
	}

	return &t, nil
}

var UpdateVariableCmd = CreateUpdateCommand(
	"update a variable",
	`Update attributes of a sync project variable (DPRSystemVariable).`,
	updateVariable,
)

func updateVariable(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSystemVariable](c, db)
}

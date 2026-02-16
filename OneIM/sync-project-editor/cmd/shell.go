package cmd

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRShell struct {
	oneim.Specials
	Displayable
	UID_DPRShell                string
	UID_QBMClrType              string
	UID_DPRSystemVariableSetDef *string
	OriginInfo                  *string
	ScriptLanguage              string
	IsAutomaticallyManaged      bool
	IsFinalized                 int
	EditedBy                    *string
	EditedSince                 *string
	SystemVariableSetDef        string   `mapstructure:",omitempty"`
	StartInfos                  []string `mapstructure:",omitempty"`
	Workflows                   []string `mapstructure:",omitempty"`
	VariableSets                []string `mapstructure:",omitempty"`
	Schemas                     []string `mapstructure:",omitempty"`
	Connections                 []string `mapstructure:",omitempty"`
	Maps                        []string `mapstructure:",omitempty"`
}

var ShellCmd = CreateBaseCommand(
	"shell",
	"sync project commands",
	`View and update sync project records (DPRShell).`,
	showShells,
)

var ShowShellCmd = CreateShowCommand(
	"show details of one or more shells",
	`View sync project records (DPRShell).`,
	nil,
	showShells,
)

func showShells(c *cobra.Command, db *sqlx.DB) error {

	shellId, _ := c.Flags().GetString("id")

	wc := "1=1"
	if len(shellId) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRShell='%s'`, shellId)
	}

	return ShowDPRObjects[DPRShell](db, wc, fillShellData)

}

func fillShellData(db *sqlx.DB, shell *DPRShell) error {
	if shell.UID_DPRSystemVariableSetDef != nil {
		shell.SystemVariableSetDef, _ = GetForeignDisplay(
			db, "DPRSystemVariableSet", *shell.UID_DPRSystemVariableSetDef)
	}
	shell.StartInfos, _ = GetChildDisplays(db, "DPRProjectionStartInfo", "UID_DPRShell", shell.UID_DPRShell)
	shell.Workflows, _ = GetChildDisplays(db, "DPRProjectionConfig", "UID_DPRShell", shell.UID_DPRShell)
	shell.VariableSets, _ = GetChildDisplays(db, "DPRSystemVariableSet", "UID_DPRShell", shell.UID_DPRShell)
	shell.Schemas, _ = GetChildDisplays(db, "DPRSchema", "UID_DPRShell", shell.UID_DPRShell)
	shell.Connections, _ = GetChildDisplays(db, "DPRSystemConnection", "UID_DPRShell", shell.UID_DPRShell)
	shell.Maps, _ = GetChildDisplays(db, "DPRSystemMap", "UID_DPRShell", shell.UID_DPRShell)

	return nil
}

var InsertShellCmd = CreateInsertCommand(
	"create a new synchronization project",
	`Create a new sync project (DPRShell).`,
	[]string{"name"},
	insertShell,
)

func insertShell(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRShell](c, db, "VI.Projector.ProjectorShell", newShell)
}

func newShell(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRShell, error) {

	origin := "Created by sped command " + time.Now().String()

	t := DPRShell{
		UID_DPRShell:   id,
		UID_QBMClrType: clrId,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		Displayable: Displayable{
			Name:        &id,
			DisplayName: &name,
		},
		OriginInfo:     &origin,
		ScriptLanguage: "VisualBasicNet",
	}

	return &t, nil
}

var UpdateShellCmd = CreateUpdateCommand(
	"update an existing synchronization project",
	`Update attributes of a sync project (DPRShell).`,
	updateShell,
)

func updateShell(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRShell](c, db)
}

func GetDefaultVariableSet(db *sqlx.DB, shellId string) (*DPRSystemVariableSet, error) {

	t, err := dbx.GetStructSingleton[DPRShell](db, shellId)
	if err != nil {
		return nil, err
	}

	if len(*t.UID_DPRSystemVariableSetDef) > 0 {
		vs, err := dbx.GetStructSingleton[DPRSystemVariableSet](db, *t.UID_DPRSystemVariableSetDef)
		if err != nil {
			return nil, err
		}

		return &vs, nil
	}

	return nil, nil
}

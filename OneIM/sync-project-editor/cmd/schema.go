package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRSchema struct {
	oneim.Specials
	Displayable
	UID_DPRSchema      string
	UID_DPRShell       string
	UID_QBMClrType     string
	FunctionalLevel    int
	IsLocked           bool
	IsPartial          bool
	NameFormat         *string
	SystemCapabilities *string
	SystemDisplay      *string
	SystemId           string
	SystemSubType      *string
	SystemType         *string
	SystemVersion      *string
	SchemaTypes        []string `mapstructure:",omitzero"`
}

var SchemaCmd = CreateBaseCommand(
	"schema",
	"sync project schema commands",
	`View and update synchronization schema (DPRSchema).`,
	showSchemas,
)

var ShowSchemaCmd = CreateShowCommand(
	"show details of one synchronization schema",
	`View sync project schema (DPRSchema).`,
	[]string{"shell"},
	showSchemas,
)

func showSchemas(c *cobra.Command, db *sqlx.DB) error {

	shellId, _ := c.Flags().GetString("shell")
	if len(shellId) == 0 {
		return errors.New("shell id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRShell='%s'`, shellId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSchema='%s'`, id)
	}

	return ShowDPRObjects[DPRSchema](db, wc, fillSchemaData)
}

func fillSchemaData(db *sqlx.DB, t *DPRSchema) error {

	t.SchemaTypes, _ = GetChildDisplays(db, "DPRSchemaType", "UID_DPRSchema", t.UID_DPRSchema)

	return nil
}

var InsertSchemaCmd = CreateInsertCommand(
	"create a new synchronization schema",
	`Create a new sync schema (DPRSchema).`,
	[]string{"shell", "name", "system-id", "clr-name"},
	insertSchema,
)

func insertSchema(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSchema](c, db, clrId, newSchema)
}

func newSchema(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchema, error) {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	sysId, _ := c.Flags().GetString("system-id")

	return newSchemaInternal(db, id, objectKey, name, shellId, sysId, clrId)
}

func newSchemaInternal(
	db *sqlx.DB,
	id string, objectKey string, name string,
	shellId string, sysId string,
	clrId string,
) (*DPRSchema, error) {

	t := DPRSchema{
		UID_DPRSchema:  id,
		UID_DPRShell:   shellId,
		UID_QBMClrType: clrId,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
		SystemId: sysId,
	}

	return &t, nil
}

var InsertOneIMSchemaCmd = CreateBaseCommand(
	"insert-oneim-schema",
	"create a new synchronization schema for the OneIM database",
	`Create a new sync schema for Identity Manager (DPRSchema).`,
	insertOneIMSchema,
)

func insertOneIMSchema(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSchema](c, db, "VI.Projector.Database.DatabaseSchema", newOneIMSchema)
}
func newOneIMSchema(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchema, error) {

	if len(name) == 0 {
		return nil, errors.New("Missing schema name")
	}

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	oneIMDB, err := dbx.GetStructSingletonByWC[DialogDatabase](db, "IsMainDatabase=1")
	if err != nil {
		return nil, err
	}

	sysId := fmt.Sprintf(`FTP#%s`, oneIMDB.UID_Database)
	systemType := "OneIM"
	systemVersion := fmt.Sprintf(`%s.0.0`, *oneIMDB.EditionVersion)
	systemDisplay := "One Identity Manager"
	nameFormat := "Identifier"
	systemCapabilities := "SupportsRevisions"

	t := DPRSchema{
		UID_DPRSchema:  id,
		UID_DPRShell:   shellId,
		UID_QBMClrType: clrId,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
		SystemId:           sysId,
		SystemType:         &systemType,
		SystemVersion:      &systemVersion,
		SystemDisplay:      &systemDisplay,
		NameFormat:         &nameFormat,
		SystemCapabilities: &systemCapabilities,
		FunctionalLevel:    2,
	}

	return &t, nil
}

var UpdateSchemaCmd = CreateUpdateCommand(
	"update an existing synchronization schema",
	`Update attributes of a sync project schema (DPRSchema).`,
	updateSchema,
)

func updateSchema(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSchema](c, db)
}

func getSchema(db *sqlx.DB, shellId string, wc string) (DPRSchema, error) {

	var t DPRSchema
	ts, err := dbx.GetStructData[DPRSchema](db, "DPRSchema", wc)
	if err != nil {
		return t, err
	}
	if len(ts) == 0 {
		return t, errors.New("Schema not found in shell: " + shellId)
	} else if len(ts) > 1 {
		return t, errors.New(fmt.Sprintf("Too many schemas in shell %s", shellId))
	}

	return ts[0], nil
}

func GetOneIMSchema(db *sqlx.DB, shellId string) (DPRSchema, error) {
	wc := fmt.Sprintf(`UID_DPRShell='%s' and SystemId like 'FTP#%%'`, shellId)
	return getSchema(db, shellId, wc)
}

func GetTargetSystemSchema(db *sqlx.DB, shellId string) (DPRSchema, error) {
	wc := fmt.Sprintf(`UID_DPRShell='%s' and not SystemId like 'FTP#%%'`, shellId)
	return getSchema(db, shellId, wc)
}

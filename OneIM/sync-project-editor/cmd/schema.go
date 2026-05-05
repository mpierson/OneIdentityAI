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
	`Create a new synchronization schema (DPRSchema) and return the UID_DPRSchema of the new schema.`,
	[]string{"shell", "name", "system-id", "clr-name"},
	insertSchema,
)

func insertSchema(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSchema](c, db, clrId, newSchemaCmd)
}

func newSchemaCmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchema, error) {

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	sysId, _ := c.Flags().GetString("system-id")

	return newSchema(id, objectKey, name, clrId, shellId, sysId, "", "", "", "", 0)
}

var InsertOneIMSchemaCmd = CreateBaseCommand(
	"insert-oneim-schema",
	"create a new synchronization schema for the OneIM database",
	`Create a new synchronization schema for Identity Manager (DPRSchema) and return the UID_Schema of the new schema.`,
	insertOneIMSchema,
)

func insertOneIMSchema(c *cobra.Command, db *sqlx.DB) error {

	var schemaId string
	newFn := func(c *cobra.Command, db *sqlx.DB, id string, objectKey string, name string, clrId string) (*DPRSchema, error) {
		schemaId = id
		t, schemaErr := newOneIMSchema(c, db, id, objectKey, name, clrId)
		if schemaErr != nil {
			return t, schemaErr
		}
		t.IsPartial = true
		return t, nil
	}
	err := ExecInsertCommand[DPRSchema](c, db, "VI.Projector.Database.DatabaseSchema", newFn)
	if err != nil {
		return err
	}

	// add standard OneIM revision table
	updater := func(db *sqlx.DB, schemaType *DPRSchemaType) error {
		schemaType.IsLocked = true
		schemaType.IsReadOnly = true
		schemaType.ShrinkLock = 1
		metaData := "ModUidDef=1"
		schemaType.MetaData = &metaData
		return nil
	}
	typeId, err := AddTypeToSchema(db, schemaId, "QBMVTableRevision", "VI.Projector.Database.DatabaseSchemaType", updater)
	if err != nil {
		return err
	}
	err = AddDBPropertiesToSchemaType(db, typeId)
	if err != nil {
		return err
	}
	_, err = AddDefaultClassToSchemaType(db, typeId)
	if err != nil {
		return err
	}

	return nil
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
	sysVersion := fmt.Sprintf(`%s.0.0`, *oneIMDB.EditionVersion)
	return newSchema(id, objectKey, name, clrId,
		shellId, sysId, "OneIM", sysVersion, "One Identity Manager", "SupportsRevisions", 2)
}

func newSchema(
	id string, objectKey string, name string, clrId string,
	shellId string, systemId string,
	systemType string, systemVersion string, systemDisplay string,
	systemCapabilities string, functionalLevel int) (*DPRSchema, error) {

	t := DPRSchema{
		UID_DPRSchema:  id,
		UID_DPRShell:   shellId,
		UID_QBMClrType: clrId,
		Specials:       oneim.NewSpecials(objectKey, "sped"),
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
			NameFormat:  &NAME_FORMAT_Identifier,
		},
		SystemId:           systemId,
		SystemType:         &systemType,
		SystemVersion:      &systemVersion,
		SystemDisplay:      &systemDisplay,
		SystemCapabilities: &systemCapabilities,
		FunctionalLevel:    functionalLevel,
	}

	return &t, nil
}

var InsertCustomSchemaCmd = CreateBaseCommand(
	"insert-target-schema",
	"create a new synchronization schema for a target system",
	`Create a new synchronization schema for a custom target system (DPRSchema) and return the UID_Schema of the new schema.`,
	insertCustomSchema,
)

func insertCustomSchema(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSchema](c, db, "VI.Projector.Powershell.PoshSchema", newCustomSchema)
}
func newCustomSchema(
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

	sysId, _ := c.Flags().GetString("system-id")
	if !IsValidIdOrName(sysId) {
		return nil, errors.New("Invalid system identifier: " + sysId)
	}

	t, err := newSchema(id, objectKey, name, clrId,
		shellId, sysId, "Posh", "345", "PoshNet40", "NativeFilterNotSupported", 0)
	if err != nil {
		return nil, err
	}

	return t, nil
}

var UpdateSchemaCmd = CreateUpdateCommand(
	"update an existing synchronization schema",
	`Update attributes of a sync project schema (DPRSchema).`,
	updateSchema,
)

func updateSchema(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSchema](c, db)
}

func getSchema(db *sqlx.DB, shellId string, wc string) (*DPRSchema, error) {

	var t DPRSchema
	ts, err := dbx.GetStructData[DPRSchema](db, "DPRSchema", wc)
	if err != nil {
		return &t, err
	}
	if len(ts) == 0 && len(shellId) > 0 {
		return &t, errors.New("Schema not found in shell: " + shellId)
	} else if len(ts) > 1 && len(shellId) > 0 {
		return &t, errors.New(fmt.Sprintf("Too many matching schemas in shell %s", shellId))
	}

	return &ts[0], nil
}

func GetSchema(db *sqlx.DB, schemaId string) (*DPRSchema, error) {
	if !IsValidIdOrName(schemaId) {
		return nil, errors.New("Invalid schema identifier: " + schemaId)
	}
	wc := fmt.Sprintf(`UID_DPRSchema='%s'`, schemaId)
	return getSchema(db, "", wc)
}

func GetOneIMSchema(db *sqlx.DB, shellId string) (*DPRSchema, error) {
	wc := fmt.Sprintf(`UID_DPRShell='%s' and SystemId like 'FTP#%%'`, shellId)
	return getSchema(db, shellId, wc)
}

func GetTargetSystemSchema(db *sqlx.DB, shellId string) (*DPRSchema, error) {
	wc := fmt.Sprintf(`UID_DPRShell='%s' and not SystemId like 'FTP#%%'`, shellId)
	return getSchema(db, shellId, wc)
}

// choose OneIM CLR or target system CLR, based on given schema type
func GetCLRForTarget(db *sqlx.DB, schemaId string, oneIMCLR string, targetCLR string) (string, error) {

	schema, err := GetSchema(db, schemaId)
	if err != nil {
		return "", err
	}

	if schema.SystemType != nil && "OneIM" == *(schema.SystemType) {
		return oneIMCLR, nil
	}
	return targetCLR, nil
}

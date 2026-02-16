package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRSchemaClass struct {
	oneim.Specials
	Displayable
	UID_DPRSchemaClass string
	UID_DPRSchemaType  string
	UID_QBMClrType     string
	Filter             *string
	IsLocked           bool
	IsObsolete         bool
	LeftMaps           []string `mapstructure:",omitempty"`
	RightMaps          []string `mapstructure:",omitempty"`
}

var SchemaClassCmd = CreateBaseCommand(
	"schema-class",
	"sync project schema class commands",
	`View and update synchronization schema class (DPRSchemaClass).`,
	showSchemaClasss,
)

var ShowSchemaClassCmd = CreateShowCommand(
	"show details of one or more synchronization schema classs",
	`View sync project schema class (DPRSchemaClass).`,
	[]string{"schema-type-id"},
	showSchemaClasss,
)

func showSchemaClasss(c *cobra.Command, db *sqlx.DB) error {

	schemaTypeId, _ := c.Flags().GetString("schema-type-id")
	if len(schemaTypeId) == 0 {
		return errors.New("schema type id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRSchemaType='%s'`, schemaTypeId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSchemaClass='%s'`, id)
	}

	return ShowDPRObjects[DPRSchemaClass](db, wc, fillSchemaClassData)
}

func fillSchemaClassData(db *sqlx.DB, t *DPRSchemaClass) error {

	t.LeftMaps, _ = GetChildDisplays(db, "DPRSystemMap", "UID_LeftDPRSchemaClass", t.UID_DPRSchemaClass)
	t.RightMaps, _ = GetChildDisplays(db, "DPRSystemMap", "UID_RightDPRSchemaClass", t.UID_DPRSchemaClass)

	return nil
}

var InsertSchemaClassCmd = CreateInsertCommand(
	"create a new synchronization schema class",
	`Create a new sync schema class (DPRSchemaClass).`,
	[]string{"schema-type-id", "name", "clr-name"},
	insertSchemaClass,
)

func insertSchemaClass(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSchemaClass](c, db, clrId, newSchemaClass_cmd)
}

func newSchemaClass_cmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchemaClass, error) {

	schemaTypeId, err := GetStructId_MustExist[DPRSchemaType](c, "id", db)
	if err != nil {
		return nil, err
	}

	return newSchemaClass(db, schemaTypeId, id, objectKey, name, name, clrId)
}

func newSchemaClass(
	db *sqlx.DB,
	schemaTypeId string,
	id string, objectKey string,
	name string, displayName string,
	clrId string,
) (*DPRSchemaClass, error) {

	t := DPRSchemaClass{
		UID_DPRSchemaClass: id,
		UID_QBMClrType:     clrId,
		UID_DPRSchemaType:  schemaTypeId,
		Specials:           oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &displayName,
		},
	}

	return &t, nil
}

var UpdateSchemaClassCmd = CreateUpdateCommand(
	"update an existing synchronization schema class",
	`Update attributes of a sync project schema class (DPRSchemaClass).`,
	updateSchemaClass,
)

func updateSchemaClass(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSchemaClass](c, db)
}

func GetSchemaClassByName(db *sqlx.DB, shellId string, schemaId string, name string) (DPRSchemaClass, error) {

	var t DPRSchemaClass

	if !IsValidIdOrName(shellId) || !IsValidIdOrName(name) {
		return t, errors.New(fmt.Sprintf("invalid input: %s %s", shellId, name))
	}

	wc := fmt.Sprintf(`
		UID_DPRSchemaType in (
			select UID_DPRSchemaType from DPRSchemaType 
				where UID_DPRSchema in (
					select UID_DPRSchema from DPRSchema 
						where UID_DPRShell='%s' and UID_DPRSchema='%s'
			)
		)
		and
		Name = '%s'`,
		shellId, schemaId, name)
	t1, err := dbx.GetStructSingletonByWC[DPRSchemaClass](db, wc)
	if err != nil {
		return t, err
	}

	return t1, nil
}

package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRSchemaProperty struct {
	oneim.Specials
	Displayable
	UID_DPRSchemaProperty string
	UID_DPRSchemaType     string
	UID_QBMClrType        string
	DataType              *string
	SchemaType            string `mapstructure:",omitzero"`
}

var SchemaPropertyCmd = CreateBaseCommand(
	"schema-property",
	"schema property commands",
	`View and update synchronization schema properties (DPRSchemaProperty).`,
	showSchemaPropertys,
)

var ShowSchemaPropertyCmd = CreateShowCommand(
	"show synchronization schema property",
	`View sync project schema property (DPRSchemaProperty).`,
	[]string{"schema-type-id"},
	showSchemaPropertys,
)

func showSchemaPropertys(c *cobra.Command, db *sqlx.DB) error {

	schemaTypeId, _ := c.Flags().GetString("schema-type-id")
	if len(schemaTypeId) == 0 {
		return errors.New("schema type id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRSchemaType='%s'`, schemaTypeId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSchemaProperty='%s'`, id)
	}

	return ShowDPRObjects[DPRSchemaProperty](db, wc, fillSchemaPropertyData)
}

func fillSchemaPropertyData(db *sqlx.DB, t *DPRSchemaProperty) error {

	t.SchemaType, _ = GetForeignDisplay(db, "DPRSchemaType", t.UID_DPRSchemaType)

	return nil
}

var InsertSchemaPropertyCmd = CreateInsertCommand(
	"create a new synchronization schema property",
	`Create a new sync schema property (DPRSchemaProperty).`,
	[]string{"schema-type-id", "name", "clr-name", "data-type"},
	insertSchemaProperty,
)

func insertSchemaProperty(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSchemaProperty](c, db, clrId, newSchemaProperty_cmd)
}

func newSchemaProperty_cmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchemaProperty, error) {

	schemaTypeId, err := GetStructId_MustExist[DPRSchemaType](c, "schema-type-id", db)
	if err != nil {
		return nil, err
	}

	dataType, err := c.Flags().GetString("data-type")
	if len(dataType) == 0 {
		return nil, errors.New("Missing data type")
	}

	return newSchemaProperty(db, schemaTypeId, id, objectKey, name, dataType, clrId)
}

func newSchemaProperty(
	db *sqlx.DB,
	schemaTypeId string,
	id string, objectKey string, name string,
	dataType string,
	clrId string,
) (*DPRSchemaProperty, error) {

	t := DPRSchemaProperty{
		UID_DPRSchemaProperty: id,
		UID_QBMClrType:        clrId,
		UID_DPRSchemaType:     schemaTypeId,
		Specials:              oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
		DataType: &dataType,
	}

	return &t, nil
}

var UpdateSchemaPropertyCmd = CreateUpdateCommand(
	"update an existing synchronization schema property",
	`Update attributes of a sync project schema property (DPRSchemaProperty).`,
	updateSchemaProperty,
)

func updateSchemaProperty(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSchemaProperty](c, db)
}

func GetAllSchemaProperties(db *sqlx.DB, schemaTypeId string) ([]DPRSchemaProperty, error) {

	if !IsValidIdOrName(schemaTypeId) {
		return nil, errors.New("invalid schema type id: " + schemaTypeId)
	}

	return dbx.GetStructData[DPRSchemaProperty](db, "DPRSchemaProperty",
		fmt.Sprintf(`UID_DPRSchemaType='%s'`, schemaTypeId))
}

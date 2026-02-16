package cmd

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

type DPRSchemaMethod struct {
	oneim.Specials
	Displayable
	UID_DPRSchemaMethod        string
	UID_DPRSchemaType          string
	UID_QBMClrType             string
	AcceptPartialLoadedObjects bool
	IsLocked                   bool
	IsNotCapableForImport      bool
	IsObsolete                 bool
	IsReadOnly                 bool
	MethodType                 *string
	SchemaType                 string `mapstructure:",omitzero"`
}

var ValidSchemaMethods = []string{"Insert", "Update", "Delete", "MarkAsOutstanding", "UnmarkAsOutstanding"}

var SchemaMethodCmd = CreateBaseCommand(
	"schema-method",
	"system schema method commands",
	`View and update synchronization method (DPRSchemaMethod).`,
	showSchemaMethods,
)

var ShowSchemaMethodCmd = CreateShowCommand(
	"show synchronization schema methods",
	`View sync project schema methods (DPRSchemaMethod).`,
	[]string{"schema-type-id"},
	showSchemaMethods,
)

func showSchemaMethods(c *cobra.Command, db *sqlx.DB) error {

	schemaTypeId, _ := c.Flags().GetString("schema-type-id")
	if len(schemaTypeId) == 0 {
		return errors.New("schema type required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRSchemaType='%s'`, schemaTypeId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSchemaMethod='%s'`, id)
	}

	return ShowDPRObjects[DPRSchemaMethod](db, wc, fillSchemaMethodData)
}

func fillSchemaMethodData(db *sqlx.DB, t *DPRSchemaMethod) error {

	t.SchemaType, _ = GetForeignDisplay(db, "DPRSchemaType", t.UID_DPRSchemaType)

	return nil
}

var InsertSchemaMethodCmd = CreateInsertCommand(
	"create a new synchronization schema method",
	`Create a new sync schema method (DPRSchemaMethod).`,
	[]string{"schema-type-id", "name", "clr-name"},
	insertSchemaMethod,
)

func insertSchemaMethod(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSchemaMethod](c, db, clrId, newSchemaMethodCmd)
}

func newSchemaMethodCmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchemaMethod, error) {

	schemaTypeId, err := GetStructId_MustExist[DPRSchemaType](c, "schema-type-id", db)
	if err != nil {
		return nil, err
	}
	return newSchemaMethod(db, id, objectKey, name, schemaTypeId, clrId)
}

func newSchemaMethod(
	db *sqlx.DB,
	id string, objectKey string, name string, schemaTypeId string,
	clrId string,
) (*DPRSchemaMethod, error) {

	if !slices.Contains(ValidSchemaMethods, name) {
		return nil, errors.New(fmt.Sprintf(`Invalid method name '%s'. Expect one of %v.`, name, ValidSchemaMethods))
	}

	t := DPRSchemaMethod{
		UID_DPRSchemaMethod: id,
		UID_QBMClrType:      clrId,
		UID_DPRSchemaType:   schemaTypeId,
		Specials:            oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
	}

	return &t, nil
}

var UpdateSchemaMethodCmd = CreateUpdateCommand(
	"update an existing synchronization schema method",
	`Update attributes of a sync project schema method (DPRSchemaMethod).`,
	updateSchemaMethod,
)

func updateSchemaMethod(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSchemaMethod](c, db)
}

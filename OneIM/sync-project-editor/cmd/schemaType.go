package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRSchemaType struct {
	oneim.Specials
	Displayable
	UID_DPRSchemaType string
	UID_DPRSchema     string
	UID_QBMClrType    string
	IsAdditional      bool
	IsInActive        bool
	IsLocked          bool
	IsMNType          bool
	IsObsolete        bool
	IsReadOnly        bool
	IsVirtual         bool
	ShrinkLock        int
	DisplayPattern    *string
	MetaData          *string
	Classes           []string `mapstructure:",omitzero"`
	SchemaProperties  []string `mapstructure:",omitzero"`
	SchemaMethods     []string `mapstructure:",omitzero"`
}

var SchemaTypeCmd = CreateBaseCommand(
	"schema-type",
	"sync project schema type commands",
	`View and update synchronization schema type (DPRSchemaType).`,
	showSchemaTypes,
)

var ShowSchemaTypeCmd = CreateShowCommand(
	"show details of one or more synchronization schema types",
	`View sync project schema type (DPRSchemaType).`,
	[]string{"schema-id"},
	showSchemaTypes,
)

func showSchemaTypes(c *cobra.Command, db *sqlx.DB) error {

	schemaId, _ := c.Flags().GetString("schema-id")
	if len(schemaId) == 0 {
		return errors.New("schema id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRSchema='%s'`, schemaId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSchemaType='%s'`, id)
	}

	return ShowDPRObjects[DPRSchemaType](db, wc, fillSchemaTypeData)
}

func fillSchemaTypeData(db *sqlx.DB, t *DPRSchemaType) error {

	//t.Map, _ = GetForeignDisplay(db, "DPRSystemMap", t.UID_DPRSystemMap)
	t.Classes, _ = GetChildDisplays(db, "DPRSchemaClass", "UID_DPRSchemaType", t.UID_DPRSchemaType)
	t.SchemaProperties, _ = GetChildDisplays(db, "DPRSchemaProperty", "UID_DPRSchemaType", t.UID_DPRSchemaType)
	t.SchemaMethods, _ = GetChildDisplays(db, "DPRSchemaMethod", "UID_DPRSchemaType", t.UID_DPRSchemaType)

	return nil
}

var InsertSchemaTypeCmd = CreateInsertCommand(
	"create a new synchronization schema type",
	`Create a new synchronization schema type (DPRSchemaType) and return the UID_DPRSchemaType of the new type.`,
	[]string{"schema-id", "name", "clr-name"},
	insertSchemaType,
)

func insertSchemaType(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSchemaType](c, db, clrId, newSchemaType)
}

func newSchemaType(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchemaType, error) {

	schemaId, _ := c.Flags().GetString("schema-id")

	t := DPRSchemaType{
		UID_DPRSchemaType: id,
		UID_QBMClrType:    clrId,
		UID_DPRSchema:     schemaId,
		Specials:          oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
	}

	return &t, nil
}

var UpdateSchemaTypeCmd = CreateUpdateCommand(
	"update an existing synchronization schema type",
	`Update attributes of a sync project schema type (DPRSchemaType).`,
	updateSchemaType,
)

func updateSchemaType(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSchemaType](c, db)
}

var AddSchemaMethodsCmd = createDPRCommand(
	"add-methods",
	"add one or more methods to a schema type",
	`Add methods to schema type (DPRSchemaMethod).`,
	[]string{"id", "clr-name"},
	addMethodsToSchemaType,
)

func addMethodsToSchemaType(c *cobra.Command, db *sqlx.DB) error {

	id, err := GetStructId_MustExist[DPRSchemaType](c, "id", db)
	if err != nil {
		return err
	}

	clrName, err := c.Flags().GetString("clr-name")
	if err != nil {
		return err
	} else if len(clrName) == 0 {
		return errors.New("CLR name required")
	}

	clrId, err := GetClrId(db, clrName)
	if err != nil {
		return err
	}

	var methods = []string{}
	allFlag, _ := c.Flags().GetBool("all")
	if allFlag {
		methods = ValidSchemaMethods
	} else {
		methodList, err := c.Flags().GetString("methods")
		if err != nil {
			return err
		}
		methods = strings.Split(methodList, " ")
	}

	for _, m := range methods {
		err = addMethodToSchemaType(db, id, m, clrId)
		if err != nil {
			return err
		}
	}

	return nil
}
func addMethodToSchemaType(db *sqlx.DB, typeId string, method string, clrId string) error {

	id, err := dbx.GetNewId(db)
	if err != nil {
		return err
	}
	objectKey := oneim.MakeObjectKey("DPRSchemaMethod", id)

	t, err := newSchemaMethod(db, id, objectKey, method, typeId, clrId)
	err = InsertDPRObject[DPRSchemaMethod](db, t)
	if err != nil {
		return err
	}

	return nil
}

var AddOneIMSchemaPropertiesCmd = createDPRCommand(
	"add-oneim-properties",
	"add all columns of the corresponding OneIM table to a schema type",
	`Add OneIM columns to schema type (DPRSchemaProperty).`,
	[]string{"id"},
	addPropertiesToSchemaType,
)

func addPropertiesToSchemaType(c *cobra.Command, db *sqlx.DB) error {

	id, err := c.Flags().GetString("id")
	if err != nil {
		return err
	}

	schemaType, err := dbx.GetStructSingleton[DPRSchemaType](db, id)
	if err != nil {
		return err
	}

	// current props, to avoid dupes
	currentProps, err := GetAllSchemaProperties(db, id)
	if err != nil {
		return err
	}
	// extract slice of property names
	currentPropNames := make([]string, len(currentProps))
	for i, p := range currentProps {
		currentPropNames[i] = *p.Name
	}

	cols, err := GetTableColumns(db, *schemaType.Name)
	for _, c := range cols {

		if !slices.Contains(currentPropNames, c[0]) {
			err = addPropertyToSchemaType(db, id, c[0], c[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func addPropertyToSchemaType(db *sqlx.DB, typeId string, cName string, cType string) error {

	id, err := dbx.GetNewId(db)
	if err != nil {
		return err
	}
	objectKey := oneim.MakeObjectKey("DPRSchemaProperty", id)

	// OneIM DB IDs
	clrId, err := GetClrId(db, "VI.Projector.Database.DatabaseSchemaProperty")
	if err != nil {
		return err
	}

	t, err := newSchemaProperty(db, typeId, id, objectKey, cName, cType, clrId)
	err = InsertDPRObject[DPRSchemaProperty](db, t)
	if err != nil {
		return err
	}

	return nil
}

var AddDefaultSchemaClassCmd = createDPRCommand(
	"add-default-class",
	"add the default 'all' class to a schema type",
	`Add the default unfiltered 'all' class (DPRSchemaClass) to a schema type and return the UID_DPRSchemaClass of the new class.`,
	[]string{"id"},
	addDefaultClassToSchemaType,
)

func addDefaultClassToSchemaType(c *cobra.Command, db *sqlx.DB) error {

	schemaTypeId, err := c.Flags().GetString("id")
	if err != nil {
		return err
	}

	schemaType, err := dbx.GetStructSingleton[DPRSchemaType](db, schemaTypeId)
	if err != nil {
		return err
	}

	id, err := dbx.GetNewId(db)
	if err != nil {
		return err
	}
	objectKey := oneim.MakeObjectKey("DPRSchemaClass", id)

	name := fmt.Sprintf(`%s_Master`, *schemaType.Name)
	displayName := fmt.Sprintf(`%s (all)`, *schemaType.Name)

	// imply clr of class from schema type CLR
	schemaTypeClrName, err := GetClrName(db, schemaType.UID_QBMClrType)
	if err != nil {
		return err
	}
	clrName := ""
	if schemaTypeClrName == "VI.Projector.Database.DatabaseSchemaType" {
		clrName = "VI.Projector.Database.DatabaseSchemaClass"
	} else {
		clrName = "VI.Projector.Schema.GenericSchemaClass"
	}

	clrId, err := GetClrId(db, clrName)
	if err != nil {
		return err
	}

	t, err := newSchemaClass(db, schemaTypeId, id, objectKey, name, displayName, clrId)
	err = InsertDPRObject[DPRSchemaClass](db, t)
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}

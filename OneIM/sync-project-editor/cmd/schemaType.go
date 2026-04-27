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
	[]string{"schema-id", "name"},
	insertSchemaType,
)

func insertSchemaType(c *cobra.Command, db *sqlx.DB) error {
	clr, _ := c.Flags().GetString("clr-name")

	if len(clr) == 0 {
		schemaId, _ := c.Flags().GetString("schema-id")
		clr, _ = GetCLRForTarget(db, schemaId,
			"VI.Projector.Database.DatabaseSchemaType", "VI.Projector.Powershell.PoshSchemaType")
	}

	return ExecInsertCommand[DPRSchemaType](c, db, clr, newSchemaTypeCmd)
}

func newSchemaTypeCmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSchemaType, error) {
	schemaId, _ := c.Flags().GetString("schema-id")
	return newSchemaType(db, id, objectKey, name, clrId, schemaId)
}

func newSchemaType(
	db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
	schemaId string,
) (*DPRSchemaType, error) {

	t := DPRSchemaType{
		UID_DPRSchemaType: id,
		UID_QBMClrType:    clrId,
		UID_DPRSchema:     schemaId,
		Specials:          oneim.NewSpecials(objectKey, "sped"),
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
			NameFormat:  &NAME_FORMAT_Identifier,
		},
		IsLocked: true,
	}

	return &t, nil
}

func AddTypeToSchema(db *sqlx.DB, schemaId string, typeName string, clr string) (string, error) {

	newFn := func(db1 *sqlx.DB, id string, objectKey string, name string, clrId string) (*DPRSchemaType, error) {
		return newSchemaType(db1, id, objectKey, name, clrId, schemaId)
	}

	_, id, err := InsertNewDPRObject(db, typeName, clr, newFn)

	return id, err
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
	[]string{"id"},
	addMethodsToSchemaType,
)

func addMethodsToSchemaType(c *cobra.Command, db *sqlx.DB) error {

	id, err := GetStructId_MustExist[DPRSchemaType](c, "id", db)
	if err != nil {
		return err
	}
	clr, _ := c.Flags().GetString("clr-name")
	if len(clr) == 0 {
		schemaType, err := dbx.GetStructSingleton[DPRSchemaType](db, id)
		if err != nil {
			return err
		}

		clr, _ = GetCLRForTarget(db, schemaType.UID_DPRSchema,
			"VI.Projector.Database.DatabaseSchemaMethod", "VI.Projector.Powershell.Schema.PoshSchemaMethod")
	}
	if len(clr) == 0 {
		return errors.New("CLR name required")
	}

	clrId, err := GetClrId(db, clr)
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

	id, objectKey, err := NewDPRKeys[DPRSchemaMethod](db)

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
	addDBPropertiesToSchemaTypeCmd,
)

func addDBPropertiesToSchemaTypeCmd(c *cobra.Command, db *sqlx.DB) error {

	id, err := c.Flags().GetString("id")
	if err != nil {
		return err
	}
	return AddDBPropertiesToSchemaType(db, id)
}

func AddDBPropertiesToSchemaType(db *sqlx.DB, schemaTypeId string) error {

	schemaType, err := dbx.GetStructSingleton[DPRSchemaType](db, schemaTypeId)
	if err != nil {
		return err
	}

	// current props, to avoid dupes
	currentProps, err := GetAllSchemaProperties(db, schemaTypeId)
	if err != nil {
		return err
	}
	// extract slice of property names
	currentPropNames := make([]string, len(currentProps))
	for i, p := range currentProps {
		currentPropNames[i] = *p.Name
	}

	cols, err := GetTableColumns(db, *schemaType.Name)
	for _, col := range cols {

		if !slices.Contains(currentPropNames, col[0]) {
			err = addDBPropertyToSchemaType(db, schemaTypeId, col[0], col[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func addDBPropertyToSchemaType(db *sqlx.DB, typeId string, cName string, cType string) error {

	id, objectKey, err := NewDPRKeys[DPRSchemaProperty](db)
	if err != nil {
		return err
	}
	clrId, err := GetClrId(db, "VI.Projector.Database.DatabaseSchemaProperty")
	if err != nil {
		return err
	}

	t, err := newSchemaProperty(typeId, id, objectKey, cName, cType, clrId)
	err = InsertDPRObject[DPRSchemaProperty](db, t)
	if err != nil {
		return err
	}

	// TODO: update unique key, secret attrs

	return nil
}

var AddDefaultSchemaClassCmd = createDPRCommand(
	"add-default-class",
	"add the default 'all' class to a schema type",
	`Add the default unfiltered 'all' class (DPRSchemaClass) to a schema type and return the UID_DPRSchemaClass of the new class.`,
	[]string{"id"},
	addDefaultClassToSchemaTypeCmd,
)

func addDefaultClassToSchemaTypeCmd(c *cobra.Command, db *sqlx.DB) error {

	schemaTypeId, err := c.Flags().GetString("id")
	if err != nil {
		return err
	}
	t, err := AddDefaultClassToSchemaType(db, schemaTypeId)
	if err != nil {
		return err
	}

	fmt.Println(t.UID_DPRSchemaClass)
	return nil
}

func AddDefaultClassToSchemaType(db *sqlx.DB, schemaTypeId string) (*DPRSchemaClass, error) {

	schemaType, err := dbx.GetStructSingleton[DPRSchemaType](db, schemaTypeId)
	if err != nil {
		return nil, err
	}

	id, objectKey, err := NewDPRKeys[DPRSchemaClass](db)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf(`%s_Master`, *schemaType.Name)
	displayName := fmt.Sprintf(`%s (all)`, *schemaType.Name)

	// imply clr of class from schema type CLR
	schemaTypeClrName, err := GetClrName(db, schemaType.UID_QBMClrType)
	if err != nil {
		return nil, err
	}
	clrName := ""
	if schemaTypeClrName == "VI.Projector.Database.DatabaseSchemaType" {
		clrName = "VI.Projector.Database.DatabaseSchemaClass"
	} else {
		clrName = "VI.Projector.Schema.GenericSchemaClass"
	}

	clrId, err := GetClrId(db, clrName)
	if err != nil {
		return nil, err
	}

	t, err := newSchemaClass(db, schemaTypeId, id, objectKey, name, displayName, clrId)
	err = InsertDPRObject[DPRSchemaClass](db, t)
	if err != nil {
		return t, err
	}

	return t, nil
}

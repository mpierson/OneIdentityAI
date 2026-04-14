package cmd

import (
	"encoding/xml"
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
	UID_DPRSchemaProperty     string
	UID_DPRSchemaType         string
	UID_QBMClrType            string
	IsAdditional              int
	IsReference               int
	IsVirtual                 bool
	UID_BaseDPRSchemaProperty *string
	SerializationBag          *string
	DataType                  *string
	SchemaType                string `mapstructure:",omitzero"`
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
	`Create a new synchronization schema property (DPRSchemaProperty) and return the UID_DPRSchemaProperty of the new property.`,
	[]string{"schema-type-id", "name", "data-type"},
	insertSchemaProperty,
)

func insertSchemaProperty(c *cobra.Command, db *sqlx.DB) error {
	clr, _ := c.Flags().GetString("clr-name")

	if len(clr) == 0 {
		schemaId, _ := c.Flags().GetString("schema-id")
		clr, _ = GetCLRForTarget(db, schemaId,
			"VI.Projector.Database.DatabaseSchemaProperty", "VI.Projector.Powershell.PoshSchemaProperty")
	}

	return ExecInsertCommand[DPRSchemaProperty](c, db, clr, newSchemaProperty_cmd)
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

type SerializationBagData struct {
	XMLName xml.Name `xml:"Data"`
	Name    string   `xml:",attr"`
	Type    string   `xml:",attr"`
	Value   string   `xml:",chardata"`
}
type SchemaPropertySerializationBag struct {
	XMLName xml.Name               `xml:"SerializationBag"`
	Version string                 `xml:",attr"`
	Datas   []SerializationBagData `xml:"Data"`
}

type ResolutionSchemaTypeInformation struct {
	XMLName                xml.Name `xml:"ResolutionSchemaTypeInformation"`
	Name                   string   `xml:",attr"`
	ResolutionPropertyName string   `xml:",attr"`
	ValuePropertyName      string   `xml:",attr"`
}
type ResolutionSchemaTypes struct {
	XMLName   xml.Name `xml:"ResolutionSchemaTypes"`
	TypeInfos []ResolutionSchemaTypeInformation
}

func InsertKeyBasedVirtualProperty(db *sqlx.DB,
	schemaTypeId string, baseAttrName string,
	lookupTable string, resolutionAttr string, valueAttr string,
) (string, error) {

	propName := fmt.Sprintf("vrt%s", baseAttrName)
	baseAttr, err := GetSchemaProperty(db, schemaTypeId, baseAttrName)
	if err != nil {
		return "", err
	}

	clrName := "VI.Projector.Schema.Properties.MultiKeyResolutionSchemaProperty"
	clrId, err := GetClrId(db, clrName)
	if err != nil {
		return "", err
	}

	id, objectKey, err := NewDPRKeys[DPRSchemaProperty](db)

	prop, err := newSchemaProperty(db, schemaTypeId, id, objectKey, propName, "String", clrId)
	if err != nil {
		return "", err
	}
	prop.UID_BaseDPRSchemaProperty = &baseAttr.UID_DPRSchemaProperty
	prop.IsVirtual = true
	prop.IsAdditional = -1
	prop.IsReference = -1

	// internal xml def for serialization bag
	rst := &ResolutionSchemaTypes{}
	rst.TypeInfos = []ResolutionSchemaTypeInformation{
		ResolutionSchemaTypeInformation{Name: lookupTable, ResolutionPropertyName: resolutionAttr, ValuePropertyName: valueAttr},
	}
	rst_str, _ := xml.MarshalIndent(rst, " ", "  ")

	// serialization bag container
	sb := &SchemaPropertySerializationBag{Version: "1.0"}
	sb.Datas = []SerializationBagData{
		SerializationBagData{Name: "ResolutionInformation", Type: "String", Value: string(rst_str)},
		SerializationBagData{Name: "IgnoreCase", Type: "Bool", Value: "False"},
		SerializationBagData{Name: "HandleResolutionFailureAsError", Type: "Bool", Value: "False"},
		SerializationBagData{Name: "UseSystemAttachedDataStore", Type: "Bool", Value: "True"},
		SerializationBagData{Name: "ReportResolutionFailure", Type: "Bool", Value: "True"},
	}
	sb_byte, _ := xml.MarshalIndent(sb, " ", "  ")
	sb_str := string(sb_byte)
	prop.SerializationBag = &sb_str

	err = InsertDPRObject[DPRSchemaProperty](db, prop)
	if err != nil {
		return "", err
	}

	return *prop.Name, nil
}

func GetSchemaProperty(db *sqlx.DB, schemaTypeId string, name string) (*DPRSchemaProperty, error) {

	wc := fmt.Sprintf("UID_DPRSchemaType='%s' and Name='%s'", schemaTypeId, name)
	t, err := dbx.GetStructSingletonByWC[DPRSchemaProperty](db, wc)
	return &t, err
}

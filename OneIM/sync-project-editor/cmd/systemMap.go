package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

// NB: sqlx assumes lowercase column names
type DPRSystemMap struct {
	oneim.Specials
	Displayable
	UID_DPRSystemMap        string
	UID_DPRShell            string
	UID_QBMClrType          string
	UID_LeftDPRSchemaClass  string
	UID_RightDPRSchemaClass string
	UID_BaseDPRSystemMap    *string
	Capabilities            *string
	IsHierarchyProjection   bool
	IsTemplate              bool
	MappingDirection        *string

	LeftClass  string   `mapstructure:",omitzero"`
	RightClass string   `mapstructure:",omitzero"`
	ParentMap  string   `mapstructure:",omitzero"`
	Rules      []string `mapstructure:",omitzero"`
}

var SystemMapCmd = CreateBaseCommand(
	"map",
	"sync map commands",
	`View and update synchronization map records (DPRSystemMap).`,
	showMaps,
)

var ShowSystemMapCmd = CreateShowCommand(
	"show details of one synchronization map",
	`View sync Map detail (DPRSystemMap).`,
	[]string{"shell"},
	showMaps,
)

func showMaps(c *cobra.Command, db *sqlx.DB) error {

	shellId, _ := c.Flags().GetString("shell")
	if len(shellId) == 0 {
		return errors.New("shell id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRShell='%s'`, shellId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSystemMap='%s'`, id)
	}

	return ShowDPRObjects[DPRSystemMap](db, wc, fillMapData)
}

func fillMapData(db *sqlx.DB, t *DPRSystemMap) error {

	t.LeftClass, _ = GetForeignDisplay(db, "DPRSchemaClass", t.UID_LeftDPRSchemaClass)
	t.RightClass, _ = GetForeignDisplay(db, "DPRSchemaClass", t.UID_RightDPRSchemaClass)

	if t.UID_BaseDPRSystemMap != nil {
		t.ParentMap, _ = GetForeignDisplay(db, "DPRSystemMap", *t.UID_BaseDPRSystemMap)
	}

	t.Rules, _ = GetChildDisplays(db, "DPRSystemMappingRule", "UID_DPRSystemMap", t.UID_DPRSystemMap)

	return nil
}

var InsertSystemMapCmd = CreateInsertCommand(
	"create a new synchronization map",
	`Create a new schema class attribute map (DPRSystemMap) and return the UID_DPRSystemMap of the new map.`,
	[]string{"shell", "name", "left-schema-class-id", "right-schema-class-id"},
	insertMap,
)

func insertMap(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemMap](c, db, "VI.Projector.Mapping.SystemMap", newMap_cmd)
}

func newMap_cmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemMap, error) {

	shellId, _ := c.Flags().GetString("shell")
	if len(shellId) == 0 {
		return nil, errors.New("shell id required")
	}

	leftClassId, _ := c.Flags().GetString("left-schema-class-id")
	if len(leftClassId) == 0 {
		return nil, errors.New("left class id required")
	}
	rightClassId, _ := c.Flags().GetString("right-schema-class-id")
	if len(rightClassId) == 0 {
		return nil, errors.New("right class id required")
	}

	mapDirection, _ := c.Flags().GetString("direction")
	if len(mapDirection) == 0 {
		return nil, errors.New("specify a map direction")
	}

	return newSystemMap(db,
		shellId,
		id, objectKey, name,
		leftClassId, rightClassId, mapDirection,
		clrId)
}

func newSystemMap(
	db *sqlx.DB,
	shellId,
	id string, objectKey string, name string,
	leftClassId string, rightClassId string, mapDirection string,
	clrId string,
) (*DPRSystemMap, error) {

	defaultCapabilities := "Default"

	t := DPRSystemMap{
		UID_DPRSystemMap:        id,
		UID_DPRShell:            shellId,
		UID_QBMClrType:          clrId,
		UID_LeftDPRSchemaClass:  leftClassId,
		UID_RightDPRSchemaClass: rightClassId,
		Specials:                oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &id,
			DisplayName: &name,
		},
		MappingDirection: &mapDirection,
		Capabilities:     &defaultCapabilities,
	}

	return &t, nil
}

var InsertSystemMapByNameCmd = CreateBaseCommand(
	"insert-by-name",
	"create a new synchronization map",
	`Create a new schema class attribute map (DPRSystemMap), referencing classes by name, and return the UID_DPRSystemMap of the new map.`,
	insertMapByName,
)

func insertMapByName(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemMap](c, db, "VI.Projector.Mapping.SystemMap", newMapByName_cmd)
}

func newMapByName_cmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemMap, error) {

	shellId, _ := c.Flags().GetString("shell")
	if len(shellId) == 0 {
		return nil, errors.New("shell id required")
	}

	leftClassName, _ := c.Flags().GetString("left")
	if len(leftClassName) == 0 {
		return nil, errors.New("left class name required")
	}
	lSchema, err := GetOneIMSchema(db, shellId)
	if err != nil {
		return nil, err
	}
	lSchemaClass, err := GetSchemaClassByName(db, shellId, lSchema.UID_DPRSchema, leftClassName)
	if err != nil {
		return nil, err
	}

	rightClassName, _ := c.Flags().GetString("right")
	if len(rightClassName) == 0 {
		return nil, errors.New("right class name required")
	}
	rSchema, err := GetTargetSystemSchema(db, shellId)
	if err != nil {
		return nil, err
	}
	rSchemaClass, err := GetSchemaClassByName(db, shellId, rSchema.UID_DPRSchema, rightClassName)
	if err != nil {
		return nil, err
	}

	mapDirection, _ := c.Flags().GetString("direction")
	if len(mapDirection) == 0 {
		return nil, errors.New("specify a map direction")
	}

	return newSystemMap(db,
		shellId,
		id, objectKey, name,
		lSchemaClass.UID_DPRSchemaClass, rSchemaClass.UID_DPRSchemaClass, mapDirection,
		clrId)

}

var UpdateSystemMapCmd = CreateUpdateCommand(
	"update an existing synchronization map",
	`Update attributes of a sync project map (DPRSystemMap).`,
	updateMap,
)

func updateMap(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSystemMap](c, db)
}

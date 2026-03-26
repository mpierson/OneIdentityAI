package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRSystemMappingRule struct {
	oneim.Specials
	Displayable
	UID_DPRSystemMappingRule       string
	UID_DPRSystemMap               string
	UID_QBMClrType                 string
	UID_MappingCondQBMClrType      *string
	PropertyLeft                   *string
	PropertyRight                  *string
	ConcurrenceBehavior            *string
	DisableMergeModeSupport        bool
	DoNotOverrideLeft              bool
	DoNotOverrideRight             bool
	HandleAsSingleValueLeft        bool
	HandleAsSingleValueRight       bool
	IgnoreCase                     bool
	IgnoreCaseDifferencesOnly      bool
	IgnoreMappingDirectionByCreate bool
	IsKeyRule                      bool
	IsRogueCorrectionEnabled       bool
	IsRogueDetectionEnabled        bool
	MappingCondition               *string
	MappingDirection               *string
	MvpOrderBehavior               *string
	PerformMappingContraProjection bool
	SortOrder                      int

	Map string `mapstructure:",omitzero"`
}

var MAP_RULE_DIRECTION_Left = "ToTheLeft"
var MAP_RULE_DIRECTION_None = "DoNotMap"

var MappingRuleCmd = CreateBaseCommand(
	"mapping-rule",
	"system mapping rule commands",
	`View and update synchronization rules in a map (DPRSystemMappingRule).`,
	showMappingRules,
)

var ShowMappingRuleCmd = CreateShowCommand(
	"show synchronization mapping rules",
	`View sync project mapping rules (DPRSystemMappingRule).`,
	[]string{"map-id"},
	showMappingRules,
)

func showMappingRules(c *cobra.Command, db *sqlx.DB) error {

	mapId, _ := c.Flags().GetString("map-id")
	if len(mapId) == 0 {
		return errors.New("map id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRSystemMap='%s'`, mapId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSystemMappingRule='%s'`, id)
	}

	return ShowDPRObjects[DPRSystemMappingRule](db, wc, fillMappingRuleData)
}

func fillMappingRuleData(db *sqlx.DB, t *DPRSystemMappingRule) error {

	t.Map, _ = GetForeignDisplay(db, "DPRSystemMap", t.UID_DPRSystemMap)

	return nil
}

var InsertMappingRuleCmd = CreateInsertCommand(
	"create a new synchronization mapping rule",
	`Create a new sync mapping rule (DPRSystemMappingRule).`,
	[]string{"map-id", "name", "left-property", "right-property", "direction"},
	insertMappingRule,
)

func insertMappingRule(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSystemMappingRule](c, db, clrId, newMappingRule)
}

func newMappingRule(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemMappingRule, error) {

	mapId, err := GetStructId_MustExist[DPRSystemMap](c, "map-id", db)
	if err != nil {
		return nil, err
	}

	leftProp, err := c.Flags().GetString("left-property")
	if err != nil {
		return nil, err
	}

	rightProp, err := c.Flags().GetString("right-property")
	if err != nil {
		return nil, err
	}

	// direction optional or n/a for some cmds
	direction, _ := c.Flags().GetString("direction")

	t := DPRSystemMappingRule{
		UID_DPRSystemMappingRule: id,
		UID_QBMClrType:           clrId,
		UID_DPRSystemMap:         mapId,
		Specials:                 oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
		PropertyLeft:     &leftProp,
		PropertyRight:    &rightProp,
		MappingDirection: &direction,
	}

	return &t, nil
}

var AddKeyBasedRuleCmd = createDPRCommand(
	"add-key-based-rule",
	"create a new key-based lookup mapping rule",
	`Create a new sync matching rule (DPRSystemMappingRule), with key-based lookup.`,
	[]string{
		"map-id", "name",
		"left-property", "right-property",
		"lookup-table", "left-key-attribute", "right-key-attribute",
	},
	insertKeyBasedMatchingRule,
)

func insertKeyBasedMatchingRule(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemMappingRule](c, db, "VI.Projector.Mapping.Rules.SinglePropertyComparisonRule", newKeyBasedMappingRule)
}

func newKeyBasedMappingRule(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemMappingRule, error) {

	var t *DPRSystemMappingRule

	mapId, err := GetStructId_MustExist[DPRSystemMap](c, "map-id", db)
	if err != nil {
		return t, err
	}
	m, err := dbx.GetStructSingleton[DPRSystemMap](db, mapId)
	if err != nil {
		return t, err
	}
	classId := m.UID_LeftDPRSchemaClass
	cls, err := dbx.GetStructSingleton[DPRSchemaClass](db, classId)
	if err != nil {
		return t, err
	}

	leftProp, err := c.Flags().GetString("left-property")
	if err != nil {
		return t, err
	}

	lookupTable, err := c.Flags().GetString("lookup-table")
	if err != nil {
		return t, err
	}
	leftKey, err := c.Flags().GetString("left-key-attribute")
	if err != nil {
		return t, err
	}
	rightKey, err := c.Flags().GetString("right-key-attribute")
	if err != nil {
		return t, err
	}

	// create a new virtual property for the left side key-based lookup
	vrtPropName, err := InsertKeyBasedVirtualProperty(db, cls.UID_DPRSchemaType, leftProp, lookupTable, rightKey, leftKey)
	if err != nil {
		return t, err
	}

	// create baseline map
	t, err = newMappingRule(c, db, id, objectKey, name, clrId)
	if err != nil {
		return t, err
	}

	// map to the new virt prop
	t.PropertyLeft = &vrtPropName
	t.MappingDirection = &MAP_RULE_DIRECTION_Left

	return t, nil
}

var InsertMatchingRuleCmd = createDPRCommand(
	"insert-matching-rule",
	"create a new synchronization object mapping rule",
	`Create a new sync matching rule (DPRSystemMappingRule).`,
	[]string{"map-id", "name", "left-property", "right-property"},
	insertMatchingRule,
)

func insertMatchingRule(c *cobra.Command, db *sqlx.DB) error {
	clrId, _ := c.Flags().GetString("clr-name")
	return ExecInsertCommand[DPRSystemMappingRule](c, db, clrId, newMatchingRule)
}

func newMatchingRule(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemMappingRule, error) {

	t, err := newMappingRule(c, db, id, objectKey, name, clrId)
	if err != nil {
		return nil, err
	}

	t.IsKeyRule = true
	t.MappingDirection = &MAP_RULE_DIRECTION_None

	return t, nil
}

var UpdateMappingRuleCmd = CreateUpdateCommand(
	"update an existing synchronization mapping rule",
	`Update attributes of a sync project mapping rule (DPRSystemMappingRule).`,
	updateMappingRule,
)

func updateMappingRule(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSystemMappingRule](c, db)
}

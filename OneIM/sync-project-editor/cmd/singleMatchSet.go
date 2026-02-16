package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

type DPRSystemObjectMatchSet struct {
	oneim.Specials
	Displayable
	UID_DPRSystemObjectMatchSet string
	UID_QBMClrType              string
	Methods                     []string `mapstructure:",omitzero"`
}

type DPRSystemObjectMatchSetMAsgn struct {
	oneim.Specials
	Displayable
	UID_DPRSystemObjectMatchSetMA string
	UID_DPRSystemObjectMatchSet   string
	UID_DPRSchemaMethod           string
	UID_QBMClrType                string
	Side                          string
	TargetProjectionDirection     *string
	Sequence                      int
	ConditionData                 *string
	UID_ConditionQBMClrType       *string
	MatchSet                      string `mapstructure:",omitzero"`
	Method                        string `mapstructure:",omitzero"`
}

var SingleMatchSetCmd = CreateBaseCommand(
	"match-set",
	"system match set commands",
	`View and update synchronization match sets (DPRSystemObjectMatchSet).`,
	showMatchSet,
)

var ShowSingleMatchSetCmd = CreateShowCommand(
	"show synchronization match set",
	`View sync project match set (DPRSystemObjectMatchSet).`,
	nil,
	showMatchSet,
)

func showMatchSet(c *cobra.Command, db *sqlx.DB) error {

	id, _ := c.Flags().GetString("id")

	wc := "1=1"
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSystemObjectMatchSet='%s'`, id)
	}

	return ShowDPRObjects[DPRSystemObjectMatchSet](db, wc, fillSingleMatchSetData)
}

func fillSingleMatchSetData(db *sqlx.DB, t *DPRSystemObjectMatchSet) error {

	methodWC := fmt.Sprintf(`
			UID_DPRSchemaMethod in (
				select UID_DPRSchemaMethod from DPRSystemObjectMatchSetMAsgn 
					where UID_DPRSystemObjectMatchSet='%s'
			)
		`, t.UID_DPRSystemObjectMatchSet)
	t.Methods, _ = GetStructDisplays(db, "DPRSchemaMethod", methodWC)

	return nil
}

var InsertSingleMatchSetCmd = CreateInsertCommand(
	"create a new synchronization match set",
	`Create a new sync match set (DPRSystemObjectMatchSet).`,
	[]string{"name"},
	insertMatchSet,
)

func insertMatchSet(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemObjectMatchSet](c, db, "VI.Projector.Projection.SystemObjectMatchingSet", newSingleMatchSetCmd)
}

func newSingleMatchSetCmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemObjectMatchSet, error) {
	return newSingleMatchSet(db, id, objectKey, name, clrId)
}

func newSingleMatchSet(
	db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemObjectMatchSet, error) {

	t := DPRSystemObjectMatchSet{
		UID_DPRSystemObjectMatchSet: id,
		UID_QBMClrType:              clrId,
		Specials:                    oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
	}

	return &t, nil
}

var UpdateSingleMatchSetCmd = CreateUpdateCommand(
	"update an existing synchronization match set",
	`Update attributes of a sync project match set (DPRSystemObjectMatchSet).`,
	updateMatchSet,
)

func updateMatchSet(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSystemObjectMatchSet](c, db)
}

func InsertNewMatchSetMethodAssignment(db *sqlx.DB, name string,
	matchSetId string, schemaMethodId string,
	side string, direction string,
) (*DPRSystemObjectMatchSetMAsgn, string, error) {

	fNew := func(db *sqlx.DB, id string, objectKey string, name string, clrId string) (*DPRSystemObjectMatchSetMAsgn, error) {
		return newSingleMatchSetMethodAssignment(db, id, objectKey, name, clrId, matchSetId, schemaMethodId, side, direction)
	}
	return InsertNewDPRObject[DPRSystemObjectMatchSetMAsgn](db, name, "VI.Projector.Projection.SchemaMethodAssignment", fNew)
}
func newSingleMatchSetMethodAssignment(db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
	matchSetId string, schemaMethodId string,
	side string, direction string,
) (*DPRSystemObjectMatchSetMAsgn, error) {

	t := DPRSystemObjectMatchSetMAsgn{
		UID_DPRSystemObjectMatchSetMA: id,
		UID_QBMClrType:                clrId,
		UID_DPRSystemObjectMatchSet:   matchSetId,
		UID_DPRSchemaMethod:           schemaMethodId,
		Side:                          side,
		TargetProjectionDirection:     &direction,
		Specials:                      oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name: &name,
		},
	}

	return &t, nil
}

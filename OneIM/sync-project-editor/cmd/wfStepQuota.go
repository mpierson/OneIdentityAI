package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

type DPRProjectionStepQuota struct {
	oneim.Specials
	Displayable
	UID_DPRProjectionStepQuota string
	UID_QBMClrType             string
	Settings                   string

	Steps []string `mapstructure:",omitzero"`
}

var StepQuotaCmd = CreateBaseCommand(
	"workflow-step-quota",
	"system workflow step quota commands",
	`View and update synchronization workflow quota (DPRProjectionStepQuota).`,
	showStepQuotas,
)

var ShowStepQuotaCmd = CreateShowCommand(
	"show synchronization workflow step quotas",
	`View sync project workflow step quotas (DPRProjectionStepQuota).`,
	[]string{"map-id"},
	showStepQuotas,
)

func showStepQuotas(c *cobra.Command, db *sqlx.DB) error {

	id, _ := c.Flags().GetString("id")

	wc := "1=1"
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRProjectionStepQuota='%s'`, id)
	}

	return ShowDPRObjects[DPRProjectionStepQuota](db, wc, fillStepQuotaData)
}

func fillStepQuotaData(db *sqlx.DB, t *DPRProjectionStepQuota) error {

	return nil
}

var InsertStepQuotaCmd = CreateInsertCommand(
	"create a new synchronization workflow step quota",
	`Create a new sync workflow step quota (DPRProjectionStepQuota).`,
	[]string{"name", "type"},
	insertStepQuota,
)

func insertStepQuota(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRProjectionStepQuota](c, db, "VI.Projector.Projection.ProjectionStepQuota", newStepQuota)
}

func newStepQuota(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRProjectionStepQuota, error) {

	settings, err := c.Flags().GetString("type")
	if err != nil {
		return nil, err
	}

	t := DPRProjectionStepQuota{
		UID_DPRProjectionStepQuota: id,
		UID_QBMClrType:             clrId,
		Specials:                   oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
		Settings: settings,
	}

	return &t, nil
}

var UpdateStepQuotaCmd = CreateUpdateCommand(
	"update an existing synchronization workflow step quota",
	`Update attributes of a sync project workflow step quota (DPRProjectionStepQuota).`,
	updateStepQuota,
)

func updateStepQuota(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRProjectionStepQuota](c, db)
}

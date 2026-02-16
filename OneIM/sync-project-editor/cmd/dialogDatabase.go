package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

type DialogDatabase struct {
	oneim.Specials
	Displayable
	UID_Database       string
	ConnectionProvider *string
	ConnectionString   *string
	CustomerName       *string
	CustomerPrefix     *string
	DataOrigin         int
	EditionDescription *string
	EditionName        *string
	EditionVersion     *string
	ElementColor       *string
	IsMainDatabase     bool
	ProductionLevel    int
}

var DialogDatabaseCmd = CreateBaseCommand(
	"database",
	"Identity Manager database record commands",
	`View database details (DialogDatabase).`,
	showDialogDatabases,
)

var ShowDialogDatabaseCmd = CreateShowCommand(
	"show details of Identity Manager database",
	`View Identity Manager database details (DialogDatabase).`,
	nil,
	showDialogDatabases,
)

func showDialogDatabases(c *cobra.Command, db *sqlx.DB) error {

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf("1=1")
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DialogDatabase='%s'`, id)
	}

	return ShowDPRObjects[DialogDatabase](db, wc, nil)
}

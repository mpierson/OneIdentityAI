package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
)

type QBMClrType struct {
	oneim.Specials
	Displayable
	UID_QBMClrType   string
	FullTypeName     string
	Assembly         string
	ExposedInterface *string
	Alias            *string
}

var ClrTypeCmd = CreateBaseCommand(
	"clr",
	"Identity Manager runtime object commands",
	`Runtime type (QBMClrType).`,
	showClrTypes,
)

var ShowClrTypeCmd = CreateShowCommand(
	"show details of Identity Manager type",
	`View Identity Manager runtime types (QBMClrType).`,
	nil,
	showClrTypes,
)

func showClrTypes(c *cobra.Command, db *sqlx.DB) error {

	id, _ := c.Flags().GetString("id")
	name, _ := c.Flags().GetString("name")

	wc := fmt.Sprintf("1=1")
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_QBMClrType='%s'`, id)
	} else if len(name) > 0 {
		wc = wc + fmt.Sprintf(` AND FullTypeName='%s'`, name)
	}

	return ShowDPRObjects[QBMClrType](db, wc, nil)
}

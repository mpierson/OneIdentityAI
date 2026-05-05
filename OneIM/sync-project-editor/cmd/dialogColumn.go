package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DialogColumn struct {
	oneim.Specials
	Displayable
	UID_DialogColumn string
	ColumnName       string
	Caption          *string
	IsUID            bool
	IsPKMember       bool
	IsCrypted        bool
	IsForeignKey     bool
	SchemaDataType   string
	UID_DialogTable  string
	Commentary       *string
	SchemaDataLen    int
}

type DialogTable struct {
	oneim.Specials
	Displayable
	UID_DialogTable string
	TableName       string
	TableType       string
	PKName1         *string
	PKName2         *string
	UsageType       *string
}

func GetDialogTable(db *sqlx.DB, tableName string) (*DialogTable, error) {

	if !IsValidIdOrName(tableName) {
		return nil, errors.New("invalid table name")
	}

	wc := fmt.Sprintf(`TableName='%s'`, tableName)
	t, err := dbx.GetStructSingletonByWC[DialogTable](db, wc)
	return &t, err
}

func GetDialogColumnsForTable(db *sqlx.DB, UID_DialogTable string) ([]DialogColumn, error) {

	if !IsValidIdOrName(UID_DialogTable) {
		return nil, errors.New("invalid table name")
	}

	wc := fmt.Sprintf(`UID_DialogTable='%s'`, UID_DialogTable)
	return dbx.GetStructData[DialogColumn](db, "DialogColumn", wc)
}

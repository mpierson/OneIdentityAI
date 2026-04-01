package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type Displayable struct {
	Name                 *string `mapstructure:",omitzero"`
	DisplayName          *string `mapstructure:",omitzero"`
	DisplayNameQualified *string `mapstructure:",omitzero"`
	Description          *string `mapstructure:",omitzero"`
}

func IsValidIdOrName(v string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(v)
}

func CheckRequiredFlag(c *cobra.Command, flagName string) error {

	f, err := c.Flags().GetString(flagName)
	if err != nil {
		return err
	}
	if len(f) == 0 {
		return errors.New("missing flag " + flagName)
	}

	return nil
}

func CheckRequiredFlags(c *cobra.Command, flagNames []string) error {

	for _, n := range flagNames {
		err := CheckRequiredFlag(c, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDBFromContext(c *cobra.Command) (*sqlx.DB, error) {

	cctx := c.Context()
	if cctx != nil {
		if db, ok := cctx.Value("db_connection").(*sqlx.DB); ok {
			return db, nil
		}
	}

	return nil, errors.New("Missing DB in context")
}

var leader = "                     "

func GetHeading(buf map[string]interface{}) string {

	name := buf["DisplayName"]
	if name == nil {
		name = buf["Name"]
	}
	if name == nil {
		name = buf["DisplayNameQualified"]
	}
	if name == nil {
		name = buf["DisplayValue"]
	}

	separator := " : "
	desc := buf["Description"]
	if desc == nil {
		desc = ""
		separator = ""
	}

	return fmt.Sprintf("%v%s%v", name, separator, desc)
}
func PrintHeading(buf map[string]interface{}, indent int) {
	fmt.Println(fmt.Sprintf("%s%s", leader[0:indent], GetHeading(buf)))
}

func PrintAttr(buf map[string]interface{}, label string, attrName string, indent int) {
	fmt.Println(fmt.Sprintf("%s%s : %v", leader[0:indent], label, buf[attrName]))
}

func (d *Displayable) GetHeader() string {

	empty := ""
	name := d.Name
	if name == nil {
		name = d.DisplayName
		if name == nil {
			name = d.DisplayNameQualified
			if name == nil {
				name = &empty
			}
		}
	}

	separator := " : "
	desc := d.Description
	if desc == nil {
		desc = &empty
		separator = ""
	}

	return fmt.Sprintf("%v%s%v", name, separator, desc)
}
func (d *Displayable) PrintHeader(indent int) {
	fmt.Println(fmt.Sprintf("%v%s", leader[0:indent], d.GetHeader()))
}

func PrintField(label string, v interface{}, indent int) {
	fmt.Println(fmt.Sprintf("%s%s : %v", leader[0:indent], label, v))
}

func PrintStruct(t interface{}, indent int) {
	mypp := pp.New()
	mypp.SetColoringEnabled(true)
	mypp.SetExportedOnly(true)
	mypp.SetOmitEmpty(true)
	mypp.Println(t)
}

func GetForeignDisplayByObjectKey(db *sqlx.DB, objectKey string) (string, error) {

	if len(objectKey) == 0 {
		return "", nil
	}

	tableName, r, err := dbx.GetSingletonTableDataByKey(db, objectKey)
	if err != nil {
		fmt.Println("failed to fetch " + objectKey + " record")
		return "", err
	}
	if r != nil {
		return fmt.Sprintf(`%s [%s]`, GetHeading(r), tableName), nil
	}

	return "", errors.New(objectKey + " record not found")
}

func GetForeignDisplay(db *sqlx.DB, table string, uid string) (string, error) {

	if len(uid) == 0 {
		return "", nil
	}

	keyCol := "UID_" + table

	r, err := dbx.GetForeignSingleton(db, table, keyCol, uid)
	if err != nil {
		fmt.Println("failed to fetch " + table + " records")
		return "", err
	}
	if r != nil {
		return GetHeading(r), nil
	}

	return "", errors.New(table + " record not found")
}

func GetChildDisplays(db *sqlx.DB, table string, keyColumn string, keyValue string) ([]string, error) {

	wc := fmt.Sprintf(`%s = '%s'`, keyColumn, keyValue)
	return GetStructDisplays(db, table, wc)
}

func GetStructDisplays(db *sqlx.DB, table string, wc string) ([]string, error) {

	buffs, err := dbx.GetBufferedTableData(db, table, wc, 100)
	if err != nil {
		return nil, err
	}

	displays := make([]string, 0)
	for _, buf := range buffs {
		displays = append(displays, GetHeading(buf))
	}

	return displays, nil
}

func ShowDPRObjects[T any](db *sqlx.DB, whereClause string, fnFillStructData func(db *sqlx.DB, t *T) error) error {
	var t T
	structs, err := dbx.GetStructData[T](db, reflect.TypeOf(t).Name(), whereClause)
	if err != nil {
		return err
	}

	for _, s := range structs {
		if fnFillStructData != nil {
			err = fnFillStructData(db, &s)
			if err != nil {
				return err
			}
		}
		PrintStruct(s, 0)
	}

	return nil
}

func GetClrId(db *sqlx.DB, clrName string) (string, error) {
	v, err := dbx.GetTableValue(db, "QBMCLRType", "UID_QBMCLRType", fmt.Sprintf(`FullTypeName='%s'`, clrName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch CLR %s\n", clrName)
	}
	return v, err
}

func GetClrName(db *sqlx.DB, clrId string) (string, error) {
	v, err := dbx.GetTableValue(db, "QBMCLRType", "FullTypeName", fmt.Sprintf(`UID_QBMClrType='%s'`, clrId))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch CLR %s\n", clrId)
	}
	return v, err
}

func InsertDPRObject[T any](db *sqlx.DB, t *T) error {

	stmt, err := dbx.GenerateInsertStmt(*t)
	if err != nil {
		return err
	}

	tx := db.MustBegin()
	_, err = tx.NamedExec(stmt, t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to execute: "+stmt)
		return err
	}
	tx.Commit()

	return nil
}

func ExecInsertCommand[T any](c *cobra.Command, db *sqlx.DB,
	clrTypeName string,
	fNew func(*cobra.Command, *sqlx.DB, string, string, string, string) (*T, error),
) error {

	fWithoutCmd := func(db1 *sqlx.DB, id string, objectKey string, name string, clrId string) (*T, error) {
		return fNew(c, db1, id, objectKey, name, clrId)
	}

	name, _ := c.Flags().GetString("name")

	_, newId, err := InsertNewDPRObject[T](db, name, clrTypeName, fWithoutCmd)
	if err != nil {
		return err
	}

	//fmt.Println(fmt.Sprintf(`UID_%s: %s`, reflect.TypeOf(*t).Name(), newId))
	fmt.Println(newId)

	return nil
}

func InsertNewDPRObject[T any](db *sqlx.DB, name string, clrTypeName string,
	fNew func(db *sqlx.DB, id string, objectKey string, name string, clrId string) (*T, error),
) (*T, string, error) {

	var t T
	id, err := dbx.GetNewId(db)
	objectKey := oneim.MakeObjectKey(reflect.TypeOf(t).Name(), id)

	clrId := ""
	if len(clrTypeName) != 0 {
		clrId, err = GetClrId(db, clrTypeName)
		if err != nil {
			return &t, id, err
		}
	}

	// allow caller to populate DPR struct
	t1, err := fNew(db, id, objectKey, name, clrId)
	if err != nil {
		return t1, id, err
	}

	// TODO: check for existing row with same name?

	// push to DB
	err = InsertDPRObject[T](db, t1)
	if err != nil {
		return t1, id, err
	}

	return t1, id, nil
}

func ExecUpdateCommand[T any](c *cobra.Command, db *sqlx.DB) error {

	id, err := c.Flags().GetString("id")
	if err != nil {
		return err
	} else if len(id) == 0 {
		return errors.New("Missing object id")
	}

	content, err := GetCommandContent(c)
	if err != nil {
		return err
	}

	t, err := UpdateDPRObject[T](db, id, content)
	if err != nil {
		return err
	}
	PrintStruct(t, 0)

	return nil
}

func GetCommandContent(c *cobra.Command) (string, error) {

	content, err := c.Flags().GetString("content")
	if err != nil {
		return "", err
	}

	// if content is '-' then read from stdin
	if len(content) > 0 && content == "-" {
		reader := c.InOrStdin()
		if reader == nil {
			return "", errors.New("provide content as parameter or use '-'")
		}

		bytes, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		content = string(bytes)
	}

	return content, nil
}

func UpdateDPRObject[T any](db *sqlx.DB, id string, content string) (T, error) {

	t, err := dbx.GetStructSingleton[T](db, id)
	if err != nil {
		return t, err
	}

	json.Unmarshal([]byte(content), &t)

	err = UpdateStruct[T](db, &t, id)
	if err != nil {
		return t, err
	}

	return t, nil
}

func UpdateStruct[T any](db *sqlx.DB, t *T, id string) error {

	stmt, err := dbx.GenerateUpdateStmt(*t, id)
	if err != nil {
		return err
	}

	tx := db.MustBegin()
	_, err = tx.NamedExec(stmt, t)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func CreateBaseCommand(name string,
	shortDescr string, longDescr string,
	showFn func(c *cobra.Command, db *sqlx.DB) error) *cobra.Command {

	return createDPRCommand(name, shortDescr, longDescr, nil, showFn)
}

func CreateShowCommand(
	shortDescr string, longDescr string,
	requiredFlags []string,
	showFn func(c *cobra.Command, db *sqlx.DB) error) *cobra.Command {

	return createDPRCommand("show", shortDescr, longDescr, requiredFlags, showFn)
}

func CreateInsertCommand(
	shortDescr string, longDescr string,
	requiredFlags []string,
	fn func(c *cobra.Command, db *sqlx.DB) error) *cobra.Command {

	return createDPRCommand("insert", shortDescr, longDescr, requiredFlags, fn)
}
func CreateUpdateCommand(shortDescr string, longDescr string,
	fn func(c *cobra.Command, db *sqlx.DB) error) *cobra.Command {

	return createDPRCommand("update", shortDescr, longDescr, []string{"id", "content"}, fn)
}

func createDPRCommand(name string,
	shortDescr string, longDescr string,
	requiredFlags []string,
	dbFn func(c *cobra.Command, db *sqlx.DB) error) *cobra.Command {

	fPreRunE := func(c *cobra.Command, args []string) error {
		return CheckRequiredFlags(c, requiredFlags)
	}

	fRunE := func(c *cobra.Command, args []string) error {
		db, err := GetDBFromContext(c)
		if err != nil {
			return err
		}

		err = dbFn(c, db)
		if err != nil {
			return err
		}

		return nil
	}

	return createCommand(name, shortDescr, longDescr, fPreRunE, fRunE)
}

func createCommand(name string,
	shortDescr string, longDescr string,
	fPreRunE func(c *cobra.Command, args []string) error,
	fRunE func(c *cobra.Command, args []string) error) *cobra.Command {

	return &cobra.Command{
		Use:     name,
		Short:   shortDescr,
		Long:    longDescr,
		PreRunE: fPreRunE,
		RunE:    fRunE,
	}
}

func GetStructId_MustExist[T any](c *cobra.Command, flag string, db *sqlx.DB) (string, error) {

	id, _ := c.Flags().GetString(flag)
	if len(id) == 0 {
		return "", errors.New(flag + " is required")
	}
	if exists, err := dbx.StructExists[T](db, id); err != nil || !exists {
		return "", errors.New("invalid " + flag)
	}

	return id, nil
}

var COLUMN_TYPE_MAP = map[string]string{
	"VARCHAR":   "string",
	"NVARCHAR":  "string",
	"NCHAR":     "string",
	"CHAR":      "string",
	"VARBINARY": "Binary",
	"INT":       "Integer",
	"BIT":       "Integer",
	"FLOAT":     "Float",
	"DATETIME":  "DateTime",
	"BOOL":      "Boolean",
}

func GetTableColumns(db *sqlx.DB, table string) ([][2]string, error) {

	// get columns with un-translated SQL types
	cols, err := dbx.GetTableColumns(db, table)
	if err != nil {
		return nil, err
	}

	if len(cols) == 0 {
		return nil, errors.New("columns not found")
	}

	// convert types to DPR standards
	for i := range cols {
		if cType, mapContainsKey := COLUMN_TYPE_MAP[cols[i][1]]; mapContainsKey {
			cols[i][1] = cType
		} else {
			return nil, errors.New("unmapped SQL type: " + cols[i][1])
		}
	}

	return cols, nil
}

func FireDBEvent(db *sqlx.DB, objectType string, whereClause string, eventName string, priority int) error {

	genProcId, err := dbx.GetNewId(db)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`exec QBM_PJobCreate_HOFireEvent
        @objecttype = '%s',
        @whereclause = '%s' , 
		@EventName = '%s',
        @priority = %v,
        @GenProcID = '%s'`, objectType, whereClause, eventName, priority, genProcId)
	_, err = db.Exec(q)

	return err
}

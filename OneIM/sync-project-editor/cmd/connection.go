package cmd

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRSystemConnection struct {
	oneim.Specials
	Displayable
	UID_DPRSystemConnection       string
	UID_DPRShell                  string
	UID_DPRSchema                 string
	UID_QBMClrType                string
	UID_ConnectorQBMClrType       string
	UID_ParamDescriptorQBMClrType *string `mapstructure:",omitempty"`
	UID_DPRSystemScope            *string
	UID_DPRSystemScopeForRefer    *string
	ConnectionParameter           string
	ConnectRetryCount             int
	ConnectRetryDelay             int
	DefaultDisplay                *string
	IsReadOnly                    bool
	JournalLogFailedObjectChanges bool
	JournalLogObjectChanges       bool
	JournalLogPropertyChanges     bool
	JournalMessageContexts        *string
	WriteJournal                  bool
	Schema                        string `mapstructure:",omitempty"`
	SystemScope                   string `mapstructure:",omitempty"`
	SystemScopeForRefer           string `mapstructure:",omitempty"`
}

var ConnectionCmd = CreateBaseCommand(
	"connection",
	"sync Connection commands",
	`View and update synchronization Connection records (DPRSystemConnection).`,
	showConnections,
)

var ShowConnectionCmd = CreateShowCommand(
	"show details of one synchronization Connection",
	`View sync Connection detail (DPRSystemConnection).`,
	[]string{"shell"},
	showConnections,
)

func showConnections(c *cobra.Command, db *sqlx.DB) error {

	shellId, _ := c.Flags().GetString("shell")
	if len(shellId) == 0 {
		return errors.New("shell id required")
	}

	id, _ := c.Flags().GetString("id")

	wc := fmt.Sprintf(`UID_DPRShell='%s'`, shellId)
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSystemConnection='%s'`, id)
	}

	return ShowDPRObjects[DPRSystemConnection](db, wc, fillConnectionData)
}

func fillConnectionData(db *sqlx.DB, t *DPRSystemConnection) error {

	t.Schema, _ = GetForeignDisplay(db, "DPRSchema", t.UID_DPRSchema)
	if t.UID_DPRSystemScope != nil {
		t.SystemScope, _ = GetForeignDisplay(db, "DPRSystemScope", *t.UID_DPRSystemScope)
	}
	if t.UID_DPRSystemScopeForRefer != nil {
		t.SystemScopeForRefer, _ = GetForeignDisplay(db, "DPRSystemScope", *t.UID_DPRSystemScopeForRefer)
	}

	return nil
}

var InsertConnectionCmd = CreateInsertCommand(
	"create a new synchronization connection",
	`Create a new synchronization connection (DPRSystemConnection), and return the UID_DPRSystemConnection of the new connection.`,
	[]string{"shell", "name", "schema-id"},
	insertConnection,
)

func insertConnection(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemConnection](c, db, "VI.Projector.Connection.SystemConnection", newConnection_cmd)
}

func newConnection_cmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemConnection, error) {

	schemaId, err := GetStructId_MustExist[DPRSchema](c, "schema-id", db)
	if err != nil {
		return nil, err
	}

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	conType, _ := c.Flags().GetString("type")
	if len(conType) == 0 {
		return nil, errors.New("Missing connection type")
	}

	var conClrType string
	var paramClrType string
	if conType == "MainConnection" {
		conClrType = "VI.Projector.Database.DatabaseConnectorDescriptor"
		paramClrType = "VI.Projector.Database.DatabaseConnectionParameterDescriptor"
	} else {
		conClrType, _ = c.Flags().GetString("connector-type")
		if len(conClrType) == 0 {
			return nil, errors.New("missing connector CLR type")
		}
	}

	// connection string is optional?
	conString, _ := c.Flags().GetString("parameters")

	return newConnection(db,
		shellId, schemaId,
		id, objectKey, name,
		conType, conClrType,
		conString, paramClrType,
		clrId)
}

func newConnection(
	db *sqlx.DB,
	shellId string, schemaId string,
	id string, objectKey string, name string,
	conType string, conClrType string,
	conParameter string, paramClrType string,
	clrId string,
) (*DPRSystemConnection, error) {

	shellName, err := dbx.GetTableValue(db, "DPRShell", "DisplayName", fmt.Sprintf(`UID_DPRShell='%s'`, shellId))
	if err != nil {
		return nil, err
	}

	displayName := fmt.Sprintf(`%s - %s`, shellName, name)

	connectorClrId, err := GetClrId(db, conClrType)
	if err != nil {
		return nil, err
	}

	var paramClrId string
	if len(paramClrType) > 0 {
		paramClrId, err = GetClrId(db, paramClrType)
		if err != nil {
			return nil, err
		}
	}

	t := DPRSystemConnection{
		UID_DPRSystemConnection: id,
		UID_DPRShell:            shellId,
		UID_QBMClrType:          clrId,
		UID_ConnectorQBMClrType: connectorClrId,
		UID_DPRSchema:           schemaId,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
		DefaultDisplay: &name,
		Displayable: Displayable{
			DisplayNameQualified: &displayName,
			Name:                 &conType,
		},
	}
	if len(conParameter) > 0 {
		t.ConnectionParameter = conParameter
	}
	if len(paramClrId) > 0 {
		t.UID_ParamDescriptorQBMClrType = &paramClrId
	}

	return &t, nil
}

var InsertOneIMConnectionCmd = CreateBaseCommand(
	"insert-oneim-connection",
	"create a new synchronization connection for Identity Manager",
	`Create a new synchronization connection (DPRSystemConnection) for Identity Manager and return the UID_DPRSystemConnection of the new connection.`,
	insertOneIMConnection,
)

func insertOneIMConnection(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemConnection](c, db, "VI.Projector.Connection.SystemConnection", newOneIMConnection)
}

func newOneIMConnection(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemConnection, error) {

	// TODO: check for existing OneIM connection

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	// connection string is optional?
	conString, _ := c.Flags().GetString("parameters")

	// lookup primary OneIM schema in given shell
	schema, err := GetOneIMSchema(db, shellId)
	if err != nil {
		return nil, err
	}

	return newConnection(db,
		shellId, schema.UID_DPRSchema,
		id, objectKey, *schema.Name,
		"MainConnection", "VI.Projector.Database.DatabaseConnectorDescriptor",
		conString, "VI.Projector.Database.DatabaseConnectionParameterDescriptor",
		clrId)
}

var InsertTargetSystemConnectionCmd = CreateBaseCommand(
	"insert-target-system-connection",
	"create a new synchronization connection for a target system",
	`Create a new synchronization connection (DPRSystemConnection) for a target system and return the UID_DPRSystemConnection of the new connection.`,
	insertTargetSystemConnection,
)

func insertTargetSystemConnection(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemConnection](c, db, "VI.Projector.Connection.SystemConnection", newTargetSystemConnection)
}

func newTargetSystemConnection(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemConnection, error) {

	// TODO: check for existing OneIM connection

	shellId, err := GetStructId_MustExist[DPRShell](c, "shell", db)
	if err != nil {
		return nil, err
	}

	// lookup target system schema
	schema, err := GetTargetSystemSchema(db, shellId)
	if err != nil {
		return nil, err
	}

	connectorType, _ := c.Flags().GetString("connector-type")
	if len(connectorType) == 0 {
		return nil, errors.New("missing connector CLR type")
	}

	// connection string is optional?
	conString, _ := c.Flags().GetString("parameters")
	conStringType := ""

	return newConnection(db,
		shellId, schema.UID_DPRSchema,
		id, objectKey, *schema.Name,
		"ConnectedSystemConnection", connectorType,
		conString, conStringType,
		clrId)
}

var UpdateConnectionCmd = CreateUpdateCommand(
	"update an existing connection",
	`Update attributes of a sync connection (DPRSystemConnection).`,
	updateConnection,
)

func updateConnection(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSystemConnection](c, db)
}

func GetAllConnections(db *sqlx.DB, shellId string) ([]DPRSystemConnection, error) {
	if !IsValidIdOrName(shellId) {
		return nil, errors.New("invalid shell id")
	}
	wc := fmt.Sprintf(`UID_DPRShell='%s'`, shellId)
	return dbx.GetStructData[DPRSystemConnection](db, "DPRSystemConnection", wc)
}

func GetDefaultTargetSystemConnection(db *sqlx.DB, shellId string) (*DPRSystemConnection, error) {
	if !IsValidIdOrName(shellId) {
		return nil, errors.New("invalid shell id")
	}

	ts, err := GetAllConnections(db, shellId)
	if err != nil {
		return nil, err
	}
	for _, t := range ts {
		if *t.Name == "ConnectedSystemConnection" {
			return &t, nil
		}
	}

	return nil, nil
}

// ----- connector XML definition ------------------------------------

func CompressConnectorXML(xmlPath string) (string, error) {
	xmlBytes, err := os.ReadFile(xmlPath)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", xmlPath, err)
	}

	return CompressConnectorXMLString(string(xmlBytes))
}

// CompressConnectorXML encodes a connector definition XML file for use in an
// Identity Manager target system connection string (DefinitionXml parameter).
//
// The encoding matches .NET's three-step process:
//  1. Base64-encode the raw XML bytes
//  2. Compress the Base64 string with raw DEFLATE (equivalent to DeflateStream)
//  3. Base64-encode the compressed bytes
//
// # Claude Code April 2026
func CompressConnectorXMLString(xmlContent string) (string, error) {

	b64 := base64.StdEncoding.EncodeToString([]byte(xmlContent))

	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return "", fmt.Errorf("creating deflate writer: %w", err)
	}
	if _, err = w.Write([]byte(b64)); err != nil {
		return "", fmt.Errorf("compressing: %w", err)
	}
	if err = w.Close(); err != nil {
		return "", fmt.Errorf("flushing deflate writer: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

var CompressConnectorDefinitionCmd = createCommand(
	"compress-connector-definition",
	"compress the given connector definition xml, suitable for synchronization connection parameter string",
	`Returns a compressed and Base64 encoded version of given input string, according to requirements for connector definition XML embedded in a synchronization project connection parameter string.`,
	nil,
	compressDefinitionCmd,
)

func compressDefinitionCmd(c *cobra.Command, args []string) error {
	xmlString, err := c.Flags().GetString("xml")
	if err != nil {
		return err
	} else if len(xmlString) == 0 {
		return errors.New("missing xml content")
	}

	compressedXml, err := CompressConnectorXMLString(xmlString)
	if err != nil {
		return nil
	} else if len(compressedXml) == 0 {
		return errors.New("failed to compress xml")
	}

	fmt.Println(compressedXml)
	return nil
}

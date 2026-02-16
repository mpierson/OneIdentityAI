package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRRootObjConnectionInfo struct {
	oneim.Specials
	UID_DPRRootObjConnectionInfo string
	UID_DPRSystemConnection      string
	UID_DPRSystemVariableSet     string
	UID_QBMServer                string
	ObjectKeyRoot                string
	IsOffline                    bool
	IsOfflineModeAvailable       bool
	Connection                   string `mapstructure:",omitzero"`
	VariableSet                  string `mapstructure:",omitzero"`
	RootObject                   string `mapstructure:",omitzero"`
	Server                       string `mapstructure:",omitzero"`
}

type QBMServer struct {
	oneim.Specials
	UID_QBMServer string
	Ident_Server  string
	FQDN          *string
	IPV4          *string
	IPV6          *string
}

func fillRootObjData(db *sqlx.DB, t *DPRRootObjConnectionInfo) error {

	t.Connection, _ = GetForeignDisplay(db, "DPRSystemConnection", t.UID_DPRSystemConnection)
	t.VariableSet, _ = GetForeignDisplay(db, "DPRSystemVariableSet", t.UID_DPRSystemVariableSet)
	t.RootObject, _ = GetForeignDisplayByObjectKey(db, t.ObjectKeyRoot)
	t.Server, _ = GetQBMServerDisplay(db, t.UID_QBMServer)

	return nil
}

func GetRootObjDisplay(db *sqlx.DB, rootObjId string) (string, error) {
	t, err := dbx.GetStructSingleton[DPRRootObjConnectionInfo](db, rootObjId)
	if err != nil {
		return "", err
	}

	err = fillRootObjData(db, &t)
	if err != nil {
		return "", err
	}

	display := fmt.Sprintf(`%s / %s / %s`, t.Connection, t.VariableSet, t.RootObject)
	return display, nil
}

func NewRootObj(db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
	connectionId string,
	variableSetId string,
	serverId string,
	objectKeyRoot string,
) (*DPRRootObjConnectionInfo, error) {

	t := DPRRootObjConnectionInfo{
		UID_DPRRootObjConnectionInfo: id,
		UID_DPRSystemConnection:      connectionId,
		UID_DPRSystemVariableSet:     variableSetId,
		UID_QBMServer:                serverId,
		ObjectKeyRoot:                objectKeyRoot,
		Specials: oneim.Specials{
			XObjectKey: objectKey,
		},
	}

	return &t, nil
}

func InsertNewRootObj(db *sqlx.DB,
	connectionId string,
	variableSetId string,
	serverId string,
	rootObjectKey string) (*DPRRootObjConnectionInfo, error) {

	fNew := func(db *sqlx.DB, id string, objectKey string, name string, clrId string) (*DPRRootObjConnectionInfo, error) {
		return NewRootObj(db, id, objectKey, name, clrId, connectionId, variableSetId, serverId, rootObjectKey)
	}

	t, _, err := InsertNewDPRObject[DPRRootObjConnectionInfo](db, "", "", fNew)
	return t, err

}

func GetQBMServerDisplay(db *sqlx.DB, serverId string) (string, error) {
	t, err := dbx.GetStructSingleton[QBMServer](db, serverId)
	if err != nil {
		return "", err
	}

	err = fillQBMServerData(db, &t)
	if err != nil {
		return "", err
	}

	display := fmt.Sprintf(`%s / %s`, t.Ident_Server, *t.FQDN)
	return display, nil
}

func fillQBMServerData(db *sqlx.DB, t *QBMServer) error {

	return nil
}

func GetServerByName(db *sqlx.DB, serverName string) (*QBMServer, error) {

	ts, err := dbx.GetStructData[QBMServer](db, "QBMServer", fmt.Sprintf(`Ident_Server='%s'`, serverName))
	if err != nil {
		return nil, err
	}
	if len(ts) == 0 {
		return nil, errors.New("object not found: " + serverName)
	} else if len(ts) > 1 {
		return nil, errors.New(fmt.Sprintf("Too many rows returned for %s in QBMServer", serverName))
	}

	return &ts[0], nil

}

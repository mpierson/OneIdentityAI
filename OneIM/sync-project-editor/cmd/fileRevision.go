package cmd

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type QBMFileRevision struct {
	oneim.Specials
	UID_QBMFileRevision string
	FileSize            int64
	FileName            string
	FileVersion         *string
	DeployTargets       []string `mapstructure:",omitzero"`
	FilePath            *string  `mapstructure:"-"`
}
type QBMFileRevision_write struct {
	QBMFileRevision
	FileContent []byte
	HashValue   []byte
}

type QBMDeployTarget struct {
	oneim.Specials
	UID_QBMDeployTarget   string
	Ident_QBMDeployTarget string
	DisplayValue          *string
}

type QBMFileHasDeployTarget struct {
	oneim.Specials
	UID_QBMFileHasDeployTarget string
	UID_QBMFileRevision        string
	ObjectKeyDeployTarget      string
}

var FileCmd = CreateBaseCommand(
	"file",
	"binary file commands",
	`Manage binary files associated with sync project (QBMFileRevision).`,
	showFiles,
)

var ShowFilesCmd = CreateShowCommand(
	"show details of custom binary files",
	`View custom binary files in Identity Manager (QBMFileRevision).`,
	[]string{},
	showFiles,
)

func showFiles(c *cobra.Command, db *sqlx.DB) error {

	id, _ := c.Flags().GetString("id")

	wc := "UID_QBMFileRevision like 'CCC-%'"
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_QBMFileRevision='%s'`, id)
	}

	return ShowDPRObjects[QBMFileRevision](db, wc, fillFileData)
}

func fillFileData(db *sqlx.DB, t *QBMFileRevision) error {
	wc := fmt.Sprintf(
		`XObjectKey in (select ObjectKeyDeployTarget from QBMFileHasDeployTarget where UID_QBMFileRevision='%s')`,
		t.UID_QBMFileRevision,
	)
	t.DeployTargets, _ = GetStructDisplays(db, "QBMDeployTarget", wc)
	return nil
}

var InsertFileCmd = CreateInsertCommand(
	"insert a file in the database",
	`Insert a file in the Identity Manager database (QBMFileRevision), and return the UID_QBMFileRevision of the new entry.`,
	[]string{"file", "file-version"},
	insertFile,
)

func insertFile(c *cobra.Command, db *sqlx.DB) error {

	// ----- insert file record -------------

	fWithoutCmd := func(db1 *sqlx.DB, id string, objectKey string, name string, clrId string) (*QBMFileRevision, error) {
		return newFileRevision(c, db1, id, objectKey, name, clrId)
	}

	t, newId, err := InsertNewDPRObject[QBMFileRevision](db, "placeholder", "", fWithoutCmd)
	if err != nil {
		return err
	}

	fmt.Println("uploading file")

	dat, err := os.ReadFile(*t.FilePath)
	if err != nil {
		return err
	}
	err = UpdateFileContent(db, t, dat)
	if err != nil {
		return err
	}

	fmt.Println("file uploaded, now pushing data to job servers")

	// ----- add Server assignment ---------

	_, err = assignDeployTarget(db, t, "Server")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*360)
	defer cancel()

	// wait for job server task(s) to complete
	err = WaitForTaskStart(db, ctx, "TaskName='CheckAndUpdate'", 2*time.Second)
	if err != nil {
		return err
	}
	fmt.Println("job server update started")
	err = WaitForTaskFinish(db, ctx, "TaskName='CheckAndUpdate'", 10*time.Second)
	if err != nil {
		return err
	}

	fmt.Println("job server update running, verification will start shortly")

	// ---- trigger module info refresh -----

	err = UpdateModuleInfo(db)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*90)
	defer cancel()
	err = WaitForTaskFinish(db, ctx, "TaskName='CallMethod' AND JobChainName like '%%UpdateModuleInfo%%'", 10*time.Second)
	if err != nil {
		return err
	}

	// ----- add Jobserver assignment

	_, err = assignDeployTarget(db, t, "Jobserver")
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*360)
	defer cancel()

	// wait for job server task(s) to complete
	err = WaitForTaskStart(db, ctx, "TaskName='CheckAndUpdate'", 2*time.Second)
	if err != nil {
		return err
	}
	fmt.Println("still verifiying")
	err = WaitForTaskFinish(db, ctx, "TaskName='CheckAndUpdate'", 10*time.Second)
	if err != nil {
		return err
	}

	fmt.Println("done")
	fmt.Println("file identifier: " + newId)

	return nil
}

func newFileRevision(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*QBMFileRevision, error) {

	filePath, _ := c.Flags().GetString("file")

	// check for existing
	ts, err := GetFileRevisions(db, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if len(ts) > 0 {
		return nil, errors.New("file already exists in Identity Manager")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileVersion, _ := c.Flags().GetString("file-version")

	t := QBMFileRevision{
		UID_QBMFileRevision: id,
		Specials:            oneim.NewSpecials(objectKey, "sped"),
		FileName:            filepath.Base(filePath),
		FileSize:            fileInfo.Size(),
		FileVersion:         &fileVersion,
		FilePath:            &filePath,
	}

	return &t, nil
}

func GetFileRevisions(db *sqlx.DB, fileName string) ([]QBMFileRevision, error) {
	if !IsValidIdOrName(fileName) {
		return nil, errors.New("Invalid file name " + fileName)
	}

	wc := fmt.Sprintf(`FileName='%s'`, fileName)
	return dbx.GetStructData[QBMFileRevision](db, "QBMFileRevision", wc)
}

func UpdateFileContent(db *sqlx.DB, t *QBMFileRevision, content []byte) error {

	hash, err := getFileHash(content)

	q := `UPDATE QBMFileRevision SET FileContent = :content, HashValue = :hash WHERE UID_QBMFileRevision = :id`
	params := map[string]interface{}{
		"id":      t.UID_QBMFileRevision,
		"content": content,
		"hash":    hash,
	}

	_, err = db.NamedExec(q, params)
	if err != nil {
		return err
	}

	return nil
}

func getFileHash(content []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(content); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

var UpdateModuleInfoCmd = createDPRCommand(
	"update-module-info",
	"update internal accounting of custom files",
	`Trigger an update of the custom module info XML blob with file metadata.`,
	[]string{},
	updateModuleInfo_internal,
)

func updateModuleInfo_internal(c *cobra.Command, db *sqlx.DB) error {
	return UpdateModuleInfo(db)
}

func UpdateModuleInfo(db *sqlx.DB) error {
	genProcId, err := dbx.GetNewId(db)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(` exec QBM_PJobCreate_HOCallMethod 
		@objecttype = 'QBMModuleDef' , 
		@whereclause = 'UID_ModuleDef = ''CCC-Moduledefinition''' , 
		@save = 1 , 
		@MethodName = 'UpdateModuleInfoXML' , 
		@priority = 10,
		@GenProcID = '%s'`, genProcId)
	_, err = db.Exec(q)
	return err
}

//		@ObjectKeysAffected = DEFAULT,

var AssignDeployTargetCmd = createDPRCommand(
	"assign-deploy-target",
	"assign a deploy target",
	`Assign a deploy target to an existing file. (QBMFileHasDeployTarget)`,
	[]string{"id", "deploy-target"},
	assignDeployTargetC,
)

func assignDeployTargetC(c *cobra.Command, db *sqlx.DB) error {

	id, err := GetStructId_MustExist[QBMFileRevision](c, "id", db)
	if err != nil {
		return err
	}
	fr, err := dbx.GetStructSingleton[QBMFileRevision](db, id)
	if err != nil {
		return err
	}

	targetName, err := c.Flags().GetString("deploy-target")
	if err != nil {
		return err
	}

	newId, err := assignDeployTarget(db, &fr, targetName)
	if err != nil {
		return err
	}

	fmt.Println(newId)
	return nil
}

func assignDeployTarget(db *sqlx.DB, fr *QBMFileRevision, targetName string) (string, error) {

	if !IsValidIdOrName(targetName) {
		return "", errors.New("Invalid deploy target " + targetName)
	}
	wc := fmt.Sprintf(`Ident_QBMDeployTarget='%s'`, targetName)
	dt, err := dbx.GetStructSingletonByWC[QBMDeployTarget](db, wc)
	if err != nil {
		return "", err
	}

	fWithoutCmd := func(db1 *sqlx.DB, id string, objectKey string, name string, clrId string) (*QBMFileHasDeployTarget, error) {
		t := QBMFileHasDeployTarget{
			UID_QBMFileHasDeployTarget: id,
			Specials:                   oneim.NewSpecials(objectKey, "sped"),
			UID_QBMFileRevision:        fr.UID_QBMFileRevision,
			ObjectKeyDeployTarget:      dt.XObjectKey,
		}

		return &t, nil
	}

	_, newId, err := InsertNewDPRObject[QBMFileHasDeployTarget](db, "placeholder", "", fWithoutCmd)
	if err != nil {
		return newId, err
	}

	return newId, nil
}

var SyncFilesCmd = createDPRCommand(
	"sync-files",
	"sync files on a specific job server",
	`Trigger a file synchronization event for the given job server (QBMServer)`,
	[]string{"server-name"},
	syncFiles,
)

func syncFiles(c *cobra.Command, db *sqlx.DB) error {

	serverName, err := c.Flags().GetString("server-name")
	if err != nil {
		return err
	}
	if !IsValidIdOrName(serverName) {
		return errors.New("Invalid server name " + serverName)
	}
	wc := fmt.Sprintf(`Ident_Server='%s'`, serverName)
	server, err := dbx.GetStructSingletonByWC[QBMServer](db, wc)
	if err != nil {
		return err
	}

	wc = fmt.Sprintf(`XObjectKey=''%s''`, server.Specials.XObjectKey)
	return FireDBEvent(db, "QBMServer", wc, "CheckVersion", 2)
}

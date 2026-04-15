package cmd

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRJournal struct {
	UID_DPRJournal             string
	UID_DPRProjectionStartInfo *string
	JobId                      *string
	ProjectionState            *string
	CreationTime               sql.NullTime
	CompletionTime             sql.NullTime
	Messages                   []DPRJournalMessage `mapstructure:",omitzero"`
}
type DPRJournalMessage struct {
	UID_DPRJournalMessage string
	UID_DPRJournal        string
	MessageString         string
	MessageType           string
	MessageContext        *string
	SequenceNumber        int
}

func FillJournalData(db *sqlx.DB, t *DPRJournal) error {

	wc := fmt.Sprintf(`UID_DPRJournal = '%s'`, t.UID_DPRJournal)
	var err error
	t.Messages, err = dbx.GetStructData[DPRJournalMessage](db, "DPRJournalMessage", wc)
	if err != nil {
		return err
	}

	return nil
}

func getJournalMessagesInt(db *sqlx.DB, UID_DPRJournal string) (*DPRJournal, error) {

	t, err := dbx.GetStructSingleton[DPRJournal](db, UID_DPRJournal)
	if err != nil {
		return &t, err
	}
	err = FillJournalData(db, &t)
	if err != nil {
		return &t, err
	}

	return &t, nil
}
func GetJournalMessages(db *sqlx.DB, UID_DPRStartInfo string, UID_Job string) (*DPRJournal, error) {

	wc := fmt.Sprintf(`UID_DPRProjectionStartInfo='%s' AND JobId='%s'`, UID_DPRStartInfo, UID_Job)
	t, err := dbx.GetStructSingletonByWC[DPRJournal](db, wc)
	if err != nil {
		return &t, err
	}

	err = FillJournalData(db, &t)
	if err != nil {
		return &t, err
	}

	return &t, nil
}

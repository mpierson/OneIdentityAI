package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"pso.oneidentity.com/sped/oneim"
	"pso.oneidentity.com/sped/oneim/dbx"
)

type DPRSystemObjectMatchSets struct {
	oneim.Specials
	Displayable
	UID_DPRSystemObjectMatchSets  string
	UID_QBMClrType                string
	UID_DifferenceLeftToRightSet  *string
	UID_DifferenceRightToLeftSet  *string
	UID_IntersectionDifferenceSet *string
	UID_IntersectionEqualitySet   *string

	DifferenceLeftToRightSet  string `mapstructure:",omitzero"`
	DifferenceRightToLeftSet  string `mapstructure:",omitzero"`
	IntersectionDifferenceSet string `mapstructure:",omitzero"`
	IntersectionEqualitySet   string `mapstructure:",omitzero"`

	WorkflowSteps []string `mapstructure:",omitzero"`
}

var MatchSetsCmd = CreateBaseCommand(
	"match-sets",
	"system match set commands",
	`View and update synchronization match sets (DPRSystemObjectMatchSets).`,
	showMatchSets,
)

var ShowMatchSetsCmd = CreateShowCommand(
	"show synchronization match sets",
	`View sync project match sets (DPRSystemObjectMatchSets).`,
	nil,
	showMatchSets,
)

func showMatchSets(c *cobra.Command, db *sqlx.DB) error {

	id, _ := c.Flags().GetString("id")

	wc := "1=1"
	if len(id) > 0 {
		wc = wc + fmt.Sprintf(` AND UID_DPRSystemObjectMatchSets='%s'`, id)
	}

	return ShowDPRObjects[DPRSystemObjectMatchSets](db, wc, fillMatchSetData)
}

func fillMatchSetData(db *sqlx.DB, t *DPRSystemObjectMatchSets) error {

	t.DifferenceLeftToRightSet = getObjectMatchSetDisplay(db, t.UID_DifferenceLeftToRightSet)
	t.DifferenceRightToLeftSet = getObjectMatchSetDisplay(db, t.UID_DifferenceRightToLeftSet)
	t.IntersectionDifferenceSet = getObjectMatchSetDisplay(db, t.UID_IntersectionDifferenceSet)
	t.IntersectionEqualitySet = getObjectMatchSetDisplay(db, t.UID_IntersectionEqualitySet)

	t.WorkflowSteps, _ = GetChildDisplays(db, "DPRProjectionConfigStep", "UID_DPRSystemObjectMatchSets", t.UID_DPRSystemObjectMatchSets)

	return nil
}

func getObjectMatchSetDisplay(db *sqlx.DB, id *string) string {
	if id == nil {
		return ""
	} else {
		val, err := GetForeignDisplay(db, "DPRSystemObjectMatchSet", *id)
		if err != nil {
			return ""
		} else {
			return val
		}
	}
}

var InsertMatchSetsCmd = CreateInsertCommand(
	"create a new synchronization match set",
	`Create a new sync match set (DPRSystemObjectMatchSets).`,
	[]string{"name"},
	insertMatchSets,
)

func insertMatchSets(c *cobra.Command, db *sqlx.DB) error {
	return ExecInsertCommand[DPRSystemObjectMatchSets](c, db, "VI.Projector.Projection.SystemObjectMatchingSets", newMatchSetsCmd)
}

func newMatchSetsCmd(
	c *cobra.Command, db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemObjectMatchSets, error) {
	return NewMatchSets(db, id, objectKey, name, clrId)
}

func NewMatchSets(
	db *sqlx.DB,
	id string, objectKey string, name string,
	clrId string,
) (*DPRSystemObjectMatchSets, error) {

	t := DPRSystemObjectMatchSets{
		UID_DPRSystemObjectMatchSets: id,
		UID_QBMClrType:               clrId,
		Specials:                     oneim.Specials{XObjectKey: objectKey},
		Displayable: Displayable{
			Name:        &name,
			DisplayName: &name,
		},
	}

	return &t, nil
}

// return new match set collection with all four match sets
func CreateNewMatchSetsWithDefaults(db *sqlx.DB, stepId string) (*DPRSystemObjectMatchSets, error) {

	var t1 DPRSystemObjectMatchSets

	matchSetsId, err := dbx.GetNewId(db)
	if err != nil {
		return &t1, err
	}
	objectKey := oneim.MakeObjectKey("DPRSystemObjectMatchSets", matchSetsId)
	clrId, _ := GetClrId(db, "VI.Projector.Projection.SystemObjectMatchingSets")
	// new match set collection will re-use step id (this appears to be the only way to tie the set back to wf)
	name := fmt.Sprintf(`Step%s`, stepId)

	t, err := NewMatchSets(db, matchSetsId, objectKey, name, clrId)
	if err != nil {
		return &t1, err
	}

	err = InsertDPRObject[DPRSystemObjectMatchSets](db, t)
	if err != nil {
		return t, err
	}

	err = addDefaultMatchSets(db, t, matchSetsId)

	return t, err
}

var UpdateMatchSetsCmd = CreateUpdateCommand(
	"update an existing synchronization match set",
	`Update attributes of a sync project match set (DPRSystemObjectMatchSets).`,
	updateMatchSets,
)

func updateMatchSets(c *cobra.Command, db *sqlx.DB) error {
	return ExecUpdateCommand[DPRSystemObjectMatchSets](c, db)
}

var AddDefaultMatchSetsCmd = createDPRCommand(
	"add-default-sets",
	"add the default match sets to collection",
	`Add default match sets to a collection (DPRSystemObjectMatchSet -> DPRSystemObjectMatchSets).`,
	[]string{"id"},
	addDefaultMatchSets_cmd,
)

func addDefaultMatchSets_cmd(c *cobra.Command, db *sqlx.DB) error {

	id, err := GetStructId_MustExist[DPRSystemObjectMatchSets](c, "id", db)
	if err != nil {
		return err
	}
	t, err := dbx.GetStructSingleton[DPRSystemObjectMatchSets](db, id)
	if err != nil {
		return err
	}

	return addDefaultMatchSets(db, &t, id)
}

func addDefaultMatchSets(db *sqlx.DB, t *DPRSystemObjectMatchSets, setsId string) error {

	id1, _ := addMatchSetToCollection(db, setsId, "DifferenceLeftToRight")
	t.UID_DifferenceLeftToRightSet = &id1
	id2, _ := addMatchSetToCollection(db, setsId, "DifferenceRightToLeft")
	t.UID_DifferenceRightToLeftSet = &id2
	id3, _ := addMatchSetToCollection(db, setsId, "IntersectionWithDifferences")
	t.UID_IntersectionDifferenceSet = &id3
	id4, _ := addMatchSetToCollection(db, setsId, "IntersectionWithoutDifferences")
	t.UID_IntersectionEqualitySet = &id4

	err := UpdateStruct[DPRSystemObjectMatchSets](db, t, setsId)
	if err != nil {
		return err
	}

	return nil
}

func addMatchSetToCollection(db *sqlx.DB, collectionId string, matchSetName string) (string, error) {

	id, err := dbx.GetNewId(db)
	if err != nil {
		return "", err
	}
	objectKey := oneim.MakeObjectKey("DPRSystemObjectMatchSet", id)

	clrId, _ := GetClrId(db, "VI.Projector.Projection.SystemObjectMatchingSet")

	t, err := newSingleMatchSet(db, id, objectKey, matchSetName, clrId)
	err = InsertDPRObject[DPRSystemObjectMatchSet](db, t)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetMatchSetByName(db *sqlx.DB, collectionId string, matchSetName string) (*DPRSystemObjectMatchSet, error) {

	if !IsValidIdOrName(collectionId) {
		return nil, errors.New("invalid collection id: " + matchSetName)
	}
	collection, err := dbx.GetStructSingleton[DPRSystemObjectMatchSets](db, collectionId)
	if err != nil {
		return nil, err
	}

	if !IsValidIdOrName(matchSetName) {
		return nil, errors.New("invalid match set: " + matchSetName)
	}
	matchSetId := ""
	switch matchSetName {
	case "DifferenceLeftToRight":
		matchSetId = *collection.UID_DifferenceLeftToRightSet
	case "DifferenceRightToLeft":
		matchSetId = *collection.UID_DifferenceRightToLeftSet
	case "IntersectionWithDifferences":
		matchSetId = *collection.UID_IntersectionDifferenceSet
	case "IntersectionWithoutDifferences":
		matchSetId = *collection.UID_IntersectionEqualitySet
	default:
		return nil, errors.New("invalid match set: " + matchSetName)
	}
	if len(matchSetId) == 0 {
		return nil, errors.New("failed to lookup match set")
	}

	matchSet, err := dbx.GetStructSingleton[DPRSystemObjectMatchSet](db, matchSetId)
	return &matchSet, err
}

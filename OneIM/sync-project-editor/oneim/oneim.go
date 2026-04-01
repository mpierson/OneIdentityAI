package oneim

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

type Specials struct {
	XDateInserted sql.NullTime
	XUserInserted sql.NullString
	XDateUpdated  sql.NullTime
	XUserUpdated  sql.NullString
	XObjectKey    string
}

func NewSpecials(XObjectKey string, user string) Specials {
	ts := time.Now()
	return Specials{
		XDateInserted: sql.NullTime{Time: ts.UTC(), Valid: true},
		XUserInserted: sql.NullString{String: user, Valid: true},
		XDateUpdated:  sql.NullTime{Time: ts.UTC(), Valid: true},
		XUserUpdated:  sql.NullString{String: user, Valid: true},
		XObjectKey:    XObjectKey,
	}
}

// attributes with create/update metadata
var ATTR_Metadata = []string{"XDateInserted", "XUserInserted", "XDateUpdated", "XUserUpdated"}

// for use in f_attr, to include metadata as element attrs in xml
var MAP_Metadata = makeMetaMap()

func makeMetaMap() map[string]string {
	m := make(map[string]string, len(ATTR_Metadata))
	for _, e := range ATTR_Metadata {
		m[e] = e
	}
	return m
}

// deconstruct object key (many-to-many tables not yet supported)
func GetKeyParts(objectKey string) (t string, ids []string) {

	// table name
	re_t := regexp.MustCompile(`<Key><T>([A-Za-z-]+)</T>`)
	if !re_t.MatchString(objectKey) {
		return "", nil
	}
	t = re_t.FindStringSubmatch(objectKey)[1]

	// primary keys
	re_p := regexp.MustCompile(`<P>([A-Za-z0-9-]+)</P>`)
	if !re_p.MatchString(objectKey) {
		return "", nil
	}

	match_p := re_p.FindAllStringSubmatch(objectKey, -1)

	// extract keys from matches -- each match element includes: [0] = matched string, [1] = subgroup in match
	ids = make([]string, len(match_p), len(match_p))
	for i, e := range match_p {
		ids[i] = e[1]
	}

	return t, ids
}

func MakeObjectKey(table string, id string) string {
	return fmt.Sprintf(`<Key><T>%s</T><P>%s</P></Key>`, table, id)
}

func GetNonNullFieldNames(t interface{}) ([]string, error) {

	// decode struct into map
	var result = make(map[string]interface{}, 0)
	config := &mapstructure.DecoderConfig{
		Squash:    true,
		DecodeNil: false,
		Result:    &result,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, nil
	}
	err = decoder.Decode(t)
	if err != nil {
		return nil, err
	}

	// extract names, check for invalid null values
	names := make([]string, 0)
	for k, v := range result {
		switch t := v.(type) {
		case map[string]interface{}:
			isValid, exists := t["Valid"]
			if exists {
				switch ivT := isValid.(type) {
				case bool:
					if ivT {
						names = append(names, k)
					}
				default:
					names = append(names, k)
				}
			} else {
				names = append(names, k)
			}
		default:
			names = append(names, k)
		}
	}

	return names, nil
}

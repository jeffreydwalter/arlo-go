package util

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func PrettyPrint(data interface{}) string {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprint("error:", err)
	}
	return fmt.Sprint(string(j))
}

/*
type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}
*/

func Decode(b interface{}, s interface{}) error {
	if err := mapstructure.Decode(b, s); err != nil {
		return errors.Wrap(err, "failed to create "+reflect.TypeOf(s).String()+" object")
	}
	return nil
}

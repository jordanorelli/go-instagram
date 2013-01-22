package instagram

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"strconv"
	"time"
)

type Timestamp struct{ time.Time }

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 0)
	if err != nil {
		log.Println("Unable to parse timestamp from json")
	}
	*t = Timestamp{time.Unix(i, 0)}
	return nil
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

func (t *Timestamp) GetBSON() (interface{}, error) {
	return t.Time, nil
}

func (t *Timestamp) SetBSON(raw bson.Raw) error {
	tim := new(time.Time)
	err := raw.Unmarshal(&tim)
	if err != nil {
		return err
	}
	*t = Timestamp{*tim}
	return nil
}

type NString string

func (s *NString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(s))
}

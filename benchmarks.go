package main

import (
	"encoding/gob"
	"encoding/json"
	"math"
	"time"
)

var sampleValue_old = Sample{
	Name:    "testing",
	Age:     18,
	IsHuman: true,
	Nested: Nested{
		Inner: 123,
		Attrs: []string{"aa", "bb", "cc"},
	},
	LastVal: 777,
}

func MarshallingSimpser() {
	Marshal(&sampleValue)
}

func MarshallingJSON() {
	json.Marshal(&sampleValue)
}

func init() {
	gob.Register(&SmallStruct{})
	gob.Register(&Sample{})
}

// SmallStruct is a test structure of small size.
type SmallStruct struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

var sampleValue = SmallStruct{
	Name:     "some pretty long name",
	BirthDay: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
	Phone:    "1234567890",
	Siblings: 3,
	Spouse:   true,
	Money:    123456.789,
}

func (s *SmallStruct) Marshal(w *Writer) {
	w.Object()
	w.String("name")
	w.String(s.Name)
	w.String("birthDay")
	w.Int64(s.BirthDay.UnixMilli())
	w.String("phone")
	w.String(s.Phone)
	w.String("siblings")
	w.Int(s.Siblings)
	w.String("spouse")
	w.Bool(s.Spouse)
	w.String("money")
	w.Float64(s.Money)
	w.End()
}

func (s *SmallStruct) Unmarshal(r *Reader, rv *Value) {
	r.IterateObject(rv, func(key string, v *Value) {
		switch key {
		case "name":
			s.Name = v.VString
		case "birthDay":
			s.BirthDay = time.UnixMilli(v.VInt64)
		case "phone":
			s.Phone = v.VString
		case "siblings":
			s.Siblings = int(v.VInt32)
		case "spouse":
			s.Spouse = v.VBool
		case "money":
			s.Money = math.Float64frombits(uint64(v.VInt64))
		default:
			panic("unknown field")
		}
	})
}

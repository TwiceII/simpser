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

var marshalledSampleValue_2 []byte
var marshalledSampleValue []byte = []byte{2, 4, 0, 0, 0, 4, 110, 97, 109, 101, 4, 0, 0, 0, 21, 115, 111, 109, 101, 32, 112, 114, 101, 116, 116, 121, 32, 108, 111, 110, 103, 32, 110, 97, 109, 101, 4, 0, 0, 0, 8, 98, 105, 114, 116, 104, 68, 97, 121, 6, 0, 0, 1, 138, 232, 136, 228, 0, 4, 0, 0, 0, 5, 112, 104, 111, 110, 101, 4, 0, 0, 0, 10, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 4, 0, 0, 0, 8, 115, 105, 98, 108, 105, 110, 103, 115, 5, 0, 0, 0, 3, 4, 0, 0, 0, 6, 115, 112, 111, 117, 115, 101, 8, 1, 4, 0, 0, 0, 5, 109, 111, 110, 101, 121, 7, 64, 254, 36, 12, 159, 190, 118, 201, 1}

func init() {
	gob.Register(&SmallStruct{})
	gob.Register(&Sample{})
	//dt, err := Marshal(&sampleValue)
	//if err != nil {
	//	panic(err)
	//}
	//marshalledSampleValue = dt
	//fmt.Println("marshalled: ", marshalledSampleValue)
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

func (s *SmallStruct) Unmarshal(r *Reader, rv Value) {
	r.IterateObject(rv, func(key string, v Value) {
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

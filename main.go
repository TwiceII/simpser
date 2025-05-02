package main

import (
	"fmt"
)

type ValueType byte

func (vt ValueType) String() string {
	switch vt {
	case ValueType_Error:
		return "Error"
	case ValueType_End:
		return "End"
	case ValueType_Object:
		return "Object"
	case ValueType_Array:
		return "Array"
	case ValueType_Int64:
		return "Int64"
	case ValueType_Bool:
		return "Bool"
	case ValueType_String:
		return "String"
	default:
		return fmt.Sprintf("Unknown(%d)", vt)
	}
}

const (
	ValueType_Error ValueType = iota
	ValueType_End
	ValueType_Object
	ValueType_Array
	ValueType_String
	ValueType_Int32
	ValueType_Int64
	ValueType_Float64
	ValueType_Bool
)

type Value struct {
	Type  ValueType
	Depth int

	VInt64  int64
	VInt32  int32
	VString string
	VBool   bool
}

// ------------------ reader ------------------

// ------------------ inspection ------------------

func printIndent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
}

//func PrintReaderObject(r *Reader) {
//	obj := r.Read()
//	PrintValue(r, obj, 0)
//}

//func PrintValue(r *Reader, val Value, depth int) {
//	var k, v Value
//	var count int
//	switch val.Type {
//	case ValueType_Object:
//		fmt.Print("{\n")
//		for {
//			if ok := r.iterObject(&val, &k, &v); !ok {
//				break
//			}
//			if count > 0 {
//				fmt.Print(",\n")
//			}
//			count++
//			printIndent(depth + 1)
//			PrintValue(r, k, depth+1)
//			fmt.Print(": ")
//			PrintValue(r, v, depth+1)
//		}
//		if count > 0 {
//			fmt.Print("\n")
//		}
//		printIndent(depth)
//		fmt.Print("}")
//		break
//	case ValueType_Array:
//		fmt.Print("[\n")
//		for {
//			if ok := r.iterArray(&val, &v); !ok {
//				break
//			}
//			if count > 0 {
//				fmt.Print(",\n")
//			}
//			count++
//			printIndent(depth + 1)
//			PrintValue(r, v, depth+1)
//		}
//		if count > 0 {
//			fmt.Print("\n")
//		}
//		printIndent(depth)
//		fmt.Print("]")
//		break
//	case ValueType_Int64:
//		fmt.Printf("%d", val.VInt64)
//		break
//	case ValueType_Bool:
//		if *val.VBool {
//			fmt.Print("true")
//		} else {
//			fmt.Print("false")
//		}
//		break
//	case ValueType_String:
//		fmt.Printf("%q", *val.VString)
//		break
//	default:
//		fmt.Printf("ERROR: unknown value type: %d", val.Type)
//	}
//}

// ------------------ map ------------------

//func WriteVal(w *Writer, v any) {
//	switch v := v.(type) {
//	case map[string]any:
//		w.Object()
//		for k, v := range v {
//			w.String(&k)
//			WriteVal(w, v)
//		}
//		w.End()
//	case []any:
//		w.Array()
//		for _, item := range v {
//			WriteVal(w, item)
//		}
//		w.End()
//	case string:
//		w.String(&v)
//	case int:
//		w.Int(&v)
//	case bool:
//		w.Bool(&v)
//	default:
//		panic(fmt.Errorf("unsupported type %T", v))
//	}
//}

// ------------------ main ------------------

type Nested struct {
	Inner int64
	Attrs []string
}

type Sample struct {
	Name    string
	Age     int64
	IsHuman bool
	Nested  Nested
	LastVal int64
}

func (s *Sample) Marshal_old(w *Writer) {
	w.Object()
	w.Int64(s.Age)
	w.String("name")
	w.String(s.Name)
	w.String("age")
	w.Int64(s.Age)
	w.End()
}

func (s *Sample) Marshal(w *Writer) {
	// @@TODO: ergonomics
	//w.Obj().
	//	String("name", s.Name).
	//	Int64("age", s.Age).
	//	Bool("isHuman", s.IsHuman).
	//	Int64("lastVal", s.LastVal).
	//	End()

	w.Object()
	w.String("name")
	w.String(s.Name)
	w.String("age")
	w.Int64(s.Age)
	w.String("isHuman")
	w.Bool(s.IsHuman)

	// nested
	w.String("nested")
	w.Object() // < nested
	w.String("inner")
	w.Int64(s.Nested.Inner)
	w.String("attrs")
	w.Array() // < array
	for _, attr := range s.Nested.Attrs {
		w.String(attr)
	}
	w.End() // > array
	w.End() // > nested

	w.String("lastVal")
	w.Int64(s.LastVal)

	w.End()
}

func (n *Nested) Unmarshal(r *Reader, rv Value) {
	r.IterateObject(rv, func(key string, v Value) {
		switch key {
		case "inner":
			n.Inner = v.VInt64
			break
		case "attrs":
			r.IterateArray(v, func(v Value) {
				n.Attrs = append(n.Attrs, v.VString)
			})
			break
		}
	})
}

func (s *Sample) Unmarshal(r *Reader, rv Value) {
	r.IterateObject(rv, func(key string, v Value) {
		switch key {
		case "name":
			s.Name = v.VString
			break
		case "age":
			s.Age = v.VInt64
			break
		case "isHuman":
			s.IsHuman = v.VBool
			break
		case "nested":
			s.Nested.Unmarshal(r, v)
			break
		case "lastVal":
			s.LastVal = v.VInt64
			break
		}
	})
}

func main() {
	dt, err := Marshal(&sampleValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("marshalled: ", dt)
	var s2 SmallStruct
	if err := Unmarshal(dt, &s2); err != nil {
		panic("panic on unmarshalling: " + err.Error())
	}
	fmt.Println("unmarshalled: ", s2)
}

func main_3() {
	fmt.Println("simpser start")

	s := Sample{
		Name:    "testing",
		Age:     18,
		IsHuman: true,
		Nested: Nested{
			Inner: 123,
			Attrs: []string{"aa", "bb", "cc"},
		},
		LastVal: 777,
	}
	//if err := NewWriter(&bb).Encode(&s); err != nil {
	//	panic("panic on marshalling: " + err.Error())
	//}
	//fmt.Println("marshalled: ", bb.Bytes())

	dt, err := Marshal(&s)
	if err != nil {
		panic(err)
	}
	fmt.Println("marshalled: ", dt)

	var s2 Sample
	if err := Unmarshal(dt, &s2); err != nil {
		panic("panic on unmarshalling: " + err.Error())
	}
	fmt.Println("unmarshalled: ", s2)

	//PrintReaderObject(NewReader(&bb))

	//var s2 Sample
	//if err := NewReader(&bb).Decode(&s2); err != nil {
	//	fmt.Println("error on unmarshalling: " + err.Error())
	//}
	//fmt.Println("unmarshalled: ", s2)
}

//func main_2() {
//	fmt.Println("simpser map")
//	var bb bytes.Buffer
//	w := &Writer{
//		wr: &bb,
//	}
//	m := map[string]any{
//		"name":     "testing_from_map",
//		"age":      18,
//		"isHuman":  true,
//		"lastVal":  999,
//		"otherVal": "dafaq",
//		"nested": map[string]any{
//			"inner": 123,
//			"attrs": []any{"a", "b", "c"},
//		},
//	}
//	WriteVal(w, m)
//	fmt.Println("marshalled: ", bb.Bytes())
//	fmt.Println(string(bb.Bytes()))
//
//	var s2 Sample
//	if err := NewReader(&bb).Decode(&s2); err != nil {
//		fmt.Println("error on unmarshalling: " + err.Error())
//	}
//	fmt.Println("unmarshalled: ", s2)
//
//	//r2 := &Reader{
//	//	rd: &bb,
//	//}
//	//s2 := &Sample{}
//	//if err := s2.Unmarshal(r2); err != nil {
//	//	panic("panic on unmarshalling: " + err.Error())
//	//}
//	//fmt.Println("unmarshalled: ", s2)
//
//	//fmt.Println("marshalled: ", bb.Bytes())
//	//fmt.Println(string(bb.Bytes()))
//}

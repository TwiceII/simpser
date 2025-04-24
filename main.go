package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type ValueType byte

const (
	ValueType_Error ValueType = iota
	ValueType_End
	ValueType_Object
	ValueType_Array
	ValueType_Int64
	ValueType_Bool
	ValueType_String
)

type Writer struct {
	wr io.Writer
}

type Value struct {
	Type  ValueType
	Depth int

	VInt64  int64
	VString string
	VBool   bool
}

func writeUint64(w *Writer, v uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	w.wr.Write(b)
}

func writeInt64(w *Writer, v int64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	w.wr.Write(b)
}

func writeInt32(w *Writer, v int32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	w.wr.Write(b)
}

func Write(w *Writer, v Value) {
	// write tag byte
	w.wr.Write([]byte{byte(v.Type)})
	switch v.Type {
	case ValueType_Int64:
		writeInt64(w, v.VInt64)
		break
	case ValueType_String:
		writeInt32(w, int32(len(v.VString)))
		w.wr.Write([]byte(v.VString))
		break
	case ValueType_Bool:
		var bv byte
		if v.VBool {
			bv = 1
		}
		w.wr.Write([]byte{bv})
		break
	}
}

func WriteInt(w *Writer, v int) {
	Write(w, Value{
		Type:   ValueType_Int64,
		VInt64: int64(v),
	})
}

func WriteInt64(w *Writer, v int64) {
	Write(w, Value{
		Type:   ValueType_Int64,
		VInt64: v,
	})
}

func WriteString(w *Writer, v string) {
	Write(w, Value{
		Type:    ValueType_String,
		VString: v,
	})
}

func WriteBool(w *Writer, v bool) {
	Write(w, Value{
		Type:  ValueType_Bool,
		VBool: v,
	})
}

func WriteObject(w *Writer) {
	Write(w, Value{
		Type: ValueType_Object,
	})
}

func WriteArray(w *Writer) {
	Write(w, Value{
		Type: ValueType_Array,
	})
}

func WriteEnd(w *Writer) {
	Write(w, Value{
		Type: ValueType_End,
	})
}

// ------------------ reader ------------------

type ReaderInterface interface {
	io.Reader
	io.ByteReader
}

type Reader struct {
	rd    ReaderInterface
	depth int
	len   int
	cur   int
}

func Read(r *Reader) (Value, error) {
	var v Value
	// read type byte
	var (
		typeB byte
		err   error
	)
	if typeB, err = r.rd.ReadByte(); err != nil {
		return v, err
	}
	v.Type = ValueType(typeB)
	switch v.Type {
	case ValueType_End:
		r.depth--
		break
	case ValueType_Object, ValueType_Array:
		r.depth++
		v.Depth = r.depth
		break
	case ValueType_Int64:
		if err := binary.Read(r.rd, binary.BigEndian, &v.VInt64); err != nil {
			return v, err
		}
		break
	case ValueType_String:
		var sLen int32
		if err := binary.Read(r.rd, binary.BigEndian, &sLen); err != nil {
			return v, err
		}
		b := make([]byte, sLen)
		if _, err := r.rd.Read(b); err != nil {
			return v, err
		}
		v.VString = string(b)
		break
	case ValueType_Bool:
		var (
			bv  byte
			err error
		)
		if bv, err = r.rd.ReadByte(); err != nil {
			return v, err
		}
		if bv == 1 {
			v.VBool = true
		} else {
			v.VBool = false
		}
		break
	default:
		return v, fmt.Errorf("unknown value type: %d", v.Type)
	}

	return v, nil
}

func discardUntilDepth(r *Reader, depth int) {
	for {
		if r.depth == depth {
			break
		}

		_, err := Read(r)
		if err != nil {
			break
		}
	}
}

func iterObject(r *Reader, obj, k, v *Value) (bool, error) {
	discardUntilDepth(r, obj.Depth)
	var err error
	*k, err = Read(r)
	if err != nil {
		return false, err
	}
	if k.Type == ValueType_End {
		return false, nil
	}
	*v, err = Read(r)
	if err != nil {
		return false, err
	}
	return true, nil
}

func iterArray(r *Reader, arr, av *Value) (bool, error) {
	discardUntilDepth(r, arr.Depth)
	var err error
	*av, err = Read(r)
	if err != nil {
		return false, err
	}
	if av.Type == ValueType_End {
		return false, nil
	}
	return true, nil
}

func readIntoObj(r *Reader, kvFn func(key string, v *Value)) error {
	var err error
	var obj, k, v Value
	if obj, err = Read(r); err != nil {
		return err
	}
	for {
		if ok, err := iterObject(r, &obj, &k, &v); err != nil {
			return err
		} else if !ok {
			break
		}
		// assert key is string
		if k.Type != ValueType_String {
			panic(fmt.Errorf("key is not string: %s\n", k.Type))
		}
		kvFn(k.VString, &v)
	}
	return nil
}

// ------------------ inspection ------------------

func PrintIndent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
}

func PrintValue(r *Reader, val Value, depth int) {
	var k, v Value
	var count int
	switch val.Type {
	case ValueType_Object:
		fmt.Print("{\n")
		for {
			if ok, err := iterObject(r, &val, &k, &v); err != nil {
				panic(err)
			} else if !ok {
				break
			}
			count++
			if count > 0 {
				fmt.Print(",\n")
			}
			PrintIndent(depth + 1)
			PrintValue(r, k, depth+1)
			fmt.Print(": ")
			PrintValue(r, v, depth+1)
		}
		if count > 0 {
			fmt.Print("\n")
		}
		PrintIndent(depth)
		fmt.Print("}")
		break
	case ValueType_Array:
		fmt.Print("[\n")
		for {
			if ok, err := iterArray(r, &val, &v); err != nil {
				panic(err)
			} else if !ok {
				break
			}
			count++
			if count > 0 {
				fmt.Print(",\n")
				PrintIndent(depth + 1)
				PrintValue(r, v, depth+1)
			}
		}
		if count > 0 {
			fmt.Print("\n")
		}
		PrintIndent(depth)
		fmt.Print("]")
		break
	case ValueType_Int64:
		fmt.Printf("%d", val.VInt64)
		break
	case ValueType_Bool:
		if val.VBool {
			fmt.Print("true")
		} else {
			fmt.Print("false")
		}
		break
	case ValueType_String:
		fmt.Printf("%q", val.VString)
		break
	default:
		fmt.Printf("ERROR: unknown value type: %d", val.Type)
	}
}

// ------------------ map ------------------

func WriteVal(w *Writer, v any) {
	switch v := v.(type) {
	case map[string]any:
		WriteObject(w)
		for k, v := range v {
			WriteString(w, k)
			WriteVal(w, v)
		}
		WriteEnd(w)
	case []any:
		WriteArray(w)
		for _, item := range v {
			WriteVal(w, item)
		}
		WriteEnd(w)
	case string:
		WriteString(w, v)
	case int:
		WriteInt(w, v)
	case bool:
		WriteBool(w, v)
	default:
		panic(fmt.Errorf("unsupported type %T", v))
	}
}

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
}

func (s *Sample) Marshal(w *Writer) error {
	WriteObject(w)
	WriteString(w, "name")
	WriteString(w, s.Name)
	WriteString(w, "age")
	WriteInt64(w, s.Age)
	WriteString(w, "isHuman")
	WriteBool(w, s.IsHuman)
	WriteEnd(w)
	return nil
}

func (n *Nested) Unmarshal(r *Reader) error {
	return readIntoObj(r, func(key string, v *Value) {
		switch key {
		case "inner":
			n.Inner = v.VInt64
		case "attrs":
			if v.Type != ValueType_Array {
				panic(fmt.Errorf("attrs is not array: %s\n", v.Type))
			}
			for {
				var av Value
				if ok, err := iterArray(r, v, &av); err != nil {
					panic(err)
				} else if !ok {
					break
				}
				n.Attrs = append(n.Attrs, av.VString)
			}
		}
	})
}

func (s *Sample) Unmarshal(r *Reader) error {
	return readIntoObj(r, func(key string, v *Value) {
		switch key {
		case "name":
			s.Name = v.VString
		case "age":
			s.Age = v.VInt64
		case "isHuman":
			s.IsHuman = v.VBool
		case "nested":
			var n Nested
			n.Unmarshal(r)
			s.Nested = n
		}
	})
}

func main2() {
	fmt.Println("simpser start")

	var bb bytes.Buffer
	s := &Sample{
		Name:    "testing",
		Age:     18,
		IsHuman: true,
	}
	w := &Writer{
		wr: &bb,
	}

	if err := s.Marshal(w); err != nil {
		panic("panic on marshalling: " + err.Error())
	}
	fmt.Println("marshalled: ", bb.Bytes())

	//r := &Reader{
	//	rd: &bb,
	//}
	//obj, err := Read(r)
	//if err != nil {
	//	panic("panic on reading: " + err.Error())
	//}
	//PrintValue(r, obj, 0)

	r2 := &Reader{
		rd: &bb,
	}
	s2 := &Sample{}
	if err := s2.Unmarshal(r2); err != nil {
		panic("panic on unmarshalling: " + err.Error())
	}
	fmt.Println("unmarshalled: ", s2)
}

func main() {
	fmt.Println("simpser map")
	var bb bytes.Buffer
	w := &Writer{
		wr: &bb,
	}
	m := map[string]any{
		"name":    "testing",
		"age":     18,
		"isHuman": true,
		"nested": map[string]any{
			"inner": 123,
			"attrs": []any{"a", "b", "c"},
		},
	}
	WriteVal(w, m)
	fmt.Println("marshalled: ", bb.Bytes())
	fmt.Println(string(bb.Bytes()))

	//r2 := &Reader{
	//	rd: &bb,
	//}
	//s2 := &Sample{}
	//if err := s2.Unmarshal(r2); err != nil {
	//	panic("panic on unmarshalling: " + err.Error())
	//}
	//fmt.Println("unmarshalled: ", s2)

	//fmt.Println("marshalled: ", bb.Bytes())
	//fmt.Println(string(bb.Bytes()))
}

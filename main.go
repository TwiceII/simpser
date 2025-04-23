package main

import (
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

type Reader struct {
	rd    io.Reader
	depth int
	len   int
	cur   int
}

func Read(r *Reader) (Value, error) {
	var v Value
	// read type byte
	if _, err := r.rd.Read([]byte{byte(v.Type)}); err != nil {
		return v, err
	}
	switch v.Type {
	case ValueType_End:
		r.depth--
		break
	case ValueType_Object:
	case ValueType_Array:
		r.depth++
		v.Depth = r.depth
		break
	case ValueType_Int64:
		if err := binary.Read(r.rd, binary.BigEndian, &v.VInt64); err != nil {
			return v, err
		}
		break
	case ValueType_String:
		var sLen int
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
		var bv byte
		if _, err := r.rd.Read([]byte{bv}); err != nil {
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
		v, err := Read(r)
		if err != nil {
			break
		}
		if v.Depth == depth {
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

// ------------------ main ------------------

type Sample struct {
	Name    string
	Age     int64
	IsHuman bool
}

func (s *Sample) Marshal(w *Writer) {
	WriteObject(w)
	WriteString(w, s.Name)
	WriteInt64(w, s.Age)
	WriteBool(w, s.IsHuman)
	WriteEnd(w)
}

func (s *Sample) Unmarshal(r *Reader) error {
	var err error
	var obj, k, v Value
	if _, err = Read(r); err != nil {
		return err
	}
	for {
		if ok, err := iterObject(r, &obj, &k, &v); err != nil {
			return err
		} else if !ok {
			break
		}
		switch k.Type {
		case ValueType_String:
			switch k.VString {
			case "name":
				s.Name = v.VString
				break
			case "age":
				s.Age = v.VInt64
				break
			case "is_human":
				s.IsHuman = v.VBool
				break
			default:
				return fmt.Errorf("unknown field: %s", k.VString)
			}
		default:
			return fmt.Errorf("unknown key type: %d", k.Type)
		}
	}
	return nil
}

func main() {
	fmt.Println("simpser start")
}

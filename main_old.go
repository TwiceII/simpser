package main

//
//import (
//	"bytes"
//	"encoding/binary"
//	"encoding/json"
//	"fmt"
//	"io"
//)
//
//type ValueType byte
//
//func (vt ValueType) String() string {
//	switch vt {
//	case ValueType_Error:
//		return "Error"
//	case ValueType_End:
//		return "End"
//	case ValueType_Object:
//		return "Object"
//	case ValueType_Array:
//		return "Array"
//	case ValueType_Int64:
//		return "Int64"
//	case ValueType_Bool:
//		return "Bool"
//	case ValueType_String:
//		return "String"
//	default:
//		return fmt.Sprintf("Unknown(%d)", vt)
//	}
//}
//
//const (
//	ValueType_Error ValueType = iota
//	ValueType_End
//	ValueType_Object
//	ValueType_Array
//	ValueType_Int64
//	ValueType_Bool
//	ValueType_String
//)
//
//type Value struct {
//	Type  ValueType
//	Depth int
//
//	VInt64  int64
//	VString string
//	VBool   bool
//}
//
//type WriterInterface interface {
//	io.Writer
//	io.ByteWriter
//}
//
//type Writer struct {
//	wr  WriterInterface
//	err error
//}
//
//type Marshalable interface {
//	Marshal(w *Writer)
//}
//
//func NewWriter(w WriterInterface) *Writer {
//	return &Writer{
//		wr: w,
//	}
//}
//
//func (w *Writer) Encode(v Marshalable) error {
//	if w.err != nil {
//		return w.err
//	}
//	v.Marshal(w)
//	return w.Err()
//}
//
//func (w *Writer) Err() error {
//	return w.err
//}
//
//func (w *Writer) write(b []byte) {
//	if w.err != nil {
//		return
//	}
//	_, w.err = w.wr.Write(b)
//}
//
//func (w *Writer) writeUint64(v uint64) {
//	var b [8]byte
//	binary.BigEndian.PutUint64(b[:], v)
//	w.write(b[:])
//}
//
//func (w *Writer) writeInt64(v int64) {
//	var b [8]byte
//	binary.BigEndian.PutUint64(b[:], uint64(v))
//	w.write(b[:])
//}
//
//func (w *Writer) writeInt32(v int32) {
//	var b [4]byte
//	binary.BigEndian.PutUint32(b[:], uint32(v))
//	w.write(b[:])
//}
//
//func (w *Writer) Write(v Value) {
//	// write tag byte
//	w.wr.WriteByte(byte(v.Type))
//	switch v.Type {
//	case ValueType_Int64:
//		w.writeInt64(v.VInt64)
//		break
//	case ValueType_String:
//		w.writeInt32(int32(len(v.VString)))
//		w.write([]byte(v.VString))
//		break
//	case ValueType_Bool:
//		var bv byte
//		if v.VBool {
//			bv = 1
//		}
//		w.wr.WriteByte(bv)
//		break
//	}
//}
//
//func (w *Writer) Int(v int) {
//	w.Write(Value{
//		Type:   ValueType_Int64,
//		VInt64: int64(v),
//	})
//}
//
//func (w *Writer) Int64(v int64) {
//	w.Write(Value{
//		Type:   ValueType_Int64,
//		VInt64: v,
//	})
//}
//
//func (w *Writer) String(v string) {
//	w.Write(Value{
//		Type:    ValueType_String,
//		VString: v,
//	})
//}
//
//func (w *Writer) Bool(v bool) {
//	w.Write(Value{
//		Type:  ValueType_Bool,
//		VBool: v,
//	})
//}
//
//func (w *Writer) Object() {
//	w.Write(Value{
//		Type: ValueType_Object,
//	})
//}
//
//func (w *Writer) Array() {
//	w.Write(Value{
//		Type: ValueType_Array,
//	})
//}
//
//func (w *Writer) End() {
//	w.Write(Value{
//		Type: ValueType_End,
//	})
//}
//
//type Obj struct {
//	w *Writer
//}
//
//func (w *Writer) Obj() *Obj {
//	w.Object()
//	return &Obj{w: w}
//}
//
////func (o *Obj) Nested() *Obj {
////	o.w.Object()
////	return o
////}
//
//func (o *Obj) String(k string, v string) *Obj {
//	o.w.String(k)
//	o.w.String(v)
//	return o
//}
//
//func (o *Obj) Int64(k string, v int64) *Obj {
//	o.w.String(k)
//	o.w.Int64(v)
//	return o
//}
//
//func (o *Obj) Bool(k string, v bool) *Obj {
//	o.w.String(k)
//	o.w.Bool(v)
//	return o
//}
//
//func (o *Obj) End() {
//	o.w.End()
//}
//
//// ------------------ reader ------------------
//
//type ReaderInterface interface {
//	io.Reader
//	io.ByteReader
//}
//
//type Unmarshalable interface {
//	Unmarshal(r *Reader, rv *Value)
//}
//
//type Reader struct {
//	rd    ReaderInterface
//	depth int
//	err   error
//}
//
//func NewReader(r ReaderInterface) *Reader {
//	return &Reader{
//		rd: r,
//	}
//}
//
//func (r *Reader) Err() error {
//	return r.err
//}
//
//func (r *Reader) Decode(v Unmarshalable) error {
//	if r.err != nil {
//		return r.err
//	}
//	rv := r.Read()
//	v.Unmarshal(r, &rv)
//	return r.Err()
//}
//
//func (r *Reader) Read() Value {
//	if r.err != nil {
//		return Value{} // default is ValueType_Error
//	}
//	var v Value
//	// read type byte
//	var (
//		typeB byte
//		err   error
//	)
//	if typeB, err = r.rd.ReadByte(); err != nil {
//		r.err = err
//		return v
//	}
//	v.Type = ValueType(typeB)
//	switch v.Type {
//	case ValueType_End:
//		r.depth--
//		break
//	case ValueType_Object, ValueType_Array:
//		r.depth++
//		v.Depth = r.depth
//		break
//	case ValueType_Int64:
//		if err := binary.Read(r.rd, binary.BigEndian, &v.VInt64); err != nil {
//			r.err = err
//			return v
//		}
//		break
//	case ValueType_String:
//		var sLen int32
//		if err := binary.Read(r.rd, binary.BigEndian, &sLen); err != nil {
//			r.err = err
//			return v
//		}
//		b := make([]byte, sLen)
//		if _, err := r.rd.Read(b); err != nil {
//			r.err = err
//			return v
//		}
//		v.VString = string(b)
//		break
//	case ValueType_Bool:
//		var (
//			bv  byte
//			err error
//		)
//		if bv, err = r.rd.ReadByte(); err != nil {
//			r.err = err
//			return v
//		}
//		if bv == 1 {
//			v.VBool = true
//		} else {
//			v.VBool = false
//		}
//		break
//	default:
//		r.err = fmt.Errorf("unknown value type: %d", v.Type)
//		return v
//	}
//
//	return v
//}
//
//func (r *Reader) discardUntilDepth(depth int) {
//	for {
//		if r.depth == depth {
//			break
//		}
//
//		r.Read()
//		if r.err != nil {
//			break
//		}
//	}
//}
//
//func (r *Reader) iterObject(obj, k, v *Value) bool {
//	r.discardUntilDepth(obj.Depth)
//	*k = r.Read()
//	if k.Type == ValueType_Error || k.Type == ValueType_End {
//		return false
//	}
//	*v = r.Read()
//	if v.Type == ValueType_Error {
//		return false
//	}
//
//	return true
//}
//
//func (r *Reader) iterArray(arr, av *Value) bool {
//	r.discardUntilDepth(arr.Depth)
//	*av = r.Read()
//	if av.Type == ValueType_Error || av.Type == ValueType_End {
//		return false
//	}
//	return true
//}
//
//func (r *Reader) IterateObject(obj *Value, kvFn func(key string, v *Value)) {
//	var k, v Value
//	for {
//		if ok := r.iterObject(obj, &k, &v); !ok {
//			break
//		}
//		// assert key is string
//		if k.Type != ValueType_String {
//			panic(fmt.Errorf("key is not string: %s\n", k.Type))
//		}
//		kvFn(k.VString, &v)
//	}
//}
//
//func (r *Reader) IterateArray(arr *Value, avFn func(v *Value)) {
//	var av Value
//	for {
//		if ok := r.iterArray(arr, &av); !ok {
//			break
//		}
//		avFn(&av)
//	}
//}
//
//// ------------------ inspection ------------------
//
//func printIndent(depth int) {
//	for i := 0; i < depth; i++ {
//		fmt.Print("  ")
//	}
//}
//
//func PrintReaderObject(r *Reader) {
//	obj := r.Read()
//	PrintValue(r, obj, 0)
//}
//
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
//		if val.VBool {
//			fmt.Print("true")
//		} else {
//			fmt.Print("false")
//		}
//		break
//	case ValueType_String:
//		fmt.Printf("%q", val.VString)
//		break
//	default:
//		fmt.Printf("ERROR: unknown value type: %d", val.Type)
//	}
//}
//
//// ------------------ map ------------------
//
//func WriteVal(w *Writer, v any) {
//	switch v := v.(type) {
//	case map[string]any:
//		w.Object()
//		for k, v := range v {
//			w.String(k)
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
//		w.String(v)
//	case int:
//		w.Int(v)
//	case bool:
//		w.Bool(v)
//	default:
//		panic(fmt.Errorf("unsupported type %T", v))
//	}
//}
//
//// ------------------ main ------------------
//
//type Nested struct {
//	Inner int64
//	Attrs []string
//}
//
//type Sample struct {
//	Name    string
//	Age     int64
//	IsHuman bool
//	Nested  Nested
//	LastVal int64
//}
//
//func (s *Sample) Marshal(w *Writer) {
//	// @@TODO: ergonomics
//	//w.Obj().
//	//	String("name", s.Name).
//	//	Int64("age", s.Age).
//	//	Bool("isHuman", s.IsHuman).
//	//	Int64("lastVal", s.LastVal).
//	//	End()
//
//	w.Object()
//	w.String("name")
//	w.String(s.Name)
//	w.String("age")
//	w.Int64(s.Age)
//	w.String("isHuman")
//	w.Bool(s.IsHuman)
//
//	// nested
//	w.String("nested")
//	w.Object() // < nested
//	w.String("inner")
//	w.Int64(s.Nested.Inner)
//	w.String("attrs")
//	w.Array() // < array
//	for _, attr := range s.Nested.Attrs {
//		w.String(attr)
//	}
//	w.End() // > array
//	w.End() // > nested
//
//	w.String("lastVal")
//	w.Int64(s.LastVal)
//
//	w.End()
//}
//
//func (n *Nested) Unmarshal(r *Reader, rv *Value) {
//	r.IterateObject(rv, func(key string, v *Value) {
//		switch key {
//		case "inner":
//			n.Inner = v.VInt64
//			break
//		case "attrs":
//			r.IterateArray(v, func(v *Value) {
//				n.Attrs = append(n.Attrs, v.VString)
//			})
//			break
//		}
//	})
//}
//
//func (s *Sample) Unmarshal(r *Reader, rv *Value) {
//	r.IterateObject(rv, func(key string, v *Value) {
//		switch key {
//		case "name":
//			s.Name = v.VString
//			break
//		case "age":
//			s.Age = v.VInt64
//			break
//		case "isHuman":
//			s.IsHuman = v.VBool
//			break
//		case "nested":
//			s.Nested.Unmarshal(r, v)
//			break
//		case "lastVal":
//			s.LastVal = v.VInt64
//			break
//		}
//	})
//}
//
//func main_2() {
//	fmt.Println("simpser start")
//
//	var bb bytes.Buffer
//	s := Sample{
//		Name:    "testing",
//		Age:     18,
//		IsHuman: true,
//		Nested: Nested{
//			Inner: 123,
//			Attrs: []string{"aa", "bb", "cc"},
//		},
//		LastVal: 777,
//	}
//	if err := NewWriter(&bb).Encode(&s); err != nil {
//		panic("panic on marshalling: " + err.Error())
//	}
//	fmt.Println("marshalled: ", bb.Bytes())
//
//	//PrintReaderObject(NewReader(&bb))
//
//	var s2 Sample
//	if err := NewReader(&bb).Decode(&s2); err != nil {
//		fmt.Println("error on unmarshalling: " + err.Error())
//	}
//	fmt.Println("unmarshalled: ", s2)
//}
//
//func main() {
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
//
//var sampleValue = Sample{
//	Name:    "testing",
//	Age:     18,
//	IsHuman: true,
//	Nested: Nested{
//		Inner: 123,
//		Attrs: []string{"aa", "bb", "cc"},
//	},
//	LastVal: 777,
//}
//
//func MarshallingSimpser() {
//	var bb bytes.Buffer
//	if err := NewWriter(&bb).Encode(&sampleValue); err != nil {
//		panic("panic on marshalling: " + err.Error())
//	}
//}
//
//func MarshallingJSON() {
//	var bb bytes.Buffer
//	if err := json.NewEncoder(&bb).Encode(&sampleValue); err != nil {
//		panic("panic on marshalling: " + err.Error())
//	}
//}

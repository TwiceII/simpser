package main

import (
	"sync"
	"unsafe"
)

type Writer struct {
	data []byte
}

func (w *Writer) Data() []byte {
	return w.data
}

func Marshal(v interface{}) ([]byte, error) {
	m := v.(Marshalable)
	w := NewWriter()
	m.Marshal(w)
	return w.Data(), nil
}

func (w *Writer) EncodeData(v Marshalable) []byte {
	v.Marshal(w)
	return w.data
}

type Marshalable interface {
	Marshal(w *Writer)
}

var buffersPool sync.Pool

func getBufferData() []byte {
	if v := buffersPool.Get(); v != nil {
		bd := v.([]byte)
		bd = bd[:0]
		return bd
	}
	return make([]byte, 0, 1024)
}

func NewWriter() *Writer {
	bd := getBufferData()
	defer buffersPool.Put(bd)

	return &Writer{
		data: bd,
		//data: make([]byte, 0, 128),
	}
}

func (w *Writer) Int(v int) {
	w.Int32(int32(v))
}

func (w *Writer) writeBytes(bs ...byte) {
	w.data = append(w.data, bs...)
}

func (w *Writer) EnsureSpace(n int) {
	var desiredSize = len(w.data) + n
	if len(w.data) < desiredSize {
		if cap(w.data) >= desiredSize {
			w.data = w.data[:desiredSize]
		} else {
			w.data = append(w.data, make([]byte, desiredSize-len(w.data))...)
		}
	}
}

func (w *Writer) writeUint32(v uint32) {
	w.writeBytes(byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (w *Writer) writeUint64(v uint64) {
	w.writeBytes(byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (w *Writer) Int32(v int32) {
	w.writeBytes(byte(ValueType_Int32))
	w.writeUint32(uint32(v))
}

func (w *Writer) Int64(v int64) {
	w.writeBytes(byte(ValueType_Int64))
	w.writeUint64(uint64(v))
}

func (w *Writer) Float64(v float64) {
	w.writeBytes(byte(ValueType_Float64))
	w.writeUint64(*(*uint64)(unsafe.Pointer(&v)))
}

func (w *Writer) String(v string) {
	w.writeBytes(byte(ValueType_String))
	size := len(v)
	w.writeUint32(uint32(size))
	var pos = len(w.data)
	w.EnsureSpace(size)
	copy(w.data[pos:pos+size], v)

	//w.writeString(v)
	//w.writeBytes(unsafe.Slice(unsafe.StringData(v), size)...)
}

//func (w *Writer) String(v string) {
//	w.writeBytes(byte(ValueType_String))
//	size := len(v)
//	w.writeUint32(uint32(size))
//	w.writeString(v)
//}

func (w *Writer) Bool(v bool) {
	w.writeBytes(byte(ValueType_Bool))
	if v {
		w.writeBytes(1)
	} else {
		w.writeBytes(0)
	}
}

func (w *Writer) Object() {
	w.writeBytes(byte(ValueType_Object))
}

func (w *Writer) Array() {
	w.writeBytes(byte(ValueType_Array))
}

func (w *Writer) End() {
	w.writeBytes(byte(ValueType_End))
}

//type Obj struct {
//	w *Writer
//}
//
//func (w *Writer) Obj() *Obj {
//	w.Object()
//	return &Obj{w: w}
//}

//func (o *Obj) Nested() *Obj {
//	o.w.Object()
//	return o
//}

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

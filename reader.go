package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type Unmarshalable interface {
	Unmarshal(r *Reader, rv *Value)
}

type Reader struct {
	data  []byte
	cur   int
	depth int
	err   error
}

func Unmarshal(data []byte, uv Unmarshalable) error {
	r := &Reader{
		data:  data,
		cur:   0,
		depth: 0,
	}
	v := r.Read()

	uv.Unmarshal(r, v)
	return r.err
}

func (r *Reader) readByte() byte {
	b := r.data[r.cur]
	r.cur++
	return b
}

func (r *Reader) readBytes(n int) []byte {
	b := r.data[r.cur : r.cur+n]
	r.cur += n
	return b
}

func (r *Reader) Read() *Value {
	var v Value
	// read type byte
	var (
		typeB byte
	)
	typeB = r.readByte()
	v.Type = ValueType(typeB)
	switch v.Type {
	case ValueType_End:
		r.depth--
		break
	case ValueType_Object, ValueType_Array:
		r.depth++
		v.Depth = r.depth
		break
	case ValueType_Int32:
		v.VInt32 = int32(binary.BigEndian.Uint32(r.readBytes(4)))
		break
	case ValueType_Int64, ValueType_Float64:
		v.VInt64 = int64(binary.BigEndian.Uint64(r.readBytes(8)))
		break
	case ValueType_String:
		sLen := int(binary.BigEndian.Uint32(r.readBytes(4)))
		v.VString = unsafe.String(&r.data[r.cur], sLen)
		r.cur += sLen
		break
	case ValueType_Bool:
		b := r.readByte()
		if b == 1 {
			v.VBool = true
		} else {
			v.VBool = false
		}
		break
	default:
		r.err = fmt.Errorf("unknown value type: %d", v.Type)
		return &v
	}

	return &v
}

func (r *Reader) discardUntilDepth(depth int) {
	for {
		if r.depth == depth {
			break
		}

		r.Read()
		if r.err != nil {
			break
		}
	}
}

func (r *Reader) iterObject(obj, k, v *Value) bool {
	r.discardUntilDepth(obj.Depth)
	*k = *r.Read()
	if k.Type == ValueType_Error || k.Type == ValueType_End {
		return false
	}
	*v = *r.Read()
	if v.Type == ValueType_Error {
		return false
	}

	return true
}

func (r *Reader) iterArray(arr, av *Value) bool {
	r.discardUntilDepth(arr.Depth)
	av = r.Read()
	if av.Type == ValueType_Error || av.Type == ValueType_End {
		return false
	}
	return true
}

func (r *Reader) IterateObject(obj *Value, kvFn func(key string, v *Value)) {
	var k, v Value
	for {
		if ok := r.iterObject(obj, &k, &v); !ok {
			break
		}
		// assert key is string
		if k.Type != ValueType_String {
			panic(fmt.Errorf("key is not string: %s\n", k.Type))
		}
		kvFn(k.VString, &v)
	}
}

func (r *Reader) IterateArray(arr *Value, avFn func(v *Value)) {
	var av Value
	for {
		if ok := r.iterArray(arr, &av); !ok {
			break
		}
		avFn(&av)
	}
}

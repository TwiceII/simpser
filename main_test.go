package main

import (
	"encoding/json"
	"testing"
)

//func BenchmarkMarshallingSimpser(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		MarshallingSimpser()
//	}
//}
//
//func BenchmarkMarshallingJSON(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		MarshallingJSON()
//	}
//}
//

func BenchmarkUnmarshalSimpser(b *testing.B) {
	//dt, err := Marshal(&sampleValue)
	//if err != nil {
	//	b.Fatal(err)
	//}
	//fmt.Println("dt::::")
	//fmt.Println(dt)
	var s2 SmallStruct
	for i := 0; i < b.N; i++ {
		Unmarshal(marshalledSampleValue, &s2)
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	//dt, err := Marshal(&sampleValue)
	//if err != nil {
	//	b.Fatal(err)
	//}
	var s2 SmallStruct
	for i := 0; i < b.N; i++ {
		json.Unmarshal(marshalledSampleValue, &s2)
	}
}

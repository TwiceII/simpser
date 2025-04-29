package main

import "testing"

func BenchmarkMarshallingSimpser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MarshallingSimpser()
	}
}

func BenchmarkMarshallingJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MarshallingJSON()
	}
}

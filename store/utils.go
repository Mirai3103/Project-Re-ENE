package store

import (
	"encoding/binary"
	"math"
)

func Float32ToBytes(f32 []float32) []byte {
	b := make([]byte, len(f32)*4)
	for i, v := range f32 {
		bits := math.Float32bits(v)
		binary.LittleEndian.PutUint32(b[i*4:], bits)
	}
	return b
}

func BytesToFloat32(b []byte) []float32 {
	if len(b)%4 != 0 {
		panic("invalid float32 blob: length must be divisible by 4")
	}
	n := len(b) / 4
	f32 := make([]float32, n)
	for i := 0; i < n; i++ {
		bits := binary.LittleEndian.Uint32(b[i*4:])
		f32[i] = math.Float32frombits(bits)
	}
	return f32
}

func L2DistanceF32(a, b []float32) float64 {
	if len(a) != len(b) {
		panic("dimension mismatch")
	}
	var sum float64
	for i := range a {
		d := float64(a[i] - b[i])
		sum += d * d
	}
	return math.Sqrt(sum)
}

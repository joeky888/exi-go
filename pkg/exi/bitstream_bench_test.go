package exi

import "testing"

// BenchmarkBitStreamReadBits benchmarks the original ReadBits implementation
func BenchmarkBitStreamReadBits(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.Run("Original_1bit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			bs.Init(data, 0)
			for j := 0; j < 100; j++ {
				_, _ = bs.ReadBits(1)
			}
		}
	})

	b.Run("Optimized_1bit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			bs.Init(data, 0)
			for j := 0; j < 100; j++ {
				_, _ = bs.ReadBitsOptimized(1)
			}
		}
	})

	b.Run("Original_6bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			bs.Init(data, 0)
			for j := 0; j < 100; j++ {
				_, _ = bs.ReadBits(6)
			}
		}
	})

	b.Run("Optimized_6bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			bs.Init(data, 0)
			for j := 0; j < 100; j++ {
				_, _ = bs.ReadBitsOptimized(6)
			}
		}
	})

	b.Run("Original_16bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			bs.Init(data, 0)
			for j := 0; j < 100; j++ {
				_, _ = bs.ReadBits(16)
			}
		}
	})

	b.Run("Optimized_16bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			bs.Init(data, 0)
			for j := 0; j < 100; j++ {
				_, _ = bs.ReadBitsOptimized(16)
			}
		}
	})
}

// BenchmarkBitStreamWriteBits benchmarks write operations
func BenchmarkBitStreamWriteBits(b *testing.B) {
	b.Run("Original_6bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			buf := make([]byte, 1024)
			bs.Init(buf, 0)
			for j := 0; j < 100; j++ {
				_ = bs.WriteBits(6, 0x3F)
			}
		}
	})

	b.Run("Optimized_6bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			buf := make([]byte, 1024)
			bs.Init(buf, 0)
			for j := 0; j < 100; j++ {
				_ = bs.WriteBitsOptimized(6, 0x3F)
			}
		}
	})

	b.Run("Original_16bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			buf := make([]byte, 1024)
			bs.Init(buf, 0)
			for j := 0; j < 100; j++ {
				_ = bs.WriteBits(16, 0xFFFF)
			}
		}
	})

	b.Run("Optimized_16bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := &BitStream{}
			buf := make([]byte, 1024)
			bs.Init(buf, 0)
			for j := 0; j < 100; j++ {
				_ = bs.WriteBitsOptimized(16, 0xFFFF)
			}
		}
	})
}

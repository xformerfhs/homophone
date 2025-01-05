//
// SPDX-FileCopyrightText: Copyright 2023-2025 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.1.0
//
// Change history:
//    2023-01-21: V1.0.0: Created.
//    2025-01-05: V1.1.0: use math/rand/v2.
//

package compressedinteger

import (
	"bytes"
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Known conversion values
// Each integer value corresponds to the byte slice value with the same index

var intValues = []int{0, 0x3f,
	0x40, 0x403f,
	0x4040, 0x40403f,
	0x404040, 0x4040403f}

var byteSliceValues = [][]byte{{0}, {0x3f},
	{0x40, 0x00}, {0x7f, 0xff},
	{0x80, 0x00, 0x00}, {0xbf, 0xff, 0xff},
	{0xc0, 0x00, 0x00, 0x00}, {0xff, 0xff, 0xff, 0xff}}

// ******** Test functions ********

// TestBoundariesFromInt tests all integer boundary cases.
func TestBoundariesFromInt(t *testing.T) {
	var err error
	var i int
	var n int
	var c []byte

	for i, n = range intValues {
		c, err = FromInt(n)
		if err != nil {
			t.Fatalf("error converting integer %d: %v", n, err)
		}

		if bytes.Compare(c, byteSliceValues[i]) != 0 {
			t.Fatalf("conversion of %d (0x%x) resulted in 0x%x instead of 0x%x", n, n, c, byteSliceValues[i])
		}
	}
}

// TestBoundariesFromBytes tests all byte slice boundary cases.
func TestBoundariesFromBytes(t *testing.T) {
	var err error
	var i int
	var n int
	var c []byte

	for i, c = range byteSliceValues {
		n, _, err = ToInt(c)
		if err != nil {
			t.Fatalf("error converting byte slice 0x%x: %v", c, err)
		}

		if n != intValues[i] {
			t.Fatalf("conversion of 0x%x resulted in %d instead of %d", c, n, intValues[i])
		}
	}
}

// TestRandomConversion tests random integer conversions in all byte slice ranges.
func TestRandomIntConversion(t *testing.T) {
	supValues := []int32{0x40, 0x4040, 0x404040, 0x40404040}
	var err error
	var n int
	var c []byte
	var rn int

	for _, supValue := range supValues {
		for i := 0; i < 100; i++ {
			// Convert random integer to byte slice
			n = int(rand.Int32N(supValue))
			c, err = FromInt(n)
			if err != nil {
				t.Fatalf("error converting integer %d: %v", n, err)
			}

			// Convert the resulting byte slice back to integer
			rn, _, err = ToInt(c)
			if err != nil {
				t.Fatalf("error converting byte slice 0x%x: %v", c, err)
			}

			// Test if the re-converted integer has the same value as the original integer
			if n != rn {
				t.Fatalf("converting %d (0x%x) gave 0x%x and when converted back resulted in %d (0x%x)", n, n, c, rn, rn)
			}
		}
	}
}

// TestRandomUInt32Conversion tests random uint32 conversions in all byte slice ranges.
func TestRandomUInt32Conversion(t *testing.T) {
	supValues := []int32{0x40, 0x4040, 0x404040, 0x40404040}
	var err error
	var n uint32
	var c []byte
	var rn uint32

	for _, supValue := range supValues {
		for i := 0; i < 100; i++ {
			// Convert random integer to byte slice
			n = uint32(rand.Int32N(supValue))
			c, err = FromUInt32(n)
			if err != nil {
				t.Fatalf("error converting integer %d: %v", n, err)
			}

			// Convert the resulting byte slice back to integer
			rn, _, err = ToUInt32(c)
			if err != nil {
				t.Fatalf("error converting byte slice 0x%x: %v", c, err)
			}

			// Test if the re-converted integer has the same value as the original integer
			if n != rn {
				t.Fatalf("converting %d (0x%x) gave 0x%x and when converted back resulted in %d (0x%x)", n, n, c, rn, rn)
			}
		}
	}
}

// TestInvalidIntegers tests if the correct errors are returned for invalid integer values
func TestInvalidIntegers(t *testing.T) {
	expectIntegerConversionError(t, math.MinInt32, ErrIntIsNegative)
	expectIntegerConversionError(t, -5, ErrIntIsNegative)
	expectIntegerConversionError(t, MinAllowedInt-1, ErrIntIsNegative)
	expectIntegerConversionError(t, MaxAllowedInt+1, ErrIntIsTooLarge)
	expectIntegerConversionError(t, math.MaxInt32, ErrIntIsTooLarge)
}

// TestInvalidBytes tests if the correct errors are returned for invalid byte slices
func TestInvalidBytes(t *testing.T) {
	expectByteSliceConversionError(t, nil, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0x40}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0x7f}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0x80}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0x80, 0x00}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xbf}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xbf, 0x00}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xc0}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xc0, 0x00}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xc0, 0x00, 0x00}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xff}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xff, 0x00}, ErrSliceTooSmall)
	expectByteSliceConversionError(t, []byte{0xff, 0x00, 0x00}, ErrSliceTooSmall)
}

// ******** Benchmark functions ********

// BenchmarkFrom1ByteInteger benchmarks conversions for integers that generate 1 byte compressed representations
func BenchmarkFrom1ByteInteger(b *testing.B) {
	i := intValues[1]
	for n := 0; n < b.N; n++ {
		_, _ = FromInt(i)
	}
}

// BenchmarkFrom2ByteInteger benchmarks conversions for integers that generate 2 bytes compressed representations
func BenchmarkFrom2ByteInteger(b *testing.B) {
	i := intValues[3]
	for n := 0; n < b.N; n++ {
		_, _ = FromInt(i)
	}
}

// BenchmarkFrom3ByteInteger benchmarks conversions for integers that generate 3 bytes compressed representations
func BenchmarkFrom3ByteInteger(b *testing.B) {
	i := intValues[5]
	for n := 0; n < b.N; n++ {
		_, _ = FromInt(i)
	}
}

// BenchmarkFrom4ByteInteger benchmarks conversions for integers that generate 4 bytes compressed representations
func BenchmarkFrom4ByteInteger(b *testing.B) {
	i := intValues[7]
	for n := 0; n < b.N; n++ {
		_, _ = FromInt(i)
	}
}

// BenchmarkFrom1ByteSlice benchmarks conversions from 1 byte compressed representations
func BenchmarkFrom1ByteSlice(b *testing.B) {
	c := byteSliceValues[1]
	for n := 0; n < b.N; n++ {
		_, _, _ = ToInt(c)
	}
}

// BenchmarkFrom2ByteSlice benchmarks conversions from 2 bytes compressed representations
func BenchmarkFrom2ByteSlice(b *testing.B) {
	c := byteSliceValues[3]
	for n := 0; n < b.N; n++ {
		_, _, _ = ToInt(c)
	}
}

// BenchmarkFrom3ByteSlice benchmarks conversions from 3 bytes compressed representations
func BenchmarkFrom3ByteSlice(b *testing.B) {
	c := byteSliceValues[5]
	for n := 0; n < b.N; n++ {
		_, _, _ = ToInt(c)
	}
}

// BenchmarkFrom4ByteSlice benchmarks conversions from 4 bytes compressed representations
func BenchmarkFrom4ByteSlice(b *testing.B) {
	c := byteSliceValues[7]
	for n := 0; n < b.N; n++ {
		_, _, _ = ToInt(c)
	}
}

// ******** Private methods ********

// expectIntegerConversionError tests if a given integer conversion return the expected error
func expectIntegerConversionError(t *testing.T, n int, expectedError error) {
	_, err := FromInt(n)
	if err == nil {
		t.Fatalf("no error converting integer %d, expected %v", n, expectedError)
	} else {
		if !errors.Is(err, expectedError) {
			t.Fatalf("converting integer %d resulted in error %v instead of %v", n, err, expectedError)
		}
	}
}

// expectByteSliceConversionError tests if a given byte slice conversion return the expected error
func expectByteSliceConversionError(t *testing.T, c []byte, expectedError error) {
	_, _, err := ToInt(c)
	if err == nil {
		t.Fatalf("no error converting bytes 0x%x, expected %v", c, expectedError)
	} else {
		if !errors.Is(err, expectedError) {
			t.Fatalf("converting bytes 0x%x resulted in error %v instead of %v", c, err, expectedError)
		}
	}
}

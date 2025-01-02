//
// SPDX-FileCopyrightText: Copyright 2023 Frank Schwab
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
// Version: 1.0.0
//
// Change history:
//    2023-03-25: V1.0.0: Created.
//

package compressedinteger

import (
	"errors"
)

// ======== Public constants ========

// Minimum and maximum allowed value

const MinAllowedInt = 0
const MaxAllowedInt = 0x40404040 - 1

// List of errors

var ErrIntIsNegative = errors.New(`integer is negative`)
var ErrIntIsTooLarge = errors.New(`integer is too large`)
var ErrSliceTooSmall = errors.New(`byte slice has not enough bytes for compressed integer`)

// ======== Private constants ========

// Constant for conversion

const offsetValue = 0x40

// Constants for masks

const noLengthMaskForByte = offsetValue - 1
const byteMaskForInteger = 0xff

// Slice manipulation constants

const resultSliceLength = 4
const resultMaxIndex = resultSliceLength - 1
const lengthBitsShiftValue = 6

// ======== Private global variables ========

// result is the global result buffer (i.e. all returned byte slices point here)
var result = make([]byte, resultSliceLength)

// ======== Public functions ========

// FromUInt32 converts an uint32 to a compressed representation byte slice
func FromUInt32(i uint32) ([]byte, error) {
	// Check if value is within allowed range
	if i > MaxAllowedInt {
		return nil, ErrIntIsTooLarge
	}

	// Convert to byte array

	// This loop subtracts the offset from each byte
	temp := i
	actIndex := resultMaxIndex
	for temp >= offsetValue {
		b := temp & byteMaskForInteger
		temp >>= 8
		if b >= offsetValue {
			b -= offsetValue
		} else {
			b += 256 - offsetValue
			temp--
		}

		result[actIndex] = byte(b)
		actIndex--
	}

	// Add length flag
	result[actIndex] = byte(temp | (uint32(resultMaxIndex-actIndex) << lengthBitsShiftValue))

	return result[actIndex:resultSliceLength], nil
}

// FromInt converts an int to a compressed representation byte slice
func FromInt(i int) ([]byte, error) {
	// Check if value is within allowed range
	if i < MinAllowedInt {
		return nil, ErrIntIsNegative
	}

	return FromUInt32(uint32(i))
}

// ToUInt32 converts a compressed representation into an uint32
func ToUInt32(p []byte) (uint32, int, error) {
	pLen := len(p)
	if pLen == 0 {
		return 0, 0, ErrSliceTooSmall
	}

	expectedLength := ExpectedLength(p[0]) // Calculate expected length

	if pLen < expectedLength {
		return 0, 0, ErrSliceTooSmall
	}

	// Decompress the byte array
	t := uint32(p[0] & noLengthMaskForByte)

	for i := 1; i < expectedLength; i++ {
		t = ((t << 8) | uint32(p[i])) + offsetValue
	}

	return t, expectedLength, nil
}

// ToInt converts a compressed representation into an int
func ToInt(p []byte) (int, int, error) {
	i, l, err := ToUInt32(p)
	return int(i), l, err
}

func ExpectedLength(b byte) int {
	return int(b>>lengthBitsShiftValue) + 1
}

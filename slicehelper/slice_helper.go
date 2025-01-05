//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
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
// Version: 1.5.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-03-17: V1.1.0: Add FillToCap.
//    2024-04-26: V1.2.0: Add SimpleFill.
//    2024-06-03: V1.3.0: Add CutTail.
//    2024-08-09: V1.4.0: Rename CutTail to SplitTail.
//    2025-01-05: V1.5.0: Add RemoveNoOrder, remove SimpleFill and SplitTail.
//

// Package slicehelper implements helper functions for slices.
package slicehelper

import (
	"homophone/constraints"
)

// ******** Private constants ********

// powerFillThresholdLen is the slice length where PowerFill is more efficient than SimpleFill.
const powerFillThresholdLen = 74

// ******** Public functions ********

// FillToCap fills a slice with a value in an efficient way up to its capacity.
func FillToCap[S ~[]T, T any](s S, v T) {
	sLen := cap(s)

	if sLen > 0 {
		doFill(s[:sLen], v, sLen)
	}
}

// ClearNumber clears a number type slice.
func ClearNumber[S ~[]T, T constraints.Number](a S) {
	FillToCap(a, 0)
}

// RemoveNoOrder removes the element with the specified index.
// The order of the elements is *not* preserved.
func RemoveNoOrder[S ~[]T, T any](s S, i int) S {
	lastIndex := len(s)

	// Just return the slice if it contains no elements.
	if lastIndex == 0 {
		return s
	}

	// Copy last element into removed element.
	lastIndex--
	if i != lastIndex {
		s[i] = s[lastIndex]
	}

	// Return a slice that is shortened by 1.
	return s[:lastIndex]
}

// ******** Private functions ********

// doFill fills a slice in an optimal way.
func doFill[S ~[]T, T any](s S, v T, l int) {
	if l >= powerFillThresholdLen {
		doPowerFill(s, v, l)
	} else {
		doSimpleFill(s, v, l)
	}
}

// doSimpleFill fills a slice in a simple way.
func doSimpleFill[S ~[]T, T any](s S, v T, l int) {
	for i := 0; i < l; i++ {
		s[i] = v
	}
}

// doPowerFill fills a slice in an efficient way.
func doPowerFill[S ~[]T, T any](s S, v T, l int) {
	// Put the value into the first slice element
	s[0] = v

	// Incrementally duplicate the value into the rest of the slice
	for j := 1; j < l; j <<= 1 {
		copy(s[j:], s[:j])
	}
}

//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
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
//    2025-02-09: V1.0.0: Created.
//

// Package equalshandler makes handling equal values easier.
package equalshandler

import (
	"cmp"
	"slices"
)

// ******** Public types ********

// Handler contains the data needed to manage equal counts.
type Handler[T cmp.Ordered] struct {
	valueToIndexMap map[T][]int
}

// ******** Creation methods ********

// NewFromValues creates a new equalshandler from the supplied counts slice.
func NewFromValues[T cmp.Ordered](values []T) *Handler[T] {
	valueToIndexMap := buildCountToIndexMap(values)

	removeSingleEntries[T](valueToIndexMap)

	return &Handler[T]{valueToIndexMap: valueToIndexMap}
}

// ******** Public type methods ********

// SortedLengths yields a slice with the lengths that have been found and that are
// greater than one. The list is sorted in descending order.
func (h *Handler[T]) SortedLengths() []T {
	result := make([]T, 0, len(h.valueToIndexMap))
	for k := range h.valueToIndexMap {
		result = append(result, k)
	}

	slices.Sort(result)
	slices.Reverse(result)

	return result
}

// Length returns the number of entries in this handler.
func (h *Handler[T]) Length() int {
	return len(h.valueToIndexMap)
}

// Entries returns a slice with the indices of the entries with this value.
func (h *Handler[T]) Entries(value T) []int {
	return h.valueToIndexMap[value]
}

// ******** Private functions ********

// buildCountToIndexMap builds a map that maps from counts
// to indices with this count.
func buildCountToIndexMap[T cmp.Ordered](values []T) map[T][]int {
	valueToIndexMap := make(map[T][]int)

	for i, value := range values {
		valueToIndexMap[value] = append(valueToIndexMap[value], i)
	}

	return valueToIndexMap
}

// removeSingleEntries removes all entries that only have a count of 1.
func removeSingleEntries[T cmp.Ordered](valueToIndexMap map[T][]int) {
	deleteList := make([]T, 0, len(valueToIndexMap))
	for k, v := range valueToIndexMap {
		if len(v) == 1 {
			deleteList = append(deleteList, k)
		}
	}

	for _, k := range deleteList {
		delete(valueToIndexMap, k)
	}
}

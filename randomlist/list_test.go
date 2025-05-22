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
//    2025-01-02: V1.0.0: Created.
//

package randomlist

import "testing"

func TestZeroLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(`Accessing an element of a zero length list did not panic`)
		}
	}()

	x := New[byte](nil)
	if x.Len() != 0 {
		t.Fatal(`Zero elements slice does not have zero length`)
	}

	x.RandomElement()
}

func TestInt(t *testing.T) {
	for testLen := 1; testLen <= 100; testLen++ {
		testSlice := make([]int, testLen)
		for i := 0; i < testLen; i++ {
			testSlice[i] = i
		}
		testRandomList := New(testSlice)
		count := 0
		for i := 0; i < testLen; i++ {
			n := testRandomList.RandomElement()
			if n == i {
				count++
			}
		}
		if testLen > 4 && count == testLen {
			t.Fatal(`Random list does not have any randomness`)
		}
	}
}

func TestFloat64(t *testing.T) {
	for testLen := 1; testLen <= 100; testLen++ {
		testSlice := make([]float64, testLen)
		for i := 0; i < testLen; i++ {
			testSlice[i] = float64(i)
		}
		testRandomList := New(testSlice)
		count := 0
		for i := 0; i < testLen; i++ {
			n := testRandomList.RandomElement()
			if n == float64(i) {
				count++
			}
		}
		if testLen > 4 && count == testLen {
			t.Fatal(`Random list does not have any randomness`)
		}
	}
}

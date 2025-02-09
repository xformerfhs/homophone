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

package distributor

import (
	"fmt"
	"homophone/equalshandler"
	"homophone/slicehelper"
	"math/rand/v2"
)

// randomAdjustment makes random adjustments to the seats until the total count matches the wanted count.
func randomAdjustment(
	seats []uint,
	distributedSeatCount uint,
	wantedSeatCount uint) {
	equalsList := equalshandler.NewFromValues(seats)
	candidateLengths := equalsList.SortedLengths()
	diffCount := distributedSeatCount - wantedSeatCount
	actCandidate := 0

	for actCandidate < len(candidateLengths) {
		actIndices := equalsList.Entries(candidateLengths[actCandidate])
		actIndicesLen := len(actIndices)

		for diffCount != 0 && actIndicesLen != 0 {
			i := rand.IntN(actIndicesLen)
			si := actIndices[i]
			s := seats[si]

			if diffCount > 0 {
				s--
				diffCount--
			} else {
				s++
				diffCount++
			}

			seats[si] = s

			if diffCount != 0 {
				actIndices = slicehelper.RemoveNoOrder(actIndices, i)
				actIndicesLen--
			}
		}

		if diffCount == 0 {
			break
		}

		actCandidate++
	}

	if diffCount != 0 {
		panic(fmt.Sprintf(`unable to find a matching distribution (diff=%d)`, diffCount))
	}
}

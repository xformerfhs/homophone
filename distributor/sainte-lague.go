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
// Version: 2.0.0
//
// Change history:
//    2025-02-08: V1.0.0: Created.
//    2025-02-09: V2.0.0: Use generic interface.
//

// Package distributor contains functions to distribute counts to a number of seats.
package distributor

import (
	"homophone/constraints"
	"math"
)

// ******** Public functions ********

// SainteLagueDistribution implements the Sainte-Laguë method for distributing a number of counts
// to a number of seats.
func SainteLagueDistribution[T constraints.Integer](counts []T, totalCount uint, wantedSeatCount uint) []uint {
	divisor := float64(totalCount) / float64(wantedSeatCount)
	lastDivisor := 0.0
	intSeats := make([]uint, len(counts))
	floatSeats := make([]float64, len(counts))
	distributedSeatCount := uint(0)

	for {
		distributedSeatCount = 0

		for i, count := range counts {
			floatSeatsCount := float64(count) / divisor
			floatSeats[i] = floatSeatsCount
			intSeatsCount := uint(math.Round(floatSeatsCount))
			intSeats[i] = intSeatsCount
			distributedSeatCount += intSeatsCount
		}

		if distributedSeatCount != wantedSeatCount {
			if nearlyEqual(lastDivisor, divisor) {
				// Unable to find a distribution because of too many equal counts.
				// Make random adjustments. This is the only way to fix this.
				randomAdjustment(intSeats, distributedSeatCount, wantedSeatCount)
				break
			} else {
				lastDivisor = divisor
				divisor = nextDivisor(counts, floatSeats, distributedSeatCount < wantedSeatCount)
			}
		} else {
			break
		}
	}

	return intSeats
}

// ******** Private functions ********

// nextDivisor calculates the next number of seats in the correct direction.
func nextDivisor[T constraints.Integer](counts []T, floatSeats []float64, needMore bool) float64 {
	if needMore {
		return nextDivisorUp(counts, floatSeats)
	} else {
		return nextDivisorDown(counts, floatSeats)
	}
}

// nextDivisorUp calculates the next divisor when more seats are needed.
func nextDivisorUp[T constraints.Integer](counts []T, floatSeats []float64) float64 {
	maxSeatBoundary := 0.0
	secondSeatBoundary := 0.0
	for i, floatSeat := range floatSeats {
		if floatSeat > 0.0 {
			seatBoundary := float64(counts[i]) / nextHalfSeatUp(floatSeat)

			if seatBoundary > maxSeatBoundary {
				secondSeatBoundary = maxSeatBoundary
				maxSeatBoundary = seatBoundary
			} else {
				if seatBoundary > secondSeatBoundary {
					secondSeatBoundary = seatBoundary
				}
			}
		}
	}

	if secondSeatBoundary != 0.0 {
		return (maxSeatBoundary + secondSeatBoundary) * 0.5
	} else {
		return maxSeatBoundary
	}
}

// nextDivisorDown calculates the next divisor when fewer seats are needed.
func nextDivisorDown[T constraints.Integer](counts []T, floatSeats []float64) float64 {
	minSeatBoundary := math.MaxFloat64
	secondSeatBoundary := math.MaxFloat64
	for i, floatSeat := range floatSeats {
		if floatSeat > 0.0 {
			seatBoundary := float64(counts[i]) / nextHalfSeatDown(floatSeat)

			if seatBoundary < minSeatBoundary {
				secondSeatBoundary = minSeatBoundary
				minSeatBoundary = seatBoundary
			} else {
				if seatBoundary < secondSeatBoundary {
					secondSeatBoundary = seatBoundary
				}
			}
		}
	}

	if secondSeatBoundary != math.MaxFloat64 {
		return (minSeatBoundary + secondSeatBoundary) * 0.5
	} else {
		return minSeatBoundary
	}
}

// nextHalfSeatUp gets the next number of seats rounded up to the next
// half between two integers.
func nextHalfSeatUp(floatSeat float64) float64 {
	doubleSeat := int(math.Ceil(floatSeat + floatSeat))
	if doubleSeat&1 == 0 {
		doubleSeat++
	}

	return float64(doubleSeat) * 0.5
}

// nextHalfSeatDown gets the next number of seats rounded down to the next
// half between two integers.
func nextHalfSeatDown(floatSeat float64) float64 {
	doubleSeat := int(math.Floor(floatSeat + floatSeat))
	if doubleSeat&1 == 0 {
		doubleSeat--
	}

	return float64(doubleSeat) * 0.5
}

// nearlyEqual compares two float64 numbers if they are equal or nearly equal.
// Floating point numbers should not be tested for equality, because different
// calculation methods yield slightly different values for the same result.
func nearlyEqual(a, b float64) bool {
	if a == b {
		return true
	} else {
		return math.Abs(a-b) < 0.000001
	}
}

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
//    2025-02-08: V1.0.0: Created.
//

// Package distributor contains functions to distribute the counts to
// a number of "seats".
package distributor

import (
	"math"
)

// ******** Public functions ********

// SainteLagueDistribution implements the Sainte-Laguë method for distributing a number of counts
// to a number of seats. There are cases when this method does not find a matching distribution.
// In this case a distribution that is nearly correct is returned.
// The caller can check this with the second returned variable that contains the number of seats
// that have been distributed.
func SainteLagueDistribution(counts []int, totalCount int, wantedSeatCount uint16) ([]uint16, uint16) {
	divisor := float64(totalCount) / float64(wantedSeatCount)
	intSeats := make([]uint16, len(counts))
	floatSeats := make([]float64, len(counts))
	distributedSeatCount := uint16(0)
	lastDivisor := 0.0

	for {
		distributedSeatCount = 0

		for i, count := range counts {
			floatSeatsCount := float64(count) / divisor
			floatSeats[i] = floatSeatsCount
			intSeatsCount := uint16(math.Round(floatSeatsCount))
			intSeats[i] = intSeatsCount
			distributedSeatCount += intSeatsCount
		}

		if distributedSeatCount != wantedSeatCount {
			if nearlyEqual(lastDivisor, divisor) {
				break
			} else {
				lastDivisor = divisor
				divisor = nextDivisor(counts, floatSeats, distributedSeatCount < wantedSeatCount)
			}
		} else {
			break
		}
	}

	return intSeats, distributedSeatCount
}

// ******** Private functions ********

// nextDivisor calculates the next number of seats in the correct direction.
func nextDivisor(counts []int, floatSeats []float64, needMore bool) float64 {
	if needMore {
		return nextDivisorUp(counts, floatSeats)
	} else {
		return nextDivisorDown(counts, floatSeats)
	}
}

// nextDivisorUp calculates the next divisor when more seats are needed.
func nextDivisorUp(counts []int, floatSeats []float64) float64 {
	maxSeatBoundary := 0.0
	secondSeatBoundary := 0.0
	for i, floatSeat := range floatSeats {
		if floatSeat != 0.0 {
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

// nextDivisorDown calculates the next divisor when less seats are needed.
func nextDivisorDown(counts []int, floatSeats []float64) float64 {
	minSeatBoundary := math.MaxFloat64
	secondSeatBoundary := math.MaxFloat64
	for i, floatSeat := range floatSeats {
		if floatSeat != 0.0 {
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

// nextHalfSeatDown gets the next number of seats rounded up to the next
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

// nearlyEqual compares two float64 numbers if they are nearer than 0.000001.
func nearlyEqual(a float64, b float64) bool {
	if a == b {
		return true
	} else {
		return math.Abs(a-b) < 0.000001
	}
}

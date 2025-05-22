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
// Version: 2.1.1
//
// Change history:
//    2025-02-08: V1.0.0: Created.
//    2025-02-09: V2.0.0: Use generic interface.
//    2025-02-10: V2.1.0: Fixed cut off criteria.
//    2025-02-24: V2.1.1: Reorder diff test for better efficiency.
//

// Package distributor contains functions to distribute counts to
// a number of seats.
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
	intSeats := make([]uint, len(counts))
	floatSeats := make([]float64, len(counts))
	distributedSeatCount := uint(0)
	minSeatCountDiff := uint(math.MaxInt32)
	minDivisor := math.MaxFloat64
	lastSeatCountDiff := uint(0)

	for {
		distributedSeatCount = seatsForDivisor(counts, divisor, floatSeats, intSeats)

		actSeatCountDiff := absUintDiff(distributedSeatCount, wantedSeatCount)

		// If the seat count difference is zero we are done.
		if actSeatCountDiff == 0 {
			break
		}

		// Otherwise, remember when the smallest difference occurred.
		if actSeatCountDiff < minSeatCountDiff {
			minDivisor = divisor
			minSeatCountDiff = actSeatCountDiff
		}

		if lastSeatCountDiff != 0 && actSeatCountDiff >= lastSeatCountDiff {
			if distributedSeatCount > minSeatCountDiff {
				distributedSeatCount = seatsForDivisor(counts, minDivisor, floatSeats, intSeats)
			}

			// Unable to find a distribution because of too many equal counts.
			// Make random adjustments. This is the only way to fix this.
			randomAdjustment(intSeats, distributedSeatCount, wantedSeatCount)
			break
		} else {
			lastSeatCountDiff = actSeatCountDiff
			divisor = nextDivisor(counts, floatSeats, distributedSeatCount < wantedSeatCount)
		}
	}

	return intSeats
}

// ******** Private functions ********

func seatsForDivisor[T constraints.Integer](counts []T, divisor float64, floatSeats []float64, intSeats []uint) uint {
	distributedSeatCount := uint(0)

	for i, count := range counts {
		floatSeatsCount := float64(count) / divisor
		floatSeats[i] = floatSeatsCount
		intSeatsCount := uint(math.Round(floatSeatsCount))
		intSeats[i] = intSeatsCount
		distributedSeatCount += intSeatsCount
	}

	return distributedSeatCount
}

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
				if seatBoundary != maxSeatBoundary && seatBoundary > secondSeatBoundary {
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
		if floatSeat > 0.5 {
			seatBoundary := float64(counts[i]) / nextHalfSeatDown(floatSeat)

			if seatBoundary < minSeatBoundary {
				secondSeatBoundary = minSeatBoundary
				minSeatBoundary = seatBoundary
			} else {
				if seatBoundary != minSeatBoundary && seatBoundary < secondSeatBoundary {
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

// absUintDiff calculates the absolute value of the difference of two Uints.
func absUintDiff(a, b uint) uint {
	if a >= b {
		return a - b
	} else {
		return b - a
	}
}

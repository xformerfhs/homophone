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
// Version: 1.1.0
//
// Change history:
//    2025-02-08: V1.0.0: Created.
//    2025-02-10: V1.1.0: Refactor constants out of statements.
//

// Package distributor_test contains the tests for the count distributors.
package distributor_test

import (
	"homophone/constraints"
	"homophone/distributor"
	"slices"
	"testing"
)

// ******** Private constants ********

const formatExpectedGot = `Expected %v, got %v`

const formatExpectedRuns = `Expected %d runs, got %d: %v`

const formatExpectedSpecificRuns = `Expected %d times %d seats, got %v`

// ******** Test functions ********

func TestSainteLagueDistributionSachsen(t *testing.T) {
	counts := []uint{749_216, 719_274, 277_173, 172_002, 119_964, 104_888}
	seatsCount := uint(119)
	expectedSeats := []uint{41, 40, 15, 10, 7, 6}
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	if !slices.Equal(expectedSeats, seats) {
		t.Errorf(formatExpectedGot, expectedSeats, seats)
	}
}

func TestSainteLagueMallersdorfPfaffenberg(t *testing.T) {
	counts := []int16{28_206, 18_251, 10_000, 9_229, 1_487}
	seatsCount := uint(20)
	expectedSeats := []uint{8, 6, 3, 3, 0}
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	if !slices.Equal(expectedSeats, seats) {
		t.Errorf(formatExpectedGot, expectedSeats, seats)
	}
}

func TestSainteLagueOnlyOne(t *testing.T) {
	counts := []uint8{5}
	seatsCount := uint(37)
	expectedSeats := []uint{37}
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	if !slices.Equal(expectedSeats, seats) {
		t.Errorf(formatExpectedGot, expectedSeats, seats)
	}
}

func TestSainteLagueEquals(t *testing.T) {
	counts := []uint{7, 7, 7, 7, 7}
	seatsCount := uint(12)
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	runMap := runs(seats)
	if len(runMap) != 2 {
		t.Errorf(formatExpectedRuns, 2, len(runMap), runMap)
	}
	if runMap[2] != 3 {
		t.Errorf(formatExpectedSpecificRuns, 3, 2, runMap[2])
	}
	if runMap[3] != 2 {
		t.Errorf(formatExpectedSpecificRuns, 2, 3, runMap[3])
	}
}

func TestSainteLagueZeroesAndOnes(t *testing.T) {
	counts := []uint{0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	seatsCount := uint(52)
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	runMap := runs(seats)
	if len(runMap) != 3 {
		t.Errorf(formatExpectedRuns, 3, len(runMap), runMap)
	}
	if runMap[10] != 3 {
		t.Errorf(formatExpectedSpecificRuns, 3, 10, runMap[10])
	}
	if runMap[11] != 2 {
		t.Errorf(formatExpectedSpecificRuns, 2, 11, runMap[11])
	}
}

func TestSainteLagueRandom(t *testing.T) {
	counts := []uint16{40, 12, 18, 14, 94, 17, 13, 17, 50, 16, 6, 28, 14, 27, 45, 14, 13, 36, 43, 70, 31, 13, 14, 14, 14, 15}
	seatsCount := uint(26)
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	runMap := runs(seats)
	if len(runMap) != 3 {
		t.Errorf(formatExpectedRuns, 3, len(runMap), runMap)
	}
	if runMap[0] != 5 {
		t.Errorf(formatExpectedSpecificRuns, 5, 0, runMap[0])
	}
	if runMap[1] != 16 {
		t.Errorf(formatExpectedSpecificRuns, 16, 1, runMap[1])
	}
	if runMap[2] != 5 {
		t.Errorf(formatExpectedSpecificRuns, 5, 2, runMap[2])
	}
}

// ******** Private functions ********

// total returns the total count.
func total[T constraints.Integer](counts []T) uint {
	result := uint(0)

	for _, count := range counts {
		result += uint(count)
	}

	return result
}

func runs(seats []uint) map[uint]uint {
	result := make(map[uint]uint)
	for _, seat := range seats {
		result[seat]++
	}
	return result
}

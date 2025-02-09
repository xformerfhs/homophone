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

// Package distributor_test contains the tests for the count distributors.
package distributor_test

import (
	"homophone/constraints"
	"homophone/distributor"
	"slices"
	"testing"
)

func TestSainteLagueDistributionSachsen(t *testing.T) {
	counts := []uint{749_216, 719_274, 277_173, 172_002, 119_964, 104_888}
	seatsCount := uint(119)
	expectedSeats := []uint{41, 40, 15, 10, 7, 6}
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	if !slices.Equal(expectedSeats, seats) {
		t.Errorf("Expected %v, got %v", expectedSeats, seats)
	}
}

func TestSainteLagueMallersdorfPfaffenberg(t *testing.T) {
	counts := []int16{28_206, 18_251, 10_000, 9_229, 1_487}
	seatsCount := uint(20)
	expectedSeats := []uint{8, 6, 3, 3, 0}
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	if !slices.Equal(expectedSeats, seats) {
		t.Errorf("Expected %v, got %v", expectedSeats, seats)
	}
}

func TestSainteLagueOnlyOne(t *testing.T) {
	counts := []uint8{5}
	seatsCount := uint(37)
	expectedSeats := []uint{37}
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	if !slices.Equal(expectedSeats, seats) {
		t.Errorf("Expected %v, got %v", expectedSeats, seats)
	}
}

func TestSainteLagueEquals(t *testing.T) {
	counts := []uint{7, 7, 7, 7, 7}
	seatsCount := uint(12)
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	runMap := runs(seats)
	if len(runMap) != 2 {
		t.Errorf("Expected 2 runs, got %v: %v", len(runMap), runMap)
	}
	if runMap[2] != 3 {
		t.Errorf("Expected 3 times 2 seats, got %v", runMap[2])
	}
	if runMap[3] != 2 {
		t.Errorf("Expected 2 times 3 seats, got %v", runMap[2])
	}
}

func TestSainteLagueZeroesAndOnes(t *testing.T) {
	counts := []uint{0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	seatsCount := uint(52)
	seats := distributor.SainteLagueDistribution(counts, total(counts), seatsCount)
	runMap := runs(seats)
	if len(runMap) != 3 {
		t.Errorf("Expected 3 runs, got %v: %v", len(runMap), runMap)
	}
	if runMap[10] != 3 {
		t.Errorf("Expected 3 times 10 seats, got %v", runMap[2])
	}
	if runMap[11] != 2 {
		t.Errorf("Expected 2 times 11 seats, got %v", runMap[2])
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

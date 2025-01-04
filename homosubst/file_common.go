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
// Version: 1.0.0
//
// Change history:
//    2025-01-03: V1.0.0: Created.
//

package homosubst

// ******** Private constants ********

// fileMagic is the magic bytes of a substitution file.
var fileMagic = []byte(`HFDF`)

// actVersion is the current version number.
const actVersion byte = 0

// substFileLength is the length of the data in a substitution file.
const substFileLength = 136

// generator is the starter value for the integrity key generation.
var generator = []byte{
	0xfe, 0xb9, 0x66, 0x43,
	0x18, 0x5b, 0x51, 0xdf,
	0x86, 0x99, 0xe5, 0x09,
	0xa4, 0xdc, 0x0d, 0xad,
	0x82, 0xed, 0xc4, 0x30,
}

// salt is the salt needed for the integrity key generation.
var salt = []byte{
	0x74, 0xbc, 0x06, 0x3e,
	0x56, 0x17, 0xda, 0xd4,
	0xf2, 0xc7, 0x91, 0x37,
	0x2a, 0xe2, 0xbf, 0x32,
}

// additionalData is the additional data needed for the integrity check.
var additionalData = []byte(`HoTzpLoZ`)

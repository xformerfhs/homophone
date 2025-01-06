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
// Version: 1.0.1
//
// Change history:
//    2025-01-02: V1.0.0: Created.
//    2025-01-05: V1.0.1: Added forgotten colon in message.
//

package main

import (
	"fmt"
	"homophone/homosubst"
)

// doEncryption encryptions the contents of a file.
func doEncryption(clearFileName string, encryptedFileName string, substitutionFileName string, keepOthers bool) int {
	fmt.Printf("Source file: '%s'\n", clearFileName)

	substitutor, err := homosubst.NewSubstitutor(clearFileName)
	if err != nil {
		return printErrorf(`Error creating substitutor: %v`, err)
	}

	fmt.Println(`Substitutions:`)
	substitutor.Print()

	err = substitutor.Encrypt(clearFileName, encryptedFileName, keepOthers)
	if err != nil {
		return printErrorf(`Error encrypting file: %v`, err)
	}
	fmt.Printf("Encrypted file: '%s'\n", outFileName)

	err = substitutor.Save(substitutionFileName)
	if err != nil {
		return printErrorf(`Error saving substitution file: %v`, err)
	}

	fmt.Printf("Substitution file: '%s'\n", substFileName)

	return rcOK
}

// doDecryption decrypts the contents of an encrypted file.
func doDecryption(encryptedFileName string, decryptedFileName string, substitutionFileName string) int {
	fmt.Printf("Encrypted file: '%s'\n", encryptedFileName)

	substitutor, err := homosubst.NewFromFile(substitutionFileName)
	if err != nil {
		return printErrorf(`Error loading substitution file: %v`, err)
	}
	fmt.Printf("Loaded substitution file: '%s'\n", substitutionFileName)

	fmt.Println(`Substitutions:`)
	substitutor.Print()

	err = substitutor.Decrypt(encryptedFileName, decryptedFileName)
	if err != nil {
		return printErrorf(`Error decrypting file: %v`, err)
	}

	fmt.Printf("Decrypted file: '%s'\n", decryptedFileName)

	return rcOK
}

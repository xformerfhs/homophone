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
//    2025-01-04: V1.0.0: Created.
//    2025-01-05: V1.1.0: Correct handling of additional arguments that are not flags.
//

package main

import (
	"flag"
	"fmt"
	"homophone/filehelper"
	"os"
)

// ******** Private variables ********

// Option values.

// They have to be global in order to modularize the main program.
// Otherwise, there would have been an awful lot of parameters to pass to functions.

// inFileName is the name of the input file.
var inFileName string

// outFileName is the name of the output file.
var outFileName string

// substFileName is the name of the substitution file.
var substFileName string

// keepOthers indicates that characters that are not in the range A-Z should be kept.
var keepOthers bool

// Flag sets.

// decryptCommand is the [flag.Flagset] for decryption.
var decryptCommand *flag.FlagSet

// encryptCommand is the [flag.Flagset] for encryption.
var encryptCommand *flag.FlagSet

// Program information.

// myName is the program name.
var myName string

// ******** Init function ********

// init initializes the command line variables.
func init() {
	myName = filehelper.RealBaseName(os.Args[0])
}

// ******** Private functions ********

// defineCommandLineFlags defines the command line flags.
func defineCommandLineFlags() {
	encryptCommand = flag.NewFlagSet(`encrypt`, flag.ExitOnError)
	encryptCommand.StringVar(&inFileName, `in`, ``, "Clear text file `path`")
	encryptCommand.StringVar(&outFileName, `out`, ``, "Encrypted file `path`")
	encryptCommand.StringVar(&substFileName, `key`, ``, "Key file `path`")
	encryptCommand.BoolVar(&keepOthers, `keep`, false, `Keep characters that are not in the range A-Z (default: do not keep)`)

	decryptCommand = flag.NewFlagSet(`decrypt`, flag.ExitOnError)
	decryptCommand.StringVar(&inFileName, `in`, ``, "Encrypted file `path`")
	decryptCommand.StringVar(&outFileName, `out`, ``, "Decrypted file `path`")
	decryptCommand.StringVar(&substFileName, `key`, ``, "Key file `path`")

	flag.Usage = myUsage
}

// parseDecryption parses the command line of a "decrypt" command.
func parseDecryption() int {
	err := decryptCommand.Parse(os.Args[2:])
	if err != nil {
		return rcHelpOrError(err)
	}

	return checkDecryptionFlags()
}

// parseEncryption parses the command line of an "encrypt" command.
func parseEncryption() int {
	err := encryptCommand.Parse(os.Args[2:])
	if err != nil {
		return rcHelpOrError(err)
	}

	return checkEncryptionFlags()
}

// checkDecryptionFlags checks the decryption flags.
func checkDecryptionFlags() int {
	rc := checkFlagsCommon(`encrypted`, decryptCommand.Args())
	if rc != rcOK {
		return rc
	}

	if len(outFileName) == 0 {
		outFileName = buildDecryptOutFilePath(inFileName)
	}

	return rcOK
}

// checkEncryptionFlags checks the encryption flags.
func checkEncryptionFlags() int {
	rc := checkFlagsCommon(`clear text`, encryptCommand.Args())
	if rc != rcOK {
		return rc
	}

	if len(outFileName) == 0 {
		outFileName = buildEncryptOutFilePath(inFileName)
	}

	return rcOK
}

// checkFlagsCommon does the checks common to all commands.
func checkFlagsCommon(typeName string, additionalArgs []string) int {
	if len(additionalArgs) > 0 {
		return printUsageErrorf(`Arguments without flags present: %s`, additionalArgs)
	}

	if len(inFileName) == 0 {
		return printUsageErrorf(`Name of %s file is missing`, typeName)
	}

	if len(substFileName) == 0 {
		substFileName = buildSubstFilePath(inFileName)
	}

	return rcOK
}

// myUsage is the function that is called by flag.Usage. It prints the usage information.
func myUsage() {
	errWriter := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(errWriter, "\n'%s' implements a homomorphic encryption.\n", myName)
	_, _ = fmt.Fprintln(errWriter, `The characters A-Z are mapped to characters in the range A-Z and a-z.
Characters in the range a-z are converted to upper case and then mapped to the range A-Z and a-z.
First, the frequency of the characters in the plaintext is analyzed and from this it is calculated which character in the plaintext must be replaced by how many characters in the ciphertext.
I.e. the more often a character appears in plain text, the more characters replace it in the encrypted text.
The character substitution table is written to a key file which is required for decryption.`)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintf(errWriter, "The usage is: %s {command} {options...}\n", myName)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `The following commands are available:`)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `decrypt: Decrypt an encrypted file`)
	decryptCommand.PrintDefaults()
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `If the 'key' file path is not specified the name 'infilebasename_ext.subst' is used.`)
	_, _ = fmt.Fprintln(errWriter, `If the 'out' file path is not specified the name 'infilebasename_homophone.txt' is used.`)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `Options can be started with either '-' or '--'`)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `encrypt: Encrypt a file`)
	encryptCommand.PrintDefaults()
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `If the 'key' file path is not specified the name 'infilebasename_ext.subst' is used.`)
	_, _ = fmt.Fprintln(errWriter, `If the 'out' file path is not specified the name 'infilebasename_decrypted.txt' is used.`)
	_, _ = fmt.Fprintln(errWriter, `If 'keep' is specified, all characters not in the range A-Z are kept and copied to the output file`)
	_, _ = fmt.Fprintln(errWriter, `If 'keep' is not specified, only characters in the range A-Z are copied to the output file. All others are discarded`)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `Options can be started with either '-' or '--'`)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `version: Print version information`)
	_, _ = fmt.Fprintln(errWriter)
	_, _ = fmt.Fprintln(errWriter, `help: Print this usage information`)
	_, _ = fmt.Fprintln(errWriter)
}

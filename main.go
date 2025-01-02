package main

import (
	"fmt"
	"homophone/homosubst"
	"os"
	"strings"
)

func main() {
	os.Exit(realMain(os.Args[1:]))
}

const (
	rcOK  = 0
	rcErr = 1
)

func realMain(args []string) int {
	numArgs := len(args)
	if numArgs < 1 {
		return printErrorUsage(`Not enough arguments`)
	}

	cmd := strings.ToLower(args[0])
	switch cmd[0] {
	case 'e':
		return doEncryption(args[1], numArgs > 1)

	case 'd':
		if len(args) < 3 {
			return printErrorUsage(`Not enough arguments`)
		}
		return doDecryption(args[1], args[2])

	default:
		return printErrorfUsage(`Unknown command: '%s'`, cmd)
	}
}

func doEncryption(fileName string, noOther bool) int {
	fmt.Printf("Source file: '%s'\n", fileName)

	substitutor, err := homosubst.NewSubstitutor(fileName)
	if err != nil {
		return printErrorf(`Error creating substitutor: %v`, err)
	}

	fmt.Println(`Substitutions:`)
	substitutor.Print()

	var outFileName string
	outFileName, err = substitutor.Encrypt(noOther)
	if err != nil {
		return printErrorf(`Error encrypting file: %v`, err)
	}
	fmt.Printf("Encrypted file: '%s'\n", outFileName)

	var substFileName string
	substFileName, err = substitutor.Save()
	if err != nil {
		return printErrorf(`Error saving substitution file: %v`, err)
	}

	fmt.Printf("Substitution file: '%s'\n", substFileName)

	return rcOK
}

func doDecryption(fileName string, substitutionFileName string) int {
	fmt.Printf("Source file: '%s'\n", fileName)

	substitutor, err := homosubst.NewLoad(substitutionFileName)
	if err != nil {
		return printErrorf(`Error loading substitution file: %v`, err)
	}
	fmt.Printf("Loaded substitution file '%s'\n", substitutionFileName)

	fmt.Println(`Substitutions:`)
	substitutor.Print()

	var outFileName string
	outFileName, err = substitutor.Decrypt(fileName)
	if err != nil {
		return printErrorf(`Error decrypting file: %v`, err)
	}

	fmt.Printf("Decrypted file: '%s'\n", outFileName)

	return rcOK
}

func printErrorf(msgFormat string, args ...any) int {
	return printError(fmt.Sprintf(msgFormat, args...))
}

func printError(msg string) int {
	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintln(os.Stderr, msg)

	return rcErr
}

func printErrorfUsage(msgFormat string, args ...any) int {
	return printErrorUsage(fmt.Sprintf(msgFormat, args...))
}

func printErrorUsage(msg string) int {
	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintln(os.Stderr, msg)

	return printUsage()
}

func printUsage() int {
	_, _ = fmt.Fprintln(os.Stderr, "\nUsage:\n   encrypt {file} [{noOther}]\n   decrypt {file} {subst file}")
	return rcErr
}

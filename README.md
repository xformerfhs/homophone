# homophone

A program for a simple homophonic encryption of text files.

[![Go Report Card](https://goreportcard.com/badge/github.com/xformerfhs/hashvalue)](https://goreportcard.com/report/github.com/xformerfhs/homophone)
[![License](https://img.shields.io/github/license/xformerfhs/hashvalue)](https://github.com/xformerfhs/hashvalue/blob/main/LICENSE)

## Introduction

In my cryptography courses, I talk about various "classic" encryption methods.
One of them is the ["homophonic" encryption](https://en.wikipedia.org/wiki/Substitution_cipher#Homophonic).
This is a substitution cipher with a twist.

In natural languages certain characters are much more common than others.
In indoeuropean languages the letter `e` is much more prevalent than other characters.
A simple [substitution cipher](https://en.wikipedia.org/wiki/Substitution_cipher) preserves the character frequencies, so the substitution character for `e` would stand out.

In order to flatten the statistical distribution of the characters one can map the frequent characters with multiple substitutions.
For this to work the destination alphabet has to have much more characters than the source alphabet.
E.g. the rarely occurring character `B` is replaced by `j`, but the frequently occurring `E` is replaced by multiple characters like `yWSqXCBb`.
Any of these characters would code for an `E`.

To get the right number of substitution characters one has to count the frequencies of the characters in the source text and adjust the number of substitution characters so that each substitution characters for each source character has about the same probability of occurring.

## Implementation

In this program all characters are mapped to upper case before encryption.
This makes the source alphabet have 26 characters: `ABCDEFGHIJKLMNOPQRSTUVWXYZ`.
The destination alphabet is set to `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`.
So 26 source characters are mapped to 62 destination characters.

On encryption the correct number of substitution characters is calculated for each of the source characters.
Then each of the substitution characters is randomly assigned.
Then each character in the clear text is replaced with a random character of its substitution list.

The "key" for this encryption are the substitution lists for the characters.
This key is saved in a separate file.

On decryption the substitution lists are read from the key file and each characters of the encrypted file is replaced by the original character.
Since all characters are converted to uppercase before encryption, the decrypted file characters are all uppercase.

## Call

The program is called like this:

```
homophone <command> [options...]
```

There are the following commands:

| Command     | Meaning                                                                                                         |
|-------------|-----------------------------------|
| `decrypt`   | Decrypt an encrypted file.        |
| `encrypt`   | Encrypt a clear text file.        |
| `help`      | Show usage information.           |
| `version`   | Show program version information. |

Two of these commands use options, namely `decrypt` and `encrypt`:

### Decrypt

The options for the `decrypt` command are the following:

```
homophone decrypt -in <encrypted file path> [-out <decrypted file path>] [-key <key file path>] 
```

| Option | Meaning                                                                   |
|--------|---------------------------------------------------------------------------|
| `in`   | Path of the encrypted file (input, required).                             |
| `out`  | Path of the file that will receive the decrypted text (output, optional). |
| `key`  | Path of the key file (input, optional).                                   |

The options can be started with either `--` or `-`.

If the `out` file path is not specified it is set to `<infile-path>/<infile-basename-without-homophone>_decrypted.<infile-extension>`.
If the `key` file path is not specified it is set to `<infile-path>/<infile-basename-without-homophone>_<infile_extension>.subst`.

E.g., if the name of the input file is `something_homophone.txt` the default name of the output file is `something_decrypted.txt` and the default name for the key file is `something_txt.subst`.

### Encrypt

The options for the `encrypt` command are the following:

```
homophone encrypt -in <clear text file path> [-out <encrypted file path>] [-key <key file path>] [-keep] 
```

| Option | Meaning                                                                                         |
|--------|-------------------------------------------------------------------------------------------------|
| `in`   | Path of the clear text file (input, required).                                                  |
| `out`  | Path of the file that will receive the encrypted text (output, optional).                       |
| `key`  | Path of the key file (output, optional).                                                        |
| `keep` | Characters that are not in range `A-Z` after conversion to uppercase are preserved (optional).  |

If `keep` is not specified characters that are not in range `A-Z` after conversion to upper case are discarded.

The options can be started with either `--` or `-`.

If the `out` file path is not specified it is set to `<infile-path>/<infile-basename>_homophone.<infile-extension>`.
If the `key` file path is not specified it is set to `<infile-path>/<infile-basename>_<infile_extension>.subst`.

E.g., if the name of the input file is `something.txt` the default name of the output file is `something_homophone.txt` and the default name for the key file is `something_txt.subst`.

### Examples

In the first example a text file with the name `message.txt` is encrypted:

```
homophone encrypt -in message.txt
```

This prints the following output:

```
Source file: 'message.txt'
Substitutions:
   A: ONjK
   B: s
   C: 4J
   D: hn
   E: PtLTGgVd
   F: p5
   G: E
   H: 20i
   I: wyS
   J: 9
   K: X
   L: FbQ
   M: r
   N: MxZ
   O: uzaB
   P: c
   Q: v
   R: 68fI
   S: U1D
   T: 7HWle
   U: qR
   V: Y
   W: Co
   X: k
   Y: mA
   Z: 3
Encrypted file: 'message_homophone.txt'
Substitution file: 'message_txt.subst'
```

The first line prints the name of the source file.
Then the substitution mapping is printed.
To the left is the clear text character and to the right the characters it is substituted with.
If there are no letters after the colon the clear text character does not appear in the clear text.
This is list is the key and is saved in the `key` file.

The encrypted file looks like this:

```
JiKIQgDToLbFPfOqe2a8BpjUR4JG1D5qQszuXNsuRWH0dVO6bA0y1laImzp72teAcdC8wlT64jrtqcCS7iNHfKSZwxEGkP8JyUV5B6mzRMEWAcLoIwlgfDmTUSMe2tVOIFA17igq1PfDBpHmcLo8yWd6UCGILXxaoZNDlAcgC6w7V81CdQFPf4KrTRcoye0W2tUGxHLM4gjDjcIKJ7SJGn6SbFuZdEzuhynPNbTOhUHBNZaeiVfQKlt8W0twnPjoO1ruhw5SdhlBILjnMaCy1H0V7wrTpzfObFEuznrgxWa4BrGeuHigNSha5W2Gy64zRx78ml2SUCBfXDaqeK1WiT6PvRw8VnDLYdZHAJ0OINJlt8U4uqMeyME7iT5ypedVZUcjJLDsGWogPxCzfh1KMnH2tpRQQD7BclzESYtO4bK1UwJjF5qQFbSxPa5HmcVwZ4yhTZlNQFmNMuW0dILYGxrB6gpjraqU7AcSxEtkLfJwDPy1e0VvRS4Xs8BoZ5zk9Rrc1uYT6WigbK3AhzEmuqrOANbIdjnmXMaoB6Cu8XGhyezRlpBIH2L5wfD7lyrT70dVkG64wUPWigvqSJXsfaoZpuk9RrcUBYt8e2dFO3AnzEyDKxgkKrcQPapjcNMEIOrK1VxHtZJG4BM7NSxwZEjbQH2LFTleP8Dz5W0LOQciKsdl
```

The encrypted file can be decrypted with the following call:

```
homophone decrypt -in message_homophone.txt
```
This prints the following output:

```
Encrypted file: 'message_homophone.txt'
Loaded substitution file 'message_txt.subst'
Substitutions:
   A: ONjK
   B: s
   C: 4J
   D: hn
   E: PtLTGgVd
   F: p5
   G: E
   H: 20i
   I: wyS
   J: 9
   K: X
   L: FbQ
   M: r
   N: MxZ
   O: uzaB
   P: c
   Q: v
   R: 68fI
   S: U1D
   T: 7HWle
   U: qR
   V: Y
   W: Co
   X: k
   Y: mA
   Z: 3
Decrypted file: 'message_decrypted.txt'
```

In the next example all options are used for encryption:

```
homophone encrypt -in message.txt -out secret.txt -key innocent.jpg -keep
```

This prints the following output:

```
Source file: 'message.txt'
Substitutions:
   A: W6Lm
   B: 3
   C: Hv
   D: qd
   E: TMiyzYOF
   F: Ea
   G: 5
   H: 94N
   I: JkV
   J: j
   K: u
   L: t2B
   M: c
   N: oCX
   O: Anhp
   P: G
   Q: 7
   R: Ix8r
   S: wbQ
   T: gSK10
   U: se
   V: U
   W: ZR
   X: l
   Y: fP
   Z: D
Encrypted file: 'secret.txt'
Substitution file: 'innocent.jpg'
```

Note that the substitutions are different.
The program randomly assigns substitutions for each encryption call.

The encrypted file `secret.txt` looks like this:

```
v9W82yw F ZMtBYx, ms14hr na L beHHzQbEeB 3Apu 63ps0 KNi T6I2P NkQSnIf AE g9O SfGYRxV1zr, vWci sG RJg4 m K8LVXkC5 ylOIvJwF ahr PAeo5 0PGTZ8JSMxQ. (fYw, ko g4y i68tP 1900b, K9O ebzxQ hE 0fGFRrV1TIw ZMxM uCnZX mQ 1fGTRIJKzrw. RYt2O8 HLci sG ZVSN gNy bFX0Tovi Wb L GI6HSkHY q8VB2).

pCF 5pnd kdMm tOWqw gA mXh04zr

B6Kyx, 19i JqML ZWQ cAdkaJFq Sp IOmd:

CnR Vw 09Y gVcy Eh8 LBt 5Anq cTo 1h vpcz Kn g4y 6kd ha SNFJx vpsCKrP.

19JQ RAxub he0 WQ 0NM IT7sV8Oq wYUioSP H4Lr6HKzIb (vAeX1kX5 g9M EkagiFC wGmvOQ 3YKZzTo Zp8db WXq SNy as22 Q1nG 0p 5VUz m HtWwbJH6B Eet2 BJoF na SfGO.

VCvkdTC1L2tf, WXAK4Yx iUMo chry E6chsb gPGkX5 YlzrvJQT Vw

04i 7ekHu 3xnRC aAl jscGw pUy8 gNF BmDf qA5

Ppe cLP WBIO6df uohZ, nI RhruMd JS ns0 Ep8 19F aVxQK gVcM.
KNO ylz8Hkbi "S4T 7eJvu 3rARX Ehl jecGb pUYx 19Y tmDf qA5" Jw LC Fl6cG2y na W Gmo5ILc, m QiC0MoHO vnXKLkXVC5 WBt 09T 2zg1y8Q pE S4i 62GNm3zg.
```

Because of the `keep` option all blanks, paragraphs, brackets and punctuations marks are preserved.

Because of the preservation of the text structure and word boundaries, this text is much easier to analyse and decrypt without the key, than the text without the kept structure.

The file can be decrypted with the following call:

```
homophone decrypt -in secret.txt -out clear.txt -key innocent.jpg
```

This produces the following output:

```
Encrypted file: 'secret.txt'
Loaded substitution file 'innocent.jpg'
Substitutions:
   A: W6Lm
   B: 3
   C: Hv
   D: qd
   E: TMiyzYOF
   F: Ea
   G: 5
   H: 94N
   I: JkV
   J: j
   K: u
   L: t2B
   M: c
   N: oCX
   O: Anhp
   P: G
   Q: 7
   R: Ix8r
   S: wbQ
   T: gSK10
   U: se
   V: U
   W: ZR
   X: l
   Y: fP
   Z: D
Decrypted file: 'clear.txt'
```

### Return codes

The possible return codes are the following:

| Code | Meaning                   |
|------|---------------------------|
| `0`  | Successful processing     |
| `1`  | Error in the command line |
| `2`  | Error while processing    |

## Program build

You must have Go installed to create the program.
This creates a directory that is specified in the `GOPATH` environment variable.
Under Windows, this is the home directory, e.g. `D:\Users\username\go`.
Under Linux this is `${HOME}/go`.
In that directory there is a subdirectory `src`.

To create the program, the source code must be stored under `%GOPATH%\src\hashvalue` or `${HOME}/go/src/hashvalue`.
Then one has to start the batch file `gb.bat` or the shell script `gb`, which builds the executables.
These scripts expect the UPX program to be in a specific location.
This location can be adapted to the local path.
If UPX is not available, no compression is performed.

As a result, the files `homophone` for Linux and `homophone.exe` for Windows are created.

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).

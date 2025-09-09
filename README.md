# helperFunctions
Miscellaneous helper functions
___

**IMPORTANT NOTE**
This package will enter major refactoring soon; this current document is a bit outdated as it is.


This package provides functions that I have been using all over my various tools but the obvious struck: instead of duplicating code on my filesystem, I should have a single source of truth.

## Overview

This package is currently divided in 4 files; if the scope of the project expands (and I expect it to expand, at some point), we might be talking about sub-packages here.

For now, we have 4 files, with a specific area of responsibility:

| File                    | Area / coverage                                                                                                                                    |
|-------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| prompt4Types.go          | Returns a value of a given type (string, bool, int, etc), with a prompt<br>This might seem useless, but I've needed something like that many times |
| terminal.go             | Terminal-related functions, such as colouring output, getting terminal size (WxH), etc                                                             |
| encodeDecodePassword.go | Functions related to encoding/decoding strings, prompting for (non-echoed) passwords, etc)                                                         |
| misc.go                 | Various minor functions                                                                                                                            | 

This is basic, for now. I really intend on expanding on this package, and eventually have it documented on pkg.go.dev/

## Installation
A simple `go get github.com/jeanfrancoisgratton/functionHelpers` and we're done.

## Usage, per file/package
(it is "per file" for now, as I'm not dividing it in sub-packages)

### FILE: misc.go
Two functions, SI(), and ReverseString().

#### SI()
SI came from a need to have comma-separated number out of a number that wasn't.

For instance, you have the number `123456789`. In a comma-separated format (SI-notation, SI for "Système International"); it becomes, in a string: `123,456,789`

You input a number, it will come back as an SI-formatted string.

#### ReverseString()
The name says it all:

You input a string ("abcdef"), it returns its reverse ("fedcba").


### FILE: prompt4Types.go
Many functions there, all following the same pattern...

The functions:

GetStringValFromPrompt()<br>
GetIntValFromPrompt()<br>
GetBoolValFromPrompt()<br>
GetStringSliceFromPrompt()<br>

The structure is intuitive: Get<DATATYPE><SCALAR/SLICE>FromPrompt()

DATATYPE is either String, Int, Bool... for now

SCALAR/SLICE is either Val (Scalar, single value) or Slice. This affects the return data type.

All take a single parameter, the "prompt", ie "the message/question you want to pass to the user"

### FILE: terminal.go
This one is basically (for now) used to colour the output, or get the TTY size, or clear a TTY.

To get the current TTY size, you call GetTerminalSize(), it'll return the WIDTHxHEIGHT (2 int values)

To clear a terminal TTY is quite easy: just call ClearTTY(), no params.

The next functions are "colour functions:"<br>
Red(sentence)<br>
Green(sentence)<br>
White(sentence)<br>
Yellow(sentence)<br>
Blue(sentence)<br>

Where "sentence" is a variable of type string that you want to return in the appropriate colour.


### FILE: encodeDecodePassword.go
3 functions in here :

**GetPassword(sentence)** : you will be prompted (the sentence string type in parameters) to enter a password. __The password is not echoed on the TTY; it will be returned in a non-encoded form as a string_. If you wish to return it as an encoded form, you can use the next function...

**EncodeString(string to be encoded, secret key)**
The string to be encoded will be returned. If secret key is empty, a default secret key is provided.
Of course, this is not the suggested behaviour...

_Important note_: secret key has to be **exactly** 32 bytes long, otherwise the default key is used.

**DecodeString(string to be decoded, secret key)**
The string to be decoded will be returned. If secret key is empty, a default secret key is provided.
Of course, this is not the suggested behaviour...

_Important note_: secret key has to be **exactly** 32 bytes long, otherwise the default key is used.



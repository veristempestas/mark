// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Generating random text: a Markov chain algorithm

Based on the program presented in the "Design and Implementation" chapter
of The Practice of Programming (Kernighan and Pike, Addison-Wesley 1999).
See also Computer Recreations, Scientific American 260, 122 - 125 (1989).

A Markov chain algorithm generates text by creating a statistical model of
potential textual suffixes for a given prefix. Consider this text:

	I am not a number! I am a free man!

Our Markov chain algorithm would arrange this text into this set of prefixes
and suffixes, or "chain": (This table assumes a prefix length of two words.)

	Prefix       Suffix

	"" ""        I
	"" I         am
	I am         a
	I am         not
	a free       man!
	am a         free
	am not       a
	a number!    I
	number! I    am
	not a        number!

To generate text using this table we select an initial prefix ("I am", for
example), choose one of the suffixes associated with that prefix at random
with probability determined by the input statistics ("a"),
and then create a new prefix by removing the first word from the prefix
and appending the suffix (making the new prefix is "am a"). Repeat this process
until we can't find any suffixes for the current prefix or we exceed the word
limit. (The word limit is necessary as the chain table may contain cycles.)

Our version of this program reads text from standard input, parsing it into a
Markov chain, and writes generated text to standard output.
The prefix and output lengths can be specified using the -prefix and -words
flags on the command-line.
*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// Parse commandline arguments
	command := os.Args[1]

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	if command == "Generate" || command == "generate" {
		modefile := os.Args[2]
		n, err := strconv.Atoi(os.Args[3])
		if err != nil {
			panic("Error while parsing to get the number of words for generating a text")
		}
		numWords := n
		tableFile := modefile
		table := ReadTableFromFile(tableFile)
		text := table.Generate(numWords)
		fmt.Println(text)

	} else if command == "Read" || command == "read" {
		p, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic("Error while parsing to get a prefix length")
		}
		prefixLen := p
		outputFile := os.Args[3]
		files := os.Args[4:]

		// build a Chain from the input files
		c := NewChain(prefixLen)
		for _, file := range files {
			c.BuildFromFile(file)
		}
		table := NewFreqTable(prefixLen)
		table.Build(*c)
		table.WriteMapToFile(outputFile)
	} else {
		panic("I don't recognize this command, try Generate/generate/Read/read ")
	}
}

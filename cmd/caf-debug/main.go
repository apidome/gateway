package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

func main() {
	query := `
	{
		#commet
		empireHero: hero(episode: EMPIRE) {
		  name
		}
		jediHero: hero(episode: JEDI) {
		  name
		}
	}
	`

	var sc scanner.Scanner

	sc.Mode = scanner.ScanIdents | scanner.ScanInts |
		scanner.ScanFloats | scanner.ScanChars |
		scanner.ScanStrings | scanner.ScanRawStrings |
		scanner.ScanComments

	sc.Init(strings.NewReader(query))

	for sc.Scan() != scanner.EOF {
		fmt.Println(sc.TokenText())
	}
}

package main

import (
	"fmt"

	css "github.com/a-h/templ-css"
)

type Input struct {
	CSS    string
	Inline bool
}

var inputs = []Input{
	{CSS: `font-color: red;`, Inline: true},
	{CSS: `font-color: {{ red }};`, Inline: true},
	{CSS: `a:hover {
		background-color: yellow;
}`, Inline: false},
	{CSS: `@media print {.class{width:5px;}}`, Inline: false},
}

func main() {
	for _, input := range inputs {
		tokens, err := css.Parse(input.CSS, input.Inline)
		if err != nil {
			fmt.Printf("Failed to parse CSS: %v\n", err)
			continue
		}
		fmt.Println(input.CSS)
		fmt.Println(css.PrintTokens(tokens))
		for _, token := range tokens {
			fmt.Printf("%#v\n", token)
		}
		fmt.Println()
	}
}

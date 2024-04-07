package css

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tdewolff/parse/v2/css"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		inline   bool
		expected []Token
	}{
		{
			name:   "attribute key/value",
			input:  `font-color: red;`,
			inline: true,
			expected: []Token{
				NewToken(0, css.IdentToken, "font-color", false),
				NewToken(10, css.ColonToken, ":", false),
				NewToken(11, css.WhitespaceToken, " ", false),
				NewToken(12, css.IdentToken, "red", false),
				NewToken(15, css.SemicolonToken, ";", false),
			},
		},
		{
			name:   "attribute key/value expression",
			input:  `font-color: {{ red }};`,
			inline: true,
			expected: []Token{
				NewToken(0, css.IdentToken, "font-color", false),
				NewToken(10, css.ColonToken, ":", false),
				NewToken(11, css.WhitespaceToken, " ", false),
				NewToken(12, css.IdentToken, "{{ red }}", true),
				NewToken(21, css.SemicolonToken, ";", false),
			},
		},
		{
			name: "hover pseudo class",
			input: `a:hover {
  background-color: yellow;
}`,
			inline: false,
			expected: []Token{
				NewToken(0, css.IdentToken, "a", false),
				NewToken(1, css.ColonToken, ":", false),
				NewToken(2, css.IdentToken, "hover", false),
				NewToken(7, css.WhitespaceToken, " ", false),
				NewToken(8, css.LeftBraceToken, "{", false),
				NewToken(9, css.WhitespaceToken, "\n  ", false),
				NewToken(11, css.IdentToken, "background-color", false),
				NewToken(26, css.ColonToken, ":", false),
				NewToken(27, css.WhitespaceToken, " ", false),
				NewToken(28, css.IdentToken, "yellow", false),
				NewToken(34, css.SemicolonToken, ";", false),
				NewToken(35, css.WhitespaceToken, "\n", false),
				NewToken(36, css.RightBraceToken, "}", false),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Parse(test.input, test.inline)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			if diff := cmp.Diff(test.expected, actual); diff != "" {
				t.Error(test.input)
				t.Error(actual)
				t.Error(diff)
			}

			fmt.Println(PrintTokens(actual))
		})
	}
}

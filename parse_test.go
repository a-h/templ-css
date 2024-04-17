package css

import (
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
				NewCSSToken(0, css.IdentToken, "font-color"),
				NewCSSToken(10, css.ColonToken, ":"),
				NewCSSToken(11, css.WhitespaceToken, " "),
				NewCSSToken(12, css.IdentToken, "red"),
				NewCSSToken(15, css.SemicolonToken, ";"),
			},
		},
		{
			name:   "attribute key/value expression",
			input:  `font-color: {{ red }};`,
			inline: true,
			expected: []Token{
				NewCSSToken(0, css.IdentToken, "font-color"),
				NewCSSToken(10, css.ColonToken, ":"),
				NewCSSToken(11, css.WhitespaceToken, " "),
				NewGoToken(12, "{{ ", "red", " }}"),
				NewCSSToken(21, css.SemicolonToken, ";"),
			},
		},
		{
			name: "hover pseudo class",
			input: `a:hover {
  background-color: yellow;
}`,
			inline: false,
			expected: []Token{
				NewCSSToken(0, css.IdentToken, "a"),
				NewCSSToken(1, css.ColonToken, ":"),
				NewCSSToken(2, css.IdentToken, "hover"),
				NewCSSToken(7, css.WhitespaceToken, " "),
				NewCSSToken(8, css.LeftBraceToken, "{"),
				NewCSSToken(9, css.WhitespaceToken, "\n  "),
				NewCSSToken(12, css.IdentToken, "background-color"),
				NewCSSToken(28, css.ColonToken, ":"),
				NewCSSToken(29, css.WhitespaceToken, " "),
				NewCSSToken(30, css.IdentToken, "yellow"),
				NewCSSToken(36, css.SemicolonToken, ";"),
				NewCSSToken(37, css.WhitespaceToken, "\n"),
				NewCSSToken(38, css.RightBraceToken, "}"),
			},
		},
		{
			name:   "media query",
			input:  `@media print {.class{width:5px;}}`,
			inline: false,
			expected: []Token{
				NewCSSToken(0, css.AtKeywordToken, "@media"),
				NewCSSToken(6, css.WhitespaceToken, " "),
				NewCSSToken(7, css.IdentToken, "print"),
				NewCSSToken(12, css.WhitespaceToken, " "),
				NewCSSToken(13, css.LeftBraceToken, "{"),
				NewCSSToken(14, css.DelimToken, "."),
				NewCSSToken(15, css.IdentToken, "class"),
				NewCSSToken(20, css.LeftBraceToken, "{"),
				NewCSSToken(21, css.IdentToken, "width"),
				NewCSSToken(26, css.ColonToken, ":"),
				NewCSSToken(27, css.DimensionToken, "5px"),
				NewCSSToken(30, css.SemicolonToken, ";"),
				NewCSSToken(31, css.RightBraceToken, "}"),
				NewCSSToken(32, css.RightBraceToken, "}"),
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
				t.Error(test.input)
				t.Error(PrintTokens(actual))
			}
		})
	}
}

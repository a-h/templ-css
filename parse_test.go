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
				NewToken(0, css.IdentToken, "font-color", false),
				NewToken(10, css.ColonToken, ":", false),
				NewToken(11, css.IdentToken, "red", false),
				NewToken(14, css.SemicolonToken, ";", false),
			},
		},
		{
			name:   "attribute key/value expression",
			input:  `font-color: {{ red }};`,
			inline: true,
			expected: []Token{
				NewToken(0, css.IdentToken, "font-color", false),
				NewToken(10, css.ColonToken, ":", false),
				NewToken(11, css.IdentToken, "{{ red }}", true),
				NewToken(20, css.SemicolonToken, ";", false),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Parse(test.input, test.inline)

			if diff := cmp.Diff(test.expected, actual); diff != "" {
				t.Error(test.input)
				t.Error(diff)
			}
		})
	}
}

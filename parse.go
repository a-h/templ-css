package css

import (
	"bytes"
	"fmt"
	"io"

	"github.com/a-h/templ/parser/v2/goexpression"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

func NewToken(pos int, tt css.TokenType, content string, isGo bool) Token {
	return Token{
		Pos:       pos,
		TokenType: tt,
		Content:   content,
		Go:        isGo,
	}
}

type Token struct {
	Pos       int
	TokenType css.TokenType
	Content   string
	Go        bool
}

func peek(input string, pos int) (r rune) {
	if pos >= len(input) {
		return 0
	}
	return rune(input[pos])
}

const lbrace = '{'

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func Parse(input string, inline bool) (tokens []Token, err error) {
	pi := parse.NewInput(bytes.NewBufferString(input))
	p := css.NewParser(pi, inline)
	for {
		pos := p.Offset()
		gt, tt, data := p.Next()
		if gt == css.ErrorGrammar {
			if p.Err() == io.EOF {
				return tokens, nil
			}
			return nil, fmt.Errorf("failed to parse CSS: %w", p.Err())
		}

		// Skip tokens that can't contain Go expressions.
		if !(gt == css.AtRuleGrammar || gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar || gt == css.DeclarationGrammar) {
			tokens = append(tokens, NewToken(pos, tt, string(data), false))
			continue
		}

		// Handle at-rules and declarations.
		tokens = append(tokens, NewToken(pos, tt, string(data), false))
		pos += len(string(data))
		if gt == css.DeclarationGrammar {
			// Add a colon.
			tokens = append(tokens, NewToken(pos, css.ColonToken, ":", false))
			pos += len(":")
			// Collect whitespace.
			for wsPos := pos; wsPos < len(input); wsPos++ {
				peeked := peek(input, wsPos)
				if !isWhitespace(peeked) {
					break
				}
				tokens = append(tokens, NewToken(pos, css.WhitespaceToken, string(peeked), false))
				pos += len(string(peeked))
			}
		}

		// Handle values.
		skipUntil := -1
	values:
		for _, val := range p.Values() {
			if skipUntil > -1 && pos < skipUntil {
				pos += len(string(val.Data))
				continue values
			}
			skipUntil = -1

			next := peek(input, pos+1)
			if val.TokenType == css.LeftBraceToken && next == lbrace {
				// We've got a Go expression.
				// Skip the next character, which is the opening brace, then read a Go expression until we hit the closing brace.
				_, end, err := goexpression.Expression(input[pos+1:])
				if err != nil {
					return nil, fmt.Errorf("failed to parse Go expression: %w", err)
				}
				expr := input[pos : pos+end+2]
				tokens = append(tokens, NewToken(pos, css.IdentToken, expr, true))
				// Skip the closing brace.
				skipUntil = pos + 1 + len(expr)
				pos += len(string(val.Data))
				continue values
			}
			tokens = append(tokens, NewToken(pos, tt, string(val.Data), false))
			pos += len(string(val.Data))
			continue values
		}

		if gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar {
			tokens = append(tokens, NewToken(pos, css.LeftBraceToken, string(lbrace), false))
			pos += len(string(lbrace))
			continue
		}
		// Add a semicolon to the end of the declaration.
		if gt == css.AtRuleGrammar || gt == css.DeclarationGrammar {
			tokens = append(tokens, NewToken(pos, css.SemicolonToken, ";", false))
			pos += len(";")
			continue
		}
	}
}

func PrintTokens(tokens []Token) string {
	var b bytes.Buffer
	for _, t := range tokens {
		b.WriteString(t.Content)
	}
	return b.String()
}

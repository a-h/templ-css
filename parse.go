package css

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ/parser/v2/goexpression"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

func NewCSSToken(pos int, tt css.TokenType, content string) Token {
	return CSSToken{
		Position:  pos,
		TokenType: tt,
		Content:   content,
	}
}

type Token interface {
	Pos() int
	String() string
}

var _ Token = CSSToken{}
var _ Token = GoToken{}

type CSSToken struct {
	Position  int
	TokenType css.TokenType
	Content   string
}

func (t CSSToken) Pos() int {
	return t.Position
}

func (t CSSToken) String() string {
	return t.Content
}

func NewGoToken(pos int, prefix, expr, suffix string) Token {
	return GoToken{
		Position: pos,
		Prefix:   prefix,
		Expr:     expr,
		Suffix:   suffix,
	}
}

type GoToken struct {
	Position int
	Prefix   string
	Expr     string
	Suffix   string
}

func (t GoToken) Pos() int {
	return t.Position
}

func (t GoToken) String() string {
	return t.Prefix + t.Expr + t.Suffix
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

func readWhitespace(input string, pos int) string {
	var sb strings.Builder
	for wsPos := pos; wsPos < len(input); wsPos++ {
		peeked := peek(input, wsPos)
		if !isWhitespace(peeked) {
			break
		}
		sb.WriteRune(peeked)
	}
	return sb.String()
}

func isEmptyWhitespace(tt css.TokenType, data []byte) bool {
	return tt == css.WhitespaceToken && len(data) == 0
}

func Parse(input string, inline bool) (tokens []Token, err error) {
	pi := parse.NewInput(bytes.NewBufferString(input))

	p := css.NewParser(pi, inline)
	pos := p.Offset()
loop:
	for {
		gt, tt, data := p.Next()
		if gt == css.ErrorGrammar {
			if p.Err() == io.EOF {
				break loop
			}
			return nil, fmt.Errorf("failed to parse CSS: %w", p.Err())
		}

		if gt == css.AtRuleGrammar || gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar || gt == css.DeclarationGrammar {
			if !isEmptyWhitespace(tt, data) {
				tokens = append(tokens, NewCSSToken(pos, tt, string(data)))
				pos += len(data)
			}
			if gt == css.DeclarationGrammar {
				tokens = append(tokens, NewCSSToken(pos, css.ColonToken, string(":")))
				pos++
			}
			var ws string
			if ws = readWhitespace(input, pos); ws != "" {
				tokens = append(tokens, NewCSSToken(pos, css.WhitespaceToken, ws))
				pos += len(ws)
			}
			// Read values.
			skipValueToIndex := -1
			for i, val := range p.Values() {
				if skipValueToIndex > -1 && i < skipValueToIndex {
					continue
				}
				if i == 0 && val.TokenType == css.WhitespaceToken && ws != "" {
					// Skip leading whitespace, it's already added.
					continue
				}
				if goToken, isGoToken := getGoToken(input, pos); isGoToken {
					startPos := pos
					tokens = append(tokens, goToken)
					pos += len(goToken.String())
					// Work out how many upcoming values to skip.
					var j int
					for j = i + 1; j < len(p.Values()); j++ {
						startPos += len(p.Values()[j].Data)
						if startPos >= pos {
							break
						}
					}
					skipValueToIndex = j
					continue
				}
				tokens = append(tokens, NewCSSToken(pos, val.TokenType, string(val.Data)))
				pos += len(val.Data)
			}
			// Read whitespace between values and semicolon.
			if ws := readWhitespace(input, pos); ws != "" {
				tokens = append(tokens, NewCSSToken(pos, css.WhitespaceToken, ws))
				pos += len(ws)
			}
			// Add braces / semicolon.
			if gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar {
				tokens = append(tokens, NewCSSToken(pos, css.LeftBraceToken, "{"))
				pos++
			} else if gt == css.AtRuleGrammar || gt == css.DeclarationGrammar {
				tokens = append(tokens, NewCSSToken(pos, css.SemicolonToken, ";"))
				pos++
			} else if gt == css.EndAtRuleGrammar || gt == css.EndRulesetGrammar {
				tokens = append(tokens, NewCSSToken(pos, css.RightBraceToken, "}"))
				pos++
			}
			// Read whitespace after braces / semicolon.
			if ws := readWhitespace(input, pos); ws != "" {
				tokens = append(tokens, NewCSSToken(pos, css.WhitespaceToken, ws))
				pos += len(ws)
			}
		} else {
			ws := readWhitespace(input, pos)
			if ws != "" {
				tokens = append(tokens, NewCSSToken(pos, css.WhitespaceToken, ws))
				pos += len(ws)
			}
			tokens = append(tokens, NewCSSToken(pos, tt, string(data)))
			pos += len(data)
		}
	}

	return tokens, nil
}

func getGoToken(input string, pos int) (expr GoToken, isGo bool) {
	expr.Position = pos
	if peek(input, pos) != lbrace {
		return
	}
	if peek(input, pos+1) != lbrace {
		return
	}
	// Read prefix.
	expr.Prefix = "{{" + readWhitespace(input, pos+2)
	pos += len(expr.Prefix)
	// Read expression.
	start, end, err := goexpression.Expression(input[pos:])
	if err != nil {
		return expr, false
	}
	expr.Expr = input[pos+start : pos+end]
	pos += end
	// Read suffix.
	expr.Suffix = readWhitespace(input, pos) + "}}"
	pos += len(expr.Suffix)
	// Return.
	return expr, true
}

func PrintTokens(tokens []Token) string {
	var b bytes.Buffer
	for _, t := range tokens {
		b.WriteString(t.String())
	}
	return b.String()
}

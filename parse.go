package css

import (
	"bytes"
	"strings"

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

type goExprState int

const (
	goExptStateInitial goExprState = iota
	goExprStateLBraceOpen
	goExprStateLBrace
	goExprStateInside
	goExprStateRBrace
	goExprStateRBraceClose
)

func Parse(input string, inline bool) (tokens []Token) {
	p := css.NewParser(parse.NewInput(bytes.NewBufferString(input)), inline)
	for {
		pos := p.Offset()
		gt, tt, data := p.Next()
		if gt == css.ErrorGrammar {
			//TODO: Handle errors.
			break
		} else if gt == css.AtRuleGrammar || gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar || gt == css.DeclarationGrammar {
			tokens = append(tokens, NewToken(pos, tt, string(data), false))
			pos += len(string(data))
			if gt == css.DeclarationGrammar {
				tokens = append(tokens, NewToken(pos, css.ColonToken, ":", false))
				pos += len(":")
			}
			var braceDepth int
			var goExpr strings.Builder
			for i, val := range p.Values() {
				var prev css.Token
				if i > 0 {
					prev = p.Values()[i-1]
				}
				if val.TokenType == css.LeftBraceToken && prev.TokenType == css.LeftBraceToken {
					braceDepth++
					pos += len("{")
					continue
				}
				if val.TokenType == css.RightBraceToken && prev.TokenType == css.RightBraceToken {
					braceDepth--
					pos += len("}")
					if braceDepth == 0 {
						tokens = append(tokens, NewToken(pos, css.IdentToken, goExpr.String(), true))
						goExpr.Reset()
					}
					continue
				}
				if braceDepth == 2 {
					// Inside double braces, we're a Go expression.
					goExpr.Write(val.Data)
					continue
				}
				tokens = append(tokens, NewToken(pos, tt, string(val.Data), false))
				pos += len(string(val.Data))
			}
			if gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar {
				tokens = append(tokens, NewToken(pos, css.LeftBraceToken, "{", false))
				pos += len("{")
			} else if gt == css.AtRuleGrammar || gt == css.DeclarationGrammar {
				tokens = append(tokens, NewToken(pos, css.SemicolonToken, ";", false))
				pos += len(";")
			}
		} else {
			tokens = append(tokens, NewToken(pos, tt, string(data), false))
		}
	}
	return
}

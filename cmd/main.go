package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/alecthomas/kong"
	"github.com/bep/godartsass/v2"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type CLI struct {
	Generate GenerateCommand `cmd:"generate" help:"Generate a Go file from a CSS file."`
}

type GenerateCommand struct {
	FileName string `help:"File to process."`
	Package  string `help:"Package name." default:"main"`
}

func (c GenerateCommand) Run(ctx context.Context) (err error) {
	var r io.ReadCloser

	extension := filepath.Ext(c.FileName)
	switch extension {
	case ".scss":
		src, err := os.ReadFile(c.FileName)
		if err != nil {
			return fmt.Errorf("could not read file: %w", err)
		}
		css, err := convertScssToCss(string(src))
		if err != nil {
			return fmt.Errorf("could not convert scss to css: %w", err)
		}
		r = io.NopCloser(strings.NewReader(css))
	case ".css":
		r, err = os.Open(c.FileName)
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file type: %s", extension)
	}
	defer r.Close()

	lexer := css.NewLexer(parse.NewInput(r))
	classes := make(map[string]bool)
	var insideSelector bool

	for {
		tt, text := lexer.Next()
		if tt == css.ErrorToken {
			break
		}

		// '{' indicates selector end.
		if tt == css.LeftBraceToken {
			insideSelector = false
		}

		if tt == css.IdentToken && insideSelector && len(text) > 0 {
			classes[string(text)] = true
		}

		// '.' indicates selector start.
		if tt == css.DelimToken && text[0] == '.' { // Look for '.'
			insideSelector = true
		}
	}

	// Collect and sort class names
	classList := make([]string, len(classes))
	var i int
	for class := range classes {
		classList[i] = class
		i++
	}
	sort.Strings(classList)

	// Print results.
	fmt.Println(writeGoCode(c.Package, classList))

	return nil
}

var getSassTranspilerOnce func() (*godartsass.Transpiler, error) = sync.OnceValues(func() (*godartsass.Transpiler, error) {
	return godartsass.Start(godartsass.Options{})
})

func convertScssToCss(src string) (css string, err error) {
	tp, err := getSassTranspilerOnce()
	if err != nil {
		return "", err
	}
	result, err := tp.Execute(godartsass.Args{
		Source: src,
	})
	return result.CSS, err
}

func writeGoCode(pkg string, classes []string) string {
	var sb strings.Builder

	sb.WriteString("package ")
	sb.WriteString(pkg)
	sb.WriteString("\n\n")

	for _, class := range classes {
		sb.WriteString("const ")
		sb.WriteString(convertToGoName(class))
		sb.WriteString(" = ")
		sb.WriteString(strconv.Quote(class))
		sb.WriteByte('\n')
	}

	return sb.String()
}

func convertToGoName(s string) string {
	// A Go identifier must begin with a letter (a-z or A-Z) or an underscore (_) and can be followed by any combination of letters, digits (0-9), and underscores.

	// Replace invalid characters with underscores.
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r
		}
		if r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, s)
	if len(s) == 0 {
		return "_"
	}
	if s[0] >= '0' && s[0] <= '9' {
		s = "_" + s
	}

	return uppercaseFirstLetter(s)
}

func uppercaseFirstLetter(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}

func main() {
	var cli CLI
	ctx := context.Background()
	kctx := kong.Parse(&cli, kong.UsageOnError(), kong.BindTo(ctx, (*context.Context)(nil)))
	if err := kctx.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// Define a plugin interface for CSS handling.
type CSSPluginInput struct {
	FileName string
	Package  string
	CSS      string
}

type CSSPluginOutput struct {
	CSS    string
	GoCode string
}

type CSSPlugin interface {
	Process(input CSSPluginInput) (CSSPluginOutput, error)
}

var extensionToPlugin = map[string]CSSPlugin{
	// Converts SCSS to CSS and then generates Go code from the CSS.
	".scss": SCSSPlugin{},
	// Generates scoped CSS and Go code from CSS.
	".module.css": ModuleCSSPlugin{},
	// Generates Go code from CSS.
	".css": CSSCodeGenPlugin{},
}

func (c GenerateCommand) Run(ctx context.Context) (err error) {
	extension := filepath.Ext(c.FileName)
	if strings.HasSuffix(c.FileName, ".module.css") {
		extension = ".module.css"
	}

	plugin, ok := extensionToPlugin[extension]
	if !ok {
		return fmt.Errorf("unsupported file type: %s", extension)
	}

	cssBytes, err := os.ReadFile(c.FileName)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	pluginInput := CSSPluginInput{
		FileName: c.FileName,
		Package:  c.Package,
		CSS:      string(cssBytes),
	}
	pluginOutput, err := plugin.Process(pluginInput)
	if err != nil {
		return fmt.Errorf("could not process file: %w", err)
	}

	fmt.Println(pluginOutput.CSS)
	fmt.Println(pluginOutput.GoCode)

	return nil
}

// CSSCodeGenPlugin generates Go code from CSS files.
// It generates constants for each class name in the CSS file.
type CSSCodeGenPlugin struct {
}

func (p CSSCodeGenPlugin) Process(input CSSPluginInput) (output CSSPluginOutput, err error) {
	lexer := css.NewLexer(parse.NewInputString(input.CSS))
	classes := make(map[string]string)
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
			classes[string(text)] = string(text)
		}

		// '.' indicates selector start.
		if tt == css.DelimToken && text[0] == '.' { // Look for '.'
			insideSelector = true
		}
	}

	output.CSS = input.CSS
	output.GoCode = generateCode(input.Package, classes)

	return output, nil
}

type SCSSPlugin struct {
}

func (p SCSSPlugin) Process(input CSSPluginInput) (output CSSPluginOutput, err error) {
	tp, err := getSassTranspilerOnce()
	if err != nil {
		return output, err
	}
	result, err := tp.Execute(godartsass.Args{
		Source: input.CSS,
	})
	if err != nil {
		return output, fmt.Errorf("could not convert scss to css: %w", err)
	}
	output.CSS = result.CSS

	// Use the CSSCodeGenPlugin to generate Go code from the standard CSS output of the SCSS transpiler.
	plugin := CSSCodeGenPlugin{}
	return plugin.Process(CSSPluginInput{
		FileName: input.FileName,
		Package:  input.Package,
		CSS:      output.CSS,
	})
}

type ModuleCSSPlugin struct {
}

func (p ModuleCSSPlugin) Process(input CSSPluginInput) (output CSSPluginOutput, err error) {
	prefix := hex.EncodeToString(sha256.New().Sum([]byte(input.CSS))[0:16])
	var builder strings.Builder
	classNamesToPrefixedClassNames := make(map[string]string)

	lexer := css.NewLexer(parse.NewInputString(input.CSS))
	var insideSelector bool

	for {
		tt, text := lexer.Next()

		// Handle parsing errors explicitly
		if tt == css.ErrorToken {
			if lexer.Err() != nil && lexer.Err() != io.EOF {
				return output, fmt.Errorf("CSS parsing error: %v", lexer.Err())
			}
			break
		}

		// End of selector
		if tt == css.LeftBraceToken {
			insideSelector = false
		}

		// Prefix class names inside selectors
		if tt == css.IdentToken && insideSelector {
			newClassName := "templ_css_" + prefix + "_" + string(text)
			builder.WriteString(newClassName)
			classNamesToPrefixedClassNames[string(text)] = newClassName
		} else {
			builder.WriteString(string(text)) // Keep other content unchanged
		}

		// Detect class selectors
		if tt == css.DelimToken && text[0] == '.' {
			insideSelector = true
		}
	}

	output.CSS = builder.String()
	output.GoCode = generateCode(input.Package, classNamesToPrefixedClassNames)

	return output, nil
}

var getSassTranspilerOnce func() (*godartsass.Transpiler, error) = sync.OnceValues(func() (*godartsass.Transpiler, error) {
	return godartsass.Start(godartsass.Options{})
})

func generateCode(pkgs string, nameToPrefixedClassName map[string]string) string {
	var sb strings.Builder

	sb.WriteString("package ")
	sb.WriteString(pkgs)
	sb.WriteString("\n\n")

	classNames := make([]string, len(nameToPrefixedClassName))
	var i int
	for className := range nameToPrefixedClassName {
		classNames[i] = className
		i++
	}
	sort.Strings(classNames)

	for _, className := range classNames {
		goName := convertToGoName(className)
		prefixedClassName := nameToPrefixedClassName[className]
		sb.WriteString("const ")
		sb.WriteString(goName)
		sb.WriteString(" = ")
		sb.WriteString(strconv.Quote(prefixedClassName))
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

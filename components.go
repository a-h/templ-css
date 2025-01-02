package templcss

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/a-h/templ/safehtml"
)

type prop struct {
	Key   string
	Value string
}

func Var(name, value string) prop {
	k, v := safehtml.SanitizeCSS(name, value)
	return prop{
		Key:   k,
		Value: v,
	}
}

func UnsanitizedVar(name, value string) prop {
	return prop{
		Key:   name,
		Value: value,
	}
}

func Variables(properties ...prop) templ.Component {
	jsonString, err := templ.JSONString(properties)
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err != nil {
			return err
		}
		io.WriteString(w, "<script type=\"text/javascript\" data-variables=\"")
		io.WriteString(w, jsonString)
		io.WriteString(w, "\">\n")
		io.WriteString(w, "const props = JSON.parse(document.currentScript.getAttribute(\"data-variables\"));\n")
		io.WriteString(w, "props.forEach(p => document.documentElement.style.setProperty(p.name, p.value));\n")
		io.WriteString(w, "</script>\n")
		return nil
	})
}

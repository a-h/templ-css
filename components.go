package templcss

import (
	"context"
	"html"
	"io"

	"github.com/a-h/templ"
	"github.com/a-h/templ/safehtml"
)

type prop struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

func New(property, value string) prop {
	p, v := safehtml.SanitizeCSS(property, value)
	return prop{
		Property: p,
		Value:    v,
	}
}

func Unsanitized(property, value string) prop {
	return prop{
		Property: property,
		Value:    value,
	}
}

func Set(properties ...prop) templ.Component {
	jsonString, err := templ.JSONString(properties)
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err != nil {
			return err
		}
		io.WriteString(w, "<script type=\"text/javascript\" data-variables=\"")
		io.WriteString(w, html.EscapeString(jsonString))
		io.WriteString(w, "\">\n")
		io.WriteString(w, "const props = JSON.parse(document.currentScript.getAttribute(\"data-variables\"));\n")
		io.WriteString(w, "const r = document.querySelector(\":root\");\n")
		io.WriteString(w, "props.forEach(p => { r.style.setProperty(p.property, p.value); alert(JSON.stringify(p)) });\n")
		io.WriteString(w, "</script>\n")
		return nil
	})
}

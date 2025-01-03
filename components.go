package templcss

import (
	"context"
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

		script := `<script type="text/javascript">
			(() => {
				const r = document.querySelector(":root");
				const props = ` + jsonString + `;
				props.forEach(({ property, value }) => r.style.setProperty(property, value));
			})();
		</script>`

		_, err = io.WriteString(w, script)
		return err
	})
}

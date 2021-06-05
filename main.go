package main

import (
	"fmt"

	"github.com/a-h/templ-css/css"
)

// It's possible to register global rules and classes at startup.
func init() {
	css.RegisterGlobalRule("h1", css.FontSize(css.FontSizeLarge))
	css.RegisterGlobalClass("name", css.Color(css.ColorRed), css.BackgroundColor(css.RGB(255, 99, 71)))
	css.RegisterGlobalClass("button", css.Color(css.ColorBlack), css.BackgroundColor(css.ColorWhiteSmoke))
}

// Templ components can register classes using css.Class, which returns a unique ID for the CSS class for this application.

// At the moment, it returns component_0 which could be minified further to `c0` (class, and the zero-th component), or perhaps
// to a partial hash to make it more unique (to avoid problems with stale CSS?).

// See componentstyles.go for an example of this. It's used in `component.templ`. It uses `string(componentClass)` because templ
// doesn't have the concept of a "css.ClassID" yet.

func main() {
	// This approach allows a global stylesheet to be rendered.
	fmt.Println(css.Stylesheet())

	// By updating css.StyleSheet to be a `css.RenderStylesheet` function that takes a context, it can update the
	// context to track that the stylesheet has been output to the HTTP response.

	// Thay way, when templ renders an element, it can check whether the class has been given to the client. If not...
	// then it can apply appropriate `style=""` attributes.

	// This means that templ.Components can carry styles with them, _and_ benefit from minified CSS, like vue.js or
	// CSS-in-js with React.

	// If you run the example, you'll also see that the result is minified CSS.
}

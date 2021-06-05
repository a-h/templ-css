package main

import "github.com/a-h/templ-css/css"

// This class is only used by the `component.templ`.
// The rendering code could work out whether the global stylesheet has been rendered or not..
// If it has, then the component would just use the class, if not, it would populate the style attribute with the required properties.
var componentClass = css.Class("component", css.BackgroundColor(css.ColorWhite))

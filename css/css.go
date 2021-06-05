package css

import (
	"fmt"
	"strings"
)

// styles is a singleton of all style rules that templ knows about.
var styles = []Rule{}

// RegisterGlobalRule defines a CSS rule that applies to the selector.
func RegisterGlobalRule(selector string, properties ...Declaration) {
	styles = append(styles, Rule{
		Selector:     selector,
		Declarations: properties,
	})
}

// RegisterGlobalClass defines a CSS class with a specific name.
func RegisterGlobalClass(name string, properties ...Declaration) {
	if !strings.HasPrefix(name, ".") {
		name = "." + name
	}
	styles = append(styles, Rule{
		Selector:     name,
		Declarations: properties,
	})
}

// Stylesheet returns the global CSS stylesheet.
func Stylesheet() string {
	var sb strings.Builder
	for i := 0; i < len(styles); i++ {
		r := styles[i]
		sb.WriteString(r.String())
	}
	return sb.String()
}

// Rule applied to the document.
type Rule struct {
	Selector     string
	Declarations []Declaration
}

// String CSS representation of the rule.
func (r Rule) String() string {
	//TODO: Should write this directly to a writer instead of to an internal builder.
	var sb strings.Builder
	sb.WriteString(r.Selector)
	sb.WriteRune('{')
	for i := 0; i < len(r.Declarations); i++ {
		d := r.Declarations[i]
		sb.WriteString(d.Property)
		sb.WriteRune(':')
		sb.WriteString(d.Value)
		if i < len(r.Declarations)-1 {
			sb.WriteRune(';')
		}
	}
	sb.WriteRune('}')
	return sb.String()
}

// Declaration is a CSS property and value pair.
type Declaration struct {
	Property string
	Value    string
}

// ClassID is a unique ID representing a "local" class used by a component.
type ClassID string

// Class creates a unqiue class name and registers it with the global styles.
func Class(name string, properties ...Declaration) (id ClassID) {
	id = ClassID(fmt.Sprintf("%s_%d", name, len(styles)))
	r := Rule{
		Selector:     "." + string(id),
		Declarations: properties,
	}
	styles = append(styles, r)
	return
}

// Property is used when the css package doesn't have a helper function yet.
// This would need to validate the values passed in (for security, e.g. remote URLs).
func Property(name, value string) Declaration {
	return Declaration{
		Property: name,
		Value:    value,
	}
}

// All the standard CSS properties would be here, with helpers.

// BackgroundColor of the element.
func BackgroundColor(v ColorValue) Declaration {
	return Declaration{
		Property: "background-color",
		Value:    string(v),
	}
}

// Color of the element.
func Color(v ColorValue) Declaration {
	return Declaration{
		Property: "color",
		Value:    string(v),
	}
}

type ColorValue string

func RGB(r, g, b int) ColorValue {
	//TODO: Be smart enough to work out that 'gold' is shorter than #FFD700, and use that instead.
	return ColorValue(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}

const (
	ColorLightSalmon          ColorValue = "#ffa07a" // #FFA07A
	ColorSalmon                          = "salmon"  // #FA8072
	ColorDarkSalmon                      = "#e9967a" // #E9967A
	ColorLightCoral                      = "#f08080" // #F08080
	ColorIndianRed                       = "#cd5c5c" // #CD5C5C
	ColorCrimson                         = "#dc143c" // #DC143C
	ColorFireBrick                       = "#b22222" // #B22222
	ColorRed                             = "red"     // #FF0000
	ColorDarkRed                         = "#8b0000" // #8B0000
	ColorCoral                           = "coral"   // #FF7F50
	ColorTomato                          = "tomato"  // #FF6347
	ColorOrangeRed                       = "#ff4500" // #FF4500
	ColorGold                            = "gold"    // #FFD700
	ColorOrange                          = "orange"  // #FFA500
	ColorDarkorange                      = "#ff8c00" // #FF8C00
	ColorLightYellow                     = "#ffffe0" // #FFFFE0
	ColorLemonChiffon                    = "#fffacd" // #FFFACD
	ColorLightGoldenRodYellow            = "#fafad2" // #FAFAD2
	ColorPapayaWhip                      = "#ffefd5" // #FFEFD5
	ColorMoccasin                        = "#ffe4b5" // #FFE4B5
	ColorPeachPuff                       = "#ffdab9" // #FFDAB9
	ColorPaleGoldenRod                   = "#eee8aa" // #EEE8AA
	ColorKhaki                           = "khaki"   // #F0E68C
	ColorDarkKhaki                       = "#bdb76b" // #BDB76B
	ColorYellow                          = "yellow"  // #FFFF00
	ColorLawnGreen                       = "#7cfc00" // #7CFC00
	ColorChartreuse                      = "#7fff00" // #7FFF00
	ColorLimeGreen                       = "#32cd32" // #32CD32
	ColorLime                            = "lime"    // #00FF00
	ColorForestGreen                     = "#228b22" // #228B22
	ColorGreen                           = "green"   // #008000
	ColorDarkGreen                       = "#006400" // #006400
	ColorGreenYellow                     = "#adff2f" // #ADFF2F
	ColorYellowGreen                     = "#9acd32" // #9ACD32
	ColorSpringGreen                     = "#00ff7f" // #00FF7F
	ColorMediumSpringGreen               = "#00fa9a" // #00FA9A
	ColorLightGreen                      = "#90ee90" // #90EE90
	ColorPaleGreen                       = "#98fb98" // #98FB98
	ColorDarkSeaGreen                    = "#8fbc8f" // #8FBC8F
	ColorMediumSeaGreen                  = "#3cb371" // #3CB371
	ColorSeaGreen                        = "#2e8b57" // #2E8B57
	ColorOlive                           = "olive"   // #808000
	ColorDarkOliveGreen                  = "#556b2f" // #556B2F
	ColorOliveDrab                       = "#6b8e23" // #6B8E23
	ColorLightcyan                       = "#e0ffff" // #E0FFFF
	ColorCyan                            = "cyan"    // #00FFFF
	ColorAqua                            = "aqua"    // #00FFFF
	ColorAquamarine                      = "#7fffd4" // #7FFFD4
	ColorMediumaquamarine                = "#66cdaa" // #66CDAA
	ColorPaleTurquoise                   = "#afeeee" // #AFEEEE
	ColorTurquoise                       = "#40e0d0" // #40E0D0
	ColorMediumTurquoise                 = "#48d1cc" // #48D1CC
	ColorDarkTurquoise                   = "#00ced1" // #00CED1
	ColorLightSeaGreen                   = "#20b2aa" // #20B2AA
	ColorCadetBlue                       = "#5f9ea0" // #5F9EA0
	ColorDarkCyan                        = "#008b8b" // #008B8B
	ColorTeal                            = "teal"    // #008080
	ColorPowderBlue                      = "#b0e0e6" // #B0E0E6
	ColorLightBlue                       = "#add8e6" // #ADD8E6
	ColorLightSkyBlue                    = "#87cefa" // #87CEFA
	ColorSkyBlue                         = "#87ceeb" // #87CEEB
	ColorDeepSkyBlue                     = "#00bfff" // #00BFFF
	ColorLightsteelBlue                  = "#b0c4de" // #B0C4DE
	ColorDodgerBlue                      = "#1e90ff" // #1E90FF
	ColorCornflowerBlue                  = "#6495ed" // #6495ED
	ColorSteelBlue                       = "#4682b4" // #4682B4
	ColorRoyalBlue                       = "#4169e1" // #4169E1
	ColorBlue                            = "blue"    // #0000FF
	ColorMediumBlue                      = "#0000cd" // #0000CD
	ColorDarkBlue                        = "#00008b" // #00008B
	ColorNavy                            = "navy"    // #000080
	ColorMidnightBlue                    = "#191970" // #191970
	ColorMediumSlateBlue                 = "#7b68ee" // #7B68EE
	ColorSlateBlue                       = "#6a5acd" // #6A5ACD
	ColorDarkSlateBlue                   = "#483d8b" // #483D8B
	ColorLavender                        = "#e6e6fa" // #E6E6FA
	ColorThistle                         = "#d8bfd8" // #D8BFD8
	ColorPlum                            = "plum"    // #DDA0DD
	ColorViolet                          = "violet"  // #EE82EE
	ColorOrchid                          = "orchid"  // #DA70D6
	ColorFuchsia                         = "#ff00ff" // #FF00FF
	ColorMagenta                         = "#ff00ff" // #FF00FF
	ColorMediumOrchid                    = "#ba55d3" // #BA55D3
	ColorMediumPurple                    = "#9370db" // #9370DB
	ColorBlueViolet                      = "#8a2be2" // #8A2BE2
	ColorDarkViolet                      = "#9400d3" // #9400D3
	ColorDarkOrchid                      = "#9932cc" // #9932CC
	ColorDarkMagenta                     = "#8b008b" // #8B008B
	ColorPurple                          = "purple"  // #800080
	ColorIndigo                          = "indigo"  // #4B0082
	ColorPink                            = "pink"    // #FFC0CB
	ColorLightPink                       = "#ffb6c1" // #FFB6C1
	ColorHotPink                         = "#ff69b4" // #FF69B4
	ColorDeepPink                        = "#ff1493" // #FF1493
	ColorPaleVioletRed                   = "#db7093" // #DB7093
	ColorMediumVioletRed                 = "#c71585" // #C71585
	ColorWhite                           = "white"   // #FFFFFF
	ColorSnow                            = "snow"    // #FFFAFA
	ColorHoneydew                        = "#f0fff0" // #F0FFF0
	ColorMintCream                       = "#f5fffa" // #F5FFFA
	ColorAzure                           = "azure"   // #F0FFFF
	ColorAliceBlue                       = "#f0f8ff" // #F0F8FF
	ColorGhostWhite                      = "#f8f8ff" // #F8F8FF
	ColorWhiteSmoke                      = "#f5f5f5" // #F5F5F5
	ColorSeashell                        = "#fff5ee" // #FFF5EE
	ColorBeige                           = "beige"   // #F5F5DC
	ColorOldLace                         = "#fdf5e6" // #FDF5E6
	ColorFloralWhite                     = "#fffaf0" // #FFFAF0
	ColorIvory                           = "ivory"   // #FFFFF0
	ColorAntiqueWhite                    = "#faebd7" // #FAEBD7
	ColorLinen                           = "linen"   // #FAF0E6
	ColorLavenderBlush                   = "#fff0f5" // #FFF0F5
	ColorMistyRose                       = "#ffe4e1" // #FFE4E1
	ColorGainsboro                       = "#dcdcdc" // #DCDCDC
	ColorLightGray                       = "#d3d3d3" // #D3D3D3
	ColorSilver                          = "silver"  // #C0C0C0
	ColorDarkGray                        = "#a9a9a9" // #A9A9A9
	ColorGray                            = "gray"    // #808080
	ColorDimGray                         = "#696969" // #696969
	ColorLightSlateGray                  = "#778899" // #778899
	ColorSlateGray                       = "#708090" // #708090
	ColorDarkslateGray                   = "#2f4f4f" // #2F4F4F
	ColorBlack                           = "black"   // #000000 //TODO: Some names have even shorter hex, e.g. #000
	ColorCornsilk                        = "#fff8dc" // #FFF8DC
	ColorBlanchedAlmond                  = "#ffebcd" // #FFEBCD
	ColorBisque                          = "bisque"  // #FFE4C4
	ColorNavajoWhite                     = "#ffdead" // #FFDEAD
	ColorWheat                           = "wheat"   // #F5DEB3
	ColorBurlywood                       = "#deb887" // #DEB887
	ColorTan                             = "tan"     // #D2B48C
	ColorRosyBrown                       = "#bc8f8f" // #BC8F8F
	ColorSandyBrown                      = "#f4a460" // #F4A460
	ColorGoldenRod                       = "#daa520" // #DAA520
	ColorPeru                            = "peru"    // #CD853F
	ColorChocolate                       = "#d2691e" // #D2691E
	ColorSaddleBrown                     = "#8b4513" // #8B4513
	ColorSienna                          = "sienna"  // #A0522D
	ColorBrown                           = "brown"   // #A52A2A
	ColorMaroon                          = "maroon"  // #800000
)

func FontSize(v FontSizeValue) Declaration {
	return Declaration{
		Property: "color",
		Value:    string(v),
	}
}

type FontSizeValue string

func FontSizePoint(pt int) FontSizeValue {
	return FontSizeValue(fmt.Sprintf("%dpt", pt))
}

const (
	FontSizeXXSmall FontSizeValue = "xx-small"
	FontSizeXSmall                = "x-small"
	FontSizeSmall                 = "small"
	FontSizeMedium                = "medium"
	FontSizeLarge                 = "large"
	FontSizeXLarge                = "x-large"
	FontSizeXXLarge               = "xx-large"
)

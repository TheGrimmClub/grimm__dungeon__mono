// Package assets embeds static text assets (ASCII art, banners) into the binary
// so grimm ships as a single self-contained executable (req R010).
package assets

import _ "embed"

// Banner is the GRIMM ASCII title shown at startup.
//
//go:embed banner.txt
var Banner string

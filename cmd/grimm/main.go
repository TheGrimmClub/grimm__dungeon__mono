// Command grimm is TheGrimmClub's interactive learning dungeon. This entry
// point stays intentionally thin: all wiring lives in internal/app.
package main

import (
	"fmt"
	"os"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "grimm:", err)
		os.Exit(1)
	}
}

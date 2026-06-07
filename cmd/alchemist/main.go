// Command alchemist is the standalone "git as potion-brewing" tool students use
// to version their work. It is a thin entry point over internal/alchemist, the
// same library grimm exposes in-game as "/alchemist" (decision D006). It brews
// in the current working directory.
package main

import (
	"fmt"
	"os"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/alchemist"
)

func main() {
	out, err := alchemist.Dispatch(alchemist.New("."), os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "alchemist:", err)
		os.Exit(1)
	}
	fmt.Println(out)
}

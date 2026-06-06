// Package content embeds the dungeon's authored data (multi-document YAML) so
// grimm ships as a single binary (req R010). The game loads it via an fs.FS,
// keeping the loader decoupled from where the files live.
package content

import "embed"

// FS holds all authored content. The world loader globs "world/*.yaml".
//
//go:embed world/*.yaml
var FS embed.FS

// WorldGlob is the glob the world loader uses against FS.
const WorldGlob = "world/*.yaml"

// Package puzzle is the heart of the curriculum: a puzzle is a prompt plus an
// acceptance Check, and a Check is one of three kinds (decision D003):
//
//   - answer     — a riddle: normalized text or regex match.
//   - artifact   — proof of a real-world action: a file exists, or a command
//     (e.g. `git ls-remote <url>`) succeeds with expected output.
//   - behavioral — "choose your path": run the student's solution (Python, Go,
//     or a CLI one-liner) against I/O cases; the language is theirs.
//
// Everything in the game that gates progress (locked doors, loot, battles)
// hangs off this abstraction.
package puzzle

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Result is the outcome of a check. Detail is optional, check-specific feedback
// (in German) the engine may show alongside the puzzle's own hint.
type Result struct {
	Passed bool
	Detail string
}

// Input is what the player supplied plus the context a check may need.
type Input struct {
	Answer  string // the raw text the player gave to `solve`
	WorkDir string // the student's working directory (artifact/behavioral)
}

// Check verifies an attempt. Implementations must be safe to call repeatedly.
type Check interface {
	Verify(ctx context.Context, in Input) Result
}

// Spec is the YAML form of a check, embedded in a puzzle document.
type Spec struct {
	Kind string `yaml:"kind"` // answer | artifact | behavioral

	// answer
	Answer string   `yaml:"answer"`
	Accept []string `yaml:"accept"` // extra accepted answers
	Regex  string   `yaml:"regex"`

	// artifact (all configured conditions must hold)
	File    string `yaml:"file"`    // a file that must exist (relative to WorkDir)
	Command string `yaml:"command"` // shell command; {answer}/{workdir} expand
	Expect  string `yaml:"expect"`  // substring required in the command's stdout

	// behavioral
	Run   string   `yaml:"run"`   // command template; {answer}/{workdir} expand
	Cases []IOCase `yaml:"cases"` // stdin -> expected stdout
}

// IOCase is one input/output example for a behavioral check.
type IOCase struct {
	In  string `yaml:"in"`
	Out string `yaml:"out"`
}

// Build constructs the Check described by a Spec.
func Build(s Spec) (Check, error) {
	switch s.Kind {
	case "answer":
		return newAnswerCheck(s)
	case "artifact":
		return artifactCheck{file: s.File, command: s.Command, expect: s.Expect}, nil
	case "behavioral":
		if len(s.Cases) == 0 {
			return nil, fmt.Errorf("puzzle: behavioral check needs at least one case")
		}
		return behavioralCheck{run: s.Run, cases: s.Cases}, nil
	default:
		return nil, fmt.Errorf("puzzle: unknown check kind %q", s.Kind)
	}
}

// --- answer ---

type answerCheck struct {
	accepted []string // normalized
	re       *regexp.Regexp
}

func newAnswerCheck(s Spec) (Check, error) {
	c := answerCheck{}
	if s.Answer != "" {
		c.accepted = append(c.accepted, normalize(s.Answer))
	}
	for _, a := range s.Accept {
		c.accepted = append(c.accepted, normalize(a))
	}
	if s.Regex != "" {
		re, err := regexp.Compile("(?i)" + s.Regex)
		if err != nil {
			return nil, fmt.Errorf("puzzle: bad answer regex: %w", err)
		}
		c.re = re
	}
	if len(c.accepted) == 0 && c.re == nil {
		return nil, fmt.Errorf("puzzle: answer check needs `answer`, `accept` or `regex`")
	}
	return c, nil
}

func (c answerCheck) Verify(_ context.Context, in Input) Result {
	if c.re != nil && c.re.MatchString(strings.TrimSpace(in.Answer)) {
		return Result{Passed: true}
	}
	got := normalize(in.Answer)
	for _, a := range c.accepted {
		if a == got {
			return Result{Passed: true}
		}
	}
	return Result{Passed: false}
}

// normalize lower-cases, collapses whitespace and trims surrounding punctuation
// so "Der Rabe!" and "rabe" match.
func normalize(s string) string {
	s = strings.ToLower(strings.Join(strings.Fields(s), " "))
	return strings.Trim(s, " .,!?;:»«\"'()")
}

// expand substitutes {answer} and {workdir} placeholders in a command template.
func expand(tmpl string, in Input) string {
	r := strings.NewReplacer("{answer}", in.Answer, "{workdir}", in.WorkDir)
	return r.Replace(tmpl)
}

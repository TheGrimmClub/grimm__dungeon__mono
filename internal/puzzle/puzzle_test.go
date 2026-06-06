package puzzle

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func verify(t *testing.T, s Spec, in Input) Result {
	t.Helper()
	c, err := Build(s)
	if err != nil {
		t.Fatalf("Build(%+v): %v", s, err)
	}
	return c.Verify(context.Background(), in)
}

func TestAnswerCheck(t *testing.T) {
	s := Spec{Kind: "answer", Answer: "der Rabe", Accept: []string{"rabe"}}

	for _, ans := range []string{"der Rabe", "  RABE ", "Rabe!", "rabe"} {
		if !verify(t, s, Input{Answer: ans}).Passed {
			t.Errorf("answer %q should pass", ans)
		}
	}
	if verify(t, s, Input{Answer: "Taube"}).Passed {
		t.Error("wrong answer should fail")
	}
}

func TestAnswerCheckRegex(t *testing.T) {
	s := Spec{Kind: "answer", Regex: `^4\s*2$`}
	if !verify(t, s, Input{Answer: "42"}).Passed {
		t.Error("regex should match 42")
	}
	if verify(t, s, Input{Answer: "24"}).Passed {
		t.Error("regex should not match 24")
	}
}

func TestArtifactCheckFile(t *testing.T) {
	dir := t.TempDir()
	s := Spec{Kind: "artifact", File: "README.md"}

	if verify(t, s, Input{WorkDir: dir}).Passed {
		t.Error("should fail before the file exists")
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !verify(t, s, Input{WorkDir: dir}).Passed {
		t.Error("should pass once the file exists")
	}
}

func TestArtifactCheckCommand(t *testing.T) {
	// A command that must succeed and emit an expected substring; {answer} is the
	// "proof" the player supplies.
	s := Spec{Kind: "artifact", Command: `echo {answer}`, Expect: "grimm"}
	if !verify(t, s, Input{Answer: "hello-grimm-world"}).Passed {
		t.Error("command artifact should pass when stdout contains the expected text")
	}
	if verify(t, s, Input{Answer: "nope"}).Passed {
		t.Error("command artifact should fail when stdout lacks the expected text")
	}
}

func TestBehavioralCheckShell(t *testing.T) {
	// Doubling puzzle: the answer is a shell one-liner; runs against I/O cases.
	s := Spec{Kind: "behavioral", Cases: []IOCase{
		{In: "", Out: "42"},
	}}
	if !verify(t, s, Input{Answer: "echo 42"}).Passed {
		t.Error("correct shell solution should pass")
	}
	if res := verify(t, s, Input{Answer: "echo 41"}); res.Passed {
		t.Error("wrong output should fail")
	} else if res.Detail == "" {
		t.Error("failing case should explain expected vs got")
	}
}

func TestBehavioralCheckStdin(t *testing.T) {
	// run template reads stdin and doubles it.
	s := Spec{Kind: "behavioral",
		Run:   `read n; echo $((n*2))`,
		Cases: []IOCase{{In: "5\n", Out: "10"}, {In: "21\n", Out: "42"}},
	}
	if !verify(t, s, Input{}).Passed {
		t.Error("doubling run should pass both cases")
	}
}

func TestBuildRejectsUnknownKind(t *testing.T) {
	if _, err := Build(Spec{Kind: "wizardry"}); err == nil {
		t.Error("unknown kind should error")
	}
}

package runner

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRunCapturesStdout(t *testing.T) {
	res := Run(context.Background(), "", "", 0, "sh", "-c", "echo hallo welt")
	if res.Err != nil {
		t.Fatalf("Run error: %v", res.Err)
	}
	if strings.TrimSpace(res.Stdout) != "hallo welt" {
		t.Errorf("stdout = %q, want 'hallo welt'", res.Stdout)
	}
	if res.ExitCode != 0 {
		t.Errorf("exit = %d, want 0", res.ExitCode)
	}
}

func TestRunFeedsStdin(t *testing.T) {
	res := Run(context.Background(), "", "5\n", 0, "sh", "-c", "read n; echo $((n*2))")
	if got := strings.TrimSpace(res.Stdout); got != "10" {
		t.Errorf("stdout = %q, want 10", got)
	}
}

func TestRunNonZeroExit(t *testing.T) {
	res := Run(context.Background(), "", "", 0, "sh", "-c", "exit 3")
	if res.Err != nil {
		t.Errorf("unexpected start error: %v", res.Err)
	}
	if res.ExitCode != 3 {
		t.Errorf("exit = %d, want 3", res.ExitCode)
	}
}

func TestRunTimeout(t *testing.T) {
	res := Run(context.Background(), "", "", 200*time.Millisecond, "sh", "-c", "sleep 5")
	if !res.TimedOut {
		t.Errorf("expected TimedOut, got %+v", res)
	}
}

func TestRunMissingBinary(t *testing.T) {
	res := Run(context.Background(), "", "", 0, "definitely-not-a-real-binary-xyz")
	if res.Err == nil {
		t.Error("expected a start error for a missing binary")
	}
}

func TestAvailable(t *testing.T) {
	if !Available("sh") {
		t.Error("sh should be available")
	}
	if Available("definitely-not-a-real-binary-xyz") {
		t.Error("bogus binary should not be available")
	}
}

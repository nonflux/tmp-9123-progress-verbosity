//go:build replay

// Replays a recorded Claude stream-json file through fullsend's own
// production parser and renderer at one progress level, so a GitHub
// Actions run page shows exactly what `fullsend run` prints at that level.
// The workflow copies this file into internal/runtime of the fullsend
// checkout under test. It is not a live agent run.
package runtime

import (
	"os"
	"testing"

	"github.com/fullsend-ai/fullsend/internal/config"
	"github.com/fullsend-ai/fullsend/internal/ui"
)

func TestReplayProgress(t *testing.T) {
	path := os.Getenv("REPLAY_STREAM")
	level := os.Getenv("REPLAY_LEVEL") // "" = key unset
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open stream: %v", err)
	}
	defer f.Close()

	// Same wiring as ClaudeRuntime.Run: the default renderer at the level.
	renderer := newEventRendererAt(ui.New(os.Stdout), level)
	if err := parseClaudeStream(f, renderer.Handle); err != nil {
		t.Fatalf("parse stream: %v", err)
	}

	// A transcript error as the runner reports one after a failed run,
	// emitted the way internal/cli/run.go does at each level. The message
	// stands in for agent-written text.
	errs := []TranscriptError{{
		Source:       "triage-replay.jsonl",
		IsError:      true,
		ErrorMessage: "REPLAY-AGENT-TEXT: the agent quoted the issue here",
		Subtype:      "error_during_execution",
	}}
	if config.ProgressIsQuiet(level) {
		EmitQuietTranscriptErrors(os.Stderr, errs)
	} else {
		ClaudeRuntime{}.EmitTranscriptErrors(os.Stderr, errs)
	}
}

package audio_test

import (
	"testing"

	"github.com/dkmaker/claude-speak/internal/audio"
)

func TestNewPlayer(t *testing.T) {
	// Test that we can create and close a player
	p, err := audio.NewPlayer()
	if err != nil {
		t.Fatalf("NewPlayer: %v", err)
	}
	p.Close()
}

func TestPlayInvalidMP3(t *testing.T) {
	p, err := audio.NewPlayer()
	if err != nil {
		t.Fatalf("NewPlayer: %v", err)
	}
	defer p.Close()

	err = p.PlayMP3([]byte("not valid mp3 data"))
	if err == nil {
		t.Error("expected error for invalid MP3 data")
	}
}

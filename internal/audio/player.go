package audio

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

// Player manages audio playback using oto.
type Player struct {
	ctx *oto.Context
}

// NewPlayer creates an audio player. Call Close when done.
func NewPlayer() (*Player, error) {
	// go-mp3 always outputs: 16-bit signed LE, stereo (2 channels)
	// Sample rate varies per MP3 but 44100 is most common.
	// We'll create context on first play to match the actual sample rate.
	return &Player{}, nil
}

// ensureContext creates or reuses the oto context for the given sample rate.
func (p *Player) ensureContext(sampleRate int) error {
	if p.ctx != nil {
		return nil
	}
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 2, // go-mp3 always outputs stereo
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return fmt.Errorf("create audio context: %w", err)
	}
	<-ready
	p.ctx = ctx
	return nil
}

// PlayMP3 decodes MP3 data and plays it synchronously. Returns when playback is complete.
func (p *Player) PlayMP3(data []byte) error {
	decoder, err := mp3.NewDecoder(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode mp3: %w", err)
	}

	if err := p.ensureContext(decoder.SampleRate()); err != nil {
		return err
	}

	player := p.ctx.NewPlayer(decoder)
	defer player.Close()

	// Start playback (NewPlayer already wraps the io.Reader decoder)
	player.Play()

	// Wait for playback to complete
	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

// Close releases audio resources.
func (p *Player) Close() {
	// oto.Context doesn't have a Close method — it lives for the process lifetime.
	// This is fine because the worker daemon owns the player.
}

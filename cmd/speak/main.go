package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dkmaker/claude-speak/internal/audio"
	"github.com/dkmaker/claude-speak/internal/elevenlabs"
	"github.com/dkmaker/claude-speak/internal/queue"
	"github.com/dkmaker/claude-speak/internal/worker"
)

var Version = "dev"

func main() {
	log.SetFlags(0)
	log.SetPrefix("speak: ")

	if len(os.Args) < 2 {
		os.Exit(0) // No args, no-op (matches current speak.sh behavior)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("get home directory: %v", err)
	}
	ttsDir := filepath.Join(home, ".claude", "tts")
	os.MkdirAll(ttsDir, 0755)

	switch os.Args[1] {
	case "--version":
		fmt.Println(Version)
	case "--daemon":
		runDaemon(ttsDir)
	case "--stop":
		if err := worker.StopByPID(ttsDir); err != nil {
			log.Fatalf("stop worker: %v", err)
		}
	default:
		// Client mode: enqueue message and ensure daemon is running
		message := strings.Join(os.Args[1:], " ")
		enqueueAndEnsureDaemon(ttsDir, message)
	}
}

func enqueueAndEnsureDaemon(ttsDir, message string) {
	// Check for API key before spawning daemon
	if os.Getenv("ELEVENLABS_API_KEY") == "" {
		fmt.Fprintln(os.Stderr, "Error: ELEVENLABS_API_KEY environment variable is not set")
		os.Exit(1)
	}

	q := queue.New(ttsDir)
	if err := q.Enqueue(message); err != nil {
		log.Fatalf("enqueue: %v", err)
	}

	// Start daemon if not running
	if _, running := worker.IsRunning(ttsDir); !running {
		exe, err := os.Executable()
		if err != nil {
			log.Fatalf("find executable: %v", err)
		}
		cmd := exec.Command(exe, "--daemon")
		cmd.Stdout = nil
		cmd.Stderr = nil
		// Detach from parent process
		detachProcess(cmd)
		if err := cmd.Start(); err != nil {
			log.Fatalf("start daemon: %v", err)
		}
		// Don't wait — daemon runs in background
	}
}

func runDaemon(ttsDir string) {
	// Redirect stderr to log file
	logFile := filepath.Join(ttsDir, "speak.log")
	logF, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		os.Stderr = logF
		defer logF.Close()
	}

	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		log.Fatal("ELEVENLABS_API_KEY not set")
	}

	// Create API client
	opts := []elevenlabs.Option{}
	if v := os.Getenv("ELEVENLABS_VOICE_ID"); v != "" {
		opts = append(opts, elevenlabs.WithVoiceID(v))
	}
	if m := os.Getenv("ELEVENLABS_MODEL"); m != "" {
		opts = append(opts, elevenlabs.WithModel(m))
	}
	client := elevenlabs.NewClient(apiKey, opts...)

	// Create audio player
	player, err := audio.NewPlayer()
	if err != nil {
		log.Fatalf("create audio player: %v", err)
	}
	defer player.Close()

	// Create speak function
	speakFn := func(text string) error {
		data, err := client.Synthesize(text)
		if err != nil {
			return fmt.Errorf("synthesize: %w", err)
		}
		if err := player.PlayMP3(data); err != nil {
			return fmt.Errorf("play: %w", err)
		}
		return nil
	}

	// Run worker
	q := queue.New(ttsDir)
	w := worker.New(q, speakFn)
	w.Run()
}


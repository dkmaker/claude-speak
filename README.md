# claude-speak

Cross-platform TTS binary for Claude Code. Claude speaks progress updates aloud while you work, using the ElevenLabs text-to-speech API.

Supports **Linux**, **macOS**, and **Windows**.

## Installation

### As a Claude Code plugin

Install from the [my-claude-plugins](https://github.com/dkmaker/my-claude-plugins) marketplace:

```
/plugin install my-claude-plugins/claude-speak
```

The plugin automatically downloads the correct binary for your platform on first session start.

### Manual installation

Download the binary for your platform from [Releases](https://github.com/dkmaker/claude-speak/releases/latest):

| Platform | Binary |
|----------|--------|
| Linux x86_64 | `speak-linux-amd64` |
| macOS Apple Silicon | `speak-darwin-arm64` |
| Windows x86_64 | `speak-windows-amd64.exe` |

Place it somewhere in your PATH and rename to `speak` (or `speak.exe` on Windows).

## ElevenLabs API Key

You need an ElevenLabs API key for text-to-speech.

### 1. Create an account

Sign up at [elevenlabs.io](https://elevenlabs.io). The free tier includes a generous monthly character quota.

### 2. Get your API key

1. Log in to [elevenlabs.io](https://elevenlabs.io)
2. Click your profile icon (bottom-left)
3. Select **API Keys**
4. Click **Create API Key** and copy it

### 3. Set the environment variable

Add to your shell profile (`~/.bashrc`, `~/.zshrc`, or equivalent):

```bash
export ELEVENLABS_API_KEY="your-api-key-here"
```

On Windows, set it via System Environment Variables or PowerShell:

```powershell
[Environment]::SetEnvironmentVariable("ELEVENLABS_API_KEY", "your-api-key-here", "User")
```

### Optional configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `ELEVENLABS_VOICE_ID` | `21m00Tcm4TlvDq8ikWAM` (Rachel) | Voice to use. Browse voices at [elevenlabs.io/voices](https://elevenlabs.io/voices). |
| `ELEVENLABS_MODEL` | `eleven_flash_v2_5` | TTS model. Flash is fastest; use `eleven_multilingual_v2` for non-English. |

## Usage

```bash
# Speak a message (returns immediately, plays async)
speak "Build complete"

# Check version
speak --version

# Stop the background worker
speak --stop
```

The first call starts a background worker daemon that processes a message queue. The daemon automatically exits after 60 seconds of inactivity.

## How it works

```
speak "text" ──▶ enqueue ──▶ worker daemon ──▶ ElevenLabs API ──▶ audio playback
  (returns           │            │                   │                │
   immediately)      ▼            ▼                   ▼                ▼
               file queue    reads FIFO         returns MP3      oto v3 player
              with locking   idle timeout       go-mp3 decode    (ALSA/CoreAudio/WASAPI)
```

## Building from source

Requires Go 1.24+ and platform audio libraries:

- **Linux**: `sudo apt-get install libasound2-dev`
- **macOS**: Xcode Command Line Tools (CoreAudio)
- **Windows**: No extra dependencies (WASAPI built-in)

```bash
git clone https://github.com/dkmaker/claude-speak.git
cd claude-speak
make build
```

## License

MIT

# Audio Toggle

A simple Go program to toggle between your Logitech headphones and desktop speakers on Arch Linux with PipeWire/PulseAudio.

## Installation

### 1. Build the program

```bash
go build -o audio-toggle main.go
```

### 2. Install to local bin directory

```bash
# Create personal bin directory if it doesn't exist
mkdir -p ~/.local/bin

# Copy the executable
cp audio-toggle ~/.local/bin/

# Make sure it's executable
chmod +x ~/.local/bin/audio-toggle
```

## Usage

Simply run from anywhere:

```bash
audio-toggle
```

The program will:

- Toggle between your Logitech PRO X headphones and desktop speakers
- Show which device it switched to
- Gracefully handle if a device is unplugged

## Requirements

- Arch Linux with PipeWire or PulseAudio
- `pactl` command available (usually included with audio setup)
- Go compiler (only needed for building)

## How it works

The program uses `pactl` commands to:

1. Get the current default audio sink
2. Find available Logitech headphones and desktop speakers
3. Switch to the other device
4. Provide user-friendly feedback

## Troubleshooting

If devices aren't detected, check what's available:

```bash
pactl list short sinks
```

The program matches devices by patterns, so it should work even if device IDs change when unplugging/replugging USB devices.

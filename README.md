# MCaster

A simple command-line tool for testing multicast connectivity by sending and receiving UDP multicast packets with timestamps and sequence numbers.

## Features

- 🚀 **Send multicast packets** with configurable intervals
- 📥 **Receive multicast packets** and display timing information
- 🌐 **Interface binding** for multi-homed systems
- ⚙️ **Flexible configuration** via CLI flags, environment variables, or config files
- 📊 **Network delay measurement** for received packets


## Why not just use \<insert tool that already does this\>?

There are lots of other tools that provide similar functionality. This was largely an experiment in vibe-coding. The entirety of the application was written via assistance from Claude Code.

For a while now I had been meaning to see how far I could get with simple(ish) tooling via vibe-coding. I needed to pick a project to do it with, and I just happened to be in the middle of a multicast refresher so figured why not.

All things considered it worked out pretty well. Getting an MVP to work with Linux and MacOS was pretty straightforward, maybe an hour of work. There was another hour maybe two of me going, "that was too easy", then playing the role of PM armed with feature creep pistols aimed at Claude.  

We also wasted a few additional hours trying to get a cross-platform version that included Windows before giving up.

After some peer review from unaware friends who are more Go-savvy than I, the general consensus seems to be that it's reasonably clean, readable code.

## Installation

### From Source

```bash
git clone https://github.com/hyposcaler-bot/mcaster.git
cd mcaster
make build
```

### Using Go Install

```bash
go install github.com/hyposcaler-bot/mcaster/cmd/mcaster@latest
```

## Quick Start

### Send multicast packets
```bash
# Send to default group (239.23.23.23:2323)
./bin/mcaster send

# Send to specific group with custom interval
./bin/mcaster send -g 224.0.1.1:8080 -t 500ms
```

### Receive multicast packets
```bash
# Receive from default group
./bin/mcaster receive

# Receive from specific group via specific interface
./bin/mcaster receive -g 224.0.1.1:8080 -i eth0
```

## Usage

### Commands

- `send` - Send multicast packets continuously
- `receive` - Listen for and display received packets

### Global Flags

- `-g, --group` - Multicast group address:port (default: "239.23.23.23:2323")
- `-i, --interface` - Network interface name (optional)
- `-d, --dport` - Destination port (overrides port in group address; default: 0 = use group port)
- `--config` - Config file path (default: $HOME/.mcaster.yaml)

### Send-specific Flags

- `-t, --interval` - Send interval (default: 1s)
- `--ttl` - TTL (Time To Live) for multicast packets (default: 1, range: 1-255)
- `-s, --sport` - Source port for sending packets (default: 0 = random, range: 0-65535)

### Examples

```bash
# Basic usage
mcaster send
mcaster receive

# Custom multicast group
mcaster send --group 224.0.1.1:8080
mcaster receive --group 224.0.1.1:8080

# Bind to specific interface
mcaster send --interface eth0
mcaster receive --interface eth0

# Fast sending interval
mcaster send --interval 100ms

# Custom TTL for cross-router testing
mcaster send --ttl 32

# Send from specific source port
mcaster send --sport 12345

# Send to specific destination port
mcaster send --dport 8080
mcaster receive --dport 8080

# Using environment variables
MULTICAST_GROUP=224.0.1.1:8080 mcaster send
MULTICAST_INTERFACE=eth0 mcaster receive
MULTICAST_TTL=16 MULTICAST_SPORT=12345 MULTICAST_DPORT=8080 mcaster send
```

## Configuration

### Environment Variables

- `MULTICAST_GROUP` - Multicast group address:port
- `MULTICAST_INTERFACE` - Network interface name
- `MULTICAST_INTERVAL` - Send interval (sender only)
- `MULTICAST_TTL` - TTL for multicast packets (sender only)
- `MULTICAST_SPORT` - Source port for sending packets (sender only)
- `MULTICAST_DPORT` - Destination port (overrides group port)

### Configuration File

Create `~/.mcaster.yaml`:

```yaml
group: "239.23.23.23:2323"
interface: "eth0"
interval: "500ms"
ttl: 16
sport: 12345
dport: 8080
```

## Output Format

### Sender Output
```
🚀 Starting multicast sender to 239.23.23.23:2323
📡 Sending packets every 1s (TTL: 1, source port: 54321)
⏹️  Press Ctrl+C to stop

📤 [15:04:05.123] Sent packet #1
📤 [15:04:06.124] Sent packet #2
```

### Receiver Output
```
🎯 Starting multicast receiver on 239.23.23.23:2323
👂 Waiting for packets...

📥 [15:04:05.125] Received packet #1 from hostname (192.168.1.100:54321) - delay: 2ms
📥 [15:04:06.126] Received packet #2 from hostname (192.168.1.100:54321) - delay: 2ms
```

## Common Use Cases

### Testing Network Connectivity

1. **Basic connectivity test**:
   ```bash
   # Terminal 1
   mcaster receive
   
   # Terminal 2
   mcaster send
   ```

2. **Cross-subnet testing**:
   ```bash
   # On receiver host
   mcaster receive -g 224.0.1.1:8080 -i eth0
   
   # On sender host
   mcaster send -g 224.0.1.1:8080 -i eth0
   ```

3. **Performance testing**:
   ```bash
   # High-frequency sending
   mcaster send -t 10ms
   ```

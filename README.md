# MCaster

A simple command-line tool for testing multicast connectivity by sending and receiving UDP multicast packets with timestamps and sequence numbers.

## Features

- 🚀 **Send multicast packets** with configurable intervals
- 📥 **Receive multicast packets** and display timing information
- 🌐 **Interface binding** for multi-homed systems
- ⚙️ **Flexible configuration** via CLI flags, environment variables, or config files
- 📊 **Network delay measurement** for received packets
- 🎯 **Professional CLI** built with Cobra framework

## Installation

### From Source

```bash
git clone https://github.com/yourusername/mcaster.git
cd mcaster
make build
```

### Using Go Install

```bash
go install github.com/yourusername/mcaster/cmd/mcaster@latest
```

## Quick Start

### Send multicast packets
```bash
# Send to default group (224.1.1.1:9999)
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

- `-g, --group` - Multicast group address:port (default: "224.1.1.1:9999")
- `-i, --interface` - Network interface name (optional)
- `--config` - Config file path (default: $HOME/.mcaster.yaml)

### Send-specific Flags

- `-t, --interval` - Send interval (default: 1s)

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

# Using environment variables
MULTICAST_GROUP=224.0.1.1:8080 mcaster send
MULTICAST_INTERFACE=eth0 mcaster receive
```

## Configuration

### Environment Variables

- `MULTICAST_GROUP` - Multicast group address:port
- `MULTICAST_INTERFACE` - Network interface name
- `MULTICAST_INTERVAL` - Send interval (sender only)

### Configuration File

Create `~/.mcaster.yaml`:

```yaml
group: "224.0.1.1:8080"
interface: "eth0"
interval: "500ms"
```

## Output Format

### Sender Output
```
🚀 Starting multicast sender to 224.1.1.1:9999
📡 Sending packets every 1s
⏹️  Press Ctrl+C to stop

📤 [15:04:05.123] Sent packet #1
📤 [15:04:06.124] Sent packet #2
```

### Receiver Output
```
🎯 Starting multicast receiver on 224.1.1.1:9999
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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

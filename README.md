# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and alert on any new or closed ports:

```bash
portwatch start --interval 60 --notify
```

Take a snapshot of currently open ports to use as a baseline:

```bash
portwatch snapshot
```

View detected changes since the last snapshot:

```bash
portwatch diff
```

### Example Output

```
[portwatch] Baseline loaded: 12 open ports
[portwatch] ALERT: New port detected → 0.0.0.0:8080 (PID 4821)
[portwatch] ALERT: Port closed → 0.0.0.0:3000
```

## Configuration

portwatch looks for a config file at `~/.portwatch.yaml`:

```yaml
interval: 30
notify: true
ignore:
  - 22
  - 80
  - 443
```

## License

MIT © 2024 yourusername
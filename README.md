# pulsectl

> Simple health-check scheduler that polls HTTP endpoints and reports status over time.

---

## Installation

```bash
go install github.com/yourusername/pulsectl@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/pulsectl.git
cd pulsectl && go build -o pulsectl .
```

---

## Usage

Define your endpoints in a `config.yaml` file:

```yaml
interval: 30s
endpoints:
  - name: API Server
    url: https://api.example.com/health
  - name: Auth Service
    url: https://auth.example.com/ping
```

Run the scheduler:

```bash
pulsectl --config config.yaml
```

Sample output:

```
[OK]   API Server     https://api.example.com/health     200  (42ms)
[FAIL] Auth Service   https://auth.example.com/ping      503  (120ms)
```

Use `--json` to emit structured JSON logs, or `--interval` to override the polling interval at runtime:

```bash
pulsectl --config config.yaml --interval 10s --json
```

---

## Flags

| Flag         | Default       | Description                        |
|--------------|---------------|------------------------------------|
| `--config`   | `config.yaml` | Path to configuration file         |
| `--interval` | from config   | Override polling interval          |
| `--json`     | `false`       | Output results as JSON             |
| `--timeout`  | `5s`          | HTTP request timeout per endpoint  |

---

## License

MIT © 2024 yourusername
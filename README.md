# PSBL-Watch 🏀

PSBL-Watch is a small Go service that monitors Puget Sound Basketball League's *Hoops-On-Demand* page and sends real-time Pushover notifications when a target pick-up game becomes available or its roster changes.

The program is designed to run unattended (every 3 minutes) in a minimal container, making it perfect for hosting on a small VPS, Raspberry Pi, or Kubernetes cluster.

---

## Features

• Scrapes the public Hoops-On-Demand web page on an adjustable schedule (default: **every 3 minutes**).  
• Detects when the roster for a specific game opens up, fills up, or otherwise changes.  
• Sends rich-formatted push notifications via the [Pushover](https://pushover.net/) API.  
• Caches the previous state in-memory to avoid duplicate alerts.  
• Ships as a single statically-linked binary (≈ 3 MB) or a 10 MB distroless Docker image.

---

## Quick start (Docker)

```bash
# Build the image
$ docker build -t psbl-watch .

# Run the container (replace the tokens with your Pushover credentials)
$ docker run -d --name psbl-watch \
    -e PUSHOVER_APP_TOKEN=xxx \
    -e PUSHOVER_USER_TOKEN=yyy \
    psbl-watch
```

By default the container logs to stdout/stderr. Use `docker logs -f psbl-watch` to tail the output.

---

## Local development

Prerequisites: **Go 1.22+**

```bash
git clone https://github.com/<your-fork>/psbl-watch.git
cd psbl-watch

# Configure environment variables
export PUSHOVER_APP_TOKEN=xxx
export PUSHOVER_USER_TOKEN=yyy

# Run once immediately
go run ./cmd/psbl-watch
```

Unit tests can be executed with:

```bash
go test ./...
```

---

## Configuration

The following environment variables are honoured:

| Variable | Description |
|----------|-------------|
| `PUSHOVER_APP_TOKEN`  | Your application API token obtained from Pushover. **Required** |
| `PUSHOVER_USER_TOKEN` | The user / group token that should receive notifications. **Required** |

### Target game

The game that is being watched is currently hard-coded in `cmd/psbl-watch/main.go`:

```go
store.TargetGame{
    Gym:   "Seattle Central College #1",
    Type:  "Saturday Morning Hoops",
    Level: "Recreational-CoEd",
}
```

Feel free to change these values or make them configurable through environment variables/flags if you need to watch a different court or league level.

---

## Deployment tips

• **Cron schedule** – Adjust the cron expression in `cmd/psbl-watch/main.go` (`@every 3m`) if you need a different polling interval.  
• **Timezone** – The Docker image sets `TZ=America/Los_Angeles`; tweak this build-arg if you operate in another zone.

---

## License

Distributed under the MIT License. See `LICENSE` for more information. 
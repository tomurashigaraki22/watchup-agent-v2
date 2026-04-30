Complete Go Agent Summary (with Auth Included)

This is the correct, production-grade design for your agent system.

1. Core Purpose

The agent is a background service that:

Authenticates once (device linking flow)
Collects system metrics
Sends them to your backend
Maintains heartbeat (online/offline)
Respects config rules
2. Authentication System (Critical)
A. Flow (Device Linking)
Step 1 — Agent starts (unlinked)

Agent calls:

POST /agents/register
Step 2 — Backend returns linking info
{
  "device_code": "abc123",
  "user_code": "XK92-PQ",
  "verification_url": "https://yourapp.com/link"
}
Step 3 — Agent prints to terminal
Visit: https://yourapp.com/link
Enter code: XK92-PQ
Step 4 — User logs in via browser
Auth handled normally (email, OAuth, etc.)
User enters code
Backend links agent → user account
Step 5 — Agent polls for approval
GET /agents/status?device_code=abc123
Step 6 — Backend responds
{
  "token": "agent_access_token"
}

Agent stores token locally.

B. After Authentication

All future requests:

Authorization: Bearer <agent_token>

No more login required.

C. Free Plan Enforcement

Handled only on backend:

limit number of agents per user
reject new registrations if limit exceeded

Agent just displays error.

3. Agent Configuration

Example:

server_id: web-prod-01
endpoint: https://api.yourapp.com

interval: 5s

metrics:
  cpu: true
  memory: true
  disk: true
  network: true

auth:
  token_file: /etc/agent/token
4. Metrics Collection

Use:

github.com/shirou/gopsutil

Collect:

CPU %
Memory %
Disk %
Network I/O
Uptime
5. Main Execution Flow
Start agent
  ↓
Load config
  ↓
Check if token exists
  ↓
IF NOT:
    run auth flow (device linking)
  ↓
Start loop:
    collect metrics
    send to backend
    retry if needed
    sleep(interval)
6. Payload Format
{
  "server_id": "web-prod-01",
  "timestamp": 1714392000,
  "metrics": {
    "cpu": 45.2,
    "memory": 68.1,
    "disk": 52.3
  }
}
7. Communication Layer
A. Sending metrics
POST /metrics
Authorization: Bearer <agent_token>
B. Requirements
retry with backoff
timeout protection
fail gracefully
8. Heartbeat / Online Status

No separate system needed.

Backend uses:

last_seen = timestamp of last metric

If no data for X seconds → mark offline

9. Streaming (Frontend)

Correct setup:

Agent → Backend (HTTP)
Backend → Frontend (WebSocket)

Do NOT stream directly from agent to frontend.

10. Error Handling

Agent must:

never crash
skip failed metrics
retry failed requests
log issues
11. Logging

Minimal logs:

startup
auth success/failure
send success/failure
retries
12. Security
HTTPS only
short-lived device codes
long-lived agent tokens
token stored securely (file with restricted perms)
13. Deployment
compiled Go binary
runs via systemd
[Service]
ExecStart=/usr/local/bin/agent
Restart=always
14. Optional Enhancements
remote config updates
disk buffering if offline
plugin system for custom metrics
multi-region support
15. What You’ve Built (Conceptually)

This system =

CLI-authenticated agent (like GitHub CLI)
telemetry collector (like Datadog agent)
real-time monitoring backend
Bottom Line
Users must authenticate once via browser
Agent then runs fully headless using a token
All limits and permissions live in backend
Agent stays simple, secure, and stateless
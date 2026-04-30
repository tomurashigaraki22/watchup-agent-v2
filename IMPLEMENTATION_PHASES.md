# WatchUp Agent - Implementation Phases

## Project Overview
A production-grade Go agent that collects system metrics and sends them to a backend API. The agent uses device linking authentication (similar to GitHub CLI) and runs as a headless background service.

---

## 📋 Phase 0: Project Setup & Foundation

### Goals
- Set up Go project structure
- Install dependencies
- Create basic configuration system

### Tasks
- [ ] Initialize Go module (`go mod init watchup-agent`)
- [ ] Install core dependencies:
  ```bash
  go get github.com/shirou/gopsutil/v3
  go get gopkg.in/yaml.v3
  ```
- [ ] Create project structure:
  ```
  watchup-agent/
  ├── cmd/
  │   └── agent/
  │       └── main.go
  ├── internal/
  │   ├── auth/
  │   ├── config/
  │   ├── metrics/
  │   └── client/
  ├── config.example.yaml
  └── README.md
  ```
- [ ] Create basic `config.yaml` structure:
  ```yaml
  server_id: ""
  endpoint: "https://api.yourapp.com"
  interval: 5s
  metrics:
    cpu: true
    memory: true
    disk: true
    network: true
  auth:
    token_file: "./agent_token"
  ```

### Deliverables
- Working Go project with proper structure
- Configuration loading system
- Basic logging setup

### Estimated Time: 2-4 hours

---

## 🔐 Phase 1: Authentication System (Device Linking)

### Goals
- Implement device linking flow
- Handle token storage and retrieval
- Manage authentication state

### Tasks

#### 1.1 Device Registration
- [ ] Create `internal/auth/device.go`
- [ ] Implement `POST /agents/register` call
- [ ] Parse response:
  ```go
  type DeviceCodeResponse struct {
      DeviceCode      string `json:"device_code"`
      UserCode        string `json:"user_code"`
      VerificationURL string `json:"verification_url"`
      ExpiresIn       int    `json:"expires_in"`
  }
  ```
- [ ] Display linking instructions to user:
  ```
  🔗 Link this agent to your account:
  
  Visit: https://yourapp.com/link
  Enter code: XK92-PQ
  
  Waiting for authorization...
  ```

#### 1.2 Token Polling
- [ ] Implement polling mechanism (`GET /agents/status?device_code=abc123`)
- [ ] Handle polling states:
  - Pending (keep polling)
  - Approved (save token)
  - Expired (restart flow)
  - Denied (exit with error)
- [ ] Add exponential backoff (5s, 10s, 15s intervals)

#### 1.3 Token Management
- [ ] Create `internal/auth/token.go`
- [ ] Implement secure token storage:
  ```go
  func SaveToken(token string, filepath string) error
  func LoadToken(filepath string) (string, error)
  func TokenExists(filepath string) bool
  ```
- [ ] Set file permissions (0600 - owner read/write only)
- [ ] Add token validation

#### 1.4 HTTP Client with Auth
- [ ] Create `internal/client/client.go`
- [ ] Implement authenticated HTTP client:
  ```go
  type Client struct {
      BaseURL    string
      Token      string
      HTTPClient *http.Client
  }
  ```
- [ ] Add `Authorization: Bearer <token>` header to all requests
- [ ] Handle 401 responses (token expired/invalid)

### Deliverables
- Complete device linking flow
- Secure token storage
- Authenticated HTTP client

### Estimated Time: 6-8 hours

---

## 📊 Phase 2: Core Metrics Collection (Must-Have)

### Goals
- Collect essential system metrics
- Create clean metric data structures
- Test metric collection independently

### Tasks

#### 2.1 CPU Metrics
- [ ] Create `internal/metrics/cpu.go`
- [ ] Implement CPU percentage collection:
  ```go
  func GetCPUUsage() (float64, error) {
      percentages, err := cpu.Percent(time.Second, false)
      return percentages[0], err
  }
  ```
- [ ] Add per-core CPU stats (optional)

#### 2.2 Memory Metrics
- [ ] Create `internal/metrics/memory.go`
- [ ] Collect memory stats:
  ```go
  type MemoryMetrics struct {
      Total       uint64  `json:"total"`
      Used        uint64  `json:"used"`
      Available   uint64  `json:"available"`
      UsedPercent float64 `json:"used_percent"`
  }
  ```
- [ ] Use `mem.VirtualMemory()`

#### 2.3 Disk Metrics
- [ ] Create `internal/metrics/disk.go`
- [ ] Collect disk usage for all partitions:
  ```go
  type DiskMetrics struct {
      Partition   string  `json:"partition"`
      Total       uint64  `json:"total"`
      Used        uint64  `json:"used"`
      Free        uint64  `json:"free"`
      UsedPercent float64 `json:"used_percent"`
  }
  ```
- [ ] Use `disk.Usage()` and `disk.Partitions()`

#### 2.4 Network Metrics (Bandwidth)
- [ ] Create `internal/metrics/network.go`
- [ ] Collect network I/O:
  ```go
  type NetworkMetrics struct {
      Interface   string `json:"interface"`
      BytesSent   uint64 `json:"bytes_sent"`
      BytesRecv   uint64 `json:"bytes_recv"`
      PacketsSent uint64 `json:"packets_sent"`
      PacketsRecv uint64 `json:"packets_recv"`
  }
  ```
- [ ] Use `net.IOCounters(true)` for per-interface stats
- [ ] Calculate bandwidth (bytes/sec) by comparing snapshots

#### 2.5 System Uptime
- [ ] Add uptime collection:
  ```go
  func GetUptime() (uint64, error) {
      return host.Uptime()
  }
  ```

#### 2.6 Metrics Aggregator
- [ ] Create `internal/metrics/collector.go`
- [ ] Implement main collector:
  ```go
  type MetricsPayload struct {
      ServerID  string                 `json:"server_id"`
      Timestamp int64                  `json:"timestamp"`
      Metrics   map[string]interface{} `json:"metrics"`
  }
  
  func CollectAll(config *Config) (*MetricsPayload, error)
  ```
- [ ] Handle partial failures (skip failed metrics, don't crash)

### Deliverables
- All core metrics collectors working
- Unified metrics payload structure
- Error handling for failed collections

### Estimated Time: 8-10 hours

---

## 🔄 Phase 3: Main Agent Loop & Communication

### Goals
- Implement main execution loop
- Send metrics to backend
- Handle retries and failures

### Tasks

#### 3.1 Main Loop
- [ ] Create main agent loop in `cmd/agent/main.go`:
  ```go
  func main() {
      // Load config
      // Check/perform auth
      // Start metrics loop
      for {
          metrics := CollectAll()
          SendMetrics(metrics)
          time.Sleep(interval)
      }
  }
  ```
- [ ] Add graceful shutdown (handle SIGTERM, SIGINT)
- [ ] Implement context cancellation

#### 3.2 HTTP Communication
- [ ] Implement `POST /metrics` in `internal/client/metrics.go`:
  ```go
  func (c *Client) SendMetrics(payload *MetricsPayload) error
  ```
- [ ] Add request timeout (30s default)
- [ ] Set proper headers:
  ```
  Content-Type: application/json
  Authorization: Bearer <token>
  ```

#### 3.3 Retry Logic
- [ ] Implement exponential backoff:
  ```go
  func SendWithRetry(payload *MetricsPayload, maxRetries int) error {
      // Retry: 1s, 2s, 4s, 8s, 16s
  }
  ```
- [ ] Max retries: 5
- [ ] Log each retry attempt
- [ ] Skip and continue on final failure (don't crash)

#### 3.4 Error Handling
- [ ] Never crash the agent
- [ ] Log errors to stdout/file
- [ ] Continue operation on non-fatal errors
- [ ] Handle network failures gracefully

### Deliverables
- Working agent loop
- Reliable metric transmission
- Robust error handling

### Estimated Time: 6-8 hours

---

## 🌐 Phase 4: Network Monitoring (Extended)

### Goals
- Add active connections tracking
- Implement port monitoring
- Add network latency checks

### Tasks

#### 4.1 Active Connections
- [ ] Create `internal/metrics/connections.go`
- [ ] Count active connections:
  ```go
  type ConnectionMetrics struct {
      Total     int `json:"total"`
      TCP       int `json:"tcp"`
      UDP       int `json:"udp"`
      Listening int `json:"listening"`
  }
  ```
- [ ] Use `net.Connections("all")`
- [ ] Group by state (ESTABLISHED, LISTEN, etc.)

#### 4.2 Port Monitoring
- [ ] Add port check configuration:
  ```yaml
  ports:
    - port: 80
      name: "HTTP"
    - port: 443
      name: "HTTPS"
    - port: 5432
      name: "PostgreSQL"
  ```
- [ ] Implement TCP connection check:
  ```go
  func CheckPort(host string, port int, timeout time.Duration) bool {
      conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
      if err != nil {
          return false
      }
      conn.Close()
      return true
  }
  ```
- [ ] Return port status (up/down) for each configured port

#### 4.3 Network Latency
- [ ] Add latency check configuration:
  ```yaml
  latency_checks:
    - host: "8.8.8.8"
      name: "Google DNS"
    - host: "api.yourapp.com"
      name: "API Server"
  ```
- [ ] Implement ICMP ping or HTTP check:
  ```go
  func CheckLatency(host string) (time.Duration, error) {
      start := time.Now()
      conn, err := net.DialTimeout("tcp", host+":80", 5*time.Second)
      if err != nil {
          return 0, err
      }
      conn.Close()
      return time.Since(start), nil
  }
  ```

### Deliverables
- Connection tracking
- Port monitoring system
- Latency measurements

### Estimated Time: 6-8 hours

---

## 🔍 Phase 5: Process Monitoring

### Goals
- List running processes
- Track top processes by CPU/memory
- Count processes and detect zombies

### Tasks

#### 5.1 Process List
- [ ] Create `internal/metrics/processes.go`
- [ ] Get all processes:
  ```go
  func GetProcessList() ([]*ProcessInfo, error) {
      procs, _ := process.Processes()
      // Return simplified list
  }
  ```
- [ ] Limit to top N processes (don't send thousands)

#### 5.2 Top Processes
- [ ] Sort processes by CPU usage:
  ```go
  type ProcessInfo struct {
      PID         int32   `json:"pid"`
      Name        string  `json:"name"`
      CPUPercent  float64 `json:"cpu_percent"`
      MemoryMB    uint64  `json:"memory_mb"`
      Status      string  `json:"status"`
  }
  ```
- [ ] Return top 10 by CPU
- [ ] Return top 10 by memory

#### 5.3 Process Counts
- [ ] Count total processes
- [ ] Count by status (running, sleeping, etc.)
- [ ] Detect zombie processes:
  ```go
  func CountZombies() (int, error) {
      procs, _ := process.Processes()
      zombies := 0
      for _, p := range procs {
          status, _ := p.Status()
          if status == "Z" {
              zombies++
          }
      }
      return zombies, nil
  }
  ```

### Deliverables
- Process monitoring system
- Top processes tracking
- Zombie detection

### Estimated Time: 4-6 hours

---

## 🛡️ Phase 6: Service Monitoring

### Goals
- Check if critical services are running
- Monitor service health via ports or systemd

### Tasks

#### 6.1 Service Configuration
- [ ] Add service monitoring config:
  ```yaml
  services:
    - name: "MySQL"
      type: "port"
      port: 3306
    - name: "Nginx"
      type: "systemd"
      unit: "nginx.service"
    - name: "Redis"
      type: "port"
      port: 6379
  ```

#### 6.2 Port-Based Checks
- [ ] Reuse port checking from Phase 4
- [ ] Return service status (up/down)

#### 6.3 Systemd Checks (Linux only)
- [ ] Execute `systemctl is-active <service>`:
  ```go
  func CheckSystemdService(unit string) (bool, error) {
      cmd := exec.Command("systemctl", "is-active", unit)
      err := cmd.Run()
      return err == nil, nil
  }
  ```
- [ ] Handle non-Linux systems gracefully

### Deliverables
- Service health monitoring
- Multi-method checks (port + systemd)

### Estimated Time: 4-6 hours

---

## 🚀 Phase 7: Deployment & Production Readiness

### Goals
- Build production binary
- Create systemd service
- Add installation script

### Tasks

#### 7.1 Build System
- [ ] Create build script:
  ```bash
  #!/bin/bash
  go build -o watchup-agent cmd/agent/main.go
  ```
- [ ] Add version information:
  ```go
  var Version = "1.0.0"
  ```
- [ ] Create release builds for multiple platforms

#### 7.2 Systemd Service
- [ ] Create `watchup-agent.service`:
  ```ini
  [Unit]
  Description=WatchUp Monitoring Agent
  After=network.target

  [Service]
  Type=simple
  User=watchup
  ExecStart=/usr/local/bin/watchup-agent
  Restart=always
  RestartSec=10

  [Install]
  WantedBy=multi-user.target
  ```

#### 7.3 Installation Script
- [ ] Create `install.sh`:
  ```bash
  #!/bin/bash
  # Copy binary
  sudo cp watchup-agent /usr/local/bin/
  sudo chmod +x /usr/local/bin/watchup-agent
  
  # Create config directory
  sudo mkdir -p /etc/watchup-agent
  sudo cp config.yaml /etc/watchup-agent/
  
  # Install systemd service
  sudo cp watchup-agent.service /etc/systemd/system/
  sudo systemctl daemon-reload
  sudo systemctl enable watchup-agent
  
  # Start service
  sudo systemctl start watchup-agent
  ```

#### 7.4 Documentation
- [ ] Write comprehensive README.md
- [ ] Document configuration options
- [ ] Add troubleshooting guide
- [ ] Create quick start guide

### Deliverables
- Production-ready binary
- Systemd integration
- Installation automation
- Complete documentation

### Estimated Time: 4-6 hours

---

## 🎯 Phase 8: Testing & Refinement

### Goals
- Test all components
- Fix bugs
- Optimize performance

### Tasks

#### 8.1 Unit Tests
- [ ] Test metric collectors
- [ ] Test authentication flow
- [ ] Test HTTP client
- [ ] Test configuration loading

#### 8.2 Integration Tests
- [ ] Test full agent flow
- [ ] Test retry logic
- [ ] Test error handling
- [ ] Test graceful shutdown

#### 8.3 Performance Testing
- [ ] Measure CPU usage of agent
- [ ] Measure memory footprint
- [ ] Test with different intervals (1s, 5s, 30s)
- [ ] Ensure no memory leaks

#### 8.4 Edge Cases
- [ ] Test with no network
- [ ] Test with invalid config
- [ ] Test with expired token
- [ ] Test with backend downtime

### Deliverables
- Test suite
- Performance benchmarks
- Bug fixes

### Estimated Time: 6-8 hours

---

## 📈 Phase 9: Optional Enhancements (Future)

### Goals
- Add advanced features based on needs

### Possible Features
- [ ] Remote configuration updates
- [ ] Disk buffering when offline
- [ ] Plugin system for custom metrics
- [ ] Multi-region support
- [ ] Compression for large payloads
- [ ] Metrics aggregation (send averages instead of raw data)
- [ ] Custom metric thresholds and alerts
- [ ] Agent auto-update mechanism

### Estimated Time: Variable (2-4 hours per feature)

---

## 🎨 Phase 10: Advanced Monitoring (Heavy - Optional)

### ⚠️ Warning
These features significantly increase complexity and resource usage. Only implement if absolutely necessary.

### Tasks

#### 10.1 SSH Login Monitoring
- [ ] Parse `/var/log/auth.log`
- [ ] Track failed login attempts
- [ ] Track successful logins
- [ ] Send security alerts

#### 10.2 Sudo Command Tracking
- [ ] Parse sudo commands from auth log
- [ ] Track who ran what command
- [ ] Send audit trail

#### 10.3 Log Tailing
- [ ] Implement log file tailing
- [ ] Add log filtering
- [ ] Limit log volume sent to backend
- [ ] Handle log rotation

### Deliverables
- Security monitoring
- Log ingestion system

### Estimated Time: 12-16 hours

---

## 📊 Summary Timeline

| Phase | Description | Time | Priority |
|-------|-------------|------|----------|
| 0 | Project Setup | 2-4h | Critical |
| 1 | Authentication | 6-8h | Critical |
| 2 | Core Metrics | 8-10h | Critical |
| 3 | Main Loop | 6-8h | Critical |
| 4 | Network Monitoring | 6-8h | High |
| 5 | Process Monitoring | 4-6h | High |
| 6 | Service Monitoring | 4-6h | Medium |
| 7 | Deployment | 4-6h | High |
| 8 | Testing | 6-8h | High |
| 9 | Enhancements | Variable | Low |
| 10 | Advanced Monitoring | 12-16h | Low |

**Total Core Development Time: 40-56 hours**

---

## 🎯 Recommended Implementation Order

### Week 1: Foundation
1. Phase 0: Project Setup
2. Phase 1: Authentication
3. Phase 2: Core Metrics (CPU, Memory, Disk)

### Week 2: Core Features
4. Phase 2: Network Metrics
5. Phase 3: Main Loop & Communication
6. Phase 4: Extended Network Monitoring

### Week 3: Extended Features
7. Phase 5: Process Monitoring
8. Phase 6: Service Monitoring
9. Phase 7: Deployment Setup

### Week 4: Polish
10. Phase 8: Testing & Refinement
11. Documentation
12. Production deployment

---

## 🔑 Key Success Metrics

- [ ] Agent runs continuously without crashes
- [ ] Metrics sent successfully every interval
- [ ] Authentication works on first run
- [ ] Agent restarts automatically on failure
- [ ] CPU usage < 5%
- [ ] Memory usage < 50MB
- [ ] No memory leaks over 24h+ runtime

---

## 📝 Notes

### What NOT to Build (Yet)
- Direct WebSocket streaming to frontend
- Real-time log streaming
- Complex log parsing
- Temperature monitoring
- Custom alerting (backend handles this)

### Architecture Principles
1. **Agent stays simple** - collect and send only
2. **Backend does heavy lifting** - analysis, storage, alerts
3. **Fail gracefully** - never crash, always retry
4. **Security first** - HTTPS only, secure token storage
5. **Observable** - good logging for debugging

### Configuration Best Practices
- Use environment variables for sensitive data
- Support both file and env var config
- Validate configuration on startup
- Provide sensible defaults

---

## 🚦 Getting Started

1. Start with Phase 0 (Project Setup)
2. Complete Phases 1-3 before moving to extended features
3. Test each phase thoroughly before moving on
4. Deploy early and iterate
5. Add advanced features only when core is stable

Good luck! 🚀

# WatchUp Agent - Complete Monitoring Capabilities

## Overview

The WatchUp Agent is a comprehensive system monitoring solution that collects real-time metrics from servers and sends them to the WatchUp backend platform. This document outlines all monitoring capabilities, data structures, and backend integration requirements.

---

## 🔧 System Requirements

- **Operating Systems**: Windows, Linux, macOS
- **Go Version**: 1.19 or later
- **Network**: HTTPS access to WatchUp backend
- **Permissions**: Standard user permissions (no root required for basic metrics)

---

## 📊 Complete Monitoring Capabilities

### 🖥️ **CPU Monitoring**
**Status**: ✅ Implemented (Phase 2)

| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `usage_percent` | Overall CPU utilization | `float64` | Percentage (0-100) |
| `per_core` | Per-core CPU usage | `[]float64` | Percentage per core |
| `cores` | Number of logical cores | `int` | Count |

**Additional CPU Info** (collected once):
- CPU model name
- Vendor ID
- Clock speed (MHz)
- Physical cores

**Collection Method**: 1-second sampling window using gopsutil
**Update Frequency**: Every collection interval (configurable)

---

### 💾 **Memory Monitoring**
**Status**: ✅ Implemented (Phase 2)

| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `total` | Total system RAM | `uint64` | Bytes |
| `used` | Used memory | `uint64` | Bytes |
| `available` | Available memory | `uint64` | Bytes |
| `used_percent` | Memory utilization | `float64` | Percentage (0-100) |
| `free` | Free memory | `uint64` | Bytes |
| `cached` | Cached memory (Linux) | `uint64` | Bytes |
| `buffers` | Buffer memory (Linux) | `uint64` | Bytes |

**Swap Memory**:
- Total swap space
- Used swap space
- Swap utilization percentage

**Collection Method**: Real-time via gopsutil
**Update Frequency**: Every collection interval

---

### 💿 **Disk Monitoring**
**Status**: ✅ Implemented (Phase 2)

**Disk Usage** (per partition):
| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `partition` | Mount point/drive letter | `string` | Path |
| `filesystem` | Filesystem type | `string` | Type (NTFS, ext4, etc.) |
| `total` | Total disk space | `uint64` | Bytes |
| `used` | Used disk space | `uint64` | Bytes |
| `free` | Free disk space | `uint64` | Bytes |
| `used_percent` | Disk utilization | `float64` | Percentage (0-100) |

**Disk I/O Statistics**:
- Total read bytes
- Total write bytes
- Read operation count
- Write operation count

**Collection Method**: All mounted partitions (excludes virtual filesystems)
**Update Frequency**: Every collection interval
**Filtering**: Automatically skips tmpfs, proc, sys, and other virtual filesystems

---

### 🌐 **Network Monitoring**
**Status**: ✅ Implemented (Phase 2)

**Per-Interface Metrics**:
| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `interface` | Network interface name | `string` | Name (eth0, wlan0, etc.) |
| `bytes_sent` | Bytes transmitted | `uint64` | Bytes |
| `bytes_recv` | Bytes received | `uint64` | Bytes |
| `packets_sent` | Packets transmitted | `uint64` | Count |
| `packets_recv` | Packets received | `uint64` | Count |
| `errin` | Input errors | `uint64` | Count |
| `errout` | Output errors | `uint64` | Count |
| `dropin` | Input drops | `uint64` | Count |
| `dropout` | Output drops | `uint64` | Count |

**Network Summary**:
- Total bytes sent (all interfaces)
- Total bytes received (all interfaces)
- Total packets sent/received
- Active interface count

**Collection Method**: Per-interface statistics via gopsutil
**Update Frequency**: Every collection interval
**Filtering**: Excludes loopback and virtual interfaces (Docker, VPN, etc.)

---

### 🔗 **Network Connections Monitoring**
**Status**: ✅ Implemented (Phase 4) - *Optional*

| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `total` | Total active connections | `int` | Count |
| `tcp` | TCP connections | `int` | Count |
| `udp` | UDP connections | `int` | Count |
| `listening` | Listening connections | `int` | Count |
| `established` | Established connections | `int` | Count |
| `by_state` | Connections by state | `map[string]int` | Count per state |
| `by_family` | IPv4/IPv6 breakdown | `map[string]int` | Count per family |

**Additional Connection Data**:
- List of listening ports
- Connection counts per process (when available)

**Collection Method**: System connection table via gopsutil
**Update Frequency**: Every collection interval
**Configuration**: Enable with `metrics.connections: true`

---

### 🚪 **Port Monitoring**
**Status**: ✅ Implemented (Phase 4) - *Configurable*

| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `port` | Port number | `int` | Port number |
| `name` | Service name | `string` | Custom name |
| `host` | Target host | `string` | Hostname/IP |
| `is_up` | Port accessibility | `bool` | true/false |
| `response_time_ms` | Connection time | `time.Duration` | Milliseconds |
| `error` | Error message (if down) | `string` | Error description |
| `checked_at` | Check timestamp | `time.Time` | RFC3339 |

**Port Summary Statistics**:
- Total ports monitored
- Ports up/down count
- Average response time

**Collection Method**: TCP connection attempts
**Update Frequency**: Every collection interval
**Configuration**: Define in `ports` array in config
**Concurrency**: All ports checked simultaneously
**Timeout**: Configurable per port (default 5s)

**Example Configuration**:
```yaml
ports:
  - port: 80
    name: "HTTP"
    host: "localhost"
    timeout: "5s"
  - port: 443
    name: "HTTPS"
    host: "mysite.com"
    timeout: "3s"
```

---

### ⚡ **Network Latency Monitoring**
**Status**: ✅ Implemented (Phase 4) - *Configurable*

| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `host` | Target host/URL | `string` | Hostname/URL |
| `name` | Check name | `string` | Custom name |
| `type` | Check type | `string` | "tcp" or "http" |
| `latency_ms` | Response latency | `time.Duration` | Milliseconds |
| `is_successful` | Check success | `bool` | true/false |
| `error` | Error message (if failed) | `string` | Error description |
| `checked_at` | Check timestamp | `time.Time` | RFC3339 |

**Latency Summary Statistics**:
- Total checks performed
- Successful/failed counts
- Average/min/max latency

**Check Types**:
1. **TCP Latency**: Measures TCP connection establishment time
2. **HTTP Latency**: Measures full HTTP request/response time

**Collection Method**: 
- TCP: `net.DialTimeout()`
- HTTP: `http.Client.Get()`
**Update Frequency**: Every collection interval
**Configuration**: Define in `latency_checks` array
**Concurrency**: All checks run simultaneously
**Timeout**: Configurable per check (default 5s)

**Example Configuration**:
```yaml
latency_checks:
  - host: "8.8.8.8"
    name: "Google DNS"
    type: "tcp"
    port: 53
    timeout: "3s"
  - host: "https://api.myservice.com"
    name: "My API"
    type: "http"
    timeout: "10s"
```

---

### 🖥️ **System Information**
**Status**: ✅ Implemented (Phase 2)

| Metric | Description | Data Type | Units |
|--------|-------------|-----------|-------|
| `uptime` | System uptime | `uint64` | Seconds |
| `hostname` | System hostname | `string` | Hostname |
| `os` | Operating system | `string` | OS name |
| `platform` | Platform details | `string` | Platform version |
| `platform_version` | OS version | `string` | Version string |
| `architecture` | System architecture | `string` | amd64, arm64, etc. |
| `boot_time` | Boot timestamp | `uint64` | Unix timestamp |

**Collection Method**: System calls via gopsutil
**Update Frequency**: Every collection interval
**Purpose**: System identification and uptime tracking

---

## 🚀 **Backend Integration**

### Authentication Flow

The agent uses OAuth2-style device linking for secure authentication:

1. **Device Registration**: `POST /agents/register`
2. **User Approval**: Web-based device linking
3. **Token Retrieval**: Polling `GET /agents/status`
4. **Authenticated Requests**: Bearer token authentication

### Required Backend Routes

#### 1. **Agent Registration**
```http
POST /agents/register
Content-Type: application/json

{
  "device_name": "Production Server Agent",
  "device_info": {
    "os": "linux",
    "arch": "amd64", 
    "version": "1.0.0",
    "hostname": "prod-server-01",
    "server_id": "web-prod-01"
  }
}
```

**Response**:
```json
{
  "device_code": "abc123def456...",
  "user_code": "XK92-PQ", 
  "verification_url": "https://yourapp.com/agent-link",
  "expires_in": 900,
  "interval": 5
}
```

**Error Responses**:
```json
// 409 Conflict - server_id already exists for this user
{
  "error": "Server ID 'web-prod-01' already exists for this user",
  "suggestion": "Use a different server_id or deactivate the existing agent"
}

// 400 Bad Request - invalid server_id format
{
  "error": "server_id must be 3-50 characters, alphanumeric and hyphens only"
}
```

#### 2. **Device Status Polling**
```http
GET /agents/status?device_code={device_code}
```

**Response** (when approved):
```json
{
  "status": "approved",
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
  "token_type": "Bearer"
}
```

#### 3. **Token Validation**
```http
GET /agents/validate
Authorization: Bearer {access_token}
```

**Response**:
```json
{
  "valid": true,
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

#### 4. **Metrics Submission** ⭐ **PRIMARY ROUTE**
```http
POST /metrics
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "server_id": "web-prod-01",
  "timestamp": 1714392000,
  "metrics": {
    "cpu": {
      "usage_percent": 45.2,
      "per_core": [42.1, 48.3, 44.7, 46.8],
      "cores": 4
    },
    "memory": {
      "total": 16777216000,
      "used": 11408506880,
      "available": 5368709120,
      "used_percent": 68.0,
      "free": 2147483648,
      "cached": 3221225472,
      "buffers": 1073741824
    },
    "disk": [
      {
        "partition": "/",
        "filesystem": "ext4",
        "total": 107374182400,
        "used": 56006140928,
        "free": 46116478976,
        "used_percent": 52.2
      }
    ],
    "network_interfaces": [
      {
        "interface": "eth0",
        "bytes_sent": 1234567890,
        "bytes_recv": 9876543210,
        "packets_sent": 1000000,
        "packets_recv": 2000000,
        "errin": 0,
        "errout": 0,
        "dropin": 0,
        "dropout": 0
      }
    ],
    "network_summary": {
      "total_bytes_sent": 1234567890,
      "total_bytes_recv": 9876543210,
      "total_packets_sent": 1000000,
      "total_packets_recv": 2000000,
      "active_interfaces": 1
    },
    "connections": {
      "total": 156,
      "tcp": 89,
      "udp": 67,
      "listening": 23,
      "established": 45,
      "by_state": {
        "LISTEN": 23,
        "ESTABLISHED": 45,
        "TIME_WAIT": 12,
        "CLOSE_WAIT": 9
      },
      "by_family": {
        "IPv4": 134,
        "IPv6": 22
      }
    },
    "listening_ports": [22, 80, 443, 3306, 5432],
    "port_checks": [
      {
        "port": 80,
        "name": "HTTP",
        "host": "localhost",
        "is_up": true,
        "response_time_ms": 5.234,
        "checked_at": "2024-04-29T10:30:00Z"
      },
      {
        "port": 443,
        "name": "HTTPS", 
        "host": "mysite.com",
        "is_up": false,
        "response_time_ms": 0,
        "error": "connection refused",
        "checked_at": "2024-04-29T10:30:00Z"
      }
    ],
    "port_summary": {
      "total_ports": 2,
      "ports_up": 1,
      "ports_down": 1,
      "avg_response_time_ms": 5.234
    },
    "latency_checks": [
      {
        "host": "8.8.8.8",
        "name": "Google DNS",
        "type": "tcp",
        "latency_ms": 23.456,
        "is_successful": true,
        "checked_at": "2024-04-29T10:30:00Z"
      },
      {
        "host": "https://api.myservice.com",
        "name": "My API",
        "type": "http", 
        "latency_ms": 156.789,
        "is_successful": true,
        "checked_at": "2024-04-29T10:30:00Z"
      }
    ],
    "latency_summary": {
      "total_checks": 2,
      "successful": 2,
      "failed": 0,
      "avg_latency_ms": 90.123,
      "min_latency_ms": 23.456,
      "max_latency_ms": 156.789
    },
    "system": {
      "uptime": 1648717,
      "hostname": "prod-server-01",
      "os": "linux",
      "platform": "ubuntu",
      "platform_version": "20.04",
      "architecture": "amd64",
      "boot_time": 1714220000
    },
    "swap": {
      "total": 2147483648,
      "used": 0,
      "free": 2147483648,
      "used_percent": 0.0
    },
    "disk_io": {
      "read_bytes": 12345678901234,
      "write_bytes": 98765432109876,
      "read_count": 1000000,
      "write_count": 500000
    }
  }
}
```

**Expected Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "success",
  "message": "Metrics received",
  "timestamp": 1714392000
}
```

### Backend Processing Requirements

#### 1. **User-Agent Ownership System**
- **Database Schema**: Each agent is linked to a specific user account
- **Server ID Uniqueness**: Per-user unique (different users can use same server_id)
- **Access Control**: Users can only access their own agents and metrics
- **Token Validation**: JWT tokens contain user_id, device_id, and server_id

**Database Structure**:
```sql
-- Agent devices table
CREATE TABLE agent_devices (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    server_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    access_token VARCHAR(500) UNIQUE,
    device_info JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    last_seen_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    
    -- Ensure server_id uniqueness per user
    UNIQUE(user_id, server_id)
);

-- Metrics data with ownership
CREATE TABLE metrics_data (
    id UUID PRIMARY KEY,
    agent_device_id UUID NOT NULL REFERENCES agent_devices(id),
    server_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL, -- Denormalized for fast queries
    timestamp TIMESTAMP NOT NULL,
    metrics_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**JWT Token Structure**:
```json
{
  "device_id": "uuid-of-agent-device",
  "user_id": "uuid-of-owner",
  "server_id": "web-prod-01",
  "iat": 1714392000,
  "exp": 1745928000
}
```

#### 2. **Data Storage**
- **Time Series Database**: Store metrics with timestamps for historical analysis
- **Metadata Storage**: Store server information, user associations, agent tokens
- **Indexing**: Index by server_id, timestamp, metric type for fast queries

#### 2. **Data Storage**
- **Time Series Database**: Store metrics with timestamps for historical analysis
- **Metadata Storage**: Store server information, user associations, agent tokens
- **Indexing**: Index by user_id, server_id, timestamp, metric type for fast queries
- **Multi-Tenancy**: All queries filtered by user_id for data isolation

#### 3. **Data Validation**
- **Schema Validation**: Validate incoming JSON against expected schema
- **Range Validation**: Ensure percentages are 0-100, bytes are positive, etc.
- **Timestamp Validation**: Ensure timestamps are recent and reasonable
- **Ownership Validation**: Verify server_id matches the authenticated agent's registration

#### 4. **Authentication & Authorization**
- **Token Validation**: Verify Bearer tokens on every request
- **Server ID Verification**: Ensure agent can only submit metrics for its registered server_id
- **User Isolation**: Users can only access their own agents and metrics
- **Rate Limiting**: Prevent abuse (e.g., 1 request per second per agent)

**Backend Validation Example**:
```javascript
// Validate metrics submission
app.post('/metrics', authenticateAgent, async (req, res) => {
    const { server_id, metrics } = req.body;
    const { device_id, user_id, server_id: tokenServerId } = req.agent;
    
    // Verify server_id matches token
    if (tokenServerId !== server_id) {
        return res.status(403).json({
            error: 'server_id mismatch',
            expected: tokenServerId,
            received: server_id
        });
    }
    
    // Store metrics with ownership
    await db.query(`
        INSERT INTO metrics_data (agent_device_id, server_id, user_id, timestamp, metrics_data)
        VALUES ($1, $2, $3, $4, $5)
    `, [device_id, server_id, user_id, new Date(timestamp * 1000), JSON.stringify(metrics)]);
    
    res.json({ status: 'success' });
});
```

#### 5. **Error Handling**
- **400 Bad Request**: Invalid JSON or missing required fields
- **401 Unauthorized**: Invalid or expired token
- **403 Forbidden**: server_id mismatch or access denied
- **409 Conflict**: server_id already exists for user (during registration)
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Database or processing errors

#### 6. **Real-time Features**
- **WebSocket Streaming**: Stream metrics to frontend dashboards
- **Alerting**: Trigger alerts based on thresholds
- **Aggregation**: Calculate averages, trends, and summaries

### Agent Communication Behavior

#### **Retry Logic**
- **Exponential Backoff**: 1s, 2s, 4s, 8s, 16s intervals
- **Max Retries**: 5 attempts per metrics submission
- **Auth Errors**: No retry on 401 (triggers re-authentication)
- **Network Errors**: Full retry sequence

#### **Error Recovery**
- **Token Expiration**: Automatic re-authentication flow
- **Network Failures**: Continue collecting, retry sending
- **Partial Failures**: Skip failed metrics, send successful ones

#### **Performance Characteristics**
- **Collection Time**: ~1-2 seconds for all metrics
- **Memory Usage**: <50MB typical
- **CPU Usage**: <5% during collection
- **Network Usage**: ~1-5KB per metrics payload

---

## 📈 **Monitoring Dashboard Integration**

### Real-time Metrics Display
The backend should provide WebSocket endpoints for real-time dashboard updates:

```javascript
// Frontend WebSocket connection
const ws = new WebSocket('wss://v2-server.watchup.site/ws/metrics');
ws.onmessage = (event) => {
  const metrics = JSON.parse(event.data);
  updateDashboard(metrics);
};
```

### Historical Data Queries
REST endpoints for historical analysis:

```http
GET /api/metrics/{server_id}?from=2024-04-29T00:00:00Z&to=2024-04-29T23:59:59Z&metric=cpu
Authorization: Bearer {user_jwt_token}
```

**Note**: Users can only query metrics from their own agents. The backend validates ownership.

### User Agent Management APIs

#### **List User's Agents**
```http
GET /api/agents
Authorization: Bearer {user_jwt_token}
```

**Response**:
```json
{
  "agents": [
    {
      "server_id": "web-prod-01",
      "device_name": "Production Web Server",
      "last_seen_at": "2024-04-29T10:30:00Z",
      "is_active": true,
      "created_at": "2024-04-29T09:00:00Z",
      "metrics_count": 1440
    },
    {
      "server_id": "db-server-main",
      "device_name": "Main Database Server",
      "last_seen_at": "2024-04-29T10:29:45Z",
      "is_active": true,
      "created_at": "2024-04-28T14:30:00Z",
      "metrics_count": 2880
    }
  ]
}
```

#### **Get Specific Agent Details**
```http
GET /api/agents/{server_id}
Authorization: Bearer {user_jwt_token}
```

**Response**:
```json
{
  "server_id": "web-prod-01",
  "device_name": "Production Web Server",
  "device_info": {
    "os": "linux",
    "arch": "amd64",
    "hostname": "prod-server-01",
    "version": "1.0.0"
  },
  "created_at": "2024-04-29T09:00:00Z",
  "last_seen_at": "2024-04-29T10:30:00Z",
  "is_active": true,
  "metrics_summary": {
    "total_metrics_received": 1440,
    "last_cpu_usage": 45.2,
    "last_memory_usage": 68.1,
    "last_disk_usage": 52.3
  }
}
```

#### **Deactivate Agent**
```http
DELETE /api/agents/{server_id}
Authorization: Bearer {user_jwt_token}
```

**Response**:
```json
{
  "message": "Agent 'web-prod-01' has been deactivated",
  "server_id": "web-prod-01"
}
```

#### **Get Agent Metrics**
```http
GET /api/agents/{server_id}/metrics?limit=100&metric=cpu
Authorization: Bearer {user_jwt_token}
```

**Response**:
```json
{
  "server_id": "web-prod-01",
  "metrics": [
    {
      "timestamp": "2024-04-29T10:30:00Z",
      "cpu": {
        "usage_percent": 45.2,
        "cores": 4
      }
    }
  ],
  "pagination": {
    "limit": 100,
    "total": 1440,
    "has_more": true
  }
}
```

### Alerting Integration
Threshold-based alerting:

```json
{
  "server_id": "web-prod-01",
  "alert_type": "cpu_high",
  "threshold": 90.0,
  "current_value": 95.2,
  "timestamp": "2024-04-29T10:30:00Z"
}
```

---

## � **Server ID & Ownership Management**

### Server ID Rules

#### **Format Requirements**:
- **Length**: 3-50 characters
- **Characters**: Alphanumeric, hyphens (-), and underscores (_) only
- **Start/End**: Cannot start or end with hyphen or underscore
- **Case**: Case-sensitive (recommended: lowercase with hyphens)

#### **Examples**:
- ✅ `web-prod-01` - Good format
- ✅ `db_server_main` - Acceptable with underscores
- ✅ `api-gateway-1` - Clear and descriptive
- ❌ `w` - Too short (minimum 3 characters)
- ❌ `web-prod-01!` - Invalid character (!)
- ❌ `-web-prod` - Cannot start with hyphen
- ❌ `web-prod-` - Cannot end with hyphen

#### **Uniqueness Rules**:
- **Per-User Unique**: Each user can have only ONE agent with a specific server_id
- **Cross-User Allowed**: Different users can use the same server_id
- **Immutable**: server_id cannot be changed after agent linking

#### **Validation Examples**:
```javascript
// User A registers: server_id = "web-prod-01" ✅ Allowed
// User B registers: server_id = "web-prod-01" ✅ Allowed (different user)
// User A registers: server_id = "web-prod-01" ❌ Rejected (duplicate for User A)
```

### Ownership Tracking

#### **Database Relationships**:
```sql
-- Users own multiple agents
users (1) ←→ (many) agent_devices

-- Agents generate multiple metrics
agent_devices (1) ←→ (many) metrics_data

-- Complete ownership chain
users → agent_devices → metrics_data
```

#### **Access Control Matrix**:
| Action | User A | User B | Admin |
|--------|--------|--------|-------|
| View User A's agents | ✅ | ❌ | ✅ |
| View User A's metrics | ✅ | ❌ | ✅ |
| Deactivate User A's agent | ✅ | ❌ | ✅ |
| Register with existing server_id | ❌* | ✅ | ✅ |

*Only if server_id already exists for User A

#### **Backend Ownership Validation**:
```javascript
// Middleware to validate agent ownership
async function validateAgentOwnership(req, res, next) {
    const { server_id } = req.params;
    const user_id = req.user.id;
    
    const agent = await db.query(`
        SELECT id, device_name FROM agent_devices 
        WHERE user_id = $1 AND server_id = $2 AND is_active = true
    `, [user_id, server_id]);
    
    if (!agent.length) {
        return res.status(404).json({
            error: 'Agent not found',
            message: `No active agent with server_id '${server_id}' found for your account`
        });
    }
    
    req.agent = agent[0];
    next();
}

// Usage in routes
app.get('/api/agents/:server_id', authenticateUser, validateAgentOwnership, (req, res) => {
    res.json(req.agent);
});
```

### Multi-Tenant Data Isolation

#### **Query Patterns**:
```sql
-- Always filter by user_id for data isolation

-- Get user's agents
SELECT * FROM agent_devices WHERE user_id = $1 AND is_active = true;

-- Get user's metrics
SELECT m.* FROM metrics_data m
JOIN agent_devices a ON m.agent_device_id = a.id
WHERE a.user_id = $1 AND a.server_id = $2;

-- Aggregate user's metrics
SELECT 
    DATE_TRUNC('hour', m.timestamp) as hour,
    AVG((m.metrics_data->>'cpu'->>'usage_percent')::float) as avg_cpu
FROM metrics_data m
JOIN agent_devices a ON m.agent_device_id = a.id
WHERE a.user_id = $1 
GROUP BY hour
ORDER BY hour DESC;
```

#### **Performance Indexes**:
```sql
-- Essential indexes for multi-tenant queries
CREATE INDEX idx_agent_devices_user_server ON agent_devices(user_id, server_id);
CREATE INDEX idx_agent_devices_user_active ON agent_devices(user_id, is_active);
CREATE INDEX idx_metrics_data_agent_timestamp ON metrics_data(agent_device_id, timestamp);
CREATE INDEX idx_metrics_data_user_timestamp ON metrics_data(user_id, timestamp);
```

---

## �🔧 **Configuration Management**

### Agent Configuration
```yaml
# Core settings
server_id: "web-prod-01"
endpoint: "https://v2-server.watchup.site"
interval: 5s

# Metric toggles
metrics:
  cpu: true
  memory: true
  disk: true
  network: true
  connections: false  # Optional

# Extended monitoring
ports:
  - port: 80
    name: "HTTP"
    host: "localhost"
    timeout: "5s"

latency_checks:
  - host: "8.8.8.8"
    name: "Google DNS"
    type: "tcp"
    port: 53
    timeout: "3s"
```

### Backend Configuration
- **Database**: Time-series database (InfluxDB, TimescaleDB, etc.) with multi-tenant support
- **Authentication**: JWT tokens with user_id, device_id, and server_id claims
- **Rate Limiting**: Configurable per-user/per-agent limits
- **Retention**: Configurable data retention policies per user plan
- **Indexing**: Optimized for multi-tenant queries (user_id + server_id + timestamp)
- **Backup**: User data isolation in backup/restore procedures

**Environment Variables**:
```bash
# Database
DATABASE_URL=postgresql://user:pass@localhost/watchup
REDIS_URL=redis://localhost:6379

# Authentication
JWT_SECRET=your-secret-key
JWT_EXPIRATION=1y

# Rate Limiting
RATE_LIMIT_PER_USER=1000  # requests per hour per user
RATE_LIMIT_PER_AGENT=60   # requests per minute per agent

# Features
MAX_AGENTS_FREE=3         # Free plan limit
MAX_AGENTS_PRO=25         # Pro plan limit
MAX_AGENTS_ENTERPRISE=1000 # Enterprise plan limit
```

---

## 🎯 **Production Deployment**

### Agent Deployment
1. **Binary Distribution**: Single executable file
2. **Configuration**: YAML configuration file
3. **Service Management**: systemd service (Linux) or Windows Service
4. **Logging**: Structured logging to stdout/file
5. **Monitoring**: Self-monitoring and health checks

### Backend Requirements
1. **Scalability**: Handle thousands of concurrent agents
2. **Reliability**: High availability and fault tolerance
3. **Security**: HTTPS, token validation, rate limiting
4. **Performance**: Sub-second response times
5. **Monitoring**: Backend metrics and alerting

---

## 📊 **Summary**

The WatchUp Agent provides **comprehensive system monitoring** with **complete user ownership tracking**:

### **Monitoring Capabilities**:
- ✅ **8 Core Metric Categories** (CPU, Memory, Disk, Network, Connections, Ports, Latency, System)
- ✅ **50+ Individual Metrics** collected in real-time
- ✅ **Configurable Monitoring** (enable/disable features as needed)
- ✅ **Production-Ready** (authentication, retry logic, error handling)
- ✅ **Cross-Platform** (Windows, Linux, macOS)
- ✅ **Lightweight** (<5% CPU, <50MB RAM)

### **Ownership & Security**:
- ✅ **User-Agent Linking** - Each agent tied to specific user account
- ✅ **Server ID Uniqueness** - Per-user unique identifiers
- ✅ **Multi-Tenant Isolation** - Complete data separation between users
- ✅ **Access Control** - Users can only access their own agents/metrics
- ✅ **Token-Based Security** - JWT tokens with ownership claims
- ✅ **API Validation** - server_id verification on every request

### **Backend Integration**:
- ✅ **4 Core API Endpoints** - Registration, polling, validation, metrics
- ✅ **Ownership Validation** - All requests validated against user ownership
- ✅ **Error Handling** - Comprehensive error codes and messages
- ✅ **Rate Limiting** - Per-user and per-agent limits
- ✅ **Real-time Streaming** - WebSocket support for dashboards
- ✅ **Historical Queries** - Time-series data with user filtering

The agent integrates seamlessly with the WatchUp backend via secure HTTPS API, providing **real-time system visibility** with **complete ownership tracking** and **multi-tenant data isolation**.
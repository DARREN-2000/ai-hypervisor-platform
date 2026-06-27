# API Reference

The AI Hypervisor Platform exposes a RESTful API for managing virtual machines, GPUs, hosts, and tasks. All endpoints are versioned under `/api/v1`.

## Authentication

Authentication is configured via the `APIKeyRequired` and `JWTEnabled` settings in the configuration. When enabled, include the authentication token in the `Authorization` header:

```http
Authorization: Bearer <your-token>
```

## Rate Limits

Rate limiting is enforced globally and per-IP based on your deployment configuration. If you exceed the rate limit, the API will respond with `429 Too Many Requests`.

## Base URL

```
https://api.ai-hypervisor.example.com/api/v1
```

## Endpoints

### Virtual Machines

#### List VMs

**Purpose:** Retrieve a list of all virtual machines in the cluster.

**Request:**
`GET /vms`

**Response (200 OK):**
```json
[
  {
    "id": "vm-123",
    "name": "inference-node-1",
    "status": "running",
    "host_id": "host-abc",
    "cpu": 8,
    "memory_mb": 32768,
    "gpus": ["gpu-xyz"]
  }
]
```

**Example:**
Retrieving the list of VMs to display on a dashboard or CLI.

**Error codes:**
- `401 Unauthorized`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 100 requests per minute.

**Best practices:** Implement pagination if managing a large number of VMs.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/vms
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/vms", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const vms = await response.json();
```

#### Create VM

**Purpose:** Provision a new virtual machine.

**Request:**
`POST /vms`

**Response (202 Accepted):**
```json
{
  "task_id": "task-456",
  "status": "pending"
}
```

**Example:**
Provisioning a new node for training an ML model.

**Error codes:**
- `400 Bad Request`
- `401 Unauthorized`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 20 requests per minute.

**Best practices:** Verify host capacity and quotas before provisioning.

**Example curl:**
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"name": "inference-node-2", "cpu": 16, "memory_mb": 65536, "gpu_count": 2, "image": "ubuntu-22.04-cuda12"}' https://api.ai-hypervisor.example.com/api/v1/vms
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
payload = {"name": "inference-node-2", "cpu": 16, "memory_mb": 65536, "gpu_count": 2, "image": "ubuntu-22.04-cuda12"}
response = requests.post("https://api.ai-hypervisor.example.com/api/v1/vms", headers=headers, json=payload)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms", {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    name: "inference-node-2",
    cpu: 16,
    memory_mb: 65536,
    gpu_count: 2,
    image: "ubuntu-22.04-cuda12"
  })
});
const result = await response.json();
```

#### Get VM Details

**Purpose:** Retrieve details for a specific virtual machine.

**Request:**
`GET /vms/{vmId}`

**Response (200 OK):**
```json
{
  "id": "vm-123",
  "name": "inference-node-1",
  "status": "running",
  "host_id": "host-abc",
  "cpu": 8,
  "memory_mb": 32768,
  "gpus": ["gpu-xyz"]
}
```

**Example:**
Checking the status of a specific VM.

**Error codes:**
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 100 requests per minute.

**Best practices:** Cache VM details if queried frequently.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/vms/vm-123
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const vm = await response.json();
```

#### Start VM

**Purpose:** Start a stopped virtual machine.

**Request:**
`POST /vms/{vmId}/start`

**Response (202 Accepted):**
```json
{
  "task_id": "task-789",
  "status": "pending"
}
```

**Example:**
Starting a VM that was stopped to save resources.

**Error codes:**
- `400 Bad Request`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 20 requests per minute.

**Best practices:** Verify VM is currently in stopped state before initiating.

**Example curl:**
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/start
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.post("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/start", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/start", {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const result = await response.json();
```

#### Stop VM

**Purpose:** Stop a running virtual machine safely.

**Request:**
`POST /vms/{vmId}/stop`

**Response (202 Accepted):**
```json
{
  "task_id": "task-101",
  "status": "pending"
}
```

**Example:**
Stopping a VM when it is no longer needed.

**Error codes:**
- `400 Bad Request`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 20 requests per minute.

**Best practices:** Verify VM is currently in running state before initiating.

**Example curl:**
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/stop
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.post("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/stop", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/stop", {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const result = await response.json();
```

#### Reboot VM

**Purpose:** Reboot a running virtual machine.

**Request:**
`POST /vms/{vmId}/reboot`

**Response (202 Accepted):**
```json
{
  "task_id": "task-202",
  "status": "pending"
}
```

**Example:**
Rebooting a VM to apply updates or recover from a frozen state.

**Error codes:**
- `400 Bad Request`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 20 requests per minute.

**Best practices:** Avoid frequent reboots to maintain workload continuity.

**Example curl:**
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/reboot
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.post("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/reboot", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123/reboot", {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const result = await response.json();
```

#### Update VM

**Purpose:** Modify an existing virtual machine's properties.

**Request:**
`PATCH /vms/{vmId}`

**Response (200 OK):**
```json
{
  "id": "vm-123",
  "name": "inference-node-1-updated",
  "status": "running"
}
```

**Example:**
Updating the name or description of a VM.

**Error codes:**
- `400 Bad Request`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 50 requests per minute.

**Best practices:** Validate input payload fields correctly.

**Example curl:**
```bash
curl -X PATCH -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"name": "inference-node-1-updated"}' https://api.ai-hypervisor.example.com/api/v1/vms/vm-123
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
payload = {"name": "inference-node-1-updated"}
response = requests.patch("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123", headers=headers, json=payload)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123", {
  method: "PATCH",
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    name: "inference-node-1-updated"
  })
});
const result = await response.json();
```

#### Delete VM

**Purpose:** Permanently remove a virtual machine.

**Request:**
`DELETE /vms/{vmId}`

**Response (202 Accepted):**
```json
{
  "task_id": "task-303",
  "status": "pending"
}
```

**Example:**
Destroying a VM when the instance is no longer required.

**Error codes:**
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 10 requests per minute.

**Best practices:** Implement a "soft delete" or user confirmation in UI clients.

**Example curl:**
```bash
curl -X DELETE -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/vms/vm-123
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.delete("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/vms/vm-123", {
  method: "DELETE",
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const result = await response.json();
```

### GPUs

#### List GPUs

**Purpose:** Retrieve a list of all GPUs available in the cluster.

**Request:**
`GET /gpus`

**Response (200 OK):**
```json
[
  {
    "id": "gpu-xyz",
    "name": "NVIDIA A100",
    "host_id": "host-abc",
    "status": "available",
    "memory_mb": 40960
  }
]
```

**Example:**
Listing available GPUs for an orchestration engine or user dashboard.

**Error codes:**
- `401 Unauthorized`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 100 requests per minute.

**Best practices:** Add filters by host or model to narrow results.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/gpus
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/gpus", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/gpus", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const gpus = await response.json();
```

#### Get GPU Details

**Purpose:** Retrieve details for a specific GPU.

**Request:**
`GET /gpus/{gpuId}`

**Response (200 OK):**
```json
{
  "id": "gpu-xyz",
  "name": "NVIDIA A100",
  "host_id": "host-abc",
  "status": "available",
  "memory_mb": 40960
}
```

**Example:**
Inspecting GPU characteristics prior to workload submission.

**Error codes:**
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 100 requests per minute.

**Best practices:** Poll less frequently to avoid overwhelming the metrics collector.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/gpus/gpu-xyz
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/gpus/gpu-xyz", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/gpus/gpu-xyz", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const gpu = await response.json();
```

### Hosts

#### List Hosts

**Purpose:** Retrieve a list of all compute nodes (hosts) in the cluster.

**Request:**
`GET /hosts`

**Response (200 OK):**
```json
[
  {
    "id": "host-abc",
    "hostname": "node-01",
    "status": "ready",
    "cpu_cores": 64,
    "memory_mb": 262144
  }
]
```

**Example:**
Monitoring the health of underlying bare-metal machines.

**Error codes:**
- `401 Unauthorized`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 100 requests per minute.

**Best practices:** Implement caching for periodic polling.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/hosts
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/hosts", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/hosts", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const hosts = await response.json();
```

#### Get Host Details

**Purpose:** Retrieve details for a specific host node.

**Request:**
`GET /hosts/{nodeId}`

**Response (200 OK):**
```json
{
  "id": "host-abc",
  "hostname": "node-01",
  "status": "ready",
  "cpu_cores": 64,
  "memory_mb": 262144
}
```

**Example:**
Retrieving configuration specs of a specific node.

**Error codes:**
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 100 requests per minute.

**Best practices:** Use host telemetry rather than static APIs for real-time monitoring.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/hosts/host-abc
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/hosts/host-abc", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/hosts/host-abc", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const host = await response.json();
```

#### Get Host Metrics

**Purpose:** Retrieve live telemetry metrics for a host.

**Request:**
`GET /hosts/{nodeId}/metrics`

**Response (200 OK):**
```json
{
  "cpu_usage_percent": 45.2,
  "memory_usage_mb": 131072,
  "gpu_utilization_percent": 88.5
}
```

**Example:**
Displaying a live node health widget on a dashboard.

**Error codes:**
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 200 requests per minute.

**Best practices:** Prefer WebSockets for real-time telemetry if available.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/hosts/host-abc/metrics
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/hosts/host-abc/metrics", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/hosts/host-abc/metrics", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const metrics = await response.json();
```

### Tasks

#### Get Task Details

**Purpose:** Retrieve the status of an asynchronous background task.

**Request:**
`GET /tasks/{taskId}`

**Response (200 OK):**
```json
{
  "id": "task-456",
  "status": "completed",
  "created_at": "2023-10-01T12:00:00Z",
  "completed_at": "2023-10-01T12:05:00Z"
}
```

**Example:**
Polling the API to determine when a VM finishes provisioning.

**Error codes:**
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

**Authentication:** Requires valid JWT or API Key if enabled.

**Rate limits:** 300 requests per minute.

**Best practices:** Use exponential backoff when polling task status.

**Example curl:**
```bash
curl -H "Authorization: Bearer $TOKEN" https://api.ai-hypervisor.example.com/api/v1/tasks/task-456
```

**Python example:**
```python
import requests

headers = {"Authorization": f"Bearer {token}"}
response = requests.get("https://api.ai-hypervisor.example.com/api/v1/tasks/task-456", headers=headers)
print(response.json())
```

**TypeScript example:**
```typescript
const response = await fetch("https://api.ai-hypervisor.example.com/api/v1/tasks/task-456", {
  headers: {
    Authorization: `Bearer ${token}`
  }
});
const task = await response.json();
```

### WebSockets

#### Cluster Events

**Purpose:** Stream real-time events for the entire cluster.

**Request:**
`GET /ws/cluster/events`

**Response:**
Streams JSON event objects.

**Example:**
```json
{
  "type": "vm_created",
  "timestamp": "2023-10-01T12:00:00Z",
  "data": { "vm_id": "vm-123" }
}
```

**Error codes:**
- `401 Unauthorized` (if connection rejected)

**Authentication:** Requires valid JWT or API Key passed via query parameter or header depending on client capability.

**Rate limits:** Rate limits do not apply to open socket connections, but connections per user may be capped.

**Best practices:** Implement automatic reconnection logic with backoff.

**Example curl:** N/A (WebSocket)

**Python example:**
```python
import asyncio
import websockets

async def listen():
    async with websockets.connect("wss://api.ai-hypervisor.example.com/api/v1/ws/cluster/events") as ws:
        while True:
            event = await ws.recv()
            print(event)

asyncio.run(listen())
```

**TypeScript example:**
```typescript
const ws = new WebSocket("wss://api.ai-hypervisor.example.com/api/v1/ws/cluster/events?token=" + token);
ws.onmessage = (event) => {
  console.log(JSON.parse(event.data));
};
```

#### VM Metrics Stream

**Purpose:** Stream real-time metrics for a specific virtual machine.

**Request:**
`GET /ws/vm/{vmId}/metrics`

**Response:**
Streams JSON metric objects.

**Example:**
```json
{
  "timestamp": "2023-10-01T12:00:00Z",
  "cpu_usage": 12.5,
  "memory_usage_mb": 1024
}
```

**Error codes:**
- `401 Unauthorized`
- `404 Not Found` (if VM does not exist)

**Authentication:** Requires valid JWT or API Key.

**Rate limits:** N/A for open streams.

**Best practices:** Aggregate or sample metrics on the client side if the event rate is too high.

**Example curl:** N/A (WebSocket)

**Python example:**
```python
import asyncio
import websockets

async def listen():
    async with websockets.connect(f"wss://api.ai-hypervisor.example.com/api/v1/ws/vm/vm-123/metrics?token={token}") as ws:
        while True:
            metrics = await ws.recv()
            print(metrics)

asyncio.run(listen())
```

**TypeScript example:**
```typescript
const ws = new WebSocket(`wss://api.ai-hypervisor.example.com/api/v1/ws/vm/vm-123/metrics?token=${token}`);
ws.onmessage = (event) => {
  console.log(JSON.parse(event.data));
};
```

## 2024-05-24 - [Do not leak internal errors]
**Vulnerability:** Information leakage. Exposing raw internal errors via `err.Error()` in API response payload.
**Learning:** Found raw `err.Error()` being passed to HTTP clients, which can contain database connection strings, paths, system data or stack traces.
**Prevention:** Check error type (e.g. `apierrors.APIError` here which is intended to be exposed safely). If the error isn't of a safe type, log the raw error on the backend to avoid losing debugging data, and send a generic "Internal server error" message to the client.

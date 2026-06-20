## 2024-05-24 - [Do not leak internal errors]
**Vulnerability:** Information leakage. Exposing raw internal errors via `err.Error()` in API response payload.
**Learning:** Found raw `err.Error()` being passed to HTTP clients, which can contain database connection strings, paths, system data or stack traces.
**Prevention:** Check error type (e.g. `apierrors.APIError` here which is intended to be exposed safely). If the error isn't of a safe type, log the raw error on the backend to avoid losing debugging data, and send a generic "Internal server error" message to the client.

## 2024-05-24 - [Do not leak internal errors in unauthenticated health checks]
**Vulnerability:** Information leakage. Exposing raw internal errors via `err.Error()` in API response payload.
**Learning:** Found raw `err.Error()` being passed to HTTP clients, which can contain database connection strings, paths, system data or stack traces.
**Prevention:** Mask raw `err.Error()` strings in public/unauthenticated endpoints like health checks to prevent leakage of internal details.

## 2024-06-20 - [Do not leak internal errors to Task results]
**Vulnerability:** Exposing internal error stack traces and internal messages inside `task` status and `vm.ErrorMessage`.
**Learning:** Found that when a `VirtualMachine` was provisioned or failed to provision, it would copy `err.Error()` directly into the Task's `error` interface and `vm.ErrorMessage`. This is passed back to clients who query the VM task status.
**Prevention:** Always use a helper to mask the raw errors and safely extract the error message from the client-facing standard `apierrors.APIError` type, falling back to a generic `internal error` message.

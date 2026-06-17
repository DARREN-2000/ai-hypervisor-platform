## 2024-05-18 - Information Leakage in API Error Handler
**Vulnerability:** The API error handler (`internal/api/handlers/base.go`) was returning the raw `err.Error()` in a `message` field when an error was not of type `*apierrors.APIError`. This could potentially leak internal details, stack traces, or other sensitive information to the client.
**Learning:** Returning raw error messages directly to the client is a common anti-pattern that violates the "fail securely" principle and can expose internal architecture to attackers.
**Prevention:** Always log the detailed, raw error on the server side and return a generic, non-descriptive error message (e.g., "An unexpected error occurred") to the client for unhandled exceptions.

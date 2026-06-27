# Threat Model and Security Practices

This document outlines the security architecture and threat model for the AI Hypervisor Platform.

## Authentication & Authorization

- The platform supports JWT-based authentication.
- API requests must include a valid Bearer token if `JWTEnabled` is true.
- Role-Based Access Control (RBAC) is implemented to restrict access to namespaces and actions.

## Network Security

- All inter-service communication (NATS) should be secured using TLS.
- The API Server enforces standard security headers via `securityHeadersMiddleware`.
- Host Agents run on bare-metal and must be isolated via strict firewall rules, only allowing traffic from the NATS server and Prometheus scrapers.

## Secrets Management

- Database credentials and API keys should be injected via environment variables or Kubernetes Secrets.
- Never hardcode secrets in configuration files.

## Vulnerability Reporting

If you discover a security vulnerability, please follow the guidelines in `SECURITY.md`.
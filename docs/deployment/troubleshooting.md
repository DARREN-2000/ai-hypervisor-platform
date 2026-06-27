# Troubleshooting Guide

This guide provides solutions to common issues encountered when deploying or operating the AI Hypervisor Platform.

## Installation Issues

### `make build` Fails

**Symptom:**
```
go build: command not found
```
**Solution:**
Ensure that Go 1.21+ is installed and configured in your `$PATH`.

### Docker Compose Fails

**Symptom:**
```
docker-compose: command not found
```
**Solution:**
Install Docker and Docker Compose. Ensure the daemon is running.

## Runtime Issues

### API Server Cannot Connect to PostgreSQL

**Symptom:** API server logs show:
```
failed to connect to postgresql
```
**Solution:**
Verify your database settings in `config/config.yaml`. Ensure PostgreSQL is running and accessible. Check firewall rules if connecting remotely.

### Host Agent Fails to Connect to Libvirt

**Symptom:** Host Agent logs show:
```
failed to connect to libvirt
```
**Solution:**
Ensure libvirtd is running on the host node. The Host Agent needs appropriate permissions to interact with libvirt.

## API Issues

### `401 Unauthorized` Responses

**Symptom:** All API requests return `401`.
**Solution:**
Check if `JWTEnabled` or `APIKeyRequired` is `true` in your configuration. If so, ensure you are passing a valid `Authorization: Bearer <token>` header.

## Further Assistance

If your issue is not listed here, please consult our [FAQ](docs/getting-started/faq.md) or open an issue on GitHub.

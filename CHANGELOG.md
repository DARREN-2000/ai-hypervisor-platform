# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of the AI Hypervisor Platform.
- Comprehensive documentation overhaul including architecture, API, and developer guides.
- GitHub community templates (`ISSUE_TEMPLATE`, `PULL_REQUEST_TEMPLATE`, `SECURITY.md`, `CODEOWNERS`, etc.).
- RESTful API endpoints for managing VMs, GPUs, Hosts, and Tasks.
- WebSocket support for cluster events and VM telemetry.
- Core scheduler with bin-packing, spreading, and NUMA-aware allocation.
- Host Agent for libvirt and NVML integration.

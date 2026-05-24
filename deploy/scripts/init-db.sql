-- AI Hypervisor Platform Database Schema
-- PostgreSQL initialization script
-- This script creates all necessary tables and indexes

BEGIN;

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Virtual Machines table
CREATE TABLE vms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    namespace VARCHAR(255) NOT NULL DEFAULT 'default',
    flavor_name VARCHAR(255) NOT NULL,
    flavor_cpu INTEGER NOT NULL,
    flavor_memory INTEGER NOT NULL,
    flavor_disk_size INTEGER NOT NULL,
    image_id VARCHAR(255) NOT NULL,
    state VARCHAR(50) NOT NULL DEFAULT 'created',
    host_node_id UUID,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    terminated_at TIMESTAMP,
    last_status_check TIMESTAMP NOT NULL DEFAULT NOW(),
    error_message TEXT,
    CONSTRAINT fk_host_node FOREIGN KEY (host_node_id) REFERENCES host_nodes(id),
    CONSTRAINT unique_vm_name UNIQUE (namespace, name)
);

CREATE INDEX idx_vms_state ON vms(state);
CREATE INDEX idx_vms_namespace ON vms(namespace);
CREATE INDEX idx_vms_host_node_id ON vms(host_node_id);
CREATE INDEX idx_vms_created_at ON vms(created_at);

-- Host Nodes table
CREATE TABLE host_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'ready',
    capacity_cpu INTEGER NOT NULL,
    capacity_memory INTEGER NOT NULL,
    capacity_disk_gb INTEGER NOT NULL,
    capacity_gpu_slots INTEGER NOT NULL DEFAULT 0,
    allocated_cpu INTEGER NOT NULL DEFAULT 0,
    allocated_memory INTEGER NOT NULL DEFAULT 0,
    allocated_disk_gb INTEGER NOT NULL DEFAULT 0,
    allocated_gpu_slots INTEGER NOT NULL DEFAULT 0,
    allocated_vms INTEGER NOT NULL DEFAULT 0,
    metadata JSONB,
    last_heartbeat TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_host_nodes_status ON host_nodes(status);
CREATE INDEX idx_host_nodes_last_heartbeat ON host_nodes(last_heartbeat);

-- GPUs table
CREATE TABLE gpus (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    host_node_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    model VARCHAR(255) NOT NULL,
    vram INTEGER NOT NULL,
    index INTEGER NOT NULL,
    pci_address VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'available',
    metadata JSONB,
    last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_host_node FOREIGN KEY (host_node_id) REFERENCES host_nodes(id),
    CONSTRAINT unique_gpu_pci UNIQUE (host_node_id, pci_address)
);

CREATE INDEX idx_gpus_host_node_id ON gpus(host_node_id);
CREATE INDEX idx_gpus_status ON gpus(status);
CREATE INDEX idx_gpus_type ON gpus(type);

-- GPU Allocations table
CREATE TABLE gpu_allocations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    gpu_id UUID NOT NULL,
    vm_id UUID NOT NULL,
    host_node_id UUID NOT NULL,
    device_index INTEGER NOT NULL,
    pci_address VARCHAR(50),
    mig_mode BOOLEAN DEFAULT FALSE,
    mig_profile VARCHAR(50),
    allocated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_gpu FOREIGN KEY (gpu_id) REFERENCES gpus(id),
    CONSTRAINT fk_vm FOREIGN KEY (vm_id) REFERENCES vms(id),
    CONSTRAINT fk_host_node FOREIGN KEY (host_node_id) REFERENCES host_nodes(id),
    CONSTRAINT unique_gpu_allocation UNIQUE (gpu_id, vm_id)
);

CREATE INDEX idx_gpu_allocations_vm_id ON gpu_allocations(vm_id);
CREATE INDEX idx_gpu_allocations_gpu_id ON gpu_allocations(gpu_id);
CREATE INDEX idx_gpu_allocations_host_node_id ON gpu_allocations(host_node_id);

-- Tasks table for async operations
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'queued',
    vm_id UUID,
    vm_name VARCHAR(255),
    vm_namespace VARCHAR(255),
    payload JSONB,
    result JSONB,
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    CONSTRAINT fk_vm FOREIGN KEY (vm_id) REFERENCES vms(id) ON DELETE SET NULL
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_type ON tasks(type);
CREATE INDEX idx_tasks_vm_id ON tasks(vm_id);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);
CREATE INDEX idx_tasks_completed_at ON tasks(completed_at);

-- Resource Metrics table
CREATE TABLE resource_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vm_id UUID,
    host_node_id UUID,
    metric_type VARCHAR(50) NOT NULL,
    cpu_usage FLOAT,
    memory_usage INTEGER,
    disk_io_read INTEGER,
    disk_io_write INTEGER,
    network_in INTEGER,
    network_out INTEGER,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_vm FOREIGN KEY (vm_id) REFERENCES vms(id) ON DELETE CASCADE,
    CONSTRAINT fk_host_node FOREIGN KEY (host_node_id) REFERENCES host_nodes(id) ON DELETE CASCADE
);

CREATE INDEX idx_resource_metrics_vm_id ON resource_metrics(vm_id, timestamp);
CREATE INDEX idx_resource_metrics_host_node_id ON resource_metrics(host_node_id, timestamp);
CREATE INDEX idx_resource_metrics_timestamp ON resource_metrics(timestamp);

-- Audit Logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    action VARCHAR(100) NOT NULL,
    actor VARCHAR(255) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id UUID,
    changes JSONB,
    status VARCHAR(50) NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource, resource_id);
CREATE INDEX idx_audit_logs_actor ON audit_logs(actor);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- Configurations table
CREATE TABLE configurations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    config_type VARCHAR(100) NOT NULL,
    config_key VARCHAR(255) NOT NULL,
    config_value JSONB NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_config UNIQUE (config_type, config_key, version)
);

CREATE INDEX idx_configurations_type ON configurations(config_type);
CREATE INDEX idx_configurations_active ON configurations(active);

-- VM Network Configurations table
CREATE TABLE vm_network_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vm_id UUID NOT NULL,
    interface_id VARCHAR(255) NOT NULL,
    interface_name VARCHAR(255) NOT NULL,
    mac_address VARCHAR(17),
    ip_address VARCHAR(45),
    netmask VARCHAR(45),
    gateway VARCHAR(45),
    bridge_name VARCHAR(255),
    vlan_id INTEGER,
    mtu INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_vm FOREIGN KEY (vm_id) REFERENCES vms(id) ON DELETE CASCADE,
    CONSTRAINT unique_vm_interface UNIQUE (vm_id, interface_id)
);

CREATE INDEX idx_vm_network_configs_vm_id ON vm_network_configs(vm_id);
CREATE INDEX idx_vm_network_configs_mac ON vm_network_configs(mac_address);

-- VM Storage Configurations table
CREATE TABLE vm_storage_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vm_id UUID NOT NULL,
    volume_id VARCHAR(255) NOT NULL,
    volume_name VARCHAR(255) NOT NULL,
    volume_type VARCHAR(50) NOT NULL,
    source_path TEXT NOT NULL,
    target_device VARCHAR(50),
    format VARCHAR(50),
    size_gb INTEGER,
    readonly BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_vm FOREIGN KEY (vm_id) REFERENCES vms(id) ON DELETE CASCADE,
    CONSTRAINT unique_vm_volume UNIQUE (vm_id, volume_id)
);

CREATE INDEX idx_vm_storage_configs_vm_id ON vm_storage_configs(vm_id);

-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    source VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    data JSONB,
    severity VARCHAR(50) NOT NULL DEFAULT 'info',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_timestamp ON events(timestamp);
CREATE INDEX idx_events_subject ON events(subject);
CREATE INDEX idx_events_severity ON events(severity);

-- Scheduling Decisions table (for audit/analytics)
CREATE TABLE scheduling_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vm_id UUID NOT NULL,
    selected_host_id UUID NOT NULL,
    score FLOAT NOT NULL,
    reason TEXT,
    alternatives JSONB,
    decision_timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_vm FOREIGN KEY (vm_id) REFERENCES vms(id) ON DELETE CASCADE,
    CONSTRAINT fk_host_node FOREIGN KEY (selected_host_id) REFERENCES host_nodes(id)
);

CREATE INDEX idx_scheduling_decisions_vm_id ON scheduling_decisions(vm_id);
CREATE INDEX idx_scheduling_decisions_host_id ON scheduling_decisions(selected_host_id);
CREATE INDEX idx_scheduling_decisions_timestamp ON scheduling_decisions(decision_timestamp);

-- Secrets table (encrypted)
CREATE TABLE secrets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    secret_type VARCHAR(100) NOT NULL,
    secret_key VARCHAR(255) NOT NULL,
    secret_value BYTEA NOT NULL,
    encrypted_by VARCHAR(255),
    owner VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    CONSTRAINT unique_secret UNIQUE (secret_type, secret_key)
);

CREATE INDEX idx_secrets_type ON secrets(secret_type);
CREATE INDEX idx_secrets_key ON secrets(secret_key);
CREATE INDEX idx_secrets_owner ON secrets(owner);
CREATE INDEX idx_secrets_expires_at ON secrets(expires_at);

-- Create views for common queries

-- View: Current VM Status Summary
CREATE VIEW vm_status_summary AS
SELECT
    state,
    COUNT(*) as count,
    SUM(CASE WHEN host_node_id IS NOT NULL THEN 1 ELSE 0 END) as on_hosts,
    SUM(CASE WHEN error_message IS NOT NULL THEN 1 ELSE 0 END) as with_errors
FROM vms
WHERE state != 'deleted'
GROUP BY state;

-- View: Host Utilization
CREATE VIEW host_utilization AS
SELECT
    hn.id,
    hn.name,
    hn.status,
    ROUND(100.0 * hn.allocated_cpu / hn.capacity_cpu, 2) as cpu_utilization_pct,
    ROUND(100.0 * hn.allocated_memory / hn.capacity_memory, 2) as memory_utilization_pct,
    ROUND(100.0 * hn.allocated_gpu_slots / hn.capacity_gpu_slots, 2) as gpu_utilization_pct,
    hn.allocated_vms,
    hn.capacity_cpu - hn.allocated_cpu as available_cpu,
    hn.capacity_memory - hn.allocated_memory as available_memory
FROM host_nodes hn;

-- View: GPU Allocation Status
CREATE VIEW gpu_allocation_status AS
SELECT
    g.id,
    g.host_node_id,
    hn.name as host_name,
    g.model,
    g.status,
    CASE WHEN ga.id IS NOT NULL THEN TRUE ELSE FALSE END as allocated,
    ga.vm_id,
    v.name as vm_name,
    ga.allocated_at
FROM gpus g
LEFT JOIN gpu_allocations ga ON g.id = ga.gpu_id
LEFT JOIN vms v ON ga.vm_id = v.id
LEFT JOIN host_nodes hn ON g.host_node_id = hn.id;

-- Create stored procedures for common operations

-- Procedure: Allocate VM to Host
CREATE OR REPLACE FUNCTION allocate_vm_to_host(
    vm_id UUID,
    host_id UUID
) RETURNS BOOLEAN AS $$
BEGIN
    UPDATE vms SET host_node_id = host_id, updated_at = NOW() WHERE id = vm_id;
    RETURN TRUE;
EXCEPTION WHEN OTHERS THEN
    RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

-- Procedure: Release VM from Host
CREATE OR REPLACE FUNCTION release_vm_from_host(
    vm_id UUID
) RETURNS BOOLEAN AS $$
BEGIN
    UPDATE vms SET host_node_id = NULL, updated_at = NOW() WHERE id = vm_id;
    DELETE FROM gpu_allocations WHERE vm_id = vm_id;
    RETURN TRUE;
EXCEPTION WHEN OTHERS THEN
    RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

-- Set up connection pooling recommendations
COMMENT ON TABLE vms IS 'Virtual machine definitions and state';
COMMENT ON TABLE host_nodes IS 'Physical host nodes in the cluster';
COMMENT ON TABLE gpus IS 'GPU devices available in the cluster';
COMMENT ON TABLE tasks IS 'Async tasks for provisioning and operations';
COMMENT ON TABLE audit_logs IS 'Audit trail for compliance';

-- Create default configurations
INSERT INTO configurations (config_type, config_key, config_value) VALUES
    ('scheduler', 'algorithm', '"bin-packing"'),
    ('scheduler', 'max_vms_per_node', '100'),
    ('gpu', 'allocation_policy', '"bin-packing"'),
    ('gpu', 'allow_sharing', 'false'),
    ('gpu', 'max_vms_per_gpu', '1')
ON CONFLICT (config_type, config_key, version) DO NOTHING;

COMMIT;

-- Indexes for better performance
CREATE INDEX idx_vms_updated_at ON vms(updated_at DESC);
CREATE INDEX idx_tasks_status_created ON tasks(status, created_at DESC);
CREATE INDEX idx_audit_logs_resource_action ON audit_logs(resource, action, timestamp DESC);
CREATE INDEX idx_resource_metrics_composite ON resource_metrics(vm_id, timestamp DESC);

-- Set up auto-cleanup for old metrics (optional - comment out if not needed)
-- CREATE OR REPLACE FUNCTION cleanup_old_metrics()
-- RETURNS void AS $$
-- BEGIN
--   DELETE FROM resource_metrics WHERE timestamp < NOW() - INTERVAL '30 days';
-- END;
-- $$ LANGUAGE plpgsql;

-- This script should be run after creating the database
-- Example: psql -U aihypervisor -d aihypervisor -f init-db.sql

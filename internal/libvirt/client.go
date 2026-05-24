//go:build libvirt
// +build libvirt

package libvirt

import (
	"context"
	"fmt"
	"sync"

	lv "github.com/libvirt/libvirt-go-module"
	"github.com/sirupsen/logrus"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"
)

// Client provides a concurrency-safe wrapper around libvirt connections.
type Client struct {
	conn   *lv.Connect
	mu     sync.Mutex
	logger *logrus.Logger
}

// NewClient establishes a new libvirt connection.
func NewClient(cfg config.LibvirtConfig, logger *logrus.Logger) (*Client, error) {
	conn, err := lv.NewConnect(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("connect libvirt: %w", err)
	}

	return &Client{
		conn:   conn,
		logger: logger,
	}, nil
}

// Close closes the underlying libvirt connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// CreateDomain creates and starts a libvirt domain from XML.
func (c *Client) CreateDomain(ctx context.Context, domainXML string) error {
	return c.withConn(ctx, func(conn *lv.Connect) error {
		domain, err := conn.DomainCreateXML(domainXML, 0)
		if domain != nil {
			defer domain.Free()
		}
		return err
	})
}

// DefineDomain defines a libvirt domain without starting it.
func (c *Client) DefineDomain(ctx context.Context, domainXML string) error {
	return c.withConn(ctx, func(conn *lv.Connect) error {
		domain, err := conn.DomainDefineXML(domainXML)
		if domain != nil {
			defer domain.Free()
		}
		return err
	})
}

// DestroyDomain forcefully stops a libvirt domain.
func (c *Client) DestroyDomain(ctx context.Context, domainName string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		return domain.Destroy()
	})
}

// StartDomain starts a defined domain.
func (c *Client) StartDomain(ctx context.Context, domainName string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		return domain.Create()
	})
}

// StopDomain attempts a graceful shutdown of a domain.
func (c *Client) StopDomain(ctx context.Context, domainName string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		return domain.Shutdown()
	})
}

// SuspendDomain pauses a running domain.
func (c *Client) SuspendDomain(ctx context.Context, domainName string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		return domain.Suspend()
	})
}

// ResumeDomain resumes a suspended domain.
func (c *Client) ResumeDomain(ctx context.Context, domainName string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		return domain.Resume()
	})
}

// RebootDomain reboots a running domain.
func (c *Client) RebootDomain(ctx context.Context, domainName string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		return domain.Reboot(0)
	})
}

// GetDomainInfo returns current domain information.
func (c *Client) GetDomainInfo(ctx context.Context, domainName string) (*orchestrator.LibvirtDomainInfo, error) {
	var info *orchestrator.LibvirtDomainInfo
	err := c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		state, _, err := domain.GetState()
		if err != nil {
			return err
		}

		lvInfo, err := domain.GetInfo()
		if err != nil {
			return err
		}

		uuid, err := domain.GetUUIDString()
		if err != nil {
			return err
		}

		info = &orchestrator.LibvirtDomainInfo{
			Name:      domainName,
			UUID:      uuid,
			MaxMemory: lvInfo.MaxMem,
			Memory:    lvInfo.Memory,
			MaxCPU:    int(lvInfo.NrVirtCpu),
			CPU:       int(lvInfo.NrVirtCpu),
			State:     domainStateString(state),
			CPUTime:   lvInfo.CpuTime,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return info, nil
}

// ListDomains returns all domain names.
func (c *Client) ListDomains(ctx context.Context) ([]string, error) {
	var names []string
	err := c.withConn(ctx, func(conn *lv.Connect) error {
		domains, err := conn.ListAllDomains(0)
		if err != nil {
			return err
		}

		for _, domain := range domains {
			name, err := domain.GetName()
			_ = domain.Free()
			if err != nil {
				continue
			}
			names = append(names, name)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return names, nil
}

// AttachDevice attaches a device to a domain.
func (c *Client) AttachDevice(ctx context.Context, domainName string, deviceXML string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		flags := lv.DOMAIN_DEVICE_MODIFY_LIVE | lv.DOMAIN_DEVICE_MODIFY_CONFIG
		return domain.AttachDeviceFlags(deviceXML, flags)
	})
}

// DetachDevice detaches a device from a domain.
func (c *Client) DetachDevice(ctx context.Context, domainName string, deviceXML string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		flags := lv.DOMAIN_DEVICE_MODIFY_LIVE | lv.DOMAIN_DEVICE_MODIFY_CONFIG
		return domain.DetachDeviceFlags(deviceXML, flags)
	})
}

// UndefineDomain removes a libvirt domain definition.
func (c *Client) UndefineDomain(ctx context.Context, domainName string) error {
	return c.withDomain(ctx, domainName, func(domain *lv.Domain) error {
		return domain.Undefine()
	})
}

func (c *Client) withConn(ctx context.Context, fn func(*lv.Connect) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}

	return fn(c.conn)
}

func (c *Client) withDomain(ctx context.Context, domainName string, fn func(*lv.Domain) error) error {
	return c.withConn(ctx, func(conn *lv.Connect) error {
		domain, err := conn.LookupDomainByName(domainName)
		if domain != nil {
			defer domain.Free()
		}
		if err != nil {
			return err
		}
		return fn(domain)
	})
}

func domainStateString(state lv.DomainState) string {
	switch state {
	case lv.DOMAIN_RUNNING:
		return "running"
	case lv.DOMAIN_BLOCKED:
		return "blocked"
	case lv.DOMAIN_PAUSED:
		return "paused"
	case lv.DOMAIN_SHUTDOWN:
		return "shutdown"
	case lv.DOMAIN_SHUTOFF:
		return "shutoff"
	case lv.DOMAIN_CRASHED:
		return "crashed"
	case lv.DOMAIN_PMSUSPENDED:
		return "pmsuspended"
	default:
		return "unknown"
	}
}


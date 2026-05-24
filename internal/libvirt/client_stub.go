//go:build !libvirt
// +build !libvirt

package libvirt

import (
    "context"
    "fmt"

    "github.com/sirupsen/logrus"

    "github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
)

// Client is a stub used when libvirt support is not compiled in.
type Client struct {
}

// NewClient returns an error indicating libvirt support is disabled in this build.
func NewClient(cfg config.LibvirtConfig, logger *logrus.Logger) (*Client, error) {
    return nil, fmt.Errorf("libvirt support not enabled in this build; compile with -tags libvirt to enable")
}

func (c *Client) Close() error {
    return nil
}

func (c *Client) CreateDomain(ctx context.Context, domainXML string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) DefineDomain(ctx context.Context, domainXML string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) DestroyDomain(ctx context.Context, domainName string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) StartDomain(ctx context.Context, domainName string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) StopDomain(ctx context.Context, domainName string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) SuspendDomain(ctx context.Context, domainName string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) ResumeDomain(ctx context.Context, domainName string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) RebootDomain(ctx context.Context, domainName string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) GetDomainInfo(ctx context.Context, domainName string) (interface{}, error) {
    return nil, fmt.Errorf("libvirt not available in this build")
}

func (c *Client) ListDomains(ctx context.Context) ([]string, error) {
    return nil, fmt.Errorf("libvirt not available in this build")
}

func (c *Client) AttachDevice(ctx context.Context, domainName string, deviceXML string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) DetachDevice(ctx context.Context, domainName string, deviceXML string) error {
    return fmt.Errorf("libvirt not available in this build")
}

func (c *Client) UndefineDomain(ctx context.Context, domainName string) error {
    return fmt.Errorf("libvirt not available in this build")
}

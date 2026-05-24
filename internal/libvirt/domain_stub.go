//go:build !libvirt
// +build !libvirt

package libvirt

import (
    "fmt"
    "time"

    "github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

// BuildDomainXML returns a minimal domain XML suitable for CI/testing when
// libvirt XML helpers are not available. This should be replaced by the
// libvirt-backed implementation when building with -tags libvirt.
func BuildDomainXML(vm *models.VirtualMachine) (string, error) {
    if vm == nil {
        return "", fmt.Errorf("vm is required")
    }
    if vm.Name == "" {
        return "", fmt.Errorf("vm name is required")
    }
    memoryMiB := vm.Flavor.Memory * 1024
    vcpu := vm.Flavor.CPU

    now := time.Now().Format(time.RFC3339)
    // Minimal domain XML - not feature-complete.
    xml := fmt.Sprintf(`<domain type="kvm">
  <name>%s</name>
  <uuid>%s</uuid>
  <memory unit="MiB">%d</memory>
  <vcpu>%d</vcpu>
  <os>
    <type>hvm</type>
  </os>
  <devices>
    <emulator>/usr/bin/qemu-system-x86_64</emulator>
  </devices>
  <metadata>
    <generated_at>%s</generated_at>
  </metadata>
</domain>`, vm.Name, vm.ID, memoryMiB, vcpu, now)
    return xml, nil
}

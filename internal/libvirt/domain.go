//go:build libvirt
// +build libvirt

package libvirt

import (
	"fmt"
	"strings"

	libvirtxml "github.com/libvirt/libvirt-go-xml"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

// BuildDomainXML builds a libvirt domain XML definition from a VM spec.
func BuildDomainXML(vm *models.VirtualMachine) (string, error) {
	if vm == nil {
		return "", fmt.Errorf("vm is required")
	}
	if vm.Name == "" {
		return "", fmt.Errorf("vm name is required")
	}
	if vm.Flavor.CPU <= 0 || vm.Flavor.Memory <= 0 {
		return "", fmt.Errorf("vm flavor must include cpu and memory")
	}
	if len(vm.StorageConfig.Volumes) == 0 {
		return "", fmt.Errorf("vm storage volume is required")
	}

	memoryMiB := uint(vm.Flavor.Memory * 1024)
	vcpuCount := uint(vm.Flavor.CPU)
	volume := vm.StorageConfig.Volumes[0]
	if volume.TargetDev == "" {
		volume.TargetDev = "vda"
	}
	if volume.Source == "" {
		return "", fmt.Errorf("volume source is required")
	}

	domain := libvirtxml.Domain{
		Type: "kvm",
		Name: vm.Name,
		UUID: vm.ID,
		Memory: &libvirtxml.DomainMemory{
			Value: memoryMiB,
			Unit:  "MiB",
		},
		VCPU: &libvirtxml.DomainVCPU{Value: vcpuCount},
		CPU: &libvirtxml.DomainCPU{
			Mode: "host-passthrough",
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{Type: "hvm"},
		},
	}

	disk := libvirtxml.DomainDisk{
		Device: "disk",
		Driver: &libvirtxml.DomainDiskDriver{
			Name: "qemu",
			Type: volume.Format,
		},
		Source: &libvirtxml.DomainDiskSource{
			File: &libvirtxml.DomainDiskSourceFile{File: volume.Source},
		},
		Target: &libvirtxml.DomainDiskTarget{Dev: volume.TargetDev, Bus: "virtio"},
	}
	if disk.Driver.Type == "" {
		disk.Driver.Type = "qcow2"
	}

	interfaces := make([]libvirtxml.DomainInterface, 0, len(vm.NetworkConfig.Interfaces))
	for _, nic := range vm.NetworkConfig.Interfaces {
		iface := libvirtxml.DomainInterface{Model: &libvirtxml.DomainInterfaceModel{Type: "virtio"}}
		if nic.MAC != "" {
			iface.MAC = &libvirtxml.DomainInterfaceMAC{Address: nic.MAC}
		}
		if nic.Bridge != "" || strings.EqualFold(vm.NetworkConfig.Type, "bridge") {
			iface.Type = "bridge"
			iface.Source = &libvirtxml.DomainInterfaceSource{Bridge: nic.Bridge}
		} else {
			iface.Type = "network"
			iface.Source = &libvirtxml.DomainInterfaceSource{Network: "default"}
		}
		interfaces = append(interfaces, iface)
	}

	if len(interfaces) == 0 {
		interfaces = append(interfaces, libvirtxml.DomainInterface{
			Type:   "network",
			Source: &libvirtxml.DomainInterfaceSource{Network: "default"},
			Model:  &libvirtxml.DomainInterfaceModel{Type: "virtio"},
		})
	}

	hostdevs := make([]libvirtxml.DomainHostdev, 0, len(vm.GPUAllocations))
	for _, allocation := range vm.GPUAllocations {
		if allocation.PCIAddress == "" {
			continue
		}
		domainID, bus, slot, function, err := parsePCIAddress(allocation.PCIAddress)
		if err != nil {
			continue
		}
		hostdevs = append(hostdevs, libvirtxml.DomainHostdev{
			Mode:    "subsystem",
			Type:    "pci",
			Managed: "yes",
			Source: &libvirtxml.DomainHostdevSource{
				Address: &libvirtxml.DomainAddressPCI{
					Domain:   &domainID,
					Bus:      &bus,
					Slot:     &slot,
					Function: &function,
				},
			},
		})
	}

	domain.Devices = &libvirtxml.DomainDeviceList{
		Disks:      []libvirtxml.DomainDisk{disk},
		Interfaces: interfaces,
		Hostdevs:   hostdevs,
		Graphics: []libvirtxml.DomainGraphic{
			{Type: "vnc", AutoPort: "yes"},
		},
	}

	return domain.Marshal()
}

func parsePCIAddress(addr string) (string, string, string, string, error) {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return "", "", "", "", fmt.Errorf("invalid pci address")
	}

	domain := "0x0000"
	bus := ""
	slotFunc := ""

	if len(parts) == 3 {
		domain = toPCIHex(parts[0])
		bus = toPCIHex(parts[1])
		slotFunc = parts[2]
	} else {
		bus = toPCIHex(parts[0])
		slotFunc = parts[1]
	}

	slotParts := strings.Split(slotFunc, ".")
	if len(slotParts) != 2 {
		return "", "", "", "", fmt.Errorf("invalid pci address")
	}

	slot := toPCIHex(slotParts[0])
	function := toPCIHex(slotParts[1])

	return domain, bus, slot, function, nil
}

func toPCIHex(value string) string {
	value = strings.TrimPrefix(value, "0x")
	if len(value) == 1 {
		value = "0" + value
	}
	return "0x" + value
}

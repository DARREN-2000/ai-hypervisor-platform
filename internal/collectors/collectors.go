package collectors

import (
    "context"
    "math/rand"
    "time"

    "github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
    "github.com/DARREN-2000/ai-hypervisor-platform/pkg/telemetry"
)

// VMUsage is a compact structure returned by a VM fetcher.
type VMUsage struct {
    VMID   string
    HostID string
    Metrics models.ResourceMetrics
}

// GPUUsage is a compact structure returned by a GPU fetcher.
type GPUUsage struct {
    GPUID   string
    HostID  string
    Model   string
    Metrics models.GPUMetrics
}

// VMFetcher is the signature for a function that samples VM resource usage.
type VMFetcher func(ctx context.Context) ([]VMUsage, error)

// GPUFetcher is the signature for a function that samples GPU usage.
type GPUFetcher func(ctx context.Context) ([]GPUUsage, error)

// StartVMCollector runs a background loop that periodically invokes fetcher
// and reports metrics to the provided telemetry.Metrics implementation.
// It returns a cancel function to stop the collector.
func StartVMCollector(ctx context.Context, interval time.Duration, metrics telemetry.Metrics, fetcher VMFetcher) context.CancelFunc {
    if fetcher == nil {
        // provide a synthetic no-op fetcher to avoid nil panics
        fetcher = func(ctx context.Context) ([]VMUsage, error) { return []VMUsage{}, nil }
    }

    cctx, cancel := context.WithCancel(ctx)
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-cctx.Done():
                return
            case <-ticker.C:
                items, err := fetcher(cctx)
                if err != nil {
                    continue
                }
                for _, u := range items {
                    // Map our model fields into telemetry observer
                    metrics.ObserveVMResource(u.VMID, u.HostID, u.Metrics.CPU, int64(u.Metrics.Memory)*1024*1024, int64(u.Metrics.DiskIORead), int64(u.Metrics.DiskIOWrite), int64(u.Metrics.NetworkIn), int64(u.Metrics.NetworkOut))
                }
            }
        }
    }()
    return cancel
}

// StartGPUCollector runs a background loop that periodically invokes fetcher
// and reports GPU metrics via telemetry.Metrics.
func StartGPUCollector(ctx context.Context, interval time.Duration, metrics telemetry.Metrics, fetcher GPUFetcher) context.CancelFunc {
    if fetcher == nil {
        fetcher = func(ctx context.Context) ([]GPUUsage, error) { return []GPUUsage{}, nil }
    }

    cctx, cancel := context.WithCancel(ctx)
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-cctx.Done():
                return
            case <-ticker.C:
                items, err := fetcher(cctx)
                if err != nil {
                    continue
                }
                for _, u := range items {
                    metrics.ObserveGPUUsage(u.GPUID, u.HostID, u.Model, float64(u.Metrics.Utilization), int64(u.Metrics.MemoryUsed)*1024*1024, int64(u.Metrics.MemoryFree)*1024*1024, float64(u.Metrics.TemperatureC), float64(u.Metrics.PowerDraw))
                }
            }
        }
    }()
    return cancel
}

// SyntheticFetchers provides example demo fetchers used for local testing
// and initial integration. These should be replaced by real implementations
// that query libvirt, nvidia-smi, NVML, or node exporters in production.
var SyntheticVMFetcher = func(num int) VMFetcher {
    return func(ctx context.Context) ([]VMUsage, error) {
        out := make([]VMUsage, 0, num)
        for i := 0; i < num; i++ {
            vm := VMUsage{
                VMID:   "vm-" + randSeq(6),
                HostID: "host-1",
                Metrics: models.ResourceMetrics{
                    CPU:      rand.Float64() * 8,
                    Memory:   1024 + rand.Intn(16*1024),
                    DiskIORead:  rand.Intn(1024),
                    DiskIOWrite: rand.Intn(1024),
                    NetworkIn: rand.Intn(1000),
                    NetworkOut: rand.Intn(1000),
                    Timestamp: time.Now(),
                },
            }
            out = append(out, vm)
        }
        return out, nil
    }
}

var SyntheticGPUFetcher = func(num int) GPUFetcher {
    return func(ctx context.Context) ([]GPUUsage, error) {
        out := make([]GPUUsage, 0, num)
        for i := 0; i < num; i++ {
            g := GPUUsage{
                GPUID:  "gpu-" + randSeq(4),
                HostID: "host-1",
                Model:  "nvidia-a100",
                Metrics: models.GPUMetrics{
                    Utilization:  rand.Intn(100),
                    MemoryUsed:   1024 + rand.Intn(48*1024),
                    MemoryFree:   1024 + rand.Intn(48*1024),
                    TemperatureC: 40 + rand.Intn(50),
                    PowerDraw:    50 + rand.Intn(250),
                },
            }
            out = append(out, g)
        }
        return out, nil
    }
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

import { motion } from "framer-motion"
import { Layers, Activity, RefreshCw, Zap, Shield, BarChart, GitCommit, Code } from "lucide-react"
import { cn } from "@/utils/cn"

const features = [
  {
    title: "Unified Control Plane",
    description: "Manage AWS, GCP, and on-prem clusters from one dashboard. Eliminates vendor lock-in.",
    icon: <Layers size={20} />,
    className: "md:col-span-2 lg:col-span-3",
  },
  {
    title: "Zero-Config Autoscaling",
    description: "Scale from 0 to 10,000 GPUs instantly based on queue depth.",
    icon: <Activity size={20} />,
    className: "md:col-span-1 lg:col-span-2",
  },
  {
    title: "Fault-Tolerant Training",
    description: "Automatic checkpointing and job resumption on node failure.",
    icon: <RefreshCw size={20} />,
    className: "md:col-span-1 lg:col-span-3",
  },
  {
    title: "Bare-Metal Networking",
    description: "eBPF-powered networking for maximum throughput and zero-copy memory.",
    icon: <Zap size={20} />,
    className: "md:col-span-2 lg:col-span-4",
  },
  {
    title: "Granular RBAC & Identity",
    description: "Native OIDC integration and namespace isolation.",
    icon: <Shield size={20} />,
    className: "md:col-span-1 lg:col-span-2",
  },
  {
    title: "Real-Time Telemetry",
    description: "Sub-second visibility into GPU utilization, memory, and temps.",
    icon: <BarChart size={20} />,
    className: "md:col-span-1 lg:col-span-2",
  },
  {
    title: "Instant Reproducibility",
    description: "Declarative YAML configs mapped to immutable image hashes.",
    icon: <GitCommit size={20} />,
    className: "md:col-span-1 lg:col-span-2",
  },
  {
    title: "Developer-First APIs",
    description: "gRPC and REST APIs with native Go, Python, and Rust SDKs.",
    icon: <Code size={20} />,
    className: "md:col-span-2 lg:col-span-2",
  },
]

export default function Features() {
  return (
    <section id="features" className="py-32 bg-bg-secondary relative">
      <div className="container mx-auto px-6 max-w-container">

        <div className="mb-16">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-4">Primitives for the AI era.</h2>
          <p className="text-text-muted text-lg max-w-2xl">
            Everything you need to orchestrate state-of-the-art models, built into a single cohesive layer.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-8 gap-4 auto-rows-[200px]">
          {features.map((feature, idx) => (
            <motion.div
              key={idx}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: idx * 0.05 }}
              className={cn(
                "glass p-6 rounded-2xl flex flex-col justify-between group cursor-default hover:border-white/20 transition-all duration-300 relative overflow-hidden",
                feature.className
              )}
            >
              <div className="absolute inset-0 bg-gradient-to-br from-white/[0.02] to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
              <div className="w-10 h-10 rounded-full bg-white/5 flex items-center justify-center text-white border border-white/10 group-hover:scale-110 transition-transform">
                {feature.icon}
              </div>
              <div className="relative z-10">
                <h3 className="font-semibold text-lg text-white mb-2">{feature.title}</h3>
                <p className="text-sm text-text-muted line-clamp-2 group-hover:line-clamp-none transition-all">{feature.description}</p>
              </div>
            </motion.div>
          ))}
        </div>

      </div>
    </section>
  )
}

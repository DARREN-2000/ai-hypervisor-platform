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
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-4 text-white">Primitives for the AI era.</h2>
          <p className="text-text-muted text-lg max-w-2xl font-medium">
            Everything you need to orchestrate state-of-the-art models, built into a single cohesive layer.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-8 gap-4 auto-rows-[220px]">
          {features.map((feature, idx) => (
            <motion.div
              key={idx}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: idx * 0.05 }}
              className={cn(
                "glass p-6 rounded-2xl flex flex-col justify-between group cursor-default hover:border-white/[0.12] transition-all duration-500 relative overflow-hidden hover:-translate-y-1",
                feature.className
              )}
            >
              <div className="absolute inset-0 bg-gradient-to-br from-white/[0.03] to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />
              <div className="w-9 h-9 rounded-md bg-white/[0.03] flex items-center justify-center text-zinc-300 border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.02)] group-hover:text-white transition-colors duration-300">
                {feature.icon}
              </div>
              <div className="relative z-10">
                <h3 className="font-medium text-base text-zinc-100 mb-2">{feature.title}</h3>
                <p className="text-sm text-zinc-400 line-clamp-2 group-hover:line-clamp-none transition-all duration-300">{feature.description}</p>
              </div>
            </motion.div>
          ))}
        </div>

      </div>
    </section>
  )
}

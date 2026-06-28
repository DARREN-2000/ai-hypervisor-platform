import { motion } from "framer-motion"
import { Layers, Zap, Shield } from "lucide-react"

export default function Solution() {
  const pillars = [
    {
      icon: <Layers size={22} className="text-zinc-300" />,
      title: "Unified Resource Pool",
      description: "Onyx abstracts all your disparate GPUs—across clouds and on-prem—into a single, logical cluster."
    },
    {
      icon: <Zap size={22} className="text-zinc-300" />,
      title: "Zero-Config Autoscaling",
      description: "Scale from 0 to 10,000 GPUs instantly based on queue depth. Pay exactly for what you compute."
    },
    {
      icon: <Shield size={22} className="text-zinc-300" />,
      title: "Fault-Tolerant by Default",
      description: "Automatic checkpointing and job resumption. If a node dies, Onyx simply reschedules and resumes."
    }
  ]

  return (
    <section className="py-32 relative">
      <div className="container mx-auto px-6 max-w-container">

        <div className="text-center max-w-3xl mx-auto mb-20">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-6 text-white">
            One OS for all your compute.
          </h2>
          <p className="text-lg text-text-muted leading-relaxed font-medium">
            Onyx abstracts away the chaos of cluster management, networking, and workload orchestration into a single, elegant layer.
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          {pillars.map((pillar, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
              className="glass p-8 rounded-2xl relative group overflow-hidden transition-transform duration-500 hover:-translate-y-1"
            >
              <div className="absolute inset-0 bg-gradient-to-b from-white/[0.03] to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />
              <div className="w-10 h-10 rounded-md bg-white/[0.03] flex items-center justify-center mb-6 border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.02)] group-hover:bg-white/[0.06] transition-colors duration-300">
                {pillar.icon}
              </div>
              <h3 className="text-lg font-semibold mb-3 text-white">{pillar.title}</h3>
              <p className="text-text-muted text-sm leading-relaxed">{pillar.description}</p>
            </motion.div>
          ))}
        </div>

      </div>
    </section>
  )
}

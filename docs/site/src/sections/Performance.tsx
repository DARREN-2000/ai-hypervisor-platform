import { motion } from "framer-motion"

export default function Performance() {
  return (
    <section className="py-32 relative border-t border-white/5">
      <div className="container mx-auto px-6 max-w-container">
        <div className="grid md:grid-cols-2 gap-16 items-center">

          <div>
            <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-6 text-white">
              Speed is a feature.
            </h2>
            <p className="text-lg text-text-muted mb-8 leading-relaxed font-medium">
              Standard orchestrators like Kubernetes weren't built for AI. They poll state and rely on slow iptables. Onyx is event-driven and uses eBPF, resulting in microsecond scheduling latency.
            </p>
            <div className="flex gap-10">
              <div>
                <div className="text-3xl font-bold text-white mb-1 tracking-tight">{"<1ms"}</div>
                <div className="text-sm text-zinc-400 font-medium">Scheduling Latency</div>
              </div>
              <div>
                <div className="text-3xl font-bold text-white mb-1 tracking-tight">99.9%</div>
                <div className="text-sm text-zinc-400 font-medium">Network Throughput</div>
              </div>
            </div>
          </div>

          <div className="glass p-8 rounded-2xl relative overflow-hidden">
            <div className="absolute top-0 right-0 p-4 opacity-10">
              <div className="w-24 h-24 border border-white rounded-full"></div>
            </div>
            <h3 className="text-xs font-semibold text-zinc-500 mb-8 uppercase tracking-widest">Job Scheduling Time (ms)</h3>

            <div className="space-y-8">
              <div>
                <div className="flex justify-between text-sm mb-3">
                  <span className="text-zinc-100 font-medium flex items-center gap-2"><span className="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></span>Onyx</span>
                  <span className="font-mono text-emerald-400">0.8ms</span>
                </div>
                <div className="w-full bg-white/[0.04] rounded-full h-2 overflow-hidden border border-white/[0.02]">
                  <motion.div
                    initial={{ width: 0 }}
                    whileInView={{ width: "5%" }}
                    transition={{ duration: 1.5, ease: [0.16, 1, 0.3, 1] }}
                    className="bg-emerald-500 h-2 rounded-full shadow-[0_0_10px_rgba(16,185,129,0.5)]"
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between text-sm mb-3">
                  <span className="text-zinc-400 font-medium">Standard K8s Scheduler</span>
                  <span className="font-mono text-zinc-500">45.0ms</span>
                </div>
                <div className="w-full bg-white/[0.04] rounded-full h-2 overflow-hidden border border-white/[0.02]">
                  <motion.div
                    initial={{ width: 0 }}
                    whileInView={{ width: "65%" }}
                    transition={{ duration: 1.5, ease: [0.16, 1, 0.3, 1] }}
                    className="bg-white/20 h-2 rounded-full"
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between text-sm mb-3">
                  <span className="text-zinc-400 font-medium">Slurm</span>
                  <span className="font-mono text-zinc-500">120.0ms</span>
                </div>
                <div className="w-full bg-white/[0.04] rounded-full h-2 overflow-hidden border border-white/[0.02]">
                  <motion.div
                    initial={{ width: 0 }}
                    whileInView={{ width: "95%" }}
                    transition={{ duration: 1.5, ease: [0.16, 1, 0.3, 1] }}
                    className="bg-white/10 h-2 rounded-full"
                  />
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  )
}

import { motion } from "framer-motion"

export default function Performance() {
  return (
    <section className="py-32 relative border-t border-white/5">
      <div className="container mx-auto px-6 max-w-container">
        <div className="grid md:grid-cols-2 gap-16 items-center">

          <div>
            <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-6">
              Speed is a feature.
            </h2>
            <p className="text-lg text-text-muted mb-8 leading-relaxed">
              Standard orchestrators like Kubernetes weren't built for AI. They poll state and rely on slow iptables. Onyx is event-driven and uses eBPF, resulting in microsecond scheduling latency.
            </p>
            <div className="flex gap-8">
              <div>
                <div className="text-3xl font-bold text-white mb-1">{"<1ms"}</div>
                <div className="text-sm text-text-muted">Scheduling Latency</div>
              </div>
              <div>
                <div className="text-3xl font-bold text-white mb-1">99.9%</div>
                <div className="text-sm text-text-muted">Network Throughput</div>
              </div>
            </div>
          </div>

          <div className="glass p-8 rounded-2xl">
            <h3 className="text-sm font-medium text-text-muted mb-6 uppercase tracking-wider">Job Scheduling Time (ms)</h3>

            <div className="space-y-6">
              <div>
                <div className="flex justify-between text-sm mb-2">
                  <span className="text-white">Onyx</span>
                  <span className="font-mono text-green-400">0.8ms</span>
                </div>
                <div className="w-full bg-white/5 rounded-full h-3">
                  <motion.div
                    initial={{ width: 0 }}
                    whileInView={{ width: "5%" }}
                    transition={{ duration: 1, ease: "easeOut" }}
                    className="bg-green-500 h-3 rounded-full"
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between text-sm mb-2">
                  <span className="text-text-muted">Standard K8s Scheduler</span>
                  <span className="font-mono text-text-muted">45.0ms</span>
                </div>
                <div className="w-full bg-white/5 rounded-full h-3">
                  <motion.div
                    initial={{ width: 0 }}
                    whileInView={{ width: "65%" }}
                    transition={{ duration: 1, ease: "easeOut" }}
                    className="bg-white/20 h-3 rounded-full"
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between text-sm mb-2">
                  <span className="text-text-muted">Slurm</span>
                  <span className="font-mono text-text-muted">120.0ms</span>
                </div>
                <div className="w-full bg-white/5 rounded-full h-3">
                  <motion.div
                    initial={{ width: 0 }}
                    whileInView={{ width: "95%" }}
                    transition={{ duration: 1, ease: "easeOut" }}
                    className="bg-white/10 h-3 rounded-full"
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

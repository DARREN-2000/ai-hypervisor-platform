import { motion } from "framer-motion"

export default function Architecture() {
  return (
    <section className="py-32 relative overflow-hidden">
      <div className="container mx-auto px-6 max-w-container">

        <div className="text-center max-w-3xl mx-auto mb-20">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-6">
            Bare-metal performance. <br className="hidden md:block" />
            <span className="text-text-muted">Cloud-native flexibility.</span>
          </h2>
          <p className="text-lg text-text-muted leading-relaxed">
            Onyx bypasses standard Linux networking stacks using custom eBPF filters, delivering near-bare-metal latency across distributed clusters.
          </p>
        </div>

        <div className="relative mx-auto max-w-4xl rounded-2xl border border-white/10 bg-black/50 p-8 backdrop-blur-xl">
          <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,rgba(0,112,243,0.15),transparent_50%)]" />

          <div className="relative z-10 flex flex-col items-center gap-6">
            {/* Layer 1: App */}
            <motion.div
              whileHover={{ scale: 1.02 }}
              className="w-full md:w-2/3 glass p-4 text-center font-mono text-sm border-blue-500/30 text-blue-200"
            >
              Your Code (PyTorch / Ray / JAX)
            </motion.div>

            <div className="w-px h-8 bg-gradient-to-b from-blue-500/50 to-purple-500/50" />

            {/* Layer 2: Onyx */}
            <motion.div
              whileHover={{ scale: 1.02 }}
              className="w-full glass p-6 text-center shadow-[0_0_30px_rgba(121,40,202,0.2)] border-purple-500/50 relative overflow-hidden"
            >
              <div className="absolute inset-0 animate-flow opacity-20 bg-[linear-gradient(90deg,transparent_0%,rgba(255,255,255,0.2)_50%,transparent_100%)] bg-[length:200%_100%]" />
              <span className="text-xl font-bold text-white relative z-10">Onyx Control Plane & eBPF Router</span>
            </motion.div>

            <div className="flex gap-16 md:gap-32 h-8">
              <div className="w-px h-full bg-gradient-to-b from-purple-500/50 to-green-500/50" />
              <div className="w-px h-full bg-gradient-to-b from-purple-500/50 to-green-500/50 hidden md:block" />
              <div className="w-px h-full bg-gradient-to-b from-purple-500/50 to-green-500/50" />
            </div>

            {/* Layer 3: Hardware */}
            <div className="flex flex-wrap justify-center gap-4 w-full">
              {['AWS H100', 'GCP A100', 'On-Prem HGX'].map((node, i) => (
                <motion.div
                  key={i}
                  whileHover={{ scale: 1.05 }}
                  className="flex-1 min-w-[120px] glass p-4 text-center font-mono text-xs border-green-500/30 text-green-200"
                >
                  {node}
                </motion.div>
              ))}
            </div>
          </div>

        </div>
      </div>
    </section>
  )
}

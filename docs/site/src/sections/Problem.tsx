import { motion } from "framer-motion"

export default function Problem() {
  return (
    <section className="py-32 relative bg-bg-secondary">
      <div className="container mx-auto px-6 max-w-container">
        <div className="grid md:grid-cols-2 gap-16 items-center">

          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
          >
            <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-6">
              Hardware is hard. <br />
              <span className="text-text-muted">Ops shouldn't be.</span>
            </h2>
            <p className="text-lg text-text-muted mb-8 leading-relaxed">
              Building AI infrastructure today feels like assembling an airplane in freefall. You're stitching together disjointed hardware, custom bash scripts, and multiple cloud providers just to get a training job running.
            </p>
            <ul className="space-y-4">
              {[
                "Wasted GPU cycles due to inefficient scheduling",
                "Vendor lock-in across AWS, GCP, and CoreWeave",
                "OOM errors crashing 14-day training runs",
                "Days lost configuring CUDA, NCCL, and networking"
              ].map((item, i) => (
                <li key={i} className="flex items-start gap-3 text-text-main">
                  <span className="text-red-500 mt-1">✗</span>
                  {item}
                </li>
              ))}
            </ul>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            className="rounded-xl border border-red-500/20 bg-black p-6 font-mono text-sm text-red-400 overflow-hidden shadow-[0_0_40px_rgba(239,68,68,0.1)] relative"
          >
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-red-500 to-transparent"></div>
            <p>[ERROR] RuntimeError: CUDA error: out of memory</p>
            <p className="opacity-70 mt-2">CUDA kernel errors might be asynchronously reported at some other API call, so the stacktrace below might be incorrect.</p>
            <p className="opacity-50 mt-2">For debugging consider passing CUDA_LAUNCH_BLOCKING=1.</p>
            <p className="mt-4">[FATAL] Node h100-worker-4 disconnected.</p>
            <p>[WARN] Attempting to reschedule pods... Failed.</p>
            <p className="mt-4 text-white">Process exited with code 1. Training aborted after 72 hours.</p>
          </motion.div>

        </div>
      </div>
    </section>
  )
}

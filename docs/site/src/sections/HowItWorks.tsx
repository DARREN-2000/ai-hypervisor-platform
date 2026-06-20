import { motion } from "framer-motion"

export default function HowItWorks() {
  const steps = [
    { num: "01", title: "Connect your clouds", desc: "Install the Onyx agent on any Kubernetes cluster or bare-metal Linux machine. Onyx automatically discovers GPUs and negotiates network topologies." },
    { num: "02", title: "Define your workload", desc: "Write a simple YAML file or use our Python SDK to describe your training job, required resources, and Docker image." },
    { num: "03", title: "Onyx schedules it", desc: "Our deterministic scheduler finds the optimal nodes, provisions them, sets up zero-copy networking, and starts your job in milliseconds." },
  ]

  return (
    <section className="py-32 bg-bg-secondary">
      <div className="container mx-auto px-6 max-w-container">

        <div className="mb-16">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-4">Deploy with a single command.</h2>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          {steps.map((step, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: i * 0.1 }}
              className="relative"
            >
              <div className="text-6xl font-black text-white/5 mb-4">{step.num}</div>
              <h3 className="text-xl font-semibold mb-3">{step.title}</h3>
              <p className="text-text-muted text-sm leading-relaxed">{step.desc}</p>
            </motion.div>
          ))}
        </div>

      </div>
    </section>
  )
}

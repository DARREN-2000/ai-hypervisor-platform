import { motion } from "framer-motion"
import { ChevronRight, Terminal } from "lucide-react"
import { Button } from "@/components/ui/Button"
import { Badge } from "@/components/ui/Badge"

export default function Hero() {
  return (
    <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden">
      <div className="container mx-auto px-6 max-w-container relative z-10 flex flex-col items-center text-center">

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: "easeOut" }}
          className="mb-8"
        >
          <Badge variant="outline" className="px-4 py-1.5 border-white/10 bg-white/5 backdrop-blur-md rounded-full text-sm flex items-center gap-2 hover:bg-white/10 transition-colors cursor-pointer">
            <span className="text-blue-400">✨ Onyx 1.0 is now available</span>
            <span className="text-text-muted">Introducing multi-cloud roaming</span>
            <ChevronRight size={14} className="text-text-muted" />
          </Badge>
        </motion.div>

        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1, ease: "easeOut" }}
          className="text-5xl md:text-7xl lg:text-8xl font-bold tracking-tighter max-w-4xl mb-8 leading-[1.1] bg-clip-text text-transparent bg-gradient-to-b from-white to-white/50"
        >
          Run AI at scale. <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-white via-zinc-400 to-zinc-600">Without the ops tax.</span>
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2, ease: "easeOut" }}
          className="text-lg md:text-xl text-text-muted max-w-2xl mb-12 leading-relaxed font-medium"
        >
          Unify your GPU clusters, orchestrate complex ML workloads, and scale from prototype to production with a single command. Stop managing hardware. Start building intelligence.
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3, ease: "easeOut" }}
          className="flex flex-col sm:flex-row items-center gap-4 mb-24"
        >
          <Button size="lg" className="w-full sm:w-auto text-base">
            Start Building
          </Button>
          <Button variant="outline" size="lg" className="w-full sm:w-auto text-base">
            Read Documentation
          </Button>
        </motion.div>

        {/* Mockup Window */}
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.5, ease: "easeOut" }}
          className="w-full max-w-5xl relative"
        >
          <div className="absolute -inset-1 bg-gradient-to-b from-white/10 to-transparent blur-[120px] -z-10 rounded-full" />
          <div className="rounded-xl overflow-hidden border border-white/[0.08] bg-[#0A0A0A] shadow-[0_30px_100px_-20px_rgba(0,0,0,1)]">
            {/* Window Header */}
            <div className="flex items-center px-4 py-3 border-b border-white/[0.04] bg-[#111111]">
              <div className="flex gap-2">
                <div className="w-2.5 h-2.5 rounded-full bg-[#FF5F56] border border-[#E0443E]" />
                <div className="w-2.5 h-2.5 rounded-full bg-[#FFBD2E] border border-[#DEA123]" />
                <div className="w-2.5 h-2.5 rounded-full bg-[#27C93F] border border-[#1AAB29]" />
              </div>
              <div className="mx-auto flex items-center gap-2 text-[11px] text-text-muted font-mono opacity-80">
                <Terminal size={12} />
                ~ / projects / onyx
              </div>
            </div>
            {/* Window Body */}
            <div className="p-6 text-left font-mono text-sm leading-relaxed text-zinc-300 bg-[#050505] h-[300px] overflow-hidden">
              <div className="flex items-center gap-3 mb-2">
                <span className="text-zinc-500">➜</span>
                <span className="text-zinc-100 font-medium">onyx deploy --cluster=h100-pool</span>
              </div>
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 1.5 }}
              >
                <div className="text-zinc-500 mt-4">[1/4] <span className="text-zinc-400">Authenticating with control plane...</span> <span className="text-zinc-300">Done</span></div>
                <div className="text-zinc-500">[2/4] <span className="text-zinc-400">Provisioning 64x NVIDIA H100 GPUs...</span> <span className="text-zinc-300">Done</span></div>
                <div className="text-zinc-500">[3/4] <span className="text-zinc-400">Establishing eBPF network routes...</span> <span className="text-zinc-300">Done</span></div>
                <div className="text-zinc-500">[4/4] <span className="text-zinc-400">Starting distributed training job...</span> <span className="text-zinc-300">Done</span></div>
                <div className="mt-4 text-zinc-300">Success. Job <span className="text-white font-medium">'llama-3-finetune'</span> is running in namespace <span className="text-white font-medium">'prod'</span>.</div>
                <div className="mt-2 text-zinc-500">Dashboard: <a href="#" className="text-zinc-300 hover:text-white underline decoration-zinc-700 underline-offset-4">console.onyx.dev/jobs/llama-3-finetune</a></div>
              </motion.div>
            </div>
          </div>
        </motion.div>

      </div>
    </section>
  )
}

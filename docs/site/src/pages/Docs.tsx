export default function Docs() {
  return (
    <div className="container mx-auto px-6 py-24 max-w-4xl min-h-screen">
      <h1 className="text-4xl md:text-5xl font-bold tracking-tighter mb-8 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/50">
        Documentation
      </h1>
      <p className="text-lg text-text-muted mb-8 leading-relaxed">
        Everything you need to know about the AI Hypervisor Platform.
      </p>
      <div className="p-8 border border-white/10 rounded-xl bg-white/5 backdrop-blur-md">
        <h2 className="text-2xl font-semibold mb-4 text-white">Getting Started</h2>
        <p className="text-zinc-400 mb-6 leading-relaxed">Learn how to deploy and configure your first cluster.</p>
        <ul className="list-disc pl-5 text-zinc-400 space-y-2">
            <li>Installation Guide</li>
            <li>Configuration Options</li>
            <li>API Reference</li>
        </ul>
      </div>
    </div>
  )
}

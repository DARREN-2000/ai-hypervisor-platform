export default function Roadmap() {
  const items = [
    { status: "Live", title: "Multi-cloud roaming", desc: "Deploy across AWS and GCP seamlessly." },
    { status: "Q3", title: "Serverless Inference", desc: "Scale down to zero for deployed models." },
    { status: "Q4", title: "Spot Instance Orchestration", desc: "Cut costs by 70% with automatic spot interruption handling." },
  ]
  return (
    <section className="py-24 bg-bg-secondary border-t border-white/5">
      <div className="container mx-auto px-6 max-w-container">
        <h2 className="text-3xl font-bold tracking-tight mb-12">Roadmap</h2>
        <div className="grid md:grid-cols-3 gap-8">
          {items.map((item, i) => (
            <div key={i} className="border-l border-white/10 pl-6 relative">
              <div className={`absolute left-[-5px] top-1 w-2.5 h-2.5 rounded-full ${item.status === 'Live' ? 'bg-green-500 shadow-[0_0_10px_#22c55e]' : 'bg-white/20'}`} />
              <div className="text-xs font-mono text-text-muted mb-2">{item.status}</div>
              <h4 className="font-semibold text-white mb-2">{item.title}</h4>
              <p className="text-sm text-text-muted">{item.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

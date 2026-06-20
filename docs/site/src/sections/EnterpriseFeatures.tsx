import { Button } from "@/components/ui/Button"

export default function EnterpriseFeatures() {
  return (
    <section className="py-24 border-y border-white/5 relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,rgba(0,112,243,0.1),transparent_70%)]" />
      <div className="container mx-auto px-6 max-w-container relative z-10 flex flex-col md:flex-row items-center justify-between gap-12">
        <div className="max-w-xl">
          <h2 className="text-3xl font-bold mb-4">Trusted by the teams building the future.</h2>
          <p className="text-text-muted mb-8">
            Need custom SLA, dedicated support, or on-prem deployment? Our enterprise plan is designed for massive scale and stringent security requirements.
          </p>
          <Button variant="outline">Contact Sales</Button>
        </div>
        <div className="w-full md:w-1/2 grid grid-cols-2 gap-4 opacity-50 grayscale font-mono text-sm">
          {["Company A", "Startup B", "Research Lab C", "Enterprise D"].map(logo => (
            <div key={logo} className="h-20 flex items-center justify-center border border-white/10 rounded-lg">
              {logo}
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

import { Button } from "@/components/ui/Button"

export default function PricingSection() {
  return (
    <section className="py-24 relative overflow-hidden" id="pricing">
        <div className="container mx-auto px-6 max-w-container relative z-10">
            <div className="text-center mb-16">
            <h2 className="text-3xl md:text-5xl font-bold tracking-tighter mb-4 text-white">
                Simple, transparent pricing
            </h2>
            <p className="text-lg text-zinc-400 max-w-2xl mx-auto">
                No hidden fees or unexpected charges. Pay only for the resources you use.
            </p>
            </div>

            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {/* Starter Plan */}
            <div className="p-8 border border-white/10 rounded-2xl bg-[#0A0A0A] flex flex-col hover:border-white/20 transition-colors">
                <h3 className="text-xl font-bold text-white mb-2">Starter</h3>
                <div className="text-4xl font-bold text-white mb-6">$0<span className="text-lg text-zinc-500 font-normal">/mo</span></div>
                <p className="text-zinc-400 mb-8">For individuals and small projects exploring AI.</p>
                <ul className="mb-8 flex-1 space-y-3">
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Up to 4 GPUs
                </li>
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Community Support
                </li>
                </ul>
                <Button variant="outline" className="w-full">Get Started</Button>
            </div>

            {/* Pro Plan */}
            <div className="p-8 border border-blue-500/50 rounded-2xl bg-blue-900/10 flex flex-col relative shadow-[0_0_40px_-10px_rgba(59,130,246,0.3)] transform md:-translate-y-4">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-blue-500 text-white text-xs font-bold px-3 py-1 rounded-full uppercase tracking-wider">
                Most Popular
                </div>
                <h3 className="text-xl font-bold text-white mb-2">Pro</h3>
                <div className="text-4xl font-bold text-white mb-6">$299<span className="text-lg text-zinc-500 font-normal">/mo</span></div>
                <p className="text-zinc-400 mb-8">For growing teams requiring more power and priority support.</p>
                <ul className="mb-8 flex-1 space-y-3">
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Up to 64 GPUs
                </li>
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Priority Support
                </li>
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Advanced Telemetry
                </li>
                </ul>
                <Button className="w-full">Upgrade to Pro</Button>
            </div>

            {/* Enterprise Plan */}
            <div className="p-8 border border-white/10 rounded-2xl bg-[#0A0A0A] flex flex-col hover:border-white/20 transition-colors">
                <h3 className="text-xl font-bold text-white mb-2">Enterprise</h3>
                <div className="text-4xl font-bold text-white mb-6">Custom</div>
                <p className="text-zinc-400 mb-8">For large organizations with complex needs.</p>
                <ul className="mb-8 flex-1 space-y-3">
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Unlimited GPUs
                </li>
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Dedicated Account Manager
                </li>
                <li className="flex items-center gap-2 text-zinc-300">
                    <span className="text-blue-500">✓</span> Custom SLAs
                </li>
                </ul>
                <Button variant="outline" className="w-full">Contact Sales</Button>
            </div>
            </div>
        </div>
    </section>
  )
}

import { Button } from "@/components/ui/Button"

export default function CTA() {
  return (
    <section className="py-32 relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-b from-transparent to-blue-900/20 pointer-events-none" />
      <div className="container mx-auto px-6 max-w-container text-center relative z-10">
        <h2 className="text-4xl md:text-6xl font-bold tracking-tight mb-6">
          Start orchestrating today.
        </h2>
        <p className="text-xl text-text-muted mb-10 max-w-2xl mx-auto">
          Join the teams building the future of AI. Deploy your first workload in minutes.
        </p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button size="lg" className="w-full sm:w-auto">Deploy Now</Button>
          <Button variant="outline" size="lg" className="w-full sm:w-auto">Book a Demo</Button>
        </div>
      </div>
    </section>
  )
}

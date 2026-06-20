export default function TrustedTechnologies() {
  const logos = [
    { name: "PyTorch", svg: "PyTorch" },
    { name: "Ray", svg: "Ray" },
    { name: "Kubernetes", svg: "Kubernetes" },
    { name: "AWS", svg: "AWS" },
    { name: "GCP", svg: "GCP" },
    { name: "CoreWeave", svg: "CoreWeave" },
  ]

  return (
    <section className="py-12 border-b border-white/5">
      <div className="container mx-auto px-6 max-w-container">
        <p className="text-center text-sm font-medium text-text-muted mb-8 tracking-wider uppercase">
          Integrates seamlessly with your existing stack
        </p>
        <div className="flex flex-wrap justify-center items-center gap-8 md:gap-16 opacity-50 grayscale hover:grayscale-0 transition-all duration-500">
          {logos.map((logo) => (
            <div key={logo.name} className="text-xl font-bold font-mono tracking-tighter text-white">
              {logo.svg}
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

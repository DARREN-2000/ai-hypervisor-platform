
import { Lock, ShieldCheck, Key } from "lucide-react"

export default function Security() {
  return (
    <section className="py-32 bg-bg-secondary border-t border-white/5">
      <div className="container mx-auto px-6 max-w-container">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-6">Enterprise-grade by default.</h2>
          <p className="text-lg text-text-muted leading-relaxed">
            Unblock procurement and compliance teams from day one. Onyx is built with zero-trust principles, isolated execution environments, and comprehensive audit logging.
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-6">
          {[
            { icon: <Lock />, title: "SOC2 Type II", desc: "Certified and compliant infrastructure." },
            { icon: <ShieldCheck />, title: "Namespace Isolation", desc: "Hard multitenancy for isolated workloads." },
            { icon: <Key />, title: "SSO & RBAC", desc: "Integrate with Okta, Azure AD, and more." },
          ].map((item, i) => (
            <div key={i} className="glass p-6 rounded-xl flex items-start gap-4">
              <div className="p-3 bg-white/5 rounded-lg text-white">{item.icon}</div>
              <div>
                <h4 className="font-semibold text-white mb-1">{item.title}</h4>
                <p className="text-sm text-text-muted">{item.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

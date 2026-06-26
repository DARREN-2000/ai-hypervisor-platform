import { Link } from "react-router-dom"
import { Github, Twitter, Linkedin } from "lucide-react"

export default function Footer() {
  const footerLinks = {
    Product: [
      { name: "Features", href: "#features" },
      { name: "Integrations", href: "#integrations" },
      { name: "Pricing", href: "/pricing" },
      { name: "Changelog", href: "/changelog" },
    ],
    Resources: [
      { name: "Documentation", href: "/docs" },
      { name: "API Reference", href: "/api" },
      { name: "Blog", href: "/blog" },
      { name: "Community", href: "/community" },
    ],
    Company: [
      { name: "About", href: "/about" },
      { name: "Customers", href: "/customers" },
      { name: "Enterprise", href: "/enterprise" },
      { name: "Careers", href: "/careers" },
    ],
    Legal: [
      { name: "Privacy Policy", href: "/privacy" },
      { name: "Terms of Service", href: "/terms" },
      { name: "Security", href: "/security" },
    ],
  }

  return (
    <footer className="bg-bg-secondary border-t border-white/10 pt-20 pb-10">
      <div className="container mx-auto px-6 max-w-container">
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-10 mb-16">
          <div className="col-span-2 lg:col-span-1">
            <Link to="/" className="flex items-center gap-2 mb-6">
              <div className="w-6 h-6 rounded bg-white text-black flex items-center justify-center font-bold text-xs">
                O
              </div>
              <span className="font-semibold text-lg tracking-tight">Onyx</span>
            </Link>
            <p className="text-text-muted text-sm mb-6 max-w-xs">
              The operating system for AI infrastructure. Unify, orchestrate, and scale with a single command.
            </p>
            <div className="flex items-center gap-4 text-text-muted">
              <a href="https://github.com/onyx" aria-label="GitHub" className="hover:text-white transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white rounded-sm"><Github size={20} /></a>
              <a href="https://twitter.com/onyx" aria-label="Twitter" className="hover:text-white transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white rounded-sm"><Twitter size={20} /></a>
              <a href="https://linkedin.com/company/onyx" aria-label="LinkedIn" className="hover:text-white transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white rounded-sm"><Linkedin size={20} /></a>
            </div>
          </div>

          {Object.entries(footerLinks).map(([category, links]) => (
            <div key={category}>
              <h4 className="font-semibold text-white mb-4">{category}</h4>
              <ul className="space-y-3">
                {links.map((link) => (
                  <li key={link.name}>
                    <Link to={link.href} className="text-sm text-text-muted hover:text-white transition-colors">
                      {link.name}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="border-t border-white/10 pt-8 flex flex-col md:flex-row items-center justify-between gap-4">
          <p className="text-text-muted text-sm">
            © {new Date().getFullYear()} Onyx Inc. All rights reserved.
          </p>
          <div className="flex items-center gap-2 text-sm text-text-muted">
            <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
            All systems operational
          </div>
        </div>
      </div>
    </footer>
  )
}

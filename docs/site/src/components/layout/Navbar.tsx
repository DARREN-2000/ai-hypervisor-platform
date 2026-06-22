import { useState, useEffect } from "react"
import { Link, useLocation } from "react-router-dom"
import { Menu, X, Command } from "lucide-react"
import { Button } from "@/components/ui/Button"
import { cn } from "@/utils/cn"

export default function Navbar() {
  const [scrolled, setScrolled] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const location = useLocation()

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 20)
    }
    window.addEventListener("scroll", handleScroll)
    return () => window.removeEventListener("scroll", handleScroll)
  }, [])

  useEffect(() => {
    setMobileMenuOpen(false)
  }, [location])

  const navLinks = [
    { name: "Features", href: "#features" },
    { name: "Documentation", href: "/docs" },
    { name: "Pricing", href: "/pricing" },
    { name: "Enterprise", href: "/enterprise" },
  ]

  return (
    <header
      className={cn(
        "fixed top-0 left-0 right-0 z-50 transition-all duration-300 border-b border-transparent",
        scrolled ? "bg-black/50 backdrop-blur-md border-white/10 py-3" : "bg-transparent py-5"
      )}
    >
      <div className="container mx-auto px-6 max-w-container flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Link to="/" className="flex items-center gap-2 group">
            <div className="w-8 h-8 rounded bg-white text-black flex items-center justify-center font-bold text-lg group-hover:scale-105 transition-transform shadow-[0_0_15px_rgba(255,255,255,0.3)]">
              O
            </div>
            <span className="font-semibold text-lg tracking-tight">Onyx</span>
          </Link>
        </div>

        {/* Desktop Nav */}
        <nav className="hidden md:flex items-center gap-8">
          {navLinks.map((link) => (
            <Link
              key={link.name}
              to={link.href}
              className="text-sm text-text-muted hover:text-white transition-colors"
            >
              {link.name}
            </Link>
          ))}
        </nav>

        <div className="hidden md:flex items-center gap-4">
          <button
            aria-label="Search"
            className="text-text-muted hover:text-white flex items-center gap-2 text-sm px-3 py-1.5 rounded-md border border-white/10 bg-white/5 transition-colors"
          >
            <Command size={14} />
            <span>K</span>
          </button>
          <Button variant="ghost" className="text-sm">Log in</Button>
          <Button size="sm">Deploy Now</Button>
        </div>

        {/* Mobile Toggle */}
        <button
          aria-label="Toggle mobile menu"
          aria-expanded={mobileMenuOpen}
          className="md:hidden text-white"
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
        >
          {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden absolute top-full left-0 right-0 bg-black/95 backdrop-blur-xl border-b border-white/10 p-6 flex flex-col gap-4 shadow-2xl">
          {navLinks.map((link) => (
            <Link
              key={link.name}
              to={link.href}
              className="text-base text-text-muted hover:text-white py-2 border-b border-white/5"
            >
              {link.name}
            </Link>
          ))}
          <div className="flex flex-col gap-3 mt-4">
            <Button variant="outline" className="w-full justify-center">Log in</Button>
            <Button className="w-full justify-center">Deploy Now</Button>
          </div>
        </div>
      )}
    </header>
  )
}

import { Link } from "react-router-dom"
import { Button } from "@/components/ui/Button"
import { Construction } from "lucide-react"

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-4 relative">
      {/* Subtle Glow Effect */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-64 h-64 bg-accent-glow/20 rounded-full blur-[80px] -z-10"></div>

      <div className="bg-white/5 border border-white/10 rounded-2xl p-8 max-w-md w-full backdrop-blur-sm shadow-xl flex flex-col items-center">
        <div className="w-16 h-16 rounded-full bg-white/5 border border-white/10 flex items-center justify-center mb-6 text-text-muted">
           <Construction size={32} />
        </div>

        <h1 className="text-5xl font-bold text-white mb-2 tracking-tight">404</h1>
        <h2 className="text-xl font-semibold text-text-main mb-4">Page Not Found</h2>

        <p className="text-text-muted mb-8">
          The page you are looking for might have been removed, had its name changed, or is temporarily unavailable as we continue building the Onyx platform.
        </p>

        <Link to="/">
          <Button className="px-8">
            Return Home
          </Button>
        </Link>
      </div>
    </div>
  )
}

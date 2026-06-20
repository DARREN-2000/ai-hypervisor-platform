import React from "react"
import Navbar from "@/components/layout/Navbar"
import Footer from "@/components/layout/Footer"

interface MainLayoutProps {
  children: React.ReactNode
}

export default function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className="min-h-screen flex flex-col relative text-text-main bg-background overflow-x-hidden">
      {/* Global Background Glow */}
      <div className="fixed inset-0 pointer-events-none z-[-1] flex items-center justify-center opacity-30">
        <div className="w-[800px] h-[800px] bg-accent-glow rounded-full blur-[150px] opacity-20"></div>
        <div className="w-[600px] h-[600px] bg-purple-600/20 rounded-full blur-[150px] absolute top-1/4 left-1/4"></div>
      </div>

      <Navbar />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  )
}

import { HashRouter as Router, Routes, Route } from "react-router-dom"
import MainLayout from "@/layouts/MainLayout"
import Home from "@/pages/Home"

function App() {
  return (
    <Router>
      <MainLayout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="*" element={<div className="min-h-screen flex items-center justify-center pt-24"><h1 className="text-4xl font-bold">404 - Not Found</h1></div>} />
        </Routes>
      </MainLayout>
    </Router>
  )
}

export default App

import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import MainLayout from "@/layouts/MainLayout"
import Home from "@/pages/Home"
import NotFound from "@/pages/NotFound"
import Docs from "@/pages/Docs"
import Pricing from "@/pages/Pricing"
import Enterprise from "@/pages/Enterprise"

function App() {
  return (
    <Router basename="/ai-hypervisor-platform">
      <MainLayout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/docs" element={<Docs />} />
          <Route path="/pricing" element={<Pricing />} />
          <Route path="/enterprise" element={<Enterprise />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </MainLayout>
    </Router>
  )
}

export default App

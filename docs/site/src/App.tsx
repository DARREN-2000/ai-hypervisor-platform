import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import MainLayout from "@/layouts/MainLayout"
import Home from "@/pages/Home"
import NotFound from "@/pages/NotFound"

function App() {
  return (
    <Router basename="/ai-hypervisor-platform">
      <MainLayout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </MainLayout>
    </Router>
  )
}

export default App

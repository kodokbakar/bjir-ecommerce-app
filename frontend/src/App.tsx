import { Routes, Route } from 'react-router-dom'

function App() {
  return (
    <Routes>
      <Route path="/" element={<div className="p-8 text-center"><h1 className="text-3xl font-bold text-blue-600">E-Commerce Frontend</h1><p className="mt-4 text-gray-600">Setup berhasil! Tailwind CSS berfungsi.</p></div>} />
    </Routes>
  )
}

export default App

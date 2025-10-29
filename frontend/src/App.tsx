import { Route, Routes } from 'react-router-dom'
import './App.css'
import Main from './pages/Main'
import { TopTen } from './pages/TopTen'

function App() {
  return (
    <Routes>
      <Route path='/' element={<Main />} />
      <Route path='/topten/:roomName' element={<TopTen />} />
    </Routes>
  )
}

export default App

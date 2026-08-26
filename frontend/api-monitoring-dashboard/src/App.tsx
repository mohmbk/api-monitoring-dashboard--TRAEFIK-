
import { BrowserRouter as Router, Routes, Route , Navigate  } from "react-router-dom";

import './App.css'


import Navbar from './components/navbar/nav';
import Login from './components/pages/login/login';
import Dashboard from './components/pages/dashboard/dashboard';
import Signup from './components/pages/signup/signup';
import History from './components/pages/history/history';
function App() {

  const showNavbar =
    location.pathname !== "/login" &&
    location.pathname !== "/signup";

  return (
    <>
      <Router>
        {showNavbar && <Navbar/>}
        <div>
          <Routes>
            <Route path="/" element={<Navigate to="/login" />} />
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/dashboard/api/:id/history" element={<History />} />
          </Routes>
        </div>
      </Router>
    </>
  )
}

export default App

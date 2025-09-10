import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Registration from "./pages/Registration";
import Login from "./pages/Login";
import Home from "./pages/Home";

const App = () => {
  return (
    <>
      <Router>
          <Routes>
            <Route path="/" element={<Login />}/>
            <Route path="/register" element={<Registration />}/>
            <Route path="/home" element={<Home />} />
          </Routes>
      </Router>
    </>
  )
}

export default App;
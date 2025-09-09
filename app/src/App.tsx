import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Registration from "./pages/Registration";
import Login from "./pages/Login";

const App = () => {
  return (
    <>
      <Router>
          <Routes>
            <Route path="/" element={<Login />}/>
            <Route path="/register" element={<Registration />}/>
          </Routes>
      </Router>
    </>
  )
}

export default App;
import { Link } from "react-router-dom";
import "./nav.css";

function Navbar() {
  return (
    <section className="nav">
      <h1>API Monitoring System</h1>

      <ul>
        <li>
          <Link to="/dashboard" className="text">
            Dashboard
          </Link>
        </li>
      </ul>
    </section>
  );
}

export default Navbar;
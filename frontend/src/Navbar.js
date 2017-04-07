import { MenuItem, Nav, Navbar, NavDropdown, NavItem } from 'react-bootstrap';
import logo from './img/lockHollow.svg';
import './Navbar.css'

const logoutLink = (
  <NavItem href="logout">Logout</NavItem>
);

const navbarInstance = (
  <Navbar inverse collapseOnSelect fixedTop>
    <Navbar.Header>
      <Navbar.Brand>
        <a href="#">
          <span className="App-header">
            <img src={logo} className="App-logo" alt="logo" />
            Security Dashboard
          </span>
        </a>
      </Navbar.Brand>
      <Navbar.Toggle />
    </Navbar.Header>
    <Navbar.Collapse>
      {projectList}
      <Nav pullRight>
        {userLink}
      </Nav>
    </Navbar.Collapse>
  </Navbar>
);

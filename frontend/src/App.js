import React, { Component } from 'react';
import { Button, ButtonGroup, DropdownButton, MenuItem, Nav, Navbar, NavDropdown, NavItem } from 'react-bootstrap';
import logo from './img/lockHollow.svg';
import './App.css';

var loggedIn = false;
var userLink;
var projectList;
var userButton;

//var projectList = (
//   <Nav>
//     <NavDropdown eventKey={1} title="Projects" id="basicNavDropdown">
//       // <MenuItem eventKey={1.1}>Action</MenuItem>
//       // <MenuItem eventKey={1.2}>Another action</MenuItem>
//       // <MenuItem eventKey={1.3}>Something else here</MenuItem>
//       // <MenuItem divider />
//       // <MenuItem eventKey={1.3}>Separated link</MenuItem>
//     </NavDropdown>
//   </Nav>
// );

const buttonGroupInstance = (
  <ButtonGroup>
    <DropdownButton id="dropdown-btn-menu" bsStyle="success" title="Dropdown">
      <MenuItem key="1">Dropdown link</MenuItem>
      <MenuItem key="2">Dropdown link</MenuItem>
    </DropdownButton>
    <Button bsStyle="info">Middle</Button>
    <Button bsStyle="info">Right</Button>
  </ButtonGroup>
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

class App extends Component {

  render() {
    return (
      <div className="App">
        {navbarInstance}
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
        {buttonGroupInstance}
      </div>
    );
  }
}

export default App;

import React, { Component } from 'react';
import {
  MenuItem,
  Nav,
  Navbar,
  NavDropdown
} from 'react-bootstrap';

import logo from './img/lockHollow.svg';
import './Title.css'

const UserButton = (props) => (
  <NavDropdown eventKey={1} title={props.user} id="basicNavDropdown">
    <MenuItem>Logout</MenuItem>
  </NavDropdown>
)

const ProjectList = (
  <NavDropdown eventKey={1} title="Projects" id="basicNavDropdown">
    <MenuItem>Overview</MenuItem>
    <MenuItem divider />
    <MenuItem>Payment</MenuItem>
    <MenuItem>OfferAPI</MenuItem>
    <MenuItem>MobileConnect</MenuItem>
  </NavDropdown>
);

export default class Title extends Component {

  state = {
    user: 'marcelo@telenordigital.com'
  }

  render() {
    return (
      <Navbar inverse collapseOnSelect fixedTop>
        <Navbar.Header>
          <Navbar.Brand>
            <a href="#">
              <img src={logo} className="Title-logo" alt="logo" />
            </a>
          </Navbar.Brand>
          <Navbar.Brand>
            <a href="#">
              Securidash
            </a>
          </Navbar.Brand>
          <Navbar.Toggle />
        </Navbar.Header>
        <Navbar.Collapse>
          <Nav pullLeft>
            {ProjectList}
          </Nav>
          <Nav pullRight>
            <UserButton user={this.state.user}/>
          </Nav>
        </Navbar.Collapse>
      </Navbar>
    )
  }
}

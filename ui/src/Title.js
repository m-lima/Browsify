import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {
  MenuItem,
  Nav,
  Navbar,
  NavDropdown
} from 'react-bootstrap';

import * as Constants from './Constants.js'
import logo from './img/lockHollow.svg';

const UserButton = (props) => (
  <NavDropdown id={'user-dropdown'} eventKey={1} title={props.user}>
    <MenuItem eventKey={1.1} onClick={() =>
      fetch(Constants.logout, { method: 'POST', credentials: 'include' })
        .then(window.location.reload())}>
      Logout</MenuItem>
  </NavDropdown>
)

const ProjectList = (
  <NavDropdown eventKey={1} title='Projects' id='project-dropdown'>
    <MenuItem eventKey={1.1}>Overview</MenuItem>
    <MenuItem divider />
    <MenuItem eventKey={1.1}>Payment</MenuItem>
    <MenuItem eventKey={1.1}>OfferAPI</MenuItem>
    <MenuItem eventKey={1.1}>MobileConnect</MenuItem>
  </NavDropdown>
)

export default class Title extends Component {

  render() {
    return (
      <Navbar inverse collapseOnSelect fixedTop>
        <Navbar.Header>
          {this.props.user &&
            <Navbar.Brand>
              <Link to={Constants.ui}>
                <img src={logo} className='Title-logo' alt='logo' style={{ height: '100%' }} />
              </Link>
            </Navbar.Brand>
          }
          <Navbar.Brand>
            <Link to={Constants.ui}>
              Securidash
            </Link>
          </Navbar.Brand>
          <Navbar.Toggle />
        </Navbar.Header>
          {this.props.user
            ? <Navbar.Collapse>
                <Nav pullLeft>
                  {ProjectList}
                </Nav>
                <Nav pullRight>
                  <UserButton user={this.props.user}/>
                </Nav>
              </Navbar.Collapse>
            : <Navbar.Collapse>
                <Navbar.Text pullRight>
                  <Navbar.Link href={Constants.login} style={{ textDecoration: 'none' }}>Login</Navbar.Link>
                </Navbar.Text>
              </Navbar.Collapse>
          }
      </Navbar>
    )
  }
}

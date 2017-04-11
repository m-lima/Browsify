import React, { Component } from 'react';
import {
  MenuItem,
  Nav,
  Navbar,
  NavDropdown,
  Button
} from 'react-bootstrap';

import logo from './img/lockHollow.svg';
import './Title.css'

const UserButton = (props) => (
  <NavDropdown eventKey={1} title={props.user}>
    <Navbar.Form>
      <form action='https://localhost/logout' method='post'>
        <Button className='menuitem' bsSize='xsmall' bsStyle='info' type='submit'>Logout</Button>
      </form>
    </Navbar.Form>
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
)

export default class Title extends Component {

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
                <Nav pullRight>
                  <Navbar.Form>
                    <form action='https://localhost/login' method='post'>
                      <Button className='menuitem' bsStyle='info' type='submit'>Login</Button>
                    </form>
                  </Navbar.Form>
                </Nav>
              </Navbar.Collapse>
          }
      </Navbar>
    )
  }
}

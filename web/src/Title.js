import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {
  MenuItem,
  Nav,
  Navbar,
  NavDropdown,
  Checkbox,
  Image
} from 'react-bootstrap';

import * as Constants from './Constants.js'
import logo from './img/lock.svg';

const UserDropdown = (user) =>
  <span>
    <Image src={user.Avatar} alt='avatar' style={{ height: 20, marginRight: 10 }} rounded />
    {user.Email}
  </span>

const ShowHidden = (props) => (
  (props.user.Admin || props.user.CanShowHidden) &&
    <MenuItem onClick={() => {
        props.user.ShouldShowHidden = !props.user.ShouldShowHidden
        props.updater(props.user)
      }}>
      {props.user.ShouldShowHidden
        ? <Checkbox checked >Hidden</Checkbox>
        : <Checkbox>Hidden</Checkbox>
      }
    </MenuItem>
)

const ShowProtected = (props) => (
  (props.user.Admin || props.user.CanShowProtected) &&
    <MenuItem onClick={() => {
        props.user.ShouldShowProtected = !props.user.ShouldShowProtected
        props.updater(props.user)
      }}>
      {props.user.ShouldShowProtected
        ? <Checkbox checked >Protected</Checkbox>
        : <Checkbox>Protected</Checkbox>
      }
    </MenuItem>
)

const UserButton = (props) => (
  props.user
  ? <NavDropdown id='user-dropdown' title={UserDropdown(props.user)} eventKey={1}>
      {props.user.Admin && <MenuItem>Admin Panel</MenuItem>}
      {props.user.Admin && <MenuItem divider />}

      <ShowHidden user={props.user} updater={props.updater} />
      <ShowProtected user={props.user} updater={props.updater} />
      {(props.user.Admin || props.user.CanShowHidden || props.user.CanShowProtected) && <MenuItem divider />}

      <MenuItem eventKey={1.1} onClick={() =>
        fetch(Constants.logout, { method: 'POST', credentials: 'include' })
          .then(window.location.reload())}>
        Logout
      </MenuItem>
    </NavDropdown>
  : <Navbar.Text>Loading..</Navbar.Text>
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
          {this.props.authorized &&
            <Navbar.Brand>
              <Link to={Constants.ui}>
                <img src={logo} className='Title-logo' alt='logo' style={{ height: '20px' }} />
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
          {this.props.authorized
            ? <Navbar.Collapse>
                <Nav pullLeft>
                  {ProjectList}
                </Nav>
                <Nav pullRight>
                  <UserButton user={this.props.user} updater={this.props.updater} />
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

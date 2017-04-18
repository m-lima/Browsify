import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {
  MenuItem,
  Nav,
  Navbar,
  NavDropdown,
  NavItem,
  Checkbox,
  Image
} from 'react-bootstrap';

import * as Constants from './Constants.js'
import logo from './img/lock.svg';

const UserDropdown = (user) =>
  <span>
    <Image
      src={user.Avatar}
      alt=''
      style={{ height: 36, marginTop: -8, marginBottom: -8, marginRight: 10 }}
      rounded
    />
    {user.Email}
  </span>

const ShowHidden = (props) => (
  (props.user.Admin || props.user.CanShowHidden) &&
    <MenuItem onClick={() => {
        props.user.ShouldShowHidden = !props.user.ShouldShowHidden
        props.updater(props.user)
      }}>
      {props.user.ShouldShowHidden
        ? <Checkbox defaultChecked >Hidden</Checkbox>
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
        ? <Checkbox defaultChecked >Protected</Checkbox>
        : <Checkbox>Protected</Checkbox>
      }
    </MenuItem>
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

  state = {
    user: null,
    loading: false
  }

  constructor(props) {
    super(props)
    this.updateUser = this.updateUser.bind(this)
  }

  componentDidMount() {
    this.fetchUser()
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.authorized !== (this.state.user !== null)) {
      this.fetchUser()
    }
  }

  invalidateUser() {
    this.setState({ user: null, loading: false })
    this.props.authUpdater(false)
  }

  performUserFetch(url, request) {
    if (!url || !request) {
      this.invalidateUser()
      return
    }

    this.setState({ loading: true })
    fetch(url, request)
      .then(response => {
        if (response.ok) {
          response.json()
            .then(newUser => {
              this.setState({ user: newUser, loading: false })
              this.props.authUpdater(true)
            })
            .catch(this.invalidateUser)
        } else {
          this.invalidateUser()
        }
      })
      .catch(this.invalidateUser)
  }

  fetchUser() {
    this.performUserFetch(Constants.user, { method: 'GET', credentials: 'include' })
  }

  updateUser(user) {
    var req = {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      },
      body: 'ShouldShowHidden='
        + user.ShouldShowHidden
        + '&ShouldShowProtected='
        + user.ShouldShowProtected
    }

    this.props.refresher()
    this.performUserFetch(Constants.userUpdate, req)
  }

  renderUserButton() {
    if (this.state.loading) {
      return <Navbar.Text pullRight>Loading..</Navbar.Text>
    }

    if (this.state.user) {
      var user = this.state.user
      return (
        <Nav pullRight>
          <NavDropdown
            id='user-dropdown'
            title={UserDropdown(user)}
            style={{ height: 50 }}
            eventKey={1}>

            {user.Admin && <MenuItem>Admin Panel</MenuItem>}
            {user.Admin && <MenuItem divider />}

            <ShowHidden user={user} updater={this.updateUser} />
            <ShowProtected user={user} updater={this.updateUser} />
            {(user.Admin || user.CanShowHidden || user.CanShowProtected) && <MenuItem divider />}

            <MenuItem eventKey={1.1} onClick={() =>
              fetch(Constants.logout, { method: 'POST', credentials: 'include' })
                .then(this.fetchUser())}>
              Logout
            </MenuItem>
          </NavDropdown>
        </Nav>
      )
    }

    return (
      <Nav pullRight>
        <NavItem onClick={() => window.location = Constants.login}>Login</NavItem>
      </Nav>
    )
  }

  render() {
    return (
      <Navbar inverse collapseOnSelect fixedTop>
        <Navbar.Header>
          {this.state.user &&
            <Navbar.Brand>
              <Link to={Constants.ui}>
                <img src={logo} alt='logo' style={{ height: 20 }} />
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
        <Navbar.Collapse>
          {this.state.user &&
            <Nav pullLeft>
              {ProjectList}
            </Nav>
          }
          {this.renderUserButton()}
        </Navbar.Collapse>
      </Navbar>
    )
  }
}

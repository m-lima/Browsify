import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom';
import {
  Grid,
  Row
} from 'react-bootstrap';

import Title from './Title.js'
import BrowseList from './Browse.js'
import Landing from './Landing.js'
import * as Constants from './Constants.js'

class StatefulRenderer extends Component {

  state = {
    basePath: '',
    entries: null,
    status: null,
    user: null
  }

  constructor(props) {
    super(props)
    this.updateUser = this.updateUser.bind(this)
  }

  checkUser() {
    if (!this.state.user) {
      fetch(Constants.user, { method: 'GET', credentials: 'include' })
        .then(response => {
          if (response.ok) {
            response.json()
              .then(newUser => this.setState({ user: newUser }))
          }
        })
    }
  }

  fetchFiles(path) {
    fetch(Constants.api + path, { method: 'GET', credentials: 'include' })
      .then(response => {
        if (response.ok) {
          response.json().then(newEntries => {
            newEntries
            ? this.setState({ entries: newEntries, status: null })
            : this.setState({ entries: [], status: null })
          })
          .catch(err => {
            this.setState({ entries: [], status: err.message })
          })
        } else {
          this.setState({ entries: [], status: response.status })
        }})
      .catch(err =>
        this.setState({ entries: [], status: err.message })
      )
  }

  fetchData(path) {
    if (!this.props.location.pathname.startsWith(Constants.ui)) {
      return
    }

    if (path === undefined || path === '') {
      path = this.props.location.pathname
    }

    path = path.substring(4)
    if (path.length > 0 && path.charAt(path.length - 1) !== '/') {
      path += '/'
    }

    this.setState({ basePath: path, entries: null })
    this.fetchFiles(path)
    this.checkUser()
  }

  updateUser(user) {
    this.setState({ user: null })
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

    fetch(Constants.userUpdate, req)
      .then(response => {
        if (response.ok) {
          this.fetchFiles(this.state.basePath)
          response.json()
            .then(newUser => this.setState({ user: newUser }))
            .catch(err =>
              this.setState({ entries: [], status: Constants.statusUnauthorized })
            )
        } else {
          this.setState({ entries: [], status: Constants.statusUnauthorized })
        }
      })
      .catch(err =>
        this.setState({ entries: [], status: Constants.statusUnauthorized })
      )
  }

  componentDidMount() {
    if (!this.state.entries) {
      this.fetchData(this.props.location.pathname)
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps !== this.props) {
      this.fetchData(nextProps.location.pathname)
    }
  }

  render() {
    switch (this.state.status) {
      case Constants.statusNotFound:
      case null:
        return (
          <div style={{ paddingTop: 80 }} >
            <Title user={this.state.user} authorized={true} updater={this.updateUser} />
            <BrowseList basePath={this.state.basePath} entries={this.state.entries} status={this.state.status} />
          </div>
        )
      case Constants.statusUnauthorized:
        return(
          <div style={{ height: '100%' }} >
            <Title />
            <Landing />
          </div>
        )
      case Constants.statusForbidden:
        return(
          <div style={{ height: '100%', paddingTop: 80, backgroundColor: '#222222', color: 'lightgray' }} >
            <Title />
            <Grid>
              <Row>
                <h3>Unauthorized</h3>
                <p>The user you logged in with does have access to this site</p>
                <a href={Constants.login}>Retry</a>
              </Row>
            </Grid>
          </div>
        )
      default:
        return(
          <div style={{ height: '100%', paddingTop: 80, backgroundColor: '#222222', color: 'lightgray' }} >
            <Title user={this.state.user} authorized={true} updater={this.updateUser} />
            <Grid>
              <Row>
                <h3>Oops!</h3>
                <p>An error occured while processing your request</p>
                <a href={Constants.ui}>Main page</a>
              </Row>
            </Grid>
          </div>
        )
    }
  }
}

export default class App extends Component {
  render() {
    return (
      <Router>
        <Route path='/' component={StatefulRenderer} />
      </Router>
    )
  }
}

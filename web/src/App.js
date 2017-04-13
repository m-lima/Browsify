import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom';

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

  checkUser() {
    if (!this.state.user) {
      fetch(Constants.user, { method: 'GET', credentials: 'include' })
        .then(response => {
          if (response.ok) {
            response.text()
              .then(newUser => this.setState({ user: newUser }))
          } else {
            console.log('Error fetching user')
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
      .catch(err => {
        this.setState({ entries: [], status: err.message })
      })
  }

  fetchData(path) {
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
    if (this.state.status === Constants.statusUnauthorized) {
      return (
        <div style={{ height: '100%'}} >
          <Title />
          <Landing />
        </div>
      )
    } else {
      return (
        <div style={{ marginTop: 80 }} >
          <Title user={this.state.user ? this.state.user : 'Loading..'} />
          <BrowseList basePath={this.state.basePath} entries={this.state.entries} status={this.state.status} />
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

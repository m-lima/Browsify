import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom';

import Title from './Title.js'
import BrowseList from './Browse.js'
import Landing from './Landing.js'

// Update for apiPrefix (remove host)
const urlPath = 'https://localhost/api/';
const userPath = 'https://localhost/user/';

class StatefulRenderer extends Component {

  state = {
    basePath: '',
    entries: null,
    status: null,
    user: null
  }

  fetchData(path) {
    if (path === undefined || path === '') {
      path = this.props.location.pathname
    }

    path = path.substring(4)
    if (path.length > 0) {
      path += '/'
    }

    this.setState({ basePath: path, entries: null })

    fetch(urlPath + path, { method: 'GET', credentials: 'include' })
      .then(response => {
        if (response.ok) {
          response.json().then(newEntries => {
            newEntries
            ? this.setState({ basePath: path, entries: newEntries, status: null })
            : this.setState({ basePath: path, entries: [], status: null })
          })
        } else {
          throw new Error(response.status)
        }})
      .then(() => {
        if (!this.state.user) {
          fetch(userPath, { method: 'GET', credentials: 'include' })
            .then(response => {
              if (response.ok) {
                response.json().then( newUser => {
                  newUser
                  ? this.setState({ user: newUser })
                  : this.setState({ entries: [], user: null, status: '403' })
                })
              } else {
                throw new Error('403')
              }
            })
        }})
      .catch(err => {
        var newUser = this.state.user
        if (newUser && err.message === '403') {
          newUser = null
        }

        this.setState({ basePath: path, entries: [], status: err.message, user: newUser })
      })
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
    if (this.state.status === '403') {
      return (
        <div style={{ height: '100%'}} >
          <Title />
          <Landing />
        </div>
      )
    } else {
      return (
        <div>
          <Title user={this.state.user ? this.state.user.Email : 'Loading..'} />
          <div style={{ marginTop: 80 }} >
            <BrowseList basePath={this.state.basePath} entries={this.state.entries} status={this.state.status} />
          </div>
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

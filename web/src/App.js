import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom';

import Title from './Title.js'
import Main from './Main.js'
import * as Constants from './Constants.js'

class StatefulRenderer extends Component {

  state = {
    authorized: false,
    refresh: false
  }

  constructor(props) {
    super(props)
    this.updateAuth = this.updateAuth.bind(this)
    this.refresh = this.refresh.bind(this)
  }

  updateAuth(authorized) {
    this.setState({ authorized: authorized })
  }

  refresh() {
    this.setState({ refresh: !this.state.refresh })
  }

  cleanPath(path) {
    if (path === undefined || path === '') {
      path = this.props.location.pathname
    }

    if (!this.props.location.pathname.startsWith(Constants.ui)) {
      return
    }

    path = path.substring(4)
    if (path.length > 0 && path.charAt(path.length - 1) !== '/') {
      path += '/'
    }

    return path
  }

  render() {
    return (
      <div style={{ height: '100%' }}>
        <Title
          authorized={this.state.authorized}
          authUpdater={this.updateAuth}
          refresher={this.refresh}
        />
        <Main
          authorized={this.state.authorized}
          authUpdater={this.updateAuth}
          path={this.cleanPath(this.props.location.basePath)}
          refresh={this.state.refresh}
        />
      </div>
    )
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

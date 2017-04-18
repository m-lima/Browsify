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
    status: null
  }

  constructor(props) {
    super(props)
    this.updateStatus = this.updateStatus.bind(this)
  }

  updateStatus(status) {
    this.setState({ status: status })
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
        <Title status={this.state.status} updater={this.updateStatus} />
        <Main status={this.state.status} updater={this.updateStatus} path={this.cleanPath(this.props.location.basePath)} />
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

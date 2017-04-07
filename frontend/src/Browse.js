import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom';

// Update for apiPrefix (remove host)
const urlPath = 'https://localhost/api/';
const pathPrefix = '/ui/'

const LogInner = (props) => (
  <h1>Not logged in</h1>
)

const ErrorShower = (props) => (
  <h1>Not found</h1>
)

class EntryList extends Component {

  state = {
    basePath: '',
    entries: null,
    loggedIn: false,
    errorResponse: false
  }

  fetchData() {
    var path = this.props.location.pathname.substring(4)

    fetch(urlPath + path)
      .then(response => {
        if (response.ok) {
          response.json().then(newEntries => {
            this.setState({ basePath: path, entries: newEntries, loggedIn: true, errorResponse: false })
          })
        } else {
          throw new Error(response.status)
        }})
      .catch(err => {
        if (err.message === '404') {
          this.setState({ basePath: path, entries: null, loggedIn: true, errorResponse: true })
        } else {
          this.setState({ basePath: path, entries: null, loggedIn: false, errorResponse: true })
        }
      })
  }

  componentDidMount() {
    this.fetchData()
  }

  login(event) {
    console.log('Logging in')
    event.preventDefault()
    window.open(this.makeHref('https://localhost/login'))
    // fetch('https://localhost/login', { method: 'GET', credentials: 'include' })
    //   .then(res => console.log(res))
  }

  render() {
    if (!this.state.loggedIn) {
      return (
        // <a href='https://localhost/login'>Login</a>
        // <a href='#' onClick={this.login}>Login</a>
        <a href='https://localhost/login' target='_blank' onClick={(event) => {event.preventDefault(); window.open(this.makeHref('https://localhost/login'));}}>login</a>
      )
    }

    if (this.state.errorResponse) {
      return (
        // <ErrorShower />)
        <a href='#' onClick={this.login}>LLLLLLogin</a>
      )
    }

    const entries = this.state.entries
    var basePath = this.state.basePath
    if (basePath.length > 0) {
      basePath += '/'
    }

    return (
      <ul>
      {entries
        ? (this.state.entries.map((entry, index) => (
          <EntryRenderer key={index} entry={entry} base={basePath} updater={this.fetchData}/>
          )))
        : 'Loading..'}
        <il>
          <a href='#' onClick={this.login}>Login</a>
        </il>
      </ul>
    );
  }
}

const EntryRenderer = (props) => (
  <li>
  {props.entry.Directory
    ? <a href={pathPrefix + props.base + props.entry.Name} onClick={props.updater}>{props.entry.Name}</a>
    : <a href={urlPath + props.base +  props.entry.Name}>{props.entry.Name}</a>
  }
  </li>
)

export default class Browse extends Component {
  render() {
    return (
      <Router>
        <Route path="/" component={EntryList} />
      </Router>
    );
  }
}


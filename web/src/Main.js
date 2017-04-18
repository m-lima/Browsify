import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {
  Panel,
  Table,
  Grid,
  Row
} from 'react-bootstrap';

import './Main.css'
import * as Constants from './Constants.js'
import logo from './img/lockHollow.svg';

const DarkStyle = {
  height: '100%',
  paddingTop: 80,
  backgroundColor: '#222222',
  color: 'lightGray'
}

const EntryRenderer = (props) => (
  <tr>
    <td>
      {props.entry.Directory
        ? <span className="glyphicon glyphicon-folder-open" aria-hidden="true"></span>
        : <span />
      }
    </td>
    <td>
      {props.entry.Directory
        ? <Link to={Constants.ui + props.base + props.entry.Name}>{props.entry.Name}</Link>
        : <a href={Constants.api + props.base + props.entry.Name}>{props.entry.Name}</a>
      }
    </td>
    <td>{props.entry.Size}</td>
    <td>{new Date(props.entry.Date).toLocaleString()}</td>
  </tr>
)

export default class Main extends Component {

  state = {
    path: '',
    entries: [],
    loading: false,
    status: null,
    authorized: false
  }

  componentDidMount() {
    this.fetchData(this.props.path)
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.path !== this.props.path || nextProps.authorized !== this.state.authorized) {
      this.fetchData(nextProps.path)
    }
  }

  invalidateData(error) {
    var authorized = (error !== Constants.statusUnauthorized && error !== Constants.statusForbidden)
    this.setState({ entries: [], loading: false, status: error, authorized: authorized } )
    this.props.authUpdater(authorized)
  }

  fetchData(path) {
    this.setState({ path: path, loading: true, authorized: false })
    fetch(Constants.api + path, { method: 'GET', credentials: 'include' })
      .then(response => {
        if (response.ok) {
          response.json().then(newEntries => {
            this.setState({ entries: newEntries, loading: false, status: Constants.statusOK, authorized: true })
            this.props.authUpdater(true)
          })
          .catch(() => this.invalidateData(Constants.statusNotFound))
        } else {
          this.invalidateData(response.status)
        }
      })
      .catch(() => this.invalidateData(Constants.statusNotFound))
  }

  generateBreadcrumb(path) {
    var folders = path.split('/')
    var currentActive

    for (var i = 0; i < folders.length; i++) {
      if (folders[i]) {
        currentActive = i
        if (i === 0) {
          folders[i] = [ '/'+folders[i], Constants.ui + folders[i] ]
        } else {
          folders[i] = [ '/'+folders[i], folders[i-1][1] + '/' + folders[i] ]
        }
      }
    }

    return (
      <span>
        <Link to={Constants.ui}>Home</Link>
        {folders.map((folder, index) => {
          if (folder) {
            if (index === currentActive) {
              return(
                <b key={index}>
                  {folder[0]}
                </b>
              )
            }
            return(
              <Link to={folder[1]} key={index}>
                {folder[0]}
              </Link>
            )
          }
          return ('')
        })}
      </span>
    )
  }

  renderTable() {
    if (this.state.loading) {
      return <b>Loading..</b>
    }

    if (this.state.status === Constants.statusNotFound) {
      return <b>Not found</b>
    }

    if (!this.state.entries || this.state.entries.length === 0) {
      return <b>Empty folder</b>
    }

    return (
      <Table fill>
        <thead>
          <tr>
            <th></th>
            <th>Name</th>
            <th>Size</th>
            <th>Date</th>
          </tr>
        </thead>
        <tbody>
          {this.state.entries.map((entry, index) => (
            <EntryRenderer key={index} entry={entry} base={this.state.path} />
          ))}
        </tbody>
      </Table>
    )
  }

  render() {
    switch(this.state.status) {
      case Constants.statusUnauthorized:
        return (
          <div className='Landing'>
            <img src={logo} className='Landing-logo' alt='logo' />
          </div>
        )
      case Constants.statusForbidden:
        return (
          <div style={DarkStyle}>
            <Grid>
              <h3>Unauthorized</h3>
              <p>The current user does not have access to this content</p>
              <a href={Constants.login}>Retry</a>
            </Grid>
          </div>
        )
      case Constants.statusOK:
      case Constants.statusNotFound:
        return (
          <Grid style={{ paddingTop: 80 }}>
            <Row>
              <Panel header={this.generateBreadcrumb(this.state.path)}>
                {this.renderTable()}
              </Panel>
            </Row>
          </Grid>
        )
      case null:
        return (
          <div style={DarkStyle}/>
        )
      default:
        return (
          <div style={DarkStyle}>
            <Grid>
              <h3>Oops!</h3>
              <p>An error has occurred while processing your request</p>
              <a href={Constants.ui}>Return home</a>
            </Grid>
          </div>
        )
    }
  }
}

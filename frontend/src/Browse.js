import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {
  Panel,
  Table,
  Grid,
  Row
} from 'react-bootstrap';

// Update for apiPrefix (remove host)
const urlPath = 'https://localhost/api/';
const pathPrefix = '/ui/'

const StatusRenderer = (props) => (
  <div>
    {props.status === '403'
      ? <div>
          <h1>Not logged in</h1>
          <form action='https://localhost/login' method='post'>
            <button>Login</button>
          </form>
        </div>
      : props.status === '404'
        ? <h1>Not found</h1>
        : <h1>Error: {props.status}</h1>
    }
  </div>
)

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
        ? <Link to={pathPrefix + props.base + props.entry.Name}>{props.entry.Name}</Link>
        : <a href={urlPath + props.base + props.entry.Name}>{props.entry.Name}</a>
      }
    </td>
    <td>{props.entry.Size}</td>
    <td>{new Date(props.entry.Date).toLocaleString()}</td>
  </tr>
)

export default class BrowseList extends Component {

  state = {
    basePath: '',
    entries: null,
    status: 0
  }

  fetchData(path) {
    console.log(this.props.location)
    this.setState({ entries: null })
    if (path === undefined || path === '') {
      path = this.props.location.pathname
    }

    path = path.substring(4)
    if (path.length > 0) {
      path += '/'
    }

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
      .catch(err => {
        this.setState({ basePath: path, entries: [], status: err.message})
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
    const entries = this.state.entries
    if (!entries) {
      return (<h1>Loading..</h1>)
    } else {
      if (this.state.status) {
        return (<StatusRenderer status={this.state.status} login={this.login} />)
      } else {
        return (
          <Grid>
          <Row>
          <Panel header={this.state.basePath}>
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
              {entries.map((entry, index) => (
                <EntryRenderer key={index} entry={entry} base={this.state.basePath} updater={this.fetchData} caller={this}/>
              ))}
            </tbody>
          </Table>
          </Panel>
          </Row>
          </Grid>
        );
      }
    }
  }
}

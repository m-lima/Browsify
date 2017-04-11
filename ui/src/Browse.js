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

const ContentRenderer = (props) => {
  if (!props.entries) {
    return (
      <b>
        Loading..
      </b>
    )
  }

  if (props.status) {
    return (
      <b>
        {props.status}
      </b>
    )
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
        {props.entries.map((entry, index) => (
          <EntryRenderer key={index} entry={entry} base={props.path} />
        ))}
      </tbody>
    </Table>
  )
}

export default class BrowseList extends Component {

  generateBreadcrumb(path) {
    var folders = path.split('/')
    var currentActive

    for (var i = 0; i < folders.length; i++) {
      if (folders[i]) {
        currentActive = i
        if (i === 0) {
          folders[i] = [ '/'+folders[i], '/ui/'+ folders[i] ]
        } else {
          folders[i] = [ '/'+folders[i], folders[i-1][1] + '/' + folders[i] ]
        }
      }
    }

    return (
      <span>
        <Link to={'/ui'}>Home</Link>
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

  render() {
    return(
      <Grid>
        <Row>
          <Panel header={this.generateBreadcrumb(this.props.basePath)}>
            <ContentRenderer entries={this.props.entries} status={this.props.status} path={this.props.basePath} />
          </Panel>
        </Row>
      </Grid>
    )
  }
}

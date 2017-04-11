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
      <div>
        <Link to={'/ui'}>Home</Link>
        {folders.map((folder, index) => {
          if (folder) {
            if (index === currentActive) {
              return(
                <b>
                {folder[0]}
                </b>
              )
            }
            return(
              <Link to={folder[1]}>
                {folder[0]}
              </Link>
            )
          }
          return ('')
        })}
      </div>
    )

    // return (
    //   <Breadcrumb>
    //     <Breadcrumb.Item href={folders[0][1]}>
    //       Home
    //     </Breadcrumb.Item>
    //     {folders.map((folder, index) => {
    //       if (folder) {
    //         if (index === currentActive) {
    //           return(
    //             <Breadcrumb.Item active>
    //               {folder[0]}
    //             </Breadcrumb.Item>
    //           )
    //         }
    //         return(
    //           <Breadcrumb.Item href={folder[1]}>
    //             {folder[0]}
    //           </Breadcrumb.Item>
    //         )
    //       }
    //       return ('')
    //     })}
    //   </Breadcrumb>
    // )
  }

  render() {
    const entries = this.props.entries
    if (!entries) {
      return (
        <Grid>
          <Row>
            <Panel header={this.generateBreadcrumb(this.props.basePath)}>
              Loading..
            </Panel>
          </Row>
        </Grid>
      )
    } else {
      if (this.props.status) {
        return (
          <Grid>
            <Row>
              <Panel header={this.props.status} />
            </Row>
          </Grid>
        )
      } else {
        return (
          <Grid>
            <Row>
              <Panel header={this.generateBreadcrumb(this.props.basePath)}>
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
                      <EntryRenderer key={index} entry={entry} base={this.props.basePath} />
                    ))}
                  </tbody>
                </Table>
              </Panel>
            </Row>
          </Grid>
        )
      }
    }
  }
}

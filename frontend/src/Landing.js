import React, { Component } from 'react';
// import { Button } from 'react-bootstrap';
import logo from './img/lockHollow.svg';
import './Landing.css'

export default class Landing extends Component {

  render() {
    return (
      <div className="Landing">
        <div className="Landing-bundle">
          <img src={logo} className="Landing-logo" alt="logo" />
          <h1 className="Landing-intro">SecuriDash</h1>
          <Button bsStyle="info">Login</Button>
        </div>
        <div className="Landing-stripe" />
      </div>
    );
  }
}

import React, { Component } from 'react';
import { Button } from 'react-bootstrap';
import logo from './img/lockHollow.svg';
// import './Landing.css'

export default class Landing extends Component {

  render() {
    return (
      <div className='Landing'>
        <img src={logo} className='Landing-logo' alt='logo' />
        <h1 />
        <form action='https://localhost/login' method='post'>
          <Button bsStyle='info' bsSize='lg' type='submit'>Login</Button>
        </form>
      </div>
    );
  }
}

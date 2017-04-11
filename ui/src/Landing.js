import React, { Component } from 'react';
import logo from './img/lockHollow.svg';
import './Landing.css'

export default class Landing extends Component {

  render() {
    return (
      <div className='Landing'>
        <img src={logo} className='Landing-logo' alt='logo' />
      </div>
    );
  }
}

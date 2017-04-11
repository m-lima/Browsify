import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom';

import Title from './Title.js'
import BrowseList from './Browse.js'
import Landing from './Landing.js'
import './App.css';

class App extends Component {

  render() {
    return (
      <Router>
        <div>
          <Title />
          <div className='App-container'>
            <Route path="/" component={BrowseList} />
          </div>
        </div>
      </Router>
    )
  }
}

export default App

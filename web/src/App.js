import React, { Component } from 'react';
import logo from './logo.svg';
import Player from './player';
import './styles/css/App.css';

class App extends Component {
	render() {
		return (
			<div className="App">
				<header className="App-header">
					<img src={logo} className="App-logo" alt="logo" />
					<h1 className="App-title">Welcome to House Party</h1>
				</header>
				<p className="App-intro">
					Current Song
				</p>
				<Player/>
			</div>
		);
	}
}

export default App;

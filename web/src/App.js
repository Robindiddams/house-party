import React, { Component } from 'react';
import Sockette from 'sockette';
import logo from './logo.svg';
import { PlaybackControls, PauseButton } from 'react-player-controls'
import './styles/css/styles.css';

class App extends Component {

	constructor() {
		super()
		this.state = {
			songs: [],
			songNumber:1,
			currentSong:"thebeatles",
			paused:false
		};
		this.ws = new Sockette('ws://localhost:8080', {
			timeout: 5e3,
			maxAttempts: 5,
			onopen: e => console.log('Connected!', e),
			onmessage: e => console.log('Received:', e),
			onreconnect: e => console.log('Reconnecting...', e),
			onclose: e => console.log('Closed!', e),
			onerror: e => console.log('Error:', e)
		});
	}

	logHello() {
		console.log("hey");
	}
	

	render() {

		return (
			<div className="site-container">
				<header className="site-header">
					<h1 className="site-title-text">Welcome to House Party</h1>
				</header>
				<div className="site-body">
					<div className="song-block">
						<div className="song-text">
							<h1>{this.state.currentSong}</h1>
						</div>
						<PlaybackControls
							isPlayable={true}
							isPlaying={this.state.paused}
							onPlaybackChange={() => {
								if (this.state.paused) {
									this.setState({paused:false})
									this.ws.send("hello")
								} else {
									this.setState({paused:true})
								}
							}}
							showPrevious={true}
							hasPrevious={this.state.songNumber > 0}
							onPrevious={this.logHello}
							showNext={true}
							hasNext={this.state.songNumber < 5}
							onNext={this.logHello}
							/>
						</div>
				</div>
			</div>
		);
	}
}

export default App;

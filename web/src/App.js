import React, { Component } from 'react';
import Sockette from 'sockette';
import logo from './logo.svg';
import { PlaybackControls, PauseButton } from 'react-player-controls'
import './styles/css/styles.css';

class App extends Component {

	constructor(props) {
		super(props)
		this.state = {
			songs: [],
			songNumber:1,
			currentSong:"thebeatles",
			playing:true
		};
		this.ws = new Sockette('ws://localhost:8080', {
			timeout: 5e3,
			maxAttempts: 5,
			onopen: e => console.log('Connected!', e),
			onmessage: e => {
				// console.log("Recieved!", e)
				let obj = JSON.parse(e.data)
				console.log("Recieved!", obj);
				this.setState({playing:obj.playing})
				obj.title && this.setState({currentSong:obj.title})
				console.log("state!", this.state);
			},
			onreconnect: e => console.log('Reconnecting...', e),
			onclose: e => console.log('Closed!', e),
			onerror: e => console.log('Error:', e)
		});
	}

	logHello() {
		console.log("hey");
	}

	sendAction(action) {
		this.ws.json({action});
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
							isPlaying={this.state.playing}
							onPlaybackChange={() => {
								if (this.state.playing) {
									// this.setState({playing:false});
									this.sendAction("pause")
								} else {
									// this.setState({playing:true});
									this.sendAction("play")
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

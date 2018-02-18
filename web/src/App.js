import React, { Component } from 'react';
import Sockette from 'sockette';
import { PlaybackControls } from 'react-player-controls'
import './styles/css/styles.css';

class App extends Component {

	constructor(props) {
		super(props)
		const s = 'ws://' + window.location.hostname + ':8080/sock' //use for dev
		// const s = 'ws://' + window.location.host + 'sock' //use for prod
		this.state = {
			songs: [],
			hasNext: true,
			currentSong:"thebeatles",
			playing:true,
			url:"",
			connected:false
		};
		this.ws = new Sockette(s, {
			timeout: 5e3,
			maxAttempts: 5,
			onopen: e => console.log('Connected!', e),
			onmessage: e => {
				// console.log("Recieved!", e)
				let obj = JSON.parse(e.data)
				console.log("Recieved!", obj);
				this.setState({playing:obj.playing})
				obj.title && this.setState({currentSong:obj.title})
				// console.log("state!", this.state);
			},
			onreconnect: e => {
				console.log('Reconnecting...', e)
			},
			onclose: e => {
				console.log('Closed!', e)
				this.setState({connected:false})
			},
			onerror: e => {
				// console.log('Error:', e)
			}
		});

		this.sendAction = this.sendAction.bind(this);
		this.handleChange = this.handleChange.bind(this);
	}


	sendAction(action) {
		if (this.state.connected) {
			this.ws.json({action});
		}
	}

	handleChange(event) {
		const url = event.target.value
		this.setState({url});
		if (this.ValidURL(url)) {
			console.log("valid");
			//do login here to show what song it is
			this.ws.json({action:"queue", meta: url})
			this.setState({url:""})
		}
	}

	ValidURL(url) {
		if (url !== undefined || url !== '') {
			var regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=|\?v=)([^#\&\?]*).*/;
			var match = url.match(regExp);
			if (match && match[2].length === 11) {
				return true
			}
			else {
				return false
			}
		}
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
						<div className="stack-center">
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
								showPrevious={false}
								showNext={true}
								hasNext={this.state.hasNext}
								onNext={() => {
									this.sendAction("next")
								}}
								/>
							</div>
						<div className="url-paste-box">
							<input type="text" className="Input-text" value={this.state.url} onChange={this.handleChange} placeholder="paste a url here"/>
						</div>
					</div>
				</div>
			</div>
		);
	}
}

export default App;

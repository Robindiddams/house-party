import React, { Component } from 'react';
import Sockette from 'sockette';
import { PlaybackControls } from 'react-player-controls'
import './styles/css/styles.css';

class App extends Component {

	constructor(props) {
		super(props)
		this.state = {
			songs: [],
			hasNext: true,
			currentSong:"thebeatles",
			playing:true,
			url:"",
			connected:false,
			started:false,
		};
		console.log(window.location)
		this.sendJson = this.sendJson.bind(this);
		this.handleChange = this.handleChange.bind(this);
		this.playbackControl = this.playbackControl.bind(this);
	}

	componentDidMount() {
		setInterval( () => {
			this.sendJson({action:'ping'})
			.then(resp => {
				this.setState({
					playing:resp.playing,
					currentSong:resp.title,
				});
				if (resp.playing) {
					this.setState({started: true});
				}
			});
		},3000)
	}

	sendJson(body) {
		const apiurl = '/api/action';
		return fetch( apiurl, {
			body: JSON.stringify(body),
			cache: 'no-cache',
			headers: {
				'content-type': 'application/json'
			},
			method: 'POST',
			redirect: 'follow',
		})
		.then(response => {
			return response.json();
		})
		.then(myJson => {
			return myJson
		})
		.catch(e => {
			console.log(e);
		});
	}

	playbackControl() {
		let playstr = 'play';
		if (this.state.playing) {
			playstr = 'pause';
		}
		this.setState({playing:!this.state.playing})
		this.sendJson({action:playstr}).then(resp => {
			console.log(resp);
			if (resp !== undefined) {
				this.setState({
					playing:resp.playing,
					currentSong:resp.title,
				});
			}
		});
	}

	handleChange(event) {
		const url = event.target.value
		this.setState({url});
		if (this.ValidURL(url)) {
			console.log("valid");
			//do login here to show what song it is
			this.sendJson({action:"queue", meta: url})
			.then(resp => {
				console.log("resp:", resp)
				this.setState({url:""})
			});
		}
	}

	// FIXME: This sort of works :/
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
								isPlayable={this.state.started}
								isPlaying={this.state.playing}
								onPlaybackChange={this.playbackControl}
								showPrevious={false}
								showNext={true}
								hasNext={this.state.hasNext}
								onNext={() => {
									this.sendJson({action:"next"});
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

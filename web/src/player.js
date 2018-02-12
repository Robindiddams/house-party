import React, { Component } from 'react';
import { PlayButton, PlaybackControls } from 'react-player-controls'
import './styles/css/player.css';

class Player extends Component {
	render() {
		let currentSong = 1

		return (
			<div>
				<PlaybackControls
					isPlayable={true}
					isPlaying={false}
					onPlaybackChange={console.log("hey")}
					showPrevious={true}
					hasPrevious={currentSong > 0}
					onPrevious={currentSong--}
					showNext={true}
					hasNext={currentSong < 5}
					onNext={currentSong++}
				/>
			</div>
		);
	}
}

export default Player;

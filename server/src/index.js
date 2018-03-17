// const express = require('express');
// const path = require('path');
import path from 'path';
import express from 'express';
import bodyParser from 'body-parser';
import mpv from 'node-mpv';

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
let mpvPlayer = new mpv({
	'audio_only': true,
});

let playing = false;

// Serve static files from the React app
app.use(express.static(path.join(__dirname, '../../web/build')));

// Put all API endpoints under '/api'
app.post('/api/action', (req, res) => {
	console.log(`got this ${JSON.stringify(req.body)}`);
	// req.body.action 
	let state = {};
	switch (req.body.action) {
	case 'play': 
		if (!playing) {
			console.log('playing');
			mpvPlayer.play();
			playing = true;
		}
		break;
	case 'pause': 
		if (playing) {
			console.log('pausing');
			mpvPlayer.pause();
			playing = false;
		}
		break;
	case 'queue': 
		console.log('queueing');
		mpvPlayer.load(req.body.meta);
		playing = true;
		break;
	default : 
		console.log('i hope this works');
		break;
	}
	state.playing = playing;
	res.status(200).send(state);	
});

// The 'catchall' handler: for any request that doesn't
// match one above, send back React's index.html file.
app.get('*', (req, res) => {
	res.sendFile(path.join(__dirname+'../../web/build/index.html'));
});

const port = process.env.PORT || 5000;
app.listen(port);

console.log(`Password generator listening on ${port}`);
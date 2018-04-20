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

let status = {
	playing: false,
	title:null,
	length:0,
	volume:0,
};

// Serve static files from the React app
app.use(express.static(path.join(__dirname, '../../web/build')));

mpvPlayer.on('statuschange', s => {
	// console.log(s);
	status.title = s['media-title'];
	status.length = s.duration;
	status.volume = s.volume;
});
// Put all API endpoints under '/api'
app.post('/api/action', (req, res) => {
	console.log(`got this ${JSON.stringify(req.body)}\n${req.host}`);
	// req.body.action 
	// let state = {};
	switch (req.body.action) {
	case 'play': 
		if (!status.playing) {
			console.log('playing');
			mpvPlayer.play();
			status.playing = true;
		}
		break;
	case 'pause': 
		if (status.playing) {
			console.log('pausing');
			mpvPlayer.pause();
			status.playing = false;
		}
		break;
	case 'queue': 
		console.log('queueing');
		mpvPlayer.append(req.body.meta, 'append-play');
		status.playing = true;
		break;
	case 'next': 
		console.log('skipping');
		mpvPlayer.next();
		status.playing = true;
		break;
	default : 
		// console.log('ping?');
		break;
	}
	res.status(200).send(status);	
});

// The 'catchall' handler: for any request that doesn't
// match one above, send back React's index.html file.
app.get('*', (req, res) => {
	res.sendFile(path.join(__dirname+'../../web/build/index.html'));
});

const port = process.env.PORT || 5000;
app.listen(port);

console.log(`React server listening on ${port}`);
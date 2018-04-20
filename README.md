# House Party

It's a front end for mpv player that accepts youtube urls and plays them. I wanted something that i could throw on a raspberry pi and plug into my speakers and then queue up songs. It works, but I havent implemented a way to see whats queued next.

### Build instructions
You'll need mpv player
```bash
# start "server"
cd server
npm install
npm start

# react dev server
cd web/
npm install
npm run sass #compile sass
npm start
```

### Contributing
PRs welcome!
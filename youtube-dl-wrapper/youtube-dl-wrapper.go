package wrapper

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func checkDep() error {
	// Make sure we have both ffmpeg and youtube-dl
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("missing ffmpeg")
	}
	_, err = exec.LookPath("youtube-dl")
	if err != nil {
		return fmt.Errorf("missing youtube-dl")
	}
	return nil
}

// DownloadMp3 will download an video from youtube and extract
// just the audio, convert it to mp3 and return the file
func DownloadMp3(url string) (*os.File, error) {
	err := checkDep()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("youtube-dl", "--prefer-ffmpeg", "-x", "--audio-format", "mp3", url)
	b := bytes.Buffer{}
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// fmt.Println(b.String())
		return nil, fmt.Errorf("Failed to download/extract: %v", err)
	}

	// extract filename
	re := regexp.MustCompile(`\[ffmpeg\]\sDestination:\s(.*)\.mp3\s`)
	match := re.FindStringSubmatch(b.String())
	if len(match) == 2 {
		//TODO: find a way to get the full filepath
		name := match[1] + ".mp3"
		return os.Open(name)
	}
	return nil, fmt.Errorf("Unable to locate downloaded file")
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shkh/lastfm-go/lastfm"
)

const FIRST_TRACK_LINE_INDEX = 3

type Track struct {
	artist    string
	album     string
	title     string
	timestamp string
}

func main() {
	logPath := flag.String("f", "", "Path to .scrobbler.log")
	offset := flag.String("o", "0h", "Time difference from UTC (format +10h or -10.5h")
	nonInteractive := flag.String("n", "", "Non Interactive Mode: Automatically (\"keep\", \"delete\" or \"delete-on-success\") at end of program")
	auth := flag.Bool("auth", false, "First Time Authentication")
	flag.Parse()

	api := lastfm.New(API_KEY, API_SECRET)

	/* First time Authentication */
	if *auth {
		/* https://www.last.fm/api/desktopauth */

		/* Create folder to store session key */
		err := os.Mkdir(getConfigDir(), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		/* "Step 2" */
		token, err := api.GetToken()
		if err != nil {
			log.Fatal(err)
		}

		/* Step 3 */

		authURL := api.GetAuthTokenUrl(token)
		fmt.Printf("Go to %q, allow access and press ENTER\n", authURL)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')

		/* "Step 4" */
		loginErr := api.LoginWithToken(token)
		if loginErr != nil {
			log.Fatal(loginErr)
		}

		sessionKey := api.GetSessionKey()
		/* Save session key in config dir/rb-scrobbler */
		os.WriteFile(getKeyFilePath(), []byte(sessionKey), os.ModePerm)

	}

	/* When given a file, start executing here */

	if *logPath != "" {

		var tracks []Track
		var success uint
		var fails uint

		scrobblerLog, err := importLog(logPath)
		if err != nil {
			log.Fatal(err)
		}

		/* length -1 since you go out of bounds otherwise. Only iterate from where tracks
		actually show up */
		for i := FIRST_TRACK_LINE_INDEX; i < len(scrobblerLog)-1; i++ {
			newTrack, listened := logLineToTrack(scrobblerLog[i], *offset)
			if listened == true {
				tracks = append(tracks, newTrack)
			}
		}

		/* Login here, after tracks have been parsed and are ready to send */
		sessionKey, err := getSavedKey()
		if err != nil {
			log.Fatal(err)
		}
		api.SetSession(sessionKey)

		for _, track := range tracks {
			p := lastfm.P{"artist": track.artist, "album": track.album, "track": track.title, "timestamp": track.timestamp}

			_, err := api.Track.Scrobble(p)
			if err != nil {
				fmt.Printf("[FAIL] %s - %s\n", track.artist, track.title)
				fails++
			} else {
				fmt.Printf("[OK] %s - %s\n", track.artist, track.title)
				success++
			}
		}

		fmt.Printf("\nFinished: %d tracks scrobbled, %d failed, %d total\n", success, fails, len(tracks))

		/* Handling of file (manual/non interactive delete/keep) */

		switch *nonInteractive {
		case "keep":
			fmt.Printf("%q kept\n", *logPath)

		case "delete":
			deleteLogFile(logPath)

		case "delete-on-success":
			if fails == 0 {
				deleteLogFile(logPath)
			} else {
				fmt.Printf("Scrobble failures: %q not deleted.", *logPath)
				os.Exit(1)
			}

		default:
			reader := bufio.NewReader(os.Stdin)
			var input string
			fmt.Printf("Delete %q? [y/n] ", *logPath)
			input, err := reader.ReadString('\n')
			fmt.Print("\n")
			if err != nil {
				fmt.Printf("Error reading input! File %q not deleted.\n%v\n", *logPath, err)
			} else if strings.ContainsAny(input, "y") || strings.ContainsAny(input, "Y") {
				deleteLogFile(logPath)
			} else {
				fmt.Printf("%q kept.\n", *logPath)
			}
		}

	} else {
		fmt.Println("File (-f) cannot be empty!")
		os.Exit(1)
	}

}

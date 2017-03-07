package statutils

import (
	"fmt"
	"net/url"
	"time"

	"github.com/sauron/detectors"
	"github.com/sauron/session"
)

//StartSessionsStatBeholder collects and outputs statistics on sessions
func StartSessionsStatBeholder(periodSec int, sessions *sstrg.SessionsTable, detector *detectors.Detector) {
	// fire once per second
	t := time.NewTicker(time.Second * time.Duration(periodSec))

	for {
		//calcDurationStat(sessions)
		calcBotsStat(sessions, detector)
		<-t.C
	}
}

func calcBotsStat(sessions *sstrg.SessionsTable, detector *detectors.Detector) {
	fmt.Printf("\n\n\n----Active bots session ------\n")

	keysByLabel := map[string][]string{
		"human":      []string{},
		"irrelevant": []string{},
		"bot":        []string{},
		"unknown":    []string{},
		"humanlike":  []string{},
	}
	total := 0
	sessions.Lock()

	//TODO: do not iterate over map
	for key, s := range sessions.H {
		label := (*detector).GetLabel(s)

		switch label {
		case detectors.BotLabel:
			keysByLabel["bot"] = append(keysByLabel["bot"], key)
		case detectors.HumanLabel:
			keysByLabel["human"] = append(keysByLabel["human"], key)
		case detectors.UnknownLabel:
			keysByLabel["unknown"] = append(keysByLabel["unknown"], key)
		case detectors.IrrelevantLabel:
			keysByLabel["irrelevant"] = append(keysByLabel["irrelevant"], key)
		case 4:
			keysByLabel["humanlike"] = append(keysByLabel["humanlike"], key)
		}

		total++
	}

	sessions.Unlock()

	labelsToPrint := []string{"bot"}

	for _, label := range labelsToPrint {
		fmt.Printf("%s:\n\n", label)

		for _, sessionKey := range keysByLabel[label] {
			fmt.Printf("http://localhost:3000/raw?key=%s\n\n", url.QueryEscape(sessionKey))
		}
	}

	for _, label := range labelsToPrint {
		fmt.Printf("%s:\n\n", label)

		for _, sessionKey := range keysByLabel[label] {
			fmt.Printf("http://localhost:3000/raw?key=%s\n\n", url.QueryEscape(sessionKey))
		}
	}

	fmt.Printf("----\n")

	for label := range keysByLabel {
		fmt.Printf("%s: %d\n", label, len(keysByLabel[label]))
	}

	fmt.Printf("----\nTotal: %d\n", total)
}

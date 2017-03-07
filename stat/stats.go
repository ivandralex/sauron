package statutils

import (
	"fmt"
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
	fmt.Println("----Active bots session ------")

	var botsKeys = []string{}
	counts := struct {
		bots       int
		human      int
		uknown     int
		irrelevant int
	}{}
	sessions.Lock()

	//TODO: do not iterate over map
	for key, s := range sessions.H {
		label := (*detector).GetLabel(s)

		switch label {
		case detectors.BotLabel:
			fmt.Println(key)
			botsKeys = append(botsKeys, key)
			counts.bots++
		case detectors.HumanLabel:
			counts.human++
		case detectors.UnknownLabel:
			counts.uknown++
		case detectors.IrrelevantLabel:
			counts.irrelevant++
		}
	}

	sessions.Unlock()

	fmt.Printf("----\nCounts: %v\n\n", counts)
}

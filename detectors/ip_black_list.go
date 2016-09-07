package detectors

import (
	"fmt"
	"strings"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//BlackListDetector detects bot by checking if client has ip from black list
type BlackListDetector struct {
	blackListedIPs map[string]bool
}

//Init initializes human path detector
func (d *BlackListDetector) Init(configPath string) {
	d.blackListedIPs = configutil.ReadStringMap(configPath)
}

//GetLabel returns label for session by checking
func (d *BlackListDetector) GetLabel(s *sstrg.SessionData) int {
	//Convert IP to mask
	parts := strings.Split(s.IP, ".")

	if len(parts) != 4 {
		fmt.Printf("Incorrect IP: %s\n", s.IP)
		return UnknownLabel
	}

	parts[2] = "*"
	parts[3] = "*"
	mask := strings.Join(parts, ".")

	if _, ok := d.blackListedIPs[mask]; ok {
		return BotLabel
	}

	return UnknownLabel
}

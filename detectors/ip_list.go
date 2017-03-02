package detectors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//IPListDetector detects bot by checking if client has ip from black list
type IPListDetector struct {
	listedIPS map[string]bool
	label     int
}

//Init initializes human path detector
func (d *IPListDetector) Init(configPath string) {
	d.listedIPS = configutil.ReadStringMap(configPath)
}

//SetLabel sets positive label for this detector
func (d *IPListDetector) SetLabel(label int) {
	d.label = label
}

//GetLabel returns label for session by checking
func (d *IPListDetector) GetLabel(s *sstrg.SessionData) int {
	if _, ok := d.listedIPS[s.IP]; ok {
		return d.label
	}

	mask, err := ipToMask(s.IP)

	if err != nil {
		fmt.Printf("Incorrect IP: %s\n", s.IP)
	}

	if _, ok := d.listedIPS[mask]; ok {
		return d.label
	}

	return UnknownLabel
}

func ipToMask(ip string) (string, error) {
	//Convert IP to mask
	parts := strings.Split(ip, ".")

	if len(parts) != 4 {
		return "", errors.New("Incorrect IP")
	}

	parts[2] = "*"
	parts[3] = "*"
	mask := strings.Join(parts, ".")

	return mask, nil
}

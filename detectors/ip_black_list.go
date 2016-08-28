package detectors

import (
	"fmt"
	"strings"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//TODO: it's not a very good idea to use global var here
var blackListedIPs map[string]bool

//TODO: check guidelines if it's a good idea to read config here rather than get it as argument in GetLabel
func init() {
	blackListedIPs = configutil.ReadStringMap("../configs/ip_black_list.csv")
}

//GetLabelByBlackList returns label for session by checking
func GetLabelByBlackList(s *sstrg.SessionData) string {
	//Convert IP to mask
	parts := strings.Split(s.IP, ".")

	if len(parts) != 4 {
		fmt.Printf("Incorrect IP: %s\n", s.IP)
		return "0"
	}

	parts[2] = "*"
	parts[3] = "*"
	mask := strings.Join(parts, ".")

	if _, ok := blackListedIPs[mask]; ok {
		fmt.Printf("Found bot %s\n", s.IP)
		return "1"
	}

	return "0"
}

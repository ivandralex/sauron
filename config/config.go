package configutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//ReadPathsConfig config containing HTTP paths
func ReadPathsConfig(fileName string) []string {
	f, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error readin top paths ", err)
		os.Exit(1)
	}

	var targetPaths []string

	r := bufio.NewReader(f)
	for {
		str, err := r.ReadString(10)
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		//Remove trailing slash
		if strings.HasSuffix(str, "\n") {
			str = str[:len(str)-1]
		}
		if strings.HasSuffix(str, "/") {
			str = str[:len(str)-1]
		}

		targetPaths = append(targetPaths, str)
	}

	f.Close()

	return targetPaths
}

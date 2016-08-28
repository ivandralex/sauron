package configutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//ReadPathsConfig config containing HTTP paths
func ReadPathsConfig(filePath string) []string {
	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error reading %s %v\n", filePath, err)
		os.Exit(1)
	}

	var paths []string

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

		paths = append(paths, str)
	}

	f.Close()

	return paths
}

//ReadStringMap config containing HTTP paths
func ReadStringMap(filePath string) map[string]bool {
	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error reading %s %v\n", filePath, err)
		os.Exit(1)
	}

	var stringMap = make(map[string]bool)

	r := bufio.NewReader(f)
	for {
		str, err := r.ReadString(10)
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		if strings.HasSuffix(str, "\n") {
			str = str[:len(str)-1]
		}

		stringMap[str] = true
	}

	f.Close()

	return stringMap
}

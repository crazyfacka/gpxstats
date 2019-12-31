package main

import (
	"fmt"
	"os"

	"github.com/crazyfacka/gpxstats/tools"

	"github.com/tkrajina/gpxgo/gpx"
)

func main() {
	var file *os.File
	var gpxFiles []*gpx.GPX
	var data []byte
	var err error

	bytesRead := -1
	buffer := make([]byte, 100)

	gpxFilesToRead := os.Args[1:]
	if len(gpxFilesToRead) < 1 {
		fmt.Println("You must specify at least one GPX file to parse")
		os.Exit(-1)
	}

	gpxFiles = make([]*gpx.GPX, len(gpxFilesToRead))

	for i, curFileToRead := range gpxFilesToRead {
		file, err = os.Open(curFileToRead)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err.Error())
			os.Exit(-1)
		}

		for bytesRead != 0 {
			bytesRead, err = file.Read(buffer)
			data = append(data, buffer[:bytesRead]...)
		}

		gpxFiles[i], err = gpx.ParseBytes(data)
		if err != nil {
			fmt.Printf("Error parsing GPX: %s\n", err.Error())
			os.Exit(-1)
		}

		file.Close()

		bytesRead = -1
		buffer = make([]byte, 100)
		data = data[:0]
	}

	/* Print stats */
	if len(gpxFiles) == 1 {
		tools.PrintSingleStats(gpxFiles[0], os.Args[1])
	} else {
		tools.PrintCombinedStats(gpxFiles, os.Args[1:])
	}
}

package make

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"yapm/logger"

	"github.com/fatih/color"
)

func Make(args []string) {
	path := args[0]

	fmt.Println()
	logger.LogString(path, "Making", color.FgBlue, color.Bold)

	logger.Info(logger.CreateLogEntry("Extracting configs"))
	configs := readTarConfigs(path)

	logger.Info(logger.CreateLogEntry("Reading hashes.csv"))
	fmt.Println()
	for k, v := range configs.Hashes {
		logger.LogRaw(strings.Trim(k, " "), color.FgMagenta, color.Bold)
		fmt.Print(": ")
		logger.LogRawln(strings.Trim(v, " "), color.FgBlue, color.Bold)
	}

	fmt.Println()

	logger.Info(logger.CreateLogEntry("Reading location definitions"))
	for k, v := range configs.Definitions {
		logger.LogRaw(strings.Trim(k, " "), color.FgMagenta, color.Bold)
		fmt.Print(": ")
		logger.LogRawln(strings.Trim(v, " "), color.FgBlue, color.Bold)
	}

	fmt.Println()
}

func readTarConfigs(path string) ConfigBundle {
	var ret ConfigBundle = ConfigBundle{}

	fle, err := os.Open(path)
	if err != nil {
		logger.Error(logger.CreateLogEntry(err.Error()))
	}

	tread := tar.NewReader(fle)
	for {
		header, err := tread.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(logger.CreateLogEntry(err.Error()))
		}

		if header.Name == "definitions.json" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, tread)
			if err != nil {
				logger.Error(logger.CreateLogEntry(err.Error()))
			}

			json.Unmarshal(buf.Bytes(), &ret.Definitions)
		}

		if header.Name == "hashes.csv" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, tread)
			if err != nil {
				logger.Error(logger.CreateLogEntry(err.Error()))
			}

			ret.Hashes = readTarCsv(buf.String())
		}
	}

	return ret
}

func readTarCsv(content string) map[string]string {
	ret := map[string]string{}
	lines := []string{}

	cur := ""
	for _, c := range content {
		if c == '\n' {
			lines = append(lines, cur)
			cur = ""

			continue
		}

		cur = cur + string(c)
	}

	for _, str := range lines {
		fields := strings.Split(str, ",")
		ret[fields[0]] = fields[1]
	}

	return ret
}

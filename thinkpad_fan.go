package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.thinkpad-fan")

const (
	pathProcFan = "/proc/acpi/ibm/fan"
)

type TPFanPlugin struct{}

func (p TPFanPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"thinkpad.fan": mp.Graphs{
			Label: "Fan Speed",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				mp.Metrics{Name: "speed", Label: "rpm"},
			},
		},
	}
}

func (p TPFanPlugin) FetchMetrics() (map[string]interface{}, error) {
	file, err := os.Open(pathProcFan)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]interface{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parts := strings.Split(string(scanner.Text()), ":")
		if len(parts) < 1 {
			continue
		}

		if parts[0] == "speed" {
			speed, err := strconv.ParseUint(strings.TrimLeft(parts[1], "\t"), 10, 32)
			if err != nil {
				logger.Warningf("ParseInt: %v", err)
			}
			result["speed"] = speed
		}
	}

	return result, nil
}

func doMain(c *cli.Context) error {
	tpf := TPFanPlugin{}
	helper := mp.NewMackerelPlugin(tpf)
	helper.Run()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "mackerel-plugin-thinkpad-fan"
	app.Version = version
	app.Usage = "Get ThinkPad Fan rpm metrics."
	app.Author = "KOJIMA Kazunori"
	app.Email = "kjm.kznr@gmail.com"
	app.Action = doMain

	app.Run(os.Args)
}

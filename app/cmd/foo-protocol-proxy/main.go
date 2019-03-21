package main

import (
	"flag"
	"fmt"
	"github.com/ahmedkamals/foo-protocol-proxy/app"
	"github.com/ahmedkamals/foo-protocol-proxy/app/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/app/config"
	"github.com/ahmedkamals/foo-protocol-proxy/app/persistence"
	"github.com/kpango/glg"
)

func main() {
	customLevel := "Info"
	logger := glg.New()
	logger.SetMode(glg.STD).AddStdLevel(customLevel, glg.STD, false)
	fmt.Println(glg.Cyan(
		`
 _______             ______                                _     ______
(_______)           (_____ \           _                  | |   (_____ \
 _____ ___   ___     _____) )___ ___ _| |_ ___   ____ ___ | |    _____) )___ ___ _   _ _   _
|  ___) _ \ / _ \   |  ____/ ___) _ (_   _) _ \ / ___) _ \| |   |  ____/ ___) _ ( \ / ) | | |
| |  | |_| | |_| |  | |   | |  | |_| || || |_| ( (__| |_| | |   | |   | |  | |_| ) X (| |_| |
|_|   \___/ \___/   |_|   |_|   \___/  \__)___/ \____)___/ \_)  |_|   |_|   \___(_/ \_)\__  |
                                                                                      (____/ `,
	))
	config := parseConfig()
	saver, err := persistence.NewSaver(config.RecoveryPath)
	if err != nil {
		logger.Fatal(err)
	}

	dispatcher := app.NewDispatcher(config, analysis.NewAnalyzer(), saver)
	dispatcher.Start()
}

func parseConfig() config.Configuration {
	var (
		listen        = flag.String("listen", ":8002", "Listening port.")
		forward       = flag.String("forward", ":8001", "Forwarding port.")
		healthAddress = flag.String("health", "0.0.0.0:8081", "Health service address.")
		httpAddr      = flag.String("http", "0.0.0.0:8081", "HTTP service address.")
		recoveryPath  = flag.String("recovery-path", "../.data/recovery.json", "Recovery path.")
	)
	flag.Parse()

	return config.Configuration{
		Listening:     *listen,
		Forwarding:    *forward,
		HealthAddress: *healthAddress,
		HTTPAddress:   *httpAddr,
		RecoveryPath:  *recoveryPath,
	}
}

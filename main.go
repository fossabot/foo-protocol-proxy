package main

import (
	"flag"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/app"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/persistence"
	"github.com/ahmedkamals/foo-protocol-proxy/utils"
	"os"
)

func main() {
	colorable := utils.NewColorable(os.Stdout)
	println(colorable.Wrap(`
 _______             ______                                _     ______
(_______)           (_____ \           _                  | |   (_____ \
 _____ ___   ___     _____) )___ ___ _| |_ ___   ____ ___ | |    _____) )___ ___ _   _ _   _
|  ___) _ \ / _ \   |  ____/ ___) _ (_   _) _ \ / ___) _ \| |   |  ____/ ___) _ ( \ / ) | | |
| |  | |_| | |_| |  | |   | |  | |_| || || |_| ( (__| |_| | |   | |   | |  | |_| ) X (| |_| |
|_|   \___/ \___/   |_|   |_|   \___/  \__)___/ \____)___/ \_)  |_|   |_|   \___(_/ \_)\__  |
                                                                                      (____/ `,
		utils.FGBlue,
	))
	config := parseConfig()
	dispatcher := app.NewDispatcher(config, analysis.NewAnalyzer(), persistence.NewSaver(config.RecoveryPath))
	dispatcher.Run()
}

func parseConfig() config.Configuration {
	var (
		listen       = flag.String("listen", ":8002", "Listening port.")
		forward      = flag.String("forward", ":8001", "Forwarding port.")
		httpAddr     = flag.String("http", "0.0.0.0:8088", "Health service address.")
		recoveryPath = flag.String("recovery-path", "data/recovery.json", "Recovery path.")
	)
	flag.Parse()

	return config.Configuration{
		Listening:    *listen,
		Forwarding:   *forward,
		HTTPAddress:  *httpAddr,
		RecoveryPath: *recoveryPath,
	}
}

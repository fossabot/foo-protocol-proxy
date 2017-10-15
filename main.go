package main

import (
	"github.com/ahmedkamals/foo-protocol-proxy/app"
	"github.com/ahmedkamals/foo-protocol-proxy/utils"
	"os"
)

func main() {
	colorable := utils.NewColorable(os.Stdout)
	println(colorable.Wrap(
		[]string{`
  _____             ____            _                  _   ____
 |  ___|__   ___   |  _ \ _ __ ___ | |_ ___   ___ ___ | | |  _ \ _ __ _____  ___   _
 | |_ / _ \ / _ \  | |_) | '__/ _ \| __/ _ \ / __/ _ \| | | |_) | '__/ _ \ \/ / | | |
 |  _| (_) | (_) | |  __/| | | (_) | || (_) | (_| (_) | | |  __/| | | (_) >  <| |_| |
 |_|  \___/ \___/  |_|   |_|  \___/ \__\___/ \___\___/|_| |_|   |_|  \___/_/\_\\__, |
                                                                               |___/ `,
		}[0],
		utils.FGBlue,
	))
	dispatcher := new(app.Dispatcher)
	dispatcher.Run()
}

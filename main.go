package main

import (
	"bufio"
	"fmt"
	"github.com/3stadt/swdc/cmd"
	"github.com/blang/semver"
	"github.com/jawher/mow.cli"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var version = "vlatest"

type Config struct {
	CheckUpdateOnStart bool
}

func main() {
	log.SetLevel(log.WarnLevel)
	conf := &Config{}
	conf.checkForUpdate()
	swdc := conf.createApp()
	swdc.Run(os.Args)
}

func (conf *Config) createApp() *cli.Cli {
	app := cli.App("app", "Shopware docker control")
	app.Version("v version", "swdc "+version)
	app.Spec = "[-v]"
	verbose := app.Bool(cli.BoolOpt{
		Name:  "verbose",
		Value: false,
		Desc:  "Enable debug logs",
	})
	app.Before = func() {
		if *verbose {
			log.SetLevel(log.DebugLevel)
		}
	}

	app.Command("s u start up", "start containers", cmd.Start)

	return app
}

func (conf *Config) checkForUpdate() {
	if version == "vlatest" { // version is changed on compile via ldflags, see makefile
		log.Info("using development version, update check deactivated")
		return
	}

	if conf.CheckUpdateOnStart == false {
		log.Info("Auto update check is disabled in config.")
		return
	}

	ver := version[1:]
	latest, found, err := selfupdate.DetectLatest("3stadt/swdc")
	if err != nil {
		log.Error("error occurred while detecting version: ", err.Error())
		return
	}

	v, err := semver.Parse(ver)
	if err != nil {
		log.Error("could not parse current version: ", err.Error())
		return
	}

	if !found || latest.Version.LTE(v) {
		log.Info("using latest version")
		return
	}

	log.Warn("New version available")
	fmt.Println("Please note: Automatic update to a new version always uses the uncompressed binary.")
	fmt.Println("----------")
	fmt.Println(latest.ReleaseNotes)
	fmt.Println("----------")
	fmt.Print("Do you want to update to version ", latest.Version, "? (y/N): ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || strings.ToLower(strings.TrimSpace(input)) != "y" {
		fmt.Println("Skipping update")
		fmt.Printf("You can download the update manually at %s\n", latest.URL)
		return
	}

	log.Warn("Updating to latest version, please be patient...")

	ex, err := os.Executable()
	if err != nil {
		log.Error("error occurred while updating binary: ", err)
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, ex); err != nil {
		log.Error("error occurred while updating binary: ", err)
		return
	}
	log.Info("successfully updated to version ", latest.Version)
}

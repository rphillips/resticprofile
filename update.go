package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/creativeprojects/resticprofile/clog"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func confirmAndSelfUpdate(debug bool) error {
	if debug {
		selfupdate.EnableLog()
	}
	latest, found, err := selfupdate.DetectLatest("creativeprojects/resticprofile")
	if err != nil {
		return fmt.Errorf("Error occurred while detecting version: %v", err)
	}

	v := semver.MustParse(resticProfileVersion)
	if !found || latest.Version.LTE(v) {
		clog.Infof("Current version (%s) is the latest", resticProfileVersion)
		return nil
	}

	fmt.Print("Do you want to update to", latest.Version, "? (y/n): ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (input != "y\n" && input != "n\n") {
		return errors.New("Invalid input")
	}
	if input == "n\n" {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return errors.New("Could not locate executable path")
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		return fmt.Errorf("Error occurred while updating binary: %v", err)
	}
	clog.Infof("Successfully updated to version %s", latest.Version)
	return nil
}

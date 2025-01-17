package main

import (
	"bytes"
	"errors"
	"fmt"
	"lab.sda1.net/nexryai/hoyofeed/logger"
	"lab.sda1.net/nexryai/hoyofeed/upload"
	"os"
	"os/exec"
	"strconv"
)

func generateFeed(lang string) error {
	log := logger.GetLogger("FEED")

	cmd := exec.Command("hoyolab-rss-feeds", "-c", fmt.Sprintf("./config.%s.toml", lang))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Info("generating feed...")
	err := cmd.Run()
	if err != nil {
		return err
	}

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		log.Error("exit code: ", strconv.Itoa(exitCode))
		log.Error("stderr: ", stderr.String())
		return errors.New("exit code is not 0")
	}

	return errors.Join(
		upload.PutFileToAzureBlob("./genshin.xml", fmt.Sprintf("%s/genshin.xml", lang)),
		upload.PutFileToAzureBlob("./genshin.json", fmt.Sprintf("%s/genshin.json", lang)),
		upload.PutFileToAzureBlob("./starrail.xml", fmt.Sprintf("%s/starrail.xml", lang)),
		upload.PutFileToAzureBlob("./starrail.json", fmt.Sprintf("%s/starrail.json", lang)),
	)
}

func main() {
	log := logger.GetLogger("MAIN")

	log.Info("Starting...")

	err := generateFeed("ja-jp")
	if err != nil {
		log.ErrorWithDetail("FAILED:", err)
		os.Exit(1)
	}

	os.Exit(0)
}

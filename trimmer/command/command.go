package command

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"os/exec"
	"strings"
)

type (
	Command struct {
		DataDirPath      string
		YouTubeDLBinPath string
		FFmpegBinPath    string
		SoxBinPath       string
		Log              *bytes.Buffer
		ErrLog           *bytes.Buffer
	}
)

const (
	ServiceIDYouTube     = "youtube"
	ServiceIDTwitCasting = "twit-casting"
	ServiceIDTwitter     = "twitter"
)

func parseURL(videoURL string) (serviceID string, err error) {
	u, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %v", err)
	}

	switch {
	case strings.Contains(u.Host, "youtube.com") || strings.Contains(u.Host, "youtu.be"):
		return ServiceIDYouTube, nil

	case strings.Contains(u.Host, "twitcasting.tv"):
		return ServiceIDTwitCasting, nil

	case strings.Contains(u.Host, "twitter.com"):
		return ServiceIDTwitter, nil

	default:
		return "", fmt.Errorf("unknown service")
	}
}

func (c *Command) runCommand(cmd []string, stdout, stderr io.Writer) error {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = stdout
	command.Stderr = stderr

	log.Printf("start command: %v", cmd)
	defer func() {
		log.Printf("end command (%v): %v", command.ProcessState.ExitCode(), cmd)
	}()

	return command.Run()
}

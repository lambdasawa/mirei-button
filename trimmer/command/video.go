package command

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (c *Command) DownloadVideo(videoURL string) (*os.File, error) {
	defer func() {
		log.Printf("trim stdout: %v", c.Log.String())
		log.Printf("trim stderr: %v", c.ErrLog.String())
	}()

	outputPath, err := c.downloadVideo(videoURL)
	if err != nil {
		return nil, fmt.Errorf("download video: %v", err)
	}

	log.Println("read outputpath", outputPath)

	file, err := os.Open(outputPath)
	if err != nil {
		return nil, fmt.Errorf("open video: %v", err)
	}

	return file, nil
}

func (c *Command) downloadVideo(videoURL string) (string, error) {
	// read cache
	originalPath := c.getOriginalVideoPath(videoURL)
	if _, err := os.Stat(originalPath); err == nil {
		return originalPath, nil
	}

	// create base dir
	_ = os.MkdirAll(filepath.Dir(originalPath), 0755)

	// get service type
	serviceID, err := parseURL(videoURL)
	if err != nil {
		return "", fmt.Errorf("pares url: %v", err)
	}

	// result file path
	outputPath := c.getOriginalVideoPath(videoURL)

	switch serviceID {
	case ServiceIDYouTube, ServiceIDTwitCasting, ServiceIDTwitter:
		// run youtube-dl --format mp4 http://example.com/ -o tmp/hoge/original.mp4
		if err := c.runCommand(
			[]string{
				c.YouTubeDLBinPath,
				"--format",
				"mp4",
				videoURL,
				"-o",
				outputPath,
			},
			c.Log,
			c.ErrLog,
		); err != nil {
			return "", fmt.Errorf("download video: %v", err)
		}
	default:
		return "", fmt.Errorf("unknown service")
	}

	return outputPath, nil
}

func (c *Command) getOriginalVideoPath(videoURL string) string {
	return filepath.Join(
		c.DataDirPath,
		base64.URLEncoding.EncodeToString([]byte(videoURL)),
		"original.mp4",
	)
}

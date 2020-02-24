package command

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (c *Command) TrimSound(videoURL string, startMS int64, durationMS int64) ([]byte, error) {
	defer func() {
		log.Printf("trim stdout: %v", c.Log.String())
		log.Printf("trim stderr: %v", c.ErrLog.String())
	}()

	// get service type
	serviceID, err := parseURL(videoURL)
	if err != nil {
		return nil, err
	}

	// download sound
	originalSoundPath, err := c.downloadSound(serviceID, videoURL)
	if err != nil {
		return nil, fmt.Errorf("download sound: %v", err)
	}

	// trim sound
	trimmedSoundPath, err := c.trimSound(originalSoundPath, startMS, durationMS)
	if err != nil {
		return nil, fmt.Errorf("trim: %v", err)
	}

	blob, err := ioutil.ReadFile(trimmedSoundPath)
	if err != nil {
		return nil, fmt.Errorf("read sound: %v", err)
	}

	return blob, nil
}

func (c *Command) downloadSound(serviceID string, videoURL string) (soundPath string, err error) {
	// if cache exists, do not start command
	originalPath := c.getOriginalSoundPath(videoURL)
	if _, err := os.Stat(originalPath); err == nil {
		return originalPath, nil
	}

	_ = os.MkdirAll(filepath.Dir(originalPath), 0755)

	switch serviceID {
	case ServiceIDYouTube:
		return c.downloadYouTube(videoURL)
	case ServiceIDTwitCasting:
		return c.downloadTwitCas(videoURL)
	case ServiceIDTwitter:
		return c.downloadTwitter(videoURL)
	default:
		return "", fmt.Errorf("unknown service: %v", serviceID)
	}
}

func (c *Command) downloadYouTube(videoURL string) (soundPath string, err error) {
	originalPath := c.getOriginalSoundPath(videoURL)

	// run youtube-dl -x --audio-format mp3 http://example.com --output tmp/hoge/original.mp3
	if err := c.runCommand(
		[]string{
			c.YouTubeDLBinPath,
			"-x",
			"--audio-format",
			"mp3",
			videoURL,
			"--output",
			filepath.Join(filepath.Dir(originalPath), "original.(exts)"),
		},
		c.Log,
		c.ErrLog,
	); err != nil {
		return "", fmt.Errorf("donwload by youtube-dl: %v", err)
	}

	return originalPath, nil
}

func (c *Command) downloadTwitCas(videoURL string) (soundPath string, err error) {
	return c.downloadByScraping(ServiceIDTwitCasting, videoURL)
}

func (c *Command) downloadTwitter(videoURL string) (soundPath string, err error) {
	return c.downloadByScraping(ServiceIDTwitter, videoURL)
}

func (c *Command) downloadByScraping(service string, videoURL string) (soundPath string, err error) {
	// run youtube-dl -j http://example.com
	scrapeResult := new(bytes.Buffer)
	if err := c.runCommand(
		[]string{c.YouTubeDLBinPath, "-j", videoURL},
		scrapeResult,
		c.ErrLog,
	); err != nil {
		return "", fmt.Errorf("scrape by youtube-dl: %v", err)
	}

	scrapeResultJSON := map[string]interface{}{}
	if err := json.Unmarshal(scrapeResult.Bytes(), &scrapeResultJSON); err != nil {
		return "", fmt.Errorf("parse scraping result: %v", err)
	}

	scrapeResultURL := scrapeResultJSON["url"].(string)

	originalPath := c.getOriginalSoundPath(videoURL)

	// run ffmpeg -i http://example.m3u8 -vn tmp/hoge/original.mp3
	if err := c.runCommand(
		[]string{
			c.FFmpegBinPath,
			"-i",
			scrapeResultURL,
			"-vn",
			originalPath,
		},
		c.Log,
		c.ErrLog,
	); err != nil {
		return "", fmt.Errorf("download by ffmpeg: %v", err)
	}

	return originalPath, nil
}

func (c *Command) getOriginalSoundPath(videoURL string) string {
	videoID := base64.URLEncoding.EncodeToString([]byte(videoURL))

	return filepath.Join(c.DataDirPath, videoID, "original.mp3")
}

func (c *Command) getTrimmedSoundPath(originalPath string, startMS, durationMS int64) string {
	return filepath.Join(
		strings.Replace(originalPath, "original.mp3", "", -1),
		fmt.Sprint(startMS),
		fmt.Sprint(durationMS),
		"trim.mp3",
	)
}

func (c *Command) trimSound(originalSoundPath string, startMS, durationMS int64) (trimmedSoundPath string, err error) {
	// if cache exists, do not start command
	trimmedSoundPath = c.getTrimmedSoundPath(originalSoundPath, startMS, durationMS)
	if _, err := os.Stat(trimmedSoundPath); err == nil {
		return trimmedSoundPath, nil
	}

	// create base dir
	_ = os.MkdirAll(filepath.Dir(trimmedSoundPath), 0755)

	// run sox tmp/hoge/original.mp3 tmp/hoge/1000/2000/trim.mp3 trim 1000 2000
	if err := c.runCommand(
		[]string{
			c.SoxBinPath,
			originalSoundPath,
			trimmedSoundPath,
			"trim",
			fmt.Sprintf("%d.%03d", startMS/1000, startMS%1000),
			fmt.Sprintf("%d.%03d", durationMS/1000, durationMS%1000),
		},
		c.Log,
		c.ErrLog,
	); err != nil {
		return "", fmt.Errorf("trim by sox: %v", err)
	}

	return trimmedSoundPath, nil
}

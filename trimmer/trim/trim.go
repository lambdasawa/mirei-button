package trim

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type (
	Trimmer struct {
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

func (t *Trimmer) Trim(videoURL string, startMS int64, durationMS int64) (soundPath string, err error) {
	defer func() {
		log.Printf("trim stdout: %v", t.Log.String())
		log.Printf("trim stderr: %v", t.ErrLog.String())
	}()

	serviceID, err := parseURL(videoURL)
	if err != nil {
		return "", err
	}

	originalSoundPath, err := t.download(serviceID, videoURL)
	if err != nil {
		return "", fmt.Errorf("download: %v", err)
	}

	trimmedSoundPath, err := t.trim(originalSoundPath, startMS, durationMS)
	if err != nil {
		return "", fmt.Errorf("trim: %v", err)
	}

	return trimmedSoundPath, nil
}

func parseURL(videoURL string) (serviceID string, err error) {
	u, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %v", err)
	}

	switch {
	case strings.Contains(u.Host, "youtube.com"):
		return ServiceIDYouTube, nil

	case strings.Contains(u.Host, "twitcasting.tv"):
		return ServiceIDTwitCasting, nil

	case strings.Contains(u.Host, "twitter.com"):
		return ServiceIDTwitter, nil

	default:
		return "", fmt.Errorf("unknown service")
	}
}

func (t *Trimmer) download(serviceID string, videoURL string) (soundPath string, err error) {
	// if cache exists, do not start command
	originalPath := t.getOriginalPath(serviceID, videoURL)
	if _, err := os.Stat(originalPath); err == nil {
		return originalPath, nil
	}

	_ = os.MkdirAll(filepath.Dir(originalPath), 0755)

	switch serviceID {
	case ServiceIDYouTube:
		return t.downloadYouTube(videoURL)
	case ServiceIDTwitCasting:
		return t.downloadTwitCas(videoURL)
	case ServiceIDTwitter:
		return t.downloadTwitter(videoURL)
	default:
		return "", fmt.Errorf("unknown service: %v", serviceID)
	}
}

func (t *Trimmer) downloadYouTube(videoURL string) (soundPath string, err error) {
	originalPath := t.getOriginalPath(ServiceIDYouTube, videoURL)
	if err := t.runCommand(
		[]string{
			t.YouTubeDLBinPath,
			"-x",
			"--audio-format",
			"mp3",
			videoURL,
			"--output",
			filepath.Join(filepath.Dir(originalPath), "original.%(ext)s"),
		},
		t.Log,
		t.ErrLog,
	); err != nil {
		return "", fmt.Errorf("donwload by youtube-dl: %v", err)
	}

	return originalPath, nil
}

func (t *Trimmer) downloadTwitCas(videoURL string) (soundPath string, err error) {
	return t.downloadByScraping(ServiceIDTwitCasting, videoURL)
}

func (t *Trimmer) downloadTwitter(videoURL string) (soundPath string, err error) {
	return t.downloadByScraping(ServiceIDTwitter, videoURL)
}

func (t *Trimmer) downloadByScraping(service string, videoURL string) (soundPath string, err error) {
	// if cache exists, do not start command
	scrapeResult := new(bytes.Buffer)
	if err := t.runCommand(
		[]string{t.YouTubeDLBinPath, "-j", videoURL},
		scrapeResult,
		t.ErrLog,
	); err != nil {
		return "", fmt.Errorf("scrape by youtube-dl: %v", err)
	}

	scrapeResultJSON := map[string]interface{}{}
	if err := json.Unmarshal(scrapeResult.Bytes(), &scrapeResultJSON); err != nil {
		return "", fmt.Errorf("parse scraping result: %v", err)
	}

	scrapeResultURL := scrapeResultJSON["url"].(string)

	originalPath := t.getOriginalPath(service, videoURL)
	if err := t.runCommand(
		[]string{
			t.FFmpegBinPath,
			"-i",
			scrapeResultURL,
			"-vn",
			originalPath,
		},
		t.Log,
		t.ErrLog,
	); err != nil {
		return "", fmt.Errorf("download by ffmpeg: %v", err)
	}

	return originalPath, nil
}

func (t *Trimmer) getOriginalPath(serviceID string, videoURL string) string {
	videoID := base64.URLEncoding.EncodeToString([]byte(videoURL))

	return filepath.Join(t.DataDirPath, serviceID, videoID, "original.mp3")
}

func (t *Trimmer) getTrimmedPath(originalPath string, startMS, durationMS int64) string {
	return filepath.Join(
		strings.Replace(originalPath, "original.mp3", "", -1),
		fmt.Sprint(startMS),
		fmt.Sprint(durationMS),
		"trimmed.mp3",
	)
}

func (t *Trimmer) runCommand(cmd []string, stdout, stderr io.Writer) error {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = stdout
	command.Stderr = stderr

	log.Printf("start command: %v", cmd)
	defer func() {
		log.Printf("end command (%v): %v", command.ProcessState.ExitCode(), cmd)
	}()

	return command.Run()
}

func (t *Trimmer) trim(originalSoundPath string, startMS, durationMS int64) (trimmedSoundPath string, err error) {
	// if cache exists, do not start command
	trimmedSoundPath = t.getTrimmedPath(originalSoundPath, startMS, durationMS)
	if _, err := os.Stat(trimmedSoundPath); err == nil {
		return trimmedSoundPath, nil
	}

	_ = os.MkdirAll(filepath.Dir(trimmedSoundPath), 0755)

	if err := t.runCommand(
		[]string{
			t.SoxBinPath,
			originalSoundPath,
			trimmedSoundPath,
			"trim",
			fmt.Sprintf("%d.%03d", startMS/1000, startMS%1000),
			fmt.Sprintf("%d.%03d", durationMS/1000, durationMS%1000),
		},
		t.Log,
		t.ErrLog,
	); err != nil {
		return "", fmt.Errorf("trim by sox: %v", err)
	}

	return trimmedSoundPath, nil
}

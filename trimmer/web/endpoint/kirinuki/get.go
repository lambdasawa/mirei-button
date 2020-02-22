package kirinuki

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mb-trimmer/trim"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type (
	Req struct {
		URL        string `query:"url"`
		StartMS    int64  `query:"start-ms"`
		DurationMS int64  `query:"duration-ms"`
	}
)

func Get(c echo.Context) error {
	var req Req
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("read request: %v", err)
	}

	trimmer := &trim.Trimmer{
		DataDirPath:      "./tmp/",
		YouTubeDLBinPath: os.Getenv("MB_YOUTUBEDL_BIN_PATH"),
		FFmpegBinPath:    os.Getenv("MB_FFMPEG_BIN_PATH"),
		SoxBinPath:       os.Getenv("MB_SOX_BIN_PATH"),
		Log:              new(bytes.Buffer),
		ErrLog:           new(bytes.Buffer),
	}

	soundPath, err := trimmer.Trim(req.URL, req.StartMS, req.DurationMS)
	if err != nil {
		return fmt.Errorf("trim sound: %v", err)
	}

	blob, err := ioutil.ReadFile(soundPath)
	if err != nil {
		return fmt.Errorf("read trimmed sound: %v", err)
	}

	return c.Blob(http.StatusOK, "audio/mpeg", blob)
}

package video

import (
	"bytes"
	"fmt"
	"mb-trimmer/command"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	Req struct {
		URL string `query:"url"`
	}
)

func Get(c echo.Context) error {
	var req Req
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("read request: %v", err)
	}

	cmd := &command.Command{
		DataDirPath:      "./tmp/",
		YouTubeDLBinPath: os.Getenv("MB_YOUTUBEDL_BIN_PATH"),
		FFmpegBinPath:    os.Getenv("MB_FFMPEG_BIN_PATH"),
		SoxBinPath:       os.Getenv("MB_SOX_BIN_PATH"),
		Log:              new(bytes.Buffer),
		ErrLog:           new(bytes.Buffer),
	}

	blob, err := cmd.DownloadVideo(req.URL)
	if err != nil {
		return fmt.Errorf("trim sound: %v", err)
	}

	readSeeker := bytes.NewReader(blob)

	http.ServeContent(c.Response(), c.Request(), "original.mp4", time.Now(), readSeeker)

	return nil
}

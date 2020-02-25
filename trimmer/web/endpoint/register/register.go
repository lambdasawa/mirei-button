package register

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mb-trimmer/command"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
)

type (
	Req struct {
		URL        string `json:"url"`
		StartMS    int64  `json:"startMS"`
		DurationMS int64  `json:"durationMS"`
		Text       string `json:"text"`
		Tags       string `json:"tags"` // csv
	}

	Metadata struct {
		Items []MetadataItem `json:"items"`
	}

	MetadataItem struct {
		URL         string   `json:"url"`
		OriginalURL string   `json:"originalUrl"`
		StartMS     int64    `json:"startMS"`
		DurationMS  int64    `json:"durationMS"`
		Text        string   `json:"text"`
		Tags        []string `json:"tags"`
		CreatedAt   string   `json:"createdAt"`
	}
)

var (
	bucketName = os.Getenv("MB_BUCKET")
)

func Register(c echo.Context) error {
	var req Req
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("read request: %v", err)
	}

	if err := register(req); err != nil {
		return fmt.Errorf("register: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func BulkRegister(c echo.Context) error {
	tsvReader := csv.NewReader(c.Request().Body)
	tsvReader.Comma = '\t'

	records, err := tsvReader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		log.Println(record)

		url := record[0]

		start := record[1]
		t, err := time.Parse("15:04:05.999999999", start)
		if err != nil {
			return fmt.Errorf("parse start: %v", err)
		}
		startMS := int64(0)
		startMS += int64(t.Hour() * 60 * 60 * 1000)
		startMS += int64(t.Minute() * 60 * 1000)
		startMS += int64(t.Second() * 1000)
		startMS += int64(t.Nanosecond() / 1000 / 1000)

		duration := record[2]
		durationMS, err := time.ParseDuration(duration)
		if err != nil {
			return fmt.Errorf("parse duration: %v", err)
		}

		text := record[3]

		req := Req{
			URL:        url,
			StartMS:    startMS,
			DurationMS: durationMS.Milliseconds(),
			Text:       text,
		}
		log.Println(req)
		if err := register(req); err != nil {
			return fmt.Errorf("regsiter: %v", err)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func register(req Req) error {
	// new s3 service
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		return fmt.Errorf("aws session: %v", err)
	}
	s3Service := s3.New(sess)

	// create metadata
	metadata, err := createNewMetadata(s3Service, req)
	if err != nil {
		return fmt.Errorf("create metadata: %v", err)
	}
	log.Println("metadata", metadata)

	// create mp3
	mp3, err := createMP3(req)
	if err != nil {
		return fmt.Errorf("trim sound: %v", err)
	}
	log.Println("mp3", len(mp3))

	// todo create zip

	// todo upload zip

	// upload mp3
	if _, err := s3Service.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(buildKey(req)),
		Body:        bytes.NewReader(mp3),
		ACL:         aws.String(s3.BucketCannedACLPublicRead),
		ContentType: aws.String("audio/mpeg"),
	}); err != nil {
		return fmt.Errorf("put mp3 object: %v", err)
	}

	// upload metadata
	mbBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("encode metadata to json: %v", err)
	}
	if _, err := s3Service.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String("metadata.json"),
		Body:        bytes.NewReader(mbBytes),
		ACL:         aws.String(s3.BucketCannedACLPublicRead),
		ContentType: aws.String("application/json"),
	}); err != nil {
		return fmt.Errorf("put mp3 object: %v", err)
	}

	// invalidate cache
	cf := cloudfront.New(sess)
	if err := invalidateCache(cf); err != nil {
		return fmt.Errorf("invalidate cache: %v", err)
	}

	// todo post tweet

	return nil
}

func createNewMetadata(s3Service *s3.S3, req Req) (Metadata, error) {
	md := Metadata{}

	s3Output, err := s3Service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("metadata.json"),
	})

	// todo refactor
	switch err {
	case nil:
		defer s3Output.Body.Close()

		bytes, err := ioutil.ReadAll(s3Output.Body)
		if err != nil {
			return Metadata{}, fmt.Errorf("read s3 object: %v", err)
		}

		if err := json.Unmarshal(bytes, &md); err != nil {
			return Metadata{}, fmt.Errorf("parse metadata: %v", err)
		}

	default:
		switch aerr, ok := err.(awserr.Error); ok {
		case true:
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				// valid
			default:
				return Metadata{}, fmt.Errorf("get object (AWS error): %v", aerr)
			}
		default:
			return Metadata{}, fmt.Errorf("get object (unknown error): %v", err)
		}
	}

	md.Items = append(md.Items, MetadataItem{
		URL:         buildKey(req),
		OriginalURL: req.URL,
		StartMS:     req.StartMS,
		DurationMS:  req.DurationMS,
		Text:        req.Text,
		Tags:        strings.Split(req.Tags, ","),
		CreatedAt:   time.Now().Format(time.RFC3339),
	})

	return md, nil
}

func buildKey(req Req) string {
	return path.Join(
		"/",
		base64.URLEncoding.EncodeToString([]byte(req.URL)),
		fmt.Sprint(req.StartMS),
		fmt.Sprint(req.DurationMS),
		"trimmed.mp3",
	)
}

func createMP3(req Req) ([]byte, error) {
	cmd := &command.Command{
		DataDirPath:      "./tmp/",
		YouTubeDLBinPath: os.Getenv("MB_YOUTUBEDL_BIN_PATH"),
		FFmpegBinPath:    os.Getenv("MB_FFMPEG_BIN_PATH"),
		SoxBinPath:       os.Getenv("MB_SOX_BIN_PATH"),
		Log:              new(bytes.Buffer),
		ErrLog:           new(bytes.Buffer),
	}
	mp3, err := cmd.TrimSound(req.URL, req.StartMS, req.DurationMS)
	if err != nil {
		return nil, fmt.Errorf("trim sound: %v", err)
	}

	return mp3, nil
}

func invalidateCache(cfService *cloudfront.CloudFront) error {
	paths := []*string{
		aws.String("/index.html"),
		aws.String("/metadata.json"),
		aws.String("/favicon.ico"),
		aws.String("/js/*"),
		aws.String("/css/*"),
		aws.String("/img/*"),
	}
	if _, err := cfService.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(os.Getenv("MB_DISTRIBUTION_ID")),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(fmt.Sprint(time.Now().UnixNano() / 1000 / 1000)),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(int64(len(paths))),
				Items:    paths,
			},
		},
	}); err != nil {
		return fmt.Errorf("create invalidation: %v", err)
	}

	return nil
}

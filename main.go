package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
)

var (
	baseEndpoint = "https://s3.isk01.sakurastorage.jp"
	bucketName   = os.Getenv("BUCKET_NAME")
)

func main() {
	ctx := context.Background()

	r := gin.Default()

	r.GET("/:params/:key", func(c *gin.Context) {
		key := c.Param("key")
		buf, err := getImage(ctx, &bucketName, &key)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusNotFound)
			return
		}

		var width, height int
		var t bimg.ImageType

		re := regexp.MustCompile(`f=([^,]+)|w=([^,]+)|h=([^,]+)`)
		matches := re.FindAllStringSubmatch(c.Param("params"), -1)

		for _, match := range matches {
			if match[1] != "" {
				t, err = getImageType(match[1])
				if err != nil {
					log.Println(err)
					c.Status(http.StatusInternalServerError)
					return
				}
			}
			if match[2] != "" {
				width, err = strconv.Atoi(match[2])
				if err != nil {
					log.Println(err)
					c.Status(http.StatusInternalServerError)
					return
				}
			}
			if match[3] != "" {
				height, err = strconv.Atoi(match[3])
				if err != nil {
					log.Println(err)
					c.Status(http.StatusInternalServerError)
					return
				}
			}
		}

		resData, err := convert(*buf, t, width, height)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, http.DetectContentType(resData), resData)
	})

	r.Run(":8080")
}

func convert(buf []byte, t bimg.ImageType, width, height int) ([]byte, error) {
	return bimg.Resize(buf, bimg.Options{
		Type:   t,
		Width:  width,
		Height: height,
	})
}

func getImageType(format string) (bimg.ImageType, error) {
	var t bimg.ImageType
	var err error

	switch format {
	case "webp":
		t = bimg.WEBP
	case "jpg", "jpeg":
		t = bimg.JPEG
	case "png":
		t = bimg.PNG
	case "gif":
		t = bimg.GIF
	default:
		err = errors.New("unsupported format")
	}

	return t, err
}

func getImage(ctx context.Context, bucketName, key *string) (*[]byte, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = &baseEndpoint
	})

	o, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: bucketName,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}

	defer o.Body.Close()

	data, err := io.ReadAll(o.Body)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

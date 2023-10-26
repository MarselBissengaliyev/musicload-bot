package youtube

import (
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/MarselBissengaliyev/realtime-chat/downloader"
	"github.com/kkdai/youtube/v2"
	"github.com/pkg/errors"
)

const (
	expression = "^(http(s)?:\\/\\/)?((w){3}.)?(music\\.)?youtu(be|.be)?(\\.com)?\\/.+"
)

type Downloader struct {
	maxVideoDuration time.Duration
	r                *regexp.Regexp
}

func NewDownloader(maxVideoDuration int64) (*Downloader, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	return &Downloader{
		maxVideoDuration: time.Minute * time.Duration(maxVideoDuration),
		r:                r,
	}, nil
}

func (d *Downloader) Download(maxVideoDuration int64, url string) (string, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	
	if err != nil {
		return "", errors.Wrap(err, "error getting video info")
	}

	if video.Duration.Minutes() > float64(maxVideoDuration) {
		return "", downloader.ErrDurationLimitExceeded
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])

	if err != nil {
		return "", errors.Wrap(err, "error getting stream")
	}
	defer stream.Close()

	filename := video.ID
	strings.Replace(filename, " ", "\\ ", 0)

	file, err := os.Create(filename + ".mp3")
	if err != nil {
		return "", errors.Wrap(err, "error create file")
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", errors.Wrap(err, "error copy file")
	}

	return filename + ".mp3", nil
}

func (d *Downloader) IsValidURL(url string) bool {
	return d.r.MatchString(url)
}

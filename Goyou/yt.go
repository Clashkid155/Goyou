package Goyou

import (
	"context"
	"fmt"
	yb "github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"log"
)

type Details struct {
	Title   string
	Stream  yb.Format
	Size    string
	Quality string
	Thumb   yb.Thumbnail
	id      string
	video   *yb.Video
}

var client = yb.Client{}

var down = &downloader.Downloader{
	Client:    client,
	OutputDir: "downloads/",
}

func Query(s string) (youtubeDetails []Details) {

	video, err := client.GetVideo(s)
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	title := video.Title
	for i, k := range video.Formats {
		size := fmt.Sprintf("%0.1f MB\n", float64(video.Formats[i].Bitrate)*video.Duration.Seconds()/8/1024/1024)
		video.Formats.Type("")
		if video.Formats[i].QualityLabel != "" && video.Formats[i].QualityLabel != "144p" && video.Formats[i].QualityLabel != "240p" {
			youtubeDetails = append(youtubeDetails, Details{
				Title:   title,
				Stream:  k,
				Size:    size,
				Quality: k.QualityLabel,
				Thumb:   video.Thumbnails[2],
				id:      video.ID,
				video:   video,
			})
		}
	}
	return
}

func Download(vid Details) string {
	video := vid.video
	quality := vid.Quality
	name := down.OutputDir + downloade(video, quality)
	return name
}

func downloade(video *yb.Video, quality string) string {
	filename := fmt.Sprintf("%s_%s.mp4", downloader.SanitizeFilename(video.Title), quality) //      + "_" + quality + ".mp4"
	err := down.DownloadComposite(context.Background(), filename, video, quality, "")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(video.Title, " Downloaded ")
	return filename
}

package main

import (
	"context"
	"fmt"
	yb "github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"log"
)

type Details struct {
	title   string
	stream  yb.Format
	size    string
	quality string
	thumb   yb.Thumbnail
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
	//video.Formats.Sort()
	for i, k := range video.Formats {
		size := fmt.Sprintf("%0.1f MB\n", float64(video.Formats[i].Bitrate)*video.Duration.Seconds()/8/1024/1024)
		//fmt.Println(title, stream[i].QualityLabel)
		//fmt.Println(title, k.Quality)
		//fmt.Println(video.Formats[i].Quality, k.Quality)
		//fmt.Printf(size)
		//fmt.Println(video.Formats[i].ContentLength, k.ContentLength, k.MimeType)
		//fmt.Println("With  Audio", video.Formats.WithAudioChannels()[i].QualityLabel)
		//fmt.Println("Without ", video.Formats[i].QualityLabel)
		video.Formats.Type("")
		if video.Formats[i].QualityLabel != "" && video.Formats[i].QualityLabel != "144p" && video.Formats[i].QualityLabel != "240p" {
			youtubeDetails = append(youtubeDetails, Details{
				title:   title,
				stream:  k,
				size:    size,
				quality: k.QualityLabel,
				thumb:   video.Thumbnails[2],
				id:      video.ID,
				video:   video,
			})
		}
	}
	//fmt.Println(video.Formats.FindByQuality("480p"))
	return
}

func Download(vid Details) string {
	video := vid.video
	quality := vid.quality
	name := down.OutputDir + downloade(video, quality)

	//file := down.OutputDir + filename

	fmt.Println(name)
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

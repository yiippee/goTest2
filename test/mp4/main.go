//package main
//
//import (
//	"log"
//	"net/http"
//	"os"
//	"time"
//)
//
//func ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	video, err := os.Open("D:/Program Files (x86)/ffmpeg-20191023-1f327f5-win64-static/bin/test.mp4")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer video.Close()
//
//	http.ServeContent(w, r, "test.mp4", time.Now(), video)
//}
//
//func main() {
//	http.HandleFunc("/", ServeHTTP)
//	http.ListenAndServe(":8080", nil)
//}
package main

import (
"fmt"

"github.com/nareix/joy4/av"
"github.com/nareix/joy4/av/avutil"
"github.com/nareix/joy4/format"
)

func init() {
	format.RegisterAll()
}

func main() {
	file, err := avutil.Open("D:/Program Files (x86)/ffmpeg-20191023-1f327f5-win64-static/bin/test.mp4")
	if err != nil {
		fmt.Println(err)
	}
	streams, _ := file.Streams()
	for _, stream := range streams {
		if stream.Type().IsAudio() {
			astream := stream.(av.AudioCodecData)
			fmt.Println(astream.Type(), astream.SampleRate(), astream.SampleFormat(), astream.ChannelLayout())
		} else if stream.Type().IsVideo() {
			vstream := stream.(av.VideoCodecData)
			fmt.Println(vstream.Type(), vstream.Width(), vstream.Height())
		}
	}

	file.Close()
}

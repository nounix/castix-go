package omxplayer

import (
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/i-PUSH/RPi-Cast/src/utils"

	"github.com/go-chi/chi"
)

// Omxplayer data for templates
type Omxplayer struct {
	Buttons *[]Button
	Youtube *[]Youtube
}

// Button of omxplayer
type Button struct {
	Action, Text string
}

// Youtube for the scraped data
type Youtube struct {
	Vid, Title, Img string
}

// Init web controller for omxplayer
func Init(app *chi.Mux) {
	app.Get("/omxplayer", func(w http.ResponseWriter, r *http.Request) {
		utils.ExecTemplate(w, "resources/templates/omxplayer.html", Omxplayer{Buttons: getOmxBtns(), Youtube: nil})
	})
	app.Get("/omxplayer/control", omxCtrl)
	app.Get("/omxplayer/youtube/related", relatedYoutube)
	app.Get("/omxplayer/youtube/search", searchYoutube)
	app.Post("/omxplayer/play", playVideo)
}

func getOmxBtns() *[]Button {
	var omxBtns = []Button{
		{"SeekBackLarge", "fa-fast-backward"},
		{"SeekBackSmall", "fa-step-backward"},
		{"PlayPause", "fa-pause"},
		{"SeekForwardSmall", "fa-step-forward"},
		{"SeekForwardLarge", "fa-fast-forward"},
		{"DecreaseVolume", "fa-volume-down"},
		{"IncreaseVolume", "fa-volume-up"},
		{"Exit", "fa-stop"},
	}

	return &omxBtns
}

// send command to omxplayer via dbus
func omxCtrl(w http.ResponseWriter, r *http.Request) {
	omx := NewOmxCtrl()

	omx.Action(utils.GetURLParam(r, "action"))

	utils.WriteJSON(w, "status: omxplayer "+utils.GetURLParam(r, "action"))
}

// play direct link video
func playVideo(w http.ResponseWriter, r *http.Request) {
	req := struct{ Data string }{}
	utils.ReadJSON(r, &req)
	// get youtube link and play vid
	out, _ := exec.Command("youtube-dl", "-g", "-f", "mp4", req.Data).Output()
	PlayVid(string(out[:len(out)-1]))
	utils.WriteJSON(w, "status: omxplayer playing video")
}

// PlayVid and restart video on crash
func PlayVid(url string) {
	var vidPos string
	var vidDone bool

	if omx := NewOmxCtrl(); omx.CanQuit() {
		omx.Action("Exit")
	}

	chkVid := utils.DoEvery(5, func() {
		vidPos = utils.Sec2Time(NewOmxCtrl().Position())
	})

	go func() {
		for !vidDone {
			err := exec.Command("omxplayer", "--no-keys", "--blank", "--pos", vidPos, url).Run()
			if err == nil || err.Error() == "exit status 3" {
				vidDone = true
				close(chkVid)
			}
		}
	}()
}

func searchYoutube(w http.ResponseWriter, r *http.Request) {
	var data []Youtube
	var vids []string
	var titles []string
	var imgs []string

	ytSearch := "https://www.youtube.com/results?search_query=" + url.QueryEscape(utils.GetURLParam(r, "str"))

	doc, err := goquery.NewDocument(ytSearch)
	if err != nil {
		utils.LogErr(err)
	}

	// get list of videos searched for
	videos := doc.Find("ol.item-section").Find("li")

	// find video ids
	videos.Find("div[data-context-item-id]").Each(func(i int, s *goquery.Selection) {
		vids = append(vids, utils.MuteStr(s.Attr("data-context-item-id")))
	})

	// find video titles
	videos.Find("h3.yt-lockup-title").Each(func(i int, s *goquery.Selection) {
		titles = append(titles, s.Text())
	})

	// find video thumbnails
	videos.Find("img").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Attr("src")
		if strings.Contains(img, ".gif") {
			img, _ = s.Attr("data-thumb")
		}
		if img != "" {
			imgs = append(imgs, img)
		}
	})

	// merge the youtube data
	for i := range vids {
		data = append(data, Youtube{vids[i], titles[i], imgs[i]})
	}

	utils.ExecTemplate(w, "resources/templates/youtube.html", Omxplayer{Buttons: getOmxBtns(), Youtube: &data})
}

func relatedYoutube(w http.ResponseWriter, r *http.Request) {
	var data []Youtube
	var vids []string
	var titles []string
	var imgs []string

	ytURL := "https://www.youtube.com/watch?v=" + utils.GetURLParam(r, "vid")

	doc, err := goquery.NewDocument(ytURL)
	if err != nil {
		utils.LogErr(err)
	}

	// get list of videos searched for
	videos := doc.Find("div#watch7-sidebar-modules").Find("div.watch-sidebar-section")

	// get data about first related vid
	nextVid := videos.First().Find("li.related-list-item").First()
	vids = append(vids, utils.MuteStr(nextVid.Find("div.thumb-wrapper").Find("span[data-vid]").Attr("data-vid")))
	imgs = append(imgs, utils.MuteStr(nextVid.Find("div.thumb-wrapper").Find("img").Attr("data-thumb")))
	vidTitle := nextVid.Find("div.content-wrapper").Find("span.title").Text()
	vidTime := nextVid.Find("div.thumb-wrapper").Find("span.video-time").Text()
	titles = append(titles, vidTitle+" "+vidTime)

	// get data about the remaining related vids
	videos.Last().Find("ul#watch-related").Find("li.related-list-item-compact-video").Each(func(i int, s *goquery.Selection) {
		vids = append(vids, utils.MuteStr(s.Find("div.thumb-wrapper").Find("span[data-vid]").Attr("data-vid")))
		imgs = append(imgs, utils.MuteStr(s.Find("div.thumb-wrapper").Find("img").Attr("data-thumb")))
		vidTitle := s.Find("div.content-wrapper").Find("span.title").Text()
		vidTime := s.Find("div.thumb-wrapper").Find("span.video-time").Text()
		titles = append(titles, vidTitle+" "+vidTime)
	})

	// merge the youtube data
	for i := range vids {
		data = append(data, Youtube{vids[i], titles[i], imgs[i]})
	}

	utils.ExecTemplate(w, "resources/templates/youtube.html", Omxplayer{Buttons: getOmxBtns(), Youtube: &data})

	// get youtube link and play vid
	out, _ := exec.Command("youtube-dl", "-g", "-f", "mp4", ytURL).Output()
	PlayVid(string(out[:len(out)-1]))
}

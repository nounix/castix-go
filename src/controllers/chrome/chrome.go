package chrome

import (
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/i-PUSH/RPi-Cast/src/controllers/omxplayer"
	"github.com/i-PUSH/RPi-Cast/src/utils"
)

// Button of chrome
type Button struct {
	Cmd, Key, Text string
}

// Init web controller for chrome
func Init(app *chi.Mux) {
	app.Get("/chrome", index)
	app.Get("/chrome/xdotool/type", xdotoolType)
	app.Get("/chrome/xdotool/key", xdotoolKey)
	app.Get("/chrome/xdotool/click", xdotoolClick)
	app.Post("/chrome/xdotool/mousemove", xdotoolMouseMove)
	app.Get("/chrome/omxplayer/play", playVideo)
}

func index(w http.ResponseWriter, r *http.Request) {
	chromeBtns := [][]Button{
		{
			{"/chrome/xdotool/key", "Return", "Return"},
			{"/chrome/xdotool/key", "space", "Space"},
			{"/chrome/xdotool/key", "BackSpace", "BackSpace"},
			{"/chrome/xdotool/key", "Escape", "Escape"},
		},
		{
			{"/chrome/xdotool/key", "Up", "Up"},
			{"/chrome/xdotool/key", "Down", "Down"},
			{"/chrome/xdotool/key", "Left", "Left"},
			{"/chrome/xdotool/key", "Right", "Right"},
		},
		{
			{"/chrome/xdotool/key", "ctrl+c", "Copy"},
			{"/chrome/xdotool/key", "ctrl+v", "Paste"},
			{"/chrome/xdotool/key", "ctrl+w", "Tab Close"},
			{"/chrome/xdotool/key", "ctrl+t", "Tab New"},
		},
		{
			{"/chrome/xdotool/click", "2", "Mouse Middle"},
			{"/chrome/xdotool/click", "3", "Mouse Right"},
			{"/chrome/xdotool/click", "8", "Back"},
			{"/chrome/xdotool/click", "9", "Forward"},
		},
	}

	// start xorg if not already running
	out, err := exec.Command("ps", "ax").Output()
	utils.LogErr(err)
	xorg := strings.Contains(string(out[:]), "/usr/lib/xorg/Xorg")

	if !xorg {
		exec.Command("startx").Start()
	}

	utils.ExecTemplate(w, "resources/templates/chrome.html", chromeBtns)
}

func xdotoolType(w http.ResponseWriter, r *http.Request) {
	xdotool("type", utils.GetURLParam(r, "str"))
	utils.WriteJSON(w, "status: xdotool type: "+utils.GetURLParam(r, "str"))
}

func xdotoolKey(w http.ResponseWriter, r *http.Request) {
	xdotool("key", strings.Replace(utils.GetURLParam(r, "key"), " ", "+", -1))
	utils.WriteJSON(w, "status: xdotool press key: "+utils.GetURLParam(r, "key"))
}

func xdotoolClick(w http.ResponseWriter, r *http.Request) {
	xdotool("click", utils.GetURLParam(r, "key"))
	utils.WriteJSON(w, "status: xdotool press mouse key: "+utils.GetURLParam(r, "key"))
}

func xdotoolMouseMove(w http.ResponseWriter, r *http.Request) {
	req := struct {
		X1 int
		X2 int
		Y1 int
		Y2 int
	}{}

	utils.ReadJSON(r, &req)

	x := strconv.Itoa(req.X2 - req.X1)
	y := strconv.Itoa(req.Y2 - req.Y1)

	if req.X1 == req.X2 && req.Y1 == req.Y2 {
		xdotool("click", "1")
	} else {
		xdotool("mousemove_relative", "--", x, y)
	}

	utils.WriteJSON(w, "status: xdotool move mouse: "+x+","+y)
}

func xdotool(arg ...string) {
	cmd := exec.Command("xdotool", arg...)
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	cmd.Run()
}

func playVideo(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("xclip", "-o")
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	out, _ := cmd.Output()
	url := string(out[:])
	utils.WriteJSON(w, "status: chrome omxplayer play: "+url)
	omxplayer.PlayVid(url)
}

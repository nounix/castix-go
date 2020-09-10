package omxplayer

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/godbus/dbus"
	"github.com/i-PUSH/RPi-Cast/src/utils"
)

const (
	microsecond     int64 = 1000000
	rootInterface         = "org.mpris.MediaPlayer2"
	playerInterface       = "org.mpris.MediaPlayer2.Player"
	propertyGetter        = "org.freedesktop.DBus.Properties.Get"
)

type OmxCtrl struct {
	conn      *dbus.Conn
	omxPlayer dbus.BusObject
}

func NewOmxCtrl() *OmxCtrl {
	user := os.Getenv("USER")
	address, err := ioutil.ReadFile(fmt.Sprintf("/tmp/omxplayerdbus.%s", user))
	pid, err := ioutil.ReadFile(fmt.Sprintf("/tmp/omxplayerdbus.%s.pid", user))
	os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(address))
	os.Setenv("DBUS_SESSION_BUS_PID", string(pid))

	conn, err := dbus.SessionBus()
	if err != nil {
		utils.LogErr(err)
	}

	omxPlayer := conn.Object("org.mpris.MediaPlayer2.omxplayer", dbus.ObjectPath("/org/mpris/MediaPlayer2"))
	return &OmxCtrl{conn: conn, omxPlayer: omxPlayer}
}

func (ctrl *OmxCtrl) Close() error {
	return ctrl.conn.Close()
}

func (ctrl *OmxCtrl) Action(action string) {
	actionCodes := map[string]int{
		"Exit":             15,
		"PlayPause":        16,
		"DecreaseVolume":   17,
		"IncreaseVolume":   18,
		"SeekBackSmall":    19,
		"SeekForwardSmall": 20,
		"SeekBackLarge":    21,
		"SeekForwardLarge": 22,
	}
	utils.LogErr(ctrl.omxPlayer.Call(playerInterface+".Action", 0, actionCodes[action]).Err)
}

func (ctrl *OmxCtrl) CanQuit() (status bool) {
	utils.LogErr(ctrl.omxPlayer.Call(propertyGetter, 0, rootInterface, "CanQuit").Store(&status))
	return status
}

func (ctrl *OmxCtrl) PlaybackStatus() (status string) {
	utils.LogErr(ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "PlaybackStatus").Store(&status))
	return status
}

func (ctrl *OmxCtrl) GetSource() (source string) {
	utils.LogErr(ctrl.omxPlayer.Call(playerInterface+".GetSource", 0).Store(&source))
	return source
}

func (ctrl *OmxCtrl) Duration() (duration int64) {
	utils.LogErr(ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "Duration").Store(&duration))
	return duration
}

func (ctrl *OmxCtrl) Position() (pos int64) {
	utils.LogErr(ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "Position").Store(&pos))
	return pos / microsecond
}

func (ctrl *OmxCtrl) SetPosition(pos int64) (res int64) {
	utils.LogErr(ctrl.omxPlayer.Call(playerInterface+".SetPosition", 0, dbus.ObjectPath("/not/used"), pos).Store(&res))
	return res
}

func (ctrl *OmxCtrl) Seek(offset int64) (res int64) {
	utils.LogErr(ctrl.omxPlayer.Call(playerInterface+".Seek", 0, offset*microsecond).Store(&res))
	return res
}

// Package plotberry is a plugin for bawt that reads Plotly graphs
package plotberry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/gopherworks/bawt"
	log "github.com/sirupsen/logrus"
)

type PlotBerry struct {
	bot        *bawt.Bot
	totalUsers int
	pingTime   time.Duration
	celebrated bool
}

type TotalUsers struct {
	Plotberries int `json:"plotberries"`
}

func init() {
	bawt.RegisterPlugin(&PlotBerry{})
}

func (plotberry *PlotBerry) InitPlugin(bot *bawt.Bot) {

	plotberry.bot = bot
	plotberry.celebrated = true
	plotberry.pingTime = 10 * time.Second
	plotberry.totalUsers = 100001

	statchan := make(chan TotalUsers, 100)

	go plotberry.launchWatcher(statchan)
	go plotberry.launchCounter(statchan)

	bot.Listen(&bawt.Listener{
		MessageHandlerFunc: plotberry.ChatHandler,
		Name:               "Plotberry",
		Description:        "Provides an interface for watching a Plotly chart",
	})
}

func (plotberry *PlotBerry) ChatHandler(listen *bawt.Listener, msg *bawt.Message) {
	if msg.MentionsMe && msg.Contains("how many user") {
		msg.Reply(fmt.Sprintf("We got %d users!", plotberry.totalUsers))
	}
	return
}

func GetPlotberry() (*TotalUsers, error) {

	var data TotalUsers

	resp, err := http.Get("https://plot.ly/v0/plotberries")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (plotberry *PlotBerry) launchWatcher(statchan chan TotalUsers) {

	for {
		time.Sleep(plotberry.pingTime)

		data, err := GetPlotberry()

		if err != nil {
			log.Print(err)
			continue
		}

		if data.Plotberries != plotberry.totalUsers {
			statchan <- *data
		}

		plotberry.totalUsers = data.Plotberries
	}
}

func (plotberry *PlotBerry) launchCounter(statchan chan TotalUsers) {

	finalcountdown := 100000

	for data := range statchan {
		totalUsers := data.Plotberries
		mod := math.Mod(float64(totalUsers), 50) == 0
		rem := finalcountdown - totalUsers

		if plotberry.celebrated {
			continue
		}

		if mod || (rem <= 10) {

			var msg string
			switch {
			case rem == 10:
				msg = fmt.Sprintf("@all %d users till the finalcountdown!", rem)
			case rem == 9:
				msg = fmt.Sprintf("%d users!", rem)
			case rem == 8:
				msg = fmt.Sprintf("and %d", rem)
			case rem == 7:
				msg = fmt.Sprintf("we're at %d users. %d users till Mimosa time!\n", totalUsers, rem)
			case rem == 6:
				msg = fmt.Sprintf("%d...", rem)
			case rem == 5:
				msg = fmt.Sprintf("@all %d users\n I'm a freaky proud robot!", rem)
			case rem == 4:
				msg = fmt.Sprintf("%d users till finalcountdown!", rem)
			case rem == 3:
				msg = fmt.Sprintf("%d... \n", rem)
			case rem == 2:
				msg = fmt.Sprintf("%d more! humpa humpa\n", rem)
			case rem == 1:
				plotberry.bot.SendToChannel(plotberry.bot.Config.GeneralChannel, fmt.Sprintf("%d users until 100000.\nYOU'RE ALL MAGIC!", rem))
				msg = "https://31.media.tumblr.com/3b74abfa367a3ed9a2cd753cd9018baa/tumblr_miul04oqog1qkp8xio1_400.gif"
			case rem <= 0:
				msg = fmt.Sprintf("@all FINALCOUNTDOWN!!!\n We're at %d user signups!!!!! My human compatriots, taking an idea to a product with 100,000 users is an achievement few will experience in their life times. Reflect, humans, on your hard work and celebrate this success. You deserve it, and remember, Plot On!", totalUsers)
				plotberry.celebrated = true
			default:
				msg = fmt.Sprintf("We are at %d total user signups!", totalUsers)
			}
			plotberry.bot.SendToChannel(plotberry.bot.Config.GeneralChannel, msg)
		}
	}
}

package standup

import (
	"fmt"
	"net/http"

	"github.com/gopherworks/bawt"
	"github.com/gorilla/mux"
)

func (standup *Standup) InitWebPlugin(bot *bawt.Bot, privRouter *mux.Router, pubRouter *mux.Router) {
	// privRouter.HandleFunc("/plugins/standup.json", func(w http.ResponseWriter, r *http.Request) {
	// 	data := struct {
	// 		Users []*StandupUser
	// 	}{
	// 		Users: make([]*StandupUser, 0),
	// 	}
	// 	for _, value := range *standup.data {
	// 		data.Users = append(data.Users, value)
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")

	// 	err := json.NewEncoder(w).Encode(data)
	// 	if err != nil {
	// 		webReportError(w, "Error encoding data", err)
	// 	}
	// })
}

func webReportError(w http.ResponseWriter, msg string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("%s\n\n%s\n", msg, err)))
}

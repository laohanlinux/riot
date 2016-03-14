package handler

import (
	"net/http"

	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/cluster"

	"github.com/laohanlinux/mux"
)

func AdminHandlerFunc(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		doGet(w, r)
	default:
	}
}

func doGet(w http.ResponseWriter, r *http.Request) int {
	vars := mux.Vars(r)
	cmdAdmin := vars["cmd"]
	switch cmdAdmin {
	case "node":
		leaderName := cluster.SingleCluster().Leader()
		w.Write([]byte(leaderName))
		return 200
	case "peer":
		r := cluster.SingleCluster()
		peers, err := r.PeerStorage.Peers()
		if err != nil {
			logger.Info(err)
		}
		for _, peer := range peers {
			w.Write([]byte(peer + "\r\n"))
		}
		return 404
	}
	return 0
}

package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/cluster"

	"github.com/laohanlinux/mux"
)

const (
	aErrOk    = iota
	aNetErr   // net work error
	aBytesErr // the operation content is invalid format
	aNoLeaderErr
	aUnkownErr    // unkowned error
	aUnkownCmdErr //
)

type ResponseMsg struct {
	Results interface{} `json:"results, omitempty"`
	ErrCode int         `json:"error, omitempty"`
	Time    float64     `json:"time,omitempty"`
	start   time.Time
}

func (msg *ResponseMsg) JsonToBytes(errCode ...int) []byte {
	if len(errCode) > 0 {
		msg.ErrCode = errCode[0]
	}
	b, _ := json.Marshal(msg)
	return b
}

func AdminHandlerFunc(w http.ResponseWriter, r *http.Request) {
	msg := ResponseMsg{
		ErrCode: 0,
		start:   time.Now(),
	}

	defer func(msg *ResponseMsg) {
		msg.Time = time.Now().Sub(msg.start).Seconds()
		w.Write(msg.JsonToBytes())
	}(&msg)

	var err error

	switch r.Method {
	case "GET":
		msg.ErrCode, msg.Results, err = doGet(w, r)
		if err != nil {
			logger.Error(err)
		}
	case "POST":
		msg.ErrCode, err = doPost(w, r)
		if err != nil {
			logger.Error(err)
		}
	default:
	}
}

func doGet(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	cmdAdmin := vars["cmd"]
	switch cmdAdmin {
	case "leader":
		leaderName := cluster.SingleCluster().Leader()
		if leaderName == "" {
			return aNoLeaderErr, nil, fmt.Errorf("No Leader in the cluser")
		}
		return aErrOk, leaderName, nil
	case "peer":
		r := cluster.SingleCluster()
		peers, err := r.PeerStorage.Peers()
		if err != nil {
			return aUnkownErr, nil, err
		}
		var peerStr []string
		for _, peer := range peers {
			peerStr = raft.AddUniquePeer(peerStr, peer)
		}
		return aErrOk, peerStr, nil
	case "status":
		status := cluster.SingleCluster().Status()
		return aErrOk, status, nil
	default:
		return aUnkownCmdErr, nil, fmt.Errorf("%s is unkown cmd", cmdAdmin)
	}
}

func doPost(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	cmdAdmin := vars["cmd"]
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return aNetErr, err
	}
	switch cmdAdmin {
	// let the node(addr:port) join the cluster
	case "join":
		// {"addr:": "", "port":""}
		var remoteAdd = make(map[string]string)
		if err := json.Unmarshal(b, &remoteAdd); err != nil {
			return aBytesErr, err
		}
		if len(remoteAdd) < 1 {
			return aBytesErr, fmt.Errorf("post body is invalid")
		}
		// 1. Get The Leader
		leaderName := cluster.SingleCluster().Leader()
		if leaderName == "" {
			return aNoLeaderErr, fmt.Errorf("No Leader In Cluster")
		}
		logger.Info("The Leader Name is :", leaderName)
		// 2. make sure the leader is itself
		if !strings.HasPrefix(leaderName, remoteAdd["addr"]) && remoteAdd["addr"] != "" {
			return aNoLeaderErr, nil
		}
		addr := remoteAdd["ip"] + ":" + remoteAdd["port"]
		// _, err := net.ResolveIPAddr("tcp", addr)
		// if err != nil {
		// 	return aBytesErr, err
		// }
		logger.Debug(addr, "will join the cluster, leader is :", leaderName)
		future := cluster.SingleCluster().R.AddPeer(addr)
		if err := future.Error(); err != nil {
			if err == raft.ErrKnownPeer {
				return aErrOk, nil
			}
			return aBytesErr, err
		}
	}

	return aErrOk, nil
}

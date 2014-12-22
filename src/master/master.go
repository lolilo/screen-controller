package main

import (
	"github.com/gorilla/mux"
	"log"
	"master/master"
	"master/master/slaveMapHandler"
	"master/master/slaveMonitor"
	"net/http"
	"network"
	"proxy"
	"website"
)

var (
	SLAVE_BINARY_PATH = network.GetRelativeFilePath("../../bin/slave")
)

func main() {
	proxy.Start()
	slaveMap := master.GetSlaveMap()
	router := mux.NewRouter()
	website.InitiateWebsiteHandlers(slaveMap, router)
	router.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	router.HandleFunc("/get_slave_binary", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, SLAVE_BINARY_PATH)
	})

	slaveMapHandler.InitiateSlaveMapHandler(router, slaveMap)

	http.Handle("/", router)
	go slaveMonitor.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

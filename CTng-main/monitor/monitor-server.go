package monitor

import (
	"CTng/gossip"
	//"CTng/crypto"
	//"CTng/util"
	//"bytes"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"time"
	"strings"
	//"strconv"
	"github.com/gorilla/mux"
)

const PROTOCOL = "http://"
func bindMonitorContext(context *MonitorContext, fn func(context *MonitorContext, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(context, w, r)
	}
}

func handleMonitorRequests(c *MonitorContext) {
	// MUX which routes HTTP directories to functions.
	gorillaRouter := mux.NewRouter().StrictSlash(true)
	// POST functions
	gorillaRouter.HandleFunc("/monitor/checkforme", bindMonitorContext(c, requestcheck)).Methods("POST")
	gorillaRouter.HandleFunc("/monitor/get-updates", bindMonitorContext(c, requestupdate)).Methods("POST")
	gorillaRouter.HandleFunc("/monitor/recieve-gossip", bindMonitorContext(c, handle_gossip)).Methods("POST")
	gorillaRouter.HandleFunc("/monitor/recieve-gossip-from-gossiper", bindMonitorContext(c, handle_gossip_from_gossiper)).Methods("POST")
	// Start the HTTP server.
	http.Handle("/", gorillaRouter)
	// Listen on port set by config until server is stopped.
	log.Fatal(http.ListenAndServe(":"+c.Config.Port, nil))
}

func StartMonitorServer(c *MonitorContext) {
	// Check if the storage file exists in this directory
	monitortype := 1
	fmt.Println("What type of Monitor would you like to use?")
	fmt.Println("1. Normal, Sync monitor")
	fmt.Println("2. 0 wait monitor for testing")
	fmt.Scanln(&monitortype)
	if monitortype == 1{
	time_wait := gossip.Getwaitingtime();
	time.Sleep(time.Duration(time_wait)*time.Second);
	}
	InitializeMonitorStorage(c)
	err := c.LoadStorage()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// Storage File doesn't exit. Create new, empty json file.
			err = c.SaveStorage()
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	tr := &http.Transport{}
	c.Client = &http.Client{
		Transport: tr,
	}
	// Run a go routine to handle tasks that must occur every MMD
	go PeriodicTasks(c)
	// Start HTTP server loop on the main thread
	handleMonitorRequests(c)
}
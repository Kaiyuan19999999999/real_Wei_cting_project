package monitor

import (
	"CTng/gossip"
	"CTng/crypto"
	"CTng/util"
	"bytes"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	//"log"
	"net/http"
	//"time"
	//"strings"
	"strconv"
	//"github.com/gorilla/mux"
)

type Clientupdate struct{
	STHs *gossip.Gossip_Storage
	REVs *gossip.Gossip_Storage
	PoMs *gossip.Gossip_Storage
	MonitorID string
	//Period here means the update period, the client udpate object can contain more information than just the period 
	Period string
	PoMsig string
}

type Clientquery struct{
	Client_URL string
	LastUpdatePeriod string
}

//This function should be invoked after the monitor-gossiper system converges in this period
func PrepareClientupdate(c *MonitorContext,LastUpdatePeriod string) Clientupdate{
	LastUpdatePeriodint, _ := strconv.Atoi(LastUpdatePeriod)
	CurrentPeriodint, _:= strconv.Atoi(gossip.GetCurrentPeriod())
	//intialize some storages
	storage_conflict_pom := new(gossip.Gossip_Storage)
	*storage_conflict_pom  = make(gossip.Gossip_Storage)
	storage_sth_full := new(gossip.Gossip_Storage)
	*storage_sth_full  = make(gossip.Gossip_Storage)
	storage_rev_full := new(gossip.Gossip_Storage)
	*storage_rev_full  = make(gossip.Gossip_Storage)
	//load all poms and sign on it
	for _, gossipObject := range *storage_conflict_pom{
		(*storage_conflict_pom)[gossipObject.GetID()] = gossipObject
	}
	payload,_ := json.Marshal(*storage_conflict_pom)
	signature, _ := crypto.RSASign([]byte(payload), &c.Config.Crypto.RSAPrivateKey, c.Config.Crypto.SelfID)
	//load all STHs (Fully Threshold signed) from lastUpdatePeriod to the current period
	for _, gossipObject := range *storage_sth_full{
		for i := LastUpdatePeriodint; i < CurrentPeriodint; i++ {
			if gossipObject.Period == strconv.Itoa(i){
				(*storage_sth_full)[gossipObject.GetID()] = gossipObject
			}
		}
	}
	//load all REVs (Fully Threshold signed) from LastUpdatePeriod to the current period
	for _, gossipObject := range *storage_rev_full{
		for i := LastUpdatePeriodint; i < CurrentPeriodint; i++ {
			if gossipObject.Period == strconv.Itoa(i){
				(*storage_rev_full)[gossipObject.GetID()] = gossipObject
			}
		}
	}
	CTupdate := Clientupdate{
		STHs: storage_sth_full,
		REVs: storage_rev_full,
		PoMs: storage_conflict_pom,
		MonitorID: c.Config.Signer,
		Period: gossip.GetCurrentPeriod(),
		PoMsig: signature.String(),
	}
	return CTupdate
}

func requestupdate(c *MonitorContext, w http.ResponseWriter, r *http.Request){
	var ticket Clientquery
	fmt.Println(util.GREEN+"Client ticket received"+util.RESET)
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var ctupdate = PrepareClientupdate(c,ticket.LastUpdatePeriod)
	fmt.Println(ctupdate.Period)
	msg, _ := json.Marshal(ctupdate)
	resp, postErr := c.Client.Post("http://"+ticket.Client_URL+"/receive-updates", "application/json", bytes.NewBuffer(msg))
	if postErr != nil {
		fmt.Println("Error sending update to client: " + postErr.Error())
	} else {
		// Close the response, mentioned by http.Post
		// Alernatively, we could return the response from this function.
		defer resp.Body.Close()
		if c.Verbose {
			fmt.Println("Client responded with " + resp.Status)
		}
		fmt.Println(util.GREEN+"Client update Sent"+util.RESET)
	}
}

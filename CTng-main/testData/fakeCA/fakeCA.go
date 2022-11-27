package fakeCA

import (
	//"CTng/GZip"
	"CTng/crypto"
	"CTng/gossip"
	"CTng/util"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"encoding/hex"
	"github.com/gorilla/mux"
	"strconv"
)

var CA_SIZE int

type CAConfig struct {
	Signer              crypto.CTngID
	Port                string
	NRevoke             int
	MRD                 int
	Private             rsa.PrivateKey
	CRVs                [][]byte //should be array of CRVs
	Day                 int      //I use int so I don't have to round and convert timestamps but that would be ideal
	MisbehaviorInterval int
}



type SRH struct {
	RootHash  string
	TreeSize  int
	Period string
}

type Revocation struct {
	Day       int
	Delta_CRV []byte
	SRH      SRH
}
//Caution: this file is plagued with Global Variables for conciseness.
var config CAConfig
var REVS []gossip.Gossip_object
var fakeREVs []gossip.Gossip_object
var request_count int
var currentPeriod int
var CAType int

func generateRevocation(CA CAConfig,catype int, Period_num int) gossip.Gossip_object{
	// Generate a random-ish SRH, add to SRHS.
	hashmsg := "Root Hash" + fmt.Sprint((Period_num+7)*(24+catype))
	hash, _ := crypto.GenerateSHA256([]byte(hashmsg))
	// Generate delta CRV and then compress it
	first_arr := CA.CRVs[CA.Day] //this assumes we never have CRV of len 0 (fresh CA)
	CA.Day += 1
	CA.CRVs[CA.Day] = make([]byte, CA_SIZE, CA_SIZE)

	var delta_crv = make([]byte, CA_SIZE, CA_SIZE)
	// Make the dCRV here by randomly flipping Config.NRevoke bits
	for i := 0; i < CA.NRevoke; i++ {
		change := rand.Intn(len(delta_crv))
		flip := byte(1)
		flip = flip << uint(rand.Intn(8))
		delta_crv[change] = flip
	}
	//fmt.Println(SRH1,delta_crv)
	// creates the new CRV from the old one+dCRV
	for i, _ := range first_arr {
		CA.CRVs[CA.Day][i] = first_arr[i] ^ delta_crv[i]
	} //this is scuffed/slow for giant CRVs O(n), also I am assuming CRVs are same size, can modify for different sizes
	REV := Revocation{
		Day:       CA.Day,
		Delta_CRV: delta_crv,
		SRH: SRH{
			RootHash:  hex.EncodeToString(hash),
			TreeSize:  currentPeriod * 12571285,
			Period: gossip.GetCurrentPeriod(),
		},
	}
	payload3, _ := json.Marshal(REV)
	payload := string(CA.Signer)+"CRV"+string(payload3)
	signature, _ := crypto.RSASign([]byte(payload), &config.Private, config.Signer)
	//fmt.Println(payload)
	//fmt.Println(signature)
	gossipREV := gossip.Gossip_object{
		Application: "CTng",
		Type:        gossip.REV,
		Period:      strconv.Itoa(Period_num),
		Signer:      string(config.Signer),
		Signature:   [2]string{signature.String(), ""},
		Crypto_Scheme: "RSA",
		Payload:     [3]string{string(CA.Signer),"CRV",string(payload3)},
	}
	return gossipREV
}

func periodicTasks() {
	// Queue the next tasks to occur at next MRD.
	time.AfterFunc(time.Duration(config.MRD)*time.Second, periodicTasks)
	// Generate CRV and SRH
	fmt.Println("CA Running Tasks at Period", gossip.GetCurrentPeriod())
	currentPeriod++
}


//Hard code to simulate a CA server that will generate subdomain to communicate
func requestREV(w http.ResponseWriter, r *http.Request) {
	REV_index,err := strconv.Atoi(gossip.GetCurrentPeriod())
	if err == nil{}
	    //Always disconnecting CA
		if CAType == 4{
			request_count++
			fmt.Println(util.RED,"Not sending Any REVS",util.RESET)
			return
		}
		//Disconnecting CA
		if CAType == 3 && request_count%config.MisbehaviorInterval == 0 {
			// No response or any bad request response should trigger the accusation
			request_count++
			fmt.Println(util.RED,"Not sending Any REVS",util.RESET)
			return
		}
		//Split-World CA
		if CAType == 2 && request_count%config.MisbehaviorInterval == 0{
			json.NewEncoder(w).Encode(fakeREVs[REV_index])
			request_count++
			fmt.Println(util.RED,"FakeREV sent.",fakeREVs[REV_index].GetID(),"sig: ", fakeREVs[REV_index].Signature[0],util.RESET)
			return
		}
		// Normal CA
		if err == nil{}
		json.NewEncoder(w).Encode(REVS[REV_index])
		fmt.Println(util.GREEN,"REV sent",REVS[REV_index].GetID(),"sig: ", REVS[REV_index].Signature[0], util.RESET)
		request_count++
}

func fill_with_data(){
	REVS = REVS[:0]
	fakeREVs = fakeREVs[:0]
	for i:=0; i<60; i++{
		rev1 := generateRevocation(config, 1, i)
		fakeREV1 := generateRevocation(config, CAType, i)
		REVS = append(REVS, rev1)
	    fakeREVs = append(fakeREVs, fakeREV1)
	}
}

func getCAType() {
	fmt.Println("What type of CA would you like to use?")
	fmt.Println("1. Normal, behaving CA (default)")
	fmt.Println("2. Split-World (Send fake REVs after ", config.MisbehaviorInterval, "requests)")
	fmt.Println("3. Disconnecting CA (unresponsive every", config.MisbehaviorInterval, "requests)")
	fmt.Println("4. Always Unreponsive CA")
	fmt.Scanln(&CAType)
}

// Runs a fake CA server with the ability to act roguely.
func RunFakeCA(configFile string) {
	// Global Variable initialization
	CA_SIZE = 1024
	CAType = 1
	currentPeriod = 0
	request_count = 0
	REVS = make([]gossip.Gossip_object, 0, 60)
	fakeREVs = make([]gossip.Gossip_object, 0, 60)
	// Read the config file
	config = CAConfig{}
	configBytes, err := util.ReadByte(configFile)
	if err != nil {
		fmt.Println("Error reading config file: ", err)
		return
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Println("Error reading config file: ", err)
	}
	config.CRVs = make([][]byte, 999, 999)
	config.CRVs[0] = make([]byte, CA_SIZE, CA_SIZE)
	config.Day = 0
	getCAType()
	fill_with_data()
	// MUX which routes HTTP directories to functions.
	gorillaRouter := mux.NewRouter().StrictSlash(true)
	gorillaRouter.HandleFunc("/ctng/v2/get-revocation", requestREV).Methods("GET")
	http.Handle("/", gorillaRouter)
	fmt.Println("Listening on port", config.Port)
	go periodicTasks()
	http.ListenAndServe(":"+config.Port, nil)
}

package gossip


import (
	"CTng/util"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"CTng/crypto"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"time"
	//"reflect"
	//"math/rand"
	//"math"
	"encoding/binary"
)

type Gossiper interface {
	// Response to entering the 'base page' of a gossiper.
	// TODO: Create informational landing page
	homePage()
	// HTTP POST request, receive a JSON object from another gossiper or connected monitor.
	// /gossip/push-data
	handleGossip(w http.ResponseWriter, r *http.Request)
	// Respond to HTTP GET request.
	// /gossip/get-data
	handleGossipObjectRequest(w http.ResponseWriter, r *http.Request)
	// Push JSON object to connected network from this gossiper via HTTP POST.
	// /gossip/gossip-data
	gossipData()
	// TODO: Push JSON object to connected 'owner' (monitor) from this gossiper via HTTP POST.
	// Sends to an owner's /monitor/recieve-gossip endpoint.
	sendToOwner()
	// Process JSON object received from HTTP POST requests.
	processData()
	// TODO: Erase stored data after one MMD.
	eraseData()
	// HTTP server function which handles GET and POST requests.
	handleRequests()
}

// Binds the context to the functions we pass to the router.
func bindContext(context *GossiperContext, fn func(context *GossiperContext, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(context, w, r)
	}
}

func handleRequests(c *GossiperContext) {
	// MUX which routes HTTP directories to functions.
	gorillaRouter := mux.NewRouter().StrictSlash(true)
	// homePage() is ran when base directory is accessed.
	gorillaRouter.HandleFunc("/gossip/", homePage)
	// Inter-gossiper endpoints
	gorillaRouter.HandleFunc("/gossip/push-data", bindContext(c, handleGossip)).Methods("POST")
	// Monitor interaction endpoint
	gorillaRouter.HandleFunc("/gossip/gossip-data", bindContext(c, handleGossip)).Methods("POST")
	// Start the HTTP server.
	http.Handle("/", gorillaRouter)
	fmt.Println(util.BLUE+"Listening on port:", c.Config.Port, util.RESET)
	err := http.ListenAndServe(":"+c.Config.Port, nil)
	// We wont get here unless there's an error.
	log.Fatal("ListenAndServe: ", err)
	os.Exit(1)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the base page for the CTng gossiper.")
}
// handleGossip() is ran when POST is recieved at /gossip/push-data.
// It should verify the Gossip object and then send it to the network.
func handleGossip(c *GossiperContext, w http.ResponseWriter, r *http.Request) {
	// Parse sent object.
	// Converts JSON passed in the body of a POST to a Gossip_object.
	var gossip_obj Gossip_object
	err := json.NewDecoder(r.Body).Decode(&gossip_obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Verify the object is valid, if invalid we just ignore it
	// CON do not have a signature on it yet 
	err = gossip_obj.Verify(c.Config.Crypto)
	if err != nil {
		//fmt.Println("Received invalid object "+TypeString(gossip_obj.Type)+" from " + util.GetSenderURL(r) + ".")
		fmt.Println(util.RED,"Received invalid object "+TypeString(gossip_obj.Type)+ " signed by " + gossip_obj.Signer+" aka "+EntityString(gossip_obj.Signer) + ".",util.RESET)
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	switch gossip_obj.Type{
	case STH, REV:
		//fmt.Println(util.GREEN,"Received Valid object "+TypeString(gossip_obj.Type)+ " signed by " + gossip_obj.Signer+" aka "+EntityString(gossip_obj.Signer) + ".",util.RESET)
		Handle_Sign_and_Gossip(c,gossip_obj)
	case ACC:
		//fmt.Println(util.GREEN,"Received Valid object "+TypeString(gossip_obj.Type)+ " signed by " + gossip_obj.Signer+" aka "+EntityString(gossip_obj.Signer) + ".",util.RESET)
		Handle_ACC(c,gossip_obj)
	case CON:
		fmt.Println(util.GREEN,"Received Valid object "+TypeString(gossip_obj.Type)+util.RESET)
		Handle_CON(c, gossip_obj)
	case STH_FRAG, REV_FRAG, ACC_FRAG, CON_FRAG:
		fmt.Println(util.GREEN,"Received Valid object "+TypeString(gossip_obj.Type)+ " signed by " + gossip_obj.Signer+" aka "+EntityString(gossip_obj.Signer) + ".",util.RESET)
		Handle_Frag(c,gossip_obj)
	case STH_FULL, REV_FULL, ACC_FULL, CON_FULL:
		Handle_FULL(c,gossip_obj)
	}
}

func Check_conflicts_and_poms(c *GossiperContext, g Gossip_object)bool{
	if c.HasPoM(g.Payload[0],g.Period){
		return true
	}
	if IsDuplicateFromGSC(g, c.Storage_RAW){
		stored_obj, _ := GetObjectFromGSC(g.Get_Counter_ID(),c.Storage_RAW)
		//fmt.Println(util.YELLOW,stored_obj.Signature[0], g.Signature[0])
		pom_gen := DetectConflicts(c,g,stored_obj)
		//true if there is pom
		return pom_gen
	}
	return false
}

func Handle_CON(c *GossiperContext, g Gossip_object){
	pom_err := c.Has_TSS_CON_POM(g.Payload[0], g.Period)
	if pom_err{
		return
	}
	if IsDuplicateFromGSC(g, c.Storage_RAW){
		//swap for gossiper sync
		hashmsg1 := g.Payload[0]+g.Payload[1]+g.Payload[2]
		hash1, _ := crypto.GenerateSHA256([]byte(hashmsg1))
		int1 := binary.BigEndian.Uint32(hash1)
		existing_obj := (*c.Storage_RAW)[g.Get_Counter_ID()]
		hashmsg2 := existing_obj.Payload[0]+existing_obj.Payload[1]+g.Payload[2]
		hash2, _ := crypto.GenerateSHA256([]byte(hashmsg2))
		int2 := binary.BigEndian.Uint32(hash2)
		//this means it is not a duplicate because the payload is different
		if int1 > int2{
			//store object should swap in the new object using the same Gossip ID as the key
			c.StoreObject(g)
			//send it to connected gossipers
			GossipData(c,g)
		}
	}else{
		c.StoreObject(g)
	}
	f := func() {
		sig_frag, err := c.Config.Crypto.ThresholdSign(g.Payload[0]+g.Payload[1]+g.Payload[2])
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(util.BLUE, "Signing CON against", g.Payload[0], util.RESET)
		g.Type = CON_FRAG
		g.Signature[0] = sig_frag.String()
		g.Signer = c.Config.Crypto.SelfID.String()
		g.Crypto_Scheme = "bls"
		c.StoreObject(g)
		Process_TSS_Object(c,g,CON_FULL)
		GossipData(c,g)
	}
	// Delay the calling of f until gossip_wait_time has passed.
	time.AfterFunc(time.Duration(c.Config.Public.Gossip_wait_time)*time.Second, f)
}
func Handle_Sign_and_Gossip(c *GossiperContext, g Gossip_object){
	if g.Signature[0] == (*c.Storage_RAW)[g.Get_Counter_ID()].Signature[0]{
		return
	}
	if Check_conflicts_and_poms(c,g){
		return 
	}
	c.StoreObject(g)
	GossipData(c,g)
	//This handles the STHS
	if g.Type == STH {
		// The below function for creates the STH_FRAG object after Gossip_wait_time
		f := func() {
			if Check_conflicts_and_poms(c,g){
				return
			}
			pom_err := c.HasPoM(g.Signer, g.Period)
			//if there is no conflicting information/PoM send the Threshold signed version to the gossiper
			sig_frag, err := c.Config.Crypto.ThresholdSign(g.Payload[0]+g.Payload[1]+g.Payload[2])
			if err != nil {
				fmt.Println(err.Error())
			}
			if pom_err == false {
				fmt.Println(util.BLUE, "Signing STH of", g.Signer, util.RESET)
				g.Type = STH_FRAG
				g.Signature[0] = sig_frag.String()
				g.Signer = c.Config.Crypto.SelfID.String()
				g.Crypto_Scheme = "bls"
				c.StoreObject(g)
				Process_TSS_Object(c,g,STH_FULL)
				GossipData(c,g)
			} else {
				fmt.Println(util.RED, "Conflicting information/PoM found, not sending STH_FRAG", util.RESET)
			}

		}
		// Delay the calling of f until gossip_wait_time has passed.
		time.AfterFunc(time.Duration(c.Config.Public.Gossip_wait_time)*time.Second, f)
		return
	}
	//if the object is from a CA, revocation information
	//this handles revocation information
	if  g.Type == REV{
		f := func() {
			if Check_conflicts_and_poms(c,g){
				return
			}
			fmt.Println(util.BLUE, "Signing Revocation of", g.Signer, util.RESET)
			pom_err := c.HasPoM(g.Signer, g.Period)
			if pom_err == false {
				sig_frag, err := c.Config.Crypto.ThresholdSign(g.Payload[0]+g.Payload[1]+g.Payload[2])
				if err != nil {
					fmt.Println(err.Error())
				}
				g.Type = REV_FRAG
				g.Signature[0] = sig_frag.String()
				g.Signer = c.Config.Crypto.SelfID.String()
				g.Crypto_Scheme = "bls"
				c.StoreObject(g)
				Process_TSS_Object(c,g,REV_FULL)
				GossipData(c,g)
			}
		}
		time.AfterFunc(time.Duration(c.Config.Public.Gossip_wait_time)*time.Second, f)
		return
	}
}

func Handle_ACC(c* GossiperContext, g Gossip_object){
	if IsDuplicateFromGSC(g, c.Storage_RAW){
		return
	}
	//fmt.Println(util.GREEN,"ACC is not a duplicate",util.RESET)
	if Check_conflicts_and_poms(c,g){
		return 
	}
	//fmt.Println(util.GREEN, "No PoM against", g.Payload[0], util.RESET)
	c.StoreObject(g)
	GossipData(c,g)
	//fmt.Println(util.GREEN, "ACC stored", util.RESET)
	//Check if there is already an entry in ACC DB
	obj := (*c.ACC_DB)[g.GetID()] 
	//this means this is the first accusation against this entity
	if obj == nil{
		new_signers := []string{}
		new_signers = append(new_signers, g.Signer)
		new_sigs := []string{}
		new_sigs = append(new_sigs, g.Signature[0])
		new_counter := PoM_PreTSS_Counter{
			Signers: new_signers,
			Sigs: new_sigs,
			Num: 1,
		}
		(*c.ACC_DB)[g.GetID()] = &new_counter
		//fmt.Println(util.GREEN,"First accusation against ",g.Payload[0], " processed.", util.RESET)
	}else{
		obj.Num++
		(*c.ACC_DB)[g.GetID()] = obj
		if obj.Num>=c.Config.Crypto.Threshold{
			payload := g.Payload
			sig, err := c.Config.Crypto.ThresholdSign(payload[0]+payload[1]+payload[2])
			if err != nil{
				fmt.Println("Threshold Sign failed.")
				return
			}
			sig_field := [2]string{sig.String(), ""}
			acc_frag := Gossip_object{
				Application: g.Application,
				Type:        ACC_FRAG,
				Period:      GetCurrentPeriod(),
				Signer:      c.Config.Crypto.SelfID.String(),
				Timestamp:   GetCurrentTimestamp(),
				Signature:   sig_field,
				Crypto_Scheme: "BLS",
				Payload:     payload,
			}
			c.StoreObject(acc_frag)
			Process_TSS_Object(c, acc_frag, ACC_FULL)
			GossipData(c,acc_frag)
		}
	}
}


func Handle_Frag(c* GossiperContext, g Gossip_object){
	if IsDuplicateFromGSC(g, c.Storage_FRAG){
		return
	}
	if Check_conflicts_and_poms(c,g) && g.Type != CON_FRAG{
		return 
	}
	switch g.Type{
	case STH_FRAG:
		c.StoreObject(g)
		Process_TSS_Object(c, g, STH_FULL)
		GossipData(c,g)
	case REV_FRAG:
		c.StoreObject(g)
		Process_TSS_Object(c, g, REV_FULL)
		GossipData(c,g)
	case ACC_FRAG:
		c.StoreObject(g)
		Process_TSS_Object(c, g, ACC_FULL)
		GossipData(c,g)
	case CON_FRAG:
		fmt.Println(util.GREEN, "Storing CON_FRAG Signed by ", g.Signer, util.RESET)
		c.StoreObject(g)
		Process_TSS_Object(c, g, CON_FULL)
		GossipData(c,g)
	}
}

func Handle_FULL(c* GossiperContext, g Gossip_object){
	switch g.Type{
	case STH_FULL,REV_FULL:
		if IsDuplicateFromGS(g, c.Storage_FULL){
			return
		}else{
			c.StoreObject(g)
			GossipData(c,g)
		}
	case CON_FULL:
		if IsDuplicateFromGS(g, c.Storage_POM){
			return
		}else{
			c.StoreObject(g)
			GossipData(c,g)
		}
	case ACC_FULL:
		if IsDuplicateFromGS(g, c.Storage_POM_TEMP){
			return
		}else{
			c.StoreObject(g)
			GossipData(c,g)
		}
	}


}
// Sends a gossip object to all connected gossipers.
// This function assumes you are passing valid data. ALWAYS CHECK BEFORE CALLING THIS FUNCTION.
func GossipData(c *GossiperContext, gossip_obj Gossip_object) error {
	//desync
	/*
	rand.Seed(time.Now().UnixNano())
    n := rand.Intn(100) // n will be between 0 and 10
    time.Sleep(time.Duration(n)*time.Millisecond)
	*/
	// Convert gossip object to JSON
	msg, err := json.Marshal(gossip_obj)
	if err != nil {
		fmt.Println(err)
	}

	// Send the gossip object to all connected gossipers.
	for _, url := range c.Config.Connected_Gossipers {
		//fmt.Println("Attempting to sending data to", url)
		// HTTP POST the data to the url or IP address.
		resp, err := c.Client.Post("http://"+url+"/gossip/push-data", "application/json", bytes.NewBuffer(msg))
		if err != nil {
			if strings.Contains(err.Error(), "Client.Timeout") ||
				strings.Contains(err.Error(), "connection refused") {
				fmt.Println(util.RED+"Connection failed to "+url+".", util.RESET)
				// Don't accuse gossipers for inactivity.
				// defer Accuse(c, url)
			} else {
				fmt.Println(util.RED+err.Error(), "sending to "+url+".", util.RESET)
			}
			continue
		}
		// Close the response, mentioned by http.Post
		// Alernatively, we could return the response from this function.
		defer resp.Body.Close()
		//fmt.Println("Gossiped to " + url + " and recieved " + resp.Status)
	}
	return nil
}

// Sends a gossip object to the owner of the gossiper.
func SendToOwner(c *GossiperContext, obj Gossip_object) {
	// Convert gossip object to JSON
	msg, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	// Send the gossip object to the owner.
	resp, postErr := c.Client.Post("http://"+c.Config.Owner_URL+"/monitor/recieve-gossip-from-gossiper", "application/json", bytes.NewBuffer(msg))
	if postErr != nil {
		fmt.Println("Error sending object to owner: " + postErr.Error())
	} else {
		// Close the response, mentioned by http.Post
		// Alernatively, we could return the response from this function.
		defer resp.Body.Close()
		if c.Verbose {
			fmt.Println("Owner responded with " + resp.Status)
		}
	}
}

// Process a valid gossip object which is a duplicate to another one.
// If the signature/payload is identical, then we can safely ignore the duplicate.
// Otherwise, we generate a PoM for two objects sent in the same period.
func DetectConflicts(c *GossiperContext, obj Gossip_object, dup Gossip_object) bool{
	//If the object has PoM already, it is dead already
	//if c.HasPoM(obj.Payload[0],obj.Period){
		//return nil
	//}
	//If the object type is the same
	//In the same Periord
	//Signed by the same Entity
	//But the signature is different
	//MALICIOUS, you are exposed
	//note PoMs can have different signatures
	//fmt.Println(util.YELLOW, "Trying to detect conflicts", TypeString(obj.Type), TypeString(dup.Type), obj.Signature[0] == dup.Signature[0],util.RESET)
	if obj.Type == dup.Type && obj.Period == dup.Period && obj.Signer == dup.Signer && obj.Signature[0] != dup.Signature[0]{
		D2_POM:= Gossip_object{
			Application: obj.Application,
			Type:        CON,
			Period:      GetCurrentPeriod(),
			Signer:      "",
			Timestamp:   GetCurrentTimestamp(),
			Signature:   [2]string{obj.Signature[0], dup.Signature[0]},
			Payload:     [3]string{obj.Signer, obj.Payload[0]+obj.Payload[1]+obj.Payload[2],dup.Payload[0]+dup.Payload[1]+dup.Payload[2]},
		}
		_, sigerr1 := crypto.RSASigFromString(D2_POM.Signature[0])
		_, sigerr2 := crypto.RSASigFromString(D2_POM.Signature[1])
		fmt.Println(util.YELLOW, sigerr1, sigerr2, util.RESET)
		//store the object and send to monitor
		fmt.Println(util.YELLOW, "Entity: ", D2_POM.Payload[0], " is Malicious!", util.RESET)
		//SendToOwner(c,D2_POM)
		c.StoreObject(D2_POM)
		GossipData(c,D2_POM)
		f := func() {
			sig_frag, err := c.Config.Crypto.ThresholdSign(D2_POM.Payload[0]+D2_POM.Payload[1]+D2_POM.Payload[2])
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(util.BLUE, "Signing CON against", D2_POM.Payload[0], util.RESET)
			D2_POM.Type = CON_FRAG
			D2_POM.Signature[0] = sig_frag.String()
			D2_POM.Signer = c.Config.Crypto.SelfID.String()
			D2_POM.Crypto_Scheme = "bls"
			fmt.Println(util.GREEN, "Storing CON_FRAG Signed by ", D2_POM.Signer, util.RESET)
			c.StoreObject(D2_POM)
			GossipData(c,D2_POM)
			Process_TSS_Object(c,D2_POM,CON_FULL)
		}
		// Delay the calling of f until gossip_wait_time has passed.
		time.AfterFunc(time.Duration(c.Config.Public.Gossip_wait_time)*time.Second, f)
		return true
	}
	return false
}

func Process_TSS_Object(gc *GossiperContext, new_obj Gossip_object, target_type string) error{
	c := gc.Config.Crypto
	key := new_obj.GetID()
	//fmt.Println(key)
	newkey:=Gossip_ID{
		Period: key.Period,
		Type: target_type,
		Entity_URL: key.Entity_URL,
	}
	p_sig, err := crypto.SigFragmentFromString(new_obj.Signature[0])
	if err != nil {
		fmt.Println("partial sig conversion error (from string)")
		return err
	}
	//If there is already an TSS Object
	if _, ok:= (*gc.Storage_FULL)[newkey]; ok{
		fmt.Println(util.BLUE + "There already exists a "+ TypeString(target_type)+ " Object" + util.RESET)
		return nil
	} 
	//If there isn't a STH_FULL Object yet, but there exists some other sth_frag
	if val, ok := (*gc.Obj_TSS_DB)[key]; ok {
		val.Signers[val.Num] = new_obj.Signer
		if err != nil {
			fmt.Println("partial sig conversion error (from string)")
			return err
		}
		val.Partial_sigs[val.Num] = p_sig
		val.Num = val.Num + 1
		//fmt.Println("Finished updating Counters, the new number is", val.Num)
		//now we check if the number of sigs have reached the threshold
		if val.Num>=c.Threshold{
			TSS_sig, _ := c.ThresholdAggregate(val.Partial_sigs)
			TSS_sig_string,_ := TSS_sig.String()
			sigfield := new([2]string)
			(*sigfield)[0] = TSS_sig_string
			signermap := make(map[int]string)
			for i := 0; i<c.Threshold; i++{
				signermap[i] = val.Signers[i]
			}
			TSS_period := "0"
			//Set CON_FULL Period Number to 0 so that the monitor can search for it
			//Also this is fine because we only need 1 conflict_pom for each convicted entity
			if target_type != CON_FULL{
				TSS_period = new_obj.Period
			}
			TSS_FULL_obj := Gossip_object{
				Application: new_obj.Application,
				Type:        target_type,
				Period:      TSS_period,
				Signer:      "",
				Signers:     signermap,
				Timestamp:   GetCurrentTimestamp(),
				Signature:   *sigfield,
				Crypto_Scheme: "BLS",
				Payload:     new_obj.Payload,
			}
			//Store the POM
			fmt.Println(util.BLUE+TypeString(target_type)+" generated and Stored"+util.RESET)
			gc.StoreObject(TSS_FULL_obj)
			GossipData(gc,TSS_FULL_obj)
			//send to the monitor
			SendToOwner(gc,TSS_FULL_obj)
			return nil
		}
	}
	//if the this is the first STH_FRAG received
	//fmt.Println("This is the first partial sig registered")
	new_counter := new(Entity_Gossip_Object_TSS_Counter)
	*new_counter = Entity_Gossip_Object_TSS_Counter{
		Signers:     []string{new_obj.Signer,""},
		Num:      1,
		Partial_sigs: []crypto.SigFragment{p_sig,p_sig},
	}
	(*gc.Obj_TSS_DB)[key] = new_counter
	//fmt.Println("Number of counters in TSS DB is: ", len(*gc.Obj_TSS_DB))
	return nil

}

func PeriodicTasks(c *GossiperContext) {
	// Immediately queue up the next task to run at next MMD.
	// Doing this first means: no matter how long the rest of the function takes,
	// the next call will always occur after the correct amount of time.
	f := func() {
		PeriodicTasks(c)
	}
	time.AfterFunc(time.Duration(c.Config.Public.MMD)*time.Second, f)
	fmt.Println("Gossiper running on Period: ", GetCurrentPeriod())
	c.SaveStorage()
	c.WipeStorage()
}

func InitializeGossiperStorage (c* GossiperContext){
	c.StorageDirectory = "testData/gossiperdata/"+c.StorageID+"/"
	c.StorageFile = "GossipStorage.Json"
	util.CreateFile(c.StorageDirectory+c.StorageFile)

}

func StartGossiperServer(c *GossiperContext) {
	// Check if the storage file exists in this directory
	InitializeGossiperStorage(c)
	// Create the http client to be used.
	tr := &http.Transport{}
	c.Client = &http.Client{
		Transport: tr,
	}
	// HTTP Server Loop
	go PeriodicTasks(c)
	handleRequests(c)
}

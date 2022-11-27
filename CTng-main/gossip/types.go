package gossip

import (
	"CTng/config"
	"CTng/util"
	"CTng/crypto"
	//"encoding/json"
	"net/http"	
	"reflect"
	"fmt"
	"sync"
)


// Types of errors that can occur when parsing a Gossip_object
const (
	No_Sig_Match = "Signatures don't match"
	Mislabel     = "Fields mislabeled"
	Invalid_Type = "Invalid Type"
)

type Entity_Gossip_Object_TSS_Counter struct {
	Signers     []string
	Partial_sigs []crypto.SigFragment
	Num     int
}


type PoM_PreTSS_Counter struct{
	Signers []string //monitors
	Sigs []string //RSA sigs
	Num int   //the counter
}

//This DB stores the counter object for  sth_frags. rec_frags and acc_frags
type Gossip_Object_TSS_DB map[Gossip_ID]*Entity_Gossip_Object_TSS_Counter
//This DB stores accusations and Conflict PoMs
type Accusation_DB map[Gossip_ID]*PoM_PreTSS_Counter
type Conflict_DB map[Gossip_ID]*PoM_PreTSS_Counter

// Gossiper Context
// Ths type represents the current state of a gossiper HTTP server.
// This is the state of a gossiper server. It contains:
// The gossiper Configuration,
// Storage utilized by the gossiper,
// Any objects needed throughout the gossiper's lifetime (such as the http client).
type GossiperContext struct {
	Config      *config.Gossiper_config
	StorageID   string
	//RWlock
	RWlock sync.RWMutex
	//STH + REV + ACC
	Storage_RAW  *Gossip_Storage_Counter
	//STH_FRAG + REV_FRAG + ACC_FRAG + CON_FRAG
	Storage_FRAG *Gossip_Storage_Counter
	//STH_FULL + REV_FULL
	Storage_FULL *Gossip_Storage
	//CON + ACC_FULL
	Storage_POM_TEMP *Gossip_Storage
	//CON_FULL 
	Storage_POM *Gossip_Storage
	//ACC_FRAG counter + CON_FRAG counter 
	Obj_TSS_DB *Gossip_Object_TSS_DB
	//ACC Counter
	ACC_DB *Accusation_DB
	//Con Counter
	CON_DB *Conflict_DB
	//log
	G_log *Gossiper_log
	//File I/O
	StorageFile string 
	StorageDirectory string
	Client  *http.Client
	Verbose bool
}

type GossiperLogEntry struct{	
	Gossiper_URL string
	Period string
	Num_sth int
	Num_rev int
	Num_acc int
	Num_con int
	Num_sth_frag int
	Num_rev_frag int
	Num_acc_frag int
	Num_con_frag int
	Num_STH_FULL int
	Num_REV_FULL int
	Num_ACC_FULL int
	Num_CON_FULL int
}

type Gossiper_log map[string]GossiperLogEntry

func Gossip_Context_Init(config *config.Gossiper_config,Storage_ID string) *GossiperContext{
	storage_raw := new(Gossip_Storage_Counter)
	*storage_raw = make(Gossip_Storage_Counter)
	storage_frag := new(Gossip_Storage_Counter)
	*storage_frag = make(Gossip_Storage_Counter)
	storage_full := new(Gossip_Storage)
	*storage_full = make(Gossip_Storage)
	storage_pom := new(Gossip_Storage)
	*storage_pom = make(Gossip_Storage)
	storage_pom_temp := new(Gossip_Storage)
	*storage_pom_temp = make(Gossip_Storage)
	gossip_object_TSS_DB := new(Gossip_Object_TSS_DB)
	*gossip_object_TSS_DB = make(Gossip_Object_TSS_DB)
	accusation_db := new(Accusation_DB)
	*accusation_db = make(Accusation_DB)
	conflict_db := new(Conflict_DB)
	*conflict_db = make(Conflict_DB)
	g_log := new(Gossiper_log)
	*g_log = make(Gossiper_log) 
	ctx := GossiperContext{
		Config:      config,
		RWlock: sync.RWMutex{},
		//STH + REV + ACC + CON
		Storage_RAW:  storage_raw,
		//STH_FRAG + REV_FRAG + ACC_FRAG + CON_FRAG
		Storage_FRAG: storage_frag,
		//STH_FULL + REV_FULL + ACC_FULL + CON_FULL
		Storage_FULL: storage_full,
		//CON_FRAG + ACC_FULL + CON_FULL 
		Storage_POM: storage_pom,
		Storage_POM_TEMP: storage_pom_temp,
		//ACC_FRAG counter + CON_FRAG counter 
		Obj_TSS_DB: gossip_object_TSS_DB,
		//ACC Counter
		ACC_DB: accusation_db,
		CON_DB: conflict_db,
		G_log: g_log,
		StorageFile: "gossiper_data.json", // could be a parameter in the future.
		StorageID:   Storage_ID,
	}
	return &ctx
}
func CountStorageCounter(gs *Gossip_Storage_Counter, entry *GossiperLogEntry){
	for key,_ := range *gs{
		switch key.Type{
		case STH:
			entry.Num_sth++
		case REV:
			entry.Num_rev++
		case ACC:
			entry.Num_acc++
		case CON:
			entry.Num_con++
		case STH_FRAG:
			entry.Num_sth_frag++
		case REV_FRAG:
			entry.Num_rev_frag++
		case ACC_FRAG:
			entry.Num_acc_frag++
		case CON_FRAG:
			entry.Num_con_frag++
		case STH_FULL:
			entry.Num_STH_FULL++
		case REV_FULL:
			entry.Num_REV_FULL++
		case ACC_FULL:
			entry.Num_ACC_FULL++
		case CON_FULL:
			entry.Num_CON_FULL++
		}
	}
}

func CountStorage(gs *Gossip_Storage, entry *GossiperLogEntry){
	for key,_ := range *gs{
		switch key.Type{
		case STH:
			entry.Num_sth++
		case REV:
			entry.Num_rev++
		case ACC:
			entry.Num_acc++
		case CON:
			entry.Num_con++
		case STH_FRAG:
			entry.Num_sth_frag++
		case REV_FRAG:
			entry.Num_rev_frag++
		case ACC_FRAG:
			entry.Num_acc_frag++
		case CON_FRAG:
			entry.Num_con_frag++
		case STH_FULL:
			entry.Num_STH_FULL++
		case REV_FULL:
			entry.Num_REV_FULL++
		case ACC_FULL:
			entry.Num_ACC_FULL++
		case CON_FULL:
			entry.Num_CON_FULL++
		}
	}
}
// Saves the Storage object to the value in c.StorageFile.
func (c *GossiperContext) SaveStorage() error {
	newentry := GossiperLogEntry{
		Gossiper_URL: c.Config.Crypto.SelfID.String(),
		Period: GetCurrentPeriod(),
		Num_sth: 0,
		Num_rev: 0,
		Num_acc: 0,
		Num_sth_frag: 0,
		Num_rev_frag: 0,
		Num_acc_frag: 0,
		Num_con_frag: 0,
		Num_STH_FULL: 0,
		Num_REV_FULL: 0,
		Num_ACC_FULL: 0,
		Num_CON_FULL: 0,
	}
	CountStorageCounter((*c).Storage_RAW,&newentry)
	CountStorageCounter((*c).Storage_FRAG,&newentry)
	CountStorage((*c).Storage_FULL,&newentry)
	CountStorage((*c).Storage_POM_TEMP,&newentry)
	CountStorage((*c).Storage_POM,&newentry)
	(*c.G_log)[GetCurrentPeriod()] = newentry
	err := util.WriteData(c.StorageDirectory+"/"+c.StorageFile, c.G_log)
	return err
}
func WipeOneGS(g *Gossip_Storage){
	for key, _ := range *g{
		if  key.Period!=GetCurrentPeriod(){
			delete(*g,key)
		}
	}
}
func WipeOneGSC(g *Gossip_Storage_Counter){
	for key, _ := range *g{
		if  key.Period!=GetCurrentPeriod(){
			delete(*g,key)
		}
	}
}
//wipe all temp data
func (c *GossiperContext) WipeStorage(){
	WipeOneGSC(c.Storage_RAW)
	WipeOneGSC(c.Storage_FRAG)
	WipeOneGS(c.Storage_POM_TEMP)
	for key, _:= range *c.Obj_TSS_DB{
		if key.Period!=GetCurrentPeriod(){
			delete(*c.Obj_TSS_DB,key)
		}
	}
	for key, _:= range *c.ACC_DB{
		if key.Period!=GetCurrentPeriod(){
			delete(*c.ACC_DB,key)
		}
	}
	fmt.Println(util.BLUE,"Temp storage has been wiped.",util.RESET)
}
// Read every gossip object from c.StorageFile.
// Store all files in c.Storage by their ID.
/*
func (c *GossiperContext) LoadStorage() error {
	// Get the array that has been written to the storagefile.
	storageList := []Gossip_object{}
	bytes, err := util.ReadByte(c.StorageFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &storageList)
	if err != nil {
		return err
	}
	// Store the objects by their ID, based on the current period defined in the gossiper context.
	// Note that if the period changes (particularly, increases) between loads of the gossiper, some objects may be overwritten/lost.
	// So careful!
	for _, gossipObject := range storageList {
		(*c.Storage)[gossipObject.GetID()] = gossipObject
	}
	return nil
}
*/

// Stores an object in storage by its ID. Note that the ID utilizes Config.Public.Period_interval.
/*
func (c *GossiperContext) StoreObject(o Gossip_object) {
	(*c.Storage)[o.GetID()] = o
}*/

func (c *GossiperContext) StoreObject(o Gossip_object){
	switch o.Type{
	case STH,REV,ACC:
		c.RWlock.Lock()
		(*c.Storage_RAW)[o.Get_Counter_ID()] = o
		c.RWlock.Unlock()
	case CON:
		c.RWlock.Lock()
		(*c.Storage_RAW)[o.Get_Counter_ID()] = o
		(*c.Storage_POM_TEMP)[o.GetID()] = o
		c.RWlock.Unlock()
	case STH_FRAG,REV_FRAG,ACC_FRAG:
		c.RWlock.Lock()
		(*c.Storage_FRAG)[o.Get_Counter_ID()] = o
		c.RWlock.Unlock()
	case CON_FRAG:
		c.RWlock.Lock()
		(*c.Storage_FRAG)[o.Get_Counter_ID()] = o
		c.RWlock.Unlock()
	case STH_FULL,REV_FULL:
		c.RWlock.Lock()
		(*c.Storage_FULL)[o.GetID()] = o
		c.RWlock.Unlock()
	case ACC_FULL:
		c.RWlock.Lock()
		(*c.Storage_POM_TEMP)[o.GetID()] = o
		c.RWlock.Unlock()
	case CON_FULL:
		c.RWlock.Lock()
		(*c.Storage_POM)[o.GetID()] = o
		c.RWlock.Unlock()
	}
}

// Returns 2 fields: the object, and whether or not the object was successfully found.
// If the object isn't found then all fields of the Gossip_object will also be empty.
// WARNING: only checks the given storage
func GetObjectFromGS(id Gossip_ID, g *Gossip_Storage) (Gossip_object, bool) {
	obj := (*g)[id]
	if reflect.DeepEqual(obj, Gossip_object{}) {
		return obj, false
	}
	return obj, true
}

// Given a gossip object, check if the an object with the same ID exists in the storage.
//WARNING: only checks the given storage
func IsDuplicateFromGS(g Gossip_object, gs *Gossip_Storage) bool {
	id := g.GetID()
	_, exists := GetObjectFromGS(id, gs)
	return exists
}

// If the object isn't found then all fields of the Gossip_object will also be empty.
// WARNING: only checks the given storage
func GetObjectFromGSC(id Gossip_Counter_ID, g *Gossip_Storage_Counter) (Gossip_object, bool) {
	obj := (*g)[id]
	if reflect.DeepEqual(obj, Gossip_object{}) {
		return obj, false
	}
	return obj, true
}

// Given a gossip object, check if the an object with the same ID exists in the storage.
//WARNING: only checks the given storage
func IsDuplicateFromGSC(g Gossip_object, gs *Gossip_Storage_Counter) bool {
	id := g.Get_Counter_ID()
	_, exists := GetObjectFromGSC(id, gs)
	return exists
}

func (c *GossiperContext) Has_TSS_CON_POM(entity_URL string, period string) bool{
	ID := Gossip_ID{
		Period: "0",
		Type: CON_FULL,
		Entity_URL: entity_URL,
	}
	if _, ok := (*c.Storage_POM)[ID]; ok {
		return true
	}
	return false
}

func (c *GossiperContext) HasPoM(entity_URL string, period string) bool{
	//first check accusation pom
	ID := Gossip_ID{
		Period : period,
		Type : ACC_FULL,
		Entity_URL : entity_URL,
	}
	if _, ok := (*c.Storage_POM_TEMP)[ID]; ok {
		return true
	}
	//Then check Unsigned Conflict PoM generated this period
	ID = Gossip_ID{
		Period : period,
		Type : CON,
		Entity_URL : entity_URL,
	}
	if _, ok := (*c.Storage_POM_TEMP)[ID]; ok {
		return true
	}
	//Then check Partially Signed Conflict PoM generated this period
	ID = Gossip_ID{
		Period : period,
		Type : CON_FRAG,
		Entity_URL : entity_URL,
	}
	if _, ok := (*c.Storage_POM_TEMP)[ID]; ok {
		return true
	}
	//Finally check Threshold signed PoMs (from this period and from the past)
	ID = Gossip_ID{
		Period: "0",
		Type: CON_FULL,
		Entity_URL: entity_URL,
	}
	if _, ok := (*c.Storage_POM)[ID]; ok {
		return true
	}
	return false
}

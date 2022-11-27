package monitor

import (
	"CTng/config"
	"CTng/gossip"
	"CTng/util"
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
	"reflect"
)

type MonitorContext struct {
	Config      *config.Monitor_config
	Storage_TEMP *gossip.Gossip_Storage
	// Gossip objects from the gossiper will be assigned to their dedicated storage
	Storage_CONFLICT_POM *gossip.Gossip_Storage
	Storage_ACCUSATION_POM *gossip.Gossip_Storage
	Storage_STH_FULL *gossip.Gossip_Storage
	Storage_REV_FULL *gossip.Gossip_Storage
	// Utilize Storage directory: A folder for the files of each MMD.
	// Folder should be set to the current MMD "Period" String upon initialization.
	StorageDirectory string
	StorageID string
	StorageFile_CONFLICT_POM string
	StorageFile_ACCUSATION_POM string
	StorageFile_STH_FULL string
	StorageFile_REV_FULL string
	// The below could be used to prevent a Monitor from sending duplicate Accusations,
	// Currently, if a monitor accuses two entities in the same Period, it will trigger a gossip PoM.
	// Therefore, a monitor can only accuse once per Period. I believe this is a temporary solution.
	Verbose    bool
	Client     *http.Client
}

func (c *MonitorContext) Clean_Conflicting_Object(){
	GID := gossip.Gossip_ID{}
	for key, _ := range *c.Storage_STH_FULL{
		GID = gossip.Gossip_ID{
			Period: "0",
			Type: gossip.CON_FULL,
			Entity_URL: key.Entity_URL,
		}
		if _,ok := (*c.Storage_CONFLICT_POM)[GID]; ok{
			fmt.Println(util.BLUE+"Logger: "+ key.Entity_URL + "has Conflict_PoM on file, cleared the STH from this Logger this MMD"+util.RESET)
			delete(*c.Storage_STH_FULL,key)
		}
	}
	for key, _ := range *c.Storage_REV_FULL{
		GID = gossip.Gossip_ID{
			Period: "0",
			Type: gossip.CON_FULL,
			Entity_URL: key.Entity_URL,
		}
		if _,ok := (*c.Storage_CONFLICT_POM)[GID]; ok{
			fmt.Println(util.BLUE+"CA: "+ key.Entity_URL + "has Conflict_PoM on file, cleared the REV from this CA this MRD"+util.RESET)
			delete(*c.Storage_REV_FULL,key)
		}
	}
}

func (c *MonitorContext) SaveStorage() error{
	storageList_conflict_pom := []gossip.Gossip_object{}
    storageList_accusation_pom := []gossip.Gossip_object{}
	storageList_sth_full := []gossip.Gossip_object{}
	storageList_rev_full := []gossip.Gossip_object{}
	for _, gossipObject := range *c.Storage_CONFLICT_POM{
		storageList_conflict_pom = append(storageList_conflict_pom, gossipObject)
	}
	for _, gossipObject := range *c.Storage_ACCUSATION_POM {
		storageList_accusation_pom = append( storageList_accusation_pom, gossipObject)
	}
	for _, gossipObject := range *c.Storage_STH_FULL {
		storageList_sth_full = append(storageList_sth_full, gossipObject)
	}
	for _, gossipObject := range *c.Storage_REV_FULL {
		storageList_rev_full = append(storageList_rev_full, gossipObject)
	}

	err := util.WriteData(c.StorageDirectory+"/"+c.StorageFile_CONFLICT_POM, storageList_conflict_pom)
	if err!=nil{
		return err
	}
	err = util.WriteData(c.StorageDirectory+"/"+c.StorageFile_ACCUSATION_POM, storageList_accusation_pom)
	if err!=nil{
		return err
	}
	err = util.WriteData(c.StorageDirectory+"/"+c.StorageFile_STH_FULL, storageList_sth_full)
	if err!=nil{
		return err
	}
	err = util.WriteData(c.StorageDirectory+"/"+c.StorageFile_REV_FULL, storageList_rev_full)
	if err!=nil{
		return err
	}
	fmt.Println(util.BLUE,"File Storage Complete for Period: ",gossip.GetCurrentPeriod(),util.RESET)
	return nil
}

func (c *MonitorContext) LoadOneStorage(name string) error {
	storageList := []gossip.Gossip_object{}
	switch name{
	case gossip.CON_FULL:
		bytes, err := util.ReadByte(c.StorageFile_CONFLICT_POM)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &storageList)
		if err != nil {
			return err
		}
		for _, gossipObject := range storageList {
			(*c.Storage_CONFLICT_POM)[gossipObject.GetID()] = gossipObject
		}
	case gossip.ACC_FULL:
		bytes, err := util.ReadByte(c.StorageFile_ACCUSATION_POM);
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &storageList)
		if err != nil {
			return err
		}
		for _, gossipObject := range storageList {
			(*c.Storage_ACCUSATION_POM)[gossipObject.GetID()] = gossipObject
		}
	case gossip.STH_FULL:
		bytes, err := util.ReadByte(c.StorageFile_STH_FULL);
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &storageList)
		if err != nil {
			return err
		}
		for _, gossipObject := range storageList {
			(*c.Storage_STH_FULL)[gossipObject.GetID()] = gossipObject
		}
	case gossip.REV_FULL:
		bytes, err := util.ReadByte(c.StorageFile_REV_FULL);
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &storageList)
		if err != nil {
			return err
		}
		for _, gossipObject := range storageList {
			(*c.Storage_REV_FULL)[gossipObject.GetID()] = gossipObject
		}
	} 
	return errors.New("Mismatch")
}

func (c *MonitorContext) LoadStorage() error{
	err := c.LoadOneStorage(gossip.CON_FULL)
	if err != nil {
		return err
	}
	err = c.LoadOneStorage(gossip.STH_FULL)
	if err != nil {
		return err
	}
	err = c.LoadOneStorage(gossip.REV_FULL)
	if err != nil {
		return err
	}
	return nil
}



func (c *MonitorContext) GetObject(id gossip.Gossip_ID) gossip.Gossip_object{
	GType := id.Type
	switch GType{
	case gossip.CON_FULL: 
		obj := (*c.Storage_CONFLICT_POM)[id]
		return obj
	case gossip.ACC_FULL:
		obj := (*c.Storage_ACCUSATION_POM)[id]
		return obj
	case gossip.STH_FULL:
		obj := (*c.Storage_STH_FULL)[id]
		return obj
	case gossip.REV_FULL:
		obj := (*c.Storage_REV_FULL)[id]
		return obj
	case gossip.STH:
		obj := (*c.Storage_TEMP)[id]
		return obj
	case gossip.REV:
		obj := (*c.Storage_TEMP)[id]
		return obj
	}
	return gossip.Gossip_object{}

}
func (c *MonitorContext) IsDuplicate(g gossip.Gossip_object) bool {
	//no public period time for monitor :/
	id := g.GetID()
	obj := c.GetObject(id)
	return reflect.DeepEqual(obj,g)
}

func (c *MonitorContext) StoreObject(o gossip.Gossip_object) {
	switch o.Type{
		case gossip.CON_FULL: 
			(*c.Storage_CONFLICT_POM)[o.GetID()] = o
			fmt.Println(util.BLUE,"CONFLICT_POM Stored",util.RESET)
		case gossip.ACC_FULL:
			//ACCUSATION POM does not need to be stored, but this function is here for testing purposes
			(*c.Storage_ACCUSATION_POM)[o.GetID()] = o
			fmt.Println(util.BLUE,"ACCUSATION_POM Stored",util.RESET)
		case gossip.STH_FULL:
			(*c.Storage_STH_FULL)[o.GetID()] = o
			fmt.Println(util.BLUE,"STH_FULL Stored",util.RESET)
		case gossip.REV_FULL:
			(*c.Storage_REV_FULL)[o.GetID()] = o
			fmt.Println(util.BLUE,"REV_FULL Stored",util.RESET)
		default:
			(*c.Storage_TEMP)[o.GetID()] = o
		}

		

}

//wipe all temp data
func (c *MonitorContext) WipeStorage(){
	for key, _ := range *c.Storage_TEMP{
		if  key.Period!=gossip.GetCurrentPeriod(){
			delete(*c.Storage_ACCUSATION_POM,key)
		}
	}
	for key, _ := range *c.Storage_ACCUSATION_POM{
		if  key.Period!=gossip.GetCurrentPeriod(){
			delete(*c.Storage_ACCUSATION_POM,key)
		}
	}
	fmt.Println(util.BLUE,"Temp storage has been wiped.",util.RESET)
}
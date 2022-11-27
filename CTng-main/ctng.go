package main

/*
Code Ownership:
Finn - Made main function
*/
import (
	"CTng/config"
	"CTng/gossip"
	"CTng/monitor"
	"CTng/testData/fakeCA"
	"CTng/testData/fakeLogger"

	//"CTng/Logger"
	//"CTng/Client"
	"fmt"
	"os"
)

// main is run when the user types "go run ."
// it allows a user to run a gossiper, monitor, fakeLogger, or fakeCA.
// Currently unimplemented: Different object_storage locations than ./gossiper_data.json and ./monitor_data.json
// This field could be defined within the configuration files to make this more modular.
func main() {
	helpText := "Usage:\n ./CTng [gossiper|monitor] <public_config_file_path> <private_config_file_path> <crypto_config_path>\n ./Ctng [logger|ca] <fakeentity_config_path>"
	if len(os.Args) < 3 {
		fmt.Println(helpText)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "gossiper":
		// make the config object.
		conf, err := config.LoadGossiperConfig(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Println(helpText)
			panic(err)
		}
		ctx := gossip.Gossip_Context_Init(&conf, os.Args[5])
		/*
			// Space is allocated for all storage fields, and then make is run to initialize these spaces.
			storage := new(gossip.Gossip_Storage)
			*storage = make(gossip.Gossip_Storage)
			gossip_object_TSS_DB := new(gossip.Gossip_Object_TSS_DB)
			*gossip_object_TSS_DB = make(gossip.Gossip_Object_TSS_DB)
			ctx := gossip.GossiperContext{
				Config:      &conf,
				Storage:     storage,
				Obj_TSS_DB: gossip_object_TSS_DB,
				StorageFile: "gossiper_data.json", // could be a parameter in the future.
				StorageID:   os.Args[5],
			}
			ctx.Config = &conf*/
		gossip.StartGossiperServer(ctx)
		// break // break unneeded in  go.
	case "monitor":
		// make the config object.
		conf, err := config.LoadMonitorConfig(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Println(helpText)
			panic(err)
		}
		// Space is allocated for all storage fields, and then make is run to initialize these spaces.
		storage_temp := new(gossip.Gossip_Storage)
		*storage_temp = make(gossip.Gossip_Storage)
		storage_conflict_pom := new(gossip.Gossip_Storage)
		*storage_conflict_pom = make(gossip.Gossip_Storage)
		storage_accusation_pom := new(gossip.Gossip_Storage)
		*storage_accusation_pom = make(gossip.Gossip_Storage)
		storage_sth_full := new(gossip.Gossip_Storage)
		*storage_sth_full = make(gossip.Gossip_Storage)
		storage_rev_full := new(gossip.Gossip_Storage)
		*storage_rev_full = make(gossip.Gossip_Storage)
		ctx := monitor.MonitorContext{
			Config:                 &conf,
			Storage_TEMP:           storage_temp,
			Storage_CONFLICT_POM:   storage_conflict_pom,
			Storage_ACCUSATION_POM: storage_accusation_pom,
			Storage_STH_FULL:       storage_sth_full,
			Storage_REV_FULL:       storage_rev_full,
			StorageID:              os.Args[5],
		}
		ctx.Config = &conf
		monitor.StartMonitorServer(&ctx)
	case "logger":
		fakeLogger.RunFakeLogger(os.Args[2])
	case "ca":
		fakeCA.RunFakeCA(os.Args[2])

	default:
		fmt.Println(helpText)
	}
}

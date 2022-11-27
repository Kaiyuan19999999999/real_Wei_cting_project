package config

import (
	crypto "CTng/crypto"
	"fmt"
	"reflect"
	"testing"
)

func TestMonitor_confg(t *testing.T) {
	var mon_config Monitor_public_config
	var mpriv_config Monitor_config
	LoadConfiguration(&mon_config, "test/monitor_pub_config.json")
	check := reflect.ValueOf(&mon_config).Elem()
	fmt.Println("Public Config")
	for i := 0; i < check.NumField(); i++ {
		temp := check.Field(i).Interface()
		fmt.Println(temp)
	}

	fmt.Println("\nPrivate Config")
	LoadConfiguration(&mpriv_config, "test/monitor_priv_config.json")
	check = reflect.ValueOf(&mpriv_config).Elem()
	for i := 0; i < check.NumField(); i++ {
		temp := check.Field(i).Interface()
		fmt.Println(temp)
	}
}

func TestGossiper_config(t *testing.T) {
	var gos_config Gossiper_public_config
	var gpriv_config Gossiper_config

	LoadConfiguration(&gos_config, "test/gossiper_pub_config.json")
	fmt.Println("Public Config")
	check := reflect.ValueOf(&gos_config).Elem()
	for i := 0; i < check.NumField(); i++ {
		temp := check.Field(i).Interface()
		fmt.Println(temp)
	}

	fmt.Println("\nPrivate Config")
	LoadConfiguration(&gpriv_config, "test/gossiper_priv_config.json")
	check = reflect.ValueOf(&gpriv_config).Elem()
	for i := 0; i < check.NumField(); i++ {
		temp := check.Field(i).Interface()
		fmt.Println(temp)
	}
}
func TestLoadConfigs(t *testing.T) {
	entities := []crypto.CTngID{"localhost:8080", "localhost:8081", "localhost:8082", "localhost:8083"}
	configs, err := crypto.GenerateEntityCryptoConfigs(entities, 2)
	if err != nil {
		t.Error(err)
	}
	err = crypto.SaveCryptoFiles("test/", configs)
	if err != nil {
		t.Error(err)
	}
	gossiper_config, err := LoadGossiperConfig("test/gossiper_pub_config.json", "test/gossiper_priv_config.json", "test/localhost:8081.test.json")
	if err != nil {
		t.Error(err)
	}
	if (*gossiper_config.Crypto).SignScheme != "rsa" ||
		len((*gossiper_config.Public).Gossiper_URLs) == 0 {
		t.Error("Possible reading error")
	}
	fmt.Println(gossiper_config)
	fmt.Println(*gossiper_config.Public)
}

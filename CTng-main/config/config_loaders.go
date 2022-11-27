package config

import (
	crypto "CTng/crypto"
	"encoding/json"
	"io/ioutil"
	"log"
	// this is just to check that the values are being updated without having to put a bunch of print statements
)

func LoadConfiguration(config interface{}, file string) { //takes in the struct that it is updating and the file it is updating with
	// Let's first read the file
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	// Now let's unmarshall the data into `payload`
	err = json.Unmarshal(content, config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func LoadGossiperConfig(publicpath string, privatepath string, cryptopath string) (Gossiper_config, error) {
	c := new(Gossiper_config)
	c_pub := new(Gossiper_public_config)
	LoadConfiguration(c, privatepath)
	LoadConfiguration(c_pub, publicpath)
	c.Public = c_pub
	crypto, err := crypto.ReadCryptoConfig(cryptopath)
	c.Crypto = crypto
	if err != nil {
		return *c, err
	}
	return *c, nil
}

//Identical to above, but for the monitor.
func LoadMonitorConfig(publicpath string, privatepath string, cryptopath string) (Monitor_config, error) {
	c := new(Monitor_config)
	c_pub := new(Monitor_public_config)
	LoadConfiguration(c, privatepath)
	LoadConfiguration(c_pub, publicpath)
	c.Public = c_pub
	crypto, err := crypto.ReadCryptoConfig(cryptopath)
	c.Crypto = crypto
	if err != nil {
		return *c, err
	}
	return *c, nil
}

//Identical to above, but for the ca.
func LoadCAConfig(publicpath string, privatepath string, cryptopath string) (CA_config, error) {
	c := new(CA_config)
	c_pub := new(CA_public_config)
	LoadConfiguration(c, privatepath)
	LoadConfiguration(c_pub, publicpath)
	c.Public = c_pub
	crypto, err := crypto.ReadBasicCryptoConfig(cryptopath)
	c.Crypto = crypto
	if err != nil {
		return *c, err
	}
	return *c, nil
}

//Identical to above, but for the logger.
func LoadLoggerConfig(publicpath string, privatepath string, cryptopath string) (Logger_config, error) {
	c := new(Logger_config)
	c_pub := new(Logger_public_config)
	LoadConfiguration(c, privatepath)
	LoadConfiguration(c_pub, publicpath)
	c.Public = c_pub
	crypto, err := crypto.ReadBasicCryptoConfig(cryptopath)
	c.Crypto = crypto
	if err != nil {
		return *c, err
	}
	return *c, nil
}

package config

import (
	"CTng/crypto"
	"math/big"
)

// The structs that are read/written to files.
type Monitor_public_config struct {
	All_CA_URLs      []string
	All_Logger_URLs  []string
	Gossip_wait_time int
	MMD              int
	MRD              int
	Length           uint64 // max size of revocators
	Http_vers        []string
}

type Gossiper_public_config struct {
	Communiation_delay int
	Gossip_wait_time   int
	Max_push_size      int
	Period_interval    int64
	Expiration_time    int // if 0, no expiration.
	MMD                int
	MRD                int
	Gossiper_URLs      []string
	Signer_URLs        []string // List of all potential signers' DNS names.
}

type Monitor_config struct {
	Crypto_config_location string
	CA_URLs                []string
	Logger_URLs            []string
	Signer                 string
	Gossiper_URL           string
	Inbound_gossiper_port  string
	Port                   string
	Crypto                 *crypto.CryptoConfig
	Public                 *Monitor_public_config
}

type Gossiper_config struct {
	// Crypto_config_location string // Dont use: kind of confusing when considering relative paths. User will pass in absolute paths.
	Connected_Gossipers []string
	Owner_URL           string
	Port                string
	Crypto              *crypto.CryptoConfig
	Public              *Gossiper_public_config
}

// Unused, but the below info
// could be used to generate monitor+gossiper config files.
type Config_Input struct {
	Monitor_URLs []string
	/* Below is f: aka the "threshold number - 1"
		/ Each logger needs 2(f+1) connections.
		/ Each monitor needs f+1 monitor connections.
	  / Gossip_Wait_Time is threfore determined by the
		/ resulting diameter of the monitor network. */
	Max_Rogue_Parties int
	MMD               int // MRD derived from this
	CA_URLs           []string
	Logger_URLs       []string
	// Gossipers, for now, can be set to communicate on
	// the same local network as the monitor.
	Default_Gossiper_Port string
}

type CA_public_config struct {
	All_CA_URLs     []string
	All_Logger_URLs []string
	MMD             int
	MRD             int
	Http_vers       []string
	Length          uint64 // max size of revocators
	NormalizeNumber big.Int
}

type CA_config struct {
	Crypto_config_location string
	Logger_URLs            []string
	Signer                 string
	Port                   string
	Crypto                 *crypto.CryptoConfig
	Public                 *CA_public_config
}

type Logger_public_config struct {
	All_CA_URLs     []string
	All_Logger_URLs []string
	MMD             int
	MRD             int
	Http_vers       []string
	Length          uint64 // max size of revocators
}

type Logger_config struct {
	Crypto_config_location string
	CA_URLs                []string
	Signer                 string
	Port                   string
	Crypto                 *crypto.CryptoConfig
	Public                 *Logger_public_config
	// MisbehaviorInterval    int
	// LoggerType             int
	// LoggerType:
	//  1. Normal, behaving Logger (default)
	//  2. Split-World (Two different STHS on every MisbehaviorInterval MMD)
	//  3. Disconnecting Logger (unresponsive every MisbehaviorInterval MMD)
}

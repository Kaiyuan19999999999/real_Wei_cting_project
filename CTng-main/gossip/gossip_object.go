package gossip
import (
	"fmt"
	"CTng/crypto"
	"errors"
	"time"
	"strconv"
	"CTng/util"
)


// The only valid application type
const CTNG_APPLICATION = "CTng"

// Identifiers for different types of gossip that can be sent.
const (
	STH      = "http://ctng.uconn.edu/101"
	REV      = "http://ctng.uconn.edu/102"
	ACC      = "http://ctng.uconn.edu/103"
	CON      = "http://ctng.uconn.edu/104"
	STH_FRAG = "http://ctng.uconn.edu/201"
	REV_FRAG = "http://ctng.uconn.edu/202"
	ACC_FRAG = "http://ctng.uconn.edu/203"
	CON_FRAG = "http://ctng.uconn.edu/204"
	STH_FULL = "http://ctng.uconn.edu/301" 
	REV_FULL = "http://ctng.uconn.edu/302" 
	ACC_FULL = "http://ctng.uconn.edu/303"
	CON_FULL = "http://ctng.uconn.edu/304"
)

// This function prints the "name string" of each Gossip object type. It's used when printing this info to console.
func TypeString(t string) string {
	switch t {
	case STH:
		return "STH"
	case REV:
		return "REV"
	case ACC:
		return "ACC"
	case CON:
		return "CON"
	case STH_FRAG:
		return "STH_FRAG"
	case REV_FRAG:
		return "REV_FRAG"
	case ACC_FRAG:
		return "ACC_FRAG"
	case CON_FRAG:
		return "CON_FRAG"
	case STH_FULL:
		return "STH_FULL"
	case REV_FULL:
		return "REV_FULL"
	case ACC_FULL:
		return "ACC_FULL"
	case CON_FULL:
		return "CON_FULL"
	default:
		return "UNKNOWN"
	}
}

func EntityString(t string) string{
	switch t {
	case "localhost:9000":
		return "Logger 1"
	case "localhost:9001":
		return "Logger 2"
	case "localhost:9002":
		return "Logger 3"
	case "localhost:9100":
		return "CA 1"
	case "localhost:9101":
		return "CA 2"
	case "localhost:9102":
		return "CA 3"
	case "localhost:8180":
		return "Monitor 1"
	case "localhost:8181":
		return "Monitor 2"
	case "localhost:8182":
		return "Monitor 3"
	case "localhost: 8183":
		return "Monitor 4"
	case "localhost:8080":
		return "Gossiper 1"
	case "localhost:8081":
		return "Gossiper 2"
	case "localhost:8082":
		return "Gossiper 3"
	case "localhost:8083":
		return "Gossiper 4"
	default:
		return "UNKNOWN"
	}
}
type Gossip_object struct {
	Application string `json:"application"`
	Period string `json:"period"`
	Type        string `json:"type"`
	Signer string `json:"signer"`
	//**************************The number of signers should be equal to the Threshold, it just happened to be 2 in our case***************************************
	Signers map[int]string `json:"signers"`
	Signature [2]string `json:"signature"`
	// Timestamp is a UTC RFC3339 string
	Timestamp string `json:"timestamp"`
	Crypto_Scheme string `json:"crypto_scheme"`
	Payload [3]string `json:"payload,omitempty"`
}

type Gossip_ID struct{
	Period     string `json:"period"`
	Type       string `json:"type"`
	Entity_URL string `json:"entity_URL"`
}

type Gossip_Counter_ID struct{
	Period     string `json:"period"`
	Type       string `json:"type"`
	Entity_URL string `json:"entity_URL"`
	Signer     string `json:"signer"`
}


//This returns the ID of a gossip object, which is the primary key in our Gossip_Object_TSS_DB, and in our Gossip Storage
func (g Gossip_object) GetID() Gossip_ID{
	new_ID := Gossip_ID{
		Period: g.Period,
		Type: g.Type,
		Entity_URL: g.Payload[0],
	}
	return new_ID
}

func (g Gossip_object) Get_Counter_ID() Gossip_Counter_ID{
	new_ID := Gossip_Counter_ID{
		Period: g.Period,
		Type: g.Type,
		Entity_URL: g.Payload[0],
		Signer : g.Signer,
	}
	return new_ID
}

//Gossip Storage
type Gossip_Storage map[Gossip_ID]Gossip_object
type Gossip_Storage_Counter map[Gossip_Counter_ID]Gossip_object

func GetCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

//we are using 1 minitue to simulate 1 day in real life, this is subject to change 
func GetCurrentPeriod() string{
	timerfc := time.Now().UTC().Format(time.RFC3339)
	Miniutes, err := strconv.Atoi(timerfc[14:16])
	Periodnum := strconv.Itoa(Miniutes)
	if err != nil {
	}
	return Periodnum
}

func GetPrevPeriod() string{
	timerfc := time.Now().UTC().Format(time.RFC3339)
	Miniutes, err := strconv.Atoi(timerfc[14:16])
	Periodnum := strconv.Itoa(Miniutes-1)
	if err != nil {
	}
	return Periodnum
}
//this function needs to change with GetCurrentPeriod()
//this function makes sure that all monitors start querying at the beginning of each period
func Getwaitingtime() int{
	timerfc := time.Now().UTC().Format(time.RFC3339)
	Seconds, err := strconv.Atoi(timerfc[17:19])	
	if err != nil {
	}
	Seconds = 60-Seconds
	return Seconds
}

func Verify_CON(g Gossip_object, c *crypto.CryptoConfig) error{
	rsaSig1, sigerr1 := crypto.RSASigFromString(g.Signature[0])
	rsaSig2, sigerr2 := crypto.RSASigFromString(g.Signature[1])
	// Verify the signatures were made successfully
	if sigerr1 == nil && sigerr2 == nil {
		err1 := c.Verify([]byte(g.Payload[1]), rsaSig1)
		err2 := c.Verify([]byte(g.Payload[2]), rsaSig2)
		fmt.Print(util.YELLOW, err1, err2, util.RESET)
		if err1 == nil && err2 == nil {
			return nil
		} else {
			return errors.New("Message Signature Mismatch" + fmt.Sprint(err1) + fmt.Sprint(err2))
		}
	}else{
		fmt.Println(util.RED, "RSAsigConversionerror",util.RESET)
	}
	return errors.New("Message Signature Mismatch" + fmt.Sprint(sigerr1) + fmt.Sprint(sigerr2))
}
/*
func Verify_gossip_pom(g Gossip_object, c *crypto.CryptoConfig) error {
	if g.Type == CON_FRAG || g.Type == CON_FULL {
		var err1, err2 error
		if g.Signature[0] != g.Signature[1] {
			if g.Crypto_Scheme == "BLS"{
				fragsig1, sigerr1 := crypto.SigFragmentFromString(g.Signature[0])
				fragsig2, sigerr2 := crypto.SigFragmentFromString(g.Signature[1])
				// Verify the signatures were made successfully
				if sigerr1 != nil || sigerr2 != nil && !fragsig1.Sign.IsEqual(fragsig2.Sign) {
					err1 = c.FragmentVerify(g.Payload[1], fragsig1)
					err2 = c.FragmentVerify(g.Payload[2], fragsig2)
				}
			}else{
				rsaSig1, sigerr1 := crypto.RSASigFromString(g.Signature[0])
				rsaSig2, sigerr2 := crypto.RSASigFromString(g.Signature[1])
				// Verify the signatures were made successfully
				if sigerr1 != nil || sigerr2 != nil {
					err1 = c.Verify([]byte(g.Payload[1]), rsaSig1)
					err2 = c.Verify([]byte(g.Payload[2]), rsaSig2)
				}}
			}
			if err1 == nil && err2 == nil {
				return nil
			} else {
				return errors.New("Message Signature Mismatch" + fmt.Sprint(err1) + fmt.Sprint(err2))
			}
		} else {
			//if signatures are the same, there are no conflicting information
			return errors.New("This is not a valid gossip pom")
		}
}
*/
//verifies signature fragments match with payload
func Verify_PayloadFrag(g Gossip_object, c *crypto.CryptoConfig) error {
	if g.Signature[0] != "" && g.Payload[0] != "" {
		sig, _ := crypto.SigFragmentFromString(g.Signature[0])
		err := c.FragmentVerify(g.Payload[0]+g.Payload[1]+g.Payload[2], sig)
		if err != nil {
			return errors.New(No_Sig_Match)
		}
		return nil
	} else {
		return errors.New(Mislabel)
	}
}

//verifies threshold signatures match payload
func Verify_PayloadThreshold(g Gossip_object, c *crypto.CryptoConfig) error {
	if g.Signature[0] != "" && g.Payload[0] != "" {
		sig, _ := crypto.ThresholdSigFromString(g.Signature[0])
		err := c.ThresholdVerify(g.Payload[0]+g.Payload[1]+g.Payload[2], sig)
		if err != nil {
			return errors.New(No_Sig_Match)
		}
		return nil
	} else {
		return errors.New(Mislabel)
	}
}

// Verifies RSAsig matches payload, wait.... i think this just works out of the box with what we have
func Verify_RSAPayload(g Gossip_object, c *crypto.CryptoConfig) error {
	if g.Signature[0] != "" && g.Payload[0] != "" {
		sig, err := crypto.RSASigFromString(g.Signature[0])
		if err != nil {
			return errors.New(No_Sig_Match)
		}
		return c.Verify([]byte(g.Payload[0]+g.Payload[1]+g.Payload[2]), sig)

	} else {
		return errors.New(Mislabel)
	}
}

//Verifies Gossip object based on the type:
//STH and Revocations use RSA
//Trusted information Fragments use BLS SigFragments
//PoMs use Threshold signatures
func (g Gossip_object) Verify(c *crypto.CryptoConfig) error {
	// If everything Verified correctly, we return nil
	switch g.Type {
	case STH:
		return Verify_RSAPayload(g, c)
	case REV:
		return Verify_RSAPayload(g, c)
	case ACC:
		return Verify_RSAPayload(g, c)
	case CON:
		return Verify_CON(g,c)
	case STH_FRAG:
		return Verify_PayloadFrag(g, c)
	case REV_FRAG:
		return Verify_PayloadFrag(g, c)
	case ACC_FRAG:
		return Verify_PayloadFrag(g, c)
	case CON_FRAG:
		return Verify_PayloadFrag(g, c)
	case STH_FULL:
		return Verify_PayloadThreshold(g, c)
	case REV_FULL:
		return Verify_PayloadThreshold(g, c)
	case ACC_FULL:
		return Verify_PayloadThreshold(g, c)
	case CON_FULL:
		return Verify_PayloadThreshold(g, c)
	default:
		return errors.New(Invalid_Type)
	}
}

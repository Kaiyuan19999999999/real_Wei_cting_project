package monitor

import (
	"CTng/gossip"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type ClientMock struct{}

func dummyGossipObject() gossip.Gossip_object {
	return gossip.Gossip_object{
		Application: "test", Period: "0", Type: "", Signer: "",
		Signers:       make(map[int]string),
		Signature:     [2]string{"", ""},
		Timestamp:     "",
		Crypto_Scheme: "",
		Payload:       [3]string{"", "", ""},
	}
}

func (c *ClientMock) GoodRequest(req *http.Request) (*http.Request, error) {
	mockedRes := dummyGossipObject()
	b, err := json.Marshal(mockedRes)
	if err != nil {
		log.Panic("Error reading a mockedRes from mocked client", err)
	}

	return &http.Request{Body: ioutil.NopCloser(bytes.NewBuffer(b))}, nil
}

func (c *ClientMock) BadRequest(req *http.Request) (*http.Request, error) {
	mockedResBad := "bad"
	b, err := json.Marshal(mockedResBad)
	if err != nil {
		log.Panic("Error reading a mockedRes from mocked client", err)
	}

	return &http.Request{Body: ioutil.NopCloser(bytes.NewBuffer(b))}, nil
}

func TestReceiveGossip(t *testing.T) {
	monitorContext := MonitorContext{}
	req, _ := (&ClientMock{}).GoodRequest(&http.Request{})
	receiveGossip(&monitorContext, nil, req)
}

func TestPanicOnBadReceiveGossip(t *testing.T) {
	monitorContext := MonitorContext{}
	// Catch Panic
	defer func() { _ = recover() }()

	req, _ := (&ClientMock{}).BadRequest(&http.Request{})
	receiveGossip(&monitorContext, nil, req)

	t.Errorf("Expected panic")
}

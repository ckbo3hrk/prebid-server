package info

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/openrtb_ext"
)

// NewBiddersEndpoint implements /info/bidders
func NewBiddersEndpoint() httprouter.Handle {
	bidderNames := make([]string, 0, len(openrtb_ext.BidderMap))
	for bidderName := range openrtb_ext.BidderMap {
		bidderNames = append(bidderNames, bidderName)
	}

	biddersJson, err := json.Marshal(bidderNames)
	if err != nil {
		glog.Fatalf("error creating /info/bidders endpoint response: %v", err)
	}

	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(biddersJson); err != nil {
			glog.Errorf("error writing response to /info/bidders: %v", err)
		}
	})
}

// NewBiddersEndpoint implements /info/bidders/*
func NewBidderDetailsEndpoint(infos adapters.BidderInfos) httprouter.Handle {
	// Build all the responses up front, since there are a finite number and it won't use much memory.
	responses := make(map[string]json.RawMessage, len(infos))
	for bidderName, bidderInfo := range infos {
		jsonData, err := json.Marshal(bidderInfo)
		if err != nil {
			glog.Fatalf("Failed to JSON-marshal bidder-info/%s.yaml data.", bidderName)
		}
		responses[bidderName] = jsonData
	}

	// Return an endpoint which writes the responses from memory.
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		forBidder := ps.ByName("bidderName")
		if response, ok := responses[forBidder]; ok {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(response); err != nil {
				glog.Errorf("error writing response to /info/bidders/%s: %v", forBidder, err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

type infoFile struct {
	Maintainer   *maintainerInfo   `yaml:"maintainer" json:"maintainer"`
	Capabilities *capabilitiesInfo `yaml:"capabilities" json:"capabilities"`
}

type maintainerInfo struct {
	Email string `yaml:"email" json:"email"`
}

type capabilitiesInfo struct {
	App  *platformInfo `yaml:"app" json:"app"`
	Site *platformInfo `yaml:"site" json:"site"`
}

type platformInfo struct {
	MediaTypes []string `yaml:"mediaTypes" json:"mediaTypes"`
}

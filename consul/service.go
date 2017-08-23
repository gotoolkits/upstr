package consul

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	//com "github.com/gotoolkits/upstr/common"
)

const (
	API_GET_SERVICE = "/v1/catalog/service/"
)

type CatalogService struct {
	ID                       string
	Node                     string
	Address                  string
	Datacenter               string
	TaggedAddresses          map[string]string
	NodeMeta                 map[string]string
	ServiceID                string
	ServiceName              string
	ServiceAddress           string
	ServiceTags              []string
	ServicePort              int
	ServiceEnableTagOverride bool
	CreateIndex              uint64
	ModifyIndex              uint64
}

func GetSvrList(host string, sn string) ([]string, error) {

	if host == "" {
		host = os.Getenv("CONSUL_ADDR")
	}

	if !strings.Contains(host, ":") {
		host = host + ":" + "8500"
	}

	url := API_SCHEME + "://" + host + API_GET_SERVICE + sn

	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorln(err)

		//try one time again
		resp, err = http.Get(url)
		if err != nil {
			return nil, err
		}

	}
	defer resp.Body.Close()

	var out []*CatalogService
	var svrs []string

	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}

	for _, v := range out {

		if v.ServiceAddress == "" {
			svrs = append(svrs, v.Address+":"+strconv.Itoa(v.ServicePort))
		} else {
			svrs = append(svrs, v.ServiceAddress+":"+strconv.Itoa(v.ServicePort))
		}
	}

	return svrs, nil
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

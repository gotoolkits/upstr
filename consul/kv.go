package consul

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	API_GET_KV      = "/v1/kv/"
	API_SCHEME      = "http"
	DEFUALT_HOST    = "192.168.20.2:8500"
	DEFUALT_KV_PATH = "paas/ngx/upstream_name?raw"
)

func GetUpstrKV(host string) []string {

	if host == "" {
		host = DEFUALT_HOST
	}

	//http://192.168.20.2:8500/v1/kv/paas/ngx/upstream_name?raw
	url := API_SCHEME + "://" + host + API_GET_KV + DEFUALT_KV_PATH

	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	upStrName := strings.Split(string(body), ",")

	return upStrName
}

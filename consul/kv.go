package consul

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	API_GET_KV      = "/v1/kv/"
	API_SCHEME      = "http"
	DEFUALT_KV_PATH = "paas/ngx/upstream_name?raw"
)

func GetUpstrKV(host string) []string {

	if host == "" {
		host = os.Getenv("CONSUL_ADDR")
	}

	if !strings.Contains(host, ":") {
		host = host + ":" + "8500"
	}

	url := API_SCHEME + "://" + host + API_GET_KV + DEFUALT_KV_PATH

	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorln(err)
		//try again
		resp, err = http.Get(url)
		if err != nil {
			return nil
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	upStrName := strings.Split(string(body), ",")

	return upStrName
}

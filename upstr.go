package main

import (
	"net/http"

	"github.com/gotoolkits/upstr/common"
	"github.com/gotoolkits/upstr/consul"
	"github.com/labstack/echo"
)

const (
	DEFAULT_PORT = "18082"
)

type UpdateOut struct {
	Re          string `json:"status" xml:"status"`
	ErrCnt      int    `json:"errCount" xml:"errCount"`
	UpdateConut int    `json:"updateConut" xml:"updateConut"`
	UpdateTime  string `json:"UpdateTime" xml:"UpdateTime"`
}

type Info struct {
	WorkPath   string `json:"WorkPath" xml:"WorkPath"`
	ConfigPath string `json:"ConfigPath" xml:"ConfigPath"`
	Consul     string `json:"Consul" xml:"Consul"`
	KvPath     string `json:"KvPath" xml:"KvPath"`
	UpstremNum int    `json:"UpstremNum" xml:"UpstremNum"`
	Updated    int    `json:"Updated" xml:"Updated"`
	Error      int    `json:"Error" xml:"Error"`
	LastUpdate string `json:"LastUpdate" xml:"LastUpdate"`
}

var updateTime, stat string
var updated, uNum, errCount int
var m map[string][]string
var sHost, cHost, cPth, wPth string

func main() {

	jc := &common.JsonConf{}
	err := common.LoadJsonConf(jc)
	if err != nil {
		common.Log.Warningln("Load the configs failed ,Using the default configs!", err)
		sHost = DEFAULT_PORT
		cHost = consul.DEFUALT_HOST
		cPth = common.ORANGE_DEFAULT_CONF
		wPth = common.ORANGE_DEFAULT_PATH
	} else {
		sHost = jc.SrvPort
		cHost = jc.Host
		cPth = jc.ConfigPath
		wPth = jc.WorkPath
		common.Log.Infoln("Load configs from file:", sHost, cHost, cPth, wPth)
	}

	e := echo.New()
	e.GET("/", info)
	e.GET("list", list)
	e.GET("reload", reload)
	e.GET("status", status)

	e.HideBanner = true
	common.Log.Infoln("⇨ http server starting on ", ":"+sHost)
	e.Logger.Fatal(e.Start(":" + sHost))

}

// Handler get system info
func info(c echo.Context) error {

	u := &Info{
		WorkPath:   common.ORANGE_DEFAULT_PATH,
		ConfigPath: common.ORANGE_DEFAULT_CONF,
		Consul:     consul.DEFUALT_HOST,
		KvPath:     consul.DEFUALT_KV_PATH,
		UpstremNum: uNum,
		Updated:    updated,
		Error:      errCount,
		LastUpdate: updateTime,
	}

	return c.JSONPretty(http.StatusOK, u, "  ")
}

// Handler  Sync with consul,reload the config
func reload(c echo.Context) error {

	var updateCount int
	//获取consul upstream kv
	kv := consul.GetUpstrKV(cHost)
	if len(kv) < 1 {
		common.Log.Errorln("Get Consul KV Length is null!")
		errCount++
	}
	//获取已存在的配置upstream列表
	cfMap := common.GetList(cPth)
	if len(cfMap) < 1 {
		common.Log.Errorln("Get Nginx Upstream config Length is null!")
		errCount++
	}
	//判断kv值对于存在的配置列表是否有更新，如果有更新配置将更新upstream配置
	for _, k := range kv {
		if ok := common.UpstrExists(cfMap, k); !ok {
			common.Log.Infoln("Find new upstream is :", k)
			err := common.SetUpstream(k, cPth)
			if err != nil {
				common.Log.Errorln("Update the", k, "upstream failed", err)
				errCount++
			}
			common.Log.Infoln("Update the", k, "upstream successful")
			updateCount++
		}
	}
	//统计
	updated = updated + updateCount
	updateTime = common.GetTime()

	//重新reload服务,加重更新的
	if updateCount >= 1 {
		err := common.ReloadConf(wPth)
		if err != nil {
			common.Log.Errorln("Orange reload failed", err)
			errCount++
			stat = "failed"
		} else {
			common.Log.Infoln("Updated config successful!")
			stat = "successful"
		}
	} else {
		common.Log.Infoln("No need to update, configs is newest")
		stat = "nothing to do"

	}
	o := &UpdateOut{
		Re:          stat,
		ErrCnt:      errCount,
		UpdateConut: updateCount,
		UpdateTime:  updateTime,
	}
	return c.JSONPretty(http.StatusOK, o, "  ")
}

// Handler list the upstream list in configs
func list(c echo.Context) error {
	m := common.GetList(cPth)
	uNum = len(m)
	return c.JSONPretty(http.StatusOK, m, "  ")
}

// Handler upstr self status
func status(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
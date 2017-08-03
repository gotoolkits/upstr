package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gotoolkits/upstr/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// [upstream name] {service01,service02,serviec03}
const (
	NGX_SVC             = "orange"
	ORANGE_DEFAULT_PATH = "/usr/local/bin"
	ORANGE_DEFAULT_CONF = "/usr/local/orange/conf/nginx.conf"
)

var Log = logrus.New()

type Upconf struct {
	m      map[string][]string
	srvNum int
}

type JsonConf struct {
	SrvPort    string
	Host       string
	WorkPath   string
	ConfigPath string
}

type Confs struct {
	list []Upconf
}

func GetTime() string {
	t := time.Now()
	ft := t.Format("2006-01-02 15:04:05")
	return ft
}

//judge the upstream is exsited
func UpstrExists(src map[string][]string, key string) bool {
	if _, ok := src[key]; ok {
		return true
	}
	return false
}

func listenAndSrv() error {
	return nil

}

// to reload Orange config , if use the nginx ,need to change it
func ReloadConf(path string) error {

	if path == "" {
		path = ORANGE_DEFAULT_PATH
	}

	if isEx, err := PathExists(path); !isEx {
		return err
	}

	comm := path + "/" + "orange"
	cmd := exec.Command(comm, "reload")
	d, err := Execute(cmd)
	if err != nil {
		return err
	}

	if strings.Contains(string(d), "ERR") {
		return err
	}
	return nil
}

//check path exist is or not
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// execute a command
func Execute(c *exec.Cmd) ([]byte, error) {
	var err error
	var d []byte

	stdout, _ := c.StdoutPipe()
	err = c.Start()
	if err != nil {
		return nil, err
	}

	d, err = ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	err = c.Wait()
	if err != nil {
		return nil, err
	}

	return d, nil

}

func SetUpstream(un, path, host string) error {

	var line, svrsFomat, cmdArgs string
	var err error

	if path == "" {
		path = ORANGE_DEFAULT_CONF
	}

	svrStr, err := consul.GetSvrList(host, un)

	if err != nil {
		Log.Errorln("Get Consul Service list failed !")
		return err
	}

	if len(svrStr) < 1 {
		cmdArgs = fmt.Sprintf("/keepalive_timeout/a\\ \\n    upstream %s {\\n\\tzone %s 2m; \\n\\tserver 127.0.0.1:8001 down; \\n    }", un, un)

	} else {
		for _, v := range svrStr {
			line = "server " + v + " max_fails=3 fail_timeout=30s;"
			svrsFomat = svrsFomat + line + "\\n\\t"
		}
		cmdArgs = fmt.Sprintf("/keepalive_timeout/a\\ \\n    upstream %s {\\n\\tzone %s 2m; \\n\\t%s \\n    }", un, un, svrsFomat)
	}

	//fmt.Println("sed", "-i", cmdArgs, path)
	cmd := exec.Command("sed", "-i", cmdArgs, path)

	_, err = Execute(cmd)
	if err != nil {
		return err
	}

	return nil
}

func GetUpstream() {
	nList := FindAllNameRx("")
	fmt.Println(nList)
}

func GetList(p string) map[string][]string {

	var upName, svrAddr string
	var upstr map[string][]string

	upstr = make(map[string][]string)
	ss := FindAllRegex(p)

	for _, v := range ss {
		lines := SplitLines(v)

		for _, li := range lines {

			if strings.Contains(li, "upstream") {

				u := FindUpstrNameRx(li)
				if u != nil {
					upName = u[1]
				}
				//fmt.Println(upName)
				continue
			}

			if strings.Contains(li, "server") {
				s := FindSvrsRx(li)
				if s != nil {
					svrAddr = s[1]
				}
				upstr[upName] = append(upstr[upName], svrAddr)
				continue
			}
		}

	}
	return upstr
}

func LoadJsonConf(cnf *JsonConf) error {

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/upstr/")
	err := viper.ReadInConfig()
	if err != nil {
		Log.Warningln("Load Json Config failed, Please Check config file in . or /etc/upstr/")
		return err
	} else {
		cnf.SrvPort = viper.GetString("setting.port")
		cnf.Host = viper.GetString("consul.host")
		cnf.WorkPath = viper.GetString("orange.work_path")
		cnf.ConfigPath = viper.GetString("orange.config_path")
	}
	return nil
}

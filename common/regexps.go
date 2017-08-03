package common

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func FindAllRegex(filePath string) []string {

	var err error

	if filePath == "" {
		filePath = ORANGE_DEFAULT_CONF
	}

	context, err := ioutil.ReadFile(filePath)
	if err != nil {
		Log.Errorln(err)
	}

	re, err := regexp.Compile(`(?U)\s+upstream \S+\s+\{\s+[\S\s]+\s+\}`)

	if err != nil {

		Log.Errorln(err)
	}

	out := re.FindAllString(string(context), -1)
	return out
}

func FindAllNameRx(filePath string) []string {

	var err error
	var upstrName []string

	if filePath == "" {
		filePath = ORANGE_DEFAULT_CONF
	}

	context, err := ioutil.ReadFile(filePath)
	if err != nil {
		Log.Errorln(err)
	}

	re, err := regexp.Compile(`(?U)\s+upstream\s+[\S\-]+\s\{`)

	if err != nil {
		Log.Errorln(err)
		return nil
	}

	f := re.FindAllString(string(context), -1)

	fmt.Println(f)

	for _, v := range f {

		s := FindUpstrNameRx(v)
		upstrName = append(upstrName, s[1])
	}
	return upstrName
}

func FindUpstrNameRx(str string) []string {
	if len(str) < 1 {
		return nil
	}

	re, err := regexp.Compile(`(?U)upstream\s+([\S\-]+)\s\{`)
	if err != nil {
		Log.Errorln(err)
		return nil
	}

	out := re.FindStringSubmatch(str)
	return out

}

func FindSvrsRx(str string) []string {
	if len(str) < 1 {
		return nil
	}

	re, err := regexp.Compile(`(?U)server\s+(\S+)[\s;]`)
	if err != nil {
		Log.Errorln(err)
		return nil
	}

	out := re.FindStringSubmatch(str)
	return out
}

func SplitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

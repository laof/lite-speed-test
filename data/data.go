package data

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/laof/lite-speed-test/ping"
)

type Nodes struct {
	Url string
	Max int
}

func (s Nodes) Get() (HttpData, error) {
	var data HttpData
	res, err := http.Get(s.Url)

	if err != nil {
		return data, err
	}

	defer res.Body.Close()
	str, _ := io.ReadAll(res.Body)

	err = json.Unmarshal(str, &data)

	if err != nil {
		return data, err
	}

	for i, item := range data.List {
		name := reverseString(item.Name)
		datetime := reverseString(item.Datetime)
		txt := item.Data

		for _, o := range data.Decode {
			txt = strings.ReplaceAll(txt, o.K, o.V)
		}

		data.List[i].Name = name
		data.List[i].Datetime = datetime
		data.List[i].Data = txt
	}
	data.Decode = make([]Decode, 0)
	return data, nil
}

func (s Nodes) Test(data HttpData) (TestResult, error) {
	var servers []string
	var nodes []string
	var max = s.Max
	var all []string
	result := TestResult{}
	for _, item := range data.List {

		if item.Data == "" {
			continue
		}

		arr := strings.Split(item.Data, ",")
		all = append(all, item.Name)
		if max > 0 {
			if len(arr) > max {
				arr = arr[0:max]
			}
		}
		nodes = append(nodes, arr...)
		for range arr {
			servers = append(servers, item.Name)
		}
	}

	if len(nodes) == 0 {
		return result, errors.New("no data(nodes)")
	}

	ssr := strings.Join(nodes, "\n")
	res, err := ping.Test(ssr)
	if err != nil {
		return result, err
	}

	for _, n := range res.SuccessIndex {
		result.SuccessNodes = append(result.SuccessNodes, nodes[n])
	}

	result.ErrorServers = getServerByIndex(res.ErrorIndex, servers)

	for _, vvvv := range all {

		if !hasValue(result.ErrorServers, vvvv) {
			result.SuccessServers = append(result.SuccessServers, vvvv)
		}

	}

	return result, nil
}

func hasValue(str []string, value string) bool {

	for _, v := range str {
		if v == value {
			return true
		}
	}
	return false
}

func removeDuplicates(arr []string) []string {

	visited := make(map[string]bool)
	var list []string
	for _, str := range arr {
		if visited[str] {
			continue
		}
		visited[str] = true
		list = append(list, str)
	}
	return list

}

func getServerByIndex(indexArr []int, service []string) []string {
	var nodes []string

	for _, n := range indexArr {

		if !hasValue(nodes, service[n]) {
			nodes = append(nodes, service[n])
		}
	}

	return nodes

}

func reverseString(str string) string {
	// 将字符串转换为字节切片
	byteSlice := []byte(str)
	length := len(byteSlice)

	// 使用双指针进行字节切片的反转
	for i := 0; i < length/2; i++ {
		byteSlice[i], byteSlice[length-i-1] = byteSlice[length-i-1], byteSlice[i]
	}

	// 将字节切片转换为字符串并返回
	return string(byteSlice)
}

type TestResult struct {
	SuccessNodes   []string
	ErrorServers   []string
	SuccessServers []string
}

type HttpData struct {
	List   []List   `json:"list"`
	Decode []Decode `json:"decode"`
	Update string   `json:"update"`
	Conf   []string `json:"conf"`
}

type List struct {
	Name     string `json:"name"`
	Datetime string `json:"datetime"`
	Length   int    `json:"length"`
	Data     string `json:"data"`
}
type Decode struct {
	K string `json:"k"`
	V string `json:"v"`
}

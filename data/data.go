package data

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/laof/lite-speed-test/ping"
)

func Get(link string, check bool, max int) (Res, error) {
	var data Res
	res, err := http.Get(link)

	if err != nil {
		return data, err
	}

	defer res.Body.Close()
	str, _ := io.ReadAll(res.Body)

	err = json.Unmarshal(str, &data)

	if err != nil {
		return data, err
	}

	var servers []string
	var nodes []string
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

		arr := strings.Split(txt, ",")

		if max > 0 {
			if len(arr) > max {
				arr = arr[0:max]
			}
		}

		nodes = append(nodes, arr...)
		for range arr {
			servers = append(servers, name)
		}

	}

	data.Decode = make([]Decode, 0)

	if !check {
		return data, nil
	}

	ssr := strings.Join(nodes, "\n")

	testRes, er := ping.Test(ssr)
	if er != nil {
		return data, er
	}

	for _, n := range testRes.SuccessIndex {
		data.SuccessNodes = append(data.SuccessNodes, nodes[n])
	}

	data.ErrorServers = getErrorServers(testRes.ErrorIndex, servers)

	return data, nil
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

func getErrorServers(errorIndex []int, service []string) []string {
	var nodes []string

	for _, n := range errorIndex {

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

type Res struct {
	SuccessNodes []string `json:"successNodes"`
	ErrorServers []string `json:"errorServers"`
	List         []List   `json:"list"`
	Decode       []Decode `json:"decode"`
	Update       string   `json:"update"`
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

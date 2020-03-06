package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var ambariEndpoint string
var ambariUser string
var ambariPassword string

func main() {
	flag.StringVar(&ambariEndpoint, "ambari", "", "ambari endpoint")
	flag.StringVar(&ambariUser, "user", "", "ambari user")
	flag.StringVar(&ambariPassword, "password", "", "ambari password")
	flag.Parse()
	hosts, err := getAllHosts()
	if err != nil {
		panic(err.Error())
	}
	for hostName := range hosts {
		startAllComponentsOnTheHost(hostName)
	}
}

func getAllHosts() (map[string]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ambariEndpoint+"/hosts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Requested-By", "ambari")
	req.SetBasicAuth(ambariUser, ambariPassword)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res := map[string]interface{}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	hostItems, ok := res["items"]
	if !ok {
		return nil, errors.New("Error when getting all hosts: no items")
	}
	hosts := make(map[string]string)
	for _, item := range hostItems.([]interface{}) {
		Hosts, ok := item.(map[string]interface{})["Hosts"]
		if !ok {
			return nil, errors.New("Error when getting all hosts: no Hosts in the item")
		}
		hostName, ok := Hosts.(map[string]interface{})["host_name"]
		if !ok {
			return nil, errors.New("Error when getting all hosts: no host_name in the Hosts")
		}
		href, ok := item.(map[string]interface{})["href"]
		if !ok {
			return nil, errors.New("Error when getting all hosts: no href in the item")
		}
		hosts[hostName.(string)] = href.(string)
	}
	return hosts, nil
}

func startAllComponentsOnTheHost(hostName string) error {
	requestBody := map[string]interface{}{
		"RequestInfo": map[string]string{
			"context": "Auto Start All Components",
		},
		"Body": map[string]map[string]string{
			"HostRoles": {"state": "STARTED"},
		},
	}
	requestBodyStr, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/hosts/%s/host_components", ambariEndpoint, hostName), strings.NewReader(string(requestBodyStr)))
	if err != nil {
		return err
	}
	req.Header.Set("X-Requested-By", "ambari")
	req.SetBasicAuth(ambariUser, ambariPassword)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

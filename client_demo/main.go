package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	SERVER_PREFIX = "http://localhost:8080"
)

func main() {
	client := new(http.Client)

	// send 1 simple request
	sendTo := fmt.Sprintf("%s/run", SERVER_PREFIX)
	cl := clientRequest{
		String1: "fizz",
		String2: "buzz",
		Int1:    3,
		Int2:    5,
		Limit:   16,
	}
	err := SendAndCheck(client, sendTo, `{"result":"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16"}`, 200, cl)
	if err != nil {
		log.Printf("Failed simple test: %+v", err)
		return
	}
	log.Printf("Simple test ok")
	time.Sleep(time.Second)

	// check stat
	sendTo = fmt.Sprintf("%s/stat?key=SUNXIN", SERVER_PREFIX)
	resp, err := client.Get(sendTo)
	if err != nil {
		log.Printf("Failed send when check stat: %+v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("Bad response code when check stat: get %d", resp.StatusCode)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed ReadAll body when check stat: %+v", err)
		return
	}
	resp.Body.Close()

	if !strings.Contains(string(b), `with params &{String1:fizz String2:buzz Int1:3 Int2:5 Limit:16 ResponseGzip:false}, stat ok`) {
		log.Printf("Bad response when check stat: get (%s)", b)
		return
	}

	log.Printf("Check stat ok")
	time.Sleep(time.Second)

	// update conf
	sendTo = fmt.Sprintf("%s/update_conf?key=SUNXIN", SERVER_PREFIX)
	conf := configuration{
		Password:     "SUN_XIN",
		ForceGzip:    true,
		ForceGzipNum: 10,
	}
	err = SendAndCheck(client, sendTo, "", 200, conf)
	if err != nil {
		log.Printf("Failed update conf test: %+v", err)
		return
	}
	log.Printf("Update conf ok")
	time.Sleep(time.Second)

	// check stat with old password
	sendTo = fmt.Sprintf("%s/stat?key=SUNXIN", SERVER_PREFIX)
	resp, err = client.Get(sendTo)
	if err != nil {
		log.Printf("Failed send when check stat: %+v", err)
		return
	}

	if resp.StatusCode != 403 {
		log.Printf("Passwod changed but not work: get %d", resp.StatusCode)
		return
	}
	log.Printf("old password check ok")
	time.Sleep(time.Second)

	// check stat with new password
	sendTo = fmt.Sprintf("%s/stat?key=SUN_XIN", SERVER_PREFIX)
	resp, err = client.Get(sendTo)
	if err != nil {
		log.Printf("Failed send when check stat: %+v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("Check new password, but not work")
		return
	}
	log.Printf("new password check ok")
	time.Sleep(time.Second)

	// response must use gzip
	sendTo = fmt.Sprintf("%s/run", SERVER_PREFIX)
	cl = clientRequest{
		String1: "fizz",
		String2: "buzz",
		Int1:    2,
		Int2:    3,
		Limit:   10,
	}
	err = SendAndCheck(client, sendTo, `{"result":"1,fizz,buzz,fizz,5,fizzbuzz,7,fizz,buzz,fizz"}`, 200, cl)
	if err != nil {
		log.Printf("Failed gzip test: %+v", err)
		return
	}
	log.Printf("gzip test ok")
}

func SendAndCheck(client *http.Client, sendTo, expectRes string, expectCode int, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("Failed Marshal: %+v", err)
	}

	q, err := http.NewRequest("POST", sendTo, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("Failed NewRequest: %+v", err)
	}
	q.Header.Add("Accept-Encoding", "gzip")

	resp, err := client.Do(q)
	if err != nil {
		return fmt.Errorf("Failed Send: %+v", err)
	}

	if resp.StatusCode != expectCode {
		return fmt.Errorf("Send get code no %d: %d", resp.StatusCode, expectCode)
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		log.Printf("response gzip")
		rd, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("Failed NewReader gzip: %+v", err)
		}
		defer rd.Close()

		b, err = ioutil.ReadAll(rd)
		if err != nil {
			return fmt.Errorf("Failed ReadAll body: %+v", err)
		}
	} else { // no gzip
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Failed ReadAll body: %+v", err)
		}
	}
	defer resp.Body.Close()

	if expectRes == "" {
		return nil
	}

	if string(b) != expectRes {
		return fmt.Errorf("Bad response, expect %s, get %s", expectRes, b)
	}

	return nil
}

type clientRequest struct {
	// required
	String1 string `json:"string1"`
	String2 string `json:"string2"`
	Int1    int    `json:"int1"`
	Int2    int    `json:"int2"`
	Limit   int    `json:"limit"`

	// optional
	ResponseGzip bool `json:"response_gzip"`
}

type configuration struct {
	// TODO: must be more securely, ex: public_key, private_key
	Password string `json:"password"`

	ForceGzip    bool  `json:"force_gzip"`
	ForceGzipNum int64 `json:"force_gzip_num"`
}

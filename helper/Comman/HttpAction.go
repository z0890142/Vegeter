package Comman

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}, Log *logrus.Entry) error {

	response, err := json.Marshal(payload)

	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)

	return err
}

func GetTime() string {
	t := time.Now()
	NowTime := t.In(time.FixedZone("CST", 28800)).Format("2006-01-02T15:04:05")
	return NowTime
}
func GetTimev2() string {
	t := time.Now()
	NowTime := t.In(time.FixedZone("CST", 28800)).Format("2006-01-02T15:04:05.0000")
	return NowTime
}
func StartHTTPPost_Control(sensorID string, deviceID string, APIKEY string, rawdata []string, account string, Log *logrus.Entry) string {

	client := &http.Client{}

	url := "http://iot.cht.com.tw/iot/v1/device/" + deviceID + "/rawdata"
	NowTime := GetTime()
	var PostData string

	PostData = "[{\"id\":\"" + sensorID + "\",\"time\":"

	PostData = PostData + "\"" + NowTime + "\""

	PostData = PostData + ",\"value\": [\"" + rawdata[0] + "\",\"" + rawdata[1] + "\",\"" + rawdata[2] + "\",\"" + rawdata[3] + "-" + account + "\"]}]"
	Log.WithFields(logrus.Fields{
		"deviceID": deviceID,
		"PostData": PostData,
	}).Info()
	var jsonStr = []byte(PostData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("http request error")
		return "Error"

	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CK", APIKEY)

	resp, err := client.Do(req)

	if err != nil {

		Log.Error("Error opening error log : %v", err)

		return "Error"

	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		Log.Error("Error opening error log : %v", err)

		return "Error"

	}

	return "Post OK\n" + string(body)

}
func DMP_Control(deviceID string, APIKEY string, payload []byte, Log *logrus.Entry) bool {

	client := &http.Client{}

	url := "http://iot.cht.com.tw/iot/v1/device/" + deviceID + "/rawdata"
	Log.WithFields(logrus.Fields{
		"deviceID": deviceID,
		"PostData": string(payload),
	}).Info()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"deviceID": deviceID,
			"error":    err,
		}).Error("DMP Control Error")

		return false

	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CK", APIKEY)

	resp, err := client.Do(req)

	if err != nil {

		Log.WithFields(logrus.Fields{
			"deviceID": deviceID,
			"error":    err,
		}).Error("DMP Control Error")
		return false

	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)

	if err != nil {

		Log.WithFields(logrus.Fields{
			"deviceID": deviceID,
			"error":    err,
		}).Error("DMP Control Error")
		return false

	}
	Log.WithFields(logrus.Fields{
		"deviceID": deviceID,
	}).Info("DMP Control")
	return true

}
func HttpUrlEncoderRequest(action string, url string, header map[string]string, body url.Values, Log *logrus.Entry) (resp *http.Response, err error) {
	// requestBodyByte, _ := json.Marshal(body)

	Log.WithFields(logrus.Fields{
		"action": action,
		"url":    url,
		"header": header,
		"body":   body,
	}).Info("HttpUrlEncoderRequest Info")

	req, err := http.NewRequest(action, url, strings.NewReader(body.Encode()))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"action": action,
			"url":    url,
			"header": header,
			"body":   body,
			"err":    err,
		}).Error("NewRequest error")
		return
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	httpClient := http.Client{}
	resp, err = httpClient.Do(req)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"action": action,
			"url":    url,
			"header": header,
			"body":   body,
			"err":    err,
		}).Error("httpClient.Do error")
		return
	}
	return
}

func HttpRequest(action string, url string, header map[string]string, body interface{}, Log *logrus.Entry) (resp *http.Response, err error) {
	requestBodyByte, _ := json.Marshal(body)

	Log.WithFields(logrus.Fields{
		"Action":       action,
		"Url":          url,
		"Header":       header,
		"Request Body": body,
	}).Info("HttpRequest")
	var req *http.Request
	if action == "GET" {
		req, err = http.NewRequest(action, url, nil)
	} else {
		req, err = http.NewRequest(action, url, bytes.NewBuffer([]byte(requestBodyByte)))
	}
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Action":       action,
			"Url":          url,
			"Header":       header,
			"Request Body": body,
			"Error":        err,
		}).Error("Http Request Error")
		return
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	httpClient := http.Client{
		Timeout: 15 * time.Second,
		// Transport: &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// },
	}
	resp, err = httpClient.Do(req)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Action":       action,
			"Url":          url,
			"Header":       header,
			"Request Body": body,
			"Error":        err,
		}).Error("Http Request Error")
		return
	}
	return
}

func JsonUnMarshal(requestBody io.Reader, unMarshalItem interface{}, Log *logrus.Entry) error {

	body, err := ioutil.ReadAll(io.LimitReader(requestBody, 102400)) //io.LimitReader限制大小
	if err != nil {
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("ioutil error")
		return errors.New("ioutil error")
	}

	err = json.Unmarshal(body, unMarshalItem)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"request body": string(body),
			"error":        err,
		}).Error("Unmarshal error")
		return errors.New("Unmarshal error")
	}
	return nil
}

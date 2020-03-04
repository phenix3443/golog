package log

// go test -check.f 'BookCacheSuit.TestBookReadProcess'

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// func TestInit(t *testing.T) {
// 	Init()
// 	defer DLogger.Sync()
// 	defer SLogger.Sync()

// 	DLogger.Info("diagosis info log")

// 	SLogger.Info("static info log")

// }

func TestChangeDloggerLevel(t *testing.T) {
	Init()
	defer DLogger.Sync()

	DLogger.Info("this debug log should show")

	body := map[string]string{"level": "info"}
	jsonBody, _ := json.Marshal(body)
	log.Println(string(jsonBody))
	url := fmt.Sprintf("http://localhost:%d/dlogger", DLoggerPort)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	if resp, err := http.DefaultClient.Do(req); err != nil {
		log.Printf("修改日志等等级出错,err=%v", err)
	} else {
		defer resp.Body.Close()
		code := resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("statusCode=%d,body=%s", code, body)
	}

	DLogger.Debug("this debug log should not show")
	DLogger.Info("this info log should show")

}

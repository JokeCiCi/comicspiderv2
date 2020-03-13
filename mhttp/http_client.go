package mhttp

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/JokeCiCi/comicspiderv2/ecb"
)

func HttpGetAndStore(url, path string) bool {
	resp, err := http.Get(url)
	for err != nil {
		log.Println("http get err:", err)
		resp, err = http.Get(url)
		time.Sleep(time.Second)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false
	}
	b, _ := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(path, b, 0644)
	for err != nil {
		log.Println("ioutil write file err:", err)
		err = ioutil.WriteFile(path, b, 0644)
		time.Sleep(time.Second)
	}
	return true
}

func HttpGet(url string) (b []byte) {
	resp, err := http.Get(url)
	for err != nil {
		log.Println("http get err:", err)
		resp, err = http.Get(url)
		time.Sleep(time.Second)
	}
	defer resp.Body.Close()
	b, _ = ioutil.ReadAll(resp.Body)
	return
}

// HTTPEncryptGet ...
func HTTPEncryptGet(url string) string {
	b := HttpGet(url)
	return decrypt(b)
}

func decrypt(html []byte) (content string) {
	aescontent, _ := base64.StdEncoding.DecodeString(string(html))
	key := "281907eggplant89"
	decontent := []byte(ecb.Decrypt(aescontent, key))
	if len(decontent) == 0 {
		return
	}
	content = string(ecb.PKCS7UPad(decontent))
	return
}

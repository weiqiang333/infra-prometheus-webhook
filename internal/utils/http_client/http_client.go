package http_client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func HttpsPostForm(url string, data url.Values) error {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.PostForm(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed HttpsPostForm in StatusCode != 200 is %v, body: %s ", resp.StatusCode, string(body))
	}
	log.Println("Info HttpsPostForm Success", url, data, "body:", string(body))
	return nil
}

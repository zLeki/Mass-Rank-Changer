package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zLeki/Goblox/account"
	"github.com/zLeki/Goblox/csrf"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var client http.Client
type Users struct {
	Data               []struct {
		UserID                     int    `json:"userId"`
		Username                   string `json:"username"`
	} `json:"data"`
}
var roblo = http.Cookie{
Name:  ".ROBLOSECURIY",
Value: "cookie",
}
func DemoteEveryone(users Users) (jar *http.CookieJar){
	acc := account.Validate("")
	for _, user := range users.Data {
		endpoint := fmt.Sprintf("https://groups.roblox.com/v1/groups/13685622/users/%v", user.UserID)
		dataBytes := []byte(`{"roleId":78105964}`)
		req, err := http.NewRequest("PATCH", endpoint, bytes.NewReader(dataBytes))
		if err != nil {
			log.Fatal(err)
		}
		jar, err := cookiejar.New(nil)
		URI, err := url.Parse(endpoint)
		jar.SetCookies(URI, []*http.Cookie{&roblo})
		req.AddCookie(&roblo)
		getCSRF, _ := csrf.GetCSRF(acc)
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
		req.Header.Set("X-CSRF-TOKEN", getCSRF)
		client := &http.Client{Timeout: 10 * time.Second, Jar: jar}
		resp, err := client.Do(req)
		log.Println(getCSRF)
		if resp.StatusCode != 200 {
			date, _ := ioutil.ReadAll(resp.Body)
			log.Fatal(string(date))
		}else{
			log.Println(user.Username+" successfully demoted")
		}
	}
	return nil
}
func main() {
	req, err := http.NewRequest("GET", "https://groups.roblox.com/v1/groups/13685622/roles/78105882/users?sortOrder=Asc&limit=100&_=1646002139490", nil)

	if err != nil {
		panic(err)
	}

	req.AddCookie(&roblo)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	var user Users
	err = json.NewDecoder(resp.Body).Decode(&user)
	DemoteEveryone(user)
}

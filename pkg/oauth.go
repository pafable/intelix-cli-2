/*
This will get an oauth token from Intelix
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type OauthResp struct {
	AccessToken string `json:"access_token"`
	Exp         int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

const (
	oAuthUri = "https://api.labs.sophos.com/oauth2/token"
)

func getOauthToken(u string, cID string, cSec string) string {
	var o OauthResp

	/* Create a new request to pass in header for basic auth */
	postReq, err := http.NewRequest("POST", u, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		log.Fatal(err)
	}

	/* Convert client ID and client secret to base64
	client-id:client-secret */

	b64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cID, cSec)))
	postReq.Header.Add("Authorization", "Basic "+b64)
	postReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}

	resp, err := client.Do(postReq)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &o)
	if err != nil {
		log.Fatal(err)
	}

	return o.AccessToken
}

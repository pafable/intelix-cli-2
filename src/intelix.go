package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type OauthResp struct {
	AccessToken string `json:"access_token"`
	Exp         int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type ReportResponse struct {
	JobId     string `json:"jobId"`
	JobStatus string `json:"jobStatus"`
	Report    Report
}

type Report struct {
	Score int `json:"score"`
}

const (
	oAuthUri = "https://api.labs.sophos.com/oauth2/token"
)

func main() {
	clientId := os.Getenv("INTELIX_CLIENT_ID")
	clientSecret := os.Getenv("INTELIX_CLIENT_SECRET")
	var oauthToken string = getOauthToken(oAuthUri, clientId, clientSecret)
	// fmt.Println(oauthToken)

	// Create flag for "static"
	staticCMD := flag.NewFlagSet("static", flag.ExitOnError)
	staticFile := staticCMD.String("file", "", "static scan a file")

	// Create flag for "dynamic"
	dynamicCMD := flag.NewFlagSet("dynamic", flag.ExitOnError)
	dynamicFile := dynamicCMD.String("file", "", "dynamic scan a file")

	// Checks to make more than 1 arg is passed in for either static or dynamic analysis
	if len(os.Args) < 2 {
		fmt.Println("Enter subcommands [ static ], [ dynamic ], [ version ]")
		log.Fatal("Expected more args")
	}

	switch os.Args[1] {
	case "static":
		statusCheck, jID := fileCheck(staticCMD, staticFile, oauthToken, "static")

		for statusCheck != "SUCCESS" {
			_, statusCheck := getFileAnalysisReport(jID, oauthToken, "static")
			time.Sleep(10 * time.Second)
			if statusCheck == "SUCCESS" {
				break
			}
		}

		score, statusCheck := getFileAnalysisReport(jID, oauthToken, "static")
		fmt.Printf("Score: %d \nStatus: %s\n", score, statusCheck)
	case "dynamic":
		statusCheck, jID := fileCheck(dynamicCMD, dynamicFile, oauthToken, "dynamic")

		for statusCheck != "SUCCESS" {
			_, statusCheck = getFileAnalysisReport(jID, oauthToken, "dynamic")
			time.Sleep(10 * time.Second)
			if statusCheck == "SUCCESS" {
				break
			}
		}

		score, statusCheck := getFileAnalysisReport(jID, oauthToken, "dynamic")
		fmt.Printf("Score: %d \nStatus: %s\n", score, statusCheck)
	case "version":
		checkVersion()
	default:
	}
}

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

func fileCheck(sCMD *flag.FlagSet, file *string, token, analysisType string) (status, job_id string) {
	var r ReportResponse
	var uri string = "https://us.api.labs.sophos.com/analysis/file/" + analysisType + "/v1"
	sCMD.Parse(os.Args[2:])

	if *file == "" {
		sCMD.PrintDefaults()
	}

	if *file != "" {
		bodyBuff := &bytes.Buffer{}

		writer := multipart.NewWriter(bodyBuff)
		// fmt.Println(writer.FormDataContentType())

		fileWriter, err := writer.CreateFormFile("file", *file)
		if err != nil {
			log.Fatal(err)
		}

		fh, err := os.Open(*file)
		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			log.Fatal(err)
		}

		writer.Close()

		req, err := http.NewRequest("POST", uri, bodyBuff)
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		client := http.Client{}

		// Sends POST request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(resp.StatusCode)

		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(string(resp_body))

		err = json.Unmarshal(resp_body, &r)
		// fmt.Println(r)
	}

	return r.JobStatus, r.JobId
}

func checkVersion() {
	fmt.Println("Intelix CLI version 2.0.0")
}

func getFileAnalysisReport(id, token, analysisType string) (score int, status string) {
	var r ReportResponse
	var uri string = "https://us.api.labs.sophos.com/analysis/file/" + analysisType + "/v1/reports/" + id + "?report_format=json"
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", token)

	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(resp_body))

	err = json.Unmarshal(resp_body, &r)

	return r.Report.Score, r.JobStatus
}

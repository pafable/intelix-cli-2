/*
This will do a static or dynamic file analysis
*/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type ReportResponse struct {
	JobId     string `json:"jobId"`
	JobStatus string `json:"jobStatus"`
	Report    Report
}

type Report struct {
	Score int `json:"score"`
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

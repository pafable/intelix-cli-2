package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
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

func checkVersion() {
	fmt.Println("Intelix CLI version 2.0.0")
}

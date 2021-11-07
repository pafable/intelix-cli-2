# Intelix CLI 2

This is version 2 of Intelix CLI.
The Intelix CLI uses [Sophos Intelix](https://www.sophos.com/en-us/labs/intelix.aspx) API to analyze suspicious files in the commandline.

Please subscribe to Intelix at the [AWS Marketplace](https://aws.amazon.com/marketplace/pp/prodview-k4jb2agd65ses). (You will need an AWS account)

Once subscribed you will receive a client ID and client secret. Pass those in as environment variables. Use the environment variables below.

Add credentials as environment variables.
- INTELIX_CLIENT_ID
- INTELIX_CLIENT_SECRET

## Building the Intelix binary
Run the command below to create the binary.

MacOS & Linux
```
make build
```

Windows
```
go build -o intelix.exe .\pkg\filecheck.go .\pkg\main.go .\pkg\oauth.go
```

## Commands
Static file analysis
```
intelix static -file <filename>
```

Dynamic file analysis
```
intelix dynamic -file <filename>
```

## TODOs
- Add file hash analysis
- Add static URL analysis
- Add URL category lookup
- Add IP category lookup
- Add Android APK lookup
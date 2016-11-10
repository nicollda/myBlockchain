@ECHO OFF
cd c:\gogit

:build
cls
ECHO Building.....
go build ./

IF NOT ERRORLEVEL 1 (
	go clean ./
	git commit -m "commit" -a
	git push
	c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d @deploy.json "https://f7ff1a1d892c4bcc95e49b543caae4f4-vp0.us.blockchain.ibm.com:446/chaincode" > serviceKey.txt
)


pause

goto :build

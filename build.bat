@ECHO OFF
cd c:\gogit

:build
cls
ECHO Building.....
go build ./myHL/

IF NOT ERRORLEVEL 1 (
	go clean ./myHL/
	git commit -m "commit" -a
	git push
	c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d @deploy.json "https://56cdfd1d1c00471c913f22644988565a-vp2.us.blockchain.ibm.com:5002/chaincode" > serviceKey.txt
)


pause

goto :build

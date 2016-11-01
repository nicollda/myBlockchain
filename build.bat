@ECHO OFF
cd c:\gogit

:build
cls
ECHO Building.....
go build c:\gogit\myChaincode.go


IF NOT ERRORLEVEL 1 (
	git commit -m "commit" -a
	git push
	curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d @deploy.json "https://9f51870d42f942c79c8ab4e217f236d8-vp1.us.blockchain.ibm.com:444/chaincode" > serviceKey.txt
	set /p Keytext=<serviceKey.txt
	echo *********
	echo %Keytext%
)


pause

goto :build

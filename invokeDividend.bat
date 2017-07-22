@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 
c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"invoke\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"%Keytext:~52,128%\"},\"ctorMsg\": {\"function\": \"dividend\",\"args\": [\"JaimeKilled\",\"50\"]},\"secureContext\": \"dashboarduser_type1_2\"},\"id\": 1}" "https://54419584d43e4edc8e295d24ef28a84a-vp2.us.blockchain.ibm.com:5003/chaincode"
@echo. 

pause

goto :QUERY
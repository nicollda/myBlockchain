@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 
curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"query\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"%Keytext:~52,128%\"},\"ctorMsg\": {\"function\": \"holdings\",\"args\": [\"Aaron\"]},\"secureContext\": \"dashboarduser_type1_0\"},\"id\": 1}" "https://9f51870d42f942c79c8ab4e217f236d8-vp1.us.blockchain.ibm.com:444/chaincode"
@echo. 

pause

goto :QUERY
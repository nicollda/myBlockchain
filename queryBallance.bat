@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 
c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"query\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"%Keytext:~52,128%\"},\"ctorMsg\": {\"function\": \"ballance\",\"args\": [\"Aaron\"]},\"secureContext\": \"dashboarduser_type1_2\"},\"id\": 1}" "https://56cdfd1d1c00471c913f22644988565a-vp2.us.blockchain.ibm.com:5002/chaincode"
@echo. 

pause

goto :QUERY
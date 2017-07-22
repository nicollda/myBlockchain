@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 
c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"query\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"%Keytext:~52,128%\"},\"ctorMsg\": {\"function\": \"query\",\"args\": [\"a\"]},\"secureContext\": \"dashboarduser_type1_2\"},\"id\": 1}" "https://b3a3128c53d74b30a25d836da9e6bfed-vp2.us.blockchain.ibm.com:5002/chaincode"
@echo. 

pause

goto :QUERY
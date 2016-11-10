@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 
c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"query\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"%Keytext:~52,128%\"},\"ctorMsg\": {\"function\": \"users\",\"args\": [\"a\"]},\"secureContext\": \"dashboarduser_type1_0\"},\"id\": 1}" "https://f7ff1a1d892c4bcc95e49b543caae4f4-vp0.us.blockchain.ibm.com:446/chaincode"
@echo. 

pause

goto :QUERY
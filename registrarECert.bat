@ECHO OFF
set URL=https://f7ff1a1d892c4bcc95e49b543caae4f4-vp0.us.blockchain.ibm.com:446/registrar
set User=dashboarduser_type1_0

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 

echo registrar:
c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"enrollId\": \"dashboarduser_type1_0\",\"enrollSecret\": \"9ca37109cd\"}" "%URL%"

@echo. 
@echo. 
echo get ecert:
c:\curl\curl.exe -X GET --insecure --header "Accept: application/json" "%URL%/%User%/ecert" > ecert.json

@echo. 
echo get tcert:
c:\curl\curl.exe -X GET --insecure --header "Accept: application/json" "%URL%/%User%/tcert?count=1" > tcert.json

@echo. 

pause

goto :QUERY
@ECHO OFF
set URL=https://54419584d43e4edc8e295d24ef28a84a-vp2.us.blockchain.ibm.com:5003/registrar
set User=dashboarduser_type1_0

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 

echo registrar:
c:\curl\curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"enrollId\": \"%User%\",\"enrollSecret\": \"9ca37109cd\"}" "%URL%"

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
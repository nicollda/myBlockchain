@ECHO OFF
cd c:\gogit

:build
cls
ECHO Building.....
go build c:\gogit\myChaincode.go


IF NOT ERRORLEVEL 1 git commit -m "commit" 


pause

goto :build

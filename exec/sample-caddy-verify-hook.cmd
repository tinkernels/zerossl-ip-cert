@echo off

SETLOCAL

echo "ZEROSSL_HTTP_FV_HOST: %ZEROSSL_HTTP_FV_HOST%"
echo "ZEROSSL_HTTP_FV_PATH: %ZEROSSL_HTTP_FV_PATH%"
echo "ZEROSSL_HTTP_FV_PORT: %ZEROSSL_HTTP_FV_PORT%"
echo "ZEROSSL_HTTP_FV_CONTENT: %ZEROSSL_HTTP_FV_CONTENT%"

set VF_FILE="C:\Programs\caddy\caddy-zerossl-http-verify.conf"

echo :%ZEROSSL_HTTP_FV_PORT% {> %VF_FILE%
echo     respond %ZEROSSL_HTTP_FV_PATH% 200 {>> %VF_FILE%

REM CONTENT
set TMP_VF_CNTENT=%ZEROSSL_HTTP_FV_CONTENT%
set "IF_FIRSTLN=1"
:loop
for /f "tokens=1*" %%a in ("%TMP_VF_CNTENT%") do (
   set TMP_VF_CNTENT=%%b
   if %IF_FIRSTLN% == 1 (
    echo FIRST LINE %%a
    echo         body ^"%%a>> %VF_FILE%
    set "IF_FIRSTLN=0"
   ) else (
    echo OTHER LINE %%a
    echo %%a>> %VF_FILE%
   )
)
if defined TMP_VF_CNTENT goto :loop

echo ^">> %VF_FILE%
echo     }>> %VF_FILE%
echo }>> %VF_FILE%

ENDLOCAL

CD /d C:\Programs\caddy\
@echo on
caddy.exe validate
caddy.exe reload
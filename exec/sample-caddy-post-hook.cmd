@echo off

CD /d C:\Programs\caddy\
@echo on
caddy.exe validate
caddy.exe reload
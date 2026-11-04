@echo off
title C2F Admin Web App
cd /d "%~dp0src\c2f_admin_webapp"
echo Starting C2F Admin Web App on http://localhost:5005...
echo Connected to AWS RDS MySQL (C2F_NEW)
c2f_admin_webapp.exe
pause


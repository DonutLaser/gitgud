@echo off
go build -ldflags -H=windowsgui
ResourceHacker -open git-client.exe -save git-client.exe -action addskip -res assets/images/icon.ico -mask ICONGROUP,MAIN,
xcopy /s /y assets\* D:\Programo\custom\git-client\assets\
xcopy /y git-client.exe D:\Programos\custom\git-client\
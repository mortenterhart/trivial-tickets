@echo off

set server_name=trivial-tickets
set dirname=%~dp0

IF EXIST "cmd\ticketsystem" (
    call :info Checking for missing Go dependencies
    go get -t -v ./...

    IF %ERRORLEVEL% == 0 (
        cd cmd\ticketsystem

        call :info Starting Ticketsystem webserver
        go run ticketsystem.go

        cd ../..
    )
) ELSE (
    call :error cannot find the main executable cmd\ticketsystem\ticketsystem.go
    echo You might be in the wrong working directory, execute
    echo   cd %dirname%
    echo to change to the correct directory.
    exit /B 1
)

exit /B 0

:info
    echo [%server_name%] INFO %*
exit /B 0

:error
    echo [%server_name%] ERROR %*
exit /B 0
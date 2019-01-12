@echo off

REM Ticketsystem Trivial Tickets
REM
REM Matriculation numbers: 3040018, 3040018, 3478222
REM Lecture:               Programmieren II, INF16B
REM Lecturer:              Herr Prof. Dr. Helmut Neemann
REM Institute:             Duale Hochschule Baden-Württemberg Mosbach
REM
REM ---------------
REM CLI Start script


set program_name=commandLineInterface
set dirname=%~dp0

IF EXIST "cmd\commandLineTool" (
    call :info Checking for missing Go dependencies
    go get -t -v ./...

    IF %ERRORLEVEL% == 0 (

        call :info Starting CLI
        go run cmd\commandLineTool\commandLineTool.go

    )
) ELSE (
    call :error cannot find the main executable cmd\commandLineTool\commandLineTool.go
    echo You might be in the wrong working directory, execute
    echo   cd %dirname%
    echo to change to the correct directory.
    exit /B 1
)

exit /B 0

:info
    echo [%program_name%] INFO %*
exit /B 0

:error
    echo [%program_name%] ERROR %*
exit /B 0
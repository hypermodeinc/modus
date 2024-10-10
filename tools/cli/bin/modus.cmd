@echo off
setlocal EnableDelayedExpansion

:: Properties
set DIR=%~dp0
set OS=win

:: Config
set NODE_WANTED=22.9.0
set NODE_MIN=22
set NODE_INSTALL_PATH=%DIR%node-bin
set ARCH=%PROCESSOR_ARCHITECTURE%

set "BASE_DIR=%DIR:~0,-4%"
echo DIR %DIR%
echo BASE_DIR %BASE_DIR%
set "PATH_FILE=%DIR%scripts\modus-path"
set "LOCAL_FILE=%DIR%scripts\modus-local"

echo LOCAL_FILE %LOCAL_FILE%
:: Determine architecture
if /I "%ARCH%"=="AMD64" (
    set ARCH=x64
) else if /I "%ARCH%"=="ARM" (
    set ARCH=arm
) else if /I "%ARCH%"=="ARM64" (
    set ARCH=arm64
) else (
    echo Unsupported architecture: %ARCH%
    exit /b 1
)


:: Entry point
echo Starting script...
call :check_node
if errorlevel 1 (
    echo "local file"
    ::if not exist "%LOCAL_FILE%" (
        mkdir "%DIR%scripts\" >nul 2>&1
        (
            echo import { execute } from "@oclif/core";
            echo await execute({ dir: import.meta.url }^);
        ) > "%LOCAL_FILE%.mjs"
        (
            echo @echo off
            echo %NODE_INSTALL_PATH%\bin\node.exe "%%~dp0modus-local.mjs" %*
        ) > "%LOCAL_FILE%.cmd"
    ::)
    call :install_deps
    %~dp0scripts\modus-local.cmd %*
    exit
) else (
    echo "path file"
    if not exist "%PATH_FILE%" (
        mkdir "%DIR%scripts\" >nul 2>&1
        (
            echo import { execute } from "@oclif/core";
            echo await execute({ dir: import.meta.url }^);
        ) > "%PATH_FILE%.mjs"
        (
            echo @echo off
            echo node.exe "%%~dp0modus-path.mjs" %%*
        ) > "%PATH_FILE%.cmd"
    )
    call :install_deps
    %~dp0scripts\modus-path.cmd %*
    exit
)

endlocal

:: Install Node
:install_node
    set DOWNLOAD_FILE=node-v%NODE_WANTED%-%OS%-%ARCH%.zip
    set DOWNLOAD_URL=https://nodejs.org/dist/v%NODE_WANTED%/%DOWNLOAD_FILE%
    echo Installing Node from %DOWNLOAD_URL%

    set TEMP_DIR=%BASE_DIR%.modus-temp
    if exist "%TEMP_DIR%" (
        echo Removing existing temporary directory...
        rmdir /s /q "%TEMP_DIR%"
    )
    mkdir "%TEMP_DIR%" || (echo Failed to create temporary directory & exit /b 1)

    :: Download Node.js
    echo Downloading Node.js...
    curl --progress-bar --show-error --location --fail "%DOWNLOAD_URL%" --output "%TEMP_DIR%\%DOWNLOAD_FILE%"
    if errorlevel 1 (
        echo Failed to download Node.js from %DOWNLOAD_URL%
        exit /b 1
    )

    echo Unpacking archive...
    tar -xf "%TEMP_DIR%\%DOWNLOAD_FILE%" -C "%TEMP_DIR%"
    if errorlevel 1 (
        echo Failed to unpack Node.js archive
        exit /b 1
    )

    echo Cleaning up old installation...
    rmdir /s /q "%NODE_INSTALL_PATH%\bin"
    mkdir "%NODE_INSTALL_PATH%\bin"
    if errorlevel 1 (
        echo Failed to create Node.js installation directory
        exit /b 1
    )

    echo Moving Node.js files...
    del "%TEMP_DIR%\%DOWNLOAD_FILE%"
    xcopy "%TEMP_DIR%\node-v%NODE_WANTED%-%OS%-%ARCH%\*" "%NODE_INSTALL_PATH%\bin" /E /I >nul 2>&1
    rmdir /s /q "%TEMP_DIR%\node-v%NODE_WANTED%-%OS%-%ARCH%"
    if errorlevel 1 (
        echo Failed to move Node.js files to installation directory
        exit /b 1
    )

    echo Cleaning up temporary files...
    rmdir /s /q "%TEMP_DIR%"
exit /b

:: Check Node.js version
:check_node
    echo Checking Node.js installation...
    where node >nul 2>&1
    if errorlevel 1 (
        if exist "%NODE_INSTALL_PATH%\bin\node.exe" (
            echo Node.js not found in %NODE_INSTALL_PATH%
            call :install_node
            exit /b 1
        )
    ) else (
        if not exist "%NODE_INSTALL_PATH%\bin\node.exe" (
            :: Check the version of Node.js
            for /f "delims=v" %%v in ('node -v 2^>nul') do (
                set "NODE_VERSION=%%v"
                for /f "tokens=1 delims=." %%a in ("!NODE_VERSION!") do set "MAJOR_VERSION=%%a"
                if !MAJOR_VERSION! leq %NODE_MIN% (
                    call :install_node
                    exit /b 1
                )
            )
        )
    )
exit /b

:: Install dependencies
:install_deps
    if not exist "%BASE_DIR%node_modules" (
        echo Installing dependencies with npm...
        where npm >nul 2>&1
        if errorlevel 1 (
            echo "%~dp0\node-bin\bin\npm" --prefix "%BASE_DIR%" install install
            "%~dp0\node-bin\bin\npm" --prefix "%BASE_DIR%" install >nul 2>&1
        ) else (
            echo npm --prefix "%BASE_DIR%" install
            npm --prefix "%BASE_DIR%" install >nul 2>&1            
        )
    )
exit /b
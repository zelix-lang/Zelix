:: This script corresponds to the build process of the Surf project.
:: This code is released under the GNU GPL v3 license.
:: For more information, please visit: https://www.gnu.org/licenses/gpl-3.0.html
:: ----------------------------------------------------------------
:: Copyright (C) 2024 Rodrigo R. & All Surf Contributors
:: This program comes with ABSOLUTELY NO WARRANTY; for details type `surf license`.
:: This is free software, and you are welcome to redistribute it under certain conditions;
:: type `surf license --full` for details.
:: ----------------------------------------------------------------

@echo off
echo Building Surf...

:: Delete the bin folder if it exists
if exist bin (
    rmdir /s /q bin
)

:: Create the bin folder and build the project
mkdir bin
cd bin

go build -o surf
echo -> ./bin/surf

cd ..
echo Done.

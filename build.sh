# This script corresponds to the build process of the Zyro project.
# This code is released under the GNU GPL v3 license.
# For more information, please visit: https://www.gnu.org/licenses/gpl-3.0.html
# ----------------------------------------------------------------
# Copyright (C) 2024 Rodrigo R. & All Zyro Contributors
# This program comes with ABSOLUTELY NO WARRANTY; for details type `zyro license`.
# This is free software, and you are welcome to redistribute it under certain conditions;
# type `zyro license --full` for details.
# ----------------------------------------------------------------

echo Building Zyro...

# Delete the bin folder if it exists
if [ -d bin ]; then
    rm -r bin
fi

# Create the bin folder and build the project
mkdir bin

go build -o bin/zyro/
echo "-> ./bin/zyro/"

echo Done.
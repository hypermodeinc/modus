#!/bin/bash
rm -rf build
mkdir -p build

create_ignore_list() {
    if [[ -f ".packignore" ]]; then
        while IFS= read -r line; do
            if [[ -n "$line" && ! "$line" =~ ^# ]]; then
                echo "$line"
            fi
        done < .packignore
    fi
}

VERSION=$(grep -oP '"version": *"\K[0-9.]+(?=")' package.json)
if [[ $? -ne 0 || -z "$VERSION" ]]; then
    echo "Failed to get version from package.json or package.json not found."
    exit 1
fi

BASE_NAME="modus-cli-v${VERSION}"
TAR_FILE="${BASE_NAME}.tar.gz"
ZIP_FILE="${BASE_NAME}.zip"

IGNORE_LIST=$(create_ignore_list)

if [[ -z "$IGNORE_LIST" ]]; then
    tar -czf "build/${TAR_FILE}" .
else
    EXCLUDE_OPTIONS=()
    while IFS= read -r line; do
        EXCLUDE_OPTIONS+=(--exclude="$line")
    done <<< "$IGNORE_LIST"

    tar -czf "build/${TAR_FILE}" "${EXCLUDE_OPTIONS[@]}" .
fi

if [ $? -eq 0 ]; then
    echo "build/${TAR_FILE} created successfully."
else
    echo "Failed to create build/${TAR_FILE}."
    exit 1
fi

if [[ -z "$IGNORE_LIST" ]]; then
    zip -r "build/${ZIP_FILE}" .
else
    zip -r "build/${ZIP_FILE}" . -x "$(echo "$IGNORE_LIST" | tr '\n' '\n-x ')" -x "*.tar.gz" -x "*.zip"
fi

if [ $? -eq 0 ]; then
    echo "build/${ZIP_FILE} created successfully."
else
    echo "Failed to create build/${ZIP_FILE}."
    exit 1
fi

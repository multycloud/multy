#!/bin/sh
#VERSION=$(curl --silent "https://api.github.com/repos/multycloud/multy/releases/latest" |  jq -r .tag_name)
VERSION="v0.0.1-alpha"

ARCH=""
case $(uname -m) in
    "x86_64") ARCH="amd64";;
    "arm64") ARCH="arm64";;
    "aarch64") ARCH="arm64";;
    *)
        printf "Unsupported platform"
        exit 1
        ;;
esac

OS=""
EXT=""
case $(uname) in
    "Linux") OS="linux";EXT="tar.gz";;
    "Windows") OS="windows";EXT="zip";;
    "Darwin") OS="darwin";EXT="tar.gz";;
    *)
        printf "Unsupported OS"
        exit 1
        ;;
esac

DOWNLOAD_URL="https://github.com/multycloud/multy/releases/download/${VERSION}/multy-${VERSION}-${OS}-${ARCH}.${EXT}"
TARBALL_DEST="multy-${VERSION}-${OS}-${ARCH}.${EXT}"

printf "Downloading multy version %s\\n" "${VERSION}"

if curl -s -L -o "${TARBALL_DEST}" "${DOWNLOAD_URL}"; then
    printf "Extracting to %s\\n" "$HOME/.multy/bin"

    # If `~/.multy/bin exists, delete it
    if [ -e "${HOME}/.multy/bin" ]; then
        rm -rf "${HOME}/.multy/bin"
    fi

    mkdir -p "${HOME}/.multy"

    EXTRACT_DIR=$(mktemp -d multy.XXXXXXXXXX)
    tar zxf "${TARBALL_DEST}" -C "${EXTRACT_DIR}"

    cp -r "${EXTRACT_DIR}/." "${HOME}/.multy/bin/"

    rm -f "${TARBALL_DEST}"
    rm -rf "${EXTRACT_DIR}"
    printf "Installation complete. You can now use ~/.multy/bin/multy to run multy.\\n"
else
    >&2  printf "error: failed to download %s\\n" "${DOWNLOAD_URL}"
    exit 1
fi
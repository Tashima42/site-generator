#!/bin/sh

set -e

TMP_DIR=""
REPO_NAME="site-generator"
REPO_URL="https://github.com/tashima42/${REPO_NAME}"
REPO_RELEASE_URL="${REPO_URL}/releases"
INSTALL_DIR="$HOME/.local/bin"
SUFFIX=""
DOWNLOADER=""


# setup_arch set arch and suffix fatal if architecture not supported.
setup_arch() {
    case $(uname -m) in
    x86_64|amd64)
        ARCH=x86_64
        SUFFIX=$(uname -s)_${ARCH}
        ;;
    aarch64|arm64)
        ARCH=arm64
        SUFFIX=$(uname -s)_${ARCH}
        ;;
    i386)
        ARCH=i386
        SUFFIX=$(uname -s)_${ARCH}
        ;;
    *)
        fatal "unsupported architecture ${ARCH}"
        ;;
    esac
}

# setup_tmp creates a temporary directory and cleans up when done.
setup_tmp() {
    TMP_DIR=$(mktemp -d -t site-generator-install)
    cleanup() {
        code=$?
        set +e
        trap - EXIT
        rm -rf "${TMP_DIR}"
        exit "$code"
    }
    trap cleanup INT EXIT
}

# verify_downloader verifies existence of network downloader executable.
verify_downloader() {
    cmd="$(command -v "${1}")"
    if [ -z "${cmd}" ]; then
        return 1
    fi
    if [ ! -x "${cmd}" ]; then
        return 1
    fi

    DOWNLOADER=${cmd}
    return 0
}

# download downloads a file from a url using either curl or wget.
download() {
    case "${DOWNLOADER}" in
    *curl)
        cd "$1" && { curl -fsSLO "$2" ; cd -; }
    ;;
    *wget)
        wget -qO "$1" "$2"
    ;;
    esac

    if [ $? -ne 0 ]; then
        echo "error: download failed"
        exit 1
    fi
}

# download_tarball downloads the tarbal for the given version.
download_tarball() {
    TARBALL_URL="${REPO_RELEASE_URL}/download/${RELEASE_VERSION}/site-generator_${SUFFIX}.tar.gz"

    echo "downloading tarball from ${TARBALL_URL}"
    
    download "${TMP_DIR}" "${TARBALL_URL}"
}

# install_binaries installs the binaries from the downloaded tar.
install_binaries() {
    cd "${TMP_DIR}"
    echo  "${TMP_DIR}"
    tar -xf "${TMP_DIR}/site-generator_${SUFFIX}.tar.gz"
    rm "${TMP_DIR}/site-generator_${SUFFIX}.tar.gz"
    mkdir -p "${INSTALL_DIR}"
    cp "${TMP_DIR}/site-generator" "${INSTALL_DIR}"
}

{ # main
    if [ -z "$1" ]; then 
        echo "error: release version required"
        exit 1
    fi
    RELEASE_VERSION=$1

    echo "Installing Site Generator: ${RELEASE_VERSION}"

    setup_tmp
    setup_arch

    verify_downloader curl || verify_downloader wget || fatal "error: cannot find curl or wget"
    download_tarball
    install_binaries

    printf "Run command to access tools:\n\nPATH=%s:%s\n\n" "${PATH}" "${INSTALL_DIR}"

    exit 0
}
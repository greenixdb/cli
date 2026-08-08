#!/bin/sh
# Greenix CLI installer for Linux and macOS.
# Publisher: Edgicode Limited
#
#   curl -fsSL https://raw.githubusercontent.com/greenixdb/cli/main/scripts/install.sh | sh
#
# Environment overrides:
#   GREENIX_VERSION   release tag to install (default: latest)
#   GREENIX_INSTALL   install directory (default: /usr/local/bin, else ~/.local/bin)
set -eu

REPO="greenixdb/cli"
BIN_NAME="greenix"
VERSION="${GREENIX_VERSION:-latest}"

red()  { printf '\033[31m%s\033[0m\n' "$1"; }
grn()  { printf '\033[32m%s\033[0m\n' "$1"; }
info() { printf '  %s\n' "$1"; }
die()  { red "error: $1" >&2; exit 1; }

command -v curl >/dev/null 2>&1 || command -v wget >/dev/null 2>&1 || die "curl or wget is required"
command -v tar  >/dev/null 2>&1 || die "tar is required"

fetch() {
  # fetch <url> <output-file>
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$1" -o "$2"
  else
    wget -qO "$2" "$1"
  fi
}

# ---------------------------------------------------------------- detect OS
case "$(uname -s)" in
  Linux)  PLATFORM="linux" ;;
  Darwin) PLATFORM="macos" ;;
  *) die "unsupported operating system: $(uname -s)" ;;
esac

# -------------------------------------------------------------- detect arch
case "$(uname -m)" in
  x86_64|amd64)   ARCH="x64" ;;
  aarch64|arm64)  ARCH="arm64" ;;
  i386|i686)      ARCH="x86" ;;
  armv7l|armv7)   ARCH="arm" ;;
  *) die "unsupported architecture: $(uname -m)" ;;
esac

# macOS ships a universal build; it works on both Intel and Apple Silicon.
[ "$PLATFORM" = "macos" ] && ARCH="universal"
[ "$PLATFORM" = "macos" ] && [ "$(uname -m)" = "i386" ] && die "32-bit macOS is not supported"

ASSET="${BIN_NAME}-${PLATFORM}-${ARCH}.tar.gz"

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

# ----------------------------------------------------- pick install location
if [ -n "${GREENIX_INSTALL:-}" ]; then
  INSTALL_DIR="$GREENIX_INSTALL"
elif [ -w /usr/local/bin ] 2>/dev/null; then
  INSTALL_DIR="/usr/local/bin"
elif [ "$(id -u)" = "0" ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
fi

printf '\n'
grn "Installing Greenix CLI"
info "platform : ${PLATFORM}/${ARCH}"
info "release  : ${VERSION}"
info "target   : ${INSTALL_DIR}/${BIN_NAME}"
printf '\n'

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM

info "downloading ${ASSET}..."
fetch "$URL" "$TMP_DIR/$ASSET" || die "download failed: $URL"

# ------------------------------------------------------- verify the checksum
if [ "$VERSION" = "latest" ]; then
  SUMS_URL="https://github.com/${REPO}/releases/latest/download/SHA256SUMS"
else
  SUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/SHA256SUMS"
fi

if fetch "$SUMS_URL" "$TMP_DIR/SHA256SUMS" 2>/dev/null; then
  EXPECTED="$(grep " \*\{0,1\}${ASSET}\$" "$TMP_DIR/SHA256SUMS" | awk '{print $1}' | head -n1)"
  if [ -n "$EXPECTED" ]; then
    if command -v sha256sum >/dev/null 2>&1; then
      ACTUAL="$(sha256sum "$TMP_DIR/$ASSET" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
      ACTUAL="$(shasum -a 256 "$TMP_DIR/$ASSET" | awk '{print $1}')"
    else
      ACTUAL=""
    fi
    if [ -n "$ACTUAL" ]; then
      [ "$ACTUAL" = "$EXPECTED" ] || die "checksum mismatch for ${ASSET}"
      info "checksum verified"
    fi
  fi
else
  info "warning: SHA256SUMS not available, skipping verification"
fi

# --------------------------------------------------------------- unpack
tar -xzf "$TMP_DIR/$ASSET" -C "$TMP_DIR"
[ -f "$TMP_DIR/$BIN_NAME" ] || die "archive did not contain '${BIN_NAME}'"
chmod +x "$TMP_DIR/$BIN_NAME"

# Binaries fetched by this script are not browser downloads, so macOS does not
# attach a quarantine flag. Strip it anyway in case the archive came from one.
if [ "$PLATFORM" = "macos" ] && command -v xattr >/dev/null 2>&1; then
  xattr -d com.apple.quarantine "$TMP_DIR/$BIN_NAME" 2>/dev/null || true
fi

mkdir -p "$INSTALL_DIR" 2>/dev/null || true
if [ -w "$INSTALL_DIR" ]; then
  mv -f "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
elif command -v sudo >/dev/null 2>&1; then
  info "elevating with sudo to write to ${INSTALL_DIR}"
  sudo mkdir -p "$INSTALL_DIR"
  sudo mv -f "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
  sudo chmod +x "$INSTALL_DIR/$BIN_NAME"
else
  die "cannot write to ${INSTALL_DIR} (set GREENIX_INSTALL to a writable directory)"
fi

printf '\n'
grn "Greenix CLI installed to ${INSTALL_DIR}/${BIN_NAME}"

# ------------------------------------------------------------ PATH guidance
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*)
    printf '\n'
    info "Run: ${BIN_NAME} --version"
    ;;
  *)
    printf '\n'
    red "${INSTALL_DIR} is not on your PATH."
    info "Add this line to ~/.bashrc, ~/.zshrc or ~/.profile, then open a new terminal:"
    printf '\n    export PATH="%s:$PATH"\n\n' "$INSTALL_DIR"
    info "Or run it directly right now: ${INSTALL_DIR}/${BIN_NAME} --version"
    ;;
esac
printf '\n'

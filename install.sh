#!/bin/sh
# tgctl installer for macOS/Linux.
#   curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/tgctl/main/install.sh | sh
#
# Downloads the release archive that matches your OS/arch, verifies its SHA-256 against the
# release checksums.txt, and installs the binary. Override behaviour with env vars:
#   TGCTL_VERSION=v0.2.0        pin a version (default: latest release)
#   TGCTL_INSTALL_DIR=~/.local/bin   install location (default: /usr/local/bin)
set -eu

REPO="jjuanrivvera/tgctl"
BINARY="tgctl"
INSTALL_DIR="${TGCTL_INSTALL_DIR:-/usr/local/bin}"

RED=''; GREEN=''; YELLOW=''; NC=''
if [ -t 1 ]; then RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'; fi
info()  { printf "${GREEN}[info]${NC} %s\n"  "$1"; }
warn()  { printf "${YELLOW}[warn]${NC} %s\n" "$1" >&2; }
die()   { printf "${RED}[error]${NC} %s\n"   "$1" >&2; exit 1; }

command -v curl >/dev/null 2>&1 || die "curl is required"
command -v tar  >/dev/null 2>&1 || die "tar is required"

# --- detect platform: match the goreleaser archive naming (lowercase os, amd64/arm64) ---
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  linux|darwin) ;;
  *) die "unsupported OS: $os (this installer covers Linux and macOS; use 'go install' or Docker otherwise)" ;;
esac
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  aarch64|arm64) arch="arm64" ;;
  *) die "unsupported architecture: $arch" ;;
esac
info "detected platform: ${os}/${arch}"

# --- resolve version ---
version="${TGCTL_VERSION:-}"
if [ -z "$version" ]; then
  info "fetching the latest release tag..."
  version="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
  [ -n "$version" ] || die "could not determine the latest release; set TGCTL_VERSION"
fi
num="${version#v}"   # goreleaser archives are named without the leading 'v'
info "installing ${BINARY} ${version}"

archive="${BINARY}_${num}_${os}_${arch}.tar.gz"
base="https://github.com/${REPO}/releases/download/${version}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

info "downloading ${archive}"
curl -fsSL "${base}/${archive}"       -o "${tmp}/${archive}" || die "failed to download ${archive}"
curl -fsSL "${base}/checksums.txt"    -o "${tmp}/checksums.txt" || die "failed to download checksums.txt"

# --- verify the SHA-256 against checksums.txt ---
expected="$(grep " ${archive}\$" "${tmp}/checksums.txt" | awk '{print $1}' | head -1)"
[ -n "$expected" ] || die "no checksum for ${archive} in checksums.txt"
if command -v sha256sum >/dev/null 2>&1; then
  actual="$(sha256sum "${tmp}/${archive}" | awk '{print $1}')"
elif command -v shasum >/dev/null 2>&1; then
  actual="$(shasum -a 256 "${tmp}/${archive}" | awk '{print $1}')"
else
  die "need sha256sum or shasum to verify the download"
fi
[ "$expected" = "$actual" ] || die "checksum mismatch for ${archive} (expected ${expected}, got ${actual})"
info "checksum verified"

# --- extract & install ---
tar -xzf "${tmp}/${archive}" -C "$tmp"
[ -f "${tmp}/${BINARY}" ] || die "archive did not contain ${BINARY}"
chmod +x "${tmp}/${BINARY}"

if [ -w "$INSTALL_DIR" ]; then
  mv "${tmp}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  warn "elevating with sudo to write ${INSTALL_DIR}"
  sudo mv "${tmp}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi
info "installed to ${INSTALL_DIR}/${BINARY}"

# --- verify install ---
if command -v "$BINARY" >/dev/null 2>&1; then
  "$BINARY" version || true
  info "done. Next: ${BINARY} auth login"
else
  warn "${INSTALL_DIR} is not on your PATH — add it: export PATH=\"${INSTALL_DIR}:\$PATH\""
fi

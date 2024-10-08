#!/usr/bin/env bash

# Copyright 2024 Hypermode, Inc.
# Licensed under the terms of the Apache License, Version 2.0
# See the LICENSE file that accompanied this code for further details.
#
# SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
# SPDX-License-Identifier: Apache-2.0

# I'm hosting builds on my own fork until it goes live

# Config
GIT_REPO="JairusSW/modus-cli"
INSTALL_DIR="${MODUS_CLI:-"$HOME/.modus/cli"}"
VERSION="latest"

# Properties
ARCH="$(uname -m)"
OS="$(uname -s)"

get_latest_release() {
  curl -w "%{stderr}" --silent "https://api.github.com/repos/$GIT_REPO/releases/latest" |
    grep '"tag_name"' |
    sed -E 's/.*"([^"]+)".*/\1/'
}

download_release_from_repo() {
  local version="$1"
  local arch="$2"
  local os="$3"
  local tmpdir="$4"
  local filename="modus-$version-$os-$arch.tar.gz"
  local download_file="$tmpdir/$filename"
  local archive_url="https://github.com/$GIT_REPO/releases/download/$version/$filename"

  curl --progress-bar --show-error --location --fail "$archive_url" \
    --output "$download_file" && echo "$download_file"
}

echo_fexists() {
  [ -f "$1" ] && echo "$1"
}

detect_profile() {
  local shellname="$1"
  local uname="$2"

  if [ -f "$PROFILE" ]; then
    echo "$PROFILE"
    return
  fi

  case "$shellname" in
  bash)
    case $uname in
    Darwin)
      echo_fexists "$HOME/.bash_profile" || echo_fexists "$HOME/.bashrc"
      ;;
    *)
      echo_fexists "$HOME/.bashrc" || echo_fexists "$HOME/.bash_profile"
      ;;
    esac
    ;;
  zsh)
    echo "$HOME/.zshrc"
    ;;
  fish)
    echo "$HOME/.config/fish/config.fish"
    ;;
  *)
    local profiles
    case $uname in
    Darwin)
      profiles=(.profile .bash_profile .bashrc .zshrc .config/fish/config.fish)
      ;;
    *)
      profiles=(.profile .bashrc .bash_profile .zshrc .config/fish/config.fish)
      ;;
    esac

    for profile in "${profiles[@]}"; do
      echo_fexists "$HOME/$profile" && break
    done
    ;;
  esac
}

build_path_str() {
  local profile="$1" install_dir="$2"
  if [[ $profile == *.fish ]]; then
    echo -e "set -gx MODUS_CLI \"$install_dir\"\nstring match -r \".modus\" \"\$PATH\" > /dev/null; or set -gx PATH \"\$MODUS_CLI/bin\" \$PATH"
  else
    echo -e "# Modus CLI\nexport MODUS_CLI=\"$install_dir\"\nexport PATH=\"\$MODUS_CLI/bin:\$PATH\""
  fi
}

cli_dir_valid() {
  if [ -n "${MODUS_CLI-}" ] && [ -e "$MODUS_CLI" ] && ! [ -d "$MODUS_CLI" ]; then
    error "\$MODUS_CLI is set but is not a directory ($MODUS_CLI)."
    eprintf "Please check your profile scripts and environment."
    return 1
  fi
  return 0
}

update_profile() {
  local install_dir="$1" detected_profile="$(detect_profile $(basename "$SHELL") $(uname -s))"
  local path_str="$(build_path_str "$detected_profile" "$install_dir")"

  if [ -z "$detected_profile" ]; then
    echo "No user profile found."
    echo "$path_str"
    return 1
  fi

  if ! grep -q 'MODUS_CLI' "$detected_profile"; then
    printf "%s\n" "$path_str" >>"$detected_profile"
  fi
  echo "[3/5] Added modus to PATH"
}

install_version() {
  if ! cli_dir_valid; then
    exit 1
  fi

  case "$VERSION" in
  latest)
    VERSION="$(get_latest_release)"
    ;;
  *) ;;
  esac

  echo -e "${BOLD}${BLUE}Modus${RESET} Installer ${DIM}(${VERSION})${RESET}\n"

  install_release
  if [ "$?" == 0 ]; then
    update_profile "$INSTALL_DIR" &&
      echo "[4/5] Installed Modus CLI"
  fi
}

install_release() {
  echo -e "[1/5] Fetching archive for $OS $ARCH"
  download_archive="$(
    download_release
    exit "$?"
  )"
  exit_status="$?"
  if [ "$exit_status" != 0 ]; then
    error "Could not download Modus version '$VERSION'. See https://github.com/$GIT_REPO/releases/ for a list of available releases"
    return "$exit_status"
  fi

  clear_line
  clear_line

  echo "[1/5] Fetched archive for $OS $ARCH"

  install_from_file "$download_archive"
}

download_release() {
  if [ "$?" != 0 ]; then
    error "The current operating system ($OS) does not appear to be supported by Modus."
    return 1
  fi

  local download_dir="$(mktemp -d)"
  download_release_from_repo "$VERSION" "$ARCH" "$OS" "$download_dir"
}

install_from_file() {
  local archive="$1"
  local extract_to="$(dirname "$archive")"

  tar -xf "$archive" -C "$extract_to"

  rm -rf "$INSTALL_DIR"
  mkdir -p "$INSTALL_DIR"
  mv "$extract_to/modus/"* "$INSTALL_DIR"
  rm -rf "$extract_to"
  rm -f "$archive"

  echo "[2/5] Unpacked archive"
}

check_platform() {
  case $ARCH in
  aarch64) ARCH="arm64" ;;
  x86_64) ARCH="x64" ;;
  armv6l) ARCH="arm" ;;
  *) ;;
  esac

  case "$ARCH/$OS" in
  x64/Linux | arm64/Linux | arm/Linux | x64/Darwin | arm64/Darwin) return 0 ;;
  *)
    echo -e "Unsupported os $OS $ARCH"
    return 1
    ;;
  esac
}

restart_shell() {
  local shell_name="$(basename "$SHELL")"
  echo -e "[5/5] Restarted shell ${DIM}($shell_name)${RESET}\n\nThe Modus CLI has been installed! 🎉\nRun ${DIM}modus${RESET} to get started"

  case "$shell_name" in
    bash|zsh|fish) exec "$shell_name" ;;
    *)
      echo -e "[5/5] Clean up\n\nPlease restart your shell for changes to take effect"
      ;;
  esac
}

clear_line() {
  echo -ne "\033[F\033[K"
}


BOLD="\e[1m"
BLUE="\e[34;1m"
DIM="\e[2m"
RESET="\e[0m"

# This is the entry point

check_platform || exit 1

install_version

restart_shell

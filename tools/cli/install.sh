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
ARCH=""
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

# If file exists, echo it
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
    # Shells on macOS default to opening with a login shell, while Linuxes
    # default to a *non*-login shell, so if this is macOS we look for
    # `.bash_profile` first; if it's Linux, we look for `.bashrc` first. The
    # `*` fallthrough covers more than just Linux: it's everything that is not
    # macOS (Darwin). It can be made narrower later if need be.
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
    # Fall back to checking for profile file existence. Once again, the order
    # differs between macOS and everything else.
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
  echo "[3/4] Added modus to PATH"
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
      echo "[4/4] Installed Modus CLI"
  fi

  echo -e "\nThe Modus CLI has been installed! 🎉"
}

install_release() {
  echo "[1/4] Fetching archive for $OS $ARCH"
  # add progress percent above this
  download_archive="$(
    download_release
    exit "$?"
  )"
  exit_status="$?"
  if [ "$exit_status" != 0 ]; then
    error "Could not download Modus version '$VERSION'. See $(release_url) for a list of available releases"
    return "$exit_status"
  fi

  clear_line
  clear_line

  echo "[1/4] Fetched archive for $OS $ARCH"

  install_from_file "$download_archive"
}

download_release() {
  local arch="$(get_architecture)"
  local uname_str="$(uname -s)"
  if [ "$?" != 0 ]; then
    error "The current operating system ($uname_str) does not appear to be supported by Modus."
    return 1
  fi

  local download_dir="$(mktemp -d)"
  download_release_from_repo "$VERSION" "$arch" "$uname_str" "$download_dir"
}

install_from_file() {
  local archive="$1"
  local extract_to="$(dirname "$archive")"

  tar -xf "$archive" -C "$extract_to"

  rm -rf "$INSTALL_DIR"
  mkdir -p "$INSTALL_DIR"
  mv "$extract_to/modus/"* "$INSTALL_DIR"
  echo "$extract_to -> $INSTALL_DIR"
  rm -rf "$extract_to"
  rm -f "$archive"

  clear_line
  echo "[2/4] Unpacked archive"
}

get_architecture() {
  case "$(uname -m)" in
  aarch64) echo arm64 ;;
  x86_64) echo x64 ;;
  armv6l) echo arm ;;
  *)
    echo "Unsupported architecture"
    exit 1
    ;;
  esac
}

check_architecture() {
  local os="$(uname)"
  case "$ARCH/$os" in
  x64/Linux | arm64/Linux | arm/Linux | x64/Darwin | arm64/Darwin) return 0 ;;
  *)
    error "Unsupported architecture."
    return 1
    ;;
  esac
}

restart_shell() {
  case "$(basename "$SHELL")" in
  bash)
    echo -e "Run ${DIM}modus${RESET} to get started"
    exec bash
    ;;
  zsh)
    echo -e "Run ${DIM}modus${RESET} to get started"
    exec zsh
    ;;
  fish)
    echo -e "Run ${DIM}modus${RESET} to get started"
    exec fish
    ;;
  *)
    echo "Please restart your shell for changes to take effect"
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

ARCH="$(get_architecture)"

# This is the entry point

check_architecture || exit 1

install_version

restart_shell
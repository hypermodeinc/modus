#!/bin/bash

get_latest_release() {
  curl -w "%{stderr}" --silent "https://api.github.com/repos/hypermodeai/runtime/releases/latest" | \
    tr -d '\n' | \
    sed 's/.*tag_name": *"//' | \
    sed 's/".*//'
}

# Function to detect the shell profile
detect_profile() {
    local shellname="$(basename "$SHELL")"
    local uname="$(uname -s)"
    
    # Check if a file exists
    echo_fexists() {
        [ -f "$1" ] && echo "$1"
    }

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
                    profiles=( .profile .bash_profile .bashrc .zshrc .config/fish/config.fish )
                    ;;
                *)
                    profiles=( .profile .bashrc .bash_profile .zshrc .config/fish/config.fish )
                    ;;
            esac
            for profile in "${profiles[@]}"; do
                echo_fexists "$HOME/$profile" && break
            done
            ;;
    esac
}

# Get the download URL based on OS and architecture
get_download_url() {
    local os="$(uname -s)"
    local arch="$(uname -m)"
     
    case "$os" in
        Linux)
            case "$arch" in
                x86_64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-d626ac6-linux-x64.tar.gz";;
                arm64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-d626ac6-linux-arm64.tar.gz";;
                arm*) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-d626ac6-linux-arm.tar.gz";;
                *) echo "Unsupported architecture: $arch"; exit 1;;
            esac
            ;;
        Darwin)
            case "$arch" in
                x86_64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-d626ac6-darwin-x64.tar.gz";;
                arm64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-d626ac6-darwin-arm64.tar.gz";;
                *) echo "Unsupported architecture: $arch"; exit 1;;
            esac
            ;;
        *)
            echo "Unsupported OS: $os"; exit 1;;
    esac
}

# Get the download URL
URL=$(get_download_url)
INSTALL_DIR="$HOME/.modus/"
TARBALL="modus.tar.gz"

# Create the installation directory if it doesn't exist
mkdir -p "$INSTALL_DIR/bin"

# Download the tarball
echo "Downloading modus CLI from $URL..."
curl -L -o "$TARBALL" "$URL"

# Extract the tarball
echo "Extracting modus..."
tar -xzf "$TARBALL"

# Move the extracted binary to the installation directory
echo "Installing modus CLI..."
mv modus "$INSTALL_DIR/"

# Make the binary executable
chmod +x "$INSTALL_DIR/modus/bin/modus"

# Add to PATH (for the current session)
export PATH="$INSTALL_DIR/modus/bin:$PATH"

# Detect the appropriate shell profile
PROFILE=$(detect_profile)

# Append to the profile if the modus PATH isn't already added
if ! grep -q "$INSTALL_DIR/modus/bin" "$PROFILE"; then
    echo "Adding modus to your PATH in $PROFILE"
    echo "export PATH=\"$INSTALL_DIR/modus/bin:\$PATH\"" >> "$PROFILE"
else
    echo "Modus is already in your PATH."
fi

# Clean up
echo "Cleaning up..."
rm "$TARBALL"

# Verify installation
if command -v modus &> /dev/null; then
    echo "Modus CLI installed successfully!"
    echo "To activate the new PATH in the current session, run:"
    echo "source $PROFILE"
else
    echo "Installation failed."
fi

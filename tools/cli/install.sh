#!/bin/bash

get_download_url() {
    local os="$(uname -s)"
    local arch="$(uname -m)"
    
    case "$os" in
        Linux)
            case "$arch" in
                x86_64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-1aa6e02-linux-x64.tar.gz";;
                arm64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-1aa6e02-linux-arm64.tar.gz";;
                arm*) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-1aa6e02-linux-arm.tar.gz";;
                *) echo "Unsupported architecture: $arch"; exit 1;;
            esac
            ;;
        Darwin)
            case "$arch" in
                x86_64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-1aa6e02-darwin-x64.tar.gz";;
                arm64) echo "http://174.138.95.78:3000/dist/modus-v0.0.0-1aa6e02-darwin-arm64.tar.gz";;
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

# Clean up
echo "Cleaning up..."
rm "$TARBALL"

# Verify installation
if command -v modus &> /dev/null; then
    echo "Modus CLI installed successfully!"
    echo "To add modus to your PATH permanently, add the following line to your shell configuration file (e.g., ~/.bashrc):"
    echo "export PATH=\"$INSTALL_DIR/modus/bin:\$PATH\""
else
    echo "Installation failed."
fi

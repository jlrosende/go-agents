#!/bin/bash

# Substitute BIN for your bin directory.
# Substitute VERSION for the current released version.
BIN="/usr/local/bin" && \
VERSION="1.54.0" && \
sudo curl -sSL \
"https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" \
-o "${BIN}/buf" && \
sudo chmod +x "${BIN}/buf"
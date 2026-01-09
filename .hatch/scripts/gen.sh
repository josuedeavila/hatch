#!/usr/bin/env bash
#
# gen.sh - Generate editor-compatible files from .hatch source files
#
# Calls all generators in .hatch/scripts/generators/
#
# Usage: ./.hatch/scripts/gen.sh
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GENERATORS_DIR="$SCRIPT_DIR/generators"

# Colors
GREEN='\033[0;32m'
DIM='\033[2m'
NC='\033[0m'

echo -e "${DIM}Generating from .hatch/...${NC}"

# Run all generator scripts in the generators directory
if [[ -d "$GENERATORS_DIR" ]]; then
    for generator in "$GENERATORS_DIR"/*.sh; do
        if [[ -f "$generator" && -x "$generator" ]]; then
            "$generator"
        elif [[ -f "$generator" ]]; then
             # Try running with bash if not executable
             bash "$generator"
        fi
    done
else
    echo "Error: Generators directory not found at $GENERATORS_DIR"
    exit 1
fi

echo -e "${GREEN}✓${NC} Generation complete"

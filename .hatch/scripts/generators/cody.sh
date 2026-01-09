#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
CODY_MD="$PROJECT_ROOT/CODY.md"

# Colors
MAGENTA='\033[0;35m'
NC='\033[0m'

# Generate CODY.md from .hatch/src/rules
generate_cody_md() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"

    [[ ! -d "$rules_dir" ]] && return

    {
        if [[ -d "$rules_dir" ]]; then
            echo "# Project Instructions"
            echo ""
            echo "This file is auto-generated from .hatch/src/rules/ - do not edit directly."
            echo ""
            
            find "$rules_dir" -name "*.md" -type f | sort | while read -r src_file; do
                sed 's/^#/##/' "$src_file"
                echo ""
                echo ""
            done
        fi

        if [[ -d "$commands_dir" ]]; then
            echo "# Commands"
            echo ""
            echo "These are the available commands for this project."
            echo ""
            
            find "$commands_dir" -name "*.md" -type f | sort | while read -r src_file; do
                sed 's/^#/##/' "$src_file"
                echo ""
                echo ""
            done
        fi
    } > "$CODY_MD"
    
    echo -e "${MAGENTA}cody${NC} CODY.md"
}

# Copy ignore file for Cody
copy_cody_ignore() {
    local src_file="$HATCH_DIR/src/.ignore"
    local dest_file="$PROJECT_ROOT/.cody/ignore"

    [[ ! -f "$src_file" ]] && return

    mkdir -p "$PROJECT_ROOT/.cody"
    cp "$src_file" "$dest_file"
    echo -e "${MAGENTA}cody${NC} .cody/ignore"
}

generate_cody_md
copy_cody_ignore

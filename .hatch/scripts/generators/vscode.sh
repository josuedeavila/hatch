#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
INSTRUCTIONS_MD="$PROJECT_ROOT/INSTRUCTIONS.md"

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

# Generate INSTRUCTIONS.md from .hatch/src/rules
generate_instructions_md() {
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
    } > "$INSTRUCTIONS_MD"
    
    echo -e "${GREEN}general${NC} INSTRUCTIONS.md"
}

generate_instructions_md

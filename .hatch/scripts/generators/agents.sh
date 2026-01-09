#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
AGENTS_MD="$PROJECT_ROOT/AGENTS.md"

# Colors
MAGENTA='\033[0;35m'
NC='\033[0m'

# Generate AGENTS.md from .hatch/rules
generate_agents_md() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"

    [[ ! -d "$rules_dir" ]] && return

    {
        if [[ -d "$rules_dir" ]]; then
            echo "# Agent Instructions"
            echo ""
            echo "This file is auto-generated from .hatch/rules/ - do not edit directly."
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
    } > "$AGENTS_MD"
    
    echo -e "${MAGENTA}agents${NC} AGENTS.md"
}

generate_agents_md

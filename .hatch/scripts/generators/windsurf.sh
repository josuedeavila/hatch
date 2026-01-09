#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
WINDSURF_DIR="$PROJECT_ROOT/.windsurf"

# Colors
CYAN='\033[0;36m'
NC='\033[0m'

# Generate Windsurf rules from .hatch/src/rules
generate_windsurf_rules() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"
    local dest_dir="$WINDSURF_DIR"
    local dest_file="$dest_dir/rules.md"

    mkdir -p "$dest_dir"

    {
        if [[ -d "$rules_dir" ]]; then
            echo "# Project Rules"
            echo ""
            echo "These are the project rules and guidelines."
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
    } > "$dest_file"
    
    echo -e "${CYAN}windsurf${NC} .windsurf/rules.md"
}

generate_windsurf_rules

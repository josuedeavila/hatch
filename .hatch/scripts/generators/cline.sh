#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
CLINE_RULES="$PROJECT_ROOT/.clinerules"

# Colors
PURPLE='\033[0;35m'
NC='\033[0m'

# Generate .clinerules from .hatch/src/rules and .hatch/src/commands
generate_cline_rules() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"

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
    } > "$CLINE_RULES"
    
    echo -e "${PURPLE}cline${NC} .clinerules"
}

# Copy ignore file for Cline
copy_cline_ignore() {
    local src_file="$HATCH_DIR/src/.ignore"
    local dest_file="$PROJECT_ROOT/.clineignore"

    [[ ! -f "$src_file" ]] && return

    cp "$src_file" "$dest_file"
    echo -e "${PURPLE}cline${NC} .clineignore"
}

generate_cline_rules
copy_cline_ignore

#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
CONTINUE_DIR="$PROJECT_ROOT/.continue"

# Colors
BLUE='\033[0;34m'
NC='\033[0m'

# Generate Continue system prompt from .hatch/src/rules
generate_continue_prompt() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"
    local dest_dir="$CONTINUE_DIR"
    local dest_file="$dest_dir/system_prompt.md"

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
    
    echo -e "${BLUE}continue${NC} .continue/system_prompt.md"
}

# Copy ignore file for Continue
copy_continue_ignore() {
    local src_file="$HATCH_DIR/src/.ignore"
    local dest_file="$PROJECT_ROOT/.continueignore"

    [[ ! -f "$src_file" ]] && return

    cp "$src_file" "$dest_file"
    echo -e "${BLUE}continue${NC} .continueignore"
}

generate_continue_prompt
copy_continue_ignore

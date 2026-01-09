#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
OPENCODE_DIR="$PROJECT_ROOT/.opencode"

# Colors
CYAN='\033[0;36m'
NC='\033[0m'

# Generate OpenCode rules from .hatch/src/rules
generate_opencode_rules() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"
    local dest_dir="$OPENCODE_DIR"
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
    
    echo -e "${CYAN}opencode${NC} .opencode/rules.md"
}

# Copy ignore file for OpenCode
copy_opencode_ignore() {
    local src_file="$HATCH_DIR/src/.ignore"
    local dest_file="$PROJECT_ROOT/.opencodeignore"

    [[ ! -f "$src_file" ]] && return

    cp "$src_file" "$dest_file"
    echo -e "${CYAN}opencode${NC} .opencodeignore"
}

generate_opencode_rules
copy_opencode_ignore

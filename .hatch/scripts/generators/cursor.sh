#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
CURSOR_DIR="$PROJECT_ROOT/.cursor"

# Colors
CYAN='\033[0;36m'
YELLOW='\033[0;33m'
NC='\033[0m'

# Generate Cursor rules from .hatch/rules
generate_cursor_rules() {
    local src_dir="$HATCH_DIR/src/rules"
    local dest_dir="$CURSOR_DIR/rules"

    [[ ! -d "$src_dir" ]] && return

    mkdir -p "$dest_dir"

    find "$src_dir" -name "*.md" -type f | sort | while read -r src_file; do
        local rel_path="${src_file#$src_dir/}"
        local dest_file="$dest_dir/${rel_path%.md}.mdc"
        
        mkdir -p "$(dirname "$dest_file")"
        
        {
            echo "---"
            echo "alwaysApply: true"
            echo "---"
            echo ""
            cat "$src_file"
        } > "$dest_file"
        
        echo -e "${CYAN}cursor rule${NC} ${rel_path%.md}.mdc"
    done
}

# Generate Cursor commands/prompts from .hatch/commands
generate_cursor_commands() {
    local src_dir="$HATCH_DIR/src/commands"
    local dest_dir="$CURSOR_DIR/commands"

    [[ ! -d "$src_dir" ]] && return

    mkdir -p "$dest_dir"

    find "$src_dir" -name "*.md" -type f | sort | while read -r src_file; do
        local rel_path="${src_file#$src_dir/}"
        local dest_file="$dest_dir/$rel_path"
        
        mkdir -p "$(dirname "$dest_file")"
        cp "$src_file" "$dest_file"
        
        echo -e "${YELLOW}cursor command${NC} ${rel_path}"
    done
}

generate_cursor_rules
generate_cursor_commands

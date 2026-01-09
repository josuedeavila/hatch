#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
ZED_DIR="$PROJECT_ROOT/.zed"

# Colors
BLUE='\033[0;34m'
NC='\033[0m'

# Generate Zed project rules prompt from .hatch/rules
# This creates a single prompt containing all rules that can be used with /rules
generate_zed_rules() {
    local src_dir="$HATCH_DIR/src/rules"
    local dest_dir="$ZED_DIR/prompts"
    local dest_file="$dest_dir/rules.md"

    [[ ! -d "$src_dir" ]] && return

    mkdir -p "$dest_dir"

    {
        echo "# Project Rules"
        echo ""
        echo "These are the project rules and guidelines. Include this context when working on this codebase."
        echo ""
        
        find "$src_dir" -name "*.md" -type f | sort | while read -r src_file; do
            sed 's/^#/##/' "$src_file"
            echo ""
            echo ""
        done
    } > "$dest_file"
    
    echo -e "${BLUE}zed prompt${NC} prompts/rules.md"
}

# Generate Zed prompts from .hatch/commands
generate_zed_commands() {
    local src_dir="$HATCH_DIR/src/commands"
    local dest_dir="$ZED_DIR/prompts"

    [[ ! -d "$src_dir" ]] && return

    mkdir -p "$dest_dir"

    find "$src_dir" -name "*.md" -type f | sort | while read -r src_file; do
        local rel_path="${src_file#$src_dir/}"
        local dest_file="$dest_dir/$rel_path"
        
        mkdir -p "$(dirname "$dest_file")"
        
        {
            cat "$src_file"
            echo ""
            echo ""
            echo "---"
            echo "Please follow the instructions above."
        } > "$dest_file"
        
        echo -e "${BLUE}zed prompt${NC} prompts/${rel_path}"
    done
}

generate_zed_rules
generate_zed_commands

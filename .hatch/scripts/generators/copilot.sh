#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
GITHUB_DIR="$PROJECT_ROOT/.github"

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

# Generate GitHub Copilot instructions from .hatch/rules
# This creates .github/copilot-instructions.md containing all rules
generate_copilot_instructions() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"
    local dest_file="$GITHUB_DIR/copilot-instructions.md"

    [[ ! -d "$rules_dir" ]] && return

    mkdir -p "$GITHUB_DIR"

    {
        if [[ -d "$rules_dir" ]]; then
            echo "# Project Instructions"
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
    } > "$dest_file"
    
    echo -e "${GREEN}copilot${NC} copilot-instructions.md"
}

# Generate GitHub Copilot prompts from .hatch/commands
generate_copilot_prompts() {
    local src_dir="$HATCH_DIR/src/commands"
    local dest_dir="$GITHUB_DIR/prompts"

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
        
        echo -e "${GREEN}copilot prompt${NC} prompts/${rel_path}"
    done
}

generate_copilot_instructions
generate_copilot_prompts

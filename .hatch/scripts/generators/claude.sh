#!/usr/bin/env bash
set -euo pipefail

# Determine paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

HATCH_DIR="$PROJECT_ROOT/.hatch"
CLAUDE_DIR="$PROJECT_ROOT/.claude"
CLAUDE_MD="$PROJECT_ROOT/CLAUDE.md"

# Colors
MAGENTA='\033[0;35m'
NC='\033[0m'

# Generate CLAUDE.md from .hatch/rules (concatenates all rules into one file)
generate_claude_md() {
    local rules_dir="$HATCH_DIR/src/rules"
    local commands_dir="$HATCH_DIR/src/commands"

    [[ ! -d "$rules_dir" ]] && return

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
    } > "$CLAUDE_MD"
    
    echo -e "${MAGENTA}claude${NC} CLAUDE.md"
}

# Generate Claude Code commands from .hatch/commands
generate_claude_commands() {
    local src_dir="$HATCH_DIR/src/commands"
    local dest_dir="$CLAUDE_DIR/commands"

    [[ ! -d "$src_dir" ]] && return

    mkdir -p "$dest_dir"

    find "$src_dir" -name "*.md" -type f | sort | while read -r src_file; do
        local rel_path="${src_file#$src_dir/}"
        local dest_file="$dest_dir/$rel_path"
        
        mkdir -p "$(dirname "$dest_file")"
        cp "$src_file" "$dest_file"
        
        echo -e "${MAGENTA}claude command${NC} ${rel_path}"
    done
}

# Copy ignore file for Claude Code
copy_claude_ignore() {
    local src_file="$HATCH_DIR/src/.ignore"
    local dest_file="$CLAUDE_DIR/.claudeignore"

    [[ ! -f "$src_file" ]] && return

    mkdir -p "$CLAUDE_DIR"
    cp "$src_file" "$dest_file"
    echo -e "${MAGENTA}claude${NC} .claude/.claudeignore"
}

generate_claude_md
generate_claude_commands
copy_claude_ignore

#!/bin/bash

# Validate Reports Structure and Naming Conventions
# Checks for reports in wrong locations, validates kebab-case naming, and verifies structure

set -e

REPORTS_DIR="docs/reports"
ERRORS=0
WARNINGS=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 Validating Base-App Reports Structure..."
echo ""

# Function to check if filename is kebab-case
is_kebab_case() {
    local filename="$1"
    # Check if filename matches kebab-case pattern (lowercase letters, numbers, hyphens)
    if [[ "$filename" =~ ^[a-z0-9]+(-[a-z0-9]+)*\.md$ ]]; then
        return 0
    else
        return 1
    fi
}

# Function to check for date in filename (should not be there)
has_date_in_filename() {
    local filename="$1"
    # Check for common date patterns
    if [[ "$filename" =~ [0-9]{4}-[0-9]{2}-[0-9]{2} ]] || \
       [[ "$filename" =~ [0-9]{4}_[0-9]{2}_[0-9]{2} ]] || \
       [[ "$filename" =~ [0-9]{8} ]]; then
        return 0
    else
        return 1
    fi
}


# Check directory structure
echo "📁 Checking directory structure..."

REQUIRED_DIRS=(
    "audits"
    "audits/security"
    "technical"
    "implementation"
    "implementation/stages"
    "implementation/milestones"
    "analysis"
    "services"
    "services/auth"
    "services/theme"
    "services/webhook"
    "planning"
)

for dir in "${REQUIRED_DIRS[@]}"; do
    if [ ! -d "$REPORTS_DIR/$dir" ]; then
        echo -e "${RED}❌ ERROR:${NC} Missing directory: $REPORTS_DIR/$dir"
        ((ERRORS++))
    fi
done

# Check README.md exists
if [ ! -f "$REPORTS_DIR/README.md" ]; then
    echo -e "${RED}❌ ERROR:${NC} Missing README.md in reports directory"
    ((ERRORS++))
else
    echo -e "${GREEN}✅${NC} README.md exists"
fi

echo ""

# Validate reports by category
echo "📄 Validating report files..."

# Function to validate file with correct expected directory
validate_file_location() {
    local file="$1"
    local relative_path="${file#$REPORTS_DIR/}"
    local dir=$(dirname "$relative_path")
    local filename=$(basename "$file")
    
    # Determine expected directory based on actual location
    local expected_dir="$dir"
    
    # Validate kebab-case naming
    if ! is_kebab_case "$filename"; then
        echo -e "${RED}❌ ERROR:${NC} $file - Filename should be kebab-case (lowercase with hyphens)"
        echo "   Example: authentication-security-audit.md"
        ((ERRORS++))
    fi
    
    # Check for date in filename
    if has_date_in_filename "$filename"; then
        echo -e "${YELLOW}⚠️  WARNING:${NC} $file - Filename contains date (dates should be in frontmatter)"
        ((WARNINGS++))
    fi
    
    # Check for required frontmatter fields
    if [ -f "$file" ]; then
        if ! grep -q "^\\*\\*Date:" "$file"; then
            echo -e "${YELLOW}⚠️  WARNING:${NC} $file - Missing Date field in frontmatter"
            ((WARNINGS++))
        fi
        if ! grep -q "^\\*\\*Status:" "$file"; then
            echo -e "${YELLOW}⚠️  WARNING:${NC} $file - Missing Status field in frontmatter"
            ((WARNINGS++))
        fi
        if ! grep -q "^\\*\\*Category:" "$file"; then
            echo -e "${YELLOW}⚠️  WARNING:${NC} $file - Missing Category field in frontmatter"
            ((WARNINGS++))
        fi
    fi
}

# Audits (including subdirectories)
if [ -d "$REPORTS_DIR/audits" ]; then
    find "$REPORTS_DIR/audits" -name "*.md" -type f | while read file; do
        validate_file_location "$file"
    done
fi

# Technical
if [ -d "$REPORTS_DIR/technical" ]; then
    find "$REPORTS_DIR/technical" -name "*.md" -type f | while read file; do
        validate_file_location "$file"
    done
fi

# Implementation (including subdirectories)
if [ -d "$REPORTS_DIR/implementation" ]; then
    find "$REPORTS_DIR/implementation" -name "*.md" -type f | while read file; do
        validate_file_location "$file"
    done
fi

# Analysis
if [ -d "$REPORTS_DIR/analysis" ]; then
    find "$REPORTS_DIR/analysis" -name "*.md" -type f | while read file; do
        validate_file_location "$file"
    done
fi

# Services (including subdirectories)
if [ -d "$REPORTS_DIR/services" ]; then
    find "$REPORTS_DIR/services" -name "*.md" -type f | while read file; do
        validate_file_location "$file"
    done
fi

# Planning
if [ -d "$REPORTS_DIR/planning" ]; then
    find "$REPORTS_DIR/planning" -name "*.md" -type f | while read file; do
        validate_file_location "$file"
    done
fi

# Check for broken links in README
echo ""
echo "🔗 Checking for broken links in README.md..."

if [ -f "$REPORTS_DIR/README.md" ]; then
    # Extract markdown links
    grep -o '\[.*\](\.\/.*\.md)' "$REPORTS_DIR/README.md" | sed 's/.*(\(.*\))/\1/' | while read link; do
        # Resolve relative path
        link_path="$REPORTS_DIR/${link#./}"
        if [ ! -f "$link_path" ]; then
            echo -e "${YELLOW}⚠️  WARNING:${NC} Broken link in README.md: $link"
            ((WARNINGS++))
        fi
    done
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Validation Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ All checks passed!${NC}"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ No errors found${NC}"
    echo -e "${YELLOW}⚠️  Warnings: $WARNINGS${NC}"
    exit 0
else
    echo -e "${RED}❌ Errors: $ERRORS${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  Warnings: $WARNINGS${NC}"
    fi
    exit 1
fi


#!/bin/bash

# Auto Collapse
#This is a script that takes in
# - a markdown file
# - a designated output file
# - the number of lines as a threshold on when to collapse a section.
#
# Sample
# ./auto_collapse.sh test_example.md test_example_collapsible.md 10

# Input and output files
INPUT_FILE=${1:-"README.md"}
OUTPUT_FILE=${2:-"README_collapsible.md"}
THRESHOLD=${3:-10}

# Ensure output file is empty initially
rm -f "$OUTPUT_FILE"

# Variables to track sections
inside_section=false
section_lines=()
section_header=""

# Function to process and write a section
write_section() {
    local header="$1"
    local lines=("${@:2}")
    # Remove ## or other Markdown header prefixes for the summary
    local clean_header
    clean_header=$(echo "$header" | sed -E 's/^#+[[:space:]]//')

    if [[ ${#lines[@]} -gt $THRESHOLD ]]; then
        {
          echo "<details>"
          echo "<summary>${clean_header}</summary>"
          echo
          printf "%s\n" "${lines[@]}"
          echo "</details>"
        } >> "$OUTPUT_FILE"
    else
        echo "$header" >> "$OUTPUT_FILE"
        printf "%s\n" "${lines[@]}" >> "$OUTPUT_FILE"
    fi
}

# Read the Markdown file line by line
while IFS= read -r line; do
    if [[ $line =~ ^#+\  ]]; then # New section starts
        if [[ $inside_section == true ]]; then
            # Write the previous section
            write_section "$section_header" "${section_lines[@]:1}" # Exclude the header from content
        fi
        # Start a new section
        section_header="$line"
        section_lines=("$line") # Initialize section with the header
        inside_section=true
    else
        # Collect lines of the current section
        section_lines+=("$line")
    fi
done < "$INPUT_FILE"

# Process the last section
if [[ $inside_section == true ]]; then
    write_section "$section_header" "${section_lines[@]:1}" # Exclude the header from content
fi

echo "Collapsible Markdown written to $OUTPUT_FILE"

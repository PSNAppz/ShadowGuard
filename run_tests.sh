#!/bin/bash

# Default values for flags
VERBOSE=false
COVERAGE=false
PACKAGE=""
RACE=false
TIMEOUT=""

# Help message
usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  -v            Enable verbose output for each test"
  echo "  -c            Generate a coverage report"
  echo "  -p <package>  Specify package to test (default is all packages)"
  echo "  -r            Enable race condition detection"
  echo "  -t <time>     Set a custom timeout (e.g., 2m, 1h)"
  echo "  -h            Display this help message"
}

# Parse command-line options
while getopts "vcp:rt:h" opt; do
  case ${opt} in
    v ) VERBOSE=true ;;
    c ) COVERAGE=true ;;
    p ) PACKAGE=$OPTARG ;;
    r ) RACE=true ;;
    t ) TIMEOUT=$OPTARG ;;
    h ) usage; exit 0 ;;
    * ) usage; exit 1 ;;
  esac
done

# Base command
CMD="go test"

# Add flags based on options
[ "$VERBOSE" = true ] && CMD+=" -v"
[ "$RACE" = true ] && CMD+=" -race"
[ -n "$TIMEOUT" ] && CMD+=" -timeout $TIMEOUT"
[ -n "$PACKAGE" ] && CMD+=" $PACKAGE" || CMD+=" ./..."

# Coverage option
if [ "$COVERAGE" = true ]; then
  COVERAGE_FILE="coverage.out"
  CMD+=" -coverprofile=$COVERAGE_FILE"
fi

# Run the command
echo "Running command: $CMD"
$CMD

# Show coverage report if generated
if [ "$COVERAGE" = true ]; then
  echo "Coverage report:"
  go tool cover -func=$COVERAGE_FILE

  # Optional: Display HTML report
  # go tool cover -html=$COVERAGE_FILE
fi

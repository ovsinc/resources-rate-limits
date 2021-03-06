#!/bin/bash
#
# Code coverage generation

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PKG_LIST=$(go list ./... | grep -v /vendor/ | grep -v /static/ | grep -v internal/examples | grep -v internal/build)

# Create the coverage files directory
mkdir -vp "$COVERAGE_DIR";

# Create a coverage file for each package
for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package" ;
done ;

# Merge the coverage profile files
echo 'mode: count' > coverage.cov ;
tail -q -n +2 "${COVERAGE_DIR}"/*.cov >> coverage.cov ;

# Display the global code coverage
go tool cover -func=coverage.cov ;

# If needed, generate HTML report
if [ "$1" == "html" ]; then
    out="coverage.html"
    [[ -n $2 ]] && out="$2"
    go tool cover -html=coverage.cov -o "$out" ;
fi

# Remove the coverage files directory
rm -rf "$COVERAGE_DIR";
rm -f coverage.cov;

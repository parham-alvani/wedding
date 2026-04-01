# Wedding RSVP Website

mod front './wedfront'
mod back './wedback'

# Install all dependencies
install:
    just front::install
    cd wedback && go mod download

# Build everything
build:
    just front::build
    just back::build

# Start both services for development (frontend + backend)
dev:
    #!/usr/bin/env bash
    set -euo pipefail
    trap 'kill 0' EXIT
    just back::serve &
    just front::dev &
    wait

# Start both services in production
serve:
    #!/usr/bin/env bash
    set -euo pipefail
    trap 'kill 0' EXIT
    just back::serve &
    just front::serve &
    wait

# Run all tests
test:
    just back::test

# Run all linters
lint:
    just back::lint
    just front::format

# Clean all build artifacts
clean:
    just front::clean
    just back::clean

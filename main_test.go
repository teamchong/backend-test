package main

import (
	"testing"

	// Import all standard Benthos components
	_ "github.com/benthosdev/benthos/v4/public/components/all"

	// Import all Benthos Plugin packages
	_ "github.com/teamchong/backend-test/ratelimit"
)

// This example demonstrates how to create a rate limit plugin, which is
// configured by providing a struct containing the fields to be parsed from
// within the Benthos configuration.
func TestMain(t *testing.T) {
	// And then execute Benthos with:
	main()
}

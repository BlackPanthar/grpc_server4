// This file contains a test function that tests the StateTracker() function.
//
// The TestStateTracker() function creates a context, initializes the gRPC connection to the server,
// calls the StateTracker() function, and logs the response. If there is an error in calling
// the StateTracker() function, the test function logs the error using t.Error and the test fails.
package state_tracker

import (
	"testing"
)

func TestStateTracker(t *testing.T) {
	StateTracker()
}

package log

import (
	"account-service/testutils"
	"testing"
)

func TestMain(m *testing.M) {
	testutils.VerifyGoLeaks(m)
}

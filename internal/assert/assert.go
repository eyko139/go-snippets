package assert

import "testing"

func AssertEqual[T comparable](t *testing.T, actual, expected T) {
    t.Helper()

    if actual != expected {
        t.Errorf("Expected %v, got %v", expected, actual)
    }
}

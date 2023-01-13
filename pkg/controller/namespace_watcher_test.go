package controller

import (
	"testing"

	"github.com/acorn-io/baaah/pkg/router/tester"
	"github.com/cloudnautique/kifthenfi/pkg/scheme"
)

func TestNamespaceWatcher(t *testing.T) {
	tester.DefaultTest(t, scheme.Scheme, "testdata/basic", namespaceWatcherHandler)
}
func TestNamespaceWatcherSetSelector(t *testing.T) {
	tester.DefaultTest(t, scheme.Scheme, "testdata/set-selector", namespaceWatcherHandler)
}

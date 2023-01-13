package controller

import (
	"github.com/acorn-io/baaah/pkg/router"
	kiftfv1 "github.com/cloudnautique/kifthenfi/pkg/apis/kifthenfi.cloudnautique.com/v1"
)

func namespaceWatcherHandler(req router.Request, resp router.Response) error {
	nsw := req.Object.(*kiftfv1.NamespaceWatcher)
	return ApplyManifests(req, nsw, resp)
}

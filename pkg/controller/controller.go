package controller

import (
	"context"
	"time"

	"github.com/acorn-io/acorn/pkg/config"
	"github.com/acorn-io/acorn/pkg/k8sclient"
	"github.com/acorn-io/baaah"
	"github.com/acorn-io/baaah/pkg/apply"
	"github.com/acorn-io/baaah/pkg/restconfig"
	"github.com/acorn-io/baaah/pkg/router"
	kiftfv1 "github.com/cloudnautique/kifthenfi/pkg/apis/kifthenfi.cloudnautique.com/v1"
	"github.com/cloudnautique/kifthenfi/pkg/crds"
	"github.com/cloudnautique/kifthenfi/pkg/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Controller struct {
	Router *router.Router
	client client.Client
	Scheme *runtime.Scheme
	apply  apply.Apply
}

func New() (*Controller, error) {
	appScheme := scheme.Scheme
	router, err := baaah.DefaultRouter("kifthenfi", appScheme)
	if err != nil {
		return nil, err
	}

	cfg, err := restconfig.New(appScheme)
	if err != nil {
		return nil, err
	}

	client, err := k8sclient.New(cfg)
	if err != nil {
		return nil, err
	}

	apply := apply.New(client)

	router.Type(&kiftfv1.NamespaceWatcher{}).HandlerFunc(namespaceWatcherHandler)
	return &Controller{
		Router: router,
		client: client,
		Scheme: appScheme,
		apply:  apply,
	}, nil
}

func (c *Controller) Start(ctx context.Context) error {
	if err := crds.Create(ctx, c.Scheme, kiftfv1.SchemeGroupVersion); err != nil {
		return err
	}

	go func() {
		var success bool
		for i := 0; i < 6000; i++ {
			// This will error until the cache is primed
			if _, err := config.Get(ctx, c.Router.Backend()); err == nil {
				success = true
				break
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
		if !success {
			panic("couldn't initial cached client")
		}
	}()

	return c.Router.Start(ctx)
}

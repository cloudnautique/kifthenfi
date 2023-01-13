package main

import (
	"flag"
	"fmt"

	"github.com/cloudnautique/kifthenfi/pkg/controller"
	"github.com/cloudnautique/kifthenfi/pkg/version"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	versionFlag = flag.Bool("version", false, "print version")
)

func main() {

	flag.Parse()

	fmt.Printf("Version: %s", version.Get())
	if *versionFlag {
		return
	}

	ctx := signals.SetupSignalHandler()
	logrus.SetLevel(logrus.InfoLevel)

	ctrlr, err := controller.New()
	if err != nil {
		logrus.Fatal(err)
	}

	if err := ctrlr.Start(ctx); err != nil {
		logrus.Fatal(err)
	}
	<-ctx.Done()
	logrus.Fatal(ctx.Err())
}

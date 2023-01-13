package v1

import (
	kiftfv1 "github.com/cloudnautique/kifthenfi/pkg/apis/kifthenfi.cloudnautique.com"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const Version = "v1"

var SchemeGroupVersion = schema.GroupVersion{
	Group:   kiftfv1.Group,
	Version: Version,
}

func AddToScheme(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&NamespaceWatcher{},
		&NamespaceWatcherList{},
	)

	// Add common types
	scheme.AddKnownTypes(SchemeGroupVersion, &metav1.Status{})

	// Add the watch version that applies
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

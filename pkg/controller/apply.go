package controller

import (
	"github.com/acorn-io/baaah/pkg/router"
	kiftfv1 "github.com/cloudnautique/kifthenfi/pkg/apis/kifthenfi.cloudnautique.com/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func ApplyManifests(req router.Request, nsw *kiftfv1.NamespaceWatcher, resp router.Response) error {
	secret, err := nsw.GetManifestSecretAndSetStatus(req.Ctx, req.Client)
	if apierror.IsNotFound(err) {
		logrus.Infof("Waiting for secret %s", nsw.Spec.ManifestSecretName)
		return nil
	}
	if err != nil {
		return err
	}

	sel, err := metav1.LabelSelectorAsSelector(&nsw.Spec.NamespaceLabelSelector)
	if err != nil {
		return err
	}

	nsl := corev1.NamespaceList{}
	err = req.List(&nsl, &kclient.ListOptions{
		LabelSelector: sel,
	})
	if err != nil {
		return err
	}

	for _, ns := range nsl.Items {
		for key, manifest := range secret.Data {
			logrus.Infof("Applying section %s to namespace %s", key, ns.Name)

			labelsToApply := nsw.GetChildObjectLabels(key)
			obj, _ := prepObject(ns.Name, labelsToApply, manifest)

			resp.Objects(obj)
		}
	}
	return nil
}

func prepObject(namespace string, labels map[string]string, manifest []byte) (*unstructured.Unstructured, error) {
	ustr := &unstructured.Unstructured{}
	serializer := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	_, gvk, err := serializer.Decode(manifest, nil, ustr)
	if err != nil {
		return nil, err
	}

	ustr.SetNamespace(namespace)
	ustr.SetGroupVersionKind(*gvk)
	ustr.SetLabels(labels)

	return ustr, nil
}

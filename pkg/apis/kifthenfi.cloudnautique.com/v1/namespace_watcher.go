package v1

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/acorn-io/baaah/pkg/router"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const NamespaceWatcherLabelPrefix = "kifthenfi.cloudnautique.com"

func (nsw *NamespaceWatcher) GetManifestSecretAndSetStatus(ctx context.Context, client kclient.Client) (secret *corev1.Secret, err error) {
	secret = &corev1.Secret{}
	err = client.Get(ctx, router.Key(nsw.Namespace, nsw.Spec.ManifestSecretName), secret)
	if apierrors.IsNotFound(err) {
		nsw.Status.SecretStatus.Found = false
		return
	} else if err != nil {
		return
	}
	nsw.Status.SecretStatus.Found = true

	nsw.Status.SecretStatus.KeyHashes = hashManifestMap(secret.Data)

	return
}

func (nsw *NamespaceWatcher) GetChildObjectLabels(targetManifest string) map[string]string {
	labels := map[string]string{}
	for k, v := range nsw.ObjectMeta.Labels {
		labels[k] = v
	}
	labels[NamespaceWatcherLabelPrefix+"/manifest-key"] = targetManifest

	if nsw.Status.SecretStatus.KeyHashes != nil {
		if val, ok := nsw.Status.SecretStatus.KeyHashes[targetManifest]; ok {
			labels[NamespaceWatcherLabelPrefix+"/manifest-hash"] = val[:8]
		}
	}

	return labels
}

func hashManifestMap(data map[string][]byte) map[string]string {
	rData := map[string]string{}
	for key, manifest := range data {
		shaData := sha256.Sum256(manifest)
		rData[key] = hex.EncodeToString(shaData[:])
	}
	return rData
}

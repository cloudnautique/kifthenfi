// +k8s:deepcopy-gen=package

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NamespaceWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   NamespaceWatcherSpec   `json:"spec,omitempty"`
	Status NamespaceWatcherStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NamespaceWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceWatcher `json:"items"`
}

type NamespaceWatcherSpec struct {
	NamespaceLabelSelector metav1.LabelSelector `json:"namespaceLabelSelector,omitempty"`
	ManifestSecretName     string               `json:"manifestSecretName,omitempty"`
}

type NamespaceWatcherStatus struct {
	SecretStatus SecretStatus `json:"secretStatus,omitempty"`
}

type SecretStatus struct {
	Found     bool              `json:"found,omitempty"`
	KeyHashes map[string]string `json:"keyHashes,omitempty"`
}

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ESAlias is alias object from https://www.elastic.co/guide/en/elasticsearch/reference/7.x/indices-aliases.html
type ESAlias struct {
	// (Array) Array of index names used to perform the action. If the index parameter is not specified, this parameter is required.
	// +required
	// +listType=set
	Indices []string `json:"indices,omitempty"`
	// (String) Comma-separated list or wildcard expression of index alias names to add, remove, or delete. If the alias parameter is not specified, this parameter is required for the add or remove action.
	// +required
	// +listType=set
	Aliases []string `json:"aliases,omitempty"`
	// (Optional, query object) Filter query used to limit the index alias. If specified, the index alias only applies to documents returned by the filter. Filter query used to limit the index alias. If specified, the index alias only applies to documents returned by the filter. See Filtered aliases for an example.
	// +optional
	Filter string `json:"filter,omitempty"`
	// (Optional, boolean) If true, assigns the index as an alias’s write index. Defaults to false.
	// +optional
	IsWriteIndex bool `json:"is_write_index,omitempty"`
	// (Optional, string) Custom routing value used to route operations to a specific shard.
	// +optional
	Routing string `json:"routing,omitempty"`
	// (Optional, string) Custom routing value used for the alias’s indexing operations.
	// +optional
	IndexRouting string `json:"index_routing,omitempty"`
	// (Optional, string) Custom routing value used for the alias’s search operations.
	// +optional
	SearchRouting string `json:"search_routing,omitempty"`
}

// ESAliases is aliases object from https://www.elastic.co/guide/en/elasticsearch/reference/7.x/indices-aliases.html
type ESAliases struct {
	// Adds an alias to an index.
	// +optional
	Add ESAlias `json:"add,omitempty"`
	// Adds an alias to an index.
	// +optional
	Remove ESAlias `json:"remove,omitempty"`
	// Adds an alias to an index.
	// +optional
	RemoveIndex ESAlias `json:"remove_index,omitempty"`
}

// ElasticSearchTemplateSpec defines the desired state of ElasticSearchTemplate
type ElasticSearchTemplateSpec struct {
	// See more at https://www.elastic.co/guide/en/elasticsearch/reference/7.x/indices-templates.html
	// Name of ES template
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Pattern=`^[^-_+A-Z][^A-Z\\\/\*\?"\<\> ,|#]{1,254}$`
	Name string `json:"name"`
	// Should we drop template if K8S object is deleted, default false
	// +optional
	DropOnDelete bool `json:"drop_on_delete,omitempty"`

	// (Required, array of strings) Array of wildcard expressions used to match the names of indices during creation.
	// +listType=set
	IndexPatterns []string `json:"index_patterns"`
	// (Optional, alias object) Index aliases which include the index. See Update index alias.
	// +optional
	Aliases ESAliases `json:"aliases,omitempty"`

	// (Optional, index setting object) Configuration options for the index. See Index Settings.
	// +optional
	Settings ESIndexSettings `json:"settings"`
	// (Optional, mapping object) Mapping for fields in the index. If specified, this mapping can include:
	//     Field names
	//     Field datatypes
	//     Mapping parameters
	// +kubebuilder:validation:Pattern=`[^,:{}\[\]0-9.\-+Eaeflnr-u \n\r\t]`
	Mappings string `json:"mappings"`

	// (Optional, integer) Version number used to manage index templates externally. This number is not automatically generated by Elasticsearch.
	// +optional
	// +kubebuilder:validation:MinValue=1
	Version int64 `json:"version,omitempty"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// ElasticSearchTemplateStatus defines the observed state of ElasticSearchTemplate
type ElasticSearchTemplateStatus struct {
	Succeeded bool `json:"succeeded"`
	// +optional
	Name string `json:"template_name"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticSearchTemplate is the Schema for the elasticsearchtemplates API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=elasticsearchtemplates,scope=Namespaced
type ElasticSearchTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticSearchTemplateSpec   `json:"spec,omitempty"`
	Status ElasticSearchTemplateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticSearchTemplateList contains a list of ElasticSearchTemplate
type ElasticSearchTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticSearchTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticSearchTemplate{}, &ElasticSearchTemplateList{})
}

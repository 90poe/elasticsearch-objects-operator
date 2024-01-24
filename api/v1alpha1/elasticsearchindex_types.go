/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ES Index values taken from: https://www.elastic.co/guide/en/elasticsearch/reference/7.x/index-modules.html
type ESIdle struct {
	// How long a shard can not receive a search or get request until it’s considered search idle. (default is 30s)
	// +optional
	After string `json:"after,omitempty"`
}

type ESSearch struct {
	// +optional
	Idle ESIdle `json:"idle,omitempty"`
}

type ESAnalyze struct {
	// The maximum number of tokens that can be produced using _analyze API. Defaults to 10000.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxTokenCount int64 `json:"max_token_count,omitempty"`
}

type ESHighlights struct {
	// The maximum number of characters that will be analyzed for a highlight request. This setting is only applicable when highlighting is requested on a text that was indexed without offsets or term vectors. Defaults to 1000000.
	// +optional
	// +kubebuilder:validation:Minimum=1
	HighlightMaxAnalyzedOffset int64 `json:"max_analyzed_offset,omitempty"`
}

// ESIndexSettings settings for index
// NOTE: until we run on 1.17 - we can't have default values, this means we can't use bool here. We would need to use string and watch for "" to note default value
type ESIndexSettings struct {
	// Static Index settings
	// The number of primary shards that an index should have. Defaults to 1. This setting can only be set at index creation time. It cannot be changed on a closed index. Note: the number of shards are limited to 1024 per index. This limitation is a safety limit to prevent accidental creation of indices that can destabilize a cluster due to resource allocation.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1024
	NumOfShards int32 `json:"number_of_shards,omitempty"`
	// Shard structure
	// +optional
	Shard ESShard `json:"shard,omitempty"`
	// The default value compresses stored data with LZ4 compression, but this can be set to best_compression which uses DEFLATE for a higher compression ratio, at the expense of slower stored fields performance. If you are updating the compression type, the new one will be applied after segments are merged. Segment merging can be forced using force merge.
	// +optional
	// +kubebuilder:validation:Pattern=`^(default|best_compression)$`
	Codec string `json:"codec,omitempty"`
	// The number of shards a custom routing value can go to. Defaults to 1 and can only be set at index creation time. This value must be less than the index.number_of_shards unless the index.number_of_shards value is also 1. See Routing to an index partition for more details about how this setting is used.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1024
	// +kubebuilder:validation:ExclusiveMaximum=true
	RoutingPartitionSize int32 `json:"routing_partition_size,omitempty"`
	// Indicates whether cached filters are pre-loaded for nested queries. Possible values are true (default) and false.
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false)$`
	LoadFixedBitsetFiltersEagerly string `json:"load_fixed_bitset_filters_eagerly,omitempty"`
	// Indicates whether the index should be hidden by default. Hidden indices are not returned by default when using a wildcard expression. This behavior is controlled per request through the use of the expand_wildcards parameter. Possible values are true and false (default).
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false)$`
	Hidden string `json:"hidden,omitempty"`

	// Dynamic index settings
	// The number of replicas each primary shard has. Defaults to 1.
	// +optional
	// +kubebuilder:validation:Minimum=1
	NumOfReplicas int32 `json:"number_of_replicas,omitempty"`
	// Auto-expand the number of replicas based on the number of data nodes in the cluster. Set to a dash delimited lower and upper bound (e.g. 0-5) or use all for the upper bound (e.g. 0-all). Defaults to false (i.e. disabled). Note that the auto-expanded number of replicas only takes allocation filtering rules into account, but ignores any other allocation rules such as shard allocation awareness and total shards per node, and this can lead to the cluster health becoming YELLOW if the applicable rules prevent all the replicas from being allocated.
	// +optional
	AutoExpandReplicas string `json:"auto_expand_replicas,omitempty"`
	// How long a shard can not receive a search or get request until it’s considered search idle. (default is 30s)
	// +optional
	SearchIdleAfter ESSearch `json:"search,omitempty"`
	// How often to perform a refresh operation, which makes recent changes to the index visible to search. Defaults to 1s. Can be set to -1 to disable refresh. If this setting is not explicitly set, shards that haven’t seen search traffic for at least index.search.idle.after seconds will not receive background refreshes until they receive a search request. Searches that hit an idle shard where a refresh is pending will wait for the next background refresh (within 1s). This behavior aims to automatically optimize bulk indexing in the default case when no searches are performed. In order to opt out of this behavior an explicit value of 1s should set as the refresh interval.
	// +optional
	RefreshInterval string `json:"refresh_interval,omitempty"`
	// The maximum value of from + size for searches to this index. Defaults to 10000. Search requests take heap memory and time proportional to from + size and this limits that memory. See Scroll or Search After for a more efficient alternative to raising this.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxResultWindow int64 `json:"max_result_window,omitempty"`
	// The maximum value of from + size for inner hits definition and top hits aggregations to this index. Defaults to 100. Inner hits and top hits aggregation take heap memory and time proportional to from + size and this limits that memory.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxInnerResultWindow int64 `json:"max_inner_result_window,omitempty"`
	// The maximum value of window_size for rescore requests in searches of this index. Defaults to index.max_result_window which defaults to 10000. Search requests take heap memory and time proportional to max(window_size, from + size) and this limits that memory.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxRescoreWindow int64 `json:"max_rescore_window,omitempty"`
	// The maximum value of window_size for rescore requests in searches of this index. Defaults to index.max_result_window which defaults to 10000. Search requests take heap memory and time proportional to max(window_size, from + size) and this limits that memory.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxDocValueFieldsSearch int64 `json:"max_docvalue_fields_search,omitempty"`
	// The maximum number of script_fields that are allowed in a query. Defaults to 32.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxScriptFields int64 `json:"max_script_fields,omitempty"`
	// The maximum allowed difference between min_gram and max_gram for NGramTokenizer and NGramTokenFilter. Defaults to 1.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxNgramDiff int64 `json:"max_ngram_diff,omitempty"`
	// The maximum allowed difference between max_shingle_size and min_shingle_size for ShingleTokenFilter. Defaults to 3.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxShingleDiff int64 `json:"max_shingle_diff,omitempty"`

	// +optional
	Blocks ESIndexBlocks `json:"blocks,omitempty"`

	// Maximum number of refresh listeners available on each shard of the index. These listeners are used to implement refresh=wait_for.
	// The maximum allowed difference between max_shingle_size and min_shingle_size for ShingleTokenFilter. Defaults to 3.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxRefreshListeners int64 `json:"max_refresh_listeners,omitempty"`

	// +optional
	Analyze ESAnalyze `json:"analyze,omitempty"`

	// +optional
	Highlight ESHighlights `json:"highlight,omitempty"`
	// The maximum number of terms that can be used in Terms Query. Defaults to 65536.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxTermsCount int64 `json:"max_terms_count,omitempty"`
	// The maximum length of regex that can be used in Regexp Query. Defaults to 1000.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxRegexLength int64 `json:"max_regex_length,omitempty"`
	// Routing values
	// +optional
	Routing ESIndexRouting `json:"routing,omitempty"`
	// The length of time that a deleted document’s version number remains available for further versioned operations. Defaults to 60s.
	// +optional
	GCdeletes string `json:"gc_deletes,omitempty"`
	// The default ingest node pipeline for this index. Index requests will fail if the default pipeline is set and the pipeline does not exist. The default may be overridden using the pipeline parameter. The special pipeline name _none indicates no ingest pipeline should be run.
	// +optional
	DefaultPipeline string `json:"default_pipeline,omitempty"`
	// The final ingest node pipeline for this index. Index requests will fail if the final pipeline is set and the pipeline does not exist. The final pipeline always runs after the request pipeline (if specified) and the default pipeline (if it exists). The special pipeline name _none indicates no ingest pipeline will run.
	// +optional
	FinalPipeline string `json:"final_pipeline,omitempty"`
}

// ESShard would hold shard structure
type ESShard struct {
	// Whether or not shards should be checked for corruption before opening. When corruption is detected, it will prevent the shard from being opened. Accepts: true, false, checksum
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false|checksum)$`
	CheckOnStartup string `json:"check_on_startup,omitempty"`
}

// ESIndexBlocks defines block in dynamic values
type ESIndexBlocks struct {
	// Set to true to make the index and index metadata read only, false to allow writes and metadata changes.
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false)$`
	ReadOnly string `json:"read_only,omitempty"`
	// Similar to index.blocks.read_only but also allows deleting the index to free up resources. The disk-based shard allocator may add and remove this block automatically.
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false)$`
	ReadOnlyAllowDelete string `json:"read_only_allow_delete,omitempty"`
	// Set to true to disable read operations against the index.
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false)$`
	Read string `json:"read,omitempty"`
	// Set to true to disable data write operations against the index. Unlike read_only, this setting does not affect metadata. For instance, you can close an index with a write block, but not an index with a read_only block.
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false)$`
	Write string `json:"write,omitempty"`
	// Set to true to disable index metadata reads and writes.
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false)$`
	Metadata string `json:"metadata,omitempty"`
}

type ESRoutingAllocationEnable struct {
	// Controls shard allocation for this index. It can be set to:
	//     all (default) - Allows shard allocation for all shards.
	//     primaries - Allows shard allocation only for primary shards.
	//     new_primaries - Allows shard allocation only for newly-created primary shards.
	//     none - No shard allocation is allowed.
	// +optional
	// +kubebuilder:validation:Pattern=`^(all|primaries|new_primaries|none)$`
	Enable string `json:"enable,omitempty"`
}

type ESRoutingRebalanceEnable struct {
	// Enables shard rebalancing for this index. It can be set to:
	//     all (default) - Allows shard rebalancing for all shards.
	//     primaries - Allows shard rebalancing only for primary shards.
	//     replicas - Allows shard rebalancing only for replica shards.
	//     none - No shard rebalancing is allowed.
	// +optional
	// +kubebuilder:validation:Pattern=`^(all|primaries|replicas|none)$`
	Enable string `json:"enable,omitempty"`
}

// ESIndexRouting defines routing in dynamic values
type ESIndexRouting struct {
	Allocation ESRoutingAllocationEnable `json:"allocation,omitempty"`
	Rebalance  ESRoutingRebalanceEnable  `json:"rebalance,omitempty"`
}

// ElasticSearchIndexSpec defines the desired state of ElasticSearchIndex
// +k8s:openapi-gen=true
type ElasticSearchIndexSpec struct {
	// See more at https://www.elastic.co/guide/en/elasticsearch/reference/7.x/index-modules.html
	// Name of ES index
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Pattern=`^[^-_+A-Z][^A-Z\\\/\*\?"\<\> ,|#]{1,254}$`
	Name string `json:"name"`
	// Should we drop index if K8S object is deleted, default false
	// +optional
	DropOnDelete bool `json:"drop_on_delete,omitempty"`

	// Index settings
	// +optional
	Settings ESIndexSettings `json:"settings"`
	// Mappings of ES Index
	// +kubebuilder:validation:Pattern=`[{\[]{1}([,:{}\[\]0-9.\-+Eaeflnr-u \n\r\t]|".*?")+[}\]]{1}`
	Mappings string `json:"mappings"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "export GOROOT=/usr/local/go; operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// ElasticSearchIndexStatus defines the observed state of ElasticSearchIndex
// +k8s:openapi-gen=true
type ElasticSearchIndexStatus struct {
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ElasticSearchIndex is the Schema for the elasticsearchindices API
type ElasticSearchIndex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticSearchIndexSpec   `json:"spec,omitempty"`
	Status ElasticSearchIndexStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ElasticSearchIndexList contains a list of ElasticSearchIndex
type ElasticSearchIndexList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticSearchIndex `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticSearchIndex{}, &ElasticSearchIndexList{})
}

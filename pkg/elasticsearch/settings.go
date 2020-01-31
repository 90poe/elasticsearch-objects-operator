package elasticsearch

// import xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"

// import "reflect"

// //Settings represents generic ES index settings
// //nolint
// type Settings struct {
// 	//Static settings
// 	NumOfShards                   int32  `json:"index.number_of_shards,omitempty"`
// 	ShardCheckOnStartup           string `json:"index.shard.check_on_startup,omitempty"`
// 	Codec                         string `json:"index.codec,omitempty"`
// 	RoutingPartitionSize          int32  `json:"index.routing_partition_size,omitempty"`
// 	LoadFixedBitsetFiltersEagerly string `json:"index.load_fixed_bitset_filters_eagerly,omitempty"`
// 	Hidden                        string `json:"index.hidden,omitempty"`
// 	//Dynamic settings
// 	NumOfReplicas              int32  `json:"index.number_of_replicas,omitempty"`
// 	AutoExpandReplicas         string `json:"index.auto_expand_replicas,omitempty"`
// 	SearchIdleAfter            string `json:"index.search.idle.after,omitempty"`
// 	RefreshInterval            string `json:"index.refresh_interval,omitempty"`
// 	MaxResultWindow            int64  `json:"index.max_result_window,omitempty"`
// 	MaxInnerResultWindow       int64  `json:"index.max_inner_result_window,omitempty"`
// 	MaxRescoreWindow           int64  `json:"index.max_rescore_window,omitempty"`
// 	MaxDocValueFieldsSearch    int64  `json:"index.max_docvalue_fields_search,omitempty"`
// 	MaxScriptFields            int64  `json:"index.max_script_fields,omitempty"`
// 	MaxNgramDiff               int64  `json:"index.max_ngram_diff,omitempty"`
// 	MaxShingleDiff             int64  `json:"index.max_shingle_diff,omitempty"`
// 	ReadOnly                   string `json:"index.blocks.read_only,omitempty"`
// 	ReadOnlyAllowDelete        string `json:"index.blocks.read_only_allow_delete,omitempty"`
// 	Read                       string `json:"index.blocks.read,omitempty"`
// 	Write                      string `json:"index.blocks.write,omitempty"`
// 	Metadata                   string `json:"index.blocks.metadata,omitempty"`
// 	MaxRefreshListeners        int64  `json:"index.max_refresh_listeners,omitempty"`
// 	AnalyzeMaxTokenCount       int64  `json:"index.analyze.max_token_count,omitempty"`
// 	HighlightMaxAnalyzedOffset int64  `json:"index.highlight.max_analyzed_offset,omitempty"`
// 	MaxTermsCount              int64  `json:"index.max_terms_count,omitempty"`
// 	MaxRegexLength             int64  `json:"index.max_regex_length,omitempty"`
// 	AllocationEnable           string `json:"index.routing.allocation.enable,omitempty"`
// 	RebalanceEnable            string `json:"index.routing.rebalance.enable,omitempty"`
// 	GCdeletes                  string `json:"index.gc_deletes,omitempty"`
// 	DefaultPipeline            string `json:"index.default_pipeline,omitempty"`
// 	FinalPipeline              string `json:"index.final_pipeline,omitempty"`
// }

// // InitFromMap is going to initialize Settings from provided map
// func (s *Settings) InitFromMap(settings map[string]interface{}) {
// 	if val, ok := getInt64ValueFromSettings(settings, "index.number_of_shards"); ok {
// 		s.NumOfShards = int32(val)
// 	}
// 	if val, ok := getStringValueFromSettings(settings, "index.shards.check_on_startup"); ok {
// 		s.ShardCheckOnStartup = val
// 	}
// 	if val, ok := getStringValueFromSettings(settings, "index.codec"); ok {
// 		s.Codec = val
// 	}
// 	if val, ok := getInt64ValueFromSettings(settings, "index.routing_partition_size"); ok {
// 		s.RoutingPartitionSize = int32(val)
// 	}
// 	if val, ok := getStringValueFromSettings(settings, "index.load_fixed_bitset_filters_eagerly"); ok {
// 		s.LoadFixedBitsetFiltersEagerly = val
// 	}
// 	if val, ok := getStringValueFromSettings(settings, "index.hidden"); ok {
// 		s.Hidden = val
// 	}
// }

// // InitFromK8SObj would init Settings from K8S index
// func (s *Settings) InitFromK8SObj(k8sIndex *xov1alpha1.ElasticSearchIndex) {
// 	s.NumOfShards = k8sIndex.Spec.Static.NumOfShards
// 	s.ShardCheckOnStartup = k8sIndex.Spec.Static.ShardCheckOnStartup
// 	s.Codec = k8sIndex.Spec.Static.Codec
// 	s.RoutingPartitionSize = k8sIndex.Spec.Static.RoutingPartitionSize
// 	s.LoadFixedBitsetFiltersEagerly
// }

// func (s *Settings) getK8SStatic() *xov1alpha1.ESIndexStaticSpec {
// 	ret := &xov1alpha1.ESIndexStaticSpec{
// 		NumOfShards:                   s.NumOfShards,
// 		ShardCheckOnStartup:           s.ShardCheckOnStartup,
// 		Codec:                         s.Codec,
// 		RoutingPartitionSize:          s.RoutingPartitionSize,
// 		LoadFixedBitsetFiltersEagerly: s.LoadFixedBitsetFiltersEagerly,
// 		Hidden:                        s.Hidden,
// 	}
// 	return ret
// }

// //DiffStatic would compare static settings to ones in k8sIndex
// func (s *Settings) DiffStatic(k8sIndex *xov1alpha1.ElasticSearchIndex) bool {
// 	static := s.getK8SStatic()
// 	if !reflect.DeepEqual(static, k8sIndex.Spec.Static) {
// 		return true
// 	}
// 	return false
// }

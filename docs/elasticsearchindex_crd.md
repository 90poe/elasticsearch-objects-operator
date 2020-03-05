# ElasticSearch Index CRD

Example:
```
apiVersion: xo.90poe.io/v1alpha1
kind: ElasticSearchIndex
metadata:
  name: example-elasticsearchindex
  namespace: "90"
spec:
  name: dev_test_test
  drop_on_delete: true
  settings:
    number_of_shards: 55
    shards:
      check_on_startup: "false"
    codec: "default"
    number_of_replicas: 3
  mappings: |
    {
      "dynamic": false,
      "_source": {
        "enabled": true
      },
      "properties": {
        "isRead": {
          "type": "boolean",
          "index": true
        },
        "createdAt": {
          "type": "date",
          "index": true
        }
      }
    }
```

## Spec

Setting are made to be as close as possible to [ES API](https://www.elastic.co/guide/en/elasticsearch/reference/7.x/index-modules.html).

You would need to amend `spec` section.

|Spec|Type |Required|Notes|
|--------|:---:|:------:|:---|
|name|string|Yes|Name of ES index|
|drop_on_delete|bool|No|Should we drop index if K8S object is deleted, default false|
|settings|ESIndexSettings|Yes|See <a href="#ESIndexSettings">ESIndexSettings</a>|
|mappings|string|Yes|Mappings of ES Index, must be valid JSON|


## ESIndexSettings
<a name="ESIndexSettings"></a>

|Settings|Type |Required|Notes|
|--------|:---:|:------:|:---|
|number_of_shards|int32|No|The number of primary shards that an index should have. Defaults to 1. This setting can only be set at index creation time. It cannot be changed on a closed index. Note: the number of shards are limited to 1024 per index. This limitation is a safety limit to prevent accidental creation of indices that can destabilize a cluster due to resource allocation.|
|shard.check_on_startup|string|No|Whether or not shards should be checked for corruption before opening. When corruption is detected, it will prevent the shard from being opened. Accepts: true, false, checksum|
|codec|string|No|The default value compresses stored data with LZ4 compression, but this can be set to best_compression which uses DEFLATE for a higher compression ratio, at the expense of slower stored fields performance. If you are updating the compression type, the new one will be applied after segments are merged. Segment merging can be forced using force merge.|
|routing_partition_size|int32|No|The number of shards a custom routing value can go to. Defaults to 1 and can only be set at index creation time. This value must be less than the index.number_of_shards unless the index.number_of_shards value is also 1. See Routing to an index partition for more details about how this setting is used.|
|load_fixed_bitset_filters_eagerly|string|No|Indicates whether cached filters are pre-loaded for nested queries. Possible values are true (default) and false.|
|hidden|string|No|Indicates whether the index should be hidden by default. Hidden indices are not returned by default when using a wildcard expression. This behavior is controlled per request through the use of the expand_wildcards parameter. Possible values are true and false (default).|
|number_of_replicas|int32|No|The number of replicas each primary shard has. Defaults to 1.|
|auto_expand_replicas|string|No|Auto-expand the number of replicas based on the number of data nodes in the cluster. Set to a dash delimited lower and upper bound (e.g. 0-5) or use all for the upper bound (e.g. 0-all). Defaults to false (i.e. disabled). Note that the auto-expanded number of replicas only takes allocation filtering rules into account, but ignores any other allocation rules such as shard allocation awareness and total shards per node, and this can lead to the cluster health becoming YELLOW if the applicable rules prevent all the replicas from being allocated.|
|search.idle.after|string|No|How long a shard can not receive a search or get request until it’s considered search idle. (default is 30s)|
|refresh_interval|string|No|How often to perform a refresh operation, which makes recent changes to the index visible to search. Defaults to 1s. Can be set to -1 to disable refresh. If this setting is not explicitly set, shards that haven’t seen search traffic for at least index.search.idle.after seconds will not receive background refreshes until they receive a search request. Searches that hit an idle shard where a refresh is pending will wait for the next background refresh (within 1s). This behavior aims to automatically optimize bulk indexing in the default case when no searches are performed. In order to opt out of this behavior an explicit value of 1s should set as the refresh interval.|
|max_result_window|int64|No|The maximum value of from + size for searches to this index. Defaults to 10000. Search requests take heap memory and time proportional to from + size and this limits that memory. See Scroll or Search After for a more efficient alternative to raising this.|
|max_inner_result_window|int64|No|The maximum value of from + size for inner hits definition and top hits aggregations to this index. Defaults to 100. Inner hits and top hits aggregation take heap memory and time proportional to from + size and this limits that memory.|
|max_rescore_window|int64|No|The maximum value of window_size for rescore requests in searches of this index. Defaults to index.max_result_window which defaults to 10000. Search requests take heap memory and time proportional to max(window_size, from + size) and this limits that memory.|
|max_docvalue_fields_search|int64|No|The maximum value of window_size for rescore requests in searches of this index. Defaults to index.max_result_window which defaults to 10000. Search requests take heap memory and time proportional to max(window_size, from + size) and this limits that memory.|
|max_script_fields|int64|No|The maximum number of script_fields that are allowed in a query. Defaults to 32.|
|max_ngram_diff|int64|No|The maximum allowed difference between min_gram and max_gram for NGramTokenizer and NGramTokenFilter. Defaults to 1.|
|max_shingle_diff|int64|No|The maximum allowed difference between max_shingle_size and min_shingle_size for ShingleTokenFilter. Defaults to 3.|
|blocks.read_only|string|No|Set to true to make the index and index metadata read only, false to allow writes and metadata changes.|
|blocks.read_only_allow_delete|string|No|Similar to index.blocks.read_only but also allows deleting the index to free up resources. The disk-based shard allocator may add and remove this block automatically.|
|blocks.read|string|No|Set to true to disable read operations against the index.|
|blocks.write|string|No|Set to true to disable data write operations against the index. Unlike read_only, this setting does not affect metadata. For instance, you can close an index with a write block, but not an index with a read_only block.|
|blocks.metadata|string|No|Set to true to disable index metadata reads and writes.|
|max_refresh_listeners|int64|No|Maximum number of refresh listeners available on each shard of the index. These listeners are used to implement refresh=wait_for. The maximum allowed difference between max_shingle_size and min_shingle_size for ShingleTokenFilter. Defaults to 3.|
|analyze.max_token_count|int64|No|The maximum number of tokens that can be produced using _analyze API. Defaults to 10000.|
|highlight.max_analyzed_offset|int64|No|The maximum number of characters that will be analyzed for a highlight request. This setting is only applicable when highlighting is requested on a text that was indexed without offsets or term vectors. Defaults to 1000000.|
|max_terms_count|int64|No|The maximum number of terms that can be used in Terms Query. Defaults to 65536.|
|max_regex_length|int64|No|The maximum length of regex that can be used in Regexp Query. Defaults to 1000.|
|routing.allocation.enable|string|No|Controls shard allocation for this index. It can be set to: all (default) - Allows shard allocation for all shards. primaries - Allows shard allocation only for primary shards. new_primaries - Allows shard allocation only for newly-created primary shards. none - No shard allocation is allowed.|
|routing.rebalance.enable|string|No|// Enables shard rebalancing for this index. It can be set to: all (default) - Allows shard rebalancing for all shards. primaries - Allows shard rebalancing only for primary shards. replicas - Allows shard rebalancing only for replica shards. none - No shard rebalancing is allowed.|
|gc_deletes|string|No|The length of time that a deleted document’s version number remains available for further versioned operations. Defaults to 60s.|
|default_pipeline|string|No|The default ingest node pipeline for this index. Index requests will fail if the default pipeline is set and the pipeline does not exist. The default may be overridden using the pipeline parameter. The special pipeline name _none indicates no ingest pipeline should be run.|
|final_pipeline|string|No|The final ingest node pipeline for this index. Index requests will fail if the final pipeline is set and the pipeline does not exist. The final pipeline always runs after the request pipeline (if specified) and the default pipeline (if it exists). The special pipeline name _none indicates no ingest pipeline will run.|

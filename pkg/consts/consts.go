package consts

const (
	// ESIndexNameRegexp according to https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html
	ESIndexNameRegexp = `^[^-_+A-Z][^A-Z\\\/\*\?"\<\> ,|#]{1,254}$`
)

// ESStaticSettings is map which has ES settings static part
var ESStaticSettings map[string]bool

func init() {
	// init ESStaticSettings
	ESStaticSettings = map[string]bool{
		"index.number_of_shards":                  true,
		"index.shard.check_on_startup":            true,
		"index.codec":                             true,
		"index.routing_partition_size":            true,
		"index.load_fixed_bitset_filters_eagerly": true,
		"index.hidden":                            true,
	}
}

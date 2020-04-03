package elasticsearch

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
)

type TestValueSettings struct {
	settings  string
	valuePath string
	value     interface{}
	ok        bool
}

type TestKeySettings struct {
	settings string
	keys     []string
}

type TestAddManagedBy2InterfaceSettings struct {
	before string
	after  string
	err    error
}

type TestIsManagedByESOperatorSettings struct {
	mappings  string
	expecting bool
}

type TestDiffSett struct {
	src   *xov1alpha1.ESIndexSettings
	dest  string
	index bool
	res   bool
	err   error
}

func TestGetInt64ValueFromSettings(t *testing.T) {
	tests := []TestValueSettings{
		{
			settings:  `{"index": {"number_of_shards": "55"}}`,
			valuePath: "index.number_of_shards",
			value:     55,
			ok:        true,
		},
		{
			settings:  `{"index": {"number_of_shards": "55", "highlight":{"max_analyzed_offset":"10000"}}}`,
			valuePath: "index.highlight.max_analyzed_offset",
			value:     10000,
			ok:        true,
		},
		{
			settings:  `{"index": {"number_of_shards": "55fff"}}`,
			valuePath: "index.number_of_shards",
			value:     0,
			ok:        false,
		},
		{
			settings:  `{"indexBad": 55}`,
			valuePath: "index.number_of_shards",
			value:     0,
			ok:        false,
		},
		{
			settings:  `{"index": {"number_of_shards": {"some":true}}}`,
			valuePath: "index.number_of_shards",
			value:     0,
			ok:        false,
		},
		{
			settings:  `{"index": {"number_of_shards": {"some":true}}}`,
			valuePath: "index.number_of_shards.some.get",
			value:     0,
			ok:        false,
		},
	}
	for _, test := range tests {
		var sett map[string]interface{}
		err := json.Unmarshal([]byte(test.settings), &sett)
		if err != nil {
			t.Fatalf("could not unmarshal '%s'", test.settings)
		}
		ret, ok := getInt64ValueFromSettings(sett, test.valuePath)
		if ok != test.ok {
			t.Fatalf("could not fetch setting for %s", test.valuePath)
		}
		if !test.ok {
			continue
		}
		testVal, ok := test.value.(int)
		if !ok {
			t.Fatalf("can't convert test.value to int64 as it's of kinf %v",
				reflect.TypeOf(test.value).Kind())
		}
		if ret != int64(testVal) {
			t.Fatalf("values don't match %d != %d", ret, test.value.(int32))
		}
	}
}

func TestGetStringValueFromSettings(t *testing.T) {
	tests := []TestValueSettings{
		{
			settings:  `{"index": {"number_of_shards": "55"}}`,
			valuePath: "index.number_of_shards",
			value:     "55",
			ok:        true,
		},
		{
			settings:  `{"index": {"number_of_shards": "55", "highlight":{"max_analyzed_offset":10000}}}`,
			valuePath: "index.highlight.max_analyzed_offset",
			value:     "10000",
			ok:        true,
		},
		{
			settings:  `{"indexBad": 55}`,
			valuePath: "index.number_of_shards",
			value:     "",
			ok:        false,
		},
		{
			settings:  `{"index": {"number_of_shards": {"some":true}}}`,
			valuePath: "index.number_of_shards",
			value:     `map[string]interface {}{"some":true}`,
			ok:        true,
		},
		{
			settings:  `{"index": {"number_of_shards": {"some":true}}}`,
			valuePath: "index.number_of_shards.some.get",
			value:     "",
			ok:        false,
		},
	}
	for _, test := range tests {
		var sett map[string]interface{}
		err := json.Unmarshal([]byte(test.settings), &sett)
		if err != nil {
			t.Fatalf("could not unmarshal '%s'", test.settings)
		}
		ret, ok := getStringValueFromSettings(sett, test.valuePath)
		if ok != test.ok {
			t.Fatalf("could not fetch setting for %s", test.valuePath)
		}
		if !test.ok {
			continue
		}
		testVal, ok := test.value.(string)
		if ok != test.ok {
			t.Fatalf("expected %v but got %v on %s",
				test.ok, ok, test.valuePath)
		}
		if ret != testVal {
			t.Fatalf("values don't match %s != %s", ret, test.value.(string))
		}
	}
}

func TestGetKeysFromSettings(t *testing.T) {
	tests := []TestKeySettings{
		{
			settings: `{"index": {"number_of_shards": "55"}}`,
			keys:     []string{"index.number_of_shards"},
		},
		{
			settings: `{"index": {"number_of_shards": ["55"], "shard":{"num_of_tests":66}}}`,
			keys:     []string{"index.number_of_shards", "index.shard.num_of_tests"},
		},
		{
			settings: `{"index": {"arr": ["55"], "map":{"int":66,"bool":true, "str": "string"}}}`,
			keys:     []string{"index.arr", "index.map.int", "index.map.bool", "index.map.str"},
		},
		{
			settings: "{}",
			keys:     []string{},
		},
	}
	for _, test := range tests {
		var sett map[string]interface{}
		err := json.Unmarshal([]byte(test.settings), &sett)
		if err != nil {
			t.Fatalf("could not unmarshal '%s'", test.settings)
		}
		retKeys := getKeysFromSettings("", sett)
		sort.Slice(retKeys, func(i, j int) bool {
			return retKeys[i] > retKeys[j]
		})
		sort.Slice(test.keys, func(i, j int) bool {
			return test.keys[i] > test.keys[j]
		})
		if !reflect.DeepEqual(retKeys, test.keys) {
			t.Fatalf("keys are fetched incorrect. Expected '%v', got '%v'", test.keys, retKeys)
		}
	}
}

func TestAddManagedBy2Interface(t *testing.T) {
	tests := []TestAddManagedBy2InterfaceSettings{
		{
			before: `{"index": {"number_of_shards": "55"}}`,
			after:  `{"index": {"number_of_shards": "55"}, "_meta": {"managed-by":"elasticsearch-objects-operator.xo.90poe.io"}}`,
		},
		{
			before: `{"index": {"number_of_shards": "55"}, "_meta": {"managed-by":"xo.90poe.io"}}`,
			after:  `{"index": {"number_of_shards": "55"}, "_meta": {"managed-by":"elasticsearch-objects-operator.xo.90poe.io"}}`,
		},
		{
			before: "{[]]}",
			after:  "",
			err:    fmt.Errorf("can't json unmarshal mappings: invalid character '[' looking for beginning of object key string"),
		},
		{
			before: "4623456234",
			after:  "",
			err:    fmt.Errorf("expected map"),
		},
		{
			before: `{"index": {"number_of_shards": "55"}, "_meta": []}`,
			after:  "",
			err:    fmt.Errorf("invalid _meta map"),
		},
	}
	for _, test := range tests {
		after, err := addManagedBy2Interface(test.before)
		if test.err != nil {
			if fmt.Sprintf("%v", test.err) != fmt.Sprintf("%v", err) {
				t.Fatalf("After test got incorrect err. Expected '%v', got '%v'", test.err, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("could not add managed-by '%s'", test.after)
		}
		var afterTest interface{}
		err = json.Unmarshal([]byte(test.after), &afterTest)
		if err != nil {
			t.Fatalf("could not unmarshal '%s'", test.after)
		}
		if !reflect.DeepEqual(after, afterTest) {
			t.Fatalf("After test was added incorrectly. Expected '%#v', got '%#v'", afterTest, after)
		}
	}
}

func TestIsManagedByESOperator(t *testing.T) {
	tests := []TestIsManagedByESOperatorSettings{
		{
			mappings:  `{"_meta": {"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"properties": {"country": {"type": "text","index": false}}}`,
			expecting: true,
		},
		{
			mappings:  `{"_meta": {"managed-by":"xo"},"properties": {"country": {"type": "text","index": false}}}`,
			expecting: false,
		},
		{
			mappings:  `{"properties": {"country": {"type": "text","index": false}}}`,
			expecting: false,
		},
		{
			mappings:  `{"_meta": [],"properties": {"country": {"type": "text","index": false}}}`,
			expecting: false,
		},
		{
			mappings:  `{"_meta": {"is": true},"properties": {"country": {"type": "text","index": false}}}`,
			expecting: false,
		},
		{
			mappings:  `{"_meta": {"managed-by": true},"properties": {"country": {"type": "text","index": false}}}`,
			expecting: false,
		},
	}
	for _, test := range tests {
		var mappings map[string]interface{}
		err := json.Unmarshal([]byte(test.mappings), &mappings)
		if err != nil {
			t.Fatalf("can't json unmarshal mappings: %v", err)
		}
		got := isManagedByESOperator(mappings)
		if got != test.expecting {
			t.Fatalf("expected '%v' got '%v'", test.expecting, got)
		}
	}
}

func TestDiffSettings(t *testing.T) {
	tests := []TestDiffSett{
		{
			src:  &xov1alpha1.ESIndexSettings{},
			dest: "{}",
		},
		{
			src: &xov1alpha1.ESIndexSettings{
				NumOfShards: 5,
			},
			dest: `{"index": {"number_of_shards": 4}}`,
			res:  true,
		},
		{
			src: &xov1alpha1.ESIndexSettings{
				NumOfShards: 5,
			},
			dest:  `{"index": {"number_of_shards": 4}}`,
			index: true,
			err:   fmt.Errorf("can't change static setting index.number_of_shards from '4' to '5'"),
		},
		{
			src: &xov1alpha1.ESIndexSettings{
				NumOfReplicas: 5,
			},
			dest:  `{"index": {"number_of_replicas": 4}}`,
			index: true,
			res:   true,
		},
	}
	for _, test := range tests {
		var destTest map[string]interface{}
		err := json.Unmarshal([]byte(test.dest), &destTest)
		if err != nil {
			t.Fatalf("could not unmarshal '%s'", test.dest)
		}
		res, err := diffSettings(test.src, destTest, test.index)
		if test.err != nil {
			if fmt.Sprintf("%v", test.err) != fmt.Sprintf("%v", err) {
				t.Fatalf("After test got incorrect err. Expected '%v', got '%v'", test.err, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("could not diff '%v'", err)
		}
		if res != test.res {
			t.Fatalf("After test got incorrect result. Expected '%v', got '%v'", test.res, res)
		}
	}
}

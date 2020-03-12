package elasticsearch

import (
	"fmt"
	"sort"
	"testing"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
	"github.com/stretchr/testify/assert"
)

type TestCreateTempl struct {
	Templ *xov1alpha1.ElasticSearchTemplate
	R2R   Responce2Req
	Err   error
}

type TestUpdateTempl struct {
	Templ *xov1alpha1.ElasticSearchTemplate
	R2R   map[int]Responce2Req
	Err   error
	Msg   string
}

func TestCreateTemplate(t *testing.T) {
	var testDoer *TestDoer
	client := Client{}
	client.esURL = "http://localhost:80"
	client.es, testDoer = setupCreateTestClient(t)
	defer testDoer.Close()
	tests := []TestCreateTempl{
		{
			Templ: &xov1alpha1.ElasticSearchTemplate{
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          "some_templ",
					DropOnDelete:  true,
					IndexPatterns: []string{"some_index"},
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 32,
					},
					Aliases: map[string]xov1alpha1.ESAlias{
						"{index}-alias-for-{gender}": {},
					},
					Mappings: `
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
					`,
				},
			},
			R2R: Responce2Req{
				RequestURI:   "/_template/some_templ",
				ResponceCode: 200,
				Responce:     `{"acknowledged":true}`,
			},
		},
		{
			Templ: &xov1alpha1.ElasticSearchTemplate{
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          "/_template/some_templ",
					DropOnDelete:  true,
					IndexPatterns: []string{"some_index"},
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 32,
					},
					Mappings: `
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
					`,
				},
			},
			R2R: Responce2Req{
				RequestURI:   "/_template/%2F_template%2Fsome_templ",
				ResponceCode: 200,
				Responce:     `{"acknowledged":false}`,
			},
			Err: fmt.Errorf("can't acknowledge ES template creation/update"),
		},
	}
	for _, test := range tests {
		testDoer.R2rChan <- test.R2R
		err := client.CreateTemplate(test.Templ)
		if test.Err != nil {
			assert.EqualError(t, err, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.NoError(t, err)
	}
}

func TestUpdateTemplate(t *testing.T) {
	var testDoer *TestDoer
	client := Client{}
	client.esURL = "http://localhost:80"
	client.es, testDoer = setupCreateTestClient(t)
	defer testDoer.Close()
	tests := []TestUpdateTempl{
		{
			Templ: &xov1alpha1.ElasticSearchTemplate{
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          "some_templ",
					DropOnDelete:  true,
					IndexPatterns: []string{"some_index"},
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 32,
					},
					Mappings: `
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
					`,
				},
			},
			R2R: map[int]Responce2Req{
				1: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 200,
					Responce:     `{"some_templ":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}}}}`,
				},
			},
			Msg: "no changes on template named some_templ",
		},
		{
			Templ: &xov1alpha1.ElasticSearchTemplate{
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          "some_templ",
					DropOnDelete:  true,
					IndexPatterns: []string{"some_index"},
					Aliases: map[string]xov1alpha1.ESAlias{
						"{index}-alias-for-{gender}": {
							Filter: `{"term":{"product":"Elasticsearch"}}`,
						},
					},
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33,
					},
					Mappings: `
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
					`,
				},
			},
			R2R: map[int]Responce2Req{
				1: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 200,
					Responce:     `{"some_templ":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}}}}`,
				},
				2: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
			Msg: "successfully updated ES template some_templ",
		},
		//Test to see how template is re-created on name change
		{
			Templ: &xov1alpha1.ElasticSearchTemplate{
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          "some_templ",
					DropOnDelete:  true,
					IndexPatterns: []string{"some_index"},
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33,
					},
					Mappings: `
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
					`,
				},
			},
			R2R: map[int]Responce2Req{
				1: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 404,
					Responce:     `{}`,
				},
				2: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
			Msg: "successfully created updated ES template some_templ",
		},
		{
			Templ: &xov1alpha1.ElasticSearchTemplate{
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          "some_templ",
					DropOnDelete:  true,
					IndexPatterns: []string{"some_index"},
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33,
					},
					Mappings: `
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
					`,
				},
			},
			R2R: map[int]Responce2Req{
				1: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 200,
					Responce:     `{"some_templ":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`,
				},
				2: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 200,
					Responce:     `{"acknowledged":false}`,
				},
			},
			Err: fmt.Errorf("can't acknowledge ES template creation/update"),
		},
	}
	for _, test := range tests {
		r2rKeys := make([]int, 0, len(test.R2R))
		for key := range test.R2R {
			r2rKeys = append(r2rKeys, key)
		}
		sort.Ints(r2rKeys)
		for _, value := range r2rKeys {
			testDoer.R2rChan <- test.R2R[value]
		}
		msg, err := client.UpdateTemplate(test.Templ)
		if test.Err != nil {
			assert.EqualError(t, err, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, test.Msg, msg)
	}
}

func TestDeleteTemplate(t *testing.T) {
	var testDoer *TestDoer
	client := Client{}
	client.esURL = "http://localhost:80"
	client.es, testDoer = setupCreateTestClient(t)
	defer testDoer.Close()
	tests := []TestDelete{
		{
			IndexName: "some_templ",
			R2R: Responce2Req{
				RequestURI:   "/_template/some_templ",
				ResponceCode: 200,
				Responce:     `{"acknowledged":true}`,
			},
		},
		{
			IndexName: "some_templ",
			R2R: Responce2Req{
				RequestURI:   "/_template/some_templ",
				ResponceCode: 404,
				Responce:     `{"error":{"root_cause":[{"type":"index_template_missing_exception","reason":"index_template [some_templ] missing"}],"type":"index_template_missing_exception","reason":"index_template [some_templ] missing"},"status":404}`,
			},
			Err: fmt.Errorf("can't delete template %s: elastic: Error 404 (Not Found): index_template [some_templ] missing [type=index_template_missing_exception]", "some_templ"),
		},
		{
			IndexName: "some_templ",
			R2R: Responce2Req{
				RequestURI:   "/_template/some_templ",
				ResponceCode: 200,
				Responce:     `{"acknowledged":false}`,
			},
			Err: fmt.Errorf("can't acknowledge ES template deletion"),
		},
	}
	for _, test := range tests {
		testDoer.R2rChan <- test.R2R
		err := client.DeleteTemplate(test.IndexName)
		if test.Err != nil {
			assert.EqualError(t, err, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.NoError(t, err)
	}
}

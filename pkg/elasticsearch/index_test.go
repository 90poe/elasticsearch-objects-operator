package elasticsearch

import (
	"fmt"
	"sort"
	"testing"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
	"github.com/stretchr/testify/assert"
)

type TestCreateIndx struct {
	Index *xov1alpha1.ElasticSearchIndex
	R2R   Responce2Req
	Err   error
}

type TestUpdateIndx struct {
	Index *xov1alpha1.ElasticSearchIndex
	R2R   map[int]Responce2Req
	Err   error
	Msg   string
}

type TestDelete struct {
	IndexName string
	R2R       Responce2Req
	Err       error
}

func TestCreateIndex(t *testing.T) {
	var testDoer *TestDoer
	client := Client{}
	client.esURL = "http://localhost:80"
	client.es, testDoer = setupCreateTestClient(t)
	defer testDoer.Close()
	tests := []TestCreateIndx{
		{
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
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
				RequestURI:   "/some_index",
				ResponceCode: 200,
				Responce:     `{"acknowledged":true,"shards_acknowledged":true,"index":"some_test"}`,
			},
		},
		{
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
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
				RequestURI:   "/some_index",
				ResponceCode: 400,
				Responce:     `{"error":{"root_cause":[{"type":"resource_already_exists_exception","reason":"index [some_test/UZXBpYOrTUeWydYMpphiEg] already exists","index_uuid":"UZXBpYOrTUeWydYMpphiEg","index":"some_test"}],"type":"resource_already_exists_exception","reason":"index [some_test/UZXBpYOrTUeWydYMpphiEg] already exists","index_uuid":"UZXBpYOrTUeWydYMpphiEg","index":"some_test"},"status":400}`,
			},
			Err: fmt.Errorf("can't create ES index: elastic: Error 400 (Bad Request): index [some_test/UZXBpYOrTUeWydYMpphiEg] already exists [type=resource_already_exists_exception]"),
		},
		{
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
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
				RequestURI:   "/some_index",
				ResponceCode: 200,
				Responce:     `{"acknowledged":false,"shards_acknowledged":false,"index":"some_test"}`,
			},
			Err: fmt.Errorf("can't acknowledge ES index creation"),
		},
	}
	for _, test := range tests {
		testDoer.R2rChan <- test.R2R
		err := client.CreateIndex(test.Index)
		if test.Err != nil {
			assert.EqualError(t, err, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.NoError(t, err)
	}
}

func TestUpdateIndex(t *testing.T) {
	var testDoer *TestDoer
	client := Client{}
	client.esURL = "http://localhost:80"
	client.es, testDoer = setupCreateTestClient(t)
	defer testDoer.Close()
	tests := []TestUpdateIndx{
		{
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
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
					RequestURI:   "/_cat/indices/some_index?format=json&pretty=false",
					ResponceCode: 200,
					Responce:     `[{"health":"yellow","status":"open","index":"some_index","uuid":"iQXnF_YMTKqminns7h0-Zw","pri":"32","rep":"1","docs.count":"0","docs.deleted":"0","store.size":"7.1kb","pri.store.size":"7.1kb"}]`,
				},
				2: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"some_index":{"settings":{"index":{"creation_date":"1581606515721","number_of_shards":"32","number_of_replicas":"1","uuid":"iQXnF_YMTKqminns7h0-Zw","version":{"created":"7050299"},"provided_name":"some_index"}}}}`,
				},
			},
			Msg: "no changes on index named some_index",
		},
		{
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
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
					RequestURI:   "/_cat/indices/some_index?format=json&pretty=false",
					ResponceCode: 200,
					Responce:     `[{"health":"yellow","status":"open","index":"some_index","uuid":"iQXnF_YMTKqminns7h0-Zw","pri":"32","rep":"1","docs.count":"0","docs.deleted":"0","store.size":"7.1kb","pri.store.size":"7.1kb"}]`,
				},
				2: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"some_index":{"settings":{"index":{"creation_date":"1581606515721","number_of_shards":"32","number_of_replicas":"1","uuid":"iQXnF_YMTKqminns7h0-Zw","version":{"created":"7050299"},"provided_name":"some_index"}}}}`,
				},
			},
			Err: fmt.Errorf("can't change static setting index.number_of_shards from '32' to '33'"),
		},
		{
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfReplicas: 6,
						NumOfShards:   32,
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
					RequestURI:   "/_cat/indices/some_index?format=json&pretty=false",
					ResponceCode: 200,
					Responce:     `[{"health":"yellow","status":"open","index":"some_index","uuid":"iQXnF_YMTKqminns7h0-Zw","pri":"32","rep":"1","docs.count":"0","docs.deleted":"0","store.size":"7.1kb","pri.store.size":"7.1kb"}]`,
				},
				2: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"some_index":{"settings":{"index":{"creation_date":"1581606515721","number_of_shards":"32","number_of_replicas":"1","uuid":"iQXnF_YMTKqminns7h0-Zw","version":{"created":"7050299"},"provided_name":"some_index"}}}}`,
				},
				3: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
			Msg: "successfully updated ES index some_index",
		},
		{
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfReplicas: 6,
						NumOfShards:   32,
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
					RequestURI:   "/_cat/indices/some_index?format=json&pretty=false",
					ResponceCode: 200,
					Responce:     `[{"health":"yellow","status":"open","index":"some_index","uuid":"iQXnF_YMTKqminns7h0-Zw","pri":"32","rep":"1","docs.count":"0","docs.deleted":"0","store.size":"7.1kb","pri.store.size":"7.1kb"}]`,
				},
				2: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"some_index":{"settings":{"index":{"creation_date":"1581606515721","number_of_shards":"32","number_of_replicas":"1","uuid":"iQXnF_YMTKqminns7h0-Zw","version":{"created":"7050299"},"provided_name":"some_index"}}}}`,
				},
				3: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"acknowledged":false}`,
				},
			},
			Err: fmt.Errorf("can't acknowledge ES index update"),
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
		msg, err := client.UpdateIndex(test.Index)
		if test.Err != nil {
			assert.EqualError(t, err, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, test.Msg, msg)
	}
}

func TestDeleteIndex(t *testing.T) {
	var testDoer *TestDoer
	client := Client{}
	client.esURL = "http://localhost:80"
	client.es, testDoer = setupCreateTestClient(t)
	defer testDoer.Close()
	tests := []TestDelete{
		{
			IndexName: "some_index",
			R2R: Responce2Req{
				RequestURI:   "/some_index",
				ResponceCode: 200,
				Responce:     `{"acknowledged":true}`,
			},
		},
		{
			IndexName: "some_index",
			R2R: Responce2Req{
				RequestURI:   "/some_index",
				ResponceCode: 404,
				Responce:     `{"error":{"root_cause":[{"type":"index_not_found_exception","reason":"no such index [some_test]","resource.type":"index_or_alias","resource.id":"some_test","index_uuid":"_na_","index":"some_test"}],"type":"index_not_found_exception","reason":"no such index [some_test]","resource.type":"index_or_alias","resource.id":"some_test","index_uuid":"_na_","index":"some_test"},"status":404}`,
			},
			Err: fmt.Errorf("can't delete index %s: elastic: Error 404 (Not Found): no such index [some_test] [type=index_not_found_exception]", "some_index"),
		},
		{
			IndexName: "some_index",
			R2R: Responce2Req{
				RequestURI:   "/some_index",
				ResponceCode: 200,
				Responce:     `{"acknowledged":false}`,
			},
			Err: fmt.Errorf("can't acknowledge ES index deletion"),
		},
	}
	for _, test := range tests {
		testDoer.R2rChan <- test.R2R
		err := client.DeleteIndex(test.IndexName)
		if test.Err != nil {
			assert.EqualError(t, err, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.NoError(t, err)
	}
}

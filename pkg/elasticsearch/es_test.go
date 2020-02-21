package elasticsearch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"

	xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

const (
	// Used by Doer to mock ES clusters answer on "/_nodes/http" request
	ClusterNodesGETanswer = `{"_nodes":{"total":1,"successful":1,"failed":0},"cluster_name":"docker-cluster","nodes":{"9PuzCdHhT6CdQMVEQeNYsg":{"name":"8d9407c2fa07","transport_address":"172.18.0.2:9300","host":"172.18.0.2","ip":"172.18.0.2","version":"7.5.2","build_flavor":"default","build_type":"docker","build_hash":"8bec50e1e0ad29dad5653712cf3bb580cd1afcdf","roles":["ingest","master","data","ml"],"attributes":{"ml.machine_memory":"2086154240","xpack.installed":"true","ml.max_open_jobs":"20"},"http":{"bound_address":["0.0.0.0:9200"],"publish_address":"172.18.0.2:9200","max_content_length_in_bytes":104857600}}}}`
)

type Responce2Req struct {
	RequestURI   string
	ResponceCode int
	Responce     string
}

type TestCreateIndx struct {
	Index *xov1alpha1.ElasticSearchIndex
	R2R   Responce2Req
	Err   error
}

type TestCreateTempl struct {
	Templ *xov1alpha1.ElasticSearchTemplate
	R2R   Responce2Req
	Err   error
}

type TestUpdateIndx struct {
	Index *xov1alpha1.ElasticSearchIndex
	R2R   map[int]Responce2Req
	Err   error
	Msg   string
}

type TestUpdateTempl struct {
	Templ *xov1alpha1.ElasticSearchTemplate
	R2R   map[int]Responce2Req
	Err   error
	Msg   string
}

type TestDelete struct {
	IndexName string
	R2R       Responce2Req
	Err       error
}

type TestDoer struct {
	R2rChan chan Responce2Req
}

func NewTestDoer(requestNums int) *TestDoer {
	testDoer := &TestDoer{}
	testDoer.R2rChan = make(chan Responce2Req, requestNums)
	return testDoer
}

func (t *TestDoer) Close() {
	close(t.R2rChan)
}

func (t *TestDoer) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{}
	resp.Proto = req.Proto
	resp.Header = make(http.Header, 2)
	resp.Header.Add("content-type", "application/json")
	resp.Header.Add("charset", "UTF-8")
	var err error
	switch req.Method {
	case "HEAD":
		resp.StatusCode = http.StatusOK
	case "GET", "PUT", "DELETE":
		err = t.httpCall(req, resp)
	}
	return resp, err
}

func (t *TestDoer) httpCall(req *http.Request, resp *http.Response) error {
	r2r := <-t.R2rChan
	reqURI := req.URL.RequestURI()
	if reqURI != r2r.RequestURI {
		resp.StatusCode = http.StatusNotFound
		return nil
	}
	resp.StatusCode = r2r.ResponceCode
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(r2r.Responce))
	return nil
}

func setupCreateTestClient(t *testing.T) (*elastic.Client, *TestDoer) {
	var err error
	testCreateDoer := NewTestDoer(5)
	r2r := Responce2Req{
		RequestURI:   "/_nodes/http",
		ResponceCode: http.StatusOK,
		Responce:     ClusterNodesGETanswer,
	}
	testCreateDoer.R2rChan <- r2r
	client, err := elastic.NewClient(elastic.SetHttpClient(testCreateDoer))
	if err != nil {
		t.Fatal(err)
	}

	return client, testCreateDoer
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
					Responce:     `{"some_templ":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`,
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
					Responce:     `{"some_templ":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`,
				},
				2: {
					RequestURI:   "/_template/some_templ",
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
			Msg: "successfully updated ES template some_templ",
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
					Responce:     `{"some_templ":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`,
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

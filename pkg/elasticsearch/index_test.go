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

func TestCreateUpdateIndex(t *testing.T) {
	var testDoer *TestDoer
	client := Client{}
	client.esURL = "http://localhost:80"
	client.es, testDoer = setupCreateTestClient(t)
	defer testDoer.Close()
	tests := []TestUpdateIndx{
		{
			//No changes
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 32,
					},
					Mappings: `
					{
						"_meta": {
							"managed-by": "elasticsearch-objects-operator.xo.90poe.io"
						},
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
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce: `
					{
						"some_index": {
						  "aliases": {},
						  "mappings": {
							"_meta": {
							  "managed-by": "elasticsearch-objects-operator.xo.90poe.io"
							},
							"properties": {
							  "country": {
								"type": "text",
								"index": false
							  },
							  "id": {
								"type": "keyword"
							  },
							  "portCode": {
								"type": "keyword"
							  },
							  "portName": {
								"type": "text"
							  },
							  "region": {
								"type": "text",
								"index": false
							  }
							}
						  },
						  "settings": {
							"index": {
							  "creation_date": "1581606515721",
							  "number_of_shards": "32",
							  "number_of_replicas": "1",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "some_index"
							}
						  }
						}
					  }
					`,
				},
			},
			Msg: "no changes on index named some_index",
		},
		{
			//Static settings change error
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33,
					},
					Mappings: `
					{
						"_meta": {
							"managed-by": "elasticsearch-objects-operator.xo.90poe.io"
						},
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
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce: `
					{
						"some_index": {
						  "aliases": {},
						  "mappings": {
							"_meta": {
							  "managed-by": "elasticsearch-objects-operator.xo.90poe.io"
							},
							"properties": {
							  "country": {
								"type": "text",
								"index": false
							  },
							  "id": {
								"type": "keyword"
							  },
							  "portCode": {
								"type": "keyword"
							  },
							  "portName": {
								"type": "text"
							  },
							  "region": {
								"type": "text",
								"index": false
							  }
							}
						  },
						  "settings": {
							"index": {
							  "creation_date": "1581606515721",
							  "number_of_shards": "32",
							  "number_of_replicas": "1",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "some_index"
							}
						  }
						}
					  }
					`,
				},
			},
			Err: fmt.Errorf("can't change static setting index.number_of_shards from '32' to '33'"),
		},
		{
			//Not managed by us
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33,
					},
					Mappings: `
					{
						"_meta": {
							"managed-by": "elasticsearch-objects-operator.xo.90poe.io"
						},
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
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce: `
					{
						"some_index": {
						  "aliases": {},
						  "mappings": {
							"properties": {
							  "country": {
								"type": "text",
								"index": false
							  },
							  "id": {
								"type": "keyword"
							  },
							  "portCode": {
								"type": "keyword"
							  },
							  "portName": {
								"type": "text"
							  },
							  "region": {
								"type": "text",
								"index": false
							  }
							}
						  },
						  "settings": {
							"index": {
							  "creation_date": "1581606515721",
							  "number_of_shards": "32",
							  "number_of_replicas": "1",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "some_index"
							}
						  }
						}
					  }
					`,
				},
			},
			Err: fmt.Errorf("index 'some_index' is not managed by this operator"),
		},
		{
			//Not managed by us 2
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33,
					},
					Mappings: `
					{
						"_meta": {
							"managed-by": "elasticsearch-objects-operator.xo.90poe.io"
						},
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
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce: `
					{
						"some_index": {
						  "aliases": {},
						  "settings": {
							"index": {
							  "creation_date": "1581606515721",
							  "number_of_shards": "32",
							  "number_of_replicas": "1",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "some_index"
							}
						  }
						}
					  }
					`,
				},
			},
			Err: fmt.Errorf("index 'some_index' is not managed by this operator"),
		},
		{
			//Index can't be fetched 421 error
			Index: &xov1alpha1.ElasticSearchIndex{
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         "some_index",
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33,
					},
					Mappings: `
					{
						"_meta": {
							"managed-by": "elasticsearch-objects-operator.xo.90poe.io"
						},
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
					RequestURI:   "/some_index",
					ResponceCode: 421,
					Responce:     `{"error":"fiction_one"}`,
				},
			},
			Err: fmt.Errorf("can't get index details: can't get settings and mappings: elastic: Error 421 (Misdirected Request)"),
		},
		{
			//Successfull Update
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
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce: `
					{
						"some_index": {
						  "aliases": {},
						  "mappings": {
							"_meta": {
							  "managed-by": "elasticsearch-objects-operator.xo.90poe.io"
							},
							"properties": {
							  "country": {
								"type": "text",
								"index": false
							  },
							  "id": {
								"type": "keyword"
							  },
							  "portCode": {
								"type": "keyword"
							  },
							  "portName": {
								"type": "text"
							  },
							  "region": {
								"type": "text",
								"index": false
							  }
							}
						  },
						  "settings": {
							"index": {
							  "creation_date": "1581606515721",
							  "number_of_shards": "32",
							  "number_of_replicas": "1",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "some_index"
							}
						  }
						}
					  }
					`,
				},
				2: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
			Msg: "successfully updated ES index some_index",
		},
		{
			//Un-aknowledged Update
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
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce: `
					{
						"some_index": {
						  "aliases": {},
						  "mappings": {
							"_meta": {
							  "managed-by": "elasticsearch-objects-operator.xo.90poe.io"
							},
							"properties": {
							  "country": {
								"type": "text",
								"index": false
							  },
							  "id": {
								"type": "keyword"
							  },
							  "portCode": {
								"type": "keyword"
							  },
							  "portName": {
								"type": "text"
							  },
							  "region": {
								"type": "text",
								"index": false
							  }
							}
						  },
						  "settings": {
							"index": {
							  "creation_date": "1581606515721",
							  "number_of_shards": "32",
							  "number_of_replicas": "1",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "some_index"
							}
						  }
						}
					  }
					`,
				},
				2: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 200,
					Responce:     `{"acknowledged":false}`,
				},
			},
			Err: fmt.Errorf("can't acknowledge ES index update"),
		},
		{
			//Un-successfull Update
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
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce: `
					{
						"some_index": {
						  "aliases": {},
						  "mappings": {
							"_meta": {
							  "managed-by": "elasticsearch-objects-operator.xo.90poe.io"
							},
							"properties": {
							  "country": {
								"type": "text",
								"index": false
							  },
							  "id": {
								"type": "keyword"
							  },
							  "portCode": {
								"type": "keyword"
							  },
							  "portName": {
								"type": "text"
							  },
							  "region": {
								"type": "text",
								"index": false
							  }
							}
						  },
						  "settings": {
							"index": {
							  "creation_date": "1581606515721",
							  "number_of_shards": "32",
							  "number_of_replicas": "1",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "some_index"
							}
						  }
						}
					  }
					`,
				},
				2: {
					RequestURI:   "/some_index/_settings",
					ResponceCode: 500,
					Responce:     `{}`,
				},
			},
			Err: fmt.Errorf("can't update ES index: elastic: Error 500 (Internal Server Error)"),
		},
		{
			//Successfull index creation
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
					RequestURI:   "/some_index",
					ResponceCode: 404,
					Responce:     `{"error":{"root_cause":[{"type":"index_not_found_exception","reason":"no such index [some_index]","resource.type":"index_or_alias","resource.id":"some_index","index_uuid":"_na_","index":"some_index"}],"type":"index_not_found_exception","reason":"no such index [some_index]","resource.type":"index_or_alias","resource.id":"some_index","index_uuid":"_na_","index":"some_index"},"status":404}`,
				},
				2: {
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce:     `{"acknowledged":true,"shards_acknowledged":true,"index":"some_test"}`,
				},
			},
			Msg: "successfully created ES index some_index",
		},
		{
			//Unsuccessful Index creation
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
					RequestURI:   "/some_index",
					ResponceCode: 404,
					Responce:     `{"error":{"root_cause":[{"type":"index_not_found_exception","reason":"no such index [some_index]","resource.type":"index_or_alias","resource.id":"some_index","index_uuid":"_na_","index":"some_index"}],"type":"index_not_found_exception","reason":"no such index [some_index]","resource.type":"index_or_alias","resource.id":"some_index","index_uuid":"_na_","index":"some_index"},"status":404}`,
				},
				2: {
					RequestURI:   "/some_index",
					ResponceCode: 500,
					Responce:     `{}`,
				},
			},
			Err: fmt.Errorf("can't create ES index: elastic: Error 500 (Internal Server Error)"),
		},
		{
			//Un-aknowledged Index creation
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
					RequestURI:   "/some_index",
					ResponceCode: 404,
					Responce:     `{"error":{"root_cause":[{"type":"index_not_found_exception","reason":"no such index [some_index]","resource.type":"index_or_alias","resource.id":"some_index","index_uuid":"_na_","index":"some_index"}],"type":"index_not_found_exception","reason":"no such index [some_index]","resource.type":"index_or_alias","resource.id":"some_index","index_uuid":"_na_","index":"some_index"},"status":404}`,
				},
				2: {
					RequestURI:   "/some_index",
					ResponceCode: 200,
					Responce:     `{"acknowledged":false,"shards_acknowledged":true,"index":"some_test"}`,
				},
			},
			Err: fmt.Errorf("can't acknowledge ES index creation"),
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
		msg, err := client.CreateUpdateIndex(test.Index)
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

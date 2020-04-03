package elasticsearchindex

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/90poe/elasticsearch-objects-operator/pkg/consts"
	"github.com/90poe/elasticsearch-objects-operator/pkg/elasticsearch"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type TestCreateUpdateIndx struct {
	Index  *xov1alpha1.ElasticSearchIndex
	Status *xov1alpha1.ElasticSearchIndexStatus
	R2Rs   map[int]Responce2Req
}

func BeforeEachTest(t *testing.T) (*elasticsearch.Client, *TestDoer, *runtime.Scheme) {
	mockClient, testDoer := setupCreateTestClient(t)
	client, err := elasticsearch.New(
		elasticsearch.URL("http://localhost:9200"),
		elasticsearch.ESclient(mockClient),
	)
	assert.NoError(t, err)
	sc := scheme.Scheme
	sc.AddKnownTypes(xov1alpha1.SchemeGroupVersion, &xov1alpha1.ElasticSearchIndex{})
	sc.AddKnownTypes(xov1alpha1.SchemeGroupVersion, &xov1alpha1.ElasticSearchIndexList{})
	return client, testDoer, sc
}

func TestReconcile(t *testing.T) {
	var (
		name      = "test-elasticsearch"
		namespace = "operator"
	)

	client, testDoer, sc := BeforeEachTest(t)

	defer testDoer.Close()
	// Create fake K8S client
	cl := fake.NewFakeClient()

	// Create ReconcileElasticSearchIndex
	rp := &ReconcileElasticSearchIndex{
		client: cl,
		scheme: sc,
		es:     client,
	}
	// Create mock reconcile request
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	// Call Reconcile
	res, err := rp.Reconcile(req)
	assert.NoError(t, err)
	assert.Equal(t, res.Requeue, false)
}

func TestCreateUpdate(t *testing.T) {
	var (
		name      = "test-elasticsearch"
		namespace = "operator"
	)
	// Create mock reconcile request
	req := reconcile.Request{}
	tests := []TestCreateUpdateIndx{
		{
			//Succesfull create
			Index: &xov1alpha1.ElasticSearchIndex{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         name,
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
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: true,
				Operation:    consts.ESCreateOperation,
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 404,
					Responce: fmt.Sprintf(`
					{
						"error": {
						  "root_cause": [
							{
							  "type": "index_not_found_exception",
							  "reason": "no such index [%s]",
							  "resource.type": "index_or_alias",
							  "resource.id": "%s",
							  "index_uuid": "_na_",
							  "index": "%s"
							}
						  ],
						  "type": "index_not_found_exception",
						  "reason": "no such index [%s]",
						  "resource.type": "index_or_alias",
						  "resource.id": "%s",
						  "index_uuid": "_na_",
						  "index": "%s"
						},
						"status": 404
					  }
					`, name, name, name, name, name, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`{
						"acknowledged": true,
						"shards_acknowledged": true,
						"index": "%s"
					 }`, name),
				},
			},
		},
		{
			//Succesfull create for ES index which is present and managed by ES objects operator
			Index: &xov1alpha1.ElasticSearchIndex{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         name,
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards:   32,
						NumOfReplicas: 1,
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
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: true,
				Operation:    consts.ESCreateOperation,
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`
					{
						"%s": {
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
							  "number_of_replicas": "3",
							  "uuid": "iQXnF_YMTKqminns7h0-Zw",
							  "version": {
								"created": "7050299"
							  },
							  "provided_name": "%s"
							}
						  }
						}
					  }
					`, name, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/%s/_settings", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`{
						"acknowledged": true,
						"shards_acknowledged": true,
						"index": "%s"
					 }`, name),
				},
			},
		},
		{
			//Un-Succesfull create for ES index which is present and NOT managed by ES objects operator
			Index: &xov1alpha1.ElasticSearchIndex{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         name,
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
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: false,
				Operation:    consts.ESCreateOperation,
				LatestError:  "index 'test-elasticsearch' is not managed by this operator",
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`
					{
						"%s": {
						  "aliases": {},
						  "mappings": {
							"_meta": {
							  "managed-by": "xo"
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
							  "provided_name": "%s"
							}
						  }
						}
					  }
					`, name, name),
				},
			},
		},
		{
			//Un-succesfull create
			Index: &xov1alpha1.ElasticSearchIndex{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         name,
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
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: false,
				Operation:    consts.ESCreateOperation,
				LatestError:  "can't acknowledge ES index creation",
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 404,
					Responce: fmt.Sprintf(`
					{
						"error": {
						  "root_cause": [
							{
							  "type": "index_not_found_exception",
							  "reason": "no such index [%s]",
							  "resource.type": "index_or_alias",
							  "resource.id": "%s",
							  "index_uuid": "_na_",
							  "index": "%s"
							}
						  ],
						  "type": "index_not_found_exception",
						  "reason": "no such index [%s]",
						  "resource.type": "index_or_alias",
						  "resource.id": "%s",
						  "index_uuid": "_na_",
						  "index": "%s"
						},
						"status": 404
					  }
					`, name, name, name, name, name, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`{
						"acknowledged": false,
						"shards_acknowledged": true,
						"index": "%s"
					 }`, name),
				},
			},
		},
		{
			//Succesfull update
			Index: &xov1alpha1.ElasticSearchIndex{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  namespace,
					Finalizers: []string{"finalizer.elasticsearchindex.xo.90poe.io"},
				},
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         name,
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfReplicas: 5,
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
				Status: xov1alpha1.ElasticSearchIndexStatus{
					Operation:    consts.ESUpdateOperation,
					Name:         name,
					Acknowledged: true,
				},
			},
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: true,
				Operation:    consts.ESUpdateOperation,
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`
					{
						"%s": {
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
							  "provided_name": "%s"
							}
						  }
						}
					  }
					`, name, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/%s/_settings", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`{
						"acknowledged": true,
						"shards_acknowledged": true,
						"index": "%s"
					 }`, name),
				},
			},
		},
		{
			//Succesfull update after unsuccessfull update
			Index: &xov1alpha1.ElasticSearchIndex{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  namespace,
					Finalizers: []string{"finalizer.elasticsearchindex.xo.90poe.io"},
				},
				Spec: xov1alpha1.ElasticSearchIndexSpec{
					Name:         name,
					DropOnDelete: true,
					Settings: xov1alpha1.ESIndexSettings{
						NumOfReplicas: 5,
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
				Status: xov1alpha1.ElasticSearchIndexStatus{
					Operation:    consts.ESUpdateOperation,
					Name:         name,
					Acknowledged: false,
					LatestError:  "can't acknowledge ES index creation",
				},
			},
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: true,
				Operation:    consts.ESUpdateOperation,
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/%s", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`
					{
						"%s": {
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
							  "provided_name": "%s"
							}
						  }
						}
					  }
					`, name, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/%s/_settings", name),
					ResponceCode: 200,
					Responce: fmt.Sprintf(`{
						"acknowledged": true,
						"shards_acknowledged": true,
						"index": "%s"
					 }`, name),
				},
			},
		},
	}
	for _, test := range tests {
		client, testDoer, sc := BeforeEachTest(t)
		defer testDoer.Close()
		r2rKeys := make([]int, 0, len(test.R2Rs))
		for key := range test.R2Rs {
			r2rKeys = append(r2rKeys, key)
		}
		sort.Ints(r2rKeys)
		for _, value := range r2rKeys {
			testDoer.R2rChan <- test.R2Rs[value]
		}

		// Create fake K8S client
		cl := fake.NewFakeClient([]runtime.Object{test.Index}...)
		// Create ReconcileElasticSearchIndex
		rp := &ReconcileElasticSearchIndex{
			client: cl,
			scheme: sc,
			es:     client,
		}
		// Call Reconcile
		res, err := rp.Reconcile(req)
		assert.NoError(t, err)
		assert.Equal(t, res.Requeue, false)
		foundES := &xov1alpha1.ElasticSearchIndex{}
		err = cl.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, foundES)
		assert.NoError(t, err)
		assert.Equal(t, foundES.Status.Name, test.Status.Name)
		assert.Equal(t, foundES.Status.Acknowledged, test.Status.Acknowledged)
		assert.Equal(t, foundES.Status.Operation, test.Status.Operation)
		assert.Equal(t, foundES.Status.LatestError, test.Status.LatestError)
	}
}

func TestDelete(t *testing.T) {
	var (
		name      = "test-elasticsearch"
		namespace = "operator"
	)
	client, testDoer, sc := BeforeEachTest(t)
	defer testDoer.Close()
	//setup mock delete response
	r2r := Responce2Req{
		RequestURI:   fmt.Sprintf("/%s", name),
		ResponceCode: 200,
		Responce:     `{"acknowledged":true}`,
	}
	testDoer.R2rChan <- r2r

	now := metav1.NewTime(time.Now())
	esCR := &xov1alpha1.ElasticSearchIndex{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			DeletionTimestamp: &now,
			Finalizers:        []string{"finalizer.elasticsearchindex.xo.90poe.io"},
		},
		Spec: xov1alpha1.ElasticSearchIndexSpec{
			Name:         name,
			DropOnDelete: true,
			Settings: xov1alpha1.ESIndexSettings{
				NumOfShards: 10,
			},
			Mappings: "{}",
		},
		Status: xov1alpha1.ElasticSearchIndexStatus{
			Operation:    consts.ESUpdateOperation,
			Name:         name,
			Acknowledged: true,
		},
	}
	// Create fake K8S client
	cl := fake.NewFakeClient([]runtime.Object{esCR}...)
	// Create ReconcileElasticSearchIndex
	rp := &ReconcileElasticSearchIndex{
		client: cl,
		scheme: sc,
		es:     client,
	}
	// Create mock reconcile request
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	// Call Reconcile
	res, err := rp.Reconcile(req)
	assert.NoError(t, err)
	assert.Equal(t, res.Requeue, false)
	foundES := &xov1alpha1.ElasticSearchIndex{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, foundES)
	assert.NoError(t, err)
	assert.Equal(t, len(foundES.GetFinalizers()), 0)
}

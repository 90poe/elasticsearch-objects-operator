package elasticsearchtemplate

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

type TestUpdateTempl struct {
	Index  *xov1alpha1.ElasticSearchTemplate
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
	sc.AddKnownTypes(xov1alpha1.SchemeGroupVersion, &xov1alpha1.ElasticSearchTemplate{})
	sc.AddKnownTypes(xov1alpha1.SchemeGroupVersion, &xov1alpha1.ElasticSearchTemplateList{})
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

	// Create ReconcileElasticSearchTemplate
	rp := &ReconcileElasticSearchTemplate{
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

func TestUpdate(t *testing.T) {
	var (
		name      = "test-elasticsearch"
		namespace = "operator"
	)
	// Create mock reconcile request
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	tests := []TestUpdateTempl{
		{
			//Successful Update
			Index: &xov1alpha1.ElasticSearchTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  namespace,
					Finalizers: []string{"finalizer.elasticsearchtemplate.xo.90poe.io"},
				},
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          name,
					DropOnDelete:  true,
					IndexPatterns: []string{"some_index"},
					Settings: xov1alpha1.ESIndexSettings{
						NumOfShards: 33, //Diff here
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
				Status: xov1alpha1.ElasticSearchTemplateStatus{
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
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     fmt.Sprintf(`{"%s":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
		},
		{
			//Successful Update even if previously there were some errors
			Index: &xov1alpha1.ElasticSearchTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  namespace,
					Finalizers: []string{"finalizer.elasticsearchtemplate.xo.90poe.io"},
				},
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          name,
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
				Status: xov1alpha1.ElasticSearchTemplateStatus{
					Operation:    consts.ESCreateOperation,
					Name:         name,
					Acknowledged: false,
					LatestError:  "can't acknowledge ES template creation/update",
				},
			},
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: true,
				Operation:    consts.ESUpdateOperation,
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     fmt.Sprintf(`{"%s":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
		},
		{
			// Unsuccessful update
			Index: &xov1alpha1.ElasticSearchTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  namespace,
					Finalizers: []string{"finalizer.elasticsearchtemplate.xo.90poe.io"},
				},
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          name,
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
				Status: xov1alpha1.ElasticSearchTemplateStatus{
					Operation:    consts.ESUpdateOperation,
					Name:         name,
					Acknowledged: true,
				},
			},
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: false,
				Operation:    consts.ESUpdateOperation,
				LatestError:  "can't acknowledge ES template creation/update",
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     fmt.Sprintf(`{"%s":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":false}`,
				},
			},
		},
		{
			// Successful creation
			Index: &xov1alpha1.ElasticSearchTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          name,
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
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: true,
				Operation:    consts.ESCreateOperation,
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 404,
					Responce:     "{}",
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
		},
		{
			// Unsuccessful creation
			Index: &xov1alpha1.ElasticSearchTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          name,
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
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: false,
				Operation:    consts.ESCreateOperation,
				LatestError:  "can't acknowledge ES template creation/update",
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 404,
					Responce:     "{}",
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":false}`,
				},
			},
		},
		{
			// Successful no diff update (No real update)
			Index: &xov1alpha1.ElasticSearchTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  namespace,
					Finalizers: []string{"finalizer.elasticsearchtemplate.xo.90poe.io"},
				},
				Spec: xov1alpha1.ElasticSearchTemplateSpec{
					Name:          name,
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
				Status: xov1alpha1.ElasticSearchTemplateStatus{
					Operation:    consts.ESUpdateOperation,
					Name:         name,
					Acknowledged: false,
				},
			},
			Status: &xov1alpha1.ElasticSearchIndexStatus{
				Name:         name,
				Acknowledged: true,
				Operation:    consts.ESUpdateOperation,
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     fmt.Sprintf(`{"%s":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-objects-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":false}`,
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
		// Create ReconcileElasticSearchTemplate
		rp := &ReconcileElasticSearchTemplate{
			client: cl,
			scheme: sc,
			es:     client,
		}
		// Call Reconcile
		res, err := rp.Reconcile(req)
		assert.NoError(t, err)
		assert.Equal(t, res.Requeue, false)
		foundES := &xov1alpha1.ElasticSearchTemplate{}
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
		RequestURI:   fmt.Sprintf("/_template/%s", name),
		ResponceCode: 200,
		Responce:     `{"acknowledged":true}`,
	}
	testDoer.R2rChan <- r2r

	now := metav1.NewTime(time.Now())
	esCR := &xov1alpha1.ElasticSearchTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			DeletionTimestamp: &now,
			Finalizers:        []string{"finalizer.elasticsearchtemplate.xo.90poe.io"},
		},
		Spec: xov1alpha1.ElasticSearchTemplateSpec{
			Name:          name,
			DropOnDelete:  true,
			IndexPatterns: []string{"some_index"},
			Settings: xov1alpha1.ESIndexSettings{
				NumOfShards: 10,
			},
			Mappings: "{}",
		},
		Status: xov1alpha1.ElasticSearchTemplateStatus{
			Operation:    consts.ESUpdateOperation,
			Name:         name,
			Acknowledged: true,
		},
	}
	// Create fake K8S client
	cl := fake.NewFakeClient([]runtime.Object{esCR}...)
	// Create ReconcileElasticSearchTemplate
	rp := &ReconcileElasticSearchTemplate{
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
	foundES := &xov1alpha1.ElasticSearchTemplate{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, foundES)
	assert.NoError(t, err)
	assert.Equal(t, len(foundES.GetFinalizers()), 0)
}

package elasticsearchtemplate

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/90poe/elasticsearch-operator/pkg/consts"
	"github.com/90poe/elasticsearch-operator/pkg/elasticsearch"

	xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

type TestCreateTempl struct {
	Index *xov1alpha1.ElasticSearchTemplate
	R2R   Responce2Req
	Err   error
}

type TestUpdateTempl struct {
	Index *xov1alpha1.ElasticSearchTemplate
	R2Rs  map[int]Responce2Req
	Err   error
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

func TestCreate(t *testing.T) {
	var (
		name      = "test-elasticsearch"
		namespace = "operator"
	)
	// Create mock reconcile request
	req := reconcile.Request{}
	tests := []TestCreateTempl{
		{
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
			R2R: Responce2Req{
				RequestURI:   fmt.Sprintf("/_template/%s", name),
				ResponceCode: 200,
				Responce:     `{"acknowledged":true}`,
			},
		},
		{
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
			R2R: Responce2Req{
				RequestURI:   fmt.Sprintf("/_template/%s", name),
				ResponceCode: 400,
				Responce:     `{"acknowledged":false}`,
			},
			Err: fmt.Errorf("can't create or update ES template: elastic: Error 400 (Bad Request)"),
		},
	}
	for _, test := range tests {
		client, testDoer, sc := BeforeEachTest(t)
		defer testDoer.Close()
		testDoer.R2rChan <- test.R2R

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
		if test.Err != nil {
			assert.Equal(t, foundES.Status.Acknowledged, false)
			assert.Equal(t, foundES.Status.Operation, consts.ESCreateOperation)
			assert.Equal(t, foundES.Status.LatestError, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.Equal(t, foundES.Status.Name, name)
		assert.Equal(t, foundES.Status.Acknowledged, true)
		assert.Equal(t, foundES.Status.Operation, consts.ESCreateOperation)
		assert.Equal(t, foundES.Status.LatestError, "")
	}
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
					Acknowledged: true,
				},
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     fmt.Sprintf(`{"%s":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":true}`,
				},
			},
		},
		{
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
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     fmt.Sprintf(`{"%s":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`, name),
				},
				2: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     `{"acknowledged":false}`,
				},
			},
			Err: fmt.Errorf("can't acknowledge ES template creation/update"),
		},
		{
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
				},
			},
			R2Rs: map[int]Responce2Req{
				1: {
					RequestURI:   fmt.Sprintf("/_template/%s", name),
					ResponceCode: 200,
					Responce:     fmt.Sprintf(`{"%s":{"order":0,"index_patterns":["some_index"],"settings":{"index":{"number_of_shards":"32"}},"mappings":{"_meta":{"managed-by":"elasticsearch-operator.xo.90poe.io"},"_source":{"enabled":true},"dynamic":false,"properties":{"createdAt":{"index":true,"type":"date"},"isRead":{"index":true,"type":"boolean"}}},"aliases":{"add":{},"remove_index":{},"remove":{}}}}`, name),
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
		if test.Err != nil {
			assert.Equal(t, foundES.Status.Acknowledged, false)
			assert.Equal(t, foundES.Status.Operation, consts.ESUpdateOperation)
			assert.Equal(t, foundES.Status.LatestError, fmt.Sprintf("%s", test.Err))
			continue
		}
		assert.Equal(t, foundES.Status.Name, test.Index.Name)
		assert.Equal(t, foundES.Status.Acknowledged, test.Index.Status.Acknowledged)
		assert.Equal(t, foundES.Status.Operation, test.Index.Status.Operation)
		assert.Equal(t, foundES.Status.LatestError, "")

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

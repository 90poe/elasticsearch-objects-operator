package elasticsearchindex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/olivere/elastic/v7"
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
	//FOR INTERNAL REVIEW PURPOSES
	// if req.Method == "PUT" {
	// 	bytes, err := ioutil.ReadAll(req.Body)
	// 	defer req.Body.Close()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Printf("%s", string(bytes))
	// }
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

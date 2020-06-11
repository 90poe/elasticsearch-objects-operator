package utils

import (
	"context"
	"strconv"
	"testing"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	name      string = "test-es"
	namespace string = "operator"
)

var esIndex *xov1alpha1.ElasticSearchIndex = &xov1alpha1.ElasticSearchIndex{
	ObjectMeta: metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	},
	Spec: xov1alpha1.ElasticSearchIndexSpec{
		Name: name,
	},
}
var objs []runtime.Object = []runtime.Object{esIndex}

func TestPatchUtilShouldPatchIfThereIsDifference(t *testing.T) {
	// Create modified postgres
	modESIndex := esIndex.DeepCopy()
	modESIndex.Spec.Settings.NumOfReplicas = 10
	modESIndex.Status.Acknowledged = true
	modESIndex.Status.Operation = "create"

	// Create runtime scheme
	s := scheme.Scheme
	s.AddKnownTypes(xov1alpha1.SchemeGroupVersion, &xov1alpha1.ElasticSearchIndex{})

	// Create fake client to mock API calls
	cl := fake.NewFakeClient(objs...)

	// Patch object
	err := Patch(context.TODO(), cl, esIndex, modESIndex)
	if err != nil {
		t.Fatalf("could not patch object: (%v)", err)
	}

	// Check if esindex is identical to modified object
	foundESindex := &xov1alpha1.ElasticSearchIndex{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, foundESindex)
	if err != nil {
		t.Fatalf("could not get postgres: (%v)", err)
	}
	// Comparison
	if foundESindex.Spec.Settings.NumOfReplicas != modESIndex.Spec.Settings.NumOfReplicas {
		t.Fatalf("found ESIndex is not identical to modified ESIndex: NumOfReplicas == %d, expected %d",
			foundESindex.Spec.Settings.NumOfReplicas, modESIndex.Spec.Settings.NumOfReplicas)
	}
	if foundESindex.Status.Acknowledged != modESIndex.Status.Acknowledged {
		t.Fatalf("found ESIndex is not identical to modified ESIndex: Succeeded == %s, expected %s",
			strconv.FormatBool(foundESindex.Status.Acknowledged), strconv.FormatBool(modESIndex.Status.Acknowledged))
	}
	if foundESindex.Status.Operation != modESIndex.Status.Operation {
		t.Fatalf("found ESIndex is not identical to modified ESIndex: Operation == %s, expected %s",
			foundESindex.Status.Operation, modESIndex.Status.Operation)
	}
}

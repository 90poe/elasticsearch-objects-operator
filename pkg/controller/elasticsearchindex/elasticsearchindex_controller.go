package elasticsearchindex

import (
	"context"
	"fmt"

	xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"
	"github.com/90poe/elasticsearch-operator/pkg/config"
	"github.com/90poe/elasticsearch-operator/pkg/consts"
	"github.com/90poe/elasticsearch-operator/pkg/elasticsearch"
	"github.com/90poe/elasticsearch-operator/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"

	// metav1 "k8s.io/apimachinery/pkg/api/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const ()

var log = logf.Log.WithName("controller_elasticsearchindex")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ElasticSearchIndex Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	c := config.Get()
	es, err := elasticsearch.New(
		elasticsearch.URL(c.ESurl),
	)
	if err != nil {
		log.Error(err, "can't make new es client")
		return nil
	}
	return &ReconcileElasticSearchIndex{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		es:     es,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("elasticsearchindex-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ElasticSearchIndex
	err = c.Watch(&source.Kind{Type: &xov1alpha1.ElasticSearchIndex{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileElasticSearchIndex implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileElasticSearchIndex{}

// ReconcileElasticSearchIndex reconciles a ElasticSearchIndex object
type ReconcileElasticSearchIndex struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	es     elasticsearch.ES
}

// Reconcile reads that state of the cluster for a ElasticSearchIndex object and makes changes based on the state read
// and what is in the ElasticSearchIndex.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileElasticSearchIndex) Reconcile(request reconcile.Request) (_ reconcile.Result, reterr error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ElasticSearchIndex")

	// Fetch the ElasticSearchIndex instance
	instance := &xov1alpha1.ElasticSearchIndex{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	before := instance.DeepCopyObject()
	// Patch after every reconcile loop, if needed
	defer func() {
		err = utils.Patch(context.TODO(), r.client, before, instance)
		if err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()
	// General logic:
	// 1. If Operation == "" - attempt to create. Set Operation="create". If no Acknowledged, Acknowledged=false
	// 2. If Operation == "created" && Acknowledged - try to update. Set Operation = "update"
	// 3. Always set LatestError if error occured

	// deletion logic
	if !instance.GetDeletionTimestamp().IsZero() {
		if r.shouldDeleteIndex(instance, reqLogger) {
			// All errors ignored
			err = r.es.DeleteIndex(instance.Spec.Name)
			if err != nil {
				log.Error(err, fmt.Sprintf("error deleting '%s' index '%s': %v",
					instance.Name, instance.Spec.Name, err))
			} else {
				log.Info(fmt.Sprintf("successfully deleted ES index %s", instance.Spec.Name))
			}
		}
		instance.SetFinalizers(nil)

		log.Info(fmt.Sprintf("succesfully deleted CRD %s from K8S", instance.Name))
		return reconcile.Result{}, nil
	}

	//We called second time on creation event - nothing to do
	if instance.Status.Operation == consts.ESCreateOperation &&
		instance.ObjectMeta.Generation == 1 {
		return reconcile.Result{}, nil
	}

	switch instance.Status.Operation {
	case "":
		// Create index
		instance.Status.Name = instance.Spec.Name
		instance.Status.Operation = consts.ESCreateOperation
		err = r.es.CreateIndex(instance)
		if err != nil {
			instance.Status.LatestError = fmt.Sprintf("%v", err)
			log.Info(instance.Status.LatestError)
			return reconcile.Result{}, nil
		}
		instance.Status.Acknowledged = true
		log.Info(fmt.Sprintf("successfully created ES index %s", instance.Spec.Name))
	case consts.ESCreateOperation, consts.ESUpdateOperation:
		// Update index
		if instance.Status.Operation == consts.ESCreateOperation &&
			!instance.Status.Acknowledged {
			//Create operation was unsuccessful - ignore update
			log.Info(fmt.Sprintf("trying to update index '%s' which failed to create - ignoring",
				instance.Spec.Name))
			return reconcile.Result{}, nil
		}
		instance.Status.LatestError = ""
		instance.Status.Operation = consts.ESUpdateOperation
		msg, err := r.es.UpdateIndex(instance)
		if err != nil {
			instance.Status.Acknowledged = false
			instance.Status.LatestError = fmt.Sprintf("%v", err)
			log.Info(instance.Status.LatestError)
			return reconcile.Result{}, nil
		}
		if len(msg) != 0 {
			log.Info(msg)
		}
	}

	err = r.addFinalizer(instance, reqLogger)
	if err != nil {
		log.Error(err, "can't add finalizer")
		return r.requeue(instance, err)
	}

	reqLogger.Info("reconciler done", "CR.Namespace", instance.Namespace, "CR.Name", instance.Name)

	return reconcile.Result{}, nil
}

func (r *ReconcileElasticSearchIndex) addFinalizer(m *xov1alpha1.ElasticSearchIndex, reqLogger logr.Logger) error {
	if len(m.GetFinalizers()) < 1 && m.GetDeletionTimestamp() == nil {
		reqLogger.Info("adding Finalizer for ESIndex")
		m.SetFinalizers([]string{"finalizer.elasticsearchindex.90poe.io"})
	}
	return nil
}

func (r *ReconcileElasticSearchIndex) requeue(cr *xov1alpha1.ElasticSearchIndex, reason error) (reconcile.Result, error) {
	// cr.Status.Acknowledged = false
	return reconcile.Result{}, reason
}

func (r *ReconcileElasticSearchIndex) shouldDeleteIndex(cr *xov1alpha1.ElasticSearchIndex, logger logr.Logger) bool {
	// If DropOnDelete is false we don't need to check any further
	if !cr.Spec.DropOnDelete {
		return false
	}
	if !cr.Status.Acknowledged {
		// If we don't have aknowledge from Create/Update - we don't delete
		return false
	}
	// Get a list of all ES Indexes
	indices := xov1alpha1.ElasticSearchIndexList{}
	err := r.client.List(context.TODO(), &indices, &client.ListOptions{})
	if err != nil {
		logger.Info(fmt.Sprintf("%v", err))
		return true
	}

	for _, index := range indices.Items {
		// Skip index if it's the same as the one we're deleting
		if index.Name == cr.Name && index.Namespace == cr.Namespace {
			continue
		}
		// There already exists another ESIndex who has the same index
		// Let's not delete the index
		if index.Spec.Name == cr.Spec.Name {
			return false
		}
	}

	return true
}

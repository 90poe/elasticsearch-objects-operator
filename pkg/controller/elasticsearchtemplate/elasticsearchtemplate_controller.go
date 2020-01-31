package elasticsearchtemplate

import (
	"context"
	"fmt"

	xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"
	"github.com/90poe/elasticsearch-operator/pkg/config"
	"github.com/90poe/elasticsearch-operator/pkg/elasticsearch"
	"github.com/90poe/elasticsearch-operator/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
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

var log = logf.Log.WithName("controller_elasticsearchtemplate")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ElasticSearchTemplate Controller and adds it to the Manager. The Manager will set fields on the Controller
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
	return &ReconcileElasticSearchTemplate{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		es:     es,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("elasticsearchtemplate-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ElasticSearchTemplate
	err = c.Watch(&source.Kind{Type: &xov1alpha1.ElasticSearchTemplate{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileElasticSearchTemplate implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileElasticSearchTemplate{}

// ReconcileElasticSearchTemplate reconciles a ElasticSearchTemplate object
type ReconcileElasticSearchTemplate struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	es     elasticsearch.ES
}

// Reconcile reads that state of the cluster for a ElasticSearchTemplate object and makes changes based on the state read
// and what is in the ElasticSearchTemplate.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileElasticSearchTemplate) Reconcile(request reconcile.Request) (_ reconcile.Result, reterr error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ElasticSearchTemplate")

	// Fetch the ElasticSearchTemplate instance
	instance := &xov1alpha1.ElasticSearchTemplate{}
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

	// deletion logic
	if !instance.GetDeletionTimestamp().IsZero() {
		if r.shouldDeleteTemplate(instance, reqLogger) && instance.Status.Succeeded {
			err = r.es.DeleteTemplate(instance.Spec.Name, reqLogger)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		instance.SetFinalizers(nil)

		log.Info(fmt.Sprintf("succesfully deleted CRD %s from K8S", instance.Name))
		return reconcile.Result{}, nil
	}

	// creation logic
	if !instance.Status.Succeeded {
		// Create template
		err = r.es.CreateTemplate(instance, reqLogger)
		if err != nil {
			log.Error(err, "can't create template")
			return reconcile.Result{}, nil
		}
		instance.Status.Succeeded = true
		instance.Status.Name = instance.Spec.Name
	} else {
		// Update index
		err = r.es.UpdateTemplate(instance, reqLogger)
		if err != nil {
			log.Error(err, "can't update template")
			return reconcile.Result{}, nil
		}
	}
	err = r.addFinalizer(instance, reqLogger)
	if err != nil {
		return r.requeue(instance, err)
	}

	reqLogger.Info("reconciler done", "CR.Namespace", instance.Namespace, "CR.Name", instance.Name)

	return reconcile.Result{}, nil
}

func (r *ReconcileElasticSearchTemplate) addFinalizer(m *xov1alpha1.ElasticSearchTemplate, reqLogger logr.Logger) error {
	if len(m.GetFinalizers()) < 1 && m.GetDeletionTimestamp() == nil {
		reqLogger.Info("adding Finalizer for EStemplate")
		m.SetFinalizers([]string{"finalizer.elasticsearchtemplate.90poe.io"})
	}
	return nil
}

func (r *ReconcileElasticSearchTemplate) requeue(cr *xov1alpha1.ElasticSearchTemplate, reason error) (reconcile.Result, error) {
	cr.Status.Succeeded = false
	return reconcile.Result{}, reason
}

func (r *ReconcileElasticSearchTemplate) shouldDeleteTemplate(cr *xov1alpha1.ElasticSearchTemplate, logger logr.Logger) bool {
	// If DropOnDelete is false we don't need to check any further
	if !cr.Spec.DropOnDelete {
		return false
	}
	// Get a list of all ES Indexes
	templates := xov1alpha1.ElasticSearchTemplateList{}
	err := r.client.List(context.TODO(), &templates, &client.ListOptions{})
	if err != nil {
		logger.Info(fmt.Sprintf("%v", err))
		return true
	}

	for _, template := range templates.Items {
		// Skip template if it's the same as the one we're deleting
		if template.Name == cr.Name && template.Namespace == cr.Namespace {
			continue
		}
		// There already exists another ESTemplate who has the same template
		// Let's not delete the template
		if template.Spec.Name == cr.Spec.Name {
			return false
		}
	}

	return true
}

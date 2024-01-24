/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/api/v1alpha1"
	"github.com/90poe/elasticsearch-objects-operator/internal/config"
	"github.com/90poe/elasticsearch-objects-operator/internal/elasticsearch"
	"github.com/90poe/elasticsearch-objects-operator/internal/reporter"
	"github.com/go-logr/logr"
)

// ElasticSearchIndexReconciler reconciles a ElasticSearchIndex object
type ElasticSearchIndexReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	es        elasticsearch.ES
	messenger *reporter.Messenger
}

//+kubebuilder:rbac:groups=xo.90poe.io,resources=elasticsearchindices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=xo.90poe.io,resources=elasticsearchindices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=xo.90poe.io,resources=elasticsearchindices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ElasticSearchIndex object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *ElasticSearchIndexReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx).WithValues("elasticsearchindex", req.NamespacedName)

	// Fetch the ElasticSearchIndex instance
	instance := &xov1alpha1.ElasticSearchIndex{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.V(0).Info("ElasticSearchIndex resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.V(0).Info(fmt.Sprintf("Failed to get ElasticSearchIndex: %v", err))
		return reconcile.Result{}, err
	}

	return r.upsertIndex(ctx, instance, reqLogger)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ElasticSearchIndexReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c := config.Get()
	es, err := elasticsearch.New(
		elasticsearch.URL(c.ESurl),
	)
	if err != nil {
		return err
	}
	r.es = es
	// Make slack messenger
	r.messenger, err = reporter.New(c.SlackToken,
		reporter.SlackChannel(c.SlackChannel))
	if err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&xov1alpha1.ElasticSearchIndex{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: c.MaxConcurrentReconciles}).
		WithEventFilter(ignoreUpdateDeletePredicate()).
		Complete(r)
}

// upsertIndex will update or insert index in ES cluster
func (r *ElasticSearchIndexReconciler) upsertIndex(ctx context.Context, index *xov1alpha1.ElasticSearchIndex, reqLogger logr.Logger) (_ ctrl.Result, retErr error) {
	// Init status
	statusMessage := "Succeeded"
	status := metav1.ConditionTrue
	condition := ConditionsInsert
	reason := ConditionReasonCreateIndex

	// Defer function to update status
	defer func() {
		// Log status update
		reqLogger.Info(fmt.Sprintf("elasticsearch %s %s status: %s", index.Spec.Name,
			reason, statusMessage))
		// Send message to slack
		if status == metav1.ConditionFalse {
			// send message only on error
			r.messenger.Send(statusMessage, reporter.ErrorMessage)
		}
		// Remove last condition and set new one
		meta.RemoveStatusCondition(&index.Status.Conditions, condition)
		meta.SetStatusCondition(&index.Status.Conditions, metav1.Condition{
			Type:    condition,
			Status:  status,
			Reason:  reason,
			Message: statusMessage,
		})
		// we will return error of status update if it is not nil
		err := r.Status().Update(ctx, index)
		if err != nil {
			reqLogger.V(0).Info(fmt.Sprintf("Failed to update index status: %v", retErr))
			retErr = errors.Join(retErr, err)
		}
	}()

	// Check if index exists in ES cluster
	exists, err := r.es.IndexExists(index.Spec.Name)
	if err != nil {
		status = metav1.ConditionFalse
		statusMessage = fmt.Sprintf("can't check if index %s exists: %v", index.Spec.Name, err)
		return ctrl.Result{}, nil
	}

	// Create or update topic
	if exists {
		condition = ConditionsUpdate
		reason = ConditionReasonUpdateIndex
	}
	_, err = r.es.CreateUpdateIndex(index)

	if err != nil {
		status = metav1.ConditionFalse
		statusMessage = fmt.Sprintf("can't %s ES index %s: %v", reason, index.Name, err)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{
		RequeueAfter: RevisitIntervalSec * time.Second,
	}, nil
}

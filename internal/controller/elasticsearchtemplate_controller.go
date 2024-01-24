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

// ElasticSearchTemplateReconciler reconciles a ElasticSearchTemplate object
type ElasticSearchTemplateReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	es        elasticsearch.ES
	messenger *reporter.Messenger
}

//+kubebuilder:rbac:groups=xo.90poe.io,resources=elasticsearchtemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=xo.90poe.io,resources=elasticsearchtemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=xo.90poe.io,resources=elasticsearchtemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ElasticSearchTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *ElasticSearchTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx).WithValues("elasticsearchtemplate", req.NamespacedName)

	// Fetch the ElasticSearchTemplate instance
	instance := &xov1alpha1.ElasticSearchTemplate{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.V(0).Info("ElasticSearchTemplate resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.V(0).Info(fmt.Sprintf("Failed to get ElasticSearchTemplate: %v", err))
		return reconcile.Result{}, err
	}

	return r.upsertTemplate(ctx, instance, reqLogger)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ElasticSearchTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
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
		For(&xov1alpha1.ElasticSearchTemplate{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: c.MaxConcurrentReconciles}).
		WithEventFilter(ignoreUpdateDeletePredicate()).
		Complete(r)
}

// upsertTemplate will update or insert index in ES cluster
func (r *ElasticSearchTemplateReconciler) upsertTemplate(ctx context.Context, template *xov1alpha1.ElasticSearchTemplate, reqLogger logr.Logger) (_ ctrl.Result, retErr error) {
	// Init status
	statusMessage := "Succeeded"
	status := metav1.ConditionTrue
	condition := ConditionsInsert
	reason := ConditionReasonCreateTemplate

	// Defer function to update status
	defer func() {
		// Log status update
		reqLogger.Info(fmt.Sprintf("elasticsearch %s %s status: %s", template.Spec.Name,
			reason, statusMessage))
		// Send message to slack
		if status == metav1.ConditionFalse {
			// send message only on error
			r.messenger.Send(statusMessage, reporter.ErrorMessage)
		}
		// Remove last condition and set new one
		meta.RemoveStatusCondition(&template.Status.Conditions, condition)
		meta.SetStatusCondition(&template.Status.Conditions, metav1.Condition{
			Type:    condition,
			Status:  status,
			Reason:  reason,
			Message: statusMessage,
		})
		// we will return error of status update if it is not nil
		err := r.Status().Update(ctx, template)
		if err != nil {
			reqLogger.V(0).Info(fmt.Sprintf("Failed to update template status: %v", retErr))
			retErr = errors.Join(retErr, err)
		}
	}()

	// Check if index exists in ES cluster
	exists, err := r.es.TemplateExists(template.Spec.Name)
	if err != nil {
		status = metav1.ConditionFalse
		statusMessage = fmt.Sprintf("can't check if template %s exists: %v", template.Spec.Name, err)
		return ctrl.Result{}, nil
	}

	// Create or update topic
	if exists {
		condition = ConditionsUpdate
		reason = ConditionReasonUpdateIndex
	}
	_, err = r.es.CreateUpdateTemplate(template)

	if err != nil {
		status = metav1.ConditionFalse
		statusMessage = fmt.Sprintf("can't %s ES template %s: %v", reason, template.Name, err)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{
		RequeueAfter: RevisitIntervalSec * time.Second,
	}, nil
}

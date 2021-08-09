package cmd

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"go.uber.org/zap"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	argoevents "github.com/argoproj/argo-events"
	"github.com/argoproj/argo-events/common"
	"github.com/argoproj/argo-events/common/logging"
	"github.com/argoproj/argo-events/controllers/sensor"
	eventbusv1alpha1 "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1"
	eventsourcev1alpha1 "github.com/argoproj/argo-events/pkg/apis/eventsource/v1alpha1"
	sensorv1alpha1 "github.com/argoproj/argo-events/pkg/apis/sensor/v1alpha1"
)

const (
	sensorImageEnvVar = "SENSOR_IMAGE"
)

func Start(namespaced bool, managedNamespace string) {
	logger := logging.NewArgoEventsLogger().Named(sensor.ControllerName)
	sensorImage, defined := os.LookupEnv(sensorImageEnvVar)
	if !defined {
		logger.Fatalf("required environment variable '%s' not defined", sensorImageEnvVar)
	}
	opts := ctrl.Options{
		MetricsBindAddress:     fmt.Sprintf(":%d", common.ControllerMetricsPort),
		HealthProbeBindAddress: ":8081",
	}
	if namespaced {
		opts.Namespace = managedNamespace
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), opts)
	if err != nil {
		logger.Fatalw("unable to get a controller-runtime manager", zap.Error(err))
	}

	// Readyness probe
	if err := mgr.AddReadyzCheck("readiness", healthz.Ping); err != nil {
		logger.Fatalw("unable add a readiness check", zap.Error(err))
	}

	// Liveness probe
	if err := mgr.AddHealthzCheck("liveness", healthz.Ping); err != nil {
		logger.Fatalw("unable add a health check", zap.Error(err))
	}

	if err := sensorv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Fatalw("unable to add Sensor scheme", zap.Error(err))
	}

	if err := eventsourcev1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Fatalw("unable to add EventSource scheme", zap.Error(err))
	}

	if err := eventbusv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Fatalw("unable to add EventBus scheme", zap.Error(err))
	}

	// A controller with DefaultControllerRateLimiter
	c, err := controller.New(sensor.ControllerName, mgr, controller.Options{
		Reconciler: sensor.NewReconciler(mgr.GetClient(), mgr.GetScheme(), sensorImage, logger),
	})
	if err != nil {
		logger.Fatalw("unable to set up individual controller", zap.Error(err))
	}

	// Watch Sensor and enqueue Sensor object key
	if err := c.Watch(&source.Kind{Type: &sensorv1alpha1.Sensor{}}, &handler.EnqueueRequestForObject{},
		predicate.Or(
			predicate.GenerationChangedPredicate{},
			// TODO: change to use LabelChangedPredicate with controller-runtime v0.8
			predicate.Funcs{
				UpdateFunc: func(e event.UpdateEvent) bool {
					if e.ObjectOld == nil {
						return false
					}
					if e.ObjectNew == nil {
						return false
					}
					return !reflect.DeepEqual(e.ObjectNew.GetLabels(), e.ObjectOld.GetLabels())
				},
			},
		)); err != nil {
		logger.Fatalw("unable to watch Sensors", zap.Error(err))
	}

	// Watch Deployments and enqueue owning Sensor key
	if err := c.Watch(&source.Kind{Type: &appv1.Deployment{}}, &handler.EnqueueRequestForOwner{OwnerType: &sensorv1alpha1.Sensor{}, IsController: true}, predicate.GenerationChangedPredicate{}); err != nil {
		logger.Fatalw("unable to watch Deployments", zap.Error(err))
	}

	// Watch EventSources and enqueue effected Sensor keys (Sensors with filter dependencies)
	esEventsHandler := createEventSourceEventsHandler(mgr.GetClient(), logger)
	if err := c.Watch(&source.Kind{Type: &eventsourcev1alpha1.EventSource{}}, handler.EnqueueRequestsFromMapFunc(esEventsHandler),
		predicate.And(
			predicate.GenerationChangedPredicate{},
			predicate.Funcs{
				// TODO decide if sensor update on ES deletion is required
				DeleteFunc: func(e event.DeleteEvent) bool {
					return false
				},
				UpdateFunc: func(e event.UpdateEvent) bool {
					if e.ObjectOld == nil {
						return false
					}
					if e.ObjectNew == nil {
						return false
					}
					esOld, ok := e.ObjectOld.(*eventsourcev1alpha1.EventSource)
					if !ok {
						return false
					}
					esNew, ok := e.ObjectNew.(*eventsourcev1alpha1.EventSource)
					if !ok {
						return false
					}
					return !reflect.DeepEqual(esOld.Spec, esNew.Spec)
				},
			},
		)); err != nil {
		logger.Fatalw("unable to watch EventSources", zap.Error(err))
	}

	logger.Infow("starting sensor controller", "version", argoevents.GetVersion())
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Fatalw("unable to run sensor controller", zap.Error(err))
	}
}

func createEventSourceEventsHandler(cl client.Client, logger *zap.SugaredLogger) handler.MapFunc {
	return func(obj client.Object) []reconcile.Request {
		var requests []reconcile.Request
		sl := &sensorv1alpha1.SensorList{}
		if err := cl.List(context.TODO(), sl, &client.ListOptions{
			Namespace: obj.GetNamespace(),
		}); err != nil {
			logger.Fatalw("unable to list Sensors", zap.Error(err))
			return requests
		}

		for _, s := range sl.Items {
			for _, d := range s.Spec.Dependencies {
				if len(d.EventSourceFilter) != 0 {
					requests = append(requests, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name:      s.Name,
							Namespace: s.Namespace,
						},
					})
				}
			}
		}
		return requests
	}
}

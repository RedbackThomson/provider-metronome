/*
Copyright 2020 RedbackThomson.

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

package main

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"
	"go.uber.org/zap/zapcore"

	"k8s.io/client-go/tools/leaderelection/resourcelock"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"gopkg.in/alecthomas/kingpin.v2"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/redbackthomson/provider-metronome/apis"
	metronomeControllers "github.com/redbackthomson/provider-metronome/internal/controller"
)

func main() {
	var (
		app                     = kingpin.New(filepath.Base(os.Args[0]), "Metronome provider for Crossplane.").DefaultEnvars()
		debug                   = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		leaderElection          = app.Flag("leader-election", "Use leader election for the conroller manager.").Short('l').Default("false").Envar("LEADER_ELECTION").Bool()
		syncInterval            = app.Flag("sync", "How often all resources will be double-checked for drift from the desired state.").Short('s').Default("1h").Duration()
		pollInterval            = app.Flag("poll", "How often individual resources will be checked for drift from the desired state").Default("10m").Duration()
		pollStateMetricInterval = app.Flag("poll-state-metric", "State metric recording interval").Default("5s").Duration()
		maxReconcileRate        = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may checked for drift from the desired state.").Default("100").Int()

		enableManagementPolicies = app.Flag("enable-management-policies", "Enable support for Management Policies.").Default("true").Envar("ENABLE_MANAGEMENT_POLICIES").Bool()

		metronomeBaseUrl = app.Flag("metronome-base-url", "Base URL to use for all Metronome API requests").Default("https://api.metronome.com").Envar("METRONOME_BASE_URL").String()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug), UseISO8601())
	log := logging.NewLogrLogger(zl.WithName("provider-metronome"))
	// explicitly  provide a no-op logger by default, otherwise controller-runtime gives a warning
	ctrl.SetLogger(zap.New(zap.WriteTo(io.Discard)))
	if *debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	mgr, err := ctrl.NewManager(ratelimiter.LimitRESTConfig(cfg, *maxReconcileRate), ctrl.Options{
		Cache: cache.Options{
			SyncPeriod: syncInterval,
		},

		// controller-runtime uses both ConfigMaps and Leases for leader
		// election by default. Leases expire after 15 seconds, with a
		// 10 second renewal deadline. We've observed leader loss due to
		// renewal deadlines being exceeded when under high load - i.e.
		// hundreds of reconciles per second and ~200rps to the API
		// server. Switching to Leases only and longer leases appears to
		// alleviate this.
		LeaderElection:             *leaderElection,
		LeaderElectionID:           "crossplane-leader-election-provider-template",
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")

	mm := managed.NewMRMetricRecorder()
	sm := statemetrics.NewMRStateMetrics()

	metrics.Registry.MustRegister(mm)
	metrics.Registry.MustRegister(sm)

	mo := controller.MetricOptions{
		PollStateMetricInterval: *pollStateMetricInterval,
		MRMetrics:               mm,
		MRStateMetrics:          sm,
	}

	kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add Metronome APIs to scheme")
	o := controller.Options{
		Logger:                  log,
		MaxConcurrentReconciles: *maxReconcileRate,
		PollInterval:            *pollInterval,
		GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
		Features:                &feature.Flags{},
		MetricOptions:           &mo,
	}

	if *enableManagementPolicies {
		o.Features.Enable(feature.EnableBetaManagementPolicies)
		log.Info("Beta feature enabled", "flag", feature.EnableBetaManagementPolicies)
	}

	kingpin.FatalIfError(metronomeControllers.Setup(mgr, o, *metronomeBaseUrl), "Cannot setup Template controllers")
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

// UseISO8601 sets the logger to use ISO8601 timestamp format
func UseISO8601() zap.Opts {
	return func(o *zap.Options) {
		o.TimeEncoder = zapcore.ISO8601TimeEncoder
	}
}

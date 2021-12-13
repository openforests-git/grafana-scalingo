// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire gen -tags "oss"
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/google/wire"
	httpclient2 "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana/pkg/api"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/expr"
	"github.com/grafana/grafana/pkg/infra/httpclient"
	"github.com/grafana/grafana/pkg/infra/httpclient/httpclientprovider"
	"github.com/grafana/grafana/pkg/infra/kvstore"
	"github.com/grafana/grafana/pkg/infra/localcache"
	metrics2 "github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/infra/remotecache"
	"github.com/grafana/grafana/pkg/infra/serverlock"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/infra/usagestats"
	"github.com/grafana/grafana/pkg/infra/usagestats/service"
	"github.com/grafana/grafana/pkg/login/social"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/manager"
	"github.com/grafana/grafana/pkg/plugins/manager/loader"
	"github.com/grafana/grafana/pkg/plugins/manager/signature"
	"github.com/grafana/grafana/pkg/plugins/plugincontext"
	"github.com/grafana/grafana/pkg/plugins/plugindashboards"
	"github.com/grafana/grafana/pkg/server/backgroundsvcs"
	"github.com/grafana/grafana/pkg/services/accesscontrol/ossaccesscontrol"
	"github.com/grafana/grafana/pkg/services/alerting"
	"github.com/grafana/grafana/pkg/services/auth"
	"github.com/grafana/grafana/pkg/services/auth/jwt"
	"github.com/grafana/grafana/pkg/services/cleanup"
	"github.com/grafana/grafana/pkg/services/contexthandler"
	"github.com/grafana/grafana/pkg/services/dashboardsnapshots"
	"github.com/grafana/grafana/pkg/services/datasourceproxy"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/services/encryption/ossencryption"
	"github.com/grafana/grafana/pkg/services/hooks"
	"github.com/grafana/grafana/pkg/services/kmsproviders/osskmsproviders"
	"github.com/grafana/grafana/pkg/services/libraryelements"
	"github.com/grafana/grafana/pkg/services/librarypanels"
	"github.com/grafana/grafana/pkg/services/licensing"
	"github.com/grafana/grafana/pkg/services/live"
	"github.com/grafana/grafana/pkg/services/live/pushhttp"
	"github.com/grafana/grafana/pkg/services/login"
	"github.com/grafana/grafana/pkg/services/login/authinfoservice"
	"github.com/grafana/grafana/pkg/services/login/loginservice"
	"github.com/grafana/grafana/pkg/services/ngalert"
	"github.com/grafana/grafana/pkg/services/ngalert/metrics"
	"github.com/grafana/grafana/pkg/services/notifications"
	"github.com/grafana/grafana/pkg/services/oauthtoken"
	"github.com/grafana/grafana/pkg/services/pluginsettings"
	"github.com/grafana/grafana/pkg/services/provisioning"
	"github.com/grafana/grafana/pkg/services/quota"
	"github.com/grafana/grafana/pkg/services/rendering"
	"github.com/grafana/grafana/pkg/services/schemaloader"
	"github.com/grafana/grafana/pkg/services/search"
	"github.com/grafana/grafana/pkg/services/searchusers"
	"github.com/grafana/grafana/pkg/services/searchusers/filters"
	"github.com/grafana/grafana/pkg/services/secrets"
	"github.com/grafana/grafana/pkg/services/secrets/database"
	manager2 "github.com/grafana/grafana/pkg/services/secrets/manager"
	"github.com/grafana/grafana/pkg/services/serviceaccounts"
	manager3 "github.com/grafana/grafana/pkg/services/serviceaccounts/manager"
	"github.com/grafana/grafana/pkg/services/shorturls"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrations"
	"github.com/grafana/grafana/pkg/services/updatechecker"
	"github.com/grafana/grafana/pkg/services/validations"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/tsdb/azuremonitor"
	"github.com/grafana/grafana/pkg/tsdb/cloudmonitoring"
	"github.com/grafana/grafana/pkg/tsdb/cloudwatch"
	"github.com/grafana/grafana/pkg/tsdb/elasticsearch"
	"github.com/grafana/grafana/pkg/tsdb/grafanads"
	"github.com/grafana/grafana/pkg/tsdb/graphite"
	"github.com/grafana/grafana/pkg/tsdb/influxdb"
	"github.com/grafana/grafana/pkg/tsdb/legacydata"
	service2 "github.com/grafana/grafana/pkg/tsdb/legacydata/service"
	"github.com/grafana/grafana/pkg/tsdb/loki"
	"github.com/grafana/grafana/pkg/tsdb/mssql"
	"github.com/grafana/grafana/pkg/tsdb/mysql"
	"github.com/grafana/grafana/pkg/tsdb/opentsdb"
	"github.com/grafana/grafana/pkg/tsdb/postgres"
	"github.com/grafana/grafana/pkg/tsdb/prometheus"
	"github.com/grafana/grafana/pkg/tsdb/tempo"
	"github.com/grafana/grafana/pkg/tsdb/testdatasource"
)

import (
	_ "github.com/grafana/grafana/pkg/extensions"
)

// Injectors from wire.go:

func Initialize(cla setting.CommandLineArgs, opts Options, apiOpts api.ServerOptions) (*Server, error) {
	cfg, err := setting.NewCfgFromArgs(cla)
	if err != nil {
		return nil, err
	}
	routeRegisterImpl := routing.ProvideRegister(cfg)
	inProcBus := bus.ProvideBus()
	cacheService := localcache.ProvideService()
	ossMigrations := migrations.ProvideOSSMigrations()
	sqlStore, err := sqlstore.ProvideService(cfg, cacheService, inProcBus, ossMigrations)
	if err != nil {
		return nil, err
	}
	remoteCache, err := remotecache.ProvideService(cfg, sqlStore)
	if err != nil {
		return nil, err
	}
	ossPluginRequestValidator := validations.ProvideValidator()
	hooksService := hooks.ProvideService()
	ossLicensingService := licensing.ProvideService(cfg, hooksService)
	unsignedPluginAuthorizer, err := signature.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	loaderLoader, err := loader.ProvideService(ossLicensingService, cfg, unsignedPluginAuthorizer)
	if err != nil {
		return nil, err
	}
	pluginManager, err := manager.ProvideService(cfg, ossPluginRequestValidator, loaderLoader, sqlStore)
	if err != nil {
		return nil, err
	}
	renderingService, err := rendering.ProvideService(cfg, remoteCache, pluginManager)
	if err != nil {
		return nil, err
	}
	socialService := social.ProvideService(cfg)
	oauthtokenService := oauthtoken.ProvideService(socialService)
	secretsStoreImpl := database.ProvideSecretsStore(sqlStore)
	ossencryptionService := ossencryption.ProvideService()
	ossImpl := setting.ProvideProvider(cfg)
	osskmsprovidersService := osskmsproviders.ProvideService(ossencryptionService, ossImpl)
	kvStore := kvstore.ProvideService(sqlStore)
	usageStats := service.ProvideService(cfg, inProcBus, sqlStore, pluginManager, socialService, kvStore)
	secretsService, err := manager2.ProvideSecretsService(secretsStoreImpl, osskmsprovidersService, ossencryptionService, ossImpl, usageStats)
	if err != nil {
		return nil, err
	}
	datasourcesService := datasources.ProvideService(inProcBus, sqlStore, secretsService)
	serviceService := service2.ProvideService(pluginManager, oauthtokenService, datasourcesService)
	alertEngine := alerting.ProvideAlertEngine(renderingService, inProcBus, ossPluginRequestValidator, serviceService, usageStats, ossencryptionService, cfg)
	cacheServiceImpl := datasources.ProvideCacheService(cacheService, sqlStore)
	serverLockService := serverlock.ProvideService(sqlStore)
	userAuthTokenService := auth.ProvideUserAuthTokenService(sqlStore, serverLockService, cfg)
	shortURLService := shorturls.ProvideService(sqlStore)
	cleanUpService := cleanup.ProvideService(cfg, serverLockService, shortURLService)
	provisioningServiceImpl, err := provisioning.ProvideService(cfg, sqlStore, pluginManager, ossencryptionService)
	if err != nil {
		return nil, err
	}
	quotaService := quota.ProvideService(cfg, userAuthTokenService)
	ossUserProtectionImpl := authinfoservice.ProvideOSSUserProtectionService()
	implementation := authinfoservice.ProvideAuthInfoService(inProcBus, sqlStore, ossUserProtectionImpl, secretsService)
	loginserviceImplementation := loginservice.ProvideService(sqlStore, inProcBus, quotaService, implementation)
	ossAccessControlService := ossaccesscontrol.ProvideService(cfg, usageStats)
	provider := httpclientprovider.New(cfg)
	dataSourceProxyService := datasourceproxy.ProvideService(cacheServiceImpl, ossPluginRequestValidator, pluginManager, cfg, provider, oauthtokenService, datasourcesService)
	searchService := search.ProvideService(cfg, inProcBus)
	pluginsettingsService := pluginsettings.ProvideService(inProcBus, sqlStore, secretsService)
	plugincontextProvider := plugincontext.ProvideService(inProcBus, cacheService, pluginManager, cacheServiceImpl, secretsService, pluginsettingsService)
	logsService := cloudwatch.ProvideLogsService()
	grafanaLive, err := live.ProvideService(plugincontextProvider, cfg, routeRegisterImpl, logsService, pluginManager, cacheService, cacheServiceImpl, sqlStore, secretsService, usageStats)
	if err != nil {
		return nil, err
	}
	gateway := pushhttp.ProvideService(cfg, grafanaLive)
	authService, err := jwt.ProvideService(cfg, remoteCache)
	if err != nil {
		return nil, err
	}
	contextHandler := contexthandler.ProvideService(cfg, userAuthTokenService, authService, remoteCache, renderingService, sqlStore)
	schemaLoaderService, err := schemaloader.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	exprService := expr.ProvideService(cfg, pluginManager, secretsService)
	ngAlert := metrics.ProvideService()
	alertNG, err := ngalert.ProvideService(cfg, cacheServiceImpl, routeRegisterImpl, sqlStore, kvStore, exprService, dataSourceProxyService, quotaService, secretsService, ngAlert)
	if err != nil {
		return nil, err
	}
	libraryElementService := libraryelements.ProvideService(cfg, sqlStore, routeRegisterImpl)
	libraryPanelService := librarypanels.ProvideService(cfg, sqlStore, routeRegisterImpl, libraryElementService)
	notificationService, err := notifications.ProvideService(inProcBus, cfg)
	if err != nil {
		return nil, err
	}
	tracingService, err := tracing.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	internalMetricsService, err := metrics2.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	updatecheckerService := updatechecker.ProvideService(cfg)
	ossSearchUserFilter := filters.ProvideOSSSearchUserFilter()
	ossService := searchusers.ProvideUsersService(inProcBus, ossSearchUserFilter)
	httpServer, err := api.ProvideHTTPServer(apiOpts, cfg, routeRegisterImpl, inProcBus, renderingService, ossLicensingService, hooksService, cacheService, sqlStore, serviceService, alertEngine, ossPluginRequestValidator, pluginManager, pluginManager, pluginManager, pluginManager, loaderLoader, ossImpl, cacheServiceImpl, userAuthTokenService, cleanUpService, shortURLService, remoteCache, provisioningServiceImpl, loginserviceImplementation, ossAccessControlService, dataSourceProxyService, searchService, grafanaLive, gateway, plugincontextProvider, contextHandler, schemaLoaderService, alertNG, libraryPanelService, libraryElementService, notificationService, tracingService, internalMetricsService, quotaService, socialService, oauthtokenService, ossencryptionService, updatecheckerService, ossService, datasourcesService, secretsService, exprService)
	if err != nil {
		return nil, err
	}
	azuremonitorService := azuremonitor.ProvideService(cfg, provider, pluginManager)
	cloudWatchService, err := cloudwatch.ProvideService(cfg, logsService, pluginManager)
	if err != nil {
		return nil, err
	}
	elasticsearchService, err := elasticsearch.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	graphiteService, err := graphite.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	influxdbService, err := influxdb.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	lokiService, err := loki.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	opentsdbService, err := opentsdb.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	prometheusService, err := prometheus.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	tempoService, err := tempo.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	testdatasourceService, err := testdatasource.ProvideService(cfg, pluginManager)
	if err != nil {
		return nil, err
	}
	plugindashboardsService := plugindashboards.ProvideService(pluginManager, pluginManager, sqlStore)
	dashboardsnapshotsService := dashboardsnapshots.ProvideService(inProcBus, sqlStore, secretsService)
	postgresService, err := postgres.ProvideService(cfg, pluginManager)
	if err != nil {
		return nil, err
	}
	mysqlService, err := mysql.ProvideService(cfg, pluginManager, provider)
	if err != nil {
		return nil, err
	}
	mssqlService, err := mssql.ProvideService(cfg, pluginManager)
	if err != nil {
		return nil, err
	}
	grafanadsService := grafanads.ProvideService(cfg, pluginManager)
	cloudmonitoringService := cloudmonitoring.ProvideService(cfg, provider, pluginManager, datasourcesService)
	alertNotificationService := alerting.ProvideService(inProcBus, sqlStore, ossencryptionService)
	serviceAccountsService, err := manager3.ProvideServiceAccountsService(cfg, sqlStore, ossAccessControlService, routeRegisterImpl)
	if err != nil {
		return nil, err
	}
	backgroundServiceRegistry := backgroundsvcs.ProvideBackgroundServiceRegistry(httpServer, alertNG, cleanUpService, grafanaLive, gateway, notificationService, renderingService, userAuthTokenService, provisioningServiceImpl, alertEngine, pluginManager, internalMetricsService, usageStats, updatecheckerService, tracingService, remoteCache, azuremonitorService, cloudWatchService, elasticsearchService, graphiteService, influxdbService, lokiService, opentsdbService, prometheusService, tempoService, testdatasourceService, plugindashboardsService, dashboardsnapshotsService, postgresService, mysqlService, mssqlService, grafanadsService, cloudmonitoringService, pluginsettingsService, alertNotificationService, serviceAccountsService)
	server, err := New(opts, cfg, httpServer, ossAccessControlService, provisioningServiceImpl, backgroundServiceRegistry)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func InitializeForTest(cla setting.CommandLineArgs, opts Options, apiOpts api.ServerOptions) (*TestEnv, error) {
	cfg, err := setting.NewCfgFromArgs(cla)
	if err != nil {
		return nil, err
	}
	routeRegisterImpl := routing.ProvideRegister(cfg)
	inProcBus := bus.ProvideBus()
	ossMigrations := migrations.ProvideOSSMigrations()
	sqlStore, err := sqlstore.ProvideServiceForTests(ossMigrations)
	if err != nil {
		return nil, err
	}
	remoteCache, err := remotecache.ProvideService(cfg, sqlStore)
	if err != nil {
		return nil, err
	}
	ossPluginRequestValidator := validations.ProvideValidator()
	hooksService := hooks.ProvideService()
	ossLicensingService := licensing.ProvideService(cfg, hooksService)
	unsignedPluginAuthorizer, err := signature.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	loaderLoader, err := loader.ProvideService(ossLicensingService, cfg, unsignedPluginAuthorizer)
	if err != nil {
		return nil, err
	}
	pluginManager, err := manager.ProvideService(cfg, ossPluginRequestValidator, loaderLoader, sqlStore)
	if err != nil {
		return nil, err
	}
	renderingService, err := rendering.ProvideService(cfg, remoteCache, pluginManager)
	if err != nil {
		return nil, err
	}
	cacheService := localcache.ProvideService()
	socialService := social.ProvideService(cfg)
	oauthtokenService := oauthtoken.ProvideService(socialService)
	secretsStoreImpl := database.ProvideSecretsStore(sqlStore)
	ossencryptionService := ossencryption.ProvideService()
	ossImpl := setting.ProvideProvider(cfg)
	osskmsprovidersService := osskmsproviders.ProvideService(ossencryptionService, ossImpl)
	kvStore := kvstore.ProvideService(sqlStore)
	usageStats := service.ProvideService(cfg, inProcBus, sqlStore, pluginManager, socialService, kvStore)
	secretsService, err := manager2.ProvideSecretsService(secretsStoreImpl, osskmsprovidersService, ossencryptionService, ossImpl, usageStats)
	if err != nil {
		return nil, err
	}
	datasourcesService := datasources.ProvideService(inProcBus, sqlStore, secretsService)
	serviceService := service2.ProvideService(pluginManager, oauthtokenService, datasourcesService)
	alertEngine := alerting.ProvideAlertEngine(renderingService, inProcBus, ossPluginRequestValidator, serviceService, usageStats, ossencryptionService, cfg)
	cacheServiceImpl := datasources.ProvideCacheService(cacheService, sqlStore)
	serverLockService := serverlock.ProvideService(sqlStore)
	userAuthTokenService := auth.ProvideUserAuthTokenService(sqlStore, serverLockService, cfg)
	shortURLService := shorturls.ProvideService(sqlStore)
	cleanUpService := cleanup.ProvideService(cfg, serverLockService, shortURLService)
	provisioningServiceImpl, err := provisioning.ProvideService(cfg, sqlStore, pluginManager, ossencryptionService)
	if err != nil {
		return nil, err
	}
	quotaService := quota.ProvideService(cfg, userAuthTokenService)
	ossUserProtectionImpl := authinfoservice.ProvideOSSUserProtectionService()
	implementation := authinfoservice.ProvideAuthInfoService(inProcBus, sqlStore, ossUserProtectionImpl, secretsService)
	loginserviceImplementation := loginservice.ProvideService(sqlStore, inProcBus, quotaService, implementation)
	ossAccessControlService := ossaccesscontrol.ProvideService(cfg, usageStats)
	provider := httpclientprovider.New(cfg)
	dataSourceProxyService := datasourceproxy.ProvideService(cacheServiceImpl, ossPluginRequestValidator, pluginManager, cfg, provider, oauthtokenService, datasourcesService)
	searchService := search.ProvideService(cfg, inProcBus)
	pluginsettingsService := pluginsettings.ProvideService(inProcBus, sqlStore, secretsService)
	plugincontextProvider := plugincontext.ProvideService(inProcBus, cacheService, pluginManager, cacheServiceImpl, secretsService, pluginsettingsService)
	logsService := cloudwatch.ProvideLogsService()
	grafanaLive, err := live.ProvideService(plugincontextProvider, cfg, routeRegisterImpl, logsService, pluginManager, cacheService, cacheServiceImpl, sqlStore, secretsService, usageStats)
	if err != nil {
		return nil, err
	}
	gateway := pushhttp.ProvideService(cfg, grafanaLive)
	authService, err := jwt.ProvideService(cfg, remoteCache)
	if err != nil {
		return nil, err
	}
	contextHandler := contexthandler.ProvideService(cfg, userAuthTokenService, authService, remoteCache, renderingService, sqlStore)
	schemaLoaderService, err := schemaloader.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	exprService := expr.ProvideService(cfg, pluginManager, secretsService)
	ngAlert := metrics.ProvideServiceForTest()
	alertNG, err := ngalert.ProvideService(cfg, cacheServiceImpl, routeRegisterImpl, sqlStore, kvStore, exprService, dataSourceProxyService, quotaService, secretsService, ngAlert)
	if err != nil {
		return nil, err
	}
	libraryElementService := libraryelements.ProvideService(cfg, sqlStore, routeRegisterImpl)
	libraryPanelService := librarypanels.ProvideService(cfg, sqlStore, routeRegisterImpl, libraryElementService)
	notificationService, err := notifications.ProvideService(inProcBus, cfg)
	if err != nil {
		return nil, err
	}
	tracingService, err := tracing.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	internalMetricsService, err := metrics2.ProvideService(cfg)
	if err != nil {
		return nil, err
	}
	updatecheckerService := updatechecker.ProvideService(cfg)
	ossSearchUserFilter := filters.ProvideOSSSearchUserFilter()
	ossService := searchusers.ProvideUsersService(inProcBus, ossSearchUserFilter)
	httpServer, err := api.ProvideHTTPServer(apiOpts, cfg, routeRegisterImpl, inProcBus, renderingService, ossLicensingService, hooksService, cacheService, sqlStore, serviceService, alertEngine, ossPluginRequestValidator, pluginManager, pluginManager, pluginManager, pluginManager, loaderLoader, ossImpl, cacheServiceImpl, userAuthTokenService, cleanUpService, shortURLService, remoteCache, provisioningServiceImpl, loginserviceImplementation, ossAccessControlService, dataSourceProxyService, searchService, grafanaLive, gateway, plugincontextProvider, contextHandler, schemaLoaderService, alertNG, libraryPanelService, libraryElementService, notificationService, tracingService, internalMetricsService, quotaService, socialService, oauthtokenService, ossencryptionService, updatecheckerService, ossService, datasourcesService, secretsService, exprService)
	if err != nil {
		return nil, err
	}
	azuremonitorService := azuremonitor.ProvideService(cfg, provider, pluginManager)
	cloudWatchService, err := cloudwatch.ProvideService(cfg, logsService, pluginManager)
	if err != nil {
		return nil, err
	}
	elasticsearchService, err := elasticsearch.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	graphiteService, err := graphite.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	influxdbService, err := influxdb.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	lokiService, err := loki.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	opentsdbService, err := opentsdb.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	prometheusService, err := prometheus.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	tempoService, err := tempo.ProvideService(provider, pluginManager)
	if err != nil {
		return nil, err
	}
	testdatasourceService, err := testdatasource.ProvideService(cfg, pluginManager)
	if err != nil {
		return nil, err
	}
	plugindashboardsService := plugindashboards.ProvideService(pluginManager, pluginManager, sqlStore)
	dashboardsnapshotsService := dashboardsnapshots.ProvideService(inProcBus, sqlStore, secretsService)
	postgresService, err := postgres.ProvideService(cfg, pluginManager)
	if err != nil {
		return nil, err
	}
	mysqlService, err := mysql.ProvideService(cfg, pluginManager, provider)
	if err != nil {
		return nil, err
	}
	mssqlService, err := mssql.ProvideService(cfg, pluginManager)
	if err != nil {
		return nil, err
	}
	grafanadsService := grafanads.ProvideService(cfg, pluginManager)
	cloudmonitoringService := cloudmonitoring.ProvideService(cfg, provider, pluginManager, datasourcesService)
	alertNotificationService := alerting.ProvideService(inProcBus, sqlStore, ossencryptionService)
	serviceAccountsService, err := manager3.ProvideServiceAccountsService(cfg, sqlStore, ossAccessControlService, routeRegisterImpl)
	if err != nil {
		return nil, err
	}
	backgroundServiceRegistry := backgroundsvcs.ProvideBackgroundServiceRegistry(httpServer, alertNG, cleanUpService, grafanaLive, gateway, notificationService, renderingService, userAuthTokenService, provisioningServiceImpl, alertEngine, pluginManager, internalMetricsService, usageStats, updatecheckerService, tracingService, remoteCache, azuremonitorService, cloudWatchService, elasticsearchService, graphiteService, influxdbService, lokiService, opentsdbService, prometheusService, tempoService, testdatasourceService, plugindashboardsService, dashboardsnapshotsService, postgresService, mysqlService, mssqlService, grafanadsService, cloudmonitoringService, pluginsettingsService, alertNotificationService, serviceAccountsService)
	server, err := New(opts, cfg, httpServer, ossAccessControlService, provisioningServiceImpl, backgroundServiceRegistry)
	if err != nil {
		return nil, err
	}
	testEnv, err := ProvideTestEnv(server, sqlStore)
	if err != nil {
		return nil, err
	}
	return testEnv, nil
}

// wire.go:

var wireBasicSet = wire.NewSet(service2.ProvideService, wire.Bind(new(legacydata.RequestHandler), new(*service2.Service)), alerting.ProvideAlertEngine, wire.Bind(new(alerting.UsageStatsQuerier), new(*alerting.AlertEngine)), setting.NewCfgFromArgs, New, api.ProvideHTTPServer, bus.ProvideBus, wire.Bind(new(bus.Bus), new(*bus.InProcBus)), rendering.ProvideService, wire.Bind(new(rendering.Service), new(*rendering.RenderingService)), routing.ProvideRegister, wire.Bind(new(routing.RouteRegister), new(*routing.RouteRegisterImpl)), hooks.ProvideService, kvstore.ProvideService, localcache.ProvideService, updatechecker.ProvideService, service.ProvideService, wire.Bind(new(usagestats.Service), new(*service.UsageStats)), manager.ProvideService, wire.Bind(new(plugins.Client), new(*manager.PluginManager)), wire.Bind(new(plugins.Store), new(*manager.PluginManager)), wire.Bind(new(plugins.CoreBackendRegistrar), new(*manager.PluginManager)), wire.Bind(new(plugins.StaticRouteResolver), new(*manager.PluginManager)), wire.Bind(new(plugins.PluginDashboardManager), new(*manager.PluginManager)), wire.Bind(new(plugins.RendererManager), new(*manager.PluginManager)), loader.ProvideService, wire.Bind(new(plugins.Loader), new(*loader.Loader)), wire.Bind(new(plugins.ErrorResolver), new(*loader.Loader)), cloudwatch.ProvideService, cloudwatch.ProvideLogsService, cloudmonitoring.ProvideService, azuremonitor.ProvideService, postgres.ProvideService, mysql.ProvideService, mssql.ProvideService, httpclientprovider.New, wire.Bind(new(httpclient.Provider), new(*httpclient2.Provider)), serverlock.ProvideService, cleanup.ProvideService, shorturls.ProvideService, wire.Bind(new(shorturls.Service), new(*shorturls.ShortURLService)), quota.ProvideService, remotecache.ProvideService, loginservice.ProvideService, wire.Bind(new(login.Service), new(*loginservice.Implementation)), authinfoservice.ProvideAuthInfoService, wire.Bind(new(login.AuthInfoService), new(*authinfoservice.Implementation)), datasourceproxy.ProvideService, search.ProvideService, live.ProvideService, pushhttp.ProvideService, plugincontext.ProvideService, contexthandler.ProvideService, jwt.ProvideService, wire.Bind(new(models.JWTService), new(*jwt.AuthService)), plugindashboards.ProvideService, schemaloader.ProvideService, ngalert.ProvideService, librarypanels.ProvideService, wire.Bind(new(librarypanels.Service), new(*librarypanels.LibraryPanelService)), libraryelements.ProvideService, wire.Bind(new(libraryelements.Service), new(*libraryelements.LibraryElementService)), notifications.ProvideService, tracing.ProvideService, metrics2.ProvideService, testdatasource.ProvideService, opentsdb.ProvideService, social.ProvideService, influxdb.ProvideService, wire.Bind(new(social.Service), new(*social.SocialService)), oauthtoken.ProvideService, wire.Bind(new(oauthtoken.OAuthTokenService), new(*oauthtoken.Service)), tempo.ProvideService, loki.ProvideService, graphite.ProvideService, prometheus.ProvideService, elasticsearch.ProvideService, manager2.ProvideSecretsService, wire.Bind(new(secrets.Service), new(*manager2.SecretsService)), database.ProvideSecretsStore, wire.Bind(new(secrets.Store), new(*database.SecretsStoreImpl)), grafanads.ProvideService, dashboardsnapshots.ProvideService, datasources.ProvideService, pluginsettings.ProvideService, alerting.ProvideService, manager3.ProvideServiceAccountsService, wire.Bind(new(serviceaccounts.Service), new(*manager3.ServiceAccountsService)), expr.ProvideService)

var wireSet = wire.NewSet(
	wireBasicSet, sqlstore.ProvideService, metrics.ProvideService,
)

var wireTestSet = wire.NewSet(
	wireBasicSet,
	ProvideTestEnv, sqlstore.ProvideServiceForTests, metrics.ProvideServiceForTest,
)

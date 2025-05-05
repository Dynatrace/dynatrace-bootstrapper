package ruxit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	expectedConf = `[agentType]
apache on
dotnet not-global
go not-global
iis on
java on
loganalytics on
network on
nginx on
nodejs on
opentracing off
php on
plugin on
python off
sdk on
varnish off
wsmb on

[general]
addContainerImageNametoPG on
addContainerNameToPGI off
addNodejsScriptNameToPGI on
bpmInjection on
cassandraClusterNameInPG on
containerInjectionRules -1::EXCLUDE:EQUALS,KUBERNETES_CONTAINERNAME,POD;-2::EXCLUDE:CONTAINS,IMAGE_NAME,pause-amd64;-3::EXCLUDE:EQUALS,KUBERNETES_NAMESPACE,openshift-sdn;-4::EXCLUDE:ENDS,KUBERNETES_FULLPODNAME,-build
containerdInjection on
coreclrInjection on
crioInjection on
disableAgentTypeBasedInjection off
disableJBossServerNameProperty off
disableSpringBootGroupCalc off
dockerInjection on
dockerWindowsInjection off
enableEquinoxGroupCalc on
enableExecHook on
enableNodeJsAgentEnvFile off
enableNodeJsAgentEsmLoaders off
enableNodeJsAgentPreloading on
enableNodeJsMultiversionLibrary off
enableOsAgentDefaultIdCalc off
enablePhpCliServerInstrumentation off
enablePodmanInjection off
enableTibcoBWContainerEditionGroupCalc on
enableTipcoBWGroupCalc on
enableWebSphereLibertyGroupCalc off
envoyInjection on
expandJavaAtFiles off
fixDockerContainerAndImageNameInPGI off
fullStackJavaInMBProcesses on
injectionRules -1::EXCLUDE:CONTAINS,PHP_CLI_SCRIPT_PATH,;-2::EXCLUDE:EQUALS,EXE_NAME,php-cgi;-3::INCLUDE:CONTAINS,ASPNETCORE_APPL_PATH,;-4::INCLUDE:EQUALS,EXE_NAME,w3wp.exe;-49::EXCLUDE:EQUALS,EXE_NAME,filebeat;-50::EXCLUDE:EQUALS,EXE_NAME,metricbeat;-51::EXCLUDE:EQUALS,EXE_NAME,packetbeat;-52::EXCLUDE:EQUALS,EXE_NAME,auditbeat;-53::EXCLUDE:EQUALS,EXE_NAME,heartbeat;-54::EXCLUDE:EQUALS,EXE_NAME,functionbeat;-72::EXCLUDE:EQUALS,EXE_NAME,yq;-47::EXCLUDE:EQUALS,LINKAGE,static;-5::INCLUDE:EQUALS,EXE_NAME,caddy;-6::INCLUDE:EQUALS,EXE_NAME,influxd;-7::INCLUDE:EQUALS,EXE_NAME,adapter;-8::INCLUDE:EQUALS,EXE_NAME,auctioneer;-9::INCLUDE:EQUALS,EXE_NAME,bbs;-10::INCLUDE:EQUALS,EXE_NAME,cc-uploader;-11::INCLUDE:EQUALS,EXE_NAME,doppler;-12::INCLUDE:EQUALS,EXE_NAME,gorouter;-13::INCLUDE:EQUALS,EXE_NAME,locket;-14::INCLUDE:EQUALS,EXE_NAME,metron;-16::INCLUDE:EQUALS,EXE_NAME,rep;-17::INCLUDE:EQUALS,EXE_NAME,route-emitter;-18::INCLUDE:EQUALS,EXE_NAME,route-registrar;-19::INCLUDE:EQUALS,EXE_NAME,routing-api;-20::INCLUDE:EQUALS,EXE_NAME,scheduler;-21::INCLUDE:EQUALS,EXE_NAME,silk-daemon;-22::INCLUDE:EQUALS,EXE_NAME,switchboard;-23::INCLUDE:EQUALS,EXE_NAME,syslog_drain_binder;-24::INCLUDE:EQUALS,EXE_NAME,tps-watcher;-25::INCLUDE:EQUALS,EXE_NAME,trafficcontroller;-26::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/prebuild-install;-27::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/npm;-28::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/grunt;-29::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/typescript;-45::EXCLUDE:EQUALS,NODEJS_APP_NAME,yarn;-68::EXCLUDE:EQUALS,NODEJS_APP_NAME,corepack;-32::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/node-pre-gyp;-33::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/node-gyp;-34::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/gulp-cli;-35::EXCLUDE:EQUALS,NODEJS_SCRIPT_NAME,bin/pm2;-36::EXCLUDE:STARTS,CLOUD_FOUNDRY_APP_NAME,apps-manager-js;-55::EXCLUDE:EQUALS,EXE_NAME,grootfs;-56::EXCLUDE:EQUALS,EXE_NAME,tardis;-43::EXCLUDE:STARTS,EXE_PATH,/tmp/buildpacks/;-37::INCLUDE:CONTAINS,CLOUD_FOUNDRY_APP_NAME,;-38::EXCLUDE:EQUALS,KUBERNETES_CONTAINERNAME,POD;-39::EXCLUDE:CONTAINS,CONTAINER_IMAGE_NAME,pause-amd64;-44::EXCLUDE:EQUALS,EXE_NAME,oc;-58::EXCLUDE:EQUALS,EXE_NAME,calico-node;-67::EXCLUDE:EQUALS,EXE_PATH,/usr/bin/piper;-69::EXCLUDE:EQUALS,KUBERNETES_CONTAINERNAME,cassandra-operator;-70::EXCLUDE:CONTAINS,EXE_NAME,UiPath;-71::EXCLUDE:EQUALS,EXE_NAME,openhandlecollector.exe;-40::INCLUDE:CONTAINS,KUBERNETES_NAMESPACE,;-41::INCLUDE:CONTAINS,CONTAINER_NAME,;-46::EXCLUDE:EQUALS,EXE_PATH,/opt/cni/bin/host-local;-48::EXCLUDE:STARTS,EXE_NAME,mqsi;-57::EXCLUDE:STARTS,JAVA_JAR_FILE,org.eclipse.equinox.launcher;-59::EXCLUDE:EQUALS,EXE_NAME,casclient.exe;-60::EXCLUDE:EQUALS,JAVA_JAR_FILE,dynatrace_ibm_mq_connector.jar;-61::EXCLUDE:CONTAINS,EXE_NAME,Agent.Worker;-62::EXCLUDE:CONTAINS,DOTNET_COMMAND,Agent.Worker;-63::EXCLUDE:CONTAINS,EXE_NAME,Agent.Listener;-64::EXCLUDE:CONTAINS,DOTNET_COMMAND,Agent.Listener;-65::EXCLUDE:EQUALS,EXE_NAME,FlexNetJobExecutorService;-66::EXCLUDE:EQUALS,EXE_NAME,FlexNetMaintenanceRemotingService
nodejsAgentDir nodeagent
optimizedSuspendThreads on
php74Injection on
php7Injection on
php80InjectionEA on
php81Injection on
pythonInjection off
removeContainerIDfromPGI on
removeIdsFromPaths on
runcInjection on
serverAddress {https://example1.dev.dynatracelabs.com/communication;https://example2.dev.dynatracelabs.com/communication;https://example3.dev.dynatracelabs.com/communication;https://example3.dev.dynatracelabs.com/communication;https://example4.dev.dynatracelabs.com/communication;https://exampl5.dev.dynatracelabs.com/communication;https://example6.dev.dynatracelabs.com:443
staticGoInjection off
stripIdsFromKubernetesNamespace on
stripVersionFromImageName on
switchToPhpAgentNG on
tenant zib50933
tenantToken this-is-secret
trustedTimestampVerificationJavascriptAgent off
websphereClusterNameInPG on
wincInjection off

`

	expectedJson = `{"properties":[{"section":"general","key":"websphereClusterNameInPG","value":"on"},{"section":"agentType","key":"nginx","value":"on"},{"section":"general","key":"enableExecHook","value":"on"},{"section":"general","key":"php7Injection","value":"on"},{"section":"general","key":"removeIdsFromPaths","value":"on"},{"section":"general","key":"enableNodeJsAgentEnvFile","value":"off"},{"section":"general","key":"removeContainerIDfromPGI","value":"on"},{"section":"general","key":"tenantToken","value":"this-is-secret"},{"section":"general","key":"addContainerImageNametoPG","value":"on"},{"section":"general","key":"disableJBossServerNameProperty","value":"off"},{"section":"general","key":"bpmInjection","value":"on"},{"section":"general","key":"fixDockerContainerAndImageNameInPGI","value":"off"},{"section":"agentType","key":"wsmb","value":"on"},{"section":"general","key":"nodejsAgentDir","value":"nodeagent"},{"section":"general","key":"serverAddress","value":"{https://example1.dev.dynatracelabs.com/communication;https://example2.dev.dynatracelabs.com/communication;https://example3.dev.dynatracelabs.com/communication;https://example3.dev.dynatracelabs.com/communication;https://example4.dev.dynatracelabs.com/communication;https://exampl5.dev.dynatracelabs.com/communication;https://example6.dev.dynatracelabs.com:443"},{"section":"agentType","key":"opentracing","value":"off"},{"section":"general","key":"enablePhpCliServerInstrumentation","value":"off"},{"section":"general","key":"php74Injection","value":"on"},{"section":"general","key":"wincInjection","value":"off"},{"section":"agentType","key":"php","value":"on"},{"section":"general","key":"staticGoInjection","value":"off"},{"section":"general","key":"enableTipcoBWGroupCalc","value":"on"},{"section":"general","key":"injectionRules","value":"-1::EXCLUDE:CONTAINS,PHP_CLI_SCRIPT_PATH,;-2::EXCLUDE:EQUALS,EXE_NAME,php-cgi;-3::INCLUDE:CONTAINS,ASPNETCORE_APPL_PATH,;-4::INCLUDE:EQUALS,EXE_NAME,w3wp.exe;-49::EXCLUDE:EQUALS,EXE_NAME,filebeat;-50::EXCLUDE:EQUALS,EXE_NAME,metricbeat;-51::EXCLUDE:EQUALS,EXE_NAME,packetbeat;-52::EXCLUDE:EQUALS,EXE_NAME,auditbeat;-53::EXCLUDE:EQUALS,EXE_NAME,heartbeat;-54::EXCLUDE:EQUALS,EXE_NAME,functionbeat;-72::EXCLUDE:EQUALS,EXE_NAME,yq;-47::EXCLUDE:EQUALS,LINKAGE,static;-5::INCLUDE:EQUALS,EXE_NAME,caddy;-6::INCLUDE:EQUALS,EXE_NAME,influxd;-7::INCLUDE:EQUALS,EXE_NAME,adapter;-8::INCLUDE:EQUALS,EXE_NAME,auctioneer;-9::INCLUDE:EQUALS,EXE_NAME,bbs;-10::INCLUDE:EQUALS,EXE_NAME,cc-uploader;-11::INCLUDE:EQUALS,EXE_NAME,doppler;-12::INCLUDE:EQUALS,EXE_NAME,gorouter;-13::INCLUDE:EQUALS,EXE_NAME,locket;-14::INCLUDE:EQUALS,EXE_NAME,metron;-16::INCLUDE:EQUALS,EXE_NAME,rep;-17::INCLUDE:EQUALS,EXE_NAME,route-emitter;-18::INCLUDE:EQUALS,EXE_NAME,route-registrar;-19::INCLUDE:EQUALS,EXE_NAME,routing-api;-20::INCLUDE:EQUALS,EXE_NAME,scheduler;-21::INCLUDE:EQUALS,EXE_NAME,silk-daemon;-22::INCLUDE:EQUALS,EXE_NAME,switchboard;-23::INCLUDE:EQUALS,EXE_NAME,syslog_drain_binder;-24::INCLUDE:EQUALS,EXE_NAME,tps-watcher;-25::INCLUDE:EQUALS,EXE_NAME,trafficcontroller;-26::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/prebuild-install;-27::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/npm;-28::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/grunt;-29::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/typescript;-45::EXCLUDE:EQUALS,NODEJS_APP_NAME,yarn;-68::EXCLUDE:EQUALS,NODEJS_APP_NAME,corepack;-32::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/node-pre-gyp;-33::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/node-gyp;-34::EXCLUDE:ENDS,NODEJS_APP_BASE_DIR,/node_modules/gulp-cli;-35::EXCLUDE:EQUALS,NODEJS_SCRIPT_NAME,bin/pm2;-36::EXCLUDE:STARTS,CLOUD_FOUNDRY_APP_NAME,apps-manager-js;-55::EXCLUDE:EQUALS,EXE_NAME,grootfs;-56::EXCLUDE:EQUALS,EXE_NAME,tardis;-43::EXCLUDE:STARTS,EXE_PATH,/tmp/buildpacks/;-37::INCLUDE:CONTAINS,CLOUD_FOUNDRY_APP_NAME,;-38::EXCLUDE:EQUALS,KUBERNETES_CONTAINERNAME,POD;-39::EXCLUDE:CONTAINS,CONTAINER_IMAGE_NAME,pause-amd64;-44::EXCLUDE:EQUALS,EXE_NAME,oc;-58::EXCLUDE:EQUALS,EXE_NAME,calico-node;-67::EXCLUDE:EQUALS,EXE_PATH,/usr/bin/piper;-69::EXCLUDE:EQUALS,KUBERNETES_CONTAINERNAME,cassandra-operator;-70::EXCLUDE:CONTAINS,EXE_NAME,UiPath;-71::EXCLUDE:EQUALS,EXE_NAME,openhandlecollector.exe;-40::INCLUDE:CONTAINS,KUBERNETES_NAMESPACE,;-41::INCLUDE:CONTAINS,CONTAINER_NAME,;-46::EXCLUDE:EQUALS,EXE_PATH,/opt/cni/bin/host-local;-48::EXCLUDE:STARTS,EXE_NAME,mqsi;-57::EXCLUDE:STARTS,JAVA_JAR_FILE,org.eclipse.equinox.launcher;-59::EXCLUDE:EQUALS,EXE_NAME,casclient.exe;-60::EXCLUDE:EQUALS,JAVA_JAR_FILE,dynatrace_ibm_mq_connector.jar;-61::EXCLUDE:CONTAINS,EXE_NAME,Agent.Worker;-62::EXCLUDE:CONTAINS,DOTNET_COMMAND,Agent.Worker;-63::EXCLUDE:CONTAINS,EXE_NAME,Agent.Listener;-64::EXCLUDE:CONTAINS,DOTNET_COMMAND,Agent.Listener;-65::EXCLUDE:EQUALS,EXE_NAME,FlexNetJobExecutorService;-66::EXCLUDE:EQUALS,EXE_NAME,FlexNetMaintenanceRemotingService"},{"section":"agentType","key":"java","value":"on"},{"section":"general","key":"addContainerNameToPGI","value":"off"},{"section":"general","key":"expandJavaAtFiles","value":"off"},{"section":"general","key":"disableAgentTypeBasedInjection","value":"off"},{"section":"general","key":"addNodejsScriptNameToPGI","value":"on"},{"section":"general","key":"enableNodeJsAgentEsmLoaders","value":"off"},{"section":"general","key":"cassandraClusterNameInPG","value":"on"},{"section":"general","key":"disableSpringBootGroupCalc","value":"off"},{"section":"general","key":"enablePodmanInjection","value":"off"},{"section":"general","key":"enableNodeJsAgentPreloading","value":"on"},{"section":"general","key":"tenant","value":"zib50933"},{"section":"general","key":"enableTibcoBWContainerEditionGroupCalc","value":"on"},{"section":"general","key":"switchToPhpAgentNG","value":"on"},{"section":"general","key":"fullStackJavaInMBProcesses","value":"on"},{"section":"general","key":"envoyInjection","value":"on"},{"section":"general","key":"containerInjectionRules","value":"-1::EXCLUDE:EQUALS,KUBERNETES_CONTAINERNAME,POD;-2::EXCLUDE:CONTAINS,IMAGE_NAME,pause-amd64;-3::EXCLUDE:EQUALS,KUBERNETES_NAMESPACE,openshift-sdn;-4::EXCLUDE:ENDS,KUBERNETES_FULLPODNAME,-build"},{"section":"agentType","key":"sdk","value":"on"},{"section":"general","key":"containerdInjection","value":"on"},{"section":"agentType","key":"apache","value":"on"},{"section":"agentType","key":"nodejs","value":"on"},{"section":"general","key":"dockerWindowsInjection","value":"off"},{"section":"agentType","key":"varnish","value":"off"},{"section":"agentType","key":"python","value":"off"},{"section":"general","key":"dockerInjection","value":"on"},{"section":"agentType","key":"plugin","value":"on"},{"section":"general","key":"runcInjection","value":"on"},{"section":"general","key":"trustedTimestampVerificationJavascriptAgent","value":"off"},{"section":"general","key":"enableWebSphereLibertyGroupCalc","value":"off"},{"section":"general","key":"optimizedSuspendThreads","value":"on"},{"section":"agentType","key":"dotnet","value":"not-global"},{"section":"general","key":"php80InjectionEA","value":"on"},{"section":"general","key":"stripIdsFromKubernetesNamespace","value":"on"},{"section":"general","key":"pythonInjection","value":"off"},{"section":"agentType","key":"iis","value":"on"},{"section":"general","key":"coreclrInjection","value":"on"},{"section":"general","key":"crioInjection","value":"on"},{"section":"general","key":"enableEquinoxGroupCalc","value":"on"},{"section":"agentType","key":"loganalytics","value":"on"},{"section":"agentType","key":"network","value":"on"},{"section":"general","key":"enableNodeJsMultiversionLibrary","value":"off"},{"section":"agentType","key":"go","value":"not-global"},{"section":"general","key":"enableOsAgentDefaultIdCalc","value":"off"},{"section":"general","key":"php81Injection","value":"on"},{"section":"general","key":"stripVersionFromImageName","value":"on"}],"revision":3551371520520734114}
`
)

func TestConversions(t *testing.T) {
	ruxit, err := FromJson(strings.NewReader(expectedJson))
	require.NoError(t, err)

	rawString := ruxit.ToString()
	require.Equal(t, expectedConf, rawString)

	ruxit2, err := FromConf(strings.NewReader(rawString))
	require.NoError(t, err)
	require.NotEmpty(t, ruxit2)
	require.ElementsMatch(t, ruxit.Properties, ruxit2.Properties)

	rawString2 := ruxit2.ToString()
	require.Equal(t, expectedConf, rawString2)
	require.Equal(t, ruxit.ToMap(), ruxit2.ToMap())
	require.ElementsMatch(t, FromMap(ruxit.ToMap()).Properties, FromMap(ruxit2.ToMap()).Properties)
}

func TestMerge(t *testing.T) {
	t.Run("empty + override == override", func(t *testing.T) {
		source := ProcConf{}
		override := ProcConf{
			Properties: []Property{
				{
					Section: "test",
					Key:     "key",
					Value:   "value",
				},
			},
			Revision: 1,
		}

		merged := source.Merge(override)

		assert.Equal(t, override, merged)
	})
	t.Run("add", func(t *testing.T) {
		expectedProps := []Property{
			{
				Section: "test",
				Key:     "key1",
				Value:   "value1",
			},
			{
				Section: "test",
				Key:     "key2",
				Value:   "value2",
			},
		}

		source := ProcConf{
			Properties: []Property{
				expectedProps[0],
			},
			Revision: 0,
		}
		override := ProcConf{
			Properties: []Property{
				expectedProps[1],
			},
			Revision: 1,
		}

		merged := source.Merge(override)

		assert.Equal(t, override.Revision, merged.Revision)
		assert.ElementsMatch(t, expectedProps, merged.Properties)
	})

	t.Run("override + add", func(t *testing.T) {
		expectedProps := []Property{
			{
				Section: "test",
				Key:     "key1",
				Value:   "value1",
			},
			{
				Section: "test",
				Key:     "key2",
				Value:   "value2",
			},
		}

		source := ProcConf{
			Properties: []Property{
				{
					Section: expectedProps[0].Section,
					Key:     expectedProps[0].Key,
					Value:   "old value",
				},
			},
			Revision: 0,
		}
		override := ProcConf{
			Properties: expectedProps,
			Revision:   1,
		}

		merged := source.Merge(override)

		assert.Equal(t, override.Revision, merged.Revision)
		assert.ElementsMatch(t, expectedProps, merged.Properties)
	})
}

func TestSetupReadonly(t *testing.T) {
	t.Run("do adjustments according to installPath", func(t *testing.T) {
		installPath := "/absolute/path"
		expectedProps := []Property{
			{
				// updated from override
				Section: "test",
				Key:     "update",
				Value:   "updated-value",
			},
			{
				// overwritten to be absolute path
				Section: "test",
				Key:     "keyWithPath",
				Value:   "\"/absolute/path/relative/path\"",
			},
			{
				// added from override
				Section: "test",
				Key:     "add",
				Value:   "added-value",
			},
			{
				// added to be work with a readonly CodeModule bin
				Section: "general",
				Key:     "storage",
				Value:   "\"/var/lib/dynatrace/oneagent\"",
			},
		}

		sourceProps := []Property{
			{
				Section: "test",
				Key:     "update",
				Value:   "old-value",
			},
			{
				Section: "test",
				Key:     "keyWithPath",
				Value:   "\"../../relative/path\"",
			},
			{
				// will be removed, as it is not needed in readonly
				Section: "general",
				Key:     "logDir",
				Value:   "some-path",
			},
			{
				// will be removed, as it is not needed in readonly
				Section: "general",
				Key:     "dataStorageDir",
				Value:   "some-path",
			},
		}

		overrideProps := []Property{
			{
				Section: "test",
				Key:     "update",
				Value:   "updated-value",
			},
			{
				Section: "test",
				Key:     "add",
				Value:   "added-value",
			},
		}

		source := ProcConf{
			Properties: sourceProps,
			Revision:   0,
		}
		override := ProcConf{
			Properties:  overrideProps,
			Revision:    1,
			InstallPath: &installPath,
		}

		merged := source.Merge(override)

		assert.Equal(t, override.Revision, merged.Revision)
		assert.Equal(t, *override.InstallPath, *merged.InstallPath)

		expected := ProcConf{Properties: expectedProps}.ToString()
		assert.Equal(t, expected, merged.ToString())
	})

}

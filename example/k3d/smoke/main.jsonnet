local monitoring = import './monitoring/main.jsonnet';
local cortex = import 'cortex/main.libsonnet';
local loki = import 'loki-simple-scalable/loki.libsonnet';
local canary = import 'loki-canary/loki-canary.libsonnet';
local avalanche = import 'grafana-agent/smoke/avalanche/main.libsonnet';
local crow = import 'grafana-agent/smoke/crow/main.libsonnet';
local etcd = import 'grafana-agent/smoke/etcd/main.libsonnet';
local smoke = import 'grafana-agent/smoke/main.libsonnet';
local gragent = import 'grafana-agent/v2/main.libsonnet';
local k = import 'ksonnet-util/kausal.libsonnet';

local namespace = k.core.v1.namespace;
local pvc = k.core.v1.persistentVolumeClaim;
local statefulSet = k.apps.v1.statefulSet;
local volumeMount = k.core.v1.volumeMount;

local images = {
  agent: 'grafana/agent:main',
  agentctl: 'grafana/agentctl:main',
};

local new_crow(name, selector) =
  crow.new(name, namespace='smoke', config={
    args+: {
      'crow.prometheus-addr': 'http://cortex/api/prom',
      'crow.extra-selectors': selector,
    },
  });

local new_smoke(name) = smoke.new(name, namespace='smoke', config={
  mutationFrequency: '5m',
  chaosFrequency: '30m',
});

local loki_deployment = loki {
  _images+:: {
    loki: 'grafana/loki:2.5.0',
  },

  _config+:: {
    namespace: 'smoke',
    headless_service_name: 'loki-headless',
    http_listen_port: 3100,
    read_replicas: 1,
    write_replicas: 1,
    loki: {
      auth_enabled: false,
      server: {
        http_listen_port: 9095,
        grpc_listen_port: 9096,
      },
      memberlist: {
        join_members: [
          '%s.%s.svc.cluster.local' % [$._config.headless_service_name, $._config.namespace],
        ],
      },
      common: {
        path_prefix: '/loki',
        replication_factor: 1,
        storage: {
          filesystem : {
	    chunks_directory : "/tmp/loki/chunks",
            rules_directory: "/tmp/loki/rules",
          },
          //gcs: {
          //  bucket_name: "my-bucket-name",
          //},
        },
      },
      limits_config: {
        enforce_metric_name: false,
        reject_old_samples_max_age: '168h',  //1 week
        max_global_streams_per_user: 60000,
        ingestion_rate_mb: 75,
        ingestion_burst_size_mb: 100,
      },
      schema_config: {
        configs: [{
          from: '2021-09-12',
          store: 'boltdb-shipper',
          object_store: 'filesystem',
          schema: 'v11',
          index: {
            prefix: '%s_index_' % $._config.namespace,
            period: '24h',
          },
        }],
      },
    },
  },

  read_statefulset +:: 
    statefulSet.mixin.metadata.withNamespace('smoke'),
  write_statefulset +:: 
    statefulSet.mixin.metadata.withNamespace('smoke'),

  write_pvc::
    pvc.new() +
    pvc.mixin.metadata.withName('write-data') +
    pvc.mixin.metadata.withNamespace('smoke') +
    pvc.mixin.spec.resources.withRequests({ storage: '10Gi' }) +
    pvc.mixin.spec.withAccessModes(['ReadWriteOnce']) +
    pvc.mixin.spec.withStorageClassName('local-path'),

  read_pvc::
    pvc.new() +
    pvc.mixin.metadata.withName('read-data') +
    pvc.mixin.metadata.withNamespace('smoke') +
    pvc.mixin.spec.resources.withRequests({ storage: '10Gi' }) +
    pvc.mixin.spec.withAccessModes(['ReadWriteOnce']) +
    pvc.mixin.spec.withStorageClassName('local-path'),
};

local loki_canary = canary {
  _config+:: {
    namespace: "smoke",
    loki_canary_read_key: 'loki-canary-read-key',
    loki_canary_user: 104334,
    loki_canary_hostname: 'loki-ssd.grafana.net',
    loki_canary_metric_test_range: '6h',
    loki_canary_in_cluster_tls: true,
    loki_canary_in_cluster_metric_test_range: '6h',
  },

  loki_canary_args+:: {
    addr: $._config.loki_canary_hostname,
    tls: $._config.loki_canary_in_cluster_tls,
    port: 80,
    labelname: 'pod',
    interval: '500ms',
    size: 1024,
    wait: '1m',
    'metric-test-interval': '30m',
    'metric-test-range': $._config.loki_canary_in_cluster_metric_test_range,
    user: $._config.loki_canary_user,
    pass: $._config.loki_canary_read_key,
  },

  read_pvc+::
    pvc.new() +
    pvc.mixin.metadata.withName('read-data') +
    pvc.mixin.metadata.withNamespace('smoke') +
    pvc.mixin.spec.resources.withRequests({ storage: '10Gi' }) +
    pvc.mixin.spec.withAccessModes(['ReadWriteOnce']),

  local deployment = k.apps.v1.deployment,
  local statefulSet = k.apps.v1.statefulSet,
  local service = k.core.v1.service,

  loki_canary_daemonset: {},
  loki_canary_deployment: deployment.new('loki-canary', 3, [$.loki_canary_container]) +
      	                  deployment.mixin.metadata.withNamespace("smoke") +
                          k.util.antiAffinity,

  //loki_canary_container+::
  //  k.core.v1.container.withEnvMixin([
  //    k.core.v1.envVar.fromSecretRef(
  //      'LOKI_CANARY_READ',
  //      'loki-canary',
  //      'loki-canary-read'
  //    ),
  //  ]),
};

local smoke = {
  ns: namespace.new('smoke'),

  cortex: cortex.new('smoke'),

  loki: loki_deployment,

  canary: loki_canary,

  // Needed to run agent cluster
  etcd: etcd.new('smoke'),

  avalanche: avalanche.new(replicas=3, namespace='smoke', config={
    // We're going to be running a lot of these and we're not trying to test
    // for load, so reduce the cardinality and churn rate.
    metric_count: 1000,
    series_interval: 300,
    metric_interval: 600,
  }),

  smoke_test: new_smoke('smoke-test'),

  crows: [
    new_crow('crow-single', 'cluster="grafana-agent"'),
    new_crow('crow-cluster', 'cluster="grafana-agent-cluster"'),
  ],

  local metric_instances(crow_name) = [{
    name: 'crow',
    remote_write: [
      {
        url: 'http://cortex/api/prom/push',
        write_relabel_configs: [
          {
            source_labels: ['__name__'],
            regex: 'avalanche_.*',
            action: 'drop',
          },
        ],
      },
      {
        url: 'http://smoke-test:19090/api/prom/push',
        write_relabel_configs: [
          {
            source_labels: ['__name__'],
            regex: 'avalanche_.*',
            action: 'keep',
          },
        ],
      },
    ],
    scrape_configs: [
      {
        job_name: 'crow',
        metrics_path: '/validate',

        kubernetes_sd_configs: [{ role: 'pod' }],
        tls_config: {
          ca_file: '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt',
        },
        bearer_token_file: '/var/run/secrets/kubernetes.io/serviceaccount/token',

        relabel_configs: [{
          source_labels: ['__meta_kubernetes_namespace'],
          regex: 'smoke',
          action: 'keep',
        }, {
          source_labels: ['__meta_kubernetes_pod_container_name'],
          regex: crow_name,
          action: 'keep',
        }],
      },
    ],
  }, {
    name: 'avalanche',
    remote_write: [
      {
        url: 'http://cortex/api/prom/push',
        write_relabel_configs: [
          {
            source_labels: ['__name__'],
            regex: 'avalanche_.*',
            action: 'drop',
          },
        ],
      },
      {
        url: 'http://smoke-test:19090/api/prom/push',
        write_relabel_configs: [
          {
            source_labels: ['__name__'],
            regex: 'avalanche_.*',
            action: 'keep',
          },
        ],
      },
    ],
    scrape_configs: [
      {
        job_name: 'avalanche',
        kubernetes_sd_configs: [{ role: 'pod' }],
        tls_config: {
          ca_file: '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt',
        },
        bearer_token_file: '/var/run/secrets/kubernetes.io/serviceaccount/token',

        relabel_configs: [{
          source_labels: ['__meta_kubernetes_namespace'],
          regex: 'smoke',
          action: 'keep',
        }, {
          source_labels: ['__meta_kubernetes_pod_container_name'],
          regex: 'avalanche',
          action: 'keep',
        }],
      },
    ],
  }],

  normal_agent:
    gragent.new(name='grafana-agent', namespace='smoke') +
    gragent.withImagesMixin(images) +
    gragent.withStatefulSetController(
      replicas=1,
      volumeClaims=[
        pvc.new() +
        pvc.mixin.metadata.withName('agent-wal') +
        pvc.mixin.metadata.withNamespace('smoke') +
        pvc.mixin.spec.withAccessModes('ReadWriteOnce') +
        pvc.mixin.spec.resources.withRequests({ storage: '5Gi' }),
      ],
    ) +
    gragent.withVolumeMountsMixin([volumeMount.new('agent-wal', '/var/lib/agent')]) +
    gragent.withService() +
    gragent.withAgentConfig({
      server: { log_level: 'debug' },

      prometheus: {
        global: {
          scrape_interval: '1m',
          external_labels: {
            cluster: 'grafana-agent',
          },
        },
        wal_directory: '/var/lib/agent/data',
        configs: metric_instances('crow-single'),
      },
    }),

  cluster_agent:
    gragent.new(name='grafana-agent-cluster', namespace='smoke') +
    gragent.withImagesMixin(images) +
    gragent.withStatefulSetController(
      replicas=3,
      volumeClaims=[
        pvc.new() +
        pvc.mixin.metadata.withName('agent-cluster-wal') +
        pvc.mixin.metadata.withNamespace('smoke') +
        pvc.mixin.spec.withAccessModes('ReadWriteOnce') +
        pvc.mixin.spec.resources.withRequests({ storage: '5Gi' }),
      ],
    ) +
    gragent.withVolumeMountsMixin([volumeMount.new('agent-cluster-wal', '/var/lib/agent')]) +
    gragent.withService() +
    gragent.withAgentConfig({
      server: { log_level: 'debug' },

      prometheus: {
        global: {
          scrape_interval: '1m',
          external_labels: {
            cluster: 'grafana-agent-cluster',
          },
        },
        wal_directory: '/var/lib/agent/data',

        scraping_service: {
          enabled: true,
          dangerous_allow_reading_files: true,
          kvstore: {
            store: 'etcd',
            etcd: { endpoints: ['etcd:2379'] },
          },
          lifecycler: {
            ring: {
              kvstore: {
                store: 'etcd',
                etcd: { endpoints: ['etcd:2379'] },
              },
            },
          },
        },
      },
    }),

  // Spawn a syncer so our cluster gets the same scrape jobs as our
  // normal agent.
  sycner: gragent.newSyncer(
    name='grafana-agent-syncer',
    namespace='smoke',
    config={
      image: images.agentctl,
      api: 'http://grafana-agent-cluster.smoke.svc.cluster.local',
      configs: metric_instances('crow-cluster'),
    }
  ),
};

{
  monitoring: monitoring,
  smoke: smoke,
}

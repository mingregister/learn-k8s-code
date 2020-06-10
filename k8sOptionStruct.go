# kube-apiserver 
// ServerRunOptions runs a kubernetes api server.
type ServerRunOptions struct {
    GenericServerRunOptions *genericoptions.ServerRunOptions
    Etcd                    *genericoptions.EtcdOptions
    SecureServing           *genericoptions.SecureServingOptionsWithLoopback
    InsecureServing         *genericoptions.DeprecatedInsecureServingOptionsWithLoopback
    Audit                   *genericoptions.AuditOptions
    Features                *genericoptions.FeatureOptions
    Admission               *kubeoptions.AdmissionOptions
    Authentication          *kubeoptions.BuiltInAuthenticationOptions
    Authorization           *kubeoptions.BuiltInAuthorizationOptions
    CloudProvider           *kubeoptions.CloudProviderOptions
    APIEnablement           *genericoptions.APIEnablementOptions

    AllowPrivileged           bool
    EnableLogsHandler         bool
    EventTTL                  time.Duration
    KubeletConfig             kubeletclient.KubeletClientConfig
    KubernetesServiceNodePort int
    MaxConnectionBytesPerSec  int64
    ServiceClusterIPRange     net.IPNet // TODO: make this a list
    ServiceNodePortRange      utilnet.PortRange
    SSHKeyfile                string
    SSHUser                   string

    ProxyClientCertFile string
    ProxyClientKeyFile  string

    EnableAggregatorRouting bool

    MasterCount            int
    EndpointReconcilerType string

    ServiceAccountSigningKeyFile     string
    ServiceAccountIssuer             serviceaccount.TokenGenerator
    ServiceAccountTokenMaxExpiration time.Duration
}



# kube-controller-manager
// KubeControllerManagerOptions is the main context object for the kube-controller manager.
type KubeControllerManagerOptions struct {
    Generic           *cmoptions.GenericControllerManagerConfigurationOptions
    KubeCloudShared   *cmoptions.KubeCloudSharedOptions
    ServiceController *cmoptions.ServiceControllerOptions

    AttachDetachController           *AttachDetachControllerOptions
    CSRSigningController             *CSRSigningControllerOptions
    DaemonSetController              *DaemonSetControllerOptions
    DeploymentController             *DeploymentControllerOptions
    DeprecatedFlags                  *DeprecatedControllerOptions
    EndpointController               *EndpointControllerOptions
    GarbageCollectorController       *GarbageCollectorControllerOptions
    HPAController                    *HPAControllerOptions
    JobController                    *JobControllerOptions
    NamespaceController              *NamespaceControllerOptions
    NodeIPAMController               *NodeIPAMControllerOptions
    NodeLifecycleController          *NodeLifecycleControllerOptions
    PersistentVolumeBinderController *PersistentVolumeBinderControllerOptions
    PodGCController                  *PodGCControllerOptions
    ReplicaSetController             *ReplicaSetControllerOptions
    ReplicationController            *ReplicationControllerOptions
    ResourceQuotaController          *ResourceQuotaControllerOptions
    SAController                     *SAControllerOptions
    TTLAfterFinishedController       *TTLAfterFinishedControllerOptions

    SecureServing *apiserveroptions.SecureServingOptionsWithLoopback
    // TODO: remove insecure serving mode
    InsecureServing *apiserveroptions.DeprecatedInsecureServingOptionsWithLoopback
    Authentication  *apiserveroptions.DelegatingAuthenticationOptions
    Authorization   *apiserveroptions.DelegatingAuthorizationOptions

    Master     string
    Kubeconfig string
}


# kube-proxy
// Options contains everything necessary to create and run a proxy server.
type Options struct {
    // ConfigFile is the location of the proxy server's configuration file.
    ConfigFile string
    // WriteConfigTo is the path where the default configuration will be written.
    WriteConfigTo string
    // CleanupAndExit, when true, makes the proxy server clean up iptables rules, then exit.
    CleanupAndExit bool
    // CleanupIPVS, when true, makes the proxy server clean up ipvs rules before running.
    CleanupIPVS bool
    // WindowsService should be set to true if kube-proxy is running as a service on Windows.
    // Its corresponding flag only gets registered in Windows builds
    WindowsService bool
    // config is the proxy server's configuration object.
    config *kubeproxyconfig.KubeProxyConfiguration
    // watcher is used to watch on the update change of ConfigFile
    watcher filesystem.FSWatcher
    // proxyServer is the interface to run the proxy server
    proxyServer proxyRun
    // errCh is the channel that errors will be sent
    errCh chan error

    // The fields below here are placeholders for flags that can't be directly mapped into
    // config.KubeProxyConfiguration.
    //
    // TODO remove these fields once the deprecated flags are removed.

    // master is used to override the kubeconfig's URL to the apiserver.
    master string
    // healthzPort is the port to be used by the healthz server.
    healthzPort int32
    // metricsPort is the port to be used by the metrics server.
    metricsPort int32

    scheme *runtime.Scheme
    codecs serializer.CodecFactory

    // hostnameOverride, if set from the command line flag, takes precedence over the `HostnameOverride` value from the config file
    hostnameOverride string
}


# kubelet
// KubeletFlags contains configuration flags for the Kubelet.
// A configuration field should go in KubeletFlags instead of KubeletConfiguration if any of these are true:
// - its value will never, or cannot safely be changed during the lifetime of a node, or
// - its value cannot be safely shared between nodes at the same time (e.g. a hostname);
//   KubeletConfiguration is intended to be shared between nodes.
// In general, please try to avoid adding flags or configuration fields,
// we already have a confusingly large amount of them.
type KubeletFlags struct {
    KubeConfig          string
    BootstrapKubeconfig string

    // Insert a probability of random errors during calls to the master.
    ChaosChance float64
    // Crash immediately, rather than eating panics.
    ReallyCrashForTesting bool

    // TODO(mtaufen): It is increasingly looking like nobody actually uses the
    //                Kubelet's runonce mode anymore, so it may be a candidate
    //                for deprecation and removal.
    // If runOnce is true, the Kubelet will check the API server once for pods,
    // run those in addition to the pods specified by static pod files, and exit.
    RunOnce bool

    // enableServer enables the Kubelet's server
    EnableServer bool

    // HostnameOverride is the hostname used to identify the kubelet instead
    // of the actual hostname.
    HostnameOverride string
    // NodeIP is IP address of the node.
    // If set, kubelet will use this IP address for the node.
    NodeIP string

    // This flag, if set, sets the unique id of the instance that an external provider (i.e. cloudprovider)
    // can use to identify a specific node
    ProviderID string

    // Container-runtime-specific options.
    config.ContainerRuntimeOptions

    // certDirectory is the directory where the TLS certs are located (by
    // default /var/run/kubernetes). If tlsCertFile and tlsPrivateKeyFile
    // are provided, this flag will be ignored.
    CertDirectory string

    // cloudProvider is the provider for cloud services.
    // +optional
    CloudProvider string

    // cloudConfigFile is the path to the cloud provider configuration file.
    // +optional
    CloudConfigFile string

    // rootDirectory is the directory path to place kubelet files (volume
    // mounts,etc).
    RootDirectory string

    // The Kubelet will use this directory for checkpointing downloaded configurations and tracking configuration health.
    // The Kubelet will create this directory if it does not already exist.
    // The path may be absolute or relative; relative paths are under the Kubelet's current working directory.
    // Providing this flag enables dynamic kubelet configuration.
    // To use this flag, the DynamicKubeletConfig feature gate must be enabled.
    DynamicConfigDir cliflag.StringFlag

    // The Kubelet will load its initial configuration from this file.
    // The path may be absolute or relative; relative paths are under the Kubelet's current working directory.
    // Omit this flag to use the combination of built-in default configuration values and flags.
    KubeletConfigFile string

    // registerNode enables automatic registration with the apiserver.
    RegisterNode bool

    // registerWithTaints are an array of taints to add to a node object when
    // the kubelet registers itself. This only takes effect when registerNode
    // is true and upon the initial registration of the node.
    RegisterWithTaints []core.Taint

    // WindowsService should be set to true if kubelet is running as a service on Windows.
    // Its corresponding flag only gets registered in Windows builds.
    WindowsService bool

    // EXPERIMENTAL FLAGS
    // Whitelist of unsafe sysctls or sysctl patterns (ending in *).
    // +optional
    AllowedUnsafeSysctls []string
    // containerized should be set to true if kubelet is running in a container.
    Containerized bool
    // remoteRuntimeEndpoint is the endpoint of remote runtime service
    RemoteRuntimeEndpoint string
    // remoteImageEndpoint is the endpoint of remote image service
    RemoteImageEndpoint string
    // experimentalMounterPath is the path of mounter binary. Leave empty to use the default mount path
    ExperimentalMounterPath string
    // If enabled, the kubelet will integrate with the kernel memcg notification to determine if memory eviction thresholds are crossed rather than polling.
    // +optional
    ExperimentalKernelMemcgNotification bool
    // This flag, if set, enables a check prior to mount operations to verify that the required components
    // (binaries, etc.) to mount the volume are available on the underlying node. If the check is enabled
    // and fails the mount operation fails.
    ExperimentalCheckNodeCapabilitiesBeforeMount bool
    // This flag, if set, will avoid including `EvictionHard` limits while computing Node Allocatable.
    // Refer to [Node Allocatable](https://git.k8s.io/community/contributors/design-proposals/node/node-allocatable.md) doc for more information.
    ExperimentalNodeAllocatableIgnoreEvictionThreshold bool
    // Node Labels are the node labels to add when registering the node in the cluster
    NodeLabels map[string]string
    // volumePluginDir is the full path of the directory in which to search
    // for additional third party volume plugins
    VolumePluginDir string
    // lockFilePath is the path that kubelet will use to as a lock file.
    // It uses this file as a lock to synchronize with other kubelet processes
    // that may be running.
    LockFilePath string
    // ExitOnLockContention is a flag that signifies to the kubelet that it is running
    // in "bootstrap" mode. This requires that 'LockFilePath' has been set.
    // This will cause the kubelet to listen to inotify events on the lock file,
    // releasing it and exiting when another process tries to open that file.
    ExitOnLockContention bool
    // seccompProfileRoot is the directory path for seccomp profiles.
    SeccompProfileRoot string
    // bootstrapCheckpointPath is the path to the directory containing pod checkpoints to
    // run on restore
    BootstrapCheckpointPath string
    // NodeStatusMaxImages caps the number of images reported in Node.Status.Images.
    // This is an experimental, short-term flag to help with node scalability.
    NodeStatusMaxImages int32

    // DEPRECATED FLAGS
    // minimumGCAge is the minimum age for a finished container before it is
    // garbage collected.
    MinimumGCAge metav1.Duration
    // maxPerPodContainerCount is the maximum number of old instances to
    // retain per container. Each container takes up some disk space.
    MaxPerPodContainerCount int32
    // maxContainerCount is the maximum number of old instances of containers
    // to retain globally. Each container takes up some disk space.
    MaxContainerCount int32
    // masterServiceNamespace is The namespace from which the kubernetes
    // master services should be injected into pods.
    MasterServiceNamespace string
    // registerSchedulable tells the kubelet to register the node as
    // schedulable. Won't have any effect if register-node is false.
    // DEPRECATED: use registerWithTaints instead
    RegisterSchedulable bool
    // nonMasqueradeCIDR configures masquerading: traffic to IPs outside this range will use IP masquerade.
    NonMasqueradeCIDR string
    // This flag, if set, instructs the kubelet to keep volumes from terminated pods mounted to the node.
    // This can be useful for debugging volume related issues.
    KeepTerminatedPodVolumes bool
    // allowPrivileged enables containers to request privileged mode.
    // Defaults to true.
    AllowPrivileged bool
    // hostNetworkSources is a comma-separated list of sources from which the
    // Kubelet allows pods to use of host network. Defaults to "*". Valid
    // options are "file", "http", "api", and "*" (all sources).
    HostNetworkSources []string
    // hostPIDSources is a comma-separated list of sources from which the
    // Kubelet allows pods to use the host pid namespace. Defaults to "*".
    HostPIDSources []string
    // hostIPCSources is a comma-separated list of sources from which the
    // Kubelet allows pods to use the host ipc namespace. Defaults to "*".
    HostIPCSources []string
}
// KubeletConfiguration contains the configuration for the Kubelet
type KubeletConfiguration struct {
    metav1.TypeMeta

    // staticPodPath is the path to the directory containing local (static) pods to
    // run, or the path to a single static pod file.
    StaticPodPath string
    // syncFrequency is the max period between synchronizing running
    // containers and config
    SyncFrequency metav1.Duration
    // fileCheckFrequency is the duration between checking config files for
    // new data
    FileCheckFrequency metav1.Duration
    // httpCheckFrequency is the duration between checking http for new data
    HTTPCheckFrequency metav1.Duration
    // staticPodURL is the URL for accessing static pods to run
    StaticPodURL string
    // staticPodURLHeader is a map of slices with HTTP headers to use when accessing the podURL
    StaticPodURLHeader map[string][]string
    // address is the IP address for the Kubelet to serve on (set to 0.0.0.0
    // for all interfaces)
    Address string
    // port is the port for the Kubelet to serve on.
    Port int32
    // readOnlyPort is the read-only port for the Kubelet to serve on with
    // no authentication/authorization (set to 0 to disable)
    ReadOnlyPort int32
    // tlsCertFile is the file containing x509 Certificate for HTTPS.  (CA cert,
    // if any, concatenated after server cert). If tlsCertFile and
    // tlsPrivateKeyFile are not provided, a self-signed certificate
    // and key are generated for the public address and saved to the directory
    // passed to the Kubelet's --cert-dir flag.
    TLSCertFile string
    // tlsPrivateKeyFile is the file containing x509 private key matching tlsCertFile
    TLSPrivateKeyFile string
    // TLSCipherSuites is the list of allowed cipher suites for the server.
    // Values are from tls package constants (https://golang.org/pkg/crypto/tls/#pkg-constants).
    TLSCipherSuites []string
    // TLSMinVersion is the minimum TLS version supported.
    // Values are from tls package constants (https://golang.org/pkg/crypto/tls/#pkg-constants).
    TLSMinVersion string
    // rotateCertificates enables client certificate rotation. The Kubelet will request a
    // new certificate from the certificates.k8s.io API. This requires an approver to approve the
    // certificate signing requests. The RotateKubeletClientCertificate feature
    // must be enabled.
    RotateCertificates bool
    // serverTLSBootstrap enables server certificate bootstrap. Instead of self
    // signing a serving certificate, the Kubelet will request a certificate from
    // the certificates.k8s.io API. This requires an approver to approve the
    // certificate signing requests. The RotateKubeletServerCertificate feature
    // must be enabled.
    ServerTLSBootstrap bool
    // authentication specifies how requests to the Kubelet's server are authenticated
    Authentication KubeletAuthentication
    // authorization specifies how requests to the Kubelet's server are authorized
    Authorization KubeletAuthorization
    // registryPullQPS is the limit of registry pulls per second.
    // Set to 0 for no limit.
    RegistryPullQPS int32
    // registryBurst is the maximum size of bursty pulls, temporarily allows
    // pulls to burst to this number, while still not exceeding registryPullQPS.
    // Only used if registryPullQPS > 0.
    RegistryBurst int32
    // eventRecordQPS is the maximum event creations per second. If 0, there
    // is no limit enforced.
    EventRecordQPS int32
    // eventBurst is the maximum size of a burst of event creations, temporarily
    // allows event creations to burst to this number, while still not exceeding
    // eventRecordQPS. Only used if eventRecordQPS > 0.
    EventBurst int32
    // enableDebuggingHandlers enables server endpoints for log collection
    // and local running of containers and commands
    EnableDebuggingHandlers bool
    // enableContentionProfiling enables lock contention profiling, if enableDebuggingHandlers is true.
    EnableContentionProfiling bool
    // healthzPort is the port of the localhost healthz endpoint (set to 0 to disable)
    HealthzPort int32
    // healthzBindAddress is the IP address for the healthz server to serve on
    HealthzBindAddress string
    // oomScoreAdj is The oom-score-adj value for kubelet process. Values
    // must be within the range [-1000, 1000].
    OOMScoreAdj int32
    // clusterDomain is the DNS domain for this cluster. If set, kubelet will
    // configure all containers to search this domain in addition to the
    // host's search domains.
    ClusterDomain string
    // clusterDNS is a list of IP addresses for a cluster DNS server. If set,
    // kubelet will configure all containers to use this for DNS resolution
    // instead of the host's DNS servers.
    ClusterDNS []string
    // streamingConnectionIdleTimeout is the maximum time a streaming connection
    // can be idle before the connection is automatically closed.
    StreamingConnectionIdleTimeout metav1.Duration
    // nodeStatusUpdateFrequency is the frequency that kubelet computes node
    // status. If node lease feature is not enabled, it is also the frequency that
    // kubelet posts node status to master. In that case, be cautious when
    // changing the constant, it must work with nodeMonitorGracePeriod in nodecontroller.
    NodeStatusUpdateFrequency metav1.Duration
    // nodeStatusReportFrequency is the frequency that kubelet posts node
    // status to master if node status does not change. Kubelet will ignore this
    // frequency and post node status immediately if any change is detected. It is
    // only used when node lease feature is enabled.
    NodeStatusReportFrequency metav1.Duration
    // nodeLeaseDurationSeconds is the duration the Kubelet will set on its corresponding Lease.
    NodeLeaseDurationSeconds int32
    // imageMinimumGCAge is the minimum age for an unused image before it is
    // garbage collected.
    ImageMinimumGCAge metav1.Duration
    // imageGCHighThresholdPercent is the percent of disk usage after which
    // image garbage collection is always run. The percent is calculated as
    // this field value out of 100.
    ImageGCHighThresholdPercent int32
    // imageGCLowThresholdPercent is the percent of disk usage before which
    // image garbage collection is never run. Lowest disk usage to garbage
    // collect to. The percent is calculated as this field value out of 100.
    ImageGCLowThresholdPercent int32
    // How frequently to calculate and cache volume disk usage for all pods
    VolumeStatsAggPeriod metav1.Duration
    // KubeletCgroups is the absolute name of cgroups to isolate the kubelet in
    KubeletCgroups string
    // SystemCgroups is absolute name of cgroups in which to place
    // all non-kernel processes that are not already in a container. Empty
    // for no container. Rolling back the flag requires a reboot.
    SystemCgroups string
    // CgroupRoot is the root cgroup to use for pods.
    // If CgroupsPerQOS is enabled, this is the root of the QoS cgroup hierarchy.
    CgroupRoot string
    // Enable QoS based Cgroup hierarchy: top level cgroups for QoS Classes
    // And all Burstable and BestEffort pods are brought up under their
    // specific top level QoS cgroup.
    CgroupsPerQOS bool
    // driver that the kubelet uses to manipulate cgroups on the host (cgroupfs or systemd)
    CgroupDriver string
    // CPUManagerPolicy is the name of the policy to use.
    // Requires the CPUManager feature gate to be enabled.
    CPUManagerPolicy string
    // CPU Manager reconciliation period.
    // Requires the CPUManager feature gate to be enabled.
    CPUManagerReconcilePeriod metav1.Duration
    // Map of QoS resource reservation percentages (memory only for now).
    // Requires the QOSReserved feature gate to be enabled.
    QOSReserved map[string]string
    // runtimeRequestTimeout is the timeout for all runtime requests except long running
    // requests - pull, logs, exec and attach.
    RuntimeRequestTimeout metav1.Duration
    // hairpinMode specifies how the Kubelet should configure the container
    // bridge for hairpin packets.
    // Setting this flag allows endpoints in a Service to loadbalance back to
    // themselves if they should try to access their own Service. Values:
    //   "promiscuous-bridge": make the container bridge promiscuous.
    //   "hairpin-veth":       set the hairpin flag on container veth interfaces.
    //   "none":               do nothing.
    // Generally, one must set --hairpin-mode=hairpin-veth to achieve hairpin NAT,
    // because promiscuous-bridge assumes the existence of a container bridge named cbr0.
    HairpinMode string
    // maxPods is the number of pods that can run on this Kubelet.
    MaxPods int32
    // The CIDR to use for pod IP addresses, only used in standalone mode.
    // In cluster mode, this is obtained from the master.
    PodCIDR string
    // The maximum number of processes per pod.  If -1, the kubelet defaults to the node allocatable pid capacity.
    PodPidsLimit int64
    // ResolverConfig is the resolver configuration file used as the basis
    // for the container DNS resolution configuration.
    ResolverConfig string
    // cpuCFSQuota enables CPU CFS quota enforcement for containers that
    // specify CPU limits
    CPUCFSQuota bool
    // CPUCFSQuotaPeriod sets the CPU CFS quota period value, cpu.cfs_period_us, defaults to 100ms
    CPUCFSQuotaPeriod metav1.Duration
    // maxOpenFiles is Number of files that can be opened by Kubelet process.
    MaxOpenFiles int64
    // contentType is contentType of requests sent to apiserver.
    ContentType string
    // kubeAPIQPS is the QPS to use while talking with kubernetes apiserver
    KubeAPIQPS int32
    // kubeAPIBurst is the burst to allow while talking with kubernetes
    // apiserver
    KubeAPIBurst int32
    // serializeImagePulls when enabled, tells the Kubelet to pull images one at a time.
    SerializeImagePulls bool
    // Map of signal names to quantities that defines hard eviction thresholds. For example: {"memory.available": "300Mi"}.
    EvictionHard map[string]string
    // Map of signal names to quantities that defines soft eviction thresholds.  For example: {"memory.available": "300Mi"}.
    EvictionSoft map[string]string
    // Map of signal names to quantities that defines grace periods for each soft eviction signal. For example: {"memory.available": "30s"}.
    EvictionSoftGracePeriod map[string]string
    // Duration for which the kubelet has to wait before transitioning out of an eviction pressure condition.
    EvictionPressureTransitionPeriod metav1.Duration
    // Maximum allowed grace period (in seconds) to use when terminating pods in response to a soft eviction threshold being met.
    EvictionMaxPodGracePeriod int32
    // Map of signal names to quantities that defines minimum reclaims, which describe the minimum
    // amount of a given resource the kubelet will reclaim when performing a pod eviction while
    // that resource is under pressure. For example: {"imagefs.available": "2Gi"}
    EvictionMinimumReclaim map[string]string
    // podsPerCore is the maximum number of pods per core. Cannot exceed MaxPods.
    // If 0, this field is ignored.
    PodsPerCore int32
    // enableControllerAttachDetach enables the Attach/Detach controller to
    // manage attachment/detachment of volumes scheduled to this node, and
    // disables kubelet from executing any attach/detach operations
    EnableControllerAttachDetach bool
    // protectKernelDefaults, if true, causes the Kubelet to error if kernel
    // flags are not as it expects. Otherwise the Kubelet will attempt to modify
    // kernel flags to match its expectation.
    ProtectKernelDefaults bool
    // If true, Kubelet ensures a set of iptables rules are present on host.
    // These rules will serve as utility for various components, e.g. kube-proxy.
    // The rules will be created based on IPTablesMasqueradeBit and IPTablesDropBit.
    MakeIPTablesUtilChains bool
    // iptablesMasqueradeBit is the bit of the iptables fwmark space to mark for SNAT
    // Values must be within the range [0, 31]. Must be different from other mark bits.
    // Warning: Please match the value of the corresponding parameter in kube-proxy.
    // TODO: clean up IPTablesMasqueradeBit in kube-proxy
    IPTablesMasqueradeBit int32
    // iptablesDropBit is the bit of the iptables fwmark space to mark for dropping packets.
    // Values must be within the range [0, 31]. Must be different from other mark bits.
    IPTablesDropBit int32
    // featureGates is a map of feature names to bools that enable or disable alpha/experimental
    // features. This field modifies piecemeal the built-in default values from
    // "k8s.io/kubernetes/pkg/features/kube_features.go".
    FeatureGates map[string]bool
    // Tells the Kubelet to fail to start if swap is enabled on the node.
    FailSwapOn bool
    // A quantity defines the maximum size of the container log file before it is rotated. For example: "5Mi" or "256Ki".
    ContainerLogMaxSize string
    // Maximum number of container log files that can be present for a container.
    ContainerLogMaxFiles int32
    // ConfigMapAndSecretChangeDetectionStrategy is a mode in which config map and secret managers are running.
    ConfigMapAndSecretChangeDetectionStrategy ResourceChangeDetectionStrategy

    /* the following fields are meant for Node Allocatable */

    // A set of ResourceName=ResourceQuantity (e.g. cpu=200m,memory=150G,pids=100) pairs
    // that describe resources reserved for non-kubernetes components.
    // Currently only cpu and memory are supported.
    // See http://kubernetes.io/docs/user-guide/compute-resources for more detail.
    SystemReserved map[string]string
    // A set of ResourceName=ResourceQuantity (e.g. cpu=200m,memory=150G,pids=100) pairs
    // that describe resources reserved for kubernetes system components.
    // Currently cpu, memory and local ephemeral storage for root file system are supported.
    // See http://kubernetes.io/docs/user-guide/compute-resources for more detail.
    KubeReserved map[string]string
    // This flag helps kubelet identify absolute name of top level cgroup used to enforce `SystemReserved` compute resource reservation for OS system daemons.
    // Refer to [Node Allocatable](https://git.k8s.io/community/contributors/design-proposals/node/node-allocatable.md) doc for more information.
    SystemReservedCgroup string
    // This flag helps kubelet identify absolute name of top level cgroup used to enforce `KubeReserved` compute resource reservation for Kubernetes node system daemons.
    // Refer to [Node Allocatable](https://git.k8s.io/community/contributors/design-proposals/node/node-allocatable.md) doc for more information.
    KubeReservedCgroup string
    // This flag specifies the various Node Allocatable enforcements that Kubelet needs to perform.
    // This flag accepts a list of options. Acceptable options are `pods`, `system-reserved` & `kube-reserved`.
    // Refer to [Node Allocatable](https://git.k8s.io/community/contributors/design-proposals/node/node-allocatable.md) doc for more information.
    EnforceNodeAllocatable []string
}




































// Options has all the params needed to run a Scheduler
type Options struct {
    // The default values. These are overridden if ConfigFile is set or by values in InsecureServing.
    ComponentConfig kubeschedulerconfig.KubeSchedulerConfiguration

    SecureServing           *apiserveroptions.SecureServingOptionsWithLoopback
    CombinedInsecureServing *CombinedInsecureServingOptions
    Authentication          *apiserveroptions.DelegatingAuthenticationOptions
    Authorization           *apiserveroptions.DelegatingAuthorizationOptions
    Deprecated              *DeprecatedOptions

    // ConfigFile is the location of the scheduler server's configuration file.
    ConfigFile string

    // WriteConfigTo is the path where the default configuration will be written.
    WriteConfigTo string

    Master string
}}

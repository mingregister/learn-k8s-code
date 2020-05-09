### 资料
https://mp.weixin.qq.com/s/bYZJ1ipx7iBPw6JXiZ3Qug


## 启动主流程
**o.Run() -->[ NewProxyServer(o)--> ] o.runLoop() --> [ o.proxyServer.Run() --> ] s.Proxier.SyncLoop()**

### o.Run()
```
func (o *Options) Run() error {
    defer close(o.errCh)
    // 1.如果指定了 --write-config-to 参数，则将默认的配置文件写到指定文件并退出
    if len(o.WriteConfigTo) > 0 {
        return o.writeConfigFile()
    }

    // 2.初始化 ProxyServer 对象
    proxyServer, err := NewProxyServer(o)
    if err != nil {
        return err
    }

    // 3.如果启动参数 --cleanup 设置为 true，则清理 iptables 和 ipvs 规则并退出
    if o.CleanupAndExit {
        return proxyServer.CleanupAndExit()
    }

    o.proxyServer = proxyServer
    return o.runLoop()
}
```

### NewProxyServer
* 初始化 iptables、ipvs 相关的 interface

* 若启用了 ipvs 则检查内核版本、ipvs 依赖的内核模块、ipset 版本，内核模块主要包括：ip_vs，ip_vs_rr,
ip_vs_wrr,ip_vs_sh,nf_conntrack_ipv4,nf_conntrack，若没有相关模块，kube-proxy 会尝试使用 modprobe 自动加载

* 根据 proxyMode 初始化 proxier，kube-proxy 启动后只运行一种 proxier
```
func NewProxyServer(o *Options) (*ProxyServer, error) {
    return newProxyServer(o.config, o.CleanupAndExit, o.master)
}

func newProxyServer(
    config *proxyconfigapi.KubeProxyConfiguration,
    cleanupAndExit bool,
    master string) (*ProxyServer, error) {
    ......

    if c, err := configz.New(proxyconfigapi.GroupName); err == nil {
        c.Set(config)
    } else {
        return nil, fmt.Errorf("unable to register configz: %s", err)
    }

    ......

    // 1.关键依赖工具 iptables/ipvs/ipset/dbus
    var iptInterface utiliptables.Interface
    var ipvsInterface utilipvs.Interface
    var kernelHandler ipvs.KernelHandler
    var ipsetInterface utilipset.Interface
    var dbus utildbus.Interface

    // 2.执行 linux 命令行的工具
    execer := exec.New()

    // 3.初始化 iptables/ipvs/ipset/dbus 对象
    dbus = utildbus.New()
    iptInterface = utiliptables.New(execer, dbus, protocol)
    kernelHandler = ipvs.NewLinuxKernelHandler()
    ipsetInterface = utilipset.New(execer)

    // 4.检查该机器是否支持使用 ipvs 模式
    canUseIPVS, _ := ipvs.CanUseIPVSProxier(kernelHandler, ipsetInterface)
    if canUseIPVS {
        ipvsInterface = utilipvs.New(execer)
    }

    if cleanupAndExit {
        return &ProxyServer{
            ......
        }, nil
    }

    // 5.初始化 kube client 和 event client
    client, eventClient, err := createClients(config.ClientConnection, master)
    if err != nil {
        return nil, err
    }
    ......

    // 6.初始化 healthzServer
    var healthzServer *healthcheck.HealthzServer
    var healthzUpdater healthcheck.HealthzUpdater
    if len(config.HealthzBindAddress) > 0 {
        healthzServer = healthcheck.NewDefaultHealthzServer(config.HealthzBindAddress, 2*config.IPTables.SyncPeriod.Duration, recorder, nodeRef)
        healthzUpdater = healthzServer
    }

    // 7.proxier 是一个 interface，每种模式都是一个 proxier
    var proxier proxy.Provider

    // 8.根据 proxyMode 初始化 proxier
    proxyMode := getProxyMode(string(config.Mode), kernelHandler, ipsetInterface, iptables.LinuxKernelCompatTester{})
    ......

    if proxyMode == proxyModeIPTables {
        klog.V(0).Info("Using iptables Proxier.")
        if config.IPTables.MasqueradeBit == nil {
            return nil, fmt.Errorf("unable to read IPTables MasqueradeBit from config")
        }

        // 9.初始化 iptables 模式的 proxier
        proxier, err = iptables.NewProxier(
            .......
        )
        if err != nil {
            return nil, fmt.Errorf("unable to create proxier: %v", err)
        }
        metrics.RegisterMetrics()
    } else if proxyMode == proxyModeIPVS {
        // 10.判断是够启用了 ipv6 双栈
        if utilfeature.DefaultFeatureGate.Enabled(features.IPv6DualStack) {
            ......
            // 11.初始化 ipvs 模式的 proxier
            proxier, err = ipvs.NewDualStackProxier(
                ......
            )
        } else {
            proxier, err = ipvs.NewProxier(
                ......
            )
        }
        if err != nil {
            return nil, fmt.Errorf("unable to create proxier: %v", err)
        }
        metrics.RegisterMetrics()
    } else {
        // 12.初始化 userspace 模式的 proxier
        proxier, err = userspace.NewProxier(
            ......
        )
        if err != nil {
            return nil, fmt.Errorf("unable to create proxier: %v", err)
        }
    }

    iptInterface.AddReloadFunc(proxier.Sync)
    return &ProxyServer{
        ......
    }, nil
}
```

### o.runLoop()
```
func (o *Options) runLoop() error {
    // 1.watch 配置文件变化
    if o.watcher != nil {
        o.watcher.Run()
    }

    // 2.以 goroutine 方式启动 proxyServer
    go func() {
        err := o.proxyServer.Run()
        o.errCh <- err
    }()

    for {
        err := <-o.errCh
        if err != nil {
            return err
        }
    }
}
```
### o.proxyServer.Run()   
o.proxyServer.Run()中会启动已经初始化好的所有服务：

* 设定进程 OOMScore，可通过命令行配置，默认值为 --oom-score-adj="-999"
* 启动 metric server 和 healthz server，两者分别监听 10256 和 10249 端口
* 设置内核参数nf_conntrack_tcp_timeout_established
和 nf_conntrack_tcp_timeout_close_wait
* 将 proxier 注册到 serviceEventHandler 和 endpointsEventHandler 中
* 启动 informer 监听 service 和 endpoints 变化
* 执行 s.Proxier.SyncLoop()，启动 proxier 主循环

```
func (s *ProxyServer) Run() error {
    ......

    // 1.进程 OOMScore，避免进程因 oom 被杀掉，此处默认值为 -999
    var oomAdjuster *oom.OOMAdjuster
    if s.OOMScoreAdj != nil {
        oomAdjuster = oom.NewOOMAdjuster()
        if err := oomAdjuster.ApplyOOMScoreAdj(0, int(*s.OOMScoreAdj)); err != nil {
            klog.V(2).Info(err)
        }
    }
    ......

    // 2.启动 healthz server
    if s.HealthzServer != nil {
        s.HealthzServer.Run()
    }

    // 3.启动 metrics server
    if len(s.MetricsBindAddress) > 0 {
        ......
        go wait.Until(func() {
            err := http.ListenAndServe(s.MetricsBindAddress, proxyMux)
            if err != nil {
                utilruntime.HandleError(fmt.Errorf("starting metrics server failed: %v", err))
            }
        }, 5*time.Second, wait.NeverStop)
    }

    // 4.配置 conntrack，设置内核参数 nf_conntrack_tcp_timeout_established 和 nf_conntrack_tcp_timeout_close_wait
    if s.Conntracker != nil {
        max, err := getConntrackMax(s.ConntrackConfiguration)
        if err != nil {
            return err
        }
        if max > 0 {
            err := s.Conntracker.SetMax(max)
            ......
        }

        if s.ConntrackConfiguration.TCPEstablishedTimeout != nil && s.ConntrackConfiguration.TCPEstablishedTimeout.Duration > 0 {
            timeout := int(s.ConntrackConfiguration.TCPEstablishedTimeout.Duration / time.Second)
            if err := s.Conntracker.SetTCPEstablishedTimeout(timeout); err != nil {
                return err
            }
        }
        if s.ConntrackConfiguration.TCPCloseWaitTimeout != nil && s.ConntrackConfiguration.TCPCloseWaitTimeout.Duration > 0 {
            timeout := int(s.ConntrackConfiguration.TCPCloseWaitTimeout.Duration / time.Second)
            if err := s.Conntracker.SetTCPCloseWaitTimeout(timeout); err != nil {
                return err
            }
        }
    }

    ......

    // 5.启动 informer 监听 Services 和 Endpoints 或者 EndpointSlices 信息
    informerFactory := informers.NewSharedInformerFactoryWithOptions(s.Client, s.ConfigSyncPeriod,
        informers.WithTweakListOptions(func(options *metav1.ListOptions) {
            options.LabelSelector = labelSelector.String()
        }))


    // 6.将 proxier 注册到 serviceConfig、endpointsConfig 中
    serviceConfig := config.NewServiceConfig(informerFactory.Core().V1().Services(), s.ConfigSyncPeriod)
    serviceConfig.RegisterEventHandler(s.Proxier)
    go serviceConfig.Run(wait.NeverStop)

    if utilfeature.DefaultFeatureGate.Enabled(features.EndpointSlice) {
        endpointSliceConfig := config.NewEndpointSliceConfig(informerFactory.Discovery().V1alpha1().EndpointSlices(), s.ConfigSyncPeriod)
        endpointSliceConfig.RegisterEventHandler(s.Proxier)
        go endpointSliceConfig.Run(wait.NeverStop)
    } else {
        endpointsConfig := config.NewEndpointsConfig(informerFactory.Core().V1().Endpoints(), s.ConfigSyncPeriod)
        endpointsConfig.RegisterEventHandler(s.Proxier)
        go endpointsConfig.Run(wait.NeverStop)
    }

    // 7.启动 informer
    informerFactory.Start(wait.NeverStop)

    s.birthCry()

    // 8.启动 proxier 主循环
    s.Proxier.SyncLoop()
    return nil
}
```

###  s.Proxier.SyncLoop()
执行 proxier 的主循环
### 我们将依次从三个方面：配置服务路由、访问权限以及同数据库（etcd）的交互对kube-apiserver源码进行解读。

## 数据结构
主要为ServerRunOptions、Config、APIGroupInfo。

### ServerRunOptions
ServerRunOptions位于k8s.io/kubernetes/cmd/kube-apiserver/app/options/options.go，是kube-apiserver的启动参数结构体，涵括了kube-apiserver启动需要所有参数，后续各个逻辑需要的相关配置参数均来源于此。其中每项参数的含义请参考官网文档：

```

type ServerRunOptions struct {
    // kube-apiserver系列服务基础配置项
  GenericServerRunOptions *genericoptions.ServerRunOptions
    // 后端存储etcd的配置项
  Etcd                    *genericoptions.EtcdOptions
  SecureServing           *genericoptions.SecureServingOptionsWithLoopback
  InsecureServing         *genericoptions.DeprecatedInsecureServingOptionsWithLoopback
  Audit                   *genericoptions.AuditOptions
  Features                *genericoptions.FeatureOptions
  Admission               *kubeoptions.AdmissionOptions
    // 集群验证配置项
  Authentication          *kubeoptions.BuiltInAuthenticationOptions
    // 集群操作鉴权配置项
  Authorization           *kubeoptions.BuiltInAuthorizationOptions
  CloudProvider           *kubeoptions.CloudProviderOptions
  APIEnablement           *genericoptions.APIEnablementOptions
  EgressSelector          *genericoptions.EgressSelectorOptions
  ...
}
```
### Config
Config位于k8s.io/apiserver/pkg/server/config.go，是kube-apiserver系列服务的通用配置项，用于与各服务对应的额外配置组合来构成对应服务的配置结构体。
```

type Config struct {
  // https服务配置项
SecureServing *SecureServingInfo
  // 服务验证配置
Authentication AuthenticationInfo
  // 服务鉴权配置
Authorization AuthorizationInfo
LoopbackClientConfig *restclient.Config
EgressSelector *egressselector.EgressSelector
RuleResolver authorizer.RuleResolver
AdmissionControl    admission.Interface
AuditBackend audit.Backend
  // 配置服务路由
BuildHandlerChainFunc func(apiHandler http.Handler, c *Config) (secure
http.Handler)
RESTOptionsGetter genericregistry.RESTOptionsGetter
  ...
}
```

### APIGroupInfo
APIGroupInfo位于k8s.io/apiserver/pkg/server/generaticapiserver.go，定义了一个 API 组的相关信息。（以customresourcedefinitions资源为例，对应的API组name为apiextensions.k8s.io，API组支持的版本信息PrioritizedVersions为v1和v1beta1两个版本，版本 、资源 、 后端CRUD结构体的映射关系为map[v1beta1]

![img](https://mmbiz.qpic.cn/mmbiz_png/DWUm33nzh3yKMIbpfDU9FBtC9fib1V2aaXlFicNtbm1UI7EOFibebzLv7qA3PTvMiatPl3TYiaalWB1BEHY42zdBbBw/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)


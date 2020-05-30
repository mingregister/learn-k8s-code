options里面一般包含了各组件的配置???

就kube-scheduler而言，是使用ApplyTo把配置都加载到组件中去。
```go kubernetes-1.18.2\cmd\kube-scheduler\app\options\options.go
	c := &schedulerappconfig.Config{}
	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}
```

通过Flags把代码与实际的配置key对应起来。如：ConfigFile与--config
```go kubernetes-1.18.2\cmd\kube-scheduler\app\options\options.go
fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file. Flags override values in this file.")
```
并且一个FlagSet代表一组配置。通过kube-scheduler --help的输出可以观察到。

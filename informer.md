https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html  
https://engineering.bitnami.com/articles/kubewatch-an-example-of-kubernetes-custom-controller.html


There are two main components of a controller: **informer/SharedInformer and Workqueue.** Informer/SharedInformer watches for changes on the current state of Kubernetes objects and sends events to Workqueue where events are then popped up by worker(s) to process.



# PodInformer
PodInformer provides access to a shared informer and lister for Pods.
``` go
// 注意，这里是interface哦。
type PodInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.PodLister
}
```



# kubewatch 
``` go
lw := cache.NewListWatchFromClient(
      client,
      &v1.Pod{},
      api.NamespaceAll,
      fieldSelector)
```

``` go
store, controller := cache.NewInformer {
	&cache.ListWatch{},
	&v1.Pod{},
	resyncPeriod,
    cache.ResourceEventHandlerFuncs{},
```

``` go
cache.ListWatch {
	listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
		return client.Get().
			Namespace(namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec).
			FieldsSelectorParam(fieldSelector).
			Do().
			Get()
	}
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		options.Watch = true
		return client.Get().
			Namespace(namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec).
			FieldsSelectorParam(fieldSelector).
			Watch()
	}
}
```


``` go
type ResourceEventHandlerFuncs struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}
```

``` go
lw := cache.NewListWatchFromClient(…)
sharedInformer := cache.NewSharedInformer(lw, &api.Pod{}, resyncPeriod)
```


``` go
queue :=
workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
```

``` go
queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

informer := cache.NewSharedIndexInformer(
      &cache.ListWatch{
             ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
                    return client.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
             },
             WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
                    return client.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
             },
      },
      &api_v1.Pod{},
      0, //Skip resync
      cache.Indexers{},
)
```

``` go
informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
      AddFunc: func(obj interface{}) {
             key, err := cache.MetaNamespaceKeyFunc(obj)
             if err == nil {
                    queue.Add(key)
             }
      },
      DeleteFunc: func(obj interface{}) {
             key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
             if err == nil {
                    queue.Add(key)
             }
      },
})
```

``` go
// Run will start the controller.
// StopCh channel is used to send interrupt signal to stop it.
func (c *Controller) Run(stopCh <-chan struct{}) {
      // don't let panics crash the process
      defer utilruntime.HandleCrash()
      // make sure the work queue is shutdown which will trigger workers to end
      defer c.queue.ShutDown()

      c.logger.Info("Starting kubewatch controller")

      go c.informer.Run(stopCh)

      // wait for the caches to synchronize before starting the worker
      if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
             utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
             return
      }

      c.logger.Info("Kubewatch controller synced and ready")

     // runWorker will loop until "something bad" happens.  The .Until will
     // then rekick the worker after one second
      wait.Until(c.runWorker, time.Second, stopCh)
}
```


``` go
func (c *Controller) runWorker() {
// processNextWorkItem will automatically wait until there's work available
      for c.processNextItem() {
             // continue looping
      }
}

// processNextWorkItem deals with one key off the queue.  It returns false
// when it's time to quit.
func (c *Controller) processNextItem() bool {
       // pull the next work item from queue.  It should be a key we use to lookup
	// something in a cache
      key, quit := c.queue.Get()
      if quit {
             return false
      }

       // you always have to indicate to the queue that you've completed a piece of
	// work
      defer c.queue.Done(key)

      // do your work on the key.
      err := c.processItem(key.(string))

      if err == nil {
             // No error, tell the queue to stop tracking history
             c.queue.Forget(key)
      } else if c.queue.NumRequeues(key) < maxRetries {
             c.logger.Errorf("Error processing %s (will retry): %v", key, err)
             // requeue the item to work on later
c.queue.AddRateLimited(key)
      } else {
             // err != nil and too many retries
             c.logger.Errorf("Error processing %s (giving up): %v", key, err)
             c.queue.Forget(key)
             utilruntime.HandleError(err)
      }

      return true
}
```


``` go
func (c *Controller) processItem(key string) error {
      c.logger.Infof("Processing change to Pod %s", key)

      obj, exists, err := c.informer.GetIndexer().GetByKey(key)
      if err != nil {
             return fmt.Errorf("Error fetching object with key %s from store: %v", key, err)
      }

      if !exists {
             c.eventHandler.ObjectDeleted(obj)
             return nil
      }

      c.eventHandler.ObjectCreated(obj)
      return nil
}
```
``` go
func (s *Slack) ObjectCreated(obj interface{}) {
      notifySlack(s, obj, "created")
}

func (s *Slack) ObjectDeleted(obj interface{}) {
      notifySlack(s, obj, "deleted")
}

func notifySlack(s *Slack, obj interface{}, action string) {
      e := kbEvent.New(obj, action)
      api := slack.New(s.Token)
      params := slack.PostMessageParameters{}
      attachment := prepareSlackAttachment(e)

      params.Attachments = []slack.Attachment{attachment}
      params.AsUser = true
      channelID, timestamp, err := api.PostMessage(s.Channel, "", params)
      if err != nil {
             log.Printf("%s\n", err)
             return
      }

      log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
```
package cluster

import (
	"k8s.io/klog"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

type Cluster struct {
	kubeconfig string
}

func New(kubeconfig string) Cluster {

	c := Cluster{kubeconfig}
	return c
}

func (c Cluster) Start() {
	master := "" // this is a CLI flag in the example, but what is it for?

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags(master, c.kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	endpointsWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "endpoints", v1.NamespaceDefault, fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(endpointsWatcher, &v1.Endpoints{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	// We can now warm up the cache for initial synchronization.
	// Let's suppose that we knew about a pod "mypod" on our last run, therefore add it to the cache.
	// If this pod is not there anymore, the controller will be notified about the removal after the
	// cache has synchronized.
	indexer.Add(&v1.Endpoints{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "if-we-had-state-from-earlier-instances-we-could-possibly-detect-deletes-better",
			Namespace: v1.NamespaceDefault,
		},
	})

	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

}

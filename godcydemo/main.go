package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
)

type ServiceDiscovery struct {
	client     *clientv3.Client
	serverList map[string]string
	lock       sync.Mutex
}

func NewServiceDiscovery(endpoints []string) (*ServiceDiscovery, error) {
	client, err := clientv3.New(
		clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		},
	)
	return &ServiceDiscovery{
		client:     client,
		serverList: make(map[string]string),
	}, err
}

func (i *ServiceDiscovery) watchService(prefix string) error {
	resp, err := i.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		i.SetServiceList(string(kv.Key), string(kv.Value))
	}
	go i.watcher(prefix)
	return nil
}

func (i *ServiceDiscovery) watcher(prefix string) {
	wChan := i.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("监听前缀: %s", prefix)
	for wResp := range wChan {
		for _, e := range wResp.Events {
			switch e.Type {
			case mvccpb.PUT:
				i.SetServiceList(string(e.Kv.Key), string(e.Kv.Value))
			case mvccpb.DELETE:
				i.DelServiceList(string(e.Kv.Key))
			}
		}
	}
}

func (i *ServiceDiscovery) SetServiceList(key string, value string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.serverList[key] = value
}

func (i *ServiceDiscovery) DelServiceList(key string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	delete(i.serverList, key)
}

func (i *ServiceDiscovery) GetServices() []string {
	i.lock.Lock()
	defer i.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range i.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (i *ServiceDiscovery) Close() error {
	return i.client.Close()
}

func main() {
	var endpoints = []string{"localhost:2379"}
	discorvery, err := NewServiceDiscovery(endpoints)
	if err != nil {
		log.Fatalln(err)
	}
	defer discorvery.Close()
	discorvery.watchService("/web")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(discorvery.GetServices())
		}
	}
}

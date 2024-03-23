package main

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type ServiceRegister struct {
	client        *clientv3.Client
	leaseId       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	value         string
}

func NewServiceRegister(
	endpoints []string,
	key string,
	value string,
	lease int64,
) (result *ServiceRegister, err error) {
	client, err := clientv3.New(
		clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		},
	)
	if err != nil {
		return
	}
	result = &ServiceRegister{
		client: client,
		key:    key,
		value:  value,
	}

	err = result.putKeyWithLease(lease)
	return
}

func (i *ServiceRegister) putKeyWithLease(lease int64) error {
	resp, err := i.client.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	_, err = i.client.Put(context.Background(), i.key, i.value, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	leaseRespChan, err := i.client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	i.leaseId = resp.ID
	i.keepAliveChan = leaseRespChan
	return nil
}

func (i *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range i.keepAliveChan {
		log.Println("续约成功", leaseKeepResp)
	}
	log.Println("关闭续约")
}

func (i *ServiceRegister) close() error {
	if _, err := i.client.Revoke(context.Background(), i.leaseId); err != nil {
		return err
	}
	return i.client.Close()
}

func main() {
	log.Println("start")
	var endpoints = []string{"localhost:2379"}
	server, err := NewServiceRegister(endpoints, "/web/node1", "localhost:8000", 40)
	if err != nil {
		log.Fatalln(err)
	}
	//go server.ListenLeaseRespChan()
	server.ListenLeaseRespChan()
	// TODO 这个示例并没有启动一个真服务，只是注册了一条信息。
}

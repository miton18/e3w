package e3ch

import (
	"crypto/tls"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/soyking/e3ch"
	"github.com/soyking/e3w/conf"
)

// NewE3chClient create a new E3 client
func NewE3chClient(config *conf.Config) (*client.EtcdHRCHYClient, error) {
	var tlsConfig *tls.Config
	var err error
	if config.CertFile != "" && config.KeyFile != "" && config.CAFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      config.CertFile,
			KeyFile:       config.KeyFile,
			TrustedCAFile: config.CAFile,
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	clt, err := clientv3.New(clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   config.EtcdEndPoints,
		Username:    config.EtcdUsername,
		Password:    config.EtcdPassword,
		TLS:         tlsConfig,
	})
	if err != nil {
		return nil, err
	}

	client, err := client.New(clt, config.EtcdRootKey, config.DirValue)
	if err != nil {
		return nil, err
	}
	return client, client.FormatRootKey()
}

// CloneE3chClient dup etcd client
func CloneE3chClient(username, password string, client *client.EtcdHRCHYClient) (*client.EtcdHRCHYClient, error) {
	clt, err := clientv3.New(clientv3.Config{
		Endpoints: client.EtcdClient().Endpoints(),
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return client.Clone(clt), nil
}

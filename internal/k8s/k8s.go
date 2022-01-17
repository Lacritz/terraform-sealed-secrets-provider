package k8s

import (
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"net/http"
	"time"
)

type Client struct {
	RestClient *corev1.CoreV1Client
}

type Config struct {
	Host                                 string
	ClusterCACert, ClientCert, ClientKey []byte
	Transport                            http.RoundTripper
}

type Clienter interface {
	Get(ctx context.Context, controllerName, controllerNamespace, path string) ([]byte, error)
}

func NewClient(cfg *Config) (*Client, error) {
	restCfg := &rest.Config{
		Timeout: 10 * time.Second,
	}
	restCfg.Host = cfg.Host
	restCfg.CAData = cfg.ClusterCACert
	restCfg.CertData = cfg.ClientCert
	restCfg.KeyData = cfg.ClientKey
	if cfg.Transport != nil {
		restCfg.Transport = cfg.Transport
	}

	c, err := corev1.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}
	return &Client{RestClient: c}, nil
}

func (c *Client) Get(ctx context.Context, controllerName, controllerNamespace, path string) ([]byte, error) {
	resp, err := c.RestClient.
		Services(controllerNamespace).
		ProxyGet("http", controllerName, "", path, nil).
		Stream(ctx)

	if err != nil {
		return nil, fmt.Errorf("request to k8s cluster failed: %w", err)
	}
	b, err := io.ReadAll(resp)
	if err != nil {
		return nil, fmt.Errorf("unable to read response from k8 cluster: %w", err)
	}
	return b, nil
}

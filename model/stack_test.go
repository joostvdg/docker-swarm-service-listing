package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProxyConfiguration_String_minimal(t *testing.T) {
	proxyConfig1 := ProxyConfiguration{ServicePath: "/demo1", Https: false}
	expectedString := "[https=false, path=/demo1, domain=, port=0]"
	assert.Equal(t, expectedString, proxyConfig1.String())
}

func TestProxyConfiguration_String_full(t *testing.T) {
	proxyConfig1 := ProxyConfiguration{ServicePath: "/demo2", Https: true, ServiceDomain: "abc.xyz.com", ServicePort: 3433}
	expectedString := "[https=true, path=/demo2, domain=abc.xyz.com, port=3433]"
	assert.Equal(t, expectedString, proxyConfig1.String())
}

func TestService_String_no_config(t *testing.T) {
	service1ProxyConfigs := make([]ProxyConfiguration, 0)
	service1 := Service{Name: "ServiceOne", Alias: "", ProxyConfigurations: service1ProxyConfigs}

	expectedString := "ServiceOne \n "
	assert.Equal(t, expectedString, service1.String())
}

func TestService_String_with_configs(t *testing.T) {
	proxyConfig1 := ProxyConfiguration{ServicePath: "/demo1", Https: false}
	proxyConfig2 := ProxyConfiguration{ServicePath: "/demo2", Https: false}
	service1ProxyConfigs := make([]ProxyConfiguration, 2)
	service1ProxyConfigs[0] = proxyConfig1
	service1ProxyConfigs[1] = proxyConfig2
	service1 := Service{Name: "ServiceOne", Alias: "", ProxyConfigurations: service1ProxyConfigs}

	expectedString := "ServiceOne \n - [https=false, path=/demo1, domain=, port=0]\n- [https=false, path=/demo2, domain=, port=0]"
	assert.Equal(t, expectedString, service1.String())
}

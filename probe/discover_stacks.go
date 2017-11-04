package probe

import (
	"../model"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"strconv"
	"strings"
)

// DiscoverStacks attempts to discover Docker Stacks via a docker host,
//  limited to Stacks that contain Services that are proxied via the Docker Flow Proxy
func DiscoverStacks() []model.Stack {
	stacks := make([]model.Stack, 0)
	host := "unix:///var/run/docker.sock"
	labelFilter := fmt.Sprintf("%s=true", "com.df.notify")

	// TODO: make host configurable
	// See: https://github.com/vfarcic/docker-flow-swarm-listener/blob/bacbeb663d420289dd461d426d9beb2521540c62/service/service.go
	fmt.Println(" > Probing Host: " + host)

	services, err := retrieveSwarmServiceList(labelFilter)

	if err != nil {
		fmt.Println(err.Error())
	} else if len(services) == 0 {
		fmt.Printf("   > No Services found with label filter %s", labelFilter)
	} else {
		proxiedServices := FindProxiedServices(services)
		stacks = make([]model.Stack, len(proxiedServices))
		count := 0
		for stackName, services := range proxiedServices {

			if stackName == "" {
				continue
			}
			stack := model.Stack{
				Name:     stackName,
				Services: services,
			}
			stacks[count] = stack
			count++
		}

	}

	fmt.Printf(" > Found %d stacks\n", len(stacks))

	return stacks
}

func retrieveSwarmServiceList(labelFilter string) ([]swarm.Service, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	filter := filters.NewArgs()
	filter.Add("label", labelFilter) // TODO: make filter optional / parameter
	return cli.ServiceList(context.Background(), types.ServiceListOptions{Filters: filter})
}

// FindProxiedServices will try to find Services that are proxied via Docker Flow Proxy
func FindProxiedServices(services []swarm.Service) map[string][]model.Service {
	proxiedServices := make(map[string][]model.Service, len(services))
	for _, service := range services {
		proxyService := model.Service{Name: service.Spec.Name}
		stackName := "Other"
		proxyConfigurations := make(map[int]*model.ProxyConfiguration, 10)
		baseProxyConfig := &model.ProxyConfiguration{Https: false}

		for key := range service.Spec.Labels {
			if key == "com.docker.stack.namespace" {
				stackName = service.Spec.Labels[key]
				if strings.HasPrefix(proxyService.Name, stackName+"_") {
					proxyService.Name = strings.TrimPrefix(proxyService.Name, stackName+"_")
				}
			}
			processServiceConfigurations(&proxyService, baseProxyConfig, proxyConfigurations, key, service.Spec.Labels)
		}

		if len(proxyConfigurations) == 0 {
			proxyService.ProxyConfigurations = append(proxyService.ProxyConfigurations, *baseProxyConfig)
		} else {
			for _, proxyConfig := range proxyConfigurations {
				fillUpFromBase(baseProxyConfig, proxyConfig)
				proxyService.ProxyConfigurations = append(proxyService.ProxyConfigurations, *proxyConfig)
			}
		}

		proxiedServices[stackName] = append(proxiedServices[stackName], proxyService)
	}
	return proxiedServices
}

// we might have several shared properties in the base
// and maybe only a single different property in the individual proxy Configs
// so we should fill up the remaining properties from the base config
func fillUpFromBase(baseConfig *model.ProxyConfiguration, config *model.ProxyConfiguration) {
	if config.ServicePort == 0 {
		config.ServicePort = baseConfig.ServicePort
	}

	if config.ServicePath == "" {
		config.ServicePath = baseConfig.ServicePath
	}

	if config.ServiceDomain == "" {
		config.ServiceDomain = baseConfig.ServiceDomain
	}
}

func processServiceConfigurations(proxyService *model.Service, baseConfig *model.ProxyConfiguration, configs map[int]*model.ProxyConfiguration, key string, labels map[string]string) {

	tmpProxyConfig := baseConfig

	if strings.HasPrefix(key, "com.df") {
		labelName := key
		labelValue := labels[key]
		labelNameParts := strings.Split(labelName, ".")
		if len(labelNameParts) < 3 {
			return
		} else if len(labelNameParts) == 4 {
			fmt.Printf("  > Found label with a prefix: %s=%s (%s)\n", labelName, labelValue, labelNameParts[3])
			i, err := strconv.Atoi(labelNameParts[3])
			if err != nil {
				fmt.Println(err)
			} else if i < 10 {
				tmpProxyConfig = configs[i]
				if tmpProxyConfig == nil {
					fmt.Printf("  > Create new config for prefix (%d)\n", i)
					tmpProxyConfig = &model.ProxyConfiguration{Https: false}
					configs[i] = tmpProxyConfig
				} else {
					fmt.Printf("  > Already had this prefix in here (%d)\n", i)
				}
			}
		}
		labelName = labelNameParts[2]

		switch labelName {
		case "httpsOnly":
			if labelValue == "true" {
				tmpProxyConfig.Https = true
			}
		case "servicePath":
			tmpProxyConfig.ServicePath = labelValue
			if !strings.HasSuffix(tmpProxyConfig.ServicePath, "/") {
				tmpProxyConfig.ServicePath += "/"
			}
		case "serviceDomain":
			tmpProxyConfig.ServiceDomain = labelValue
		case "srcPort":
			tmpProxyConfig.ServicePort, _ = strconv.Atoi(labelValue)
		}
	}
}

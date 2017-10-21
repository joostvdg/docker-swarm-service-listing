package probe

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"../model"
	"strings"
	"strconv"
	"os"
)

func DiscoverStacks() []model.Stack {
	stacks := make([]model.Stack, 0)

	// TODO: make host configurable
	// See: https://github.com/vfarcic/docker-flow-swarm-listener/blob/bacbeb663d420289dd461d426d9beb2521540c62/service/service.go
	host := "unix:///var/run/docker.sock"
	if len(os.Getenv("DF_DOCKER_HOST")) > 0 {
		host = os.Getenv("DF_DOCKER_HOST")
	}
	fmt.Println(" > Probing Host: " + host)

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	filter := filters.NewArgs()
	filter.Add("label", fmt.Sprintf("%s=true", "com.df.notify")) // TODO: make filter optional / parameter
	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filters: filter})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		proxiedServices := make(map[string][]model.Service, len(services))
		for _, service := range services {
			proxyService := model.Service{Name: service.Spec.Name}
			stackName := "Other"
			foundService := false
			proxyConfig := &model.ProxyConfiguration{Https: false}

			for key := range service.Spec.Labels {
				if key == "com.docker.stack.namespace"  {
					stackName = service.Spec.Labels[key]
				}

				if strings.HasPrefix(key, "com.df" ) {
					labelName := key
					labelValue := service.Spec.Labels[key]
					labelNameParts := strings.Split(labelName, ".")
					if len(labelNameParts) < 3 {
						continue
					}
					labelName = labelNameParts[len(labelNameParts) - 1]

					switch labelName {
					case "httpsOnly":
						if labelValue == "true" {
							proxyConfig.Https = true
						}
						foundService = true
					case "servicePath":
						proxyConfig.ServicePath = labelValue
						foundService = true
					case "serviceDomain":
						proxyConfig.ServiceDomain = labelValue
						foundService = true
					case "port":
						foundService = true
						proxyConfig.ServicePort, _ = strconv.Atoi(labelValue)
					}
				}
			}
			if foundService {
				proxyService.ProxyConfigurations = append(proxyService.ProxyConfigurations, *proxyConfig)
			}
			proxiedServices[stackName] = append(proxiedServices[stackName], proxyService)
		}

		stacks = make([]model.Stack, len(proxiedServices))
		count := 0
		for stackName, services := range proxiedServices {

			if stackName == "" {
				continue
			}
			stack := model.Stack{
				Name: stackName,
				Services: services,
			}
			stacks[count] = stack
			count++
		}

	}

	fmt.Printf(" > Found %d stacks\n", len(stacks))


	return stacks
}
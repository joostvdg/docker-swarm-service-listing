package model

import "fmt"

// Docker stack with services
type Stack struct {
	Name     string `json:"Name"`
	Services []Service
}

// Docker service with Docker Flow Proxy configurations
type Service struct {
	Name                string `json:"Name"`
	Alias               string `json:"Alias"`
	ProxyConfigurations []ProxyConfiguration
}

// String function for printing debug info of the service and its proxy configurations
func (s *Service) String() string {
	configs := ""
	for _, config := range s.ProxyConfigurations {
		configs += fmt.Sprintf("- %s\n", config.String())
	}
	return fmt.Sprintf("%s \n %s", s.Name, configs)
}

// Docker Flow Proxy configuration
// Limited to what we want to expose
type ProxyConfiguration struct {
	Https         bool   `json:"Https"`
	ServicePath   string `json:"ServicePath"`
	ServiceDomain string `json:"ServiceDomain"`
	ServicePort   int    `json:"ServicePort"`
}

// String function for proxy config
func (pc *ProxyConfiguration) String() string {
	return fmt.Sprintf("[https=%v, path=%s, domain=%s, port=%d]", pc.Https, pc.ServicePath, pc.ServiceDomain, pc.ServicePort)
}

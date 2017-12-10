package model

import "fmt"

// Stack is a Docker stack with services
type Stack struct {
	Name     string `json:"Name"`
	Services []Service
}

// Service is a Docker service definition with Docker Flow Proxy configurations
type Service struct {
	Name                string `json:"Name"`
	Alias               string `json:"Alias"`
	ProxyConfigurations []ProxyConfiguration
}

// String for printing debug info of the service and its proxy configurations
func (s *Service) String() string {
	configs := ""
	for index, config := range s.ProxyConfigurations {
		lineEnd := "\n"
		if index == len(s.ProxyConfigurations)-1 {
			lineEnd = ""
		}
		configs += fmt.Sprintf("- %s%s", config.String(), lineEnd)
	}

	return fmt.Sprintf("%s \n %s", s.Name, configs)
}

// ProxyConfiguration is a Docker Flow Proxy configuration
// Limited to what we want to expose
type ProxyConfiguration struct {
	Https           bool   `json:"Https"`
	MainServicePath string `json:"MainServicePath"`
	ServicePath     string `json:"ServicePath"`
	ServiceDomain   string `json:"ServiceDomain"`
	ServicePort     int    `json:"ServicePort"`
}

// String is the String for proxy config
func (pc *ProxyConfiguration) String() string {
	return fmt.Sprintf("[https=%v, path=%s, domain=%s, port=%d]", pc.Https, pc.ServicePath, pc.ServiceDomain, pc.ServicePort)
}

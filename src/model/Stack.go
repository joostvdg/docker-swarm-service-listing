package model

import "fmt"

type Stack struct {
	Name string
	Services []Service
}

type Service struct {
	Name string
	Alias string
	ProxyConfigurations []ProxyConfiguration
}

func (s *Service) String() string {
	configs := ""
	for _,config := range s.ProxyConfigurations {
		configs += fmt.Sprintf("- %s\n", config.String())
	}
	return fmt.Sprintf("%s \n %s", s.Name, configs)
}

type ProxyConfiguration struct {
	Https bool
	ServicePath string
	ServiceDomain string
	ServicePort int
}

func (pc *ProxyConfiguration) String() string {
	return fmt.Sprintf("[https=%v, path=%s, domain=%s, port=%d]", pc.Https, pc.ServicePath, pc.ServiceDomain, pc.ServicePort)
}
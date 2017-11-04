package probe

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoProxiedSwarmServices(t *testing.T) {
	swarmServices := make([]swarm.Service, 0)
	proxiedServices := FindProxiedServices(swarmServices)
	assert.Equal(t, 0, len(proxiedServices))
}

func TestSimpleListOfProxiedServices(t *testing.T) {

	swarmServices := make([]swarm.Service, 1)
	swarmService1 := createSampleSwarmService()
	swarmServices[0] = swarmService1

	proxiedServices := FindProxiedServices(swarmServices)
	assert.Equal(t, 1, len(proxiedServices))
}

func createSampleSwarmService() swarm.Service {
	labels := make(map[string]string, 4)
	labels["com.df.httpsOnly"] = "true"
	labels["com.df.notify"] = "true"
	labels["com.df.distribute"] = "true"
	labels["com.df.servicePath"] = "/"
	containerSpec := swarm.ContainerSpec{
		Image: "helloworld",
	}
	taskSpec := swarm.TaskSpec{
		ContainerSpec: &containerSpec,
	}
	var numReplicas uint64
	numReplicas = 1
	replicas := swarm.ReplicatedService{
		Replicas: &numReplicas,
	}
	mode := swarm.ServiceMode{
		Replicated: &replicas,
	}
	annotations := swarm.Annotations{
		Labels: labels,
	}
	swarmService1Spec := swarm.ServiceSpec{
		Mode:         mode,
		TaskTemplate: taskSpec,
		Annotations:  annotations,
	}
	swarmService := swarm.Service{
		ID:   "abcdefg",
		Spec: swarmService1Spec,
	}
	return swarmService
}

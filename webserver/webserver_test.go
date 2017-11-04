package webserver

import (
	"../model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStacksHandlerEmptyResponse(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/stacks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	webserverData := &WebserverData{}
	handler := http.HandlerFunc(webserverData.HandleGetStacks)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	assert.Equal(t, status, http.StatusOK, fmt.Sprintf("handler returned wrong status code: got %v want %v",
		status, http.StatusOK))

	// Check the response body is what we expect.
	expected := `[]`
	actual := rr.Body.String()
	actual = strings.TrimSpace(actual)
	assert.Equal(t, expected, actual, fmt.Sprintf("handler returned unexpected body: got %v want %v",
		expected, actual))
}

func TestStacksHandlerSimpleResponse(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/stacks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	webserverData := createWebserverData()
	handler := http.HandlerFunc(webserverData.HandleGetStacks)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	assert.Equal(t, status, http.StatusOK, fmt.Sprintf("handler returned wrong status code: got %v want %v",
		status, http.StatusOK))

	// Check the response body is what we expect.
	expected := `[{"Name":"StackOne","Services":[{"Name":"ServiceOne","Alias":"","ProxyConfigurations":[{"Https":false,"ServicePath":"/demo1","ServiceDomain":"","ServicePort":0}]}]},{"Name":"StackTwo","Services":[{"Name":"ServiceTwo","Alias":"","ProxyConfigurations":[{"Https":true,"ServicePath":"/demo2","ServiceDomain":"","ServicePort":0},{"Https":true,"ServicePath":"/","ServiceDomain":"registry.example.com","ServicePort":18445}]},{"Name":"ServiceThree","Alias":"","ProxyConfigurations":[{"Https":false,"ServicePath":"/demo12","ServiceDomain":"","ServicePort":0}]}]}]`
	actual := rr.Body.String()
	actual = strings.TrimSpace(actual) // for some reason we get a unexpected \n
	assert.Equal(t, expected, actual, "They should be equal")
}
func createWebserverData() *WebserverData {
	proxyConfig1 := model.ProxyConfiguration{ServicePath: "/demo1", Https: false}
	service1ProxyConfigs := make([]model.ProxyConfiguration, 1)
	service1ProxyConfigs[0] = proxyConfig1
	service1 := model.Service{Name: "ServiceOne", Alias: "", ProxyConfigurations: service1ProxyConfigs}

	proxyConfig2 := model.ProxyConfiguration{ServicePath: "/demo2", Https: true}
	proxyConfig3 := model.ProxyConfiguration{ServicePath: "/", Https: true, ServicePort: 18445, ServiceDomain: "registry.example.com"}
	service2ProxyConfigs := make([]model.ProxyConfiguration, 2)
	service2ProxyConfigs[0] = proxyConfig2
	service2ProxyConfigs[1] = proxyConfig3
	service2 := model.Service{Name: "ServiceTwo", Alias: "", ProxyConfigurations: service2ProxyConfigs}

	proxyConfig4 := model.ProxyConfiguration{ServicePath: "/demo12", Https: false}
	service3ProxyConfigs := make([]model.ProxyConfiguration, 1)
	service3ProxyConfigs[0] = proxyConfig4
	service3 := model.Service{Name: "ServiceThree", Alias: "", ProxyConfigurations: service3ProxyConfigs}

	stack1Services := make([]model.Service, 1)
	stack1Services[0] = service1
	stack1 := model.Stack{Name: "StackOne", Services: stack1Services}

	stack2Services := make([]model.Service, 2)
	stack2Services[0] = service2
	stack2Services[1] = service3
	stack2 := model.Stack{Name: "StackTwo", Services: stack2Services}

	stacks := make([]model.Stack, 2)
	stacks[0] = stack1
	stacks[1] = stack2

	webserverData := WebserverData{Stacks: stacks}

	return &webserverData
}

// Code generated by counterfeiter. DO NOT EDIT.
package apputilsfakes

import (
	"context"
	"sync"

	"github.com/weaveworks/weave-gitops/pkg/apputils"
	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/services/app"
)

type FakeAppFactory struct {
	GetAppServiceStub        func(context.Context, string, string) (app.AppService, error)
	getAppServiceMutex       sync.RWMutex
	getAppServiceArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}
	getAppServiceReturns struct {
		result1 app.AppService
		result2 error
	}
	getAppServiceReturnsOnCall map[int]struct {
		result1 app.AppService
		result2 error
	}
	GetAppServiceForAddStub        func(context.Context, apputils.AddServiceParams) (app.AppService, error)
	getAppServiceForAddMutex       sync.RWMutex
	getAppServiceForAddArgsForCall []struct {
		arg1 context.Context
		arg2 apputils.AddServiceParams
	}
	getAppServiceForAddReturns struct {
		result1 app.AppService
		result2 error
	}
	getAppServiceForAddReturnsOnCall map[int]struct {
		result1 app.AppService
		result2 error
	}
	GetKubeServiceStub        func() (kube.Kube, error)
	getKubeServiceMutex       sync.RWMutex
	getKubeServiceArgsForCall []struct {
	}
	getKubeServiceReturns struct {
		result1 kube.Kube
		result2 error
	}
	getKubeServiceReturnsOnCall map[int]struct {
		result1 kube.Kube
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAppFactory) GetAppService(arg1 context.Context, arg2 string, arg3 string) (app.AppService, error) {
	fake.getAppServiceMutex.Lock()
	ret, specificReturn := fake.getAppServiceReturnsOnCall[len(fake.getAppServiceArgsForCall)]
	fake.getAppServiceArgsForCall = append(fake.getAppServiceArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetAppServiceStub
	fakeReturns := fake.getAppServiceReturns
	fake.recordInvocation("GetAppService", []interface{}{arg1, arg2, arg3})
	fake.getAppServiceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAppFactory) GetAppServiceCallCount() int {
	fake.getAppServiceMutex.RLock()
	defer fake.getAppServiceMutex.RUnlock()
	return len(fake.getAppServiceArgsForCall)
}

func (fake *FakeAppFactory) GetAppServiceCalls(stub func(context.Context, string, string) (app.AppService, error)) {
	fake.getAppServiceMutex.Lock()
	defer fake.getAppServiceMutex.Unlock()
	fake.GetAppServiceStub = stub
}

func (fake *FakeAppFactory) GetAppServiceArgsForCall(i int) (context.Context, string, string) {
	fake.getAppServiceMutex.RLock()
	defer fake.getAppServiceMutex.RUnlock()
	argsForCall := fake.getAppServiceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeAppFactory) GetAppServiceReturns(result1 app.AppService, result2 error) {
	fake.getAppServiceMutex.Lock()
	defer fake.getAppServiceMutex.Unlock()
	fake.GetAppServiceStub = nil
	fake.getAppServiceReturns = struct {
		result1 app.AppService
		result2 error
	}{result1, result2}
}

func (fake *FakeAppFactory) GetAppServiceReturnsOnCall(i int, result1 app.AppService, result2 error) {
	fake.getAppServiceMutex.Lock()
	defer fake.getAppServiceMutex.Unlock()
	fake.GetAppServiceStub = nil
	if fake.getAppServiceReturnsOnCall == nil {
		fake.getAppServiceReturnsOnCall = make(map[int]struct {
			result1 app.AppService
			result2 error
		})
	}
	fake.getAppServiceReturnsOnCall[i] = struct {
		result1 app.AppService
		result2 error
	}{result1, result2}
}

func (fake *FakeAppFactory) GetAppServiceForAdd(arg1 context.Context, arg2 apputils.AddServiceParams) (app.AppService, error) {
	fake.getAppServiceForAddMutex.Lock()
	ret, specificReturn := fake.getAppServiceForAddReturnsOnCall[len(fake.getAppServiceForAddArgsForCall)]
	fake.getAppServiceForAddArgsForCall = append(fake.getAppServiceForAddArgsForCall, struct {
		arg1 context.Context
		arg2 apputils.AddServiceParams
	}{arg1, arg2})
	stub := fake.GetAppServiceForAddStub
	fakeReturns := fake.getAppServiceForAddReturns
	fake.recordInvocation("GetAppServiceForAdd", []interface{}{arg1, arg2})
	fake.getAppServiceForAddMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAppFactory) GetAppServiceForAddCallCount() int {
	fake.getAppServiceForAddMutex.RLock()
	defer fake.getAppServiceForAddMutex.RUnlock()
	return len(fake.getAppServiceForAddArgsForCall)
}

func (fake *FakeAppFactory) GetAppServiceForAddCalls(stub func(context.Context, apputils.AddServiceParams) (app.AppService, error)) {
	fake.getAppServiceForAddMutex.Lock()
	defer fake.getAppServiceForAddMutex.Unlock()
	fake.GetAppServiceForAddStub = stub
}

func (fake *FakeAppFactory) GetAppServiceForAddArgsForCall(i int) (context.Context, apputils.AddServiceParams) {
	fake.getAppServiceForAddMutex.RLock()
	defer fake.getAppServiceForAddMutex.RUnlock()
	argsForCall := fake.getAppServiceForAddArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAppFactory) GetAppServiceForAddReturns(result1 app.AppService, result2 error) {
	fake.getAppServiceForAddMutex.Lock()
	defer fake.getAppServiceForAddMutex.Unlock()
	fake.GetAppServiceForAddStub = nil
	fake.getAppServiceForAddReturns = struct {
		result1 app.AppService
		result2 error
	}{result1, result2}
}

func (fake *FakeAppFactory) GetAppServiceForAddReturnsOnCall(i int, result1 app.AppService, result2 error) {
	fake.getAppServiceForAddMutex.Lock()
	defer fake.getAppServiceForAddMutex.Unlock()
	fake.GetAppServiceForAddStub = nil
	if fake.getAppServiceForAddReturnsOnCall == nil {
		fake.getAppServiceForAddReturnsOnCall = make(map[int]struct {
			result1 app.AppService
			result2 error
		})
	}
	fake.getAppServiceForAddReturnsOnCall[i] = struct {
		result1 app.AppService
		result2 error
	}{result1, result2}
}

func (fake *FakeAppFactory) GetKubeService() (kube.Kube, error) {
	fake.getKubeServiceMutex.Lock()
	ret, specificReturn := fake.getKubeServiceReturnsOnCall[len(fake.getKubeServiceArgsForCall)]
	fake.getKubeServiceArgsForCall = append(fake.getKubeServiceArgsForCall, struct {
	}{})
	stub := fake.GetKubeServiceStub
	fakeReturns := fake.getKubeServiceReturns
	fake.recordInvocation("GetKubeService", []interface{}{})
	fake.getKubeServiceMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAppFactory) GetKubeServiceCallCount() int {
	fake.getKubeServiceMutex.RLock()
	defer fake.getKubeServiceMutex.RUnlock()
	return len(fake.getKubeServiceArgsForCall)
}

func (fake *FakeAppFactory) GetKubeServiceCalls(stub func() (kube.Kube, error)) {
	fake.getKubeServiceMutex.Lock()
	defer fake.getKubeServiceMutex.Unlock()
	fake.GetKubeServiceStub = stub
}

func (fake *FakeAppFactory) GetKubeServiceReturns(result1 kube.Kube, result2 error) {
	fake.getKubeServiceMutex.Lock()
	defer fake.getKubeServiceMutex.Unlock()
	fake.GetKubeServiceStub = nil
	fake.getKubeServiceReturns = struct {
		result1 kube.Kube
		result2 error
	}{result1, result2}
}

func (fake *FakeAppFactory) GetKubeServiceReturnsOnCall(i int, result1 kube.Kube, result2 error) {
	fake.getKubeServiceMutex.Lock()
	defer fake.getKubeServiceMutex.Unlock()
	fake.GetKubeServiceStub = nil
	if fake.getKubeServiceReturnsOnCall == nil {
		fake.getKubeServiceReturnsOnCall = make(map[int]struct {
			result1 kube.Kube
			result2 error
		})
	}
	fake.getKubeServiceReturnsOnCall[i] = struct {
		result1 kube.Kube
		result2 error
	}{result1, result2}
}

func (fake *FakeAppFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getAppServiceMutex.RLock()
	defer fake.getAppServiceMutex.RUnlock()
	fake.getAppServiceForAddMutex.RLock()
	defer fake.getAppServiceForAddMutex.RUnlock()
	fake.getKubeServiceMutex.RLock()
	defer fake.getKubeServiceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAppFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ apputils.AppFactory = new(FakeAppFactory)

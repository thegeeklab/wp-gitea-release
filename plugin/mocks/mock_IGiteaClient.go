// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	io "io"

	gitea "code.gitea.io/sdk/gitea"

	mock "github.com/stretchr/testify/mock"
)

// MockIGiteaClient is an autogenerated mock type for the IGiteaClient type
type MockIGiteaClient struct {
	mock.Mock
}

type MockIGiteaClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIGiteaClient) EXPECT() *MockIGiteaClient_Expecter {
	return &MockIGiteaClient_Expecter{mock: &_m.Mock}
}

// CreateRelease provides a mock function with given fields: owner, repo, opt
func (_m *MockIGiteaClient) CreateRelease(owner string, repo string, opt gitea.CreateReleaseOption) (*gitea.Release, *gitea.Response, error) {
	ret := _m.Called(owner, repo, opt)

	if len(ret) == 0 {
		panic("no return value specified for CreateRelease")
	}

	var r0 *gitea.Release
	var r1 *gitea.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, gitea.CreateReleaseOption) (*gitea.Release, *gitea.Response, error)); ok {
		return rf(owner, repo, opt)
	}
	if rf, ok := ret.Get(0).(func(string, string, gitea.CreateReleaseOption) *gitea.Release); ok {
		r0 = rf(owner, repo, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gitea.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, gitea.CreateReleaseOption) *gitea.Response); ok {
		r1 = rf(owner, repo, opt)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitea.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, string, gitea.CreateReleaseOption) error); ok {
		r2 = rf(owner, repo, opt)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockIGiteaClient_CreateRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRelease'
type MockIGiteaClient_CreateRelease_Call struct {
	*mock.Call
}

// CreateRelease is a helper method to define mock.On call
//   - owner string
//   - repo string
//   - opt gitea.CreateReleaseOption
func (_e *MockIGiteaClient_Expecter) CreateRelease(owner interface{}, repo interface{}, opt interface{}) *MockIGiteaClient_CreateRelease_Call {
	return &MockIGiteaClient_CreateRelease_Call{Call: _e.mock.On("CreateRelease", owner, repo, opt)}
}

func (_c *MockIGiteaClient_CreateRelease_Call) Run(run func(owner string, repo string, opt gitea.CreateReleaseOption)) *MockIGiteaClient_CreateRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(gitea.CreateReleaseOption))
	})
	return _c
}

func (_c *MockIGiteaClient_CreateRelease_Call) Return(_a0 *gitea.Release, _a1 *gitea.Response, _a2 error) *MockIGiteaClient_CreateRelease_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockIGiteaClient_CreateRelease_Call) RunAndReturn(run func(string, string, gitea.CreateReleaseOption) (*gitea.Release, *gitea.Response, error)) *MockIGiteaClient_CreateRelease_Call {
	_c.Call.Return(run)
	return _c
}

// CreateReleaseAttachment provides a mock function with given fields: user, repo, release, file, filename
func (_m *MockIGiteaClient) CreateReleaseAttachment(user string, repo string, release int64, file io.Reader, filename string) (*gitea.Attachment, *gitea.Response, error) {
	ret := _m.Called(user, repo, release, file, filename)

	if len(ret) == 0 {
		panic("no return value specified for CreateReleaseAttachment")
	}

	var r0 *gitea.Attachment
	var r1 *gitea.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int64, io.Reader, string) (*gitea.Attachment, *gitea.Response, error)); ok {
		return rf(user, repo, release, file, filename)
	}
	if rf, ok := ret.Get(0).(func(string, string, int64, io.Reader, string) *gitea.Attachment); ok {
		r0 = rf(user, repo, release, file, filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gitea.Attachment)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, int64, io.Reader, string) *gitea.Response); ok {
		r1 = rf(user, repo, release, file, filename)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitea.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, string, int64, io.Reader, string) error); ok {
		r2 = rf(user, repo, release, file, filename)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockIGiteaClient_CreateReleaseAttachment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateReleaseAttachment'
type MockIGiteaClient_CreateReleaseAttachment_Call struct {
	*mock.Call
}

// CreateReleaseAttachment is a helper method to define mock.On call
//   - user string
//   - repo string
//   - release int64
//   - file io.Reader
//   - filename string
func (_e *MockIGiteaClient_Expecter) CreateReleaseAttachment(user interface{}, repo interface{}, release interface{}, file interface{}, filename interface{}) *MockIGiteaClient_CreateReleaseAttachment_Call {
	return &MockIGiteaClient_CreateReleaseAttachment_Call{Call: _e.mock.On("CreateReleaseAttachment", user, repo, release, file, filename)}
}

func (_c *MockIGiteaClient_CreateReleaseAttachment_Call) Run(run func(user string, repo string, release int64, file io.Reader, filename string)) *MockIGiteaClient_CreateReleaseAttachment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(int64), args[3].(io.Reader), args[4].(string))
	})
	return _c
}

func (_c *MockIGiteaClient_CreateReleaseAttachment_Call) Return(_a0 *gitea.Attachment, _a1 *gitea.Response, _a2 error) *MockIGiteaClient_CreateReleaseAttachment_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockIGiteaClient_CreateReleaseAttachment_Call) RunAndReturn(run func(string, string, int64, io.Reader, string) (*gitea.Attachment, *gitea.Response, error)) *MockIGiteaClient_CreateReleaseAttachment_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteReleaseAttachment provides a mock function with given fields: user, repo, release, id
func (_m *MockIGiteaClient) DeleteReleaseAttachment(user string, repo string, release int64, id int64) (*gitea.Response, error) {
	ret := _m.Called(user, repo, release, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteReleaseAttachment")
	}

	var r0 *gitea.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, int64, int64) (*gitea.Response, error)); ok {
		return rf(user, repo, release, id)
	}
	if rf, ok := ret.Get(0).(func(string, string, int64, int64) *gitea.Response); ok {
		r0 = rf(user, repo, release, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gitea.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, int64, int64) error); ok {
		r1 = rf(user, repo, release, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIGiteaClient_DeleteReleaseAttachment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteReleaseAttachment'
type MockIGiteaClient_DeleteReleaseAttachment_Call struct {
	*mock.Call
}

// DeleteReleaseAttachment is a helper method to define mock.On call
//   - user string
//   - repo string
//   - release int64
//   - id int64
func (_e *MockIGiteaClient_Expecter) DeleteReleaseAttachment(user interface{}, repo interface{}, release interface{}, id interface{}) *MockIGiteaClient_DeleteReleaseAttachment_Call {
	return &MockIGiteaClient_DeleteReleaseAttachment_Call{Call: _e.mock.On("DeleteReleaseAttachment", user, repo, release, id)}
}

func (_c *MockIGiteaClient_DeleteReleaseAttachment_Call) Run(run func(user string, repo string, release int64, id int64)) *MockIGiteaClient_DeleteReleaseAttachment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(int64), args[3].(int64))
	})
	return _c
}

func (_c *MockIGiteaClient_DeleteReleaseAttachment_Call) Return(_a0 *gitea.Response, _a1 error) *MockIGiteaClient_DeleteReleaseAttachment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIGiteaClient_DeleteReleaseAttachment_Call) RunAndReturn(run func(string, string, int64, int64) (*gitea.Response, error)) *MockIGiteaClient_DeleteReleaseAttachment_Call {
	_c.Call.Return(run)
	return _c
}

// ListReleaseAttachments provides a mock function with given fields: user, repo, release, opt
func (_m *MockIGiteaClient) ListReleaseAttachments(user string, repo string, release int64, opt gitea.ListReleaseAttachmentsOptions) ([]*gitea.Attachment, *gitea.Response, error) {
	ret := _m.Called(user, repo, release, opt)

	if len(ret) == 0 {
		panic("no return value specified for ListReleaseAttachments")
	}

	var r0 []*gitea.Attachment
	var r1 *gitea.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int64, gitea.ListReleaseAttachmentsOptions) ([]*gitea.Attachment, *gitea.Response, error)); ok {
		return rf(user, repo, release, opt)
	}
	if rf, ok := ret.Get(0).(func(string, string, int64, gitea.ListReleaseAttachmentsOptions) []*gitea.Attachment); ok {
		r0 = rf(user, repo, release, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gitea.Attachment)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, int64, gitea.ListReleaseAttachmentsOptions) *gitea.Response); ok {
		r1 = rf(user, repo, release, opt)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitea.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, string, int64, gitea.ListReleaseAttachmentsOptions) error); ok {
		r2 = rf(user, repo, release, opt)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockIGiteaClient_ListReleaseAttachments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListReleaseAttachments'
type MockIGiteaClient_ListReleaseAttachments_Call struct {
	*mock.Call
}

// ListReleaseAttachments is a helper method to define mock.On call
//   - user string
//   - repo string
//   - release int64
//   - opt gitea.ListReleaseAttachmentsOptions
func (_e *MockIGiteaClient_Expecter) ListReleaseAttachments(user interface{}, repo interface{}, release interface{}, opt interface{}) *MockIGiteaClient_ListReleaseAttachments_Call {
	return &MockIGiteaClient_ListReleaseAttachments_Call{Call: _e.mock.On("ListReleaseAttachments", user, repo, release, opt)}
}

func (_c *MockIGiteaClient_ListReleaseAttachments_Call) Run(run func(user string, repo string, release int64, opt gitea.ListReleaseAttachmentsOptions)) *MockIGiteaClient_ListReleaseAttachments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(int64), args[3].(gitea.ListReleaseAttachmentsOptions))
	})
	return _c
}

func (_c *MockIGiteaClient_ListReleaseAttachments_Call) Return(_a0 []*gitea.Attachment, _a1 *gitea.Response, _a2 error) *MockIGiteaClient_ListReleaseAttachments_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockIGiteaClient_ListReleaseAttachments_Call) RunAndReturn(run func(string, string, int64, gitea.ListReleaseAttachmentsOptions) ([]*gitea.Attachment, *gitea.Response, error)) *MockIGiteaClient_ListReleaseAttachments_Call {
	_c.Call.Return(run)
	return _c
}

// ListReleases provides a mock function with given fields: owner, repo, opt
func (_m *MockIGiteaClient) ListReleases(owner string, repo string, opt gitea.ListReleasesOptions) ([]*gitea.Release, *gitea.Response, error) {
	ret := _m.Called(owner, repo, opt)

	if len(ret) == 0 {
		panic("no return value specified for ListReleases")
	}

	var r0 []*gitea.Release
	var r1 *gitea.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, gitea.ListReleasesOptions) ([]*gitea.Release, *gitea.Response, error)); ok {
		return rf(owner, repo, opt)
	}
	if rf, ok := ret.Get(0).(func(string, string, gitea.ListReleasesOptions) []*gitea.Release); ok {
		r0 = rf(owner, repo, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gitea.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, gitea.ListReleasesOptions) *gitea.Response); ok {
		r1 = rf(owner, repo, opt)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitea.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, string, gitea.ListReleasesOptions) error); ok {
		r2 = rf(owner, repo, opt)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockIGiteaClient_ListReleases_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListReleases'
type MockIGiteaClient_ListReleases_Call struct {
	*mock.Call
}

// ListReleases is a helper method to define mock.On call
//   - owner string
//   - repo string
//   - opt gitea.ListReleasesOptions
func (_e *MockIGiteaClient_Expecter) ListReleases(owner interface{}, repo interface{}, opt interface{}) *MockIGiteaClient_ListReleases_Call {
	return &MockIGiteaClient_ListReleases_Call{Call: _e.mock.On("ListReleases", owner, repo, opt)}
}

func (_c *MockIGiteaClient_ListReleases_Call) Run(run func(owner string, repo string, opt gitea.ListReleasesOptions)) *MockIGiteaClient_ListReleases_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(gitea.ListReleasesOptions))
	})
	return _c
}

func (_c *MockIGiteaClient_ListReleases_Call) Return(_a0 []*gitea.Release, _a1 *gitea.Response, _a2 error) *MockIGiteaClient_ListReleases_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockIGiteaClient_ListReleases_Call) RunAndReturn(run func(string, string, gitea.ListReleasesOptions) ([]*gitea.Release, *gitea.Response, error)) *MockIGiteaClient_ListReleases_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIGiteaClient creates a new instance of MockIGiteaClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIGiteaClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIGiteaClient {
	mock := &MockIGiteaClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
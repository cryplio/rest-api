package mockrouter

import (
	matcher "github.com/Nivl/gomock-type-matcher"
	gomock "github.com/golang/mock/gomock"
)

// CreatedSuccess is a helper that expects a valid Created response
func (mr *MockHTTPResponseMockRecorder) CreatedSuccess(typ interface{}, runnable interface{}) *gomock.Call {
	createdCall := mr.Created(matcher.Interface(typ))
	createdCall.Return(nil)
	if runnable != nil {
		createdCall.Do(runnable)
	}
	return createdCall
}

// OkSuccess is a helper that expects a valid Created response
func (mr *MockHTTPResponseMockRecorder) OkSuccess(typ interface{}, runnable interface{}) *gomock.Call {
	createdCall := mr.Ok(matcher.Interface(typ))
	createdCall.Return(nil)
	if runnable != nil {
		createdCall.Do(runnable)
	}
	return createdCall
}

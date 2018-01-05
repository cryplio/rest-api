package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryplio/rest-api/modules/api"
	"github.com/golang/mock/gomock"

	"github.com/Nivl/go-filestorage/implementations/mockfilestorage"
	"github.com/Nivl/go-logger/implementations/mocklogger"
	"github.com/Nivl/go-mailer/implementations/mockmailer"
	"github.com/Nivl/go-reporter/implementations/mockreporter"
	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	matcher "github.com/Nivl/gomock-type-matcher"

	"github.com/Nivl/go-rest-tools/dependencies/mockdependencies"
	"github.com/stretchr/testify/assert"
)

// Test that an un-existing route returns JSON and a 404
func TestRouteNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	apiDeps := mockdependencies.NewMockDependencies(mockCtrl)
	apiDeps.EXPECT().NewFileStorage(gomock.Any()).Return(mockfilestorage.NewMockFileStorage(mockCtrl), nil)
	apiDeps.EXPECT().DB().Return(mocksqldb.NewMockConnection(mockCtrl))
	apiDeps.EXPECT().Mailer().Return(mockmailer.NewMockMailer(mockCtrl), nil)

	reporter := mockreporter.NewMockReporter(mockCtrl)
	reporter.EXPECT().AddTag("Req ID", matcher.Type("string"))
	reporter.EXPECT().AddTag("Endpoint", matcher.Type("string"))

	apiDeps.EXPECT().NewReporter().Return(reporter, nil)

	logger := mocklogger.NewMockLogger(mockCtrl)
	logger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	apiDeps.EXPECT().NewLogger().Return(logger, nil)

	rec := httptest.NewRecorder()
	api.GetRouter(apiDeps).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
}

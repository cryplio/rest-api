package portfolios

// Code generated; DO NOT EDIT.

import (
	"strings"

	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-sqldb/implementations/mocksqldb"

	gomock "github.com/golang/mock/gomock"

	"github.com/Nivl/go-types/datetime"
)

func TestJoinPortfolioSQL(t *testing.T) {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "user_id", "name" }
	totalFields := len(fields)
	output := JoinPortfolioSQL("tofind")

	assert.Equal(t, totalFields*2, strings.Count(output, "tofind."), "wrong number of fields returned")
	assert.True(t, strings.HasSuffix(output, "\""), "JoinSQL() output should end with a \"")
}

func TestGetPortfolioByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedID := "4408d5e1-b510-42cb-8ff8-788948a246dd"
	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().GetID(&Portfolio{}, expectedID, nil)

	_, err := GetPortfolioByID(mockDB, expectedID)
	assert.NoError(t, err, "GetPortfolioByID() should not have failed")
}

func TestGetAnyPortfolioByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedID := "4408d5e1-b510-42cb-8ff8-788948a246dd"
	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().GetID(&Portfolio{}, expectedID, nil)

	_, err := GetPortfolioByID(mockDB, expectedID)
	assert.NoError(t, err, "GetPortfolioByID() should not have failed")
}

func TestPortfolioSaveNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Portfolio{})

	p := &Portfolio{}
	err := p.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
}

func TestPortfolioSaveExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Portfolio{})

	p := &Portfolio{}
	id := uuid.NewV4().String()
	p.ID = id
	err := p.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	assert.Equal(t, id, p.ID, "ID should not have changed")
}

func TestPortfolioCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Portfolio{})

	p := &Portfolio{}
	err := p.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestPortfolioCreateWithID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)

	p := &Portfolio{}
	p.ID = uuid.NewV4().String()

	err := p.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
}

func TestPortfolioDoCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Portfolio{})

	p := &Portfolio{}
	err := p.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestPortfolioDoCreateWithDate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Portfolio{})

	createdAt := datetime.Now().AddDate(0, 0, 1)
	p := &Portfolio{CreatedAt: createdAt}
	err := p.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.True(t, p.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestPortfolioDoCreateFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertError(&Portfolio{}, errors.New("sql error"))

	p := &Portfolio{}
	err := p.doCreate(mockDB)

	assert.Error(t, err, "doCreate() should have fail")
}


func TestPortfolioUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Portfolio{})

	p := &Portfolio{}
	p.ID = uuid.NewV4().String()
	err := p.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestPortfolioUpdateWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	p := &Portfolio{}
	err := p.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
}


func TestPortfolioDoUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Portfolio{})

	p := &Portfolio{}
	p.ID = uuid.NewV4().String()
	err := p.doUpdate(mockDB)

	assert.NoError(t, err, "doUpdate() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestPortfolioDoUpdateWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	p := &Portfolio{}
	err := p.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should not have fail")
}

func TestPortfolioDoUpdateFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateError(&Portfolio{}, errors.New("sql error"))

	p := &Portfolio{}
	p.ID = uuid.NewV4().String()
	err := p.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should have fail")
}

func TestPortfolioDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().DeletionSuccess()

	p := &Portfolio{}
	p.ID = uuid.NewV4().String()
	err := p.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
}

func TestPortfolioDeleteWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	p := &Portfolio{}
	err := p.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
}

func TestPortfolioDeleteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().DeletionError(errors.New("sql error"))

	p := &Portfolio{}
	p.ID = uuid.NewV4().String()
	err := p.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
}

func TestPortfolioGetID(t *testing.T) {
	p := &Portfolio{}
	p.ID = uuid.NewV4().String()
	assert.Equal(t, p.ID, p.GetID(), "GetID() did not return the right ID")
}

func TestPortfolioSetID(t *testing.T) {
	p := &Portfolio{}
	p.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, p.ID, "SetID() did not set the ID")
}

func TestPortfolioIsZero(t *testing.T) {
	empty := &Portfolio{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *Portfolio
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &Portfolio{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}
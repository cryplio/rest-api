package tickers

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

func TestJoinSQL(t *testing.T) {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "name", "symbol", "unit", "marketcap", "volume_24h", "max_supply", "current_supply", "logo_url", "website", "price_usd", "percent_change_1h", "percent_change_24h", "percent_change_7d", "coinmarketcap_id" }
	totalFields := len(fields)
	output := JoinSQL("tofind")

	assert.Equal(t, totalFields*2, strings.Count(output, "tofind."), "wrong number of fields returned")
	assert.True(t, strings.HasSuffix(output, "\""), "JoinSQL() output should end with a \"")
}

func TestGetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedID := "4408d5e1-b510-42cb-8ff8-788948a246dd"
	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().GetID(&Ticker{}, expectedID, nil)

	_, err := GetByID(mockDB, expectedID)
	assert.NoError(t, err, "GetByID() should not have failed")
}

func TestGetAnyByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedID := "4408d5e1-b510-42cb-8ff8-788948a246dd"
	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().GetID(&Ticker{}, expectedID, nil)

	_, err := GetByID(mockDB, expectedID)
	assert.NoError(t, err, "GetByID() should not have failed")
}

func TestTickerSaveNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Ticker{})

	t := &Ticker{}
	err := t.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	assert.NotEmpty(t, t.ID, "ID should have been set")
}

func TestTickerSaveExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Ticker{})

	t := &Ticker{}
	id := uuid.NewV4().String()
	t.ID = id
	err := t.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	assert.Equal(t, id, t.ID, "ID should not have changed")
}

func TestTickerCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Ticker{})

	t := &Ticker{}
	err := t.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	assert.NotEmpty(t, t.ID, "ID should have been set")
	assert.NotNil(t, t.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, t.UpdatedAt, "UpdatedAt should have been set")
}

func TestTickerCreateWithID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)

	t := &Ticker{}
	t.ID = uuid.NewV4().String()

	err := t.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
}

func TestTickerDoCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Ticker{})

	t := &Ticker{}
	err := t.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	assert.NotEmpty(t, t.ID, "ID should have been set")
	assert.NotNil(t, t.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, t.UpdatedAt, "UpdatedAt should have been set")
}

func TestTickerDoCreateWithDate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Ticker{})

	createdAt := datetime.Now().AddDate(0, 0, 1)
	t := &Ticker{CreatedAt: createdAt}
	err := t.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	assert.NotEmpty(t, t.ID, "ID should have been set")
	assert.True(t, t.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, t.UpdatedAt, "UpdatedAt should have been set")
}

func TestTickerDoCreateFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertError(&Ticker{}, errors.New("sql error"))

	t := &Ticker{}
	err := t.doCreate(mockDB)

	assert.Error(t, err, "doCreate() should have fail")
}


func TestTickerUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Ticker{})

	t := &Ticker{}
	t.ID = uuid.NewV4().String()
	err := t.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	assert.NotEmpty(t, t.ID, "ID should have been set")
	assert.NotNil(t, t.UpdatedAt, "UpdatedAt should have been set")
}

func TestTickerUpdateWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	t := &Ticker{}
	err := t.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
}


func TestTickerDoUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Ticker{})

	t := &Ticker{}
	t.ID = uuid.NewV4().String()
	err := t.doUpdate(mockDB)

	assert.NoError(t, err, "doUpdate() should not have fail")
	assert.NotEmpty(t, t.ID, "ID should have been set")
	assert.NotNil(t, t.UpdatedAt, "UpdatedAt should have been set")
}

func TestTickerDoUpdateWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	t := &Ticker{}
	err := t.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should not have fail")
}

func TestTickerDoUpdateFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateError(&Ticker{}, errors.New("sql error"))

	t := &Ticker{}
	t.ID = uuid.NewV4().String()
	err := t.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should have fail")
}

func TestTickerDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().DeletionSuccess()

	t := &Ticker{}
	t.ID = uuid.NewV4().String()
	err := t.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
}

func TestTickerDeleteWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	t := &Ticker{}
	err := t.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
}

func TestTickerDeleteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().DeletionError(errors.New("sql error"))

	t := &Ticker{}
	t.ID = uuid.NewV4().String()
	err := t.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
}

func TestTickerGetID(t *testing.T) {
	t := &Ticker{}
	t.ID = uuid.NewV4().String()
	assert.Equal(t, t.ID, t.GetID(), "GetID() did not return the right ID")
}

func TestTickerSetID(t *testing.T) {
	t := &Ticker{}
	t.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, t.ID, "SetID() did not set the ID")
}

func TestTickerIsZero(t *testing.T) {
	empty := &Ticker{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *Ticker
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &Ticker{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}
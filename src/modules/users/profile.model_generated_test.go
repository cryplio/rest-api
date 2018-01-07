package users

// Code generated; DO NOT EDIT.

import (
	

	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-sqldb/implementations/mocksqldb"

	gomock "github.com/golang/mock/gomock"

	"github.com/Nivl/go-types/datetime"
)







func TestProfileSaveNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Profile{})

	p := &Profile{}
	err := p.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
}

func TestProfileSaveExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Profile{})

	p := &Profile{}
	id := uuid.NewV4().String()
	p.ID = id
	err := p.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	assert.Equal(t, id, p.ID, "ID should not have changed")
}

func TestProfileCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Profile{})

	p := &Profile{}
	err := p.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileCreateWithID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)

	p := &Profile{}
	p.ID = uuid.NewV4().String()

	err := p.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
}

func TestProfileDoCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Profile{})

	p := &Profile{}
	err := p.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileDoCreateWithDate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertSuccess(&Profile{})

	createdAt := datetime.Now().AddDate(0, 0, 1)
	p := &Profile{CreatedAt: createdAt}
	err := p.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.True(t, p.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileDoCreateFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().InsertError(&Profile{}, errors.New("sql error"))

	p := &Profile{}
	err := p.doCreate(mockDB)

	assert.Error(t, err, "doCreate() should have fail")
}


func TestProfileUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Profile{})

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileUpdateWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	p := &Profile{}
	err := p.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
}


func TestProfileDoUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateSuccess(&Profile{})

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.doUpdate(mockDB)

	assert.NoError(t, err, "doUpdate() should not have fail")
	assert.NotEmpty(t, p.ID, "ID should have been set")
	assert.NotNil(t, p.UpdatedAt, "UpdatedAt should have been set")
}

func TestProfileDoUpdateWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	p := &Profile{}
	err := p.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should not have fail")
}

func TestProfileDoUpdateFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().UpdateError(&Profile{}, errors.New("sql error"))

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should have fail")
}

func TestProfileDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().DeletionSuccess()

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
}

func TestProfileDeleteWithoutID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	p := &Profile{}
	err := p.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
}

func TestProfileDeleteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocksqldb.NewMockQueryable(mockCtrl)
	mockDB.EXPECT().DeletionError(errors.New("sql error"))

	p := &Profile{}
	p.ID = uuid.NewV4().String()
	err := p.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
}

func TestProfileGetID(t *testing.T) {
	p := &Profile{}
	p.ID = uuid.NewV4().String()
	assert.Equal(t, p.ID, p.GetID(), "GetID() did not return the right ID")
}

func TestProfileSetID(t *testing.T) {
	p := &Profile{}
	p.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, p.ID, "SetID() did not set the ID")
}

func TestProfileIsZero(t *testing.T) {
	empty := &Profile{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *Profile
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &Profile{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}
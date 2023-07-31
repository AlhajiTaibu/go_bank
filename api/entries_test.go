package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/AlhajiTaibu/simplebank/mock"
	db "github.com/AlhajiTaibu/simplebank/sqlc"
	"github.com/AlhajiTaibu/simplebank/util"
	"go.uber.org/mock/gomock"
)


func randomEntries(account db.Account) db.Entry{
	// account := randomAccount()

	return db.Entry{
		ID: util.RandomInt(1, 1000),
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
}

func TestGetEntry(t *testing.T){
	account := randomAccount()
	entry := randomEntries(account)

	ctrl := gomock.NewController(t)
	store := mock.NewMockStore(ctrl)

	store.EXPECT().GetEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(entry, nil)

	url := fmt.Sprintf("/entries/%v", entry.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	require.NoError(t, err)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	server.router.ServeHTTP(recorder, request)

	requireResponseBodyMatchForEntries(t, recorder.Body, entry)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func requireResponseBodyMatchForEntries(t *testing.T, body *bytes.Buffer, entries db.Entry){
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getEntry db.Entry
	err = json.Unmarshal(data, &getEntry)
	require.NoError(t, err)
	require.Equal(t, getEntry, entries)
}

func TestCreateEntry(t *testing.T){
	account := randomAccount()
	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	args := db.CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	store.EXPECT().CreateEntry(gomock.Any(), gomock.Eq(args)).Times(1).Return(db.Entry{}, nil)
	body, err := convertArgsToJson(args)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, "/entries", body)
	recorder := httptest.NewRecorder()
	require.NoError(t, err)

	server := NewServer(store)
	server.router.ServeHTTP(recorder, request)

	requireResponseBodyMatchForEntries(t, recorder.Body, db.Entry{})
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestDeleteEntry(t *testing.T){
	account := randomAccount()
	entry := randomEntries(account)

	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	store.EXPECT().DeleteEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(nil)

	url := fmt.Sprintf("/entries/%v", entry.ID)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	server := NewServer(store)
	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestUpdateEntry(t *testing.T){
	account := randomAccount()
	entry := randomEntries(account)

	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	args := db.UpdateEntryParams{
		ID: entry.ID,
		AccountID: entry.AccountID,
		Amount: util.RandomMoney(),
	}

	store.EXPECT().GetEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(entry, nil)

	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(entry.AccountID)).Times(1).Return(account, nil)

	store.EXPECT().UpdateEntry(gomock.Any(), gomock.Eq(args)).Times(1).Return(entry, nil)

	url := fmt.Sprintf("/entries/%v", entry.ID)
	body, err := convertArgsToJson(args)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPut, url, body)
	recorder := httptest.NewRecorder()
	require.NoError(t, err)

	server := NewServer(store)
	server.router.ServeHTTP(recorder, request)

	requireResponseBodyMatchForEntries(t, recorder.Body, entry)
	require.Equal(t, http.StatusOK, recorder.Code)
}
package api

import (
	"bytes"
	"database/sql"
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


func randomEntries() db.Entry{
	account := randomAccount()

	return db.Entry{
		ID: util.RandomInt(1, 1000),
		AccountID: account.ID,
		Amount: sql.NullInt64{Int64: util.RandomMoney(), Valid: true},
	}
}

func TestGetEntry(t *testing.T){
	entry := randomEntries()

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
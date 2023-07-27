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

func randomTransfer() db.Transfer {
	fromAccount := randomAccount()
	toAccount := randomAccount()
	return db.Transfer{
		ID: util.RandomInt(1,1000),
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: sql.NullInt64{Int64: util.RandomMoney(), Valid: true},
	}
}

func TestGetTransfer(t *testing.T){
	transfer := randomTransfer()

	ctrl := gomock.NewController(t)
	store := mock.NewMockStore(ctrl)
	
	store.EXPECT().GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(transfer, nil)

	url := fmt.Sprintf("/transfers/%v", transfer.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)	
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	server := NewServer(store)

	server.router.ServeHTTP(recorder, request)

	require.NoError(t, err)
	requireResponseBodyMatchForTransfer(t, recorder.Body, transfer)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func requireResponseBodyMatchForTransfer(t *testing.T, body *bytes.Buffer, transfer db.Transfer){
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getTransfer db.Transfer
	err = json.Unmarshal(data, &getTransfer)
	require.NoError(t, err)
	require.Equal(t, getTransfer, transfer)
}
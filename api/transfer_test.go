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

func randomTransfer() db.Transfer {
	fromAccount := randomAccount()
	toAccount := randomAccount()
	return db.Transfer{
		ID: util.RandomInt(1,1000),
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: util.RandomMoney(),
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

func requireResponseBodyMatcherForTx(t *testing.T, body *bytes.Buffer, transfer db.TransferTxResult){
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getTransfer db.TransferTxResult
	err = json.Unmarshal(data, &getTransfer)
	require.NoError(t, err)
	require.Equal(t, getTransfer, transfer)
}


func TestCreateTransfer(t *testing.T){
	fromAccount := randomAccount()
	toAccount := randomAccount()

	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	args := db.TransferMoneyTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: util.RandomMoney(),
	}

	result := db.TransferTxResult{
		Transfer: db.Transfer{},
		FromAccount: fromAccount,
		ToAccount: toAccount,
		FromEntry: db.Entry{},
		ToEntry: db.Entry{},
	}

	store.EXPECT().TransferMoneyTx(gomock.Any(), gomock.Eq(args)).Times(1).Return(result, nil)

	url := "/transfers"

	body, err := convertArgsToJson(args) 

	require.NoError(t, err)
	
	request, err := http.NewRequest(http.MethodPost, url, body)
	recorder := httptest.NewRecorder()
	server := NewServer(store)
	server.router.ServeHTTP(recorder, request)

	require.NoError(t, err)
	requireResponseBodyMatcherForTx(t, recorder.Body, result)
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestDeleteTransfer(t *testing.T){
	transfer := randomTransfer()

	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	store.EXPECT().DeleteTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(nil)

	url := fmt.Sprintf("/transfers/%v", transfer.ID)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	server := NewServer(store)
	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

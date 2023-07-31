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

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomMoney(),
	}
}

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	// build a stub
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%v", account.ID)

	request, err := http.NewRequest(http.MethodGet, url, nil)

	server.router.ServeHTTP(recorder, request)

	require.NoError(t, err)
	requireResponseBodyMatch(t, recorder.Body, account)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func requireResponseBodyMatch(t *testing.T, body *bytes.Buffer, account db.Account){
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getAccount db.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)
	require.Equal(t, getAccount, account)
}

func TestCreateAccount(t *testing.T){
	account := db.Account{}
	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	args := db.CreateAccountParams{
		Owner: util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance: 0,
	}

	store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(args)).Times(1).Return(account, nil)

	url := "/accounts"

	body, err := convertArgsToJson(args) 
	require.NoError(t, err)
	
	request, err := http.NewRequest(http.MethodPost, url, body)
	recorder := httptest.NewRecorder()
	server := NewServer(store)
	server.router.ServeHTTP(recorder, request)

	require.NoError(t, err)
	requireResponseBodyMatch(t, recorder.Body, account)
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func convertArgsToJson(args interface{}) (*bytes.Reader, error){
	json, err := json.Marshal(args)
	fmt.Println("hola:", args)
	data := bytes.NewReader(json)
	return data, err
}

func TestDeleteAccount(t *testing.T){
	account := randomAccount()

	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)

	store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(nil)

	url := fmt.Sprintf("/accounts/%v", account.ID)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	server := NewServer(store)
	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestUpdateAccount(t *testing.T){
	account := randomAccount()

	ctrl := gomock.NewController(t)

	store := mock.NewMockStore(ctrl)
	
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	args := db.UpdateAccountParams{
		ID: account.ID,
		Owner: account.Owner,
		Currency: account.Currency,
		Balance: util.RandomMoney(),
	}

	store.EXPECT().UpdateAccount(gomock.Any(), gomock.Eq(args)).Times(1).Return(account, nil)
	url := fmt.Sprintf("/accounts/%v", account.ID)
	body, err := convertArgsToJson(args)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPut, url, body)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	server := NewServer(store)

	server.router.ServeHTTP(recorder, request)
	requireResponseBodyMatch(t, recorder.Body, account)
	require.Equal(t, http.StatusOK, recorder.Code)

}
//  *@createTime    2022/3/20 22:26
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mockdb "github.com/haylove/small_bank/db/mock"
	db "github.com/haylove/small_bank/db/sqlc"
	"github.com/haylove/small_bank/util"
)

func TestServer_getAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				responseMatchAccount(t, recorder, account)
			},
		},
		{
			name:      "NotFount",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Invalid ID",
			accountID: 0,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)
			//start httptest server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request := httptest.NewRequest(http.MethodGet, url, nil)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       int64(rand.Intn(1000)),
		Owner:    util.RandOwner(),
		Balance:  util.RandMoney(),
		Currency: util.RandCurrency(),
	}
}

func responseMatchAccount(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
	all, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var resAccount db.Account
	err = json.Unmarshal(all, &resAccount)
	require.NoError(t, err)

	require.Equal(t, resAccount.ID, account.ID)
	require.Equal(t, resAccount.Owner, account.Owner)
	require.Equal(t, resAccount.Balance, account.Balance)
	require.Equal(t, resAccount.Currency, account.Currency)
}

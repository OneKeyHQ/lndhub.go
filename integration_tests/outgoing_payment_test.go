package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/getAlby/lndhub.go/common"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/stretchr/testify/assert"
)

func (suite *PaymentTestSuite) TestOutGoingPayment() {
	aliceFundingSats := 1000
	externalSatRequested := 500
	// 1 sat + 1 ppm
	suite.mlnd.fee = 1
	//fund alice account
	invoiceResponse := suite.createAddInvoiceReq(aliceFundingSats, "integration test external payment alice", suite.aliceToken)
	err := suite.mlnd.mockPaidInvoice(invoiceResponse, 0, false, nil)
	assert.NoError(suite.T(), err)

	//wait a bit for the callback event to hit
	time.Sleep(10 * time.Millisecond)

	//create external invoice
	externalInvoice := lnrpc.Invoice{
		Memo:  "integration tests: external pay from alice",
		Value: int64(externalSatRequested),
	}
	invoice, err := suite.externalLND.AddInvoice(context.Background(), &externalInvoice)
	assert.NoError(suite.T(), err)
	//pay external from alice
	payResponse := suite.createPayInvoiceReq(&ExpectedPayInvoiceRequestBody{
		Invoice: invoice.PaymentRequest,
	}, suite.aliceToken)
	assert.NotEmpty(suite.T(), payResponse.PaymentPreimage)

	// check that balance was reduced
	userId := getUserIdFromToken(suite.aliceToken)
	aliceBalance, err := suite.service.CurrentUserBalance(context.Background(), userId)
	if err != nil {
		fmt.Printf("Error when getting balance %v\n", err.Error())
	}
	assert.Equal(suite.T(), int64(aliceFundingSats)-int64(externalSatRequested+int(suite.mlnd.fee)), aliceBalance)

	// check that no additional transaction entry was created
	transactonEntries, err := suite.service.TransactionEntriesFor(context.Background(), userId)
	if err != nil {
		fmt.Printf("Error when getting transaction entries %v\n", err.Error())
	}
	// verify transaction entries data
	feeAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeFees, userId)
	incomingAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeIncoming, userId)
	outgoingAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeOutgoing, userId)
	currentAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeCurrent, userId)

	outgoingInvoices, _ := suite.service.InvoicesFor(context.Background(), userId, common.InvoiceTypeOutgoing)
	incomingInvoices, _ := suite.service.InvoicesFor(context.Background(), userId, common.InvoiceTypeIncoming)
	assert.Equal(suite.T(), 1, len(outgoingInvoices))
	assert.Equal(suite.T(), 1, len(incomingInvoices))

	assert.Equal(suite.T(), 3, len(transactonEntries))

	assert.Equal(suite.T(), int64(aliceFundingSats), transactonEntries[0].Amount)
	assert.Equal(suite.T(), currentAccount.ID, transactonEntries[0].CreditAccountID)
	assert.Equal(suite.T(), incomingAccount.ID, transactonEntries[0].DebitAccountID)
	assert.Equal(suite.T(), int64(0), transactonEntries[0].ParentID)
	assert.Equal(suite.T(), incomingInvoices[0].ID, transactonEntries[0].InvoiceID)

	assert.Equal(suite.T(), int64(externalSatRequested), transactonEntries[1].Amount)
	assert.Equal(suite.T(), outgoingAccount.ID, transactonEntries[1].CreditAccountID)
	assert.Equal(suite.T(), currentAccount.ID, transactonEntries[1].DebitAccountID)
	assert.Equal(suite.T(), int64(0), transactonEntries[1].ParentID)
	assert.Equal(suite.T(), outgoingInvoices[0].ID, transactonEntries[1].InvoiceID)

	assert.Equal(suite.T(), int64(suite.mlnd.fee), transactonEntries[2].Amount)
	assert.Equal(suite.T(), feeAccount.ID, transactonEntries[2].CreditAccountID)
	assert.Equal(suite.T(), currentAccount.ID, transactonEntries[2].DebitAccountID)
	assert.Equal(suite.T(), outgoingInvoices[0].ID, transactonEntries[2].InvoiceID)

	// make sure fee entry parent id is previous entry
	assert.Equal(suite.T(), transactonEntries[1].ID, transactonEntries[2].ParentID)

	//fetch transactions, make sure the fee is there
	// check invoices again
	req := httptest.NewRequest(http.MethodGet, "/gettxs", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", suite.aliceToken))
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)
	responseBody := &[]ExpectedOutgoingInvoice{}
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.NoError(suite.T(), json.NewDecoder(rec.Body).Decode(&responseBody))
	assert.Equal(suite.T(), int64(suite.mlnd.fee), (*responseBody)[0].Fee)
}

func (suite *PaymentTestSuite) TestOutGoingPaymentWithNegativeBalance() {
	// this will cause balance to go to -1
	aliceFundingSats := 1000
	externalSatRequested := 1000
	// 1 sat + 1 ppm
	suite.mlnd.fee = 1
	//fund alice account
	invoiceResponse := suite.createAddInvoiceReq(aliceFundingSats, "integration test external payment alice", suite.aliceToken)
	err := suite.mlnd.mockPaidInvoice(invoiceResponse, 0, false, nil)
	assert.NoError(suite.T(), err)

	//wait a bit for the callback event to hit
	time.Sleep(10 * time.Millisecond)

	externalInvoice := lnrpc.Invoice{
		Memo:  "integration tests: external pay from alice",
		Value: int64(externalSatRequested),
	}
	invoice, err := suite.externalLND.AddInvoice(context.Background(), &externalInvoice)
	assert.NoError(suite.T(), err)
	//pay external from alice
	payResponse := suite.createPayInvoiceReq(&ExpectedPayInvoiceRequestBody{
		Invoice: invoice.PaymentRequest,
	}, suite.aliceToken)
	assert.NotEmpty(suite.T(), payResponse.PaymentPreimage)

	// check that balance was reduced
	userId := getUserIdFromToken(suite.aliceToken)

	aliceBalance, err := suite.service.CurrentUserBalance(context.Background(), userId)
	if err != nil {
		fmt.Printf("Error when getting balance %v\n", err.Error())
	}
	assert.Equal(suite.T(), int64(aliceFundingSats)-(int64(externalSatRequested)+suite.mlnd.fee), aliceBalance)
	assert.Equal(suite.T(), int64(-1), aliceBalance)

	// check that no additional transaction entry was created
	transactonEntries, err := suite.service.TransactionEntriesFor(context.Background(), userId)
	if err != nil {
		fmt.Printf("Error when getting transaction entries %v\n", err.Error())
	}
	// verify transaction entries data
	feeAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeFees, userId)
	incomingAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeIncoming, userId)
	outgoingAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeOutgoing, userId)
	currentAccount, _ := suite.service.AccountFor(context.Background(), common.AccountTypeCurrent, userId)

	outgoingInvoices, _ := suite.service.InvoicesFor(context.Background(), userId, common.InvoiceTypeOutgoing)
	incomingInvoices, _ := suite.service.InvoicesFor(context.Background(), userId, common.InvoiceTypeIncoming)
	assert.Equal(suite.T(), 1, len(outgoingInvoices))
	assert.Equal(suite.T(), 1, len(incomingInvoices))

	assert.Equal(suite.T(), 3, len(transactonEntries))

	assert.Equal(suite.T(), int64(aliceFundingSats), transactonEntries[0].Amount)
	assert.Equal(suite.T(), currentAccount.ID, transactonEntries[0].CreditAccountID)
	assert.Equal(suite.T(), incomingAccount.ID, transactonEntries[0].DebitAccountID)
	assert.Equal(suite.T(), int64(0), transactonEntries[0].ParentID)
	assert.Equal(suite.T(), incomingInvoices[0].ID, transactonEntries[0].InvoiceID)

	assert.Equal(suite.T(), int64(externalSatRequested), transactonEntries[1].Amount)
	assert.Equal(suite.T(), outgoingAccount.ID, transactonEntries[1].CreditAccountID)
	assert.Equal(suite.T(), currentAccount.ID, transactonEntries[1].DebitAccountID)
	assert.Equal(suite.T(), int64(0), transactonEntries[1].ParentID)
	assert.Equal(suite.T(), outgoingInvoices[0].ID, transactonEntries[1].InvoiceID)

	assert.Equal(suite.T(), int64(suite.mlnd.fee), transactonEntries[2].Amount)
	assert.Equal(suite.T(), feeAccount.ID, transactonEntries[2].CreditAccountID)
	assert.Equal(suite.T(), currentAccount.ID, transactonEntries[2].DebitAccountID)
	assert.Equal(suite.T(), outgoingInvoices[0].ID, transactonEntries[2].InvoiceID)

	// make sure fee entry parent id is previous entry
	assert.Equal(suite.T(), transactonEntries[1].ID, transactonEntries[2].ParentID)
}

func (suite *PaymentTestSuite) TestZeroAmountInvoice() {
	aliceFundingSats := 1000
	amtToPay := 1000
	//fund alice account
	invoiceResponse := suite.createAddInvoiceReq(aliceFundingSats, "integration test zero amount payment alice", suite.aliceToken)
	err := suite.mlnd.mockPaidInvoice(invoiceResponse, 0, false, nil)
	assert.NoError(suite.T(), err)

	//wait a bit for the callback event to hit
	time.Sleep(10 * time.Millisecond)

	//create external invoice
	externalInvoice := lnrpc.Invoice{
		Memo:  "integration tests: zero amount pay from alice",
		Value: 0,
	}
	invoice, err := suite.externalLND.AddInvoice(context.Background(), &externalInvoice)
	assert.NoError(suite.T(), err)
	//pay external from alice
	payResponse := suite.createPayInvoiceReq(&ExpectedPayInvoiceRequestBody{
		Invoice: invoice.PaymentRequest,
		Amount:  amtToPay,
	}, suite.aliceToken)
	assert.NotEmpty(suite.T(), payResponse.PaymentPreimage)
	assert.Equal(suite.T(), int64(amtToPay), payResponse.Amount)
}
func (suite *PaymentTestSuite) TestSimultaneousPayment() {
	aliceFundingSats := 1000
	//fund alice account
	invoiceResponse := suite.createAddInvoiceReq(aliceFundingSats, "integration test simultaneous payment", suite.aliceToken)
	err := suite.mlnd.mockPaidInvoice(invoiceResponse, 0, false, nil)
	assert.NoError(suite.T(), err)

	//wait a bit for the callback event to hit
	time.Sleep(10 * time.Millisecond)

	//create 2 external invoices
	externalInvoice1 := lnrpc.Invoice{
		Memo:  "integration tests: simultaneous payment 1 from alice",
		Value: 1,
	}
	invoice1, err := suite.externalLND.AddInvoice(context.Background(), &externalInvoice1)
	assert.NoError(suite.T(), err)

	externalInvoice2 := lnrpc.Invoice{
		Memo:  "integration tests: simultaneous payment 2 from alice",
		Value: 2,
	}
	invoice2, err := suite.externalLND.AddInvoice(context.Background(), &externalInvoice2)
	assert.NoError(suite.T(), err)

	//pay 1st in goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		resp := suite.createPayInvoiceReq(&ExpectedPayInvoiceRequestBody{
			Invoice: invoice1.PaymentRequest,
		}, suite.aliceToken)
		fmt.Println(resp.Amount)
		wg.Done()
	}()
	//pay the 2nd one at the same time
	resp := suite.createPayInvoiceReq(&ExpectedPayInvoiceRequestBody{
		Invoice: invoice2.PaymentRequest,
	}, suite.aliceToken)
	fmt.Println(resp.Amount)
	wg.Wait()

}

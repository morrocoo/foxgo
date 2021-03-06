package wallet

import (
	"github.com/fox-one/foxgo/foxerr"
	"github.com/fox-one/foxgo/request"
	"github.com/satori/go.uuid"
)

type PaymentRequest struct {
	TraceId string `json:"traceId"`
	Amount  string `json:"amount"`
	AssetId string `json:"assetId"`
	Memo    string `json:"memo"`
}

func (pr *PaymentRequest) ensureTraceId() {
	if len(pr.TraceId) == 0 {
		pr.TraceId = uuid.Must(uuid.NewV4()).String()
	}
}

func withPayment(p request.Param, pr PaymentRequest) {
	pr.ensureTraceId()
	p.SetValue("assetId", pr.AssetId)
	p.SetValue("amount", pr.Amount)
	p.SetValue("traceId", pr.TraceId)
	p.SetValue("memo", pr.Memo)
}

type TransferRequest struct {
	PaymentRequest
	CounterUserId string `json:"counterUserId"`
}

func withTransfer(p request.Param, tr TransferRequest) {
	withPayment(p, tr.PaymentRequest)
	p.SetValue("counterUserId", tr.CounterUserId)
}

type WithdrawRequest struct {
	PaymentRequest
	PublicKey string `json:"publicKey"` // 提现地址
}

func withWithdraw(p request.Param, wr WithdrawRequest) {
	withPayment(p, wr.PaymentRequest)
	p.SetValue("publicKey", wr.PublicKey)

	if wr.AssetId == EOS {
		p.SetValue("label", wr.Memo)
	}
}

func InsufficientBalance(err error) bool {
	return foxerr.MatchCode(err, 1604)
}

func AmountTooSmall(err error) bool {
	return foxerr.MatchCode(err, 1603)
}

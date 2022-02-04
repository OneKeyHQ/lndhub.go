package lnd

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

type LightningClientWrapper interface {
	ListChannels(ctx context.Context, req *lnrpc.ListChannelsRequest, options ...grpc.CallOption) (*lnrpc.ListChannelsResponse, error)
	SendPaymentSync(ctx context.Context, req *lnrpc.SendRequest, options ...grpc.CallOption) (*lnrpc.SendResponse, error)
	AddInvoice(ctx context.Context, req *lnrpc.Invoice, options ...grpc.CallOption) (*lnrpc.AddInvoiceResponse, error)
	SubscribeInvoices(ctx context.Context, req *lnrpc.InvoiceSubscription, options ...grpc.CallOption) (SubscribeInvoicesWrapper, error)
	GetInfo(ctx context.Context, req *lnrpc.GetInfoRequest, options ...grpc.CallOption) (*lnrpc.GetInfoResponse, error)
	DecodeOffer(ctx context.Context, offer string) (*Offer, error)
}

type SubscribeInvoicesWrapper interface {
	Recv() (*lnrpc.Invoice, error)
}
type Offer struct {
	Type        string   `json:"type"`
	OfferID     string   `json:"offer_id"`
	Chains      []string `json:"chains"`
	Description string   `json:"description"`
	NodeID      string   `json:"node_id"`
	Signature   string   `json:"signature"`
	Vendor      string   `json:"vendor"`
	Valid       bool     `json:"valid"`
}

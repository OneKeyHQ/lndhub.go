package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Invoice : Invoice Model
type Invoice struct {
	ID              uint         `json:"id" bun:",pk,autoincrement"`
	Type            string       `json:"type"`
	UserID          int64        `json:"user_id"`
	User            *User        `bun:"rel:belongs-to,join:user_id=id"`
	Amount          int64        `json:"amount"`
	Memo            string       `json:"memo"`
	DescriptionHash string       `json:"description_hash"`
	PaymentRequest  string       `json:"payment_request"`
	RHash           string       `json:"r_hash"`
	State           string       `json:"state"`
	AddIndex        uint64       `json:"add_index"`
	CreatedAt       time.Time    `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt       bun.NullTime `json:"updated_at"`
	SettledAt       bun.NullTime `json:"settled_at"`
}

func (i *Invoice) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		i.UpdatedAt = bun.NullTime{Time: time.Now()}
	}
	return nil
}

var _ bun.BeforeAppendModelHook = (*Invoice)(nil)

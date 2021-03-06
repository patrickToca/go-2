package txnbuild

import (
	"github.com/stellar/go/amount"
	"github.com/stellar/go/price"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

//NewCreateOfferOp returns a ManageOffer operation to create a new offer, by
// setting the OfferID to "0".
func NewCreateOfferOp(selling, buying *Asset, amount, price string) ManageOffer {
	return ManageOffer{
		Selling: selling,
		Buying:  buying,
		Amount:  amount,
		Price:   price,
		OfferID: 0,
	}
}

//NewUpdateOfferOp returns a ManageOffer operation to update an offer.
func NewUpdateOfferOp(selling, buying *Asset, amount, price string, offerID uint64) ManageOffer {
	return ManageOffer{
		Selling: selling,
		Buying:  buying,
		Amount:  amount,
		Price:   price,
		OfferID: offerID,
	}
}

//NewDeleteOfferOp returns a ManageOffer operation to delete an offer, by
// setting the Amount to "0".
func NewDeleteOfferOp(offerID uint64) ManageOffer {
	// It turns out Stellar core doesn't care about any of these fields except the amount.
	// However, Horizon will reject ManageOffer if it is missing fields.
	// Horizon will also reject if Buying == Selling.
	// Therefore unfortunately we have to make up some dummy values here.
	return ManageOffer{
		Selling: NewNativeAsset(),
		Buying:  NewAsset("FAKE", "GBAQPADEYSKYMYXTMASBUIS5JI3LMOAWSTM2CHGDBJ3QDDPNCSO3DVAA"),
		Amount:  "0",
		Price:   "1",
		OfferID: offerID,
	}
}

// ManageOffer represents the Stellar manage offer operation. See
// https://www.stellar.org/developers/guides/concepts/list-of-operations.html
type ManageOffer struct {
	Selling *Asset
	Buying  *Asset
	Amount  string
	Price   string // TODO: Extend to include number, and n/d fraction. See package 'amount'
	OfferID uint64
	xdrOp   xdr.ManageOfferOp
}

// BuildXDR for ManageOffer returns a fully configured XDR Operation.
func (mo *ManageOffer) BuildXDR() (xdr.Operation, error) {
	xdrSelling, err := mo.Selling.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "Failed to set XDR 'Selling' field")
	}
	mo.xdrOp.Selling = xdrSelling

	xdrBuying, err := mo.Buying.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "Failed to set XDR 'Buying' field")
	}
	mo.xdrOp.Buying = xdrBuying

	xdrAmount, err := amount.Parse(mo.Amount)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "Failed to parse 'Amount'")
	}
	mo.xdrOp.Amount = xdrAmount

	xdrPrice, err := price.Parse(mo.Price)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "Failed to parse 'Price'")
	}
	mo.xdrOp.Price = xdrPrice

	mo.xdrOp.OfferId = xdr.Uint64(mo.OfferID)

	opType := xdr.OperationTypeManageOffer
	body, err := xdr.NewOperationBody(opType, mo.xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "Failed to build XDR OperationBody")
	}

	return xdr.Operation{Body: body}, nil
}

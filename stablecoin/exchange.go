package stablecoin

import (
	"github.com/pkg/errors"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	oracle "github.com/YaleOpenLab/openx/oracle"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func Exchange(recipientPK string, recipientSeed string, convAmount string) error {

	if !xlm.AccountExists(recipientPK) {
		return errors.New("Account does not exist, quitting!")
	}

	// check whether user has enough xlm to pay. If not, quit
	balance, err := xlm.GetNativeBalance(recipientPK)
	if err != nil {
		return errors.Wrap(err, "couldn't get native balance from api")
	}

	if utils.StoF(balance) <= utils.StoF(convAmount) {
		return errors.New("insufficient balance")
	}

	// TODO: add check for trust limit here
	hash, err := assets.TrustAsset(consts.Code, consts.StableCoinAddress, consts.StablecoinTrustLimit, recipientPK, recipientSeed)
	if err != nil {
		return errors.Wrap(err, "couldn't trust asset")
	}
	log.Println("tx hash for trusting stableUSD: ", hash)
	// now send coins across and see if our tracker detects it
	_, hash, err = xlm.SendXLM(consts.StablecoinPublicKey, convAmount, recipientSeed, "Exchange XLM for stablecoin")
	if err != nil {
		return errors.Wrap(err, "couldn't send xlm")
	}
	log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipientPK)
	return nil
}

// offer to exchange user's xlm balance for stableusd if the user does not have enough
// stableUSD to complete the payment
func OfferExchange(publicKey string, seed string, invAmount string) error {

	log.Println("OFFERING EXCHANGE TO INVESTOR")
	balance, err := xlm.GetAssetBalance(publicKey, consts.Code)
	if err != nil {
		// the user does not have a balance in STABLEUSD
		balance = "0"
	}

	balF := utils.StoF(balance)
	invF := utils.StoF(invAmount)
	if balF < invF {
		// user's stablecoin balance is less than the amount he wishes to invest, get stablecoin
		// equal to the amount he wishes to exchange
		diff := invF - balF + 1 // the extra 1 is to cover for fees
		// checking whether the user has enough xlm balance to cover for the exchange is done by Exchange()
		xlmBalance, err := xlm.GetNativeBalance(publicKey)
		if err != nil {
			return errors.Wrap(err, "couldn't get native balance from api")
		}

		totalUSD := oracle.ExchangeXLMforUSD(xlmBalance) // amount in stablecoin that the user would receive for diff

		if totalUSD < diff {
			return errors.New("User does not have enough funds to complete this transaction")
		}

		// now we need to exchange XLM equal to diff in stablecoin
		exchangeRate := oracle.ExchangeXLMforUSD("1")
		// 1 xlm can fetch exchangeRate USD, how much xlm does diff USD need?
		amountToExchange := diff / exchangeRate

		err = Exchange(publicKey, seed, utils.FtoS(amountToExchange))
		if err != nil {
			return errors.New("Unable to exchange XLM for USD and automate payment. Please get more STABLEUSD to fulfil the payment")
		}
	}

	return nil
}

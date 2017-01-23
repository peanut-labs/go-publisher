// Package peanutlabs is a utility package to help integrate with Peanut Labs as publisher.
// Publisher documentation: http://peanut-labs.github.io/publisher-doc/
package peanutlabs

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	// TransactionStatusComplete indicates a successfully completed action
	TransactionStatusComplete = "C"

	// TransactionStatusFailed indicates failure (quality check OR overquota)
	TransactionStatusFailed = "F"

	// TransactionStatusScreenout indicates disqualification from a survey
	TransactionStatusScreenout = "P"

	// OfferTypeOffer is the type for CPA offers
	OfferTypeOffer = "offer"

	// OfferTypeSurvey is the type for surveys
	OfferTypeSurvey = "survey"

	// NotificationResponseSuccess indicates response PL expects as acknowledgement
	// for successfull processing of the reward
	NotificationResponseSuccess = 1

	// NotificationResponseFailure indicates response PL expects IF reward cannot be processed.
	// PL will make a total of 5 attempts to renotify.
	NotificationResponseFailure = 0
)

var (
	defaultHost = "https://www.peanutlabs.com"

	errInvalidAppID        = errors.New("Invalid Application ID")
	errInvalidSecKey       = errors.New("Invalid Security Key")
	errInvalidTxnKey       = errors.New("Invalid Transaction Key")
	errInvalidEndUserID    = errors.New("Invalid EndUserID")
	errInvalidCallbackHash = errors.New("Invalid Hash for the reward notification")

	errInvalidAmount         = errors.New("Invalid Amount in callback")
	errInvalidCurrencyAmount = errors.New("Invalid Currency Amount in callback")
)

// Offer struct represents an offer/survey that user completes to get points
type Offer struct {
	ID    string
	Title string
	Type  string
}

// RewardNotification struct represents the parameters sent by PL via processing script call
type RewardNotification struct {
	// EndUserID indicates the User ID within your own application
	EndUserID string

	// PLUserID is the 3 part user ID used by PL to identify the user
	PLUserID string

	// Amount is the dollar amount earned by the publisher
	Amount float64

	//Status indicates if a transaction was successful or not
	Status string

	//TransactionID in the PL system, to be used to follow up on support
	TransactionID string

	//CurrencyAmount is the amount of virtual currency earned by the user
	CurrencyAmount float64

	//CurrencyName is the name of the currency in your system
	CurrencyName string

	//Program indicates which currency user earned, in case you've multiple
	Program string

	//Offer metadata
	Offer Offer
}

// Publisher struct represents
type Publisher struct {
	ApplicationID  int
	SecurityKey    string
	TransactionKey string
}

// GenerateUserID returns the 3 part UserID expected by PeanutLabs Reward Center (& other products)
func (p *Publisher) GenerateUserID(endUserID string) (string, error) {
	if !validEndUserID(endUserID) {
		return "", errInvalidEndUserID
	}
	usergo := md5Hash(fmt.Sprintf("%s%d%s", endUserID, p.ApplicationID, p.SecurityKey))
	return fmt.Sprintf("%s-%d-%s", endUserID, p.ApplicationID, usergo[0:10]), nil
}

// GenerateRewardCenterURL returns the RC URL with the userID
func (p *Publisher) GenerateRewardCenterURL(endUserID string) (string, error) {
	uid, err := p.GenerateUserID(endUserID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/userGreeting.php?userId=%s", defaultHost, uid), nil
}

// ProcessRewardNotification parses PL reward notification
func (p *Publisher) ProcessRewardNotification(req *http.Request) (*RewardNotification, error) {
	sk := p.SecurityKey
	tk := p.TransactionKey
	query := req.URL.Query()

	rn := &RewardNotification{}
	rn.TransactionID = query.Get("transactionId")
	rn.Offer.ID = query.Get("offerInvitationId")
	oidHash := query.Get("oidHash")
	txnHash := query.Get("txnHash")
	if !validOidHash(rn.Offer.ID, sk, oidHash) || !validTxnHash(rn.TransactionID, tk, txnHash) {
		return nil, errInvalidCallbackHash
	}

	amt, err := strconv.ParseFloat(query.Get("amt"), 64)
	if err != nil {
		return rn, errInvalidAmount
	}
	rn.Amount = amt

	rn.Offer.Title = query.Get("offerTitle")
	rn.Offer.Type = query.Get("offerType")

	rn.Status = query.Get("status")
	rn.Program = query.Get("program")
	rn.CurrencyName = query.Get("currencyName")
	rn.CurrencyAmount, err = strconv.ParseFloat(query.Get("currencyAmt"), 64)
	if err != nil {
		return rn, errInvalidCurrencyAmount
	}

	rn.PLUserID = query.Get("userId")
	rn.EndUserID = query.Get("endUserId")

	return rn, nil
}

// New creates an instance of PL, which implements functionality to integrate the Reward Center or DL etc.
func New(appID int, secKey string, txnKey string) (*Publisher, error) {
	if appID == 0 {
		return nil, errInvalidAppID
	}
	if len(secKey) == 0 {
		return nil, errInvalidSecKey
	}
	if len(txnKey) == 0 {
		return nil, errInvalidTxnKey
	}
	return &Publisher{ApplicationID: appID, SecurityKey: secKey}, nil
}

//helper methods
func validEndUserID(endUserID string) bool {
	if len(endUserID) == 0 || len(endUserID) > 200 {
		return false
	}
	return true
}

func validOidHash(oid string, secKey string, hash string) bool {
	return md5Hash(fmt.Sprintf("%s%s", oid, secKey)) == hash
}

func validTxnHash(tid string, txnKey string, hash string) bool {
	return md5Hash(fmt.Sprintf("%s%s", tid, txnKey)) == hash
}

func md5Hash(val string) string {
	sum := md5.Sum([]byte(val))
	return fmt.Sprintf("%x", sum)
}

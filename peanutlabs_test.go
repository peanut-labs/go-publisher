package peanutlabs_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	peanutlabs "github.com/peanut-labs/go-publisher"
)

func TestPublisherNew(t *testing.T) {
	tests := []struct {
		AppID  int
		SecKey string
		TxnKey string
		Valid  bool
	}{{
		AppID:  0,
		SecKey: "xxx",
		TxnKey: "xxx",
		Valid:  false,
	}, {
		AppID:  1,
		SecKey: "",
		TxnKey: "xxx",
		Valid:  false,
	}, {
		AppID:  1,
		SecKey: "xxx",
		TxnKey: "xxx",
		Valid:  true,
	}}

	for _, tt := range tests {
		p, err := peanutlabs.New(tt.AppID, tt.SecKey, tt.TxnKey)
		if !tt.Valid {
			if err == nil {
				t.Errorf("expected New to fail for AppID: %d and SecKey: %s", tt.AppID, tt.SecKey)
			}
		}

		if tt.Valid {
			if err != nil {
				t.Error(err)
			}

			if p.ApplicationID != tt.AppID || p.SecurityKey != tt.SecKey {
				t.FailNow()
			}
		}
	}
}

func TestPublisher_GenerateUserID(t *testing.T) {
	appID := 1
	secKey := "123"
	tests := []struct {
		EndUserID string
		PLUserID  string
		Valid     bool
	}{{
		EndUserID: "",
		PLUserID:  "",
		Valid:     false,
	}, {
		EndUserID: "saad",
		PLUserID:  "saad-1-bb753c1132",
		Valid:     true,
	}}

	p, err := peanutlabs.New(appID, secKey, secKey)
	if err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		uid, err := p.GenerateUserID(tt.EndUserID)
		if !tt.Valid {
			if err == nil {
				t.FailNow()
			}
		}

		if tt.Valid {
			if tt.PLUserID != uid {
				t.Errorf("expected %s got %s", tt.PLUserID, uid)
			}
		}
	}
}

func TestPublisher_GenerateRewardCenterURL(t *testing.T) {
	appID := 1
	secKey := "123"
	p, err := peanutlabs.New(appID, secKey, secKey)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		EndUserID   string
		ExpectedURL string
		Valid       bool
	}{{
		EndUserID:   "",
		ExpectedURL: "",
		Valid:       false,
	}, {
		EndUserID:   "saad",
		ExpectedURL: "https://www.peanutlabs.com/userGreeting.php?userId=saad-1-bb753c1132",
		Valid:       true,
	}}

	for _, tt := range tests {
		url, err := p.GenerateRewardCenterURL(tt.EndUserID)
		if !tt.Valid {
			if err == nil {
				t.FailNow()
			}
		}

		if tt.Valid {
			if tt.ExpectedURL != url {
				t.Errorf("expected %s got %s", tt.ExpectedURL, url)
			}
		}
	}
}

func TestPublisher_ProcessRewardNotification(t *testing.T) {
	query := `cmd=transactionComplete&userId=saad-1-bb753c1132&amt=1.0&offerInvitationId=123&status=C&oidHash=4297f44b13955235245b2497399d7a93&currencyAmt=50&transactionId=456&endUserId=saad&offerTitle=Survey&useragent=Peanut+Labs+Media&currencyName=Pointies&offerType=Survey&txnHash=250cf8b51c773f3f8dc8b4be867a9a02&program=`
	url := fmt.Sprintf("/callback?%s", query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	appID := 1
	secKey := "123"
	p, err := peanutlabs.New(appID, secKey, secKey)
	if err != nil {
		t.Fatal(err)
	}

	expected := &peanutlabs.RewardNotification{
		EndUserID:      "saad",
		PLUserID:       "saad-1-bb753c1132",
		Amount:         1.0,
		Status:         peanutlabs.TransactionStatusComplete,
		TransactionID:  "456",
		CurrencyAmount: 50.0,
		CurrencyName:   "Pointies",
		Program:        "",
		Offer: peanutlabs.Offer{
			ID:    "123",
			Title: "Survey",
			Type:  "Survey",
		},
	}
	rn, err := p.ProcessRewardNotification(req)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(rn, expected) {
		t.FailNow()
	}

}

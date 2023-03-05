package report

import (
	"regexp"
	"testing"
)

func TestIf(t *testing.T) {
	query := `
		"SELECT COUNT(1) FROM wal_wallet_log l 
		WHERE l.wallet_id = {wallet_id}\n AND title LIKE '%{keyword}%'
		#if { kind>0 }
			AND kind = {kind}
		#fi
		#if { kind = 0 }
			AND kind = 0 + {kind}
		#fi
		#if {trade_no}
		AND	(trade_no IS NULL OR outer_no LIKE '%{trade_no}%')";
		#fi
		#if {check} AND (check = 1) #fi

		#if {unchecked} AND (uncheck = {kind}) #fi
		`
	mp := map[string]interface{}{
		"wallet_id": 0,
		"keyword":   "提现",
		"kind":      0,
		"trade_no":  "F1",
		"check":     false,
		"unchecked": true,
	}
	ret := SqlFormat(query, mp)
	t.Log(ret)
}

func TestMathRegexp(t *testing.T) {
	var mathRegexp = regexp.MustCompile(`([^\s><=]+?)\s*([><=]*)\s*(\d+)\s*`)
	submatch := mathRegexp.FindAllStringSubmatch("key>0", 1)
	for _, v := range submatch {
		t.Log("----^", v[1], v[2], v[3])
	}
}

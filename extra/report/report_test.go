package report

import (
	"regexp"
	"testing"
)

func TestIf(t *testing.T) {
	query := `
		"SELECT COUNT(1) FROM wal_wallet_log l 
		WHERE l.wallet_id = {wallet_id}\n AND title LIKE '%{keyword}%'
		#if(kind>0)
		AND ({kind}=0 OR {kind}=kind)
		#end
		#if(trade_no)
		AND	(trade_no IS NULL OR outer_no LIKE '%{trade_no}%')";
		#end
		#if (check) AND (check = 1) #end

		#if (unchecked) AND (uncheck = {kind}) #end
		`
	mp := map[string]interface{}{
		"wallet_id": 0,
		"keyword":   "提现",
		"kind":      1,
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

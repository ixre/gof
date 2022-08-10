package report

import "testing"

func TestIf(t *testing.T) {
	query := `
		"SELECT COUNT(1) FROM wal_wallet_log l 
		WHERE l.wallet_id = {wallet_id}\n AND title LIKE '%{keyword}%'
		#if kind > 0
		AND ({kind}=0 OR {kind}=kind)
		#fi
		#if trade_no
		AND	(trade_no IS NULL OR outer_no LIKE '%{trade_no}%')";
		#fi
		`
	mp := map[string]interface{}{
		"kind":     -1,
		"trade_no": "F1",
	}
	ret := SqlFormat(query, mp)
	t.Log(ret)
}

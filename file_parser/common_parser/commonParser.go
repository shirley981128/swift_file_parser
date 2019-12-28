package common_parser

import (
	"log"
	"strings"
	"time"
)

const (
	//fields.
	TXREFNUM           = "20"
	RELATEDREF         = "21"
	ACCOUNTID          = "25"
	STATEMENTNUM       = "28C"
	DCFLOORLIMITIND    = "34F"
	DATETIMEIND        = "13D"
	OPENINGBALANCE_F   = "60F"
	OPENINGBALANCE_M   = "60M"
	STATEMENTLINES     = "61"
	INFOTOACCTOWNER    = "86"
	NUMANDSUMOFENTRY_D = "90D"
	NUMANDSUMOFENTRY_C = "90C"
	CLOSINGBALANCE_F   = "62F"
	CLOSINGBALANCE_M   = "62M"
	CLOSINGAVLBALANCE  = "64"
	FORWARDAVLBALANCE  = "65"

	//D/C Mark enum.
	CREDIT = "Credit"
	DEBIT  = "Debit"
	RC     = "Reversal of Credit"
	RD     = "Reversal of Debit"
)

type StatementNum struct {
	StatementNumber string
	SequenceNumber  string
}

type Balance struct {
	Mark     string
	Date     string
	Currency string
	Amount   string
}

type StatementLine struct {
	ValueDate        string
	EntryDate        string
	Mark             string
	FundsCode        string
	Amount           string
	TxType           string
	IDCode           string
	Ref4AcctOwner    string
	RefOfAcctServIns string
	SupDetails       string
	InfoToAcctOwner  string
}

type DCFloorLimitInd struct {
	Currency string
	Mark     string
	Amount   string
}

type NumAndSumOfEntry struct {
	Number   string
	Currency string
	Amount   string
}

func (numAndSumOfEntry *NumAndSumOfEntry) ParseField(str string) {
	for i := 0; i < 5; i++ {
		if str[i] < '0' || str[i] > '9' {
			numAndSumOfEntry.Number = str[:i]
			str = str[i:]
			break
		}
	}
	numAndSumOfEntry.Currency = str[:3]
	numAndSumOfEntry.Amount = str[3:]
	parts := strings.Split(numAndSumOfEntry.Amount, ",")
	numAndSumOfEntry.Amount = parts[0]
	if parts[1] != "" {
		numAndSumOfEntry.Amount += "." + parts[1]
	}
}

func ParseDateTimeInd(str string) string {
	t, err := time.Parse("0601021504-0700", str)
	if err != nil {
		log.Fatalln("parse time fail")
	}
	return t.UTC().Format("2006-01-02T15:04:05")
}

func (dcFloorLimitInd *DCFloorLimitInd) ParseField(str string) {
	dcFloorLimitInd.Currency = str[:3]
	str = str[3:]
	if str[0] == 'D' {
		dcFloorLimitInd.Mark = DEBIT
		str = str[1:]
	} else if str[0] == 'C' {
		dcFloorLimitInd.Mark = CREDIT
		str = str[1:]
	}
	dcFloorLimitInd.Amount = str
	parts := strings.Split(dcFloorLimitInd.Amount, ",")
	dcFloorLimitInd.Amount = parts[0]
	if parts[1] != "" {
		dcFloorLimitInd.Amount += "." + parts[1]
	}
}

func (statementNum *StatementNum) ParseField(str string) {
	numStr := strings.Split(str, "/")
	statementNum.StatementNumber = numStr[0]
	if len(numStr) > 1 {
		statementNum.SequenceNumber = numStr[1]
	}
}

func (balance *Balance) ParseField(str string) {
	switch str[0] {
	case 'C':
		balance.Mark = CREDIT
	case 'D':
		balance.Mark = DEBIT
	}
	balance.Date = str[1:7]
	t, err := time.Parse("060102", balance.Date)
	if err != nil {
		log.Fatalln("parse time fail")

	}
	balance.Date = t.UTC().Format("2006-01-02T15:04:05")
	balance.Currency = str[7:10]
	balance.Amount = str[10:]
	parts := strings.Split(balance.Amount, ",")
	balance.Amount = parts[0]
	if parts[1] != "" {
		balance.Amount += "." + parts[1]
	}
}

func (statementLine *StatementLine) ParseField(sl string, info string) bool {
	statementLine.ValueDate = sl[:6]
	t, err := time.Parse("060102", statementLine.ValueDate)
	if err != nil {
		log.Fatalln("parse time fail")
	}
	statementLine.ValueDate = t.UTC().Format("2006-01-02T15:04:05")
	sl = sl[6:]
	if sl[0] >= '0' && sl[0] <= '9' {
		statementLine.EntryDate = sl[:4]
		t, err := time.Parse("20060102", statementLine.ValueDate[:4]+statementLine.EntryDate)
		if err != nil {
			log.Fatalln("parse time fail")
		}
		statementLine.EntryDate = t.UTC().Format("2006-01-02T15:04:05")
		sl = sl[4:]
	}
	switch sl[:2] {
	case "CP":
		statementLine.Mark = CREDIT
	case "DP":
		statementLine.Mark = DEBIT
	case "RC":
		statementLine.Mark = RC
	case "RD":
		statementLine.Mark = RD
	}
	sl = sl[2:]
	if sl[0] < '0' || sl[0] > '9' {
		statementLine.FundsCode = string(sl[0])
		sl = sl[1:]
	}
	for i := 0; i < 15; i++ {
		if (sl[i+1] >= 'A' && sl[i+1] <= 'Z') || (sl[i+1] >= 'a' && sl[i+1] <= 'z') {
			statementLine.Amount += sl[:i+1]
			sl = sl[i+1:]
			break
		}
	}
	parts := strings.Split(statementLine.Amount, ",")
	statementLine.Amount = parts[0]
	if parts[1] != "" {
		statementLine.Amount += "." + parts[1]
	}
	statementLine.TxType = string(sl[0])
	statementLine.IDCode = sl[1:4]
	end := len(sl)
	if supDetailsStart := strings.Index(sl, "\n"); supDetailsStart != -1 {
		statementLine.SupDetails = sl[supDetailsStart+2:]
		end = supDetailsStart
	}
	if refOfAcctServInsStart := strings.Index(sl, "//"); refOfAcctServInsStart != -1 {
		statementLine.RefOfAcctServIns = sl[refOfAcctServInsStart+2 : end]
		end = refOfAcctServInsStart
	}
	statementLine.Ref4AcctOwner = sl[4:end]
	statementLine.InfoToAcctOwner = info
	return true
}

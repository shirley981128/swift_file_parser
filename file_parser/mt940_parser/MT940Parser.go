package mt940_parser

import (
	. "github.com/shirley981128/file_parser/common_parser"
	"log"
	"strings"
)

type MT940Origin struct {
	TxRefNum          string
	RelatedRef        string
	AccountId         string
	StatementNum      string
	OpeningBalance    string
	StatementLines    map[string]string
	InfoToAcctOwner   string
	ClosingBalance    string
	ClosingAvlBalance string
	ForwardAvlBalance []string
}

type MT940 struct {
	TxRefNum          string
	RelatedRef        string
	AccountId         string
	StatementNum      StatementNum
	OpeningBalance    Balance
	StatementLines    []StatementLine
	InfoToAcctOwner   string
	ClosingBalance    Balance
	ClosingAvlBalance Balance
	ForwardAvlBalance []Balance
}

func ParseMT940File(file []byte) *MT940Origin {
	var origin MT940Origin
	origin.StatementLines = make(map[string]string)
	fileStr := string(file)
	fields := strings.Split(fileStr, "\n")
	//文件的最后以"-"做结束符
	fields = fields[:len(fields)-1]
	for i := 0; i < len(fields); i++ {
		fields[i] = fields[i][1:]
		parts := strings.Split(fields[i], ":")
		j := i
		for ; j < len(fields)-1 && fields[j+1][0] != ':'; j++ {
			parts[1] += fields[j+1]
		}
		switch parts[0] {
		case TXREFNUM:
			origin.TxRefNum = parts[1]
		case RELATEDREF:
			origin.RelatedRef = parts[1]
		case ACCOUNTID:
			origin.AccountId = parts[1]
		case STATEMENTNUM:
			origin.StatementNum = parts[1]
		case OPENINGBALANCE_F, OPENINGBALANCE_M:
			origin.OpeningBalance = parts[1]
		case STATEMENTLINES:
			//将86(following 61)组装入对账行
			info := ""
			if strings.HasPrefix(fields[j+1], ":"+INFOTOACCTOWNER) {
				j++
				fields[j] = fields[j][1:]
				infoParts := strings.Split(fields[j], ":")
				for ; j < len(fields)-1 && fields[j+1][0] != ':'; j++ {
					infoParts[1] += fields[j+1]
				}
				info = infoParts[1]
			}
			origin.StatementLines[parts[1]] = info
		case INFOTOACCTOWNER:
			origin.InfoToAcctOwner = parts[1]
		case CLOSINGBALANCE_F, CLOSINGBALANCE_M:
			origin.ClosingBalance = parts[1]
		case CLOSINGAVLBALANCE:
			origin.ClosingAvlBalance = parts[1]
		case FORWARDAVLBALANCE:
			origin.ForwardAvlBalance = append(origin.ForwardAvlBalance, parts[1])
		}
		i = j
	}
	return &origin
}

func ParseMT940Field(file []byte) (*MT940, bool) {
	var mt940 MT940
	origin := ParseMT940File(file)
	if origin.TxRefNum == "" || origin.AccountId == "" || origin.StatementNum == "" || origin.OpeningBalance == "" || origin.ClosingBalance == "" {
		log.Fatalln("losing mandatory field,%+v", origin)
		return nil, false
	}
	mt940.TxRefNum = origin.TxRefNum
	mt940.AccountId = origin.AccountId
	mt940.StatementNum.ParseField(origin.StatementNum)
	mt940.OpeningBalance.ParseField(origin.OpeningBalance)
	mt940.ClosingBalance.ParseField(origin.ClosingBalance)
	if len(origin.StatementLines) > 0 {
		for statementLineStr, infoToAcctOwner := range origin.StatementLines {
			var statementLine StatementLine
			statementLine.ParseField(statementLineStr, infoToAcctOwner)
			mt940.StatementLines = append(mt940.StatementLines, statementLine)
		}
	}
	if origin.RelatedRef != "" {
		mt940.RelatedRef = origin.RelatedRef
	}
	if origin.InfoToAcctOwner != "" {
		mt940.InfoToAcctOwner = origin.InfoToAcctOwner
	}
	if origin.ClosingAvlBalance != "" {
		mt940.ClosingAvlBalance.ParseField(origin.ClosingAvlBalance)
	}
	if origin.ForwardAvlBalance != nil {
		for _, forwardAvlBalanceStr := range origin.ForwardAvlBalance {
			var forwardAvlBalance Balance
			forwardAvlBalance.ParseField(forwardAvlBalanceStr)
			mt940.ForwardAvlBalance = append(mt940.ForwardAvlBalance, forwardAvlBalance)
		}
	}

	return &mt940, true
}

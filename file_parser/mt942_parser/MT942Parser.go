package mt942_parser

import (
	. "github.com/shirley981128/file_parser/common_parser"
	"log"
	"strings"
)

type MT942Origin struct {
	TxRefNum           string
	RelatedRef         string
	AccountId          string
	StatementNum       string
	DCFloorLimitInd    []string
	DateTimeInd        string
	StatementLines     map[string]string
	InfoToAcctOwner    string
	NumAndSumOfEntry_C string
	NumAndSumOfEntry_D string
}

type MT942 struct {
	TxRefNum           string
	RelatedRef         string
	AccountId          string
	StatementNum       StatementNum
	DCFloorLimitInd    []DCFloorLimitInd
	DateTimeInd        string
	StatementLines     []StatementLine
	InfoToAcctOwner    string
	NumAndSumOfEntry_C NumAndSumOfEntry
	NumAndSumOfEntry_D NumAndSumOfEntry
}

func ParseMT942File(file []byte) *MT942Origin {
	var origin MT942Origin
	origin.StatementLines = make(map[string]string)
	fileStr := string(file)
	fields := strings.Split(fileStr, "\n")
	//文件的最后以"-"做结束符。
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
		case DCFLOORLIMITIND:
			origin.DCFloorLimitInd = append(origin.DCFloorLimitInd, parts[1])
		case DATETIMEIND:
			origin.DateTimeInd = parts[1]
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
		case NUMANDSUMOFENTRY_C:
			origin.NumAndSumOfEntry_C = parts[1]
		case NUMANDSUMOFENTRY_D:
			origin.NumAndSumOfEntry_D = parts[1]
		}
		i = j
	}
	return &origin
}

func ParseMT942Field(file []byte) (*MT942, bool) {
	var mt942 MT942
	origin := ParseMT942File(file)
	if origin.TxRefNum == "" || origin.AccountId == "" || origin.StatementNum == "" || origin.DCFloorLimitInd == nil || origin.DateTimeInd == "" {
		log.Fatalln("losing mandatory field,%+v", origin)
		return nil, false
	}
	mt942.TxRefNum = origin.TxRefNum
	mt942.AccountId = origin.AccountId
	mt942.StatementNum.ParseField(origin.StatementNum)
	mt942.DateTimeInd = ParseDateTimeInd(origin.DateTimeInd)
	for _, dcFloorLimitIndStr := range origin.DCFloorLimitInd {
		var dcFloorLimitInd DCFloorLimitInd
		dcFloorLimitInd.ParseField(dcFloorLimitIndStr)
		mt942.DCFloorLimitInd = append(mt942.DCFloorLimitInd, dcFloorLimitInd)
	}
	if origin.RelatedRef != "" {
		mt942.RelatedRef = origin.RelatedRef
	}
	if origin.InfoToAcctOwner != "" {
		mt942.InfoToAcctOwner = origin.InfoToAcctOwner
	}
	if origin.NumAndSumOfEntry_C != "" {
		mt942.NumAndSumOfEntry_C.ParseField(origin.NumAndSumOfEntry_C)
	}
	if origin.NumAndSumOfEntry_D != "" {
		mt942.NumAndSumOfEntry_D.ParseField(origin.NumAndSumOfEntry_D)
	}
	if len(origin.StatementLines) > 0 {
		for statementLineStr, infoToAcctOwner := range origin.StatementLines {
			var statementLine StatementLine
			if !statementLine.ParseField(statementLineStr, infoToAcctOwner) {
				return nil, false
			}
			mt942.StatementLines = append(mt942.StatementLines, statementLine)
		}
	}

	return &mt942, true
}

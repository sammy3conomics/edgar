package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type finDataType string

type finDataSearchInfo struct {
	finDataName finDataType
	finDataStr  string
}

var (
	//Types of financial data collected
	finDataSharesOutstanding finDataType = "Shares Outstanding"
	finDataRevenue           finDataType = "Revenue"
	finDataCostOfRevenue     finDataType = "Cost Of Revenue"
	finDataGrossMargin       finDataType = "Gross Margin"
	finDataOpsIncome         finDataType = "Operational Income"
	finDataOpsExpense        finDataType = "Operational Expense"
	finDataNetIncome         finDataType = "Net Income"
	finDataOpCashFlow        finDataType = "Operating Cash Flow"
	finDataCapEx             finDataType = "Capital Expenditure"
	finDataLDebt             finDataType = "Long-Term debt"
	finDataSDebt             finDataType = "Short-Term debt"
	finDataCLiab             finDataType = "Current Liabilities"
	finDataDeferred          finDataType = "Deferred revenue"
	finDataRetained          finDataType = "Retained Earnings"
	finDataUnknown           finDataType = "Unknown"

	//Keys to search for financial data in the filings
	finDataSearchKeys = []finDataSearchInfo{
		{finDataRevenue, "net revenue"},
		{finDataRevenue, "net sales"},
		{finDataRevenue, "total revenue"},
		{finDataRevenue, "total sales"},
		{finDataCostOfRevenue, "cost of sales"},
		{finDataCostOfRevenue, "cost of revenue"},
		{finDataGrossMargin, "gross margin"},
		{finDataSharesOutstanding, "shares outstanding"},
		{finDataOpsExpense, "operating expenses"},
		{finDataOpsIncome, "operating income"},
		{finDataOpsIncome, "operating (loss)"},
		{finDataOpsIncome, "operating loss"},
		{finDataNetIncome, "net income"},
		{finDataOpCashFlow, "operating activities"},
		{finDataCapEx, "plant and equipment"},
		{finDataCapEx, "capital expen"},
		{finDataSDebt, "current portion of long-term"},
		{finDataLDebt, "long term debt"},
		{finDataLDebt, "long-term debt"},
		{finDataCLiab, "total current liabilities"},
		{finDataDeferred, "deferred revenue"},
		{finDataRetained, "retained earnings"},
	}
)

func getFinDataType(key string) finDataType {
	key = strings.ToLower(key)
	for _, val := range finDataSearchKeys {
		lup := strings.ToLower(val.finDataStr)
		if strings.Contains(key, lup) {
			return val.finDataName
		}
	}
	return finDataUnknown
}

type Company struct {
	Ticker  string    `json:"Company"`
	Reports []*Filing `json:"Financial Reports"`
}

type Filing struct {
	Date    string           `json:"Report date"`
	FinData *FinancialReport `json:"Financial Data"`
}

type FinancialReport struct {
	DocType filingType  `json:"Filing Type"`
	Entity  *EntityData `json:"Entity Information"`
	Ops     *OpsData    `json:"Operational Information"`
	Bs      *BSData     `json:"Balance Sheet Information"`
	Cf      *CfData     `json:"Cash Flow Information"`
}

type EntityData struct {
	ShareCount int64 `json:"Shares Outstanding" required:"true"`
}

type OpsData struct {
	Revenue     int64 `json:"Revenue" required:"true"`
	CostOfSales int64 `json:"Cost Of Revenue" required:"true"`
	GrossMargin int64 `json:"Gross Margin" required:"true" generate:"true"`
	OpIncome    int64 `json:"Operational Income" required:"true"`
	OpExpense   int64 `json:"Operational Expense" required:"true"`
	NetIncome   int64 `json:"Net Income" required:"true"`
}

type CfData struct {
	OpCashFlow int64 `json:"Operating Cash Flow" required:"true"`
	CapEx      int64 `json:"Capital Expenditure" required:"true"`
}

type BSData struct {
	LDebt    int64 `json:"Long-Term debt" required:"true"`
	SDebt    int64 `json:"Short-Term debt" required:"true"`
	CLiab    int64 `json:"Current Liabilities" required:"true"`
	Deferred int64 `json:"Deferred revenue" required:"false"`
	Retained int64 `json:"Retained Earnings" required:"true"`
}

/*
func (c *Company) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Company data")
	}
	fmt.Println("COMPANY")
	return string(data)
}
*/

func (c Company) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Company data")
	}
	return string(data)
}

func (f Filing) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Filing data")
	}
	fmt.Println("FILING")
	return string(data)
}

func (f FinancialReport) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling financial data")
	}
	return string(data)
}

func (bs BSData) String() string {
	data, err := json.MarshalIndent(bs, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling balance sheet data")
	}
	return string(data)
}

func (cf CfData) String() string {
	data, err := json.MarshalIndent(cf, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling cash flow data")
	}
	return string(data)
}

func (ops OpsData) String() string {
	data, err := json.MarshalIndent(ops, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Operational information data")
	}
	return string(data)
}

func generateData(data interface{}, name string) int64 {
	switch name {
	case "GrossMargin":
		val, ok := data.(*OpsData)
		if ok {
			return val.Revenue - val.CostOfSales
		}
	}
	return 0
}

//Validate is a function to check that no field is set to 0 after parsing
func Validate(data interface{}) error {
	var err string
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tag, ok := t.Field(i).Tag.Lookup("required")
		val := v.Field(i).Int()
		if val == 0 && (ok && tag == "true") {
			tag, ok = t.Field(i).Tag.Lookup("generate")
			if ok && tag == "true" {
				num := generateData(data, t.Field(i).Name)
				if num == 0 {
					err += t.Field(i).Name + ","
				} else {
					v.Field(i).SetInt(num)
				}
			} else {
				err += t.Field(i).Name + ","
			}
		}
	}
	if len(err) > 0 {
		return errors.New("[" + err + "] " + "attributes did not parse")
	}
	return nil
}

func SetData(data interface{}, finType finDataType, val string) error {

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tag, ok := t.Field(i).Tag.Lookup("json")
		if ok && string(finType) == tag {
			if v.Field(i).Int() == 0 {
				v.Field(i).SetInt(normalizeNumber(val))
			}
			return nil
		}
	}
	return errors.New("Could not find the field to set: " + string(finType))
}

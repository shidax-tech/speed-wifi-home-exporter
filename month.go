package main

import (
	"encoding/xml"
	"fmt"
	"time"
)

type Date time.Time

func (date *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	d.DecodeElement(&s, &start)
	if parsed, err := time.Parse("2006-01-02", s); err != nil {
		return err
	} else {
		*date = Date(parsed)
		return nil
	}
}

func (d *Date) String() string {
	return (*time.Time)(d).Format("2006-01-02")
}

type MonthStatistics struct {
	CurrentMonthDownload int
	CurrentMonthUpload   int
	MonthDuration        int
	MonthLastClearTime   *Date
	CurrentDayUsed       int
	CurrentDayDuration   int
	CurrentMonthHSA      int `xml:"current_month_hsa"`
}

type MonthResponse struct {
	Response MonthStatistics `xml:"response",ommit`
}

type MonthClient struct {
	lastStat        MonthStatistics
	totalUploaded   int
	totalDownloaded int

	APIEndpoint string
}

func NewMonthClient(address string) MonthClient {
	return MonthClient{APIEndpoint: fmt.Sprintf("http://%s/api/monitoring/month_statistics", address)}
}

func (mc MonthClient) Fetch() (*MonthStatistics, error) {
	var month MonthStatistics
	if err := HTTPGetXML(mc.APIEndpoint, &month); err != nil {
		return nil, err
	}
	return &month, nil
}

func (mc MonthClient) Collect() (uploaded, downloaded int, err error) {
	m, err := mc.Fetch()
	if err != nil {
		return 0, 0, err
	}

	if mc.lastStat.MonthLastClearTime != m.MonthLastClearTime {
		mc.totalUploaded += m.CurrentMonthUpload
		mc.totalDownloaded += m.CurrentMonthDownload
	} else {
		mc.totalUploaded += m.CurrentMonthUpload - mc.lastStat.CurrentMonthUpload
		mc.totalDownloaded += m.CurrentMonthDownload - mc.lastStat.CurrentMonthDownload
	}
	mc.lastStat = *m

	return mc.totalUploaded, mc.totalDownloaded, nil
}

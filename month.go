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
	if parsed, err := time.Parse("2006-1-2", s); err != nil {
		return err
	} else {
		*date = Date(parsed)
		return nil
	}
}

func (d *Date) String() string {
	return (*time.Time)(d).Format("2006-1-2")
}

type MonthStatistics struct {
	CurrentMonthDownload int64
	CurrentMonthUpload   int64
	MonthDuration        int64
	MonthLastClearTime   *Date
	CurrentDayUsed       int64
	CurrentDayDuration   int64
	CurrentMonthHSA      int64 `xml:"current_month_hsa"`
}

type MonthResponse struct {
	Response MonthStatistics `xml:"response",ommit`
}

type MonthClient struct {
	lastStat        MonthStatistics
	totalUploaded   int64
	totalDownloaded int64

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

type Statistics struct {
	TotalUploaded     int64
	TotalDownloaded   int64
	MonthlyUploaded   int64
	MonthlyDownloaded int64
}

func (mc MonthClient) Collect() (Statistics, error) {
	m, err := mc.Fetch()
	if err != nil {
		return Statistics{}, err
	}

	if mc.lastStat.MonthLastClearTime != m.MonthLastClearTime {
		mc.totalUploaded += m.CurrentMonthUpload
		mc.totalDownloaded += m.CurrentMonthDownload
	} else {
		mc.totalUploaded += m.CurrentMonthUpload - mc.lastStat.CurrentMonthUpload
		mc.totalDownloaded += m.CurrentMonthDownload - mc.lastStat.CurrentMonthDownload
	}
	mc.lastStat = *m

	return Statistics{
		mc.totalUploaded,
		mc.totalDownloaded,
		m.CurrentMonthUpload,
		m.CurrentMonthDownload,
	}, nil
}

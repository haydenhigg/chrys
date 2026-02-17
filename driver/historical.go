package driver

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/haydenhigg/chrys/frame"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type HistoricalDriver struct {
	DataRoot string
	NameFmt  string // the fmt string for the CSV files with the frames
}

func NewHistorical(dataRoot, nameFmt string) *HistoricalDriver {
	return &HistoricalDriver{
		DataRoot: dataRoot,
		NameFmt:  nameFmt,
	}
}

func (d *HistoricalDriver) FetchFramesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*frame.Frame, error) {
	// split pair into assets
	assets := strings.SplitN(pair, "/", 2)
	base, quote := assets[0], assets[1]

	// format data file path
	dataFile := fmt.Sprintf(d.NameFmt, base, quote, int(interval.Minutes()))
	dataFilePath := filepath.Join(d.DataRoot, dataFile)

	// read data file
	file, err := os.Open(dataFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	frames := []*frame.Frame{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return frames, err
		}

		epochTime, _ := strconv.ParseInt(record[0], 10, 64)
		t := time.Unix(epochTime, 0)

		if !t.Before(since) {
			open, _ := strconv.ParseFloat(record[1], 64)
			high, _ := strconv.ParseFloat(record[2], 64)
			low, _ := strconv.ParseFloat(record[3], 64)
			close, _ := strconv.ParseFloat(record[4], 64)
			volume, _ := strconv.ParseFloat(record[5], 64)

			frames = append(frames, &frame.Frame{
				Time:   t,
				Open:   open,
				High:   high,
				Low:    low,
				Close:  close,
				Volume: volume,
			})
		}
	}

	slices.SortFunc(frames, func(a, b *frame.Frame) int {
		return a.Time.Compare(b.Time)
	})

	if since.Add(interval).Before(frames[0].Time) {
		return frames, errors.New("insufficient historical frames")
	}

	return frames, nil
}

func (d *HistoricalDriver) FetchBalances() (map[string]float64, error) {
	return map[string]float64{}, nil
}

func (d *HistoricalDriver) MarketOrder(side, pair string, quantity float64) error {
	return nil
}

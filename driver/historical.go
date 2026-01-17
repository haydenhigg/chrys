package driver

import (
	"encoding/csv"
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

type Historical struct {
	DataRoot string
	NameFmt  string // the fmt string for the CSV files with the frames
}

func NewHistorical(dataRoot, nameFmt string) *Historical {
	return &Historical{
		DataRoot: dataRoot,
		NameFmt:  nameFmt,
	}
}

func (c *Historical) FetchFramesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*frame.Frame, error) {
	// split pair into assets
	assets := strings.SplitN(pair, "/", 2)
	base, quote := assets[0], assets[1]

	// format data file path
	dataFile := fmt.Sprintf(c.NameFmt, base, quote, int(interval.Minutes()))
	dataFilePath := filepath.Join(c.DataRoot, dataFile)

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
		time := time.Unix(epochTime, 0)

		if !time.Before(since) {
			open, _ := strconv.ParseFloat(record[1], 64)
			high, _ := strconv.ParseFloat(record[2], 64)
			low, _ := strconv.ParseFloat(record[3], 64)
			close, _ := strconv.ParseFloat(record[4], 64)
			volume, _ := strconv.ParseFloat(record[5], 64)

			frames = append(frames, &frame.Frame{
				Time:   time,
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

	return frames, nil
}

func (c *Historical) FetchBalances() (map[string]float64, error) {
	return map[string]float64{}, nil
}

func (c *Historical) MarketOrder(side, pair string, quantity float64) error {
	return nil
}

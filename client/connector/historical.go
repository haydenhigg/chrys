package connector

import (
	"encoding/csv"
	"fmt"
	"github.com/haydenhigg/chrys/candle"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const HISTORICAL_FILE_NAME_FORMAT = "%s%s_%d.csv"

type Historical struct {
	DataRoot string
}

func (c *Historical) FetchCandlesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*candle.Candle, error) {
	symbols := strings.SplitN(pair, "/", 2)
	filePath := filepath.Join(c.DataRoot, fmt.Sprintf(
		HISTORICAL_FILE_NAME_FORMAT,
		symbols[0],
		symbols[1],
		int(interval.Minutes()),
	))

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	candles := []*candle.Candle{}

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return candles, err
		}

		timeEpoch, _ := strconv.ParseInt(record[0], 10, 64)
		time := time.Unix(timeEpoch, 0)

		if !time.Before(since) {
			open, _ := strconv.ParseFloat(record[1], 64)
			high, _ := strconv.ParseFloat(record[2], 64)
			low, _ := strconv.ParseFloat(record[3], 64)
			close, _ := strconv.ParseFloat(record[4], 64)
			volume, _ := strconv.ParseFloat(record[5], 64)

			candles = append(candles, &candle.Candle{
				Time:   time,
				Open:   open,
				High:   high,
				Low:    low,
				Close:  close,
				Volume: volume,
			})
		}
	}

	slices.SortFunc(candles, func(a, b *candle.Candle) int {
		return a.Time.Compare(b.Time)
	})

	return candles, nil
}

func (c *Historical) FetchBalances() (map[string]float64, error) {
	return map[string]float64{}, nil
}

func (c *Historical) PlaceMarketOrder(side, pair string, quantity float64) error {
	return nil
}

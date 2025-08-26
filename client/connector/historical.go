package connector

import (
	"encoding/csv"
	"fmt"
	"github.com/haydenhigg/chrys"
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

func (c *Historical) FetchFramesSince(
	series *chrys.Series,
	since time.Time,
) ([]*chrys.Frame, error) {
	filePath := filepath.Join(c.DataRoot, fmt.Sprintf(
		HISTORICAL_FILE_NAME_FORMAT,
		series.Pair.Base,
		series.Pair.Quote,
		int(series.Interval.Minutes()),
	))

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	frames := []*chrys.Frame{}

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return frames, err
		}

		timeEpoch, _ := strconv.ParseInt(record[0], 10, 64)
		time := time.Unix(timeEpoch, 0)

		if !time.Before(since) {
			open, _ := strconv.ParseFloat(record[1], 64)
			high, _ := strconv.ParseFloat(record[2], 64)
			low, _ := strconv.ParseFloat(record[3], 64)
			close, _ := strconv.ParseFloat(record[4], 64)
			volume, _ := strconv.ParseFloat(record[5], 64)

			frames = append(frames, &chrys.Frame{
				Time:   time,
				Open:   open,
				High:   high,
				Low:    low,
				Close:  close,
				Volume: volume,
			})
		}
	}

	slices.SortFunc(frames, func(a, b *chrys.Frame) int {
		return a.Time.Compare(b.Time)
	})

	return frames, nil
}

func (c *Historical) FetchBalances() (map[string]float64, error) {
	return map[string]float64{}, nil
}

func (c *Historical) PlaceMarketOrder(side, pair string, quantity float64) error {
	return nil
}

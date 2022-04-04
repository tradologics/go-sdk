package helpers

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func ParseBars(bars *map[string]interface{}) map[string]map[string][]float64 {

	type payload struct {
		timestamp float64
		data      map[string]float64
	}

	data := make(map[string]map[string][]float64)
	payloadData := make(map[string][]payload)

	for barsKey, barsValue := range *bars {
		timestamp, _ := time.Parse("2006-01-02T15:04:05.999999999", barsKey)
		assets := barsValue.(map[string]interface{})

		for assetKey, assetValue := range assets {

			payloadMap := make(map[string]float64)
			for k, v := range assetValue.(map[string]interface{}) {
				flt, _ := strconv.ParseFloat(fmt.Sprint(v), 64)
				payloadMap[k] = flt
			}

			payloadStruct := payload{timestamp: float64(timestamp.UnixMilli()), data: payloadMap}

			if _, ok := payloadData[assetKey]; ok {
				payloadData[assetKey] = append(payloadData[assetKey], payloadStruct)
			} else {
				payloadData[assetKey] = []payload{payloadStruct}
			}
		}
	}

	for asset, payload := range payloadData {
		sort.SliceStable(payload, func(i, j int) bool {
			return payload[i].timestamp < payload[j].timestamp
		})
		for _, s := range payload {
			//	fmt.Println(i, s)

			if _, ok := data[asset]["dt"]; ok {
				data[asset]["dt"] = append(data[asset]["dt"], s.timestamp)
			} else {
				dt := map[string][]float64{
					"dt": []float64{s.timestamp},
				}
				data[asset] = dt
			}

			for k, v := range s.data {
				flt, _ := strconv.ParseFloat(fmt.Sprint(v), 64)
				data[asset][k] = append(data[asset][k], flt)
			}
		}
	}

	return data
}

func Median(array []float64) float64 {
	newArray := make([]float64, len(array))
	copy(newArray, array)
	sort.Float64s(newArray)
	var middle float64
	if len(newArray)%2 == 0 {
		middle = (newArray[len(newArray)/2] + newArray[len(newArray)/2-1]) / 2
	} else {
		middle = newArray[(len(newArray)-1)/2]
	}
	return middle
}

func Mean(array *[]float64) float64 {
	return sumOfArrayElements(*array) / float64(len(*array))
}

func sumOfArrayElements(array []float64) float64 {
	sum := 0.0
	for _, v := range array {
		sum += v
	}
	return sum
}

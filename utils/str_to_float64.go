package utils

import (
	"math"
	"strconv"
)

// StringUtils
// @example
/**
var strutil StringUtils

 floatValues := []float64{123.456, -789.1011, 0.0, math.Pi, math.E}

    for _, floatValue := range floatValues {
        hexStr := strutil.Float64ToHexString(floatValue)
        fmt.Printf("Float64 value %f converts to hex string: %s\n", floatValue, hexStr)
    }
**/
type StringUtils struct{}

func hexStringToFloat64(hexStr string) (float64, error) {
	// Step 1: Parse the hexadecimal string to a uint64
	intValue, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}

	// Step 2: Convert the uint64 to a float64
	floatValue := float64(intValue)
	return floatValue, nil
}

func (s StringUtils) HexStrToFloat64(hexStr string) (float64, error) {
	return hexStringToFloat64(hexStr)
}

func float64ToHexString(floatValue float64) string {
	bits := math.Float64bits(floatValue)
	hexStr := strconv.FormatUint(bits, 16)
	return hexStr
}

func (s StringUtils) Float64ToHexString(floatValue float64) string {
	return float64ToHexString(floatValue)
}

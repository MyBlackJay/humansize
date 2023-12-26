package humansize

import (
	"errors"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	measurePattern string = `^([bB]|[bB]ytes|[kmgtpeKMGTPE]|[kmgtpeKMGTPE]?[iI]|[kmgtpeKMGTPE][iI]?[bB])?$`
	sizePattern    string = `^([0-9]+|[0-9]*\.[0-9]+)([bB]|[bB]ytes|[kmgtpeKMGTPE]|[kmgtpeKMGTPE]?[iI]|[kmgtpeKMGTPE][iI]?[bB])?$`
)

var defaultMeasure = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// ReadableSize is the representation of a compiled data size expression.
type ReadableSize struct {
	input    string
	measure  int64
	compiled *big.Int
}

// GetInput returns original data size expression.
func (rs *ReadableSize) GetInput() string {
	return rs.input
}

// GetMeasure returns the compiled data units in uint64.
func (rs *ReadableSize) GetMeasure() int64 {
	return rs.measure
}

// Get returns the compiled data size in big.Int.
func (rs *ReadableSize) Get() big.Int {
	return *rs.compiled
}

// GetCompiledUInt64 returns the compiled data size in uint64.
// Warning: Possible rounding overflow, use with relatively small numbers.
func (rs *ReadableSize) GetCompiledUInt64() uint64 {
	return rs.compiled.Uint64()
}

// compileMeasuring returns a numeric representation of a data unit in int64.
// See the constant for the allowed options.
func compileMeasuring(measure string) int64 {
	multiplier := int64(1)

	if measure == "" {
		return multiplier
	}

	switch strings.ToLower(string(measure[0])) {
	case "k":
		multiplier = 1 << 10
	case "m":
		multiplier = 1 << 20
	case "g":
		multiplier = 1 << 30
	case "t":
		multiplier = 1 << 40
	case "p":
		multiplier = 1 << 50
	case "e":
		multiplier = 1 << 60
	}

	return multiplier
}

// Compile parses a data size expression and returns, if successful, a ReadableSize object.
// For example: 100MB.
func Compile(input string) (*ReadableSize, error) {
	parser := regexp.MustCompile(sizePattern)

	if matches := parser.FindStringSubmatch(input); len(matches) == 3 {
		if sz, err := strconv.ParseFloat(matches[1], 64); err == nil {
			measure := compileMeasuring(matches[2])
			result, _ := big.NewFloat(sz).Mul(big.NewFloat(sz), big.NewFloat(float64(measure))).Int(new(big.Int))

			return &ReadableSize{
				input:    input,
				measure:  measure,
				compiled: result,
			}, nil
		}
	}

	return nil, errors.New("unsupported data size format")
}

// MustCompile parses a data size expression and returns, if successful,
// a ReadableSize object or returns panic, if an error is found.
// For example: 100MB.
func MustCompile(input string) *ReadableSize {
	res, err := Compile(input)

	if err != nil {
		panic(err)
	}

	return res
}

// ValidateMeasure parses a data size measure and returns true or false.
func ValidateMeasure(format string) bool {
	if format == "" || !regexp.MustCompile(measurePattern).MatchString(format) {
		return false
	}

	return true
}

// BytesToSize parses a number and returns a string of data size format.
// For example: 100MB.
func BytesToSize(size float64, precision uint) (string, error) {
	rounder := func() float64 {
		ratio := math.Pow(10, float64(precision))
		return math.Round(size*ratio) / ratio
	}

	if size == 0 {
		return "0B", nil
	}

	for i, v := range defaultMeasure {
		if size < 1024 || i == len(defaultMeasure)-1 {
			return strconv.FormatFloat(rounder(), 'f', int(precision), 64) + v, nil
		}
		size /= 1 << 10
	}

	return "", errors.New("unable convert")
}

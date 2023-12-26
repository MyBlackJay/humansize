package humansize

import (
	"math/big"
	"testing"
)

func TestCompileMeasure(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("Test completed successfully")
		} else {
			t.Errorf("Test failed. Expected panic")
		}
	}()
	var tests = map[string]struct {
		text     []string
		input    string
		measure  int64
		compiled uint64
		isError  bool
	}{
		"not_valid_m": {
			text:     []string{"MMB"},
			input:    "10",
			measure:  0,
			compiled: 0,
			isError:  true,
		},
	}

	for _, v := range tests {
		for _, mes := range v.text {
			input := v.input + mes
			MustCompile(input)
		}

	}
}

func TestCompile(t *testing.T) {
	var tests = map[string]struct {
		text     []string
		input    string
		measure  int64
		compiled uint64
		isError  bool
	}{
		"without_measure": {
			text:     []string{""},
			input:    "1024",
			measure:  1,
			compiled: 1024,
			isError:  false,
		},
		"empty": {
			text:     []string{"b", "B", "bytes", "Bytes"},
			input:    "",
			measure:  0,
			compiled: 0,
			isError:  true,
		},
		"b": {
			text:     []string{"b", "B", "bytes", "Bytes"},
			input:    "10",
			measure:  1,
			compiled: 10,
			isError:  false,
		},
		"k": {
			text:     []string{"k", "kb", "kib", "K", "KB", "KiB", "KIB"},
			input:    "10",
			measure:  1 << 10,
			compiled: 10 * (1 << 10),
			isError:  false,
		},
		"m": {
			text:     []string{"m", "mb", "mib", "M", "MB", "MiB", "MIB"},
			input:    "10",
			measure:  1 << 20,
			compiled: 10 * (1 << 20),
			isError:  false,
		},
		"g": {
			text:     []string{"g", "gb", "gib", "G", "GB", "GiB", "GIB"},
			input:    "10",
			measure:  1 << 30,
			compiled: 10 * (1 << 30),
			isError:  false,
		},
		"t": {
			text:     []string{"t", "tb", "tib", "T", "TB", "TiB", "TIB"},
			input:    "10",
			measure:  1 << 40,
			compiled: 10 * (1 << 40),
			isError:  false,
		},
		"p": {
			text:     []string{"p", "pb", "pib", "P", "PB", "PiB", "PIB"},
			input:    "10",
			measure:  1 << 50,
			compiled: 10 * (1 << 50),
			isError:  false,
		},
		"e": {
			text:     []string{"e", "eb", "eib", "E", "EB", "EiB", "EIB"},
			input:    "10",
			measure:  1 << 60,
			compiled: big.NewInt(10).Mul(big.NewInt(10), big.NewInt(1<<60)).Uint64(),
			isError:  false,
		},
		"float": {
			text:     []string{"MB"},
			input:    "1000.500",
			measure:  1 << 20,
			compiled: uint64(1000.5 * float64(1<<20)),
			isError:  false,
		},
	}

	for k, v := range tests {
		i := 0
		for _, mes := range v.text {
			input := v.input + mes

			res, err := Compile(input)

			switch {
			case err != nil && !v.isError:
				t.Errorf("Test %s-%d (input: %s) failed. Error was not expected.", k, i, input)
			case v.isError && err == nil:
				t.Errorf("Test %s-%d (input: %s) failed. The error was expected", k, i, input)
			case (err != nil && v.isError) || (err == nil && res.GetCompiledUInt64() == v.compiled && res.measure == v.measure && res.input == input):
				t.Logf("Test %s-%d (input: %s) completed successfully", k, i, input)
			case err == nil && (res.GetCompiledUInt64() != v.compiled || res.measure != v.measure || res.input != input):
				t.Errorf(
					"Test %s-%d (input: %s) failed. Expected: {input: %s, measure: %d, compiled: %d}. Result: {input: %s, measure: %d, compiled: %d}",
					k, i, input, input, v.measure, v.compiled, res.input, res.measure, res.compiled,
				)
			}
			i++
		}

	}
}

func TestValidateMeasure(t *testing.T) {
	tests := []struct {
		input  string
		result bool
	}{
		{"MB", true},
		{"MBN", false},
		{"mb", true},
		{"GiB", true},
		{"pib", true},
		{"b", true},
		{"Bites", false},
	}

	for i, v := range tests {
		if res := ValidateMeasure(v.input); res == v.result {
			t.Logf("Test %d (input: %s) completed successfully", i, v.input)
		} else {
			t.Errorf(
				"Test %d (input: %s) failed. Expected: %t. Result: %t",
				i, v.input, v.result, res,
			)
		}

	}

}

func TestBytesToSize(t *testing.T) {
	tests := []struct {
		input     float64
		precision uint
		result    string
	}{
		{0, 0, "0B"},
		{512, 0, "512B"},
		{2 * 1 << 20, 0, "2MB"},
		{1.5 * float64(1<<30), 1, "1.5GB"},
		{1 << 40, 1, "1.0TB"},
		{0.5 * float64(1<<50), 0, "512TB"},
		{2.596 * float64(1<<60), 2, "2.60EB"},
	}
	for i, v := range tests {
		if res, _ := BytesToSize(v.input, v.precision); res == v.result {
			t.Logf("Test %d (input: %.2f) completed successfully", i, v.input)
		} else {
			t.Errorf(
				"Test %d (input: %.2f) failed. Expected: %s. Result: %s",
				i, v.input, v.result, res,
			)
		}
	}
}

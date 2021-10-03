package writer

import (
	"fmt"
	"os"
	"testing"

	"github.com/nicholaspcr/gde3/pkg/models"
)

func TestNewWriter(t *testing.T) {
	tmp := t.TempDir()
	f, err := os.CreateTemp(tmp, "")
	if err != nil {
		t.Errorf("failed to create temp file for new Writer, %v", err)
	}
	w := NewWriter(f)
	if w == nil {
		t.Errorf("the return of the NewWriter is nil")
	}
}

// TODO CheckFilePath

func TestWriteHeader(t *testing.T) {
	tests := []struct {
		sz       int
		expected string
	}{
		{sz: 5, expected: "header;A;B;C;D;E\n"},
		{
			sz:       27,
			expected: "header;A;B;C;D;E;F;G;H;I;J;K;L;M;N;O;P;Q;R;S;T;U;V;W;X;Y;Z;AA\n",
		},
		{
			sz:       53,
			expected: "header;A;B;C;D;E;F;G;H;I;J;K;L;M;N;O;P;Q;R;S;T;U;V;W;X;Y;Z;AA;AB;AC;AD;AE;AF;AG;AH;AI;AJ;AK;AL;AM;AN;AO;AP;AQ;AR;AS;AT;AU;AV;AW;AX;AY;AZ;BA\n",
		},
	}

	for _, tt := range tests {
		tmp := t.TempDir()
		f, _ := os.CreateTemp(tmp, "")
		defer func() { f.Close() }()

		w := NewWriter(f)
		w.Comma = ';'

		err := w.WriteHeader(tt.sz)
		// checks if write was sucessful
		if err != nil {
			t.Errorf(
				"failed to write header of size %d to file %s",
				tt.sz,
				f.Name(),
			)
		}

		// checks content of file
		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Errorf(
				"Failed to read file %s after write",
				f.Name(),
			)
		}

		if string(b) != tt.expected {
			t.Errorf(
				"error expected %v, received %v",
				tt.expected,
				string(b),
			)
		}
	}
}

func TestElementsObjs(t *testing.T) {

	tests := []struct {
		elems     models.Population
		separator rune
		expected  string
		err       string
	}{
		{
			elems: models.Population{
				models.Vector{Objs: []float64{1.0, 2.0, 3.0}},
			},
			expected: fmt.Sprintf(
				"elem[0],%.8f,%.8f,%.8f\n",
				1.0, 2.0, 3.0,
			),
		},
		{
			elems: models.Population{
				models.Vector{Objs: []float64{0.01, 0.02, 0.03}},
				models.Vector{Objs: []float64{0.004, 0.005, 0.006}},
			},
			separator: ';',
			expected: fmt.Sprintf(
				"elem[0];%.8f;%.8f;%.8f\nelem[1];%.8f;%.8f;%.8f\n",
				0.01, 0.02, 0.03, 0.004, 0.005, 0.006,
			),
		},
		{
			err: "empty slice of elements",
		},
	}

	for _, tt := range tests {
		tmp := t.TempDir()
		f, _ := os.CreateTemp(tmp, "")
		defer func() { f.Close() }()

		w := NewWriter(f)
		if tt.separator != 0 {
			w.Comma = tt.separator
		}

		err := w.ElementsObjs(tt.elems)
		if err != nil && err.Error() != tt.err {
			t.Errorf(
				"failed ElementsObjs, got %v and expected %v",
				err,
				tt.err,
			)
		}

		b, _ := os.ReadFile(f.Name())

		if string(b) != tt.expected {
			t.Errorf(
				"error expected %v, got %v",
				tt.expected,
				string(b),
			)
		}
	}
}

func TestElementsVectors(t *testing.T) {

	tests := []struct {
		elems     models.Population
		separator rune
		expected  string
		err       string
	}{
		{
			elems: models.Population{
				models.Vector{X: []float64{1.0, 2.0, 3.0}},
			},
			expected: fmt.Sprintf(
				"elem[0],%.8f,%.8f,%.8f\n",
				1.0, 2.0, 3.0,
			),
		},
		{
			elems: models.Population{
				models.Vector{X: []float64{0.01, 0.02, 0.03}},
				models.Vector{X: []float64{0.004, 0.005, 0.006}},
			},
			separator: ';',
			expected: fmt.Sprintf(
				"elem[0];%.8f;%.8f;%.8f\nelem[1];%.8f;%.8f;%.8f\n",
				0.01, 0.02, 0.03, 0.004, 0.005, 0.006,
			),
		},
		{
			err: "empty slice of elements",
		},
	}

	for _, tt := range tests {
		tmp := t.TempDir()
		f, _ := os.CreateTemp(tmp, "")
		defer func() { f.Close() }()

		w := NewWriter(f)
		if tt.separator != 0 {
			w.Comma = tt.separator
		}

		err := w.ElementsVectors(tt.elems)
		if err != nil && err.Error() != tt.err {
			t.Errorf(
				"failed ElementsObjs, got %v and expected %v",
				err,
				tt.err,
			)
		}

		b, _ := os.ReadFile(f.Name())

		if string(b) != tt.expected {
			t.Errorf(
				"error expected %v, got %v",
				tt.expected,
				string(b),
			)
		}
	}
}

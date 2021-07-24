package models

// Params of the moDE
type Params struct {
	EXECS       int       `yaml:"execs"`
	NP          int       `yaml:"np"`
	M           int       `yaml:"objs"`
	DIM         int       `yaml:"dim"`
	GEN         int       `yaml:"gen"`
	FLOOR       []float64 `yaml:"floor"`
	CEIL        []float64 `yaml:"ceil"`
	CR          float64   `yaml:"cr-const"`
	F           float64   `yaml:"f-const"`
	P           float64   `yaml:"p-const"`
	DisablePlot bool      `yaml:"plot"`
}

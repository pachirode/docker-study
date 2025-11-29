package options

import (
	"encoding/json"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/pkg/flags"
)

type RunOptions struct {
	TTY    bool         `json:"tty" mapstructure:"tty"`
	MEM    string       `json:"men" mapstructure:"men"`
	CPU    string       `json:"cpu" mapstructure:"cpu"`
	Volume string       `json:"volume" mapstructure:"volume"`
	Detach bool         `json:"detach" mapstructure:"detach"`
	Name   string       `json:"name" mapstructure:"name"`
	Envs   []string     `json:"environment" mapstructure:"environment"`
	Net    string       `json:"net" mapstructure:"net"`
	Port   string       `json:"port" mapstructure:"port"`
	Log    *log.Options `json:"log" mapstructure:"log"`
}

func NewRunOptions() *RunOptions {
	opts := RunOptions{
		TTY:    false,
		MEM:    "",
		CPU:    "",
		Volume: "",
		Detach: false,
		Name:   "",
		Envs:   []string{},
		Net:    "",
		Port:   "",
		Log:    log.NewOptions(),
	}

	return &opts
}

func (opts *RunOptions) Flags() (nfs flags.NamedFlagSets) {

	opts.Log.AddFlags(nfs.GetFlagSet("logs"))

	fs := nfs.GetFlagSet("base")
	fs.BoolVarP(&opts.TTY, "tty", "t", false, "enable tty")
	fs.BoolVarP(&opts.Detach, "detach", "d", false, "enable run backend")

	fs.StringVarP(&opts.Volume, "volume", "v", "", "-v /ect/conf:/etc/conf")
	fs.StringVarP(&opts.Name, "name", "n", "", "-n test-container")
	fs.StringSliceVarP(&opts.Envs, "envs", "e", []string{}, "-e env=test")
	fs.StringVar(&opts.Net, "net", "", "--net testbr")
	fs.StringVarP(&opts.Port, "port", "p", "", "-p 8080:80")

	fsLimit := nfs.GetFlagSet("limit")
	fsLimit.StringVarP(&opts.MEM, "mem", "m", "", "-m 100m")
	fsLimit.StringVar(&opts.CPU, "cpu", "", "--cpu 100")

	return nfs
}

func (opts RunOptions) String() string {
	data, _ := json.Marshal(opts)

	return string(data)
}

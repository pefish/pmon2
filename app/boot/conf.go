package boot

import (
	"github.com/ntt360/errors"
	"github.com/pefish/pmon2/app/conf"
	"gopkg.in/yaml.v2"
	"os"
)

func Conf(confFile string) (*conf.Tpl, error) {
	d, err := os.ReadFile(confFile)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.WriteFile(confFile, []byte(`
data: "/etc/pmon2/data"
logs: "/var/log/pmon2/"
`), 0777)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			return nil, errors.WithStack(err)
		}
	}

	var c conf.Tpl
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	c.Conf = confFile

	return &c, nil
}

package configext

import (
	"github.com/aesoper101/x/jsonschemaext"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/pflag"
	"strings"
)

type PFlagProvider struct {
	p     *posflag.Posflag
	paths []jsonschemaext.Path
}

func NewPFlagProvider(rawSchema []byte, schema *jsonschema.Schema, f *pflag.FlagSet, k *koanf.Koanf) (
	*PFlagProvider,
	error,
) {
	paths, err := getSchemaPaths(rawSchema, schema)
	if err != nil {
		return nil, err
	}
	return &PFlagProvider{
		p:     posflag.Provider(f, ".", k),
		paths: paths,
	}, nil
}

func (p *PFlagProvider) ReadBytes() ([]byte, error) {
	return nil, errors.New("pflag provider does not support this method")
}

func (p *PFlagProvider) Read() (map[string]interface{}, error) {
	all, err := p.p.Read()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	knownFlags := make(map[string]interface{}, len(all))
	for k, v := range all {
		k = strings.ReplaceAll(k, ".", "-")
		for _, path := range p.paths {
			normalized := strings.ReplaceAll(path.Name, ".", "-")
			if k == normalized {
				knownFlags[k] = v
				break
			}
		}
	}
	return knownFlags, nil
}

var _ koanf.Provider = (*PFlagProvider)(nil)
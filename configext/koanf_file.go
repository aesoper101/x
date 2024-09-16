package configext

import (
	"context"
	"github.com/aesoper101/x/watcherext"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// KoanfFile implements a KoanfFile provider.
type KoanfFile struct {
	subKey string
	path   string
	parser koanf.Parser
}

// NewKoanfFile returns a file provider.
func NewKoanfFile(path string) (*KoanfFile, error) {
	return NewKoanfFileSubKey(path, "")
}

func NewKoanfFileSubKey(path, subKey string) (*KoanfFile, error) {
	kf := &KoanfFile{
		path:   filepath.Clean(path),
		subKey: subKey,
	}

	switch e := filepath.Ext(path); e {
	case ".toml":
		kf.parser = toml.Parser()
	case ".json":
		kf.parser = json.Parser()
	case ".yaml", ".yml":
		kf.parser = yaml.Parser()
	default:
		return nil, errors.Errorf("unknown config file extension: %s", e)
	}

	return kf, nil
}

// ReadBytes is not supported by KoanfFile.
func (f *KoanfFile) ReadBytes() ([]byte, error) {
	return nil, errors.New("file provider does not support this method")
}

// Read reads the file and returns the parsed configuration.
func (f *KoanfFile) Read() (map[string]interface{}, error) {
	//#nosec G304 -- false positive
	fc, err := os.ReadFile(f.path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	v, err := f.parser.Unmarshal(fc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if f.subKey == "" {
		return v, nil
	}

	path := strings.Split(f.subKey, Delimiter)
	slices.Reverse(path)

	for _, k := range path {
		v = map[string]interface{}{
			k: v,
		}
	}

	return v, nil
}

// WatchChannel watches the file and triggers a callback when it changes. It is a
// blocking function that internally spawns a goroutine to watch for changes.
func (f *KoanfFile) WatchChannel(ctx context.Context, c watcherext.EventChannel) (watcherext.Watcher, error) {
	return watcherext.WatchFile(ctx, f.path, c)
}

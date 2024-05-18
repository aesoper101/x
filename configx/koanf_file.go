package configx

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"

	"github.com/pkg/errors"
)

type KoanfFile struct {
	subKey string
	path   string
	parser koanf.Parser
}

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

func (f *KoanfFile) ReadBytes() ([]byte, error) {
	return nil, errors.New("file provider does not support this method")
}

func (f *KoanfFile) Read() (map[string]interface{}, error) {
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

func (f *KoanfFile) Watch(cb func(event interface{}, err error)) error {
	// Resolve symlinks and save the original path so that changes to symlinks
	// can be detected.
	realPath, err := filepath.EvalSymlinks(f.path)
	if err != nil {
		return err
	}
	realPath = filepath.Clean(realPath)

	// Although only a single file is being watched, fsnotify has to watch
	// the whole parent directory to pick up all events such as symlink changes.
	fDir, _ := filepath.Split(f.path)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	var (
		lastEvent     string
		lastEventTime time.Time
	)

	go func() {
	loop:
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					cb(nil, errors.New("fsnotify watch channel closed"))
					break loop
				}

				// Use a simple timer to buffer events as certain events fire
				// multiple times on some platforms.
				if event.String() == lastEvent && time.Since(lastEventTime) < time.Millisecond*5 {
					continue
				}
				lastEvent = event.String()
				lastEventTime = time.Now()

				evFile := filepath.Clean(event.Name)

				// Since the event is triggered on a directory, is this
				// one on the file being watched?
				if evFile != realPath && evFile != f.path {
					continue
				}

				// The file was removed.
				if event.Op&fsnotify.Remove != 0 {
					cb(nil, fmt.Errorf("file %s was removed", event.Name))
					break loop
				}

				// Resolve symlink to get the real path, in case the symlink's
				// target has changed.
				curPath, err := filepath.EvalSymlinks(f.path)
				if err != nil {
					cb(nil, err)
					break loop
				}
				realPath = filepath.Clean(curPath)

				// Finally, we only care about create and write.
				if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
					continue
				}

				// Trigger event.
				cb(nil, nil)

			// There's an error.
			case err, ok := <-w.Errors:
				if !ok {
					cb(nil, errors.New("fsnotify err channel closed"))
					break loop
				}

				// Pass the error to the callback.
				cb(nil, err)
				break loop
			}
		}

		_ = w.Close()
	}()

	return w.Add(fDir)
}

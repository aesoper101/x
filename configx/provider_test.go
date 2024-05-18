package configx

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type ServerConfig struct {
	Provider TestProvider `json:"provider" koanf:"provider"`
}

type TestProvider struct {
	Name   string     `json:"name" koanf:"name"`
	Config TestConfig `json:"config" koanf:"config"`
}

func (t *TestProvider) UnmarshalJSON(bytes []byte) error {
	type Data struct {
		Name   string          `json:"name" koanf:"name"`
		Config json.RawMessage `json:"config" koanf:"config"`
	}
	fmt.Println("============UnmarshalJSON==============")
	var data Data
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	switch data.Name {
	case "test1":
		config := TestConfig1{}
		err := json.Unmarshal(data.Config, &config)
		if err != nil {
			return err
		}
		t.Name = data.Name
		t.Config = &config
	case "test2":
		config := TestConfig2{}
		err := json.Unmarshal(data.Config, &config)
		if err != nil {
			return err
		}
		t.Name = data.Name
		t.Config = &config
	default:
		return fmt.Errorf("unknown provider name: %s", data.Name)
	}

	return nil
}

type TestConfig interface {
	open()
}

type TestConfig1 struct {
	Name1 string `json:"name" koanf:"name"`
}

func (t *TestConfig1) UnmarshalJSON(bytes []byte) error {
	type alias TestConfig1
	var data alias
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	*t = TestConfig1(data)
	// t.open()
	fmt.Println("UnmarshalJSON1")
	return nil
}

func (t *TestConfig1) open() {
	fmt.Println("open")
}

type TestConfig2 struct {
	Name2 string `json:"name" koanf:"name"`
}

func (t *TestConfig2) open() {
}

func TestNew(t *testing.T) {
	provider, err := New(context.Background(), WithBaseValues(map[string]interface{}{
		//"provider": map[string]interface{}{
		//	"name": "test1",
		//	"config": map[string]interface{}{
		//		"name": "test1",
		//	},
		//},
		"provider.name":        "test2",
		"provider.config.name": "test1",
	}))
	if err != nil {
		t.Error(err)
		return
	}
	var config ServerConfig
	if err := provider.Unmarshal("", &config); err != nil {
		t.Error(err)
		return
	}
	t.Log(provider.Keys())
	t.Log(config)

	bytes, err := json.Marshal(config.Provider.Config)
	if err != nil {
		return
	}
	t.Log(string(bytes))
	t.Log(reflect.TypeOf(config.Provider.Config))
}

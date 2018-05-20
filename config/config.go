package config

import (
    "fmt"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Processor *ProcessorConfig
    Handler   *HandlerConfigs
}

type ProcessorConfig struct {
    Lines       [20][]int
    Reels       [32][]string
    Prices      map[string][4]int64
    MutableCell string
    EmptyCell   string
}

type HandlerConfigs struct {
    FreeRollsByOneAddiction int
    Port                    int
    SpinURL                 string
    Secret                  string
}

const (
    processorPath = "config/processor.toml"
    handlerPath   = "config/handler.toml"
)

func New() (*Config, error) {
    pc := new(ProcessorConfig)
    if _, err := toml.DecodeFile(processorPath, pc); err != nil {
        return nil, fmt.Errorf("can't decode config file '%s': %v", processorPath, err)
    }

    hc := new(HandlerConfigs)
    if _, err := toml.DecodeFile(handlerPath, hc); err != nil {
        return nil, fmt.Errorf("can't decode config file '%s': %v", handlerPath, err)
    }

    c := new(Config)
    c.Processor = pc
    c.Handler = hc

    return c, nil
}

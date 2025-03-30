package kafka

import "strings"

type KgoKafkaConsumerConfig struct {
	Addr                 []string `json:"addr" mapstructure:"addr"`
	Domain               string   `json:"domain" mapstructure:"domain"`
	Project              string   `json:"project" mapstructure:"project"`
	Group                string   `json:"group" mapstructure:"group"`
	SubGroup             string   `json:"sub_group" mapstructure:"group"`
	WorkersEachPartition int      `json:"workers_each_partition" mapstructure:"workers_each_partition"`
	MaxRetries           int      `json:"max_retries" mapstructure:"max_retries"`
}

func (cfg *KgoKafkaConsumerConfig) GetGroup() string {
	if len(cfg.SubGroup) != 0 {
		return strings.Join([]string{cfg.Group, cfg.SubGroup}, ".")
	}
	return cfg.Group
}

type KgoKafkaPublisherConfig struct {
	Domain  string   `json:"domain" mapstructure:"domain"`
	Project string   `json:"project" mapstructure:"project"`
	Addr    []string `json:"addr" mapstructure:"addr"`
}

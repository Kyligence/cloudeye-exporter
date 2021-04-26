package collector

import (
	"io/ioutil"
	"reflect"
	"strings"

	rds "github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
	"gopkg.in/yaml.v2"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
)

type CloudAuth struct {
	ProjectName string `yaml:"project_name"`
	ProjectID   string `yaml:"project_id,omitempty"`
	DomainName  string `yaml:"domain_name,omitempty"`
	AccessKey   string `yaml:"access_key,omitempty"`
	Region      string `yaml:"region"`
	SecretKey   string `yaml:"secret_key,omitempty"`
	AuthURL     string `yaml:"auth_url"`
	UserName    string `yaml:"user_name,omitempty"`
	Password    string `yaml:"password,omitempty"`
	IsEncrypt   bool   `yaml:"is_encrypt"`
}

type Global struct {
	Port        string `yaml:"port"`
	Prefix      string `yaml:"prefix"`
	MetricPath  string `yaml:"metric_path"`
	MaxRoutines int    `yaml:"max_routines"`
}

type Custom struct {
	NamePrefix string `yaml:"name_prefix"`
	TagKey     string `yaml:"tag_key"`
	TagValue   string `yaml:"tag_value"`
	TtlMinute  int64  `yaml:"ttl_minute"`
}

type CloudConfig struct {
	Auth   CloudAuth `yaml:"auth"`
	Global Global    `yaml:"global"`
	Custom Custom    `yaml:"custom"`
}

func NewCloudConfigFromFile(file string) (*CloudConfig, error) {
	var config CloudConfig

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if !config.Auth.IsEncrypt {
		//var akDecryptErr, skDecryptErr error
		ak_decrypt, ak_err := Encrypt("Ee9TlFUNJzNuzRev", 1000, config.Auth.AccessKey)
		sk_decrypt, sk_err := Encrypt("Ee9TlFUNJzNuzRev", 1000, config.Auth.SecretKey)
		if ak_err != nil {
			return nil, ak_err
		}
		if sk_err != nil {
			return nil, sk_err
		}
		config.Auth.AccessKey = ak_decrypt
		config.Auth.SecretKey = sk_decrypt
		config.Auth.IsEncrypt = true
		data, yaml_err := yaml.Marshal(&config)
		if yaml_err != nil {
			return nil, yaml_err
		}
		err := ioutil.WriteFile(file, data, 0644)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	SetDefaultConfigValues(&config)

	return &config, err
}

func SetDefaultConfigValues(config *CloudConfig) {
	if config.Global.Port == "" {
		config.Global.Port = ":8087"
	}

	if config.Global.MetricPath == "" {
		config.Global.MetricPath = "/metrics"
	}

	if config.Global.Prefix == "" {
		config.Global.Prefix = "huaweicloud"
	}

	if config.Global.MaxRoutines == 0 {
		config.Global.MaxRoutines = 20
	}
}

var filterConfigMap map[string]map[string][]string

func InitFilterConfig(enable bool) error {
	filterConfigMap = make(map[string]map[string][]string)
	if !enable {
		return nil
	}

	data, err := ioutil.ReadFile("metric_filter_config.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &filterConfigMap)
	if err != nil {
		return err
	}
	return nil
}

func getMetricConfigMap(namespace string) map[string][]string {
	if configMap, ok := filterConfigMap[namespace]; ok {
		return configMap
	}
	return nil
}

func startWith(s string, e string) bool {
	return strings.HasPrefix(s, e)
}

func containsRDS(s []rds.RdsInstanceResponse, e rds.RdsInstanceResponse) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}

func containsLB(s []loadbalancers.LoadBalancer, e loadbalancers.LoadBalancer) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}

func containsListener(s []listeners.Listener, e listeners.Listener) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

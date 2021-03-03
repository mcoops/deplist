package scan

import (
	"testing"
)

var want map[string]string = map[string]string{
	"fluent-plugin-splunk-hec":          "1.1.2",
	"http-form_data":                    "2.3.0",
	"prometheus-client":                 "0.9.0",
	"tzinfo":                            "2.0.2",
	"fluent-plugin-rewrite-tag-filter":  "2.3.0",
	"fluent-plugin-kafka":               "0.13.1",
	"fluent-plugin-prometheus":          "1.7.3",
	"http_parser.rb":                    "0.6.0",
	"lru_redux":                         "1.1.0",
	"rest-client":                       "2.1.0",
	"tzinfo-data":                       "1.2020.1",
	"unf":                               "0.1.4",
	"aws-eventstream":                   "1.1.0",
	"yajl-ruby":                         "1.4.1",
	"jmespath":                          "1.4.0",
	"http-cookie":                       "1.0.3",
	"fluent-plugin-multi-format-parser": "1.0.0",
	"elasticsearch":                     "7.8.0",
	"fluent-plugin-elasticsearch":       "4.1.1",
	"fluent-plugin-kubernetes_metadata_filter": "2.5.2",
	"dig_rb":                           "1.0.1",
	"elasticsearch-transport":          "7.8.0",
	"ffi":                              "1.11.3",
	"mime-types":                       "3.3.1",
	"to_regexp":                        "0.2.1",
	"aws-sdk-cloudwatchlogs":           "1.38.0",
	"msgpack":                          "1.3.3",
	"typhoeus":                         "1.4.0",
	"digest-crc":                       "0.6.1",
	"http-parser":                      "1.2.1",
	"ruby-kafka":                       "1.1.0",
	"serverengine":                     "2.2.1",
	"http":                             "4.4.1",
	"ethon":                            "0.12.0",
	"multipart-post":                   "2.1.1",
	"concurrent-ruby":                  "1.1.6",
	"fluent-plugin-concat":             "2.4.0",
	"public_suffix":                    "4.0.5",
	"sigdump":                          "0.2.4",
	"syslog_protocol":                  "0.9.2",
	"aws-sigv4":                        "1.2.2",
	"fluent-plugin-record-modifier":    "2.1.0",
	"addressable":                      "2.7.0",
	"elasticsearch-api":                "7.8.0",
	"excon":                            "0.75.0",
	"fluent-plugin-cloudwatch-logs":    "0.7.6",
	"fluent-plugin-systemd":            "1.0.2",
	"ltsv":                             "0.1.2",
	"quantile":                         "0.2.1",
	"recursive-open-struct":            "1.1.2",
	"cool.io":                          "1.6.0",
	"unf_ext":                          "0.0.7.7",
	"strptime":                         "0.2.4",
	"fluent-plugin-remote-syslog":      "1.1",
	"fluent-config-regexp-type":        "1.0.0",
	"fluent-mixin-config-placeholders": "0.4.0",
	"fluentd":                          "1.7.4",
	"http-accept":                      "1.7.0",
	"systemd-journal":                  "1.3.3",
	"uuidtools":                        "2.1.5",
	"aws-sdk-core":                     "3.109.3",
	"connection_pool":                  "2.2.3",
	"domain_name":                      "0.5.20190701",
	"faraday":                          "1.0.1",
	"jsonpath":                         "1.0.5",
	"mime-types-data":                  "3.2020.0512",
	"netrc":                            "0.11.0",
	"aws-partitions":                   "1.396.0",
	"kubeclient":                       "4.8.0",
	"multi_json":                       "1.15.0",
	"net-http-persistent":              "3.1.0",
	"rake":                             "13.0.1",
	"ffi-compiler":                     "1.0.1",
}

func Test_GetRubyDeps(t *testing.T) {
	got, err := GetRubyDeps("../../test/testRepo/")

	if err != nil {
		t.Errorf("GetRubyDeps() error thrown: %+v", err)
	}

	for k, v := range want {
		if got[k] != v {
			t.Errorf("GetRubyDeps() - deps missing entry: %s = %s ", k, v)
		}
	}

	if len(got) != len(want) {
		t.Errorf("GetRubyDeps() = %d; want %d", len(got), len(want))
	}
}

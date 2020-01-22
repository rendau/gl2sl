package main

import (
	"github.com/rendau/gl2sl/internal/domain/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	cr *core.St
)

func TestMain(t *testing.M) {
	viper.SetConfigFile("conf_test.yml")
	_ = viper.ReadInConfig()

	cr = core.NewSt(
		viper.GetString("slack_webhook_url"),
		viper.GetString("slack_channel"),
		viper.GetString("graylog_link"),
	)

	exitCode := t.Run()

	os.Exit(exitCode)
}

func TestMsg(t *testing.T) {
	err := cr.HandleMessage([]byte(`
{
  "event_definition_id":"5e27eff173c7b20012b88aee",
  "event_definition_type":"aggregation-v1",
  "event_definition_title":"All",
  "event_definition_description":"",
  "job_definition_id":"5e27ef3073c7b20012b88a14",
  "job_trigger_id":"5e27f6f473c7b20012b892ad",
  "event":{
    "id":"01DZ61SATP0GHJ77GJKYWZVT46",
    "event_definition_type":"aggregation-v1",
    "event_definition_id":"5e27eff173c7b20012b88aee",
    "origin_context":"urn:graylog:message:es:graylog_0:185d48a0-3ce7-11ea-90f1-0242ac1b0003",
    "timestamp":"2020-01-22T07:16:23.156Z",
    "timestamp_processing":"2020-01-22T07:17:08.566Z",
    "timerange_start":null,
    "timerange_end":null,
    "streams":[

    ],
    "source_streams":[
      "000000000000000000000001"
    ],
    "message":"All",
    "source":"5f46a6a39d5f",
    "key_tuple":[

    ],
    "key":"",
    "priority":2,
    "alert":true,
    "fields":{

    }
  },
  "backlog":[
    {
      "index":"graylog_0",
      "message":"{\"level\":\"info\",\"ts\":1579677383.1533563,\"caller\":\"cmd/root.go:131\",\"msg\":\"Started\",\"http_listen\":\":80\"}",
      "timestamp":"2020-01-22T07:16:23.156Z",
      "source":"server",
      "stream_ids":[
        "000000000000000000000001"
      ],
      "fields":{
        "level":3,
        "gl2_remote_ip":"172.27.0.1",
        "created":"2020-01-19T19:17:01.225474494Z",
        "gl2_remote_port":37682,
        "z_caller":"cmd/root.go:131",
        "gl2_message_id":"01DZ61QYHEYGRT8MHRW168B034",
        "z_level":"info",
        "gl2_source_input":"5e08d5c379d4b40012c72d44",
        "command":"./svc",
        "image_name":"docker.pkg.github.com/medleader/medleader_api/medleader_api:latest",
        "container_name":"api",
        "z_http_listen":":80",
        "z_msg":"Started",
        "gl2_source_node":"81bf79e5-8fa1-4bc7-acf4-c53e32ba2cb6",
        "tag":"847ffaf661db",
        "z_ts":1.5796773831533563E9,
        "image_id":"sha256:af467ab4780b6eb639283c9cf6ee43977c8c9c5a5ac6e12a57349c2675676156",
        "container_id":"847ffaf661dbdef6e17342f4f90bf32dc6eb1f1ad8f9ae64b53c854a4197e976"
      },
      "id":"185d48a0-3ce7-11ea-90f1-0242ac1b0003"
    }
  ]
}
`))
	require.Nil(t, err)
}

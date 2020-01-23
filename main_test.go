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
  "event_definition_id":"5e299011fe6fcc0012cfff2e",
  "event_definition_type":"aggregation-v1",
  "event_definition_title":"All",
  "event_definition_description":"",
  "job_definition_id":"5e298fc5fe6fcc0012cffed2",
  "job_trigger_id":"5e29a115fe6fcc0012d011c9",
  "event":{
    "id":"01DZ99TFERTM76YE9S94W9TBY8",
    "event_definition_type":"aggregation-v1",
    "event_definition_id":"5e299011fe6fcc0012cfff2e",
    "origin_context":"urn:graylog:message:es:graylog_0:20b61370-3de5-11ea-9e9c-0242c0a8020a",
    "timestamp":"2020-01-23T13:34:49.373Z",
    "timestamp_processing":"2020-01-23T13:35:17.983Z",
    "timerange_start":null,
    "timerange_end":null,
    "streams":[

    ],
    "source_streams":[
      "000000000000000000000001"
    ],
    "message":"All",
    "source":"fe510147d8f7",
    "key_tuple":[
      ""
    ],
    "key":"",
    "priority":2,
    "alert":true,
    "fields":{
      "msgId":""
    }
  },
  "backlog":[
    {
      "index":"graylog_0",
      "message":"{\"level\":\"warn\",\"ts\":1579786489.372348,\"caller\":\"tst/handler.go:11\",\"msg\":\"Warn1\"}",
      "timestamp":"2020-01-23T13:34:49.373Z",
      "source":"server",
      "stream_ids":[
        "000000000000000000000001"
      ],
      "fields":{
        "level":3,
        "created":"2020-01-23T11:30:00.563152616Z",
        "gl2_remote_ip":"192.168.2.1",
        "gl2_remote_port":39408,
        "z_caller":"tst/handler.go:11",
        "gl2_message_id":"01DZ99SKH8EEXDENMASAE6KWFQ",
        "z_level":"warn",
        "gl2_source_input":"5e298d26fe6fcc0012cffbef",
        "command":"./svc",
        "image_name":"docker.pkg.github.com/medleader/medleader_api/medleader_api:latest",
        "container_name":"api",
        "z_msg":"Warn1",
        "gl2_source_node":"a0b84683-0993-4104-872c-6fa9598dd8f0",
        "tag":"ab14a377c677",
        "image_id":"sha256:71f1ab8d2774c49a9deed90bf8a9d654c936ab80e59a084d78f1dbc6351de7f4",
        "z_ts":1.579786489372348E9,
        "container_id":"ab14a377c677caeb5fe8db1fbddd8616ec9647b5e14b75359109910d4cde0bbd"
      },
      "id":"20b61370-3de5-11ea-9e9c-0242c0a8020a"
    }
  ]
}
`))
	require.Nil(t, err)
}

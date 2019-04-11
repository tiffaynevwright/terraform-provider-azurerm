package azure

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
)

func SchemaStreamAnalyticsStreamInputSerialization() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"format": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(streamanalytics.TypeAvro),
						string(streamanalytics.TypeCsv),
						string(streamanalytics.TypeJSON),
					}, false),
				},

				"field_delimiter": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						" ",
						",",
						"	",
						"|",
						";",
					}, false),
				},

				"encoding": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(streamanalytics.UTF8),
					}, false),
				},
			},
		},
	}
}

func ExpandStreamAnalyticsStreamInputSerialization(input []interface{}) (streamanalytics.BasicSerialization, error) {
	v := input[0].(map[string]interface{})

	format := streamanalytics.Type(v["format"].(string))
	encoding := v["encoding"].(string)
	fieldDelimiter := v["field_delimiter"].(string)

	switch format {
	case streamanalytics.TypeAvro:
		return streamanalytics.AvroSerialization{
			Type:       streamanalytics.TypeAvro,
			Properties: map[string]interface{}{},
		}, nil

	case streamanalytics.TypeCsv:
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `format` is set to `Csv`")
		}
		if fieldDelimiter == "" {
			return nil, fmt.Errorf("`field_delimiter` must be set when `format` is set to `Csv`")
		}
		return streamanalytics.CsvSerialization{
			Type: streamanalytics.TypeCsv,
			CsvSerializationProperties: &streamanalytics.CsvSerializationProperties{
				Encoding:       streamanalytics.Encoding(encoding),
				FieldDelimiter: utils.String(fieldDelimiter),
			},
		}, nil

	case streamanalytics.TypeJSON:
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `format` is set to `JSON`")
		}

		return streamanalytics.JSONSerialization{
			Type: streamanalytics.TypeJSON,
			JSONSerializationProperties: &streamanalytics.JSONSerializationProperties{
				Encoding: streamanalytics.Encoding(encoding),
			},
		}, nil
	}

	return nil, fmt.Errorf("Unsupported Format %q", format)
}

func FlattenStreamAnalyticsStreamInputSerialization(input streamanalytics.BasicSerialization) []interface{} {
	var encoding string
	var format string

	if _, ok := input.AsAvroSerialization(); ok {
		format = string(streamanalytics.TypeAvro)
	}

	if v, ok := input.AsCsvSerialization(); ok {
		if props := v.CsvSerializationProperties; props != nil {
			encoding = string(props.Encoding)
		}

		format = string(streamanalytics.TypeCsv)
	}

	if v, ok := input.AsJSONSerialization(); ok {
		if props := v.JSONSerializationProperties; props != nil {
			encoding = string(props.Encoding)
		}

		format = string(streamanalytics.TypeJSON)
	}

	return []interface{}{
		map[string]interface{}{
			"encoding": encoding,
			"format":   format,
		},
	}
}

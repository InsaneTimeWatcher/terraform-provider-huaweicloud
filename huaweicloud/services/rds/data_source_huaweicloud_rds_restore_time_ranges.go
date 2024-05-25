// Generated by PMS #169
package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceRdsRestoreTimeRanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsRestoreTimeRangesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of RDS instance.`,
			},
			"date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the date to be queried.`,
			},
			"restore_time": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of restoration time ranges.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the start time of the restoration time range in the UNIX timestamp format.`,
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the end time of the restoration time range in the UNIX timestamp format.`,
						},
					},
				},
			},
		},
	}
}

type RestoreTimeRangesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRestoreTimeRangesDSWrapper(d *schema.ResourceData, meta interface{}) *RestoreTimeRangesDSWrapper {
	return &RestoreTimeRangesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRdsRestoreTimeRangesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRestoreTimeRangesDSWrapper(d, meta)
	lisResTimRst, err := wrapper.ListRestoreTimes()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listRestoreTimesToSchema(lisResTimRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API RDS GET /v3/{project_id}/instances/{instance_id}/restore-time
func (w *RestoreTimeRangesDSWrapper) ListRestoreTimes() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/instances/{instance_id}/restore-time"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"date": w.Get("date"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OkCode(200).
		Request().
		Result()
}

func (w *RestoreTimeRangesDSWrapper) listRestoreTimesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("restore_time", schemas.SliceToList(body.Get("restore_time"),
			func(restoreTime gjson.Result) any {
				return map[string]any{
					"start_time": restoreTime.Get("start_time").Value(),
					"end_time":   restoreTime.Get("end_time").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
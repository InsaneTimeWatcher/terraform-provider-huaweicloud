// Generated by PMS #175
package cts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCtsTrackers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCtsTrackersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"tracker_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the tracker ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the tracker name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the tracker type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the tracker status.`,
			},
			"data_bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the data bucket name.`,
			},
			"trackers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `List of tracker information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique tracker ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tracker name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tracker type.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the tracker was created. The time is in UTC.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tracker status.`,
						},
						"kms_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the key used for trace file encryption.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the account that the tracker belongs to.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The project ID.`,
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `This parameter is returned only when the tracker status is **error**.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the LTS log group.`,
						},
						"stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the LTS log stream.`,
						},
						"data_bucket": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Information about the bucket tracked by a data tracker.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_bucket_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The OBS bucket name.`,
									},
									"search_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the logs of the tracked bucket can be searched.`,
									},
									"data_event": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The list of the bucket event operation types.`,
									},
								},
							},
						},
						"obs_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Information about the bucket to which traces are transferred.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The OBS bucket name.`,
									},
									"file_prefix_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `File name prefix to mark trace files that need to be stored in an OBS bucket.`,
									},
									"is_obs_created": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the OBS bucket is automatically created by the tracker.`,
									},
									"is_authorized_bucket": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether CTS has been granted permissions to perform operations on the OBS bucket.`,
									},
									"bucket_lifecycle": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Duration that traces are stored in the OBS bucket.`,
									},
									"compress_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The compression type.`,
									},
									"is_sort_by_service": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether to sort the path by cloud service.`,
									},
								},
							},
						},
						"is_support_validate": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether trace file verification is enabled.`,
						},
						"lts": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The LTS configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_lts_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether traces are synchronized to LTS for trace search and analysis.`,
									},
									"log_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the log group that CTS creates in LTS.`,
									},
									"log_topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the log stream that CTS creates in LTS.`,
									},
								},
							},
						},
						"is_support_trace_files_encryption": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether trace files are encrypted during transfer to an OBS bucket.`,
						},
						"is_organization_tracker": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether system tracker to apply to my organization.`,
						},
						"management_event_selector": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The management event selector.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclude_service": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The cloud service that is not dumped.`,
									},
								},
							},
						},
						"agency_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of a cloud service agency.`,
						},
					},
				},
			},
		},
	}
}

type TrackersDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTrackersDSWrapper(d *schema.ResourceData, meta interface{}) *TrackersDSWrapper {
	return &TrackersDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCtsTrackersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTrackersDSWrapper(d, meta)
	listTrackersRst, err := wrapper.ListTrackers()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listTrackersToSchema(listTrackersRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CTS GET /v3/{project_id}/trackers
func (w *TrackersDSWrapper) ListTrackers() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cts")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/trackers"
	params := map[string]any{
		"tracker_name": w.Get("name"),
		"tracker_type": w.Get("type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Filter(
			filters.New().From("trackers").
				Where("data_bucket.data_bucket_name", "=", w.Get("data_bucket_name")).
				Where("id", "=", w.Get("tracker_id")).
				Where("status", "=", w.Get("status")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *TrackersDSWrapper) listTrackersToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("trackers", schemas.SliceToList(body.Get("trackers"),
			func(trackers gjson.Result) any {
				return map[string]any{
					"id":          trackers.Get("id").Value(),
					"name":        trackers.Get("tracker_name").Value(),
					"type":        trackers.Get("tracker_type").Value(),
					"create_time": w.setTraCreTim(trackers),
					"status":      trackers.Get("status").Value(),
					"kms_id":      trackers.Get("kms_id").Value(),
					"domain_id":   trackers.Get("domain_id").Value(),
					"project_id":  trackers.Get("project_id").Value(),
					"detail":      trackers.Get("detail").Value(),
					"group_id":    trackers.Get("group_id").Value(),
					"stream_id":   trackers.Get("stream_id").Value(),
					"data_bucket": schemas.SliceToList(trackers.Get("data_bucket"),
						func(dataBucket gjson.Result) any {
							return map[string]any{
								"data_bucket_name": dataBucket.Get("data_bucket_name").Value(),
								"search_enabled":   dataBucket.Get("search_enabled").Value(),
								"data_event":       schemas.SliceToStrList(dataBucket.Get("data_event")),
							}
						},
					),
					"obs_info": schemas.SliceToList(trackers.Get("obs_info"),
						func(obsInfo gjson.Result) any {
							return map[string]any{
								"bucket_name":          obsInfo.Get("bucket_name").Value(),
								"file_prefix_name":     obsInfo.Get("file_prefix_name").Value(),
								"is_obs_created":       obsInfo.Get("is_obs_created").Value(),
								"is_authorized_bucket": obsInfo.Get("is_authorized_bucket").Value(),
								"bucket_lifecycle":     obsInfo.Get("bucket_lifecycle").Value(),
								"compress_type":        obsInfo.Get("compress_type").Value(),
								"is_sort_by_service":   obsInfo.Get("is_sort_by_service").Value(),
							}
						},
					),
					"is_support_validate": trackers.Get("is_support_validate").Value(),
					"lts": schemas.SliceToList(trackers.Get("lts"),
						func(lts gjson.Result) any {
							return map[string]any{
								"is_lts_enabled": lts.Get("is_lts_enabled").Value(),
								"log_group_name": lts.Get("log_group_name").Value(),
								"log_topic_name": lts.Get("log_topic_name").Value(),
							}
						},
					),
					"is_support_trace_files_encryption": trackers.Get("is_support_trace_files_encryption").Value(),
					"is_organization_tracker":           trackers.Get("is_organization_tracker").Value(),
					"management_event_selector": schemas.SliceToList(trackers.Get("management_event_selector"),
						func(manEveSel gjson.Result) any {
							return map[string]any{
								"exclude_service": schemas.SliceToStrList(manEveSel.Get("exclude_service")),
							}
						},
					),
					"agency_name": trackers.Get("agency_name").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*TrackersDSWrapper) setTraCreTim(data gjson.Result) string {
	return utils.FormatTimeStampUTC(int64(data.Get("create_time").Float() / 1000))
}
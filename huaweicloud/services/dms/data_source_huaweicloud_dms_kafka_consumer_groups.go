// Generated by PMS #199
package dms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const pageLimit = 10

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/groups
func DataSourceDmsKafkaConsumerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsKafkaConsumerGroupsRead,

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
				Description: `Specifies the instance ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the consumer group name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the consumer group description.`,
			},
			"lag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the the number of accumulated messages.`,
			},
			"coordinator_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the coordinator ID.`,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the consumer group status.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the consumer groups.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the consumer group name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the consumer group description.`,
						},
						"lag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of accumulated messages.`,
						},
						"coordinator_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the coordinator ID.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the consumer group status.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDmsKafkaConsumerGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	listGroupsHttpUrl := "v2/{project_id}/instances/{instance_id}/groups"
	listGroupsPath := client.Endpoint + listGroupsHttpUrl
	listGroupsPath = strings.ReplaceAll(listGroupsPath, "{project_id}", client.ProjectID)
	listGroupsPath = strings.ReplaceAll(listGroupsPath, "{instance_id}", d.Get("instance_id").(string))
	listGroupsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// pagelimit is `10`
	listGroupsPath += fmt.Sprintf("?limit=%v", pageLimit)
	listGroupsPath = buildQueryGroupsListPath(d, listGroupsPath)

	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listGroupsPath + fmt.Sprintf("&offset=%d", currentTotal)
		listGroupsResp, err := client.Request("GET", currentPath, &listGroupsOpt)
		if err != nil {
			return diag.Errorf("error retrieving groups: %s", err)
		}
		listGroupsRespBody, err := utils.FlattenResponse(listGroupsResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		groups := utils.PathSearch("groups", listGroupsRespBody, make([]interface{}, 0)).([]interface{})
		for _, group := range groups {
			// filter result
			description := utils.PathSearch("group_desc", group, "").(string)
			lag := int64(utils.PathSearch("lag", group, float64(0)).(float64))
			coordinatorID := int64(utils.PathSearch("coordinator_id", group, float64(0)).(float64))
			state := utils.PathSearch("state", group, "").(string)
			if val, ok := d.GetOk("description"); ok && description != val {
				continue
			}
			if val, ok := d.GetOk("lag"); ok && lag != val {
				continue
			}
			if val, ok := d.GetOk("coordinator_id"); ok && coordinatorID != val {
				continue
			}
			if val, ok := d.GetOk("description"); ok && state != val {
				continue
			}
			results = append(results, map[string]interface{}{
				"name":           utils.PathSearch("group_id", group, nil),
				"state":          utils.PathSearch("state", group, nil),
				"lag":            utils.PathSearch("lag", group, 0),
				"coordinator_id": utils.PathSearch("coordinator_id", group, 0),
				"description":    utils.PathSearch("group_desc", group, nil),
				"created_at": utils.FormatTimeStampRFC3339(
					int64(utils.PathSearch("createdAt", group, float64(0)).(float64))/1000, true),
			})
		}

		// `totalCount` means the number of all `groups`, and type is float64.
		currentTotal += len(groups)
		totalCount := utils.PathSearch("total", listGroupsRespBody, float64(0))
		if int(totalCount.(float64)) == currentTotal {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryGroupsListPath(d *schema.ResourceData, listGroupsPath string) string {
	if instId, ok := d.GetOk("name"); ok {
		listGroupsPath += fmt.Sprintf("&group=%s", instId)
	}
	return listGroupsPath
}

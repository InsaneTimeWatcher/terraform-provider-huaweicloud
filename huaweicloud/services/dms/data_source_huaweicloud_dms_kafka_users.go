// Generated by PMS #203
package dms

import (
	"context"
	"strings"

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

func DataSourceDmsKafkaUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsKafkaUsersRead,

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
				Description: `Specifies the user name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the user description.`,
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the user list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the username.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description.`,
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user role.`,
						},
						"default_app": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the application is the default application.`,
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

type KafkaUsersDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newKafkaUsersDSWrapper(d *schema.ResourceData, meta interface{}) *KafkaUsersDSWrapper {
	return &KafkaUsersDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDmsKafkaUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newKafkaUsersDSWrapper(d, meta)
	showInstanceUsersRst, err := wrapper.ShowInstanceUsers()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.showInstanceUsersToSchema(showInstanceUsersRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/users
func (w *KafkaUsersDSWrapper) ShowInstanceUsers() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dmsv2")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/instances/{instance_id}/users"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Filter(
			filters.New().From("users").
				Where("user_name", "contains", w.Get("name")).
				Where("user_desc", "contains", w.Get("description")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *KafkaUsersDSWrapper) showInstanceUsersToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("users", schemas.SliceToList(body.Get("users"),
			func(users gjson.Result) any {
				return map[string]any{
					"name":        users.Get("user_name").Value(),
					"description": users.Get("user_desc").Value(),
					"role":        users.Get("role").Value(),
					"default_app": users.Get("default_app").Value(),
					"created_at":  w.setUsersCreatedTime(users),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*KafkaUsersDSWrapper) setUsersCreatedTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339((data.Get("created_time").Int())/1000, true)
}

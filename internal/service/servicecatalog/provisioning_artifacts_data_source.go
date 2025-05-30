// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicecatalog

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	awstypes "github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_servicecatalog_provisioning_artifacts", name="Provisioning Artifacts")
func dataSourceProvisioningArtifacts() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceProvisioningArtifactsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(ConstraintReadTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Default:      acceptLanguageEnglish,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(acceptLanguage_Values(), false),
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provisioning_artifact_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						names.AttrCreatedTime: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrDescription: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guidance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrType: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProvisioningArtifactsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogClient(ctx)

	productID := d.Get("product_id").(string)
	input := &servicecatalog.ListProvisioningArtifactsInput{
		AcceptLanguage: aws.String(d.Get("accept_language").(string)),
		ProductId:      aws.String(productID),
	}

	output, err := conn.ListProvisioningArtifacts(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing Service Catalog Provisioning Artifacts: %s", err)
	}

	d.SetId(productID)
	if err := d.Set("provisioning_artifact_details", flattenProvisioningArtifactDetails(output.ProvisioningArtifactDetails)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting provisioning_artifact_details: %s", err)
	}

	return diags
}

func flattenProvisioningArtifactDetails(apiObjects []awstypes.ProvisioningArtifactDetail) []any {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []any

	for _, apiObject := range apiObjects {
		tfList = append(tfList, flattenProvisioningArtifactDetail(apiObject))
	}

	return tfList
}

func flattenProvisioningArtifactDetail(apiObject awstypes.ProvisioningArtifactDetail) map[string]any {
	tfMap := map[string]any{}

	if apiObject.Active != nil {
		tfMap["active"] = aws.ToBool(apiObject.Active)
	}
	if apiObject.CreatedTime != nil {
		tfMap[names.AttrCreatedTime] = aws.ToTime(apiObject.CreatedTime).String()
	}
	if apiObject.Description != nil {
		tfMap[names.AttrDescription] = aws.ToString(apiObject.Description)
	}

	tfMap["guidance"] = string(apiObject.Guidance)

	if apiObject.Id != nil {
		tfMap[names.AttrID] = aws.ToString(apiObject.Id)
	}
	if apiObject.Name != nil {
		tfMap[names.AttrName] = aws.ToString(apiObject.Name)
	}

	tfMap[names.AttrType] = string(apiObject.Type)

	return tfMap
}

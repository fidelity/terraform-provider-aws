package fsx

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKDataSource("aws_fsx_ontap_storage_virtual_machine", name="Ontap Storage Virtual Machine")
func DataSourceOntapStorageVirtualMachine() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceOntapStorageVirtualMachineRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_directory_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"netbios_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_managed_active_directory_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_ips": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"domain_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_system_administrators_group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"organizational_unit_distinguished_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iscsi": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"management": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"nfs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"smb": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": DataSoureStorageVirtualMachineFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lifecycle_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_transition_reason": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_volume_security_style": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subtype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

const (
	DSNameOntapStorageVirtualMachine = "Ontap Storage Virtual Machine Data Source"
)

func dataSourceOntapStorageVirtualMachineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &fsx.DescribeStorageVirtualMachinesInput{}

	if id, ok := d.GetOk("id"); ok {
		input.StorageVirtualMachineIds = []*string{aws.String(id.(string))}
	}

	input.Filters = BuildStorageVirtualMachineFiltersDataSource(
		d.Get("filter").(*schema.Set),
	)

	if len(input.Filters) == 0 {
		input.Filters = nil
	}

	svm, err := FindStorageVirtualMachine(ctx, conn, input)

	if err != nil {
		return sdkdiag.AppendFromErr(diags, tfresource.SingularDataSourceFindError("FSx StorageVirtualMachine", err))
	}

	d.SetId(aws.StringValue(svm.StorageVirtualMachineId))

	d.Set("arn", svm.ResourceARN)
	d.Set("endpoints", flattenOntapStorageVirtualMachineEndpoints(svm.Endpoints))
	d.Set("file_system_id", svm.FileSystemId)
	d.Set("id", svm.StorageVirtualMachineId)
	d.Set("lifecycle_status", svm.Lifecycle)
	d.Set("lifecycle_transition_reason", flattenOntapSvmLifecycleTransitionReason(svm.LifecycleTransitionReason))
	d.Set("name", svm.Name)
	d.Set("root_volume_security_style", svm.RootVolumeSecurityStyle)
	d.Set("subtype", svm.Subtype)
	d.Set("uuid", svm.UUID)

	if err := d.Set("active_directory_configuration", flattenOntapSvmActiveDirectoryConfiguration(d, svm.ActiveDirectoryConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting svm_active_directory: %s", err)
	}

	if err := d.Set("creation_time", svm.CreationTime.Format(time.RFC3339)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting creation_time: %s", err)
	}

	tags := KeyValueTags(ctx, svm.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}

func flattenOntapSvmLifecycleTransitionReason(rs *fsx.LifecycleTransitionReason) []interface{} {
	if rs == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})
	if rs.Message != nil {
		m["message"] = aws.StringValue(rs.Message)
	}

	return []interface{}{m}
}

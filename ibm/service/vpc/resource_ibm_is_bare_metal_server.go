// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerAction                  = "action"
	isBareMetalServerBandwidth               = "bandwidth"
	isBareMetalServerBootTarget              = "boot_target"
	isBareMetalServerCreatedAt               = "created_at"
	isBareMetalServerCPU                     = "cpu"
	isBareMetalServerCPUArchitecture         = "architecture"
	isBareMetalServerCPUCoreCount            = "core_count"
	isBareMetalServerCpuSocketCount          = "socket_count"
	isBareMetalServerCpuThreadPerCore        = "threads_per_core"
	isBareMetalServerCRN                     = "crn"
	isBareMetalServerDisks                   = "disks"
	isBareMetalServerDiskID                  = "id"
	isBareMetalServerDiskSize                = "size"
	isBareMetalServerDiskName                = "name"
	isBareMetalServerDiskInterfaceType       = "interface_type"
	isBareMetalServerHref                    = "href"
	isBareMetalServerMemory                  = "memory"
	isBareMetalServerTags                    = "tags"
	isBareMetalServerName                    = "name"
	isBareMetalServerNetworkInterfaces       = "network_interfaces"
	isBareMetalServerPrimaryNetworkInterface = "primary_network_interface"
	isBareMetalServerProfile                 = "profile"
	isBareMetalServerResourceGroup           = "resource_group"
	isBareMetalServerResourceType            = "resource_type"
	isBareMetalServerStatus                  = "status"
	isBareMetalServerStatusReasons           = "status_reasons"
	isBareMetalServerVPC                     = "vpc"
	isBareMetalServerZone                    = "zone"
	isBareMetalServerStatusReasonsCode       = "code"
	isBareMetalServerStatusReasonsMessage    = "message"
	isBareMetalServerStatusReasonsMoreInfo   = "more_info"
	isBareMetalServerDeleteType              = "delete_type"
	isBareMetalServerImage                   = "image"
	isBareMetalServerKeys                    = "keys"
	isBareMetalServerUserData                = "user_data"
	isBareMetalServerNicName                 = "name"
	isBareMetalServerNicPortSpeed            = "port_speed"
	isBareMetalServerNicAllowIPSpoofing      = "allow_ip_spoofing"
	isBareMetalServerNicSecurityGroups       = "security_groups"
	isBareMetalServerNicSubnet               = "subnet"
	isBareMetalServerUserAccounts            = "user_accounts"
	isBareMetalServerActionDeleting          = "deleting"
	isBareMetalServerActionDeleted           = "deleted"
	isBareMetalServerActionStatusStopping    = "stopping"
	isBareMetalServerActionStatusStopped     = "stopped"
	isBareMetalServerActionStatusStarting    = "starting"
	isBareMetalServerStatusRunning           = "running"
	isBareMetalServerStatusPending           = "pending"
	isBareMetalServerStatusRestarting        = "restarting"
	isBareMetalServerStatusFailed            = "failed"
)

func ResourceIBMIsBareMetalServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISBareMetalServerCreate,
		ReadContext:   resourceIBMISBareMetalServerRead,
		UpdateContext: resourceIBMISBareMetalServerUpdate,
		DeleteContext: resourceIBMISBareMetalServerDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
		),

		Schema: map[string]*schema.Schema{

			isBareMetalServerName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerName),
				Description:  "Bare metal server name",
			},

			isBareMetalServerAction: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerAction),
				Description:  "This restart/start/stops a bare metal server.",
			},
			isBareMetalServerBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total bandwidth (in megabits per second)",
			},
			isBareMetalServerBootTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this bare metal server disk",
			},

			isBareMetalServerCPU: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The bare metal server CPU configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerCPUArchitecture: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CPU architecture",
						},
						isBareMetalServerCPUCoreCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of cores",
						},
						isBareMetalServerCpuSocketCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of CPU sockets",
						},
						isBareMetalServerCpuThreadPerCore: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of hardware threads per core",
						},
					},
				},
			},
			isBareMetalServerCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this bare metal server",
			},
			isBareMetalServerDisks: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The disks for this bare metal server, including any disks that are associated with the boot_target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerDiskHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this bare metal server disk",
						},
						isBareMetalServerDiskID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this bare metal server disk",
						},
						isBareMetalServerDiskInterfaceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the disk. Supported values are [ nvme, sata ]",
						},
						isBareMetalServerDiskName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this disk",
						},
						isBareMetalServerDiskResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type",
						},
						isBareMetalServerDiskSize: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes)",
						},
					},
				},
			},
			isBareMetalServerHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this bare metal server",
			},
			isBareMetalServerMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of memory, truncated to whole gibibytes",
			},
			isBareMetalServerDeleteType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "hard",
				Description: "Enables stopping type of the bare metal server before deleting",
			},
			isBareMetalServerPrimaryNetworkInterface: {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network interface",
						},
						isBareMetalServerNicEnableInfraNAT: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If true, the VPC infrastructure performs any needed NAT operations. If false, the packet is passed unmodified to/from the network interface, allowing the workload to perform any needed NAT operations.",
						},
						isBareMetalServerNicInterfaceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network interface type: [ pci, vlan ]",
						},
						isBareMetalServerNicPrimaryIP: {
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    0,
							MaxItems:    1,
							Computed:    true,
							Description: "title: IPv4, The IP address. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerNicIpAddress: {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.reserved_ip"},
										Description:   "The globally unique IP address",
									},
									isBareMetalServerNicIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isBareMetalServerNicIpAutoDelete: {
										Type:          schema.TypeBool,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.reserved_ip"},
										Description:   "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isBareMetalServerNicIpName: {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.reserved_ip"},
										Description:   "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isBareMetalServerNicIpID: {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.address", "primary_network_interface.0.primary_ip.0.auto_delete", "primary_network_interface.0.primary_ip.0.name"},
										Description:   "Identifies a reserved IP by a unique property.",
									},
									isBareMetalServerNicResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
						isBareMetalServerNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isBareMetalServerNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isBareMetalServerNicPortSpeed: {
							Type:     schema.TypeInt,
							Computed: true,
						},

						isBareMetalServerNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isBareMetalServerNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						isBareMetalServerNicAllowedVlans: {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Set:         schema.HashInt,
							Description: "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) interface. A given VLAN can only be in the allowed_vlans array for one PCI type adapter per bare metal server.",
						},
					},
				},
			},

			isBareMetalServerNetworkInterfaces: {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isBareMetalServerNicName: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The user-defined name for this network interface. If unspecified, the name will be a hyphenated list of randomly-selected words",
						},
						isBareMetalServerNicHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network interface",
						},
						isBareMetalServerNicEnableInfraNAT: {
							Type:             schema.TypeBool,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "If true, the VPC infrastructure performs any needed NAT operations. If false, the packet is passed unmodified to/from the network interface, allowing the workload to perform any needed NAT operations.",
						},
						isBareMetalServerNicInterfaceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network interface type: [ pci, vlan ]",
						},
						isBareMetalServerNicPrimaryIP: {
							Type:        schema.TypeList,
							MinItems:    0,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "title: IPv4, The IP address. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerNicIpAddress: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The globally unique IP address",
									},
									isBareMetalServerNicIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isBareMetalServerNicIpAutoDelete: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isBareMetalServerNicIpName: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isBareMetalServerNicIpID: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
									isBareMetalServerNicResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
						isBareMetalServerNicSecurityGroups: {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Collection of security group ids",
						},
						isBareMetalServerNicSubnet: {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         false,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "The associated subnet",
						},
						isBareMetalServerNicAllowedVlans: {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Set:         schema.HashInt,
							Description: "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) interface. A given VLAN can only be in the allowed_vlans array for one PCI type adapter per bare metal server.",
						},

						isBareMetalServerNicAllowInterfaceToFloat: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates if the interface can float to any other server within the same resource_group. The interface will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to vlan type interfaces.",
						},

						isBareMetalServerNicVlan: {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this interface",
						},
					},
				},
			},

			isBareMetalServerKeys: {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "SSH key Ids for the bare metal server",
			},

			isBareMetalServerImage: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "image id",
			},
			isBareMetalServerProfile: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "profile name",
			},

			isBareMetalServerUserData: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "User data given for the bare metal server",
			},

			isBareMetalServerZone: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name",
			},

			isBareMetalServerVPC: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The VPC the bare metal server is to be a part of",
			},

			isBareMetalServerResourceGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Resource group name",
			},
			isBareMetalServerResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource type name",
			},

			isBareMetalServerStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bare metal server status",
			},

			isBareMetalServerStatusReasons: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerStatusReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason",
						},

						isBareMetalServerStatusReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},
						isBareMetalServerStatusReasonsMoreInfo: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason",
						},
					},
				},
			},
			isBareMetalServerTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", "tag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Tags for the Bare metal server",
			},
		},
	}
}

func ResourceIBMIsBareMetalServerValidator() *validate.ResourceValidator {
	bareMetalServerActions := "start, restart, stop"
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              bareMetalServerActions})
	ibmISBareMetalServerResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_bare_metal_server", Schema: validateSchema}
	return &ibmISBareMetalServerResourceValidator
}

func resourceIBMISBareMetalServerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.CreateBareMetalServerOptions{}
	var imageStr string
	if image, ok := d.GetOk(isBareMetalServerImage); ok {
		imageStr = image.(string)
	}
	keySet := d.Get(isBareMetalServerKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		options.Initialization = &vpcv1.BareMetalServerInitializationPrototype{
			Image: &vpcv1.ImageIdentity{
				ID: &imageStr,
			},
			Keys: keyobjs,
		}
		if userdata, ok := d.GetOk(isBareMetalServerUserData); ok {
			userdatastr := userdata.(string)
			options.Initialization.UserData = &userdatastr
		}
	}

	if name, ok := d.GetOk(isBareMetalServerName); ok {
		nameStr := name.(string)
		options.Name = &nameStr
	}

	if primnicintf, ok := d.GetOk(isBareMetalServerPrimaryNetworkInterface); ok && len(primnicintf.([]interface{})) > 0 {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isBareMetalServerNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.BareMetalServerPrimaryNetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isBareMetalServerNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}

		if primaryIpIntf, ok := primnic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
			primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})

			reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
			if ok && reservedIpIdOk.(string) != "" {
				ipid := reservedIpIdOk.(string)
				primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &ipid,
				}
			} else {
				primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

				reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
				if okAdd && reservedIpAddressOk.(string) != "" {
					reservedIpAddress := reservedIpAddressOk.(string)
					primaryip.Address = &reservedIpAddress
				}

				reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
				if okName && reservedIpNameOk.(string) != "" {
					reservedIpName := reservedIpNameOk.(string)
					primaryip.Name = &reservedIpName
				}
				reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
				if okAuto {
					reservedIpAuto := reservedIpAutoOk.(bool)
					primaryip.AutoDelete = &reservedIpAuto
				}
				if okAdd || okName || okAuto {
					primnicobj.PrimaryIP = primaryip
				}
			}
		}

		allowIPSpoofing, ok := primnic[isBareMetalServerNicAllowIPSpoofing]

		if ok && allowIPSpoofing != nil {
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if allowIPSpoofingbool {
				primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
			}
		}
		enableInfraNATbool := true
		enableInfraNAT, ok := primnic[isBareMetalServerNicEnableInfraNAT]
		if ok && enableInfraNAT != nil {
			enableInfraNATbool = enableInfraNAT.(bool)
			primnicobj.EnableInfrastructureNat = &enableInfraNATbool
		}

		secgrpintf, ok := primnic[isBareMetalServerNicSecurityGroups]
		if ok {
			secgrpSet := secgrpintf.(*schema.Set)
			if secgrpSet.Len() != 0 {
				var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
				for i, secgrpIntf := range secgrpSet.List() {
					secgrpIntfstr := secgrpIntf.(string)
					secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
						ID: &secgrpIntfstr,
					}
				}
				primnicobj.SecurityGroups = secgrpobjs
			}
		}

		if allowedVlansOk, ok := primnic[isBareMetalServerNicAllowedVlans]; ok {
			allowedVlansList := allowedVlansOk.(*schema.Set).List()

			allowedVlans := make([]int64, 0, len(allowedVlansList))
			for _, k := range allowedVlansList {
				allowedVlans = append(allowedVlans, int64(k.(int)))
			}
			primnicobj.AllowedVlans = allowedVlans
			interfaceType := "pci"
			primnicobj.InterfaceType = &interfaceType
		}
		options.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isBareMetalServerNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			interfaceType := ""
			if allowedVlansOk, ok := nic[isBareMetalServerNicAllowedVlans]; ok {
				interfaceType = "pci"
				var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByPciPrototype{}
				nicobj.InterfaceType = &interfaceType

				allowedVlansList := allowedVlansOk.(*schema.Set).List()

				allowedVlans := make([]int64, 0, len(allowedVlansList))
				for _, k := range allowedVlansList {
					allowedVlans = append(allowedVlans, int64(k.(int)))
				}
				nicobj.AllowedVlans = allowedVlans

				subnetintf, _ := nic[isBareMetalServerNicSubnet]
				subnetintfstr := subnetintf.(string)
				nicobj.Subnet = &vpcv1.SubnetIdentity{
					ID: &subnetintfstr,
				}
				name, _ := nic[isBareMetalServerNicName]
				namestr := name.(string)
				if namestr != "" {
					nicobj.Name = &namestr
				}

				enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
				enableInfraNATbool := enableInfraNAT.(bool)
				if ok {
					nicobj.EnableInfrastructureNat = &enableInfraNATbool
				}

				if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
					primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})

					reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
					if ok && reservedIpIdOk.(string) != "" {
						ipid := reservedIpIdOk.(string)
						nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
							ID: &ipid,
						}
					} else {
						primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
						reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
						if okAdd && reservedIpAddressOk.(string) != "" {
							reservedIpAddress := reservedIpAddressOk.(string)
							primaryip.Address = &reservedIpAddress
						}
						reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
						if okName && reservedIpNameOk.(string) != "" {
							reservedIpName := reservedIpNameOk.(string)
							primaryip.Name = &reservedIpName
						}
						reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
						if okAuto {
							reservedIpAuto := reservedIpAutoOk.(bool)
							primaryip.AutoDelete = &reservedIpAuto
						}
						if okAdd || okName || okAuto {
							nicobj.PrimaryIP = primaryip
						}
					}

				}

				allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
				allowIPSpoofingbool := allowIPSpoofing.(bool)
				if ok && allowIPSpoofingbool {
					nicobj.AllowIPSpoofing = &allowIPSpoofingbool
				}
				secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
				if ok {
					secgrpSet := secgrpintf.(*schema.Set)
					if secgrpSet.Len() != 0 {
						var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
						for i, secgrpIntf := range secgrpSet.List() {
							secgrpIntfstr := secgrpIntf.(string)
							secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
								ID: &secgrpIntfstr,
							}
						}
						nicobj.SecurityGroups = secgrpobjs
					}
				}
			} else {
				interfaceType = "vlan"
				var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByVlanPrototype{}
				nicobj.InterfaceType = &interfaceType

				if aitf, ok := nic[isBareMetalServerNicAllowInterfaceToFloat]; ok {
					allowInterfaceToFloat := aitf.(bool)
					nicobj.AllowInterfaceToFloat = &allowInterfaceToFloat
				}
				if vlan, ok := nic[isBareMetalServerNicVlan]; ok {
					vlanInt := int64(vlan.(int))
					nicobj.Vlan = &vlanInt
				}

				subnetintf, _ := nic[isBareMetalServerNicSubnet]
				subnetintfstr := subnetintf.(string)
				nicobj.Subnet = &vpcv1.SubnetIdentity{
					ID: &subnetintfstr,
				}
				name, _ := nic[isBareMetalServerNicName]
				namestr := name.(string)
				if namestr != "" {
					nicobj.Name = &namestr
				}

				enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
				enableInfraNATbool := enableInfraNAT.(bool)
				if ok {
					nicobj.EnableInfrastructureNat = &enableInfraNATbool
				}

				if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
					primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
					reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
					if ok && reservedIpIdOk.(string) != "" {
						ipid := reservedIpIdOk.(string)
						nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
							ID: &ipid,
						}
					} else {
						primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

						reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
						if okAdd && reservedIpAddressOk.(string) != "" {
							reservedIpAddress := reservedIpAddressOk.(string)
							primaryip.Address = &reservedIpAddress
						}

						reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
						if okName && reservedIpNameOk.(string) != "" {
							reservedIpName := reservedIpNameOk.(string)
							primaryip.Name = &reservedIpName
						}

						reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
						if okAuto {
							reservedIpAuto := reservedIpAutoOk.(bool)
							primaryip.AutoDelete = &reservedIpAuto
						}
						if okAdd || okName || okAuto {
							nicobj.PrimaryIP = primaryip
						}
					}
				}

				allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
				allowIPSpoofingbool := allowIPSpoofing.(bool)
				if ok && allowIPSpoofingbool {
					nicobj.AllowIPSpoofing = &allowIPSpoofingbool
				}
				secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
				if ok {
					secgrpSet := secgrpintf.(*schema.Set)
					if secgrpSet.Len() != 0 {
						var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
						for i, secgrpIntf := range secgrpSet.List() {
							secgrpIntfstr := secgrpIntf.(string)
							secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
								ID: &secgrpIntfstr,
							}
						}
						nicobj.SecurityGroups = secgrpobjs
					}
				}
			}
		}
	}

	if rgrp, ok := d.GetOk(isBareMetalServerResourceGroup); ok {
		rg := rgrp.(string)
		options.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	if p, ok := d.GetOk(isBareMetalServerProfile); ok {
		profile := p.(string)
		options.Profile = &vpcv1.BareMetalServerProfileIdentity{
			Name: &profile,
		}
	}

	if z, ok := d.GetOk(isBareMetalServerZone); ok {
		zone := z.(string)
		options.Zone = &vpcv1.ZoneIdentity{
			Name: &zone,
		}
	}

	if v, ok := d.GetOk(isBareMetalServerVPC); ok {
		vpc := v.(string)
		options.VPC = &vpcv1.VPCIdentity{
			ID: &vpc,
		}
	}

	bms, response, err := sess.CreateBareMetalServerWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Create bare metal server err %s\n%s", err, response))
	}
	d.SetId(*bms.ID)
	log.Printf("[INFO] Bare Metal Server : %s", *bms.ID)
	_, err = isWaitForBareMetalServerAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(err)
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isBareMetalServerTags); ok || v != "" {
		oldList, newList := d.GetChange(isBareMetalServerTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *bms.CRN)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource bare metal server (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceIBMISBareMetalServerRead(context, d, meta)
}

func resourceIBMISBareMetalServerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	err := bareMetalServerGet(context, d, meta, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func bareMetalServerGet(context context.Context, d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.GetBareMetalServerOptions{
		ID: &id,
	}
	bms, response, err := sess.GetBareMetalServerWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s): %s\n%s", id, err, response)
	}
	d.SetId(*bms.ID)
	d.Set(isBareMetalServerBandwidth, bms.Bandwidth)
	bmsBootTargetIntf := bms.BootTarget.(*vpcv1.BareMetalServerBootTarget)
	bmsBootTarget := bmsBootTargetIntf.ID
	d.Set(isBareMetalServerBootTarget, bmsBootTarget)
	cpuList := make([]map[string]interface{}, 0)
	if bms.Cpu != nil {
		currentCPU := map[string]interface{}{}
		currentCPU[isBareMetalServerCPUArchitecture] = *bms.Cpu.Architecture
		currentCPU[isBareMetalServerCPUCoreCount] = *bms.Cpu.CoreCount
		currentCPU[isBareMetalServerCpuSocketCount] = *bms.Cpu.SocketCount
		currentCPU[isBareMetalServerCpuThreadPerCore] = *bms.Cpu.ThreadsPerCore
		cpuList = append(cpuList, currentCPU)
	}
	d.Set(isBareMetalServerCPU, cpuList)
	d.Set(isBareMetalServerCRN, *bms.CRN)

	diskList := make([]map[string]interface{}, 0)
	if bms.Disks != nil {
		for _, disk := range bms.Disks {
			currentDisk := map[string]interface{}{
				isBareMetalServerDiskHref:          disk.Href,
				isBareMetalServerDiskID:            disk.ID,
				isBareMetalServerDiskInterfaceType: disk.InterfaceType,
				isBareMetalServerDiskName:          disk.Name,
				isBareMetalServerDiskResourceType:  disk.ResourceType,
				isBareMetalServerDiskSize:          disk.Size,
			}
			diskList = append(diskList, currentDisk)
		}
	}
	d.Set(isBareMetalServerDisks, diskList)
	d.Set(isBareMetalServerHref, *bms.Href)
	d.Set(isBareMetalServerMemory, *bms.Memory)
	d.Set(isBareMetalServerName, *bms.Name)
	//pni

	if bms.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic["id"] = *bms.PrimaryNetworkInterface.ID
		currentPrimNic[isBareMetalServerNicName] = *bms.PrimaryNetworkInterface.Name
		currentPrimNic[isBareMetalServerNicHref] = *bms.PrimaryNetworkInterface.Href
		currentPrimNic[isBareMetalServerNicSubnet] = *bms.PrimaryNetworkInterface.Subnet.ID
		getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
			BareMetalServerID: &id,
			ID:                bms.PrimaryNetworkInterface.ID,
		}
		bmsnic, response, err := sess.GetBareMetalServerNetworkInterfaceWithContext(context, getnicoptions)

		if err != nil {
			return fmt.Errorf("[ERROR] Error getting network interfaces attached to the bare metal server %s\n%s", err, response)
		}

		if bms.PrimaryNetworkInterface.PrimaryIP != nil {
			primaryIpList := make([]map[string]interface{}, 0)
			currentIP := map[string]interface{}{}
			if bms.PrimaryNetworkInterface.PrimaryIP.Href != nil {
				currentIP[isBareMetalServerNicIpAddress] = *bms.PrimaryNetworkInterface.PrimaryIP.Address
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.Href != nil {
				currentIP[isBareMetalServerNicIpHref] = *bms.PrimaryNetworkInterface.PrimaryIP.Href
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.Name != nil {
				currentIP[isBareMetalServerNicIpName] = *bms.PrimaryNetworkInterface.PrimaryIP.Name
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.ID != nil {
				currentIP[isBareMetalServerNicIpID] = *bms.PrimaryNetworkInterface.PrimaryIP.ID
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.ResourceType != nil {
				currentIP[isBareMetalServerNicResourceType] = *bms.PrimaryNetworkInterface.PrimaryIP.ResourceType
			}

			primaryIpList = append(primaryIpList, currentIP)
			currentPrimNic[isBareMetalServerNicPrimaryIP] = primaryIpList

			getripoptions := &vpcv1.GetSubnetReservedIPOptions{
				SubnetID: bms.PrimaryNetworkInterface.Subnet.ID,
				ID:       bms.PrimaryNetworkInterface.PrimaryIP.ID,
			}
			bmsRip, response, err := sess.GetSubnetReservedIP(getripoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error getting network interface reserved ip(%s) attached to the bare metal server primary network interface(%s): %s\n%s", *bms.PrimaryNetworkInterface.PrimaryIP.ID, *bms.PrimaryNetworkInterface.ID, err, response)
			}
			currentIP[isBareMetalServerNicIpAutoDelete] = bmsRip.AutoDelete
		}
		switch reflect.TypeOf(bmsnic).String() {
		case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
			{
				primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
				currentPrimNic[isBareMetalServerNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
				currentPrimNic[isBareMetalServerNicEnableInfraNAT] = *primNic.EnableInfrastructureNat
				currentPrimNic[isBareMetalServerNicPortSpeed] = *primNic.PortSpeed
				if len(primNic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(primNic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}

				if primNic.AllowedVlans != nil {
					var out = make([]interface{}, len(primNic.AllowedVlans), len(primNic.AllowedVlans))
					for i, v := range primNic.AllowedVlans {
						out[i] = int(v)
					}
					currentPrimNic[isBareMetalServerNicAllowedVlans] = schema.NewSet(schema.HashInt, out)
				}
			}
		case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
			{
				primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
				currentPrimNic[isBareMetalServerNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
				currentPrimNic[isBareMetalServerNicEnableInfraNAT] = *primNic.EnableInfrastructureNat

				if len(primNic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(primNic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}
			}
		}

		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isBareMetalServerPrimaryNetworkInterface, primaryNicList)
	}

	//ni

	interfacesList := make([]map[string]interface{}, 0)
	for _, intfc := range bms.NetworkInterfaces {
		if *intfc.ID != *bms.PrimaryNetworkInterface.ID {
			currentNic := map[string]interface{}{}
			currentNic["id"] = *intfc.ID
			currentNic[isBareMetalServerNicName] = *intfc.Name
			getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
				BareMetalServerID: &id,
				ID:                intfc.ID,
			}
			bmsnicintf, response, err := sess.GetBareMetalServerNetworkInterfaceWithContext(context, getnicoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error getting network interfaces attached to the bare metal server %s\n%s", err, response)
			}
			if intfc.PrimaryIP != nil {
				primaryIpList := make([]map[string]interface{}, 0)
				currentIP := map[string]interface{}{}
				if intfc.PrimaryIP.Href != nil {
					currentIP[isBareMetalServerNicIpAddress] = *intfc.PrimaryIP.Address
				}
				if intfc.PrimaryIP.Href != nil {
					currentIP[isBareMetalServerNicIpHref] = *intfc.PrimaryIP.Href
				}
				if intfc.PrimaryIP.Name != nil {
					currentIP[isBareMetalServerNicIpName] = *intfc.PrimaryIP.Name
				}
				if intfc.PrimaryIP.ID != nil {
					currentIP[isBareMetalServerNicIpID] = *intfc.PrimaryIP.ID
				}
				if intfc.PrimaryIP.ResourceType != nil {
					currentIP[isBareMetalServerNicResourceType] = *intfc.PrimaryIP.ResourceType
				}
				getripoptions := &vpcv1.GetSubnetReservedIPOptions{
					SubnetID: bms.PrimaryNetworkInterface.Subnet.ID,
					ID:       bms.PrimaryNetworkInterface.PrimaryIP.ID,
				}
				bmsRip, response, err := sess.GetSubnetReservedIP(getripoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error getting network interface reserved ip(%s) attached to the bare metal server network interface(%s): %s\n%s", *bms.PrimaryNetworkInterface.PrimaryIP.ID, *bms.PrimaryNetworkInterface.ID, err, response)
				}
				currentIP[isBareMetalServerNicIpAutoDelete] = bmsRip.AutoDelete

				primaryIpList = append(primaryIpList, currentIP)
				currentNic[isBareMetalServerNicPrimaryIP] = primaryIpList
			}

			switch reflect.TypeOf(bmsnicintf).String() {
			case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
				{
					bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
					currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
					currentNic[isBareMetalServerNicEnableInfraNAT] = *bmsnic.EnableInfrastructureNat
					currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
					currentNic[isBareMetalServerNicPortSpeed] = *bmsnic.PortSpeed
					currentNic[isBareMetalServerNicInterfaceType] = "pci"
					if len(bmsnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(bmsnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
						}
						currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
				{
					bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
					currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
					currentNic[isBareMetalServerNicEnableInfraNAT] = *bmsnic.EnableInfrastructureNat
					currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
					currentNic[isBareMetalServerNicPortSpeed] = *bmsnic.PortSpeed
					currentNic[isBareMetalServerNicInterfaceType] = "vlan"

					if len(bmsnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(bmsnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
						}
						currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			}
			interfacesList = append(interfacesList, currentNic)
		}
	}
	d.Set(isBareMetalServerNetworkInterfaces, interfacesList)

	d.Set(isBareMetalServerProfile, *bms.Profile.Name)
	if bms.ResourceGroup != nil {
		d.Set(isBareMetalServerResourceGroup, *bms.ResourceGroup.ID)
	}
	d.Set(isBareMetalServerResourceType, bms.ResourceType)
	d.Set(isBareMetalServerStatus, *bms.Status)
	statusReasonsList := make([]map[string]interface{}, 0)
	if bms.StatusReasons != nil {
		for _, sr := range bms.StatusReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR[isBareMetalServerStatusReasonsCode] = *sr.Code
				currentSR[isBareMetalServerStatusReasonsMessage] = *sr.Message
				if sr.MoreInfo != nil {
					currentSR[isBareMetalServerStatusReasonsMoreInfo] = *sr.MoreInfo
				}
				statusReasonsList = append(statusReasonsList, currentSR)
			}
		}
	}
	d.Set(isBareMetalServerStatusReasons, statusReasonsList)
	d.Set(isBareMetalServerVPC, *bms.VPC.ID)
	d.Set(isBareMetalServerZone, *bms.Zone.Name)

	tags, err := flex.GetTagsUsingCRN(meta, *bms.CRN)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource bare metal server (%s) tags: %s", d.Id(), err)
	}
	d.Set(isBareMetalServerTags, tags)

	return nil
}

func resourceIBMISBareMetalServerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id := d.Id()

	err := bareMetalServerUpdate(context, d, meta, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMISBareMetalServerRead(context, d, meta)
}

func bareMetalServerUpdate(context context.Context, d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	if d.HasChange(isBareMetalServerTags) {
		bmscrn := d.Get(isBareMetalServerCRN).(string)
		if bmscrn == "" {
			options := &vpcv1.GetBareMetalServerOptions{
				ID: &id,
			}
			bms, response, err := sess.GetBareMetalServerWithContext(context, options)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s): %s\n%s", id, err, response)
			}
			bmscrn = *bms.CRN
		}
		oldList, newList := d.GetChange(isBareMetalServerTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, bmscrn)
		if err != nil {
			log.Printf(
				"Error on update of vpc Bare metal server (%s) tags: %s", id, err)
		}
	}

	options := &vpcv1.UpdateBareMetalServerOptions{
		ID: &id,
	}
	bmsPatchModel := &vpcv1.BareMetalServerPatch{}
	flag := false

	if d.HasChange(isBareMetalServerPrimaryNetworkInterface) {
		nicId := d.Get("primary_network_interface.0.id").(string)

		if d.HasChange("primary_network_interface.0.primary_ip.0.name") || d.HasChange("primary_network_interface.0.primary_ip.0.auto_delete") {
			subnetId := d.Get("primary_network_interface.0.subnet").(string)
			ripId := d.Get("primary_network_interface.0.primary_ip.0.reserved_ip").(string)
			updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
				SubnetID: &subnetId,
				ID:       &ripId,
			}
			reservedIpPath := &vpcv1.ReservedIPPatch{}
			if d.HasChange("primary_network_interface.0.primary_ip.0.name") {
				name := d.Get("primary_network_interface.0.primary_ip.0.name").(string)
				reservedIpPath.Name = &name
			}
			if d.HasChange("primary_network_interface.0.primary_ip.0.auto_delete") {
				auto := d.Get("primary_network_interface.0.primary_ip.0.auto_delete").(bool)
				reservedIpPath.AutoDelete = &auto
			}
			reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
			if err != nil {
				return fmt.Errorf("[ERROR] Error calling reserved ip as patch \n%s", err)
			}
			updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
			_, response, err := sess.UpdateSubnetReservedIP(updateripoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating bare metal server primary network interface reserved ip(%s): %s\n%s", ripId, err, response)
			}
		}
		bmsNicUpdateOptions := &vpcv1.UpdateBareMetalServerNetworkInterfaceOptions{
			BareMetalServerID: &id,
			ID:                &nicId,
		}
		bmsNicPatchModel := &vpcv1.BareMetalServerNetworkInterfacePatch{}
		if d.HasChange("primary_network_interface.0.allowed_vlans") {
			if allowedVlansOk, ok := d.GetOk("primary_network_interface.0.allowed_vlans"); ok {
				allowedVlansList := allowedVlansOk.(*schema.Set).List()
				allowedVlans := make([]int64, 0, len(allowedVlansList))
				for _, k := range allowedVlansList {
					allowedVlans = append(allowedVlans, int64(k.(int)))
				}
				bmsNicPatchModel.AllowedVlans = allowedVlans
			}
		}
		if d.HasChange("primary_network_interface.0.allow_ip_spoofing") {

			if allowIpSpoofingOk, ok := d.GetOk("primary_network_interface.0.allow_ip_spoofing"); ok {
				allowIpSpoofing := allowIpSpoofingOk.(bool)
				if allowIpSpoofing {
					bmsNicPatchModel.AllowIPSpoofing = &allowIpSpoofing
				}
			}
		}
		if d.HasChange("primary_network_interface.0.enable_infrastructure_nat") {
			if enableNatOk, ok := d.GetOk("primary_network_interface.0.enable_infrastructure_nat"); ok {
				enableNat := enableNatOk.(bool)
				bmsNicPatchModel.EnableInfrastructureNat = &enableNat
			}
		}
		if d.HasChange("primary_network_interface.0.name") {
			if nameOk, ok := d.GetOk("primary_network_interface.0.name"); ok {
				name := nameOk.(string)
				bmsNicPatchModel.Name = &name
			}
		}
		bmsNicPatch, err := bmsNicPatchModel.AsPatch()
		if err != nil {
			return err
		}
		bmsNicUpdateOptions.BareMetalServerNetworkInterfacePatch = bmsNicPatch
		_, _, err = sess.UpdateBareMetalServerNetworkInterfaceWithContext(context, bmsNicUpdateOptions)
		if err != nil {
			return err
		}
		_, err = isWaitForBareMetalServerAvailable(sess, id, d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return err
		}
	}
	if d.HasChange(isBareMetalServerName) {
		flag = true
		nameStr := ""
		if name, ok := d.GetOk(isBareMetalServerName); ok {
			nameStr = name.(string)
		}
		bmsPatchModel.Name = &nameStr
	}
	if flag {
		bmsPatch, err := bmsPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for BareMetalServerPatch: %s", err)
		}
		options.BareMetalServerPatch = bmsPatch
		_, response, err := sess.UpdateBareMetalServerWithContext(context, options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating Bare Metal Server: %s\n%s", err, response)
		}
	}

	if d.HasChange(isBareMetalServerAction) {
		action := ""
		if actionOk, ok := d.GetOk(isBareMetalServerAction); ok {
			action = actionOk.(string)
		}
		if action == "start" {
			isBareMetalServerStart(sess, d.Id(), d, 10)
		} else if action == "stop" {
			isBareMetalServerStop(sess, d.Id(), d, 10)
		} else if action == "restart" {
			isBareMetalServerRestart(sess, d.Id(), d, 10)
		}
	}

	return nil
}

func resourceIBMISBareMetalServerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	deleteType := "hard"
	if dt, ok := d.GetOk(isBareMetalServerDeleteType); ok {
		deleteType = dt.(string)
	}
	err := bareMetalServerDelete(context, d, meta, id, deleteType)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func bareMetalServerDelete(context context.Context, d *schema.ResourceData, meta interface{}, id, deleteType string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getBmsOptions := &vpcv1.GetBareMetalServerOptions{
		ID: &id,
	}
	bms, response, err := sess.GetBareMetalServerWithContext(context, getBmsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Bare Metal Server (%s): %s\n%s", id, err, response)
	}
	if *bms.Status == "running" {

		options := &vpcv1.StopBareMetalServerOptions{
			ID:   bms.ID,
			Type: &deleteType,
		}

		response, err := sess.StopBareMetalServerWithContext(context, options)
		if err != nil && response != nil && response.StatusCode != 204 {
			return fmt.Errorf("[ERROR] Error stopping Bare Metal Server (%s): %s\n%s", id, err, response)
		}
		isWaitForBareMetalServerActionStop(sess, d.Timeout(schema.TimeoutDelete), id, d)

	}
	options := &vpcv1.DeleteBareMetalServerOptions{
		ID: &id,
	}
	response, err = sess.DeleteBareMetalServerWithContext(context, options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Bare Metal Server : %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForBareMetalServerDeleted(bmsC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isBareMetalServerActionDeleting},
		Target:     []string{"done", "", isBareMetalServerActionDeleted, isBareMetalServerStatusFailed},
		Refresh:    isBareMetalServerDeleteRefreshFunc(bmsC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isBareMetalServerDeleteRefreshFunc(bmsC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		bmsgetoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &id,
		}
		bms, response, err := bmsC.GetBareMetalServer(bmsgetoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return bms, isBareMetalServerActionDeleted, nil
			}
			return bms, "", fmt.Errorf("[ERROR] Error Getting Bare Metal Server: %s\n%s", err, response)
		}
		if *bms.Status == isBareMetalServerStatusFailed {
			return bms, *bms.Status, fmt.Errorf("[ERROR] The Bare Metal Server (%s) failed to delete: %v", *bms.ID, err)
		}
		return bms, isBareMetalServerActionDeleting, err
	}
}

func isWaitForBareMetalServerAvailable(client *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for Bare Metal Server (%s) to be available.", id)
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isBareMetalServerStatusPending, isBareMetalServerActionStatusStarting},
		Target:     []string{isBareMetalServerStatusRunning, isBareMetalServerStatusFailed},
		Refresh:    isBareMetalServerRefreshFunc(client, id, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isBareMetalServerRefreshFunc(client *vpcv1.VpcV1, id string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		bmsgetoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &id,
		}
		bms, response, err := client.GetBareMetalServer(bmsgetoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting Bare Metal Server: %s\n%s", err, response)
		}
		d.Set(isBareMetalServerStatus, *bms.Status)

		select {
		case data := <-communicator:
			return nil, "", data.(error)
		default:
			fmt.Println("no message sent")
		}

		if *bms.Status == "running" || *bms.Status == "failed" {
			// let know the isRestartStartAction() to stop
			close(communicator)
			if *bms.Status == "failed" {
				bmsStatusReason := bms.StatusReasons

				//set the status reasons
				if bms.StatusReasons != nil {
					statusReasonsList := make([]map[string]interface{}, 0)
					for _, sr := range bms.StatusReasons {
						currentSR := map[string]interface{}{}
						if sr.Code != nil && sr.Message != nil {
							currentSR[isBareMetalServerStatusReasonsCode] = *sr.Code
							currentSR[isBareMetalServerStatusReasonsMessage] = *sr.Message
							if sr.MoreInfo != nil {
								currentSR[isBareMetalServerStatusReasonsMoreInfo] = *sr.MoreInfo
							}
							statusReasonsList = append(statusReasonsList, currentSR)
						}
					}
					d.Set(isBareMetalServerStatusReasons, statusReasonsList)
				}

				out, err := json.MarshalIndent(bmsStatusReason, "", "    ")
				if err != nil {
					return bms, *bms.Status, fmt.Errorf("[ERROR] The Bare Metal Server (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted bare metal server and attempt to create the bare metal server again replacing the previous configuration", *bms.ID)
				}
				return bms, *bms.Status, fmt.Errorf("[ERROR] Bare Metal Server (%s) went into failed state during the operation \n (%+v) \n [WARNING] Running terraform apply again will remove the tainted Bare Metal Server and attempt to create the Bare Metal Server again replacing the previous configuration", *bms.ID, string(out))
			}
			return bms, *bms.Status, nil

		}
		return bms, isBareMetalServerStatusPending, nil
	}
}

func isWaitForBareMetalServerActionStop(bmsC *vpcv1.VpcV1, timeout time.Duration, id string, d *schema.ResourceData) (interface{}, error) {
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending: []string{isBareMetalServerStatusRunning, isBareMetalServerStatusPending, isBareMetalServerActionStatusStopping},
		Target:  []string{isBareMetalServerActionStatusStopped, isBareMetalServerStatusFailed, ""},
		Refresh: func() (interface{}, string, error) {
			getbmsoptions := &vpcv1.GetBareMetalServerOptions{
				ID: &id,
			}
			bms, response, err := bmsC.GetBareMetalServer(getbmsoptions)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Error Getting Bare Metal Server: %s\n%s", err, response)
			}
			select {
			case data := <-communicator:
				return nil, "", data.(error)
			default:
				fmt.Println("no message sent")
			}
			if *bms.Status == isBareMetalServerStatusFailed {
				// let know the isRestartStopAction() to stop
				close(communicator)
				return bms, *bms.Status, fmt.Errorf("[ERROR] The  Bare Metal Server %s failed to stop: %v", id, err)
			}
			return bms, *bms.Status, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isBareMetalServerRestartStopAction(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int, communicator chan interface{}) {
	subticker := time.NewTicker(time.Duration(forceTimeout) * time.Minute)
	for {
		select {

		case <-subticker.C:
			log.Println("Bare Metal Server is still in stopping state, retrying to stop with -force")
			actiontype := "hard"
			createbmssactoptions := &vpcv1.StopBareMetalServerOptions{
				ID:   &id,
				Type: &actiontype,
			}
			response, err := bmsC.StopBareMetalServer(createbmssactoptions)
			if err != nil {
				communicator <- fmt.Errorf("[ERROR] Error retrying Bare Metal Server action stop: %s\n%s", err, response)
				return
			}
		case <-communicator:
			// indicates refresh func is reached target and not proceed with the thread)
			subticker.Stop()
			return

		}
	}
}

func isBareMetalServerStart(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int) (interface{}, error) {
	createbmsactoptions := &vpcv1.StartBareMetalServerOptions{
		ID: &id,
	}
	response, err := bmsC.StartBareMetalServer(createbmsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("[ERROR] Error creating Bare Metal Server action start : %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerAvailable(bmsC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func isBareMetalServerStop(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int) (interface{}, error) {
	stoppingType := "soft"
	createbmsactoptions := &vpcv1.StopBareMetalServerOptions{
		ID:   &id,
		Type: &stoppingType,
	}
	response, err := bmsC.StopBareMetalServer(createbmsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("[ERROR] Error creating Bare Metal Server Action stop: %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerActionStop(bmsC, d.Timeout(schema.TimeoutUpdate), d.Id(), d)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func isBareMetalServerRestart(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int) (interface{}, error) {
	createbmsactoptions := &vpcv1.RestartBareMetalServerOptions{
		ID: &id,
	}
	response, err := bmsC.RestartBareMetalServer(createbmsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("[ERROR] Error creating Bare Metal Server action restart: %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerAvailable(bmsC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

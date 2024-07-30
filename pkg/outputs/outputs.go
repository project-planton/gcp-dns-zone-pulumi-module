package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpdnszone/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	ManagedZoneGcpProjectId = "gcp-project-id"
	ManagedZoneName         = "managed-zone-name"
	ManagedZoneNameservers  = "nameservers"
)

func PulumiOutputsToStackOutputsConverter(stackOutput auto.OutputMap,
	input *model.GcpDnsZoneStackInput) *model.GcpDnsZoneStackOutputs {
	return &model.GcpDnsZoneStackOutputs{
		Nameservers: autoapistackoutput.GetStringSliceVal(stackOutput, ManagedZoneNameservers),
	}
}

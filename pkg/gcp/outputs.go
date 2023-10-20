package gcp

import (
	"context"
	"github.com/pkg/errors"
	gcpdnszonestack "github.com/plantoncloud-inc/dns-zone-pulumi-blueprint/pkg/gcp/zone"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	dnsv1state "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/dnszone"
	dnszonestack "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/dnszone/stack/gcp"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/enums"
)

func Outputs(ctx context.Context, input *dnszonestack.DnsZoneGcpStackInput) (*dnszonestack.DnsZoneGcpStackOutputs, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return Get(stackOutput, input), nil
}

func Get(stackOutput map[string]interface{}, input *dnszonestack.DnsZoneGcpStackInput) *dnszonestack.DnsZoneGcpStackOutputs {
	if input.StackJob.OperationType != enums.StackOperationType_apply || stackOutput == nil {
		return &dnszonestack.DnsZoneGcpStackOutputs{}
	}
	return &dnszonestack.DnsZoneGcpStackOutputs{
		ZoneStatus: &dnsv1state.DnsZoneStatus{
			Gcp: &dnsv1state.DnsZoneGcpStatus{
				Nameservers: backend.GetStringSliceVal(stackOutput,
					gcpdnszonestack.GetManagedZoneNameserversOutputName(input.ResourceInput.DnsZone.Metadata.Name)),
			},
		},
	}
}

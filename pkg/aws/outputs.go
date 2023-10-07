package aws

import (
	dnszonestack "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/dnszone/stack/aws"
	dnsv1state "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/dnszone/state"
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/stack/rpc/enums"
	"context"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/dns-zone-pulumi-blueprint/pkg/aws/zone"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
)

func Outputs(ctx context.Context, input *dnszonestack.DnsZoneAwsStackInput) (*dnszonestack.DnsZoneAwsStackOutputs, error) {
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

func Get(stackOutput map[string]interface{}, input *dnszonestack.DnsZoneAwsStackInput) *dnszonestack.DnsZoneAwsStackOutputs {
	if input.StackJob.OperationType != enums.StackOperationType_apply || stackOutput == nil {
		return &dnszonestack.DnsZoneAwsStackOutputs{}
	}
	return &dnszonestack.DnsZoneAwsStackOutputs{
		ZoneStatus: &dnsv1state.DnsZoneStatusState{
			Aws: &dnsv1state.DnsZoneAwsStatusState{
				Nameservers: backend.GetStringSliceVal(stackOutput, zone.GetManagedZoneNameserversOutputName(input.ResourceInput.DnsZone.Metadata.Name)),
			},
		},
	}
}

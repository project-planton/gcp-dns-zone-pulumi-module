package aws

import (
	"context"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/job/enums/operationtype"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/dns-zone-pulumi-blueprint/pkg/aws/zone"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	dnsv1state "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/dnszone"
	dnszonestack "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/dnszone/stack/aws"
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
	if input.StackJob.OperationType != operationtype.StackJobOperationType_apply || stackOutput == nil {
		return &dnszonestack.DnsZoneAwsStackOutputs{}
	}
	return &dnszonestack.DnsZoneAwsStackOutputs{
		ZoneStatus: &dnsv1state.DnsZoneStatus{
			Aws: &dnsv1state.DnsZoneAwsStatus{
				Nameservers: backend.GetStringSliceVal(stackOutput, zone.GetManagedZoneNameserversOutputName(input.ResourceInput.DnsZone.Metadata.Name)),
			},
		},
	}
}

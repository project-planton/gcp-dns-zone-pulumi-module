package aws

import (
	"context"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/dns-zone-pulumi-blueprint/pkg/aws/zone"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	code2cloudv1deploydnsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/dnszone/model"
	c2cv1deploydnsstackawsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/dnszone/stack/aws/model"
)

func Outputs(ctx context.Context, input *c2cv1deploydnsstackawsmodel.DnsZoneAwsStackInput) (*c2cv1deploydnsstackawsmodel.DnsZoneAwsStackOutputs, error) {
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

func Get(stackOutput map[string]interface{}, input *c2cv1deploydnsstackawsmodel.DnsZoneAwsStackInput) *c2cv1deploydnsstackawsmodel.DnsZoneAwsStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &c2cv1deploydnsstackawsmodel.DnsZoneAwsStackOutputs{}
	}
	return &c2cv1deploydnsstackawsmodel.DnsZoneAwsStackOutputs{
		ZoneStatus: &code2cloudv1deploydnsmodel.DnsZoneStatus{
			Aws: &code2cloudv1deploydnsmodel.DnsZoneAwsStatus{
				Nameservers: backend.GetStringSliceVal(stackOutput, zone.GetManagedZoneNameserversOutputName(input.ResourceInput.DnsZone.Metadata.Name)),
			},
		},
	}
}

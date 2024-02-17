package gcp

import (
	"context"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/job/enums/operationtype"

	"github.com/pkg/errors"
	gcpdnszonestack "github.com/plantoncloud-inc/dns-zone-pulumi-blueprint/pkg/gcp/zone"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	code2cloudv1deploydnsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/dnszone/model"
	c2cv1deploydnsstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/dnszone/stack/gcp/model"
)

func Outputs(ctx context.Context, input *c2cv1deploydnsstackgcpmodel.DnsZoneGcpStackInput) (*c2cv1deploydnsstackgcpmodel.DnsZoneGcpStackOutputs, error) {
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

func Get(stackOutput map[string]interface{}, input *c2cv1deploydnsstackgcpmodel.DnsZoneGcpStackInput) *c2cv1deploydnsstackgcpmodel.DnsZoneGcpStackOutputs {
	if input.StackJob.Spec.OperationType != operationtype.StackJobOperationType_apply || stackOutput == nil {
		return &c2cv1deploydnsstackgcpmodel.DnsZoneGcpStackOutputs{}
	}
	return &c2cv1deploydnsstackgcpmodel.DnsZoneGcpStackOutputs{
		ZoneStatus: &code2cloudv1deploydnsmodel.DnsZoneStatus{
			Gcp: &code2cloudv1deploydnsmodel.DnsZoneGcpStatus{
				Nameservers: backend.GetStringSliceVal(stackOutput,
					gcpdnszonestack.GetManagedZoneNameserversOutputName(input.ResourceInput.DnsZone.Metadata.Name)),
			},
		},
	}
}

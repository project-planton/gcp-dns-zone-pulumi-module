package aws

import (
	pb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/dnszone/stack/aws"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/dns-zone-pulumi-blueprint/pkg/aws/record"
	"github.com/plantoncloud-inc/dns-zone-pulumi-blueprint/pkg/aws/zone"
	pulumiawsnativeprovider "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/automation/provider/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input     *pb.DnsZoneAwsStackInput
	AwsLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	awsNativeProvider, err := pulumiawsnativeprovider.GetNative(ctx,
		s.Input.CredentialsInput.Aws, "us-west-2")
	if err != nil {
		return errors.Wrap(err, "failed to setup aws provider")
	}

	awsClassicProvider, err := pulumiawsnativeprovider.GetClassic(ctx,
		s.Input.CredentialsInput.Aws, "us-west-2")
	if err != nil {
		return errors.Wrap(err, "failed to setup aws provider")
	}

	createdR53zone, err := zone.Resources(ctx, &zone.Input{
		AwsProvider:        awsNativeProvider,
		StackResourceInput: s.Input.ResourceInput,
		Labels:             s.AwsLabels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add zone resources")
	}

	if err := record.Resources(ctx, &record.Input{
		AwsProvider:    awsClassicProvider,
		DnsZone:        s.Input.ResourceInput.DnsZone,
		CreatedR53Zone: createdR53zone,
	}); err != nil {
		return errors.Wrap(err, "failed to add record resources")
	}
	return nil
}

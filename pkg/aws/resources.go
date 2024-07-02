package aws

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/dns-zone-pulumi-blueprint/pkg/aws/record"
	"github.com/plantoncloud/dns-zone-pulumi-blueprint/pkg/aws/zone"
	c2cv1deploydnsstackawsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/dnszone/stack/aws/model"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input     *c2cv1deploydnsstackawsmodel.DnsZoneAwsStackInput
	AwsLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	awsNativeProvider, err := pulumiawsprovider.GetNative(ctx,
		s.Input.CredentialsInput.Aws, "us-west-2")
	if err != nil {
		return errors.Wrap(err, "failed to setup aws provider")
	}

	awsClassicProvider, err := pulumiawsprovider.GetClassic(ctx,
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

package gcp

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/dns-zone-pulumi-blueprint/pkg/gcp/record"
	"github.com/plantoncloud/dns-zone-pulumi-blueprint/pkg/gcp/zone"
	c2cv1deploydnsstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/dnszone/stack/gcp/model"
	pulumigcpprovider "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/automation/provider/google"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input *c2cv1deploydnsstackgcpmodel.DnsZoneGcpStackInput
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	gcpProvider, err := pulumigcpprovider.Get(ctx, s.Input.CredentialsInput.Google)
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	domainZone, err := zone.Resources(ctx, &zone.Input{
		GcpProvider:        gcpProvider,
		StackResourceInput: s.Input.ResourceInput,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add zone resources")
	}
	if err := record.Resources(ctx, s.Input.ResourceInput.DnsZone, domainZone); err != nil {
		return errors.Wrap(err, "failed to add record resources")
	}
	return nil
}

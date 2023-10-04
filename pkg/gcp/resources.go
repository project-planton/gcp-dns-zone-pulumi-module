package gcp

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/dns-zone-pulumi-stack/pkg/gcp/record"
	"github.com/plantoncloud-inc/dns-zone-pulumi-stack/pkg/gcp/zone"
	pulumigcpprovider "github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/automation/provider/google"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, input) error {
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

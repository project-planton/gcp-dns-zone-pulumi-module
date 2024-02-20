package zone

import (
	"github.com/pkg/errors"
	commonsdnszone "github.com/plantoncloud-inc/go-commons/network/dns/zone"
	puluminameoutputaws "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/aws/output"
	c2cv1deploydnsstackawsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/dnszone/stack/aws/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	pulumiawsnative "github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/route53"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/dns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AwsProvider        *pulumiawsnative.Provider
	StackResourceInput *c2cv1deploydnsstackawsmodel.DnsZoneAwsStackResourceInput
	Labels             map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) (*route53.HostedZone, error) {
	addedManagedZone, err := addZone(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add domain")
	}
	return addedManagedZone, nil
}

func addZone(ctx *pulumi.Context, input *Input) (*route53.HostedZone, error) {
	zoneName := commonsdnszone.GetZoneName(input.StackResourceInput.DnsZone.Metadata.Name)
	z, err := route53.NewHostedZone(ctx, zoneName, &route53.HostedZoneArgs{
		Name: pulumi.String(input.StackResourceInput.DnsZone.Metadata.Name),
		//HostedZoneTags: convertLabelsToTags(input.Labels),
	}, pulumi.Provider(input.AwsProvider))

	if err != nil {
		return nil, errors.Wrapf(err, "failed to add zone for %s domain", input.StackResourceInput.DnsZone.Metadata.Name)
	}

	ctx.Export(GetManagedZoneNameOutputName(input.StackResourceInput.DnsZone.Metadata.Name), z.Name)
	ctx.Export(GetManagedZoneNameserversOutputName(input.StackResourceInput.DnsZone.Metadata.Name), z.NameServers)
	return z, nil
}

func GetManagedZoneNameOutputName(domainName string) string {
	return puluminameoutputaws.Name(dns.ManagedZone{}, commonsdnszone.GetZoneName(domainName), englishword.EnglishWord_name.String())
}

func GetManagedZoneNameserversOutputName(domainName string) string {
	return puluminameoutputaws.Name(dns.ManagedZone{}, commonsdnszone.GetZoneName(domainName), englishword.EnglishWord_nameservers.String())
}

func convertLabelsToTags(labels map[string]string) route53.HostedZoneTagArray {
	hostedZoneTagsArray := make(route53.HostedZoneTagArray, 0, len(labels))
	for k, v := range labels {
		hostedZoneTagsArray = append(hostedZoneTagsArray, route53.HostedZoneTagArgs{
			Key:   pulumi.String(k),
			Value: pulumi.String(v),
		})
	}
	return hostedZoneTagsArray
}

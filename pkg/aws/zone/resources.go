package zone

import (
	pb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/dnszone/stack/aws"
	wordpb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/english/rpc/enums"
	"github.com/pkg/errors"
	commonsdnszone "github.com/plantoncloud-inc/go-commons/network/dns/zone"
	puluminameoutputaws "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/aws/output"
	pulumiawsnative "github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/route53"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/dns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AwsProvider        *pulumiawsnative.Provider
	StackResourceInput *pb.DnsZoneAwsStackResourceInput
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
	return puluminameoutputaws.Name(dns.ManagedZone{}, commonsdnszone.GetZoneName(domainName), wordpb.Word_name.String())
}

func GetManagedZoneNameserversOutputName(domainName string) string {
	return puluminameoutputaws.Name(dns.ManagedZone{}, commonsdnszone.GetZoneName(domainName), wordpb.Word_nameservers.String())
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

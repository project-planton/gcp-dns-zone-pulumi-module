package record

import (
	"strings"

	"github.com/pkg/errors"
	puluminameawsoutput "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/aws/output"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/common/resource/network/dns/record"
	code2cloudv1deploydnsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/dnszone/model"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/route53"
	awsclassic "github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	awsclassicroute53 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AwsProvider    *awsclassic.Provider
	DnsZone        *code2cloudv1deploydnsmodel.DnsZone
	CreatedR53Zone *route53.HostedZone
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if err := addDnsRecords(ctx, input); err != nil {
		return errors.Wrap(err, "failed to add record resources")
	}
	return nil
}

func addDnsRecords(ctx *pulumi.Context, input *Input) error {
	for _, domainRecord := range input.DnsZone.Spec.Records {
		resName := record.Name(domainRecord.Name, strings.ToLower(domainRecord.RecordType.String()))
		rs, err := awsclassicroute53.NewRecord(ctx, resName, &awsclassicroute53.RecordArgs{
			ZoneId:  input.CreatedR53Zone.ID(),
			Name:    pulumi.String(domainRecord.Name),
			Ttl:     pulumi.IntPtr(int(domainRecord.TtlSeconds)),
			Type:    pulumi.String(domainRecord.RecordType),
			Records: pulumi.ToStringArray(domainRecord.Values),
		}, pulumi.Provider(input.AwsProvider))
		if err != nil {
			return errors.Wrapf(err, "failed to add %s rec", domainRecord)
		}
		ctx.Export(puluminameawsoutput.Name(rs, resName), rs.Records)
	}
	return nil
}

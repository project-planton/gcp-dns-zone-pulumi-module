package record

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/common/resource/network/dns/record"
	puluminamegcpoutput "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	code2cloudv1deploydnsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/dnszone/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/dns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, dnsDomain *code2cloudv1deploydnsmodel.DnsZone, domainZone *dns.ManagedZone) error {
	if err := addDnsRecords(ctx, dnsDomain, domainZone); err != nil {
		return errors.Wrap(err, "failed to add record resources")
	}
	return nil
}

func addDnsRecords(ctx *pulumi.Context, dnsDomain *code2cloudv1deploydnsmodel.DnsZone, domainZone *dns.ManagedZone) error {
	for _, domainRecord := range dnsDomain.Spec.Records {
		resName := record.Name(domainRecord.Name, strings.ToLower(domainRecord.RecordType.String()))
		rs, err := dns.NewRecordSet(ctx, resName, &dns.RecordSetArgs{
			ManagedZone: domainZone.Name,
			Name:        pulumi.String(domainRecord.Name),
			Project:     domainZone.Project,
			Rrdatas:     pulumi.ToStringArray(domainRecord.Values),
			Ttl:         pulumi.IntPtr(int(domainRecord.TtlSeconds)),
			Type:        pulumi.String(domainRecord.RecordType.String()),
		}, pulumi.Parent(domainZone))
		if err != nil {
			return errors.Wrapf(err, "failed to add %s rec", domainRecord)
		}
		ctx.Export(puluminamegcpoutput.Name(rs, resName), rs.Rrdatas)
	}
	return nil
}

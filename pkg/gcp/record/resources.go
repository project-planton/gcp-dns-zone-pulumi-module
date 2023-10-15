package record

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/common/resource/network/dns/record"
	puluminamegcpoutput "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	dnsv1state "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/dnszone/state"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/dns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

func Resources(ctx *pulumi.Context, dnsDomain *dnsv1state.DnsZoneState, domainZone *dns.ManagedZone) error {
	if err := addDnsRecords(ctx, dnsDomain, domainZone); err != nil {
		return errors.Wrap(err, "failed to add record resources")
	}
	return nil
}

func addDnsRecords(ctx *pulumi.Context, dnsDomain *dnsv1state.DnsZoneState, domainZone *dns.ManagedZone) error {
	for _, domainRecord := range dnsDomain.Spec.Records {
		resName := record.Name(domainRecord.Name, strings.ToLower(domainRecord.RecordType))
		rs, err := dns.NewRecordSet(ctx, resName, &dns.RecordSetArgs{
			ManagedZone: domainZone.Name,
			Name:        pulumi.String(domainRecord.Name),
			Project:     domainZone.Project,
			Rrdatas:     pulumi.ToStringArray(domainRecord.Values),
			Ttl:         pulumi.IntPtr(int(domainRecord.TtlSeconds)),
			Type:        pulumi.String(domainRecord.RecordType),
		}, pulumi.Parent(domainZone))
		if err != nil {
			return errors.Wrapf(err, "failed to add %s rec", domainRecord)
		}
		ctx.Export(puluminamegcpoutput.Name(rs, resName), rs.Rrdatas)
	}
	return nil
}

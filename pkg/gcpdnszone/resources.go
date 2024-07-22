package gcpdnszone

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	commonsdnsdomain "github.com/plantoncloud-inc/go-commons/network/dns/domain"
	commonsdnszone "github.com/plantoncloud-inc/go-commons/network/dns/zone"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpdnszone/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/dnsrecord"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/dns"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

type ResourceStack struct {
	Input *model.GcpDnsZoneStackInput
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	gcpProvider, err := pulumigoogleprovider.Get(ctx, s.Input.GcpCredential)
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	gcpDnsZone := s.Input.ApiResource

	zoneName := commonsdnszone.GetZoneName(gcpDnsZone.Metadata.Name)

	newManagedZone, err := dns.NewManagedZone(ctx, zoneName, &dns.ManagedZoneArgs{
		Name:        pulumi.String(zoneName),
		Project:     pulumi.String(gcpDnsZone.Spec.ProjectId),
		Description: pulumi.String(fmt.Sprintf("env zone for %s", gcpDnsZone.Metadata.Name)),
		DnsName:     pulumi.String(commonsdnsdomain.SuffixDot(gcpDnsZone.Metadata.Name)),
		Visibility:  pulumi.String("public"),
	}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add zone for %s domain", gcpDnsZone.Metadata.Name)
	}

	ctx.Export(GetManagedZoneNameOutputName(gcpDnsZone.Metadata.Name), newManagedZone.Name)
	ctx.Export(GetManagedZoneNameserversOutputName(gcpDnsZone.Metadata.Name), newManagedZone.NameServers)
	ctx.Export(GetManagedZoneGcpProjectIdOutputName(gcpDnsZone.Metadata.Name), newManagedZone.Project)

	if err := addIamPolicy(ctx, s.Input.ApiResource, newManagedZone); err != nil {
		return errors.Wrapf(err, "failed to add iam policy for managed zone")
	}

	for _, domainRecord := range gcpDnsZone.Spec.Records {
		resName := dnsrecord.PulumiResourceName(domainRecord.Name, strings.ToLower(domainRecord.RecordType.String()))
		rs, err := dns.NewRecordSet(ctx, resName, &dns.RecordSetArgs{
			ManagedZone: newManagedZone.Name,
			Name:        pulumi.String(domainRecord.Name),
			Project:     newManagedZone.Project,
			Rrdatas:     pulumi.ToStringArray(domainRecord.Values),
			Ttl:         pulumi.IntPtr(int(domainRecord.TtlSeconds)),
			Type:        pulumi.String(domainRecord.RecordType.String()),
		}, pulumi.Parent(newManagedZone))
		if err != nil {
			return errors.Wrapf(err, "failed to add %s rec", domainRecord)
		}
		ctx.Export(pulumigoogleprovider.PulumiOutputName(rs, resName), rs.Rrdatas)
	}
	return nil
}

// addIamPolicy creates iam policy granting gcp service accounts permissions required for managing records in the zone.
func addIamPolicy(ctx *pulumi.Context, gcpDnsZone *model.GcpDnsZone, addedManagedZone *dns.ManagedZone) error {
	//when there are no service-accounts, then there is nothing else to do
	if len(gcpDnsZone.Spec.IamServiceAccounts) == 0 {
		return nil
	}

	zoneName := commonsdnszone.GetZoneName(gcpDnsZone.Metadata.Name)
	// todo: the correct resource to use is https://cloud.google.com/dns/docs/zones/iam-per-resource-zones#gcloud
	// but the resource is not yet available in the gcp provider.
	// as a temporary workaround, granting dns admin role to all the service accounts to the entire project.
	// this method grants much broader permissions which allow the service account to control all the zones in the project.
	_, err := projects.NewIAMBinding(ctx, zoneName, &projects.IAMBindingArgs{
		Members: pulumi.StringArray(getIamBindingMembers(gcpDnsZone.Spec.IamServiceAccounts)),
		Project: addedManagedZone.Project,
		Role:    pulumi.String(standard.Dns_admin),
	}, pulumi.Parent(addedManagedZone))
	if err != nil {
		return errors.Wrapf(err, "failed to add project iam binding resource")
	}
	return nil
}

func getIamBindingMembers(iamGcpServiceAccounts []string) []pulumi.StringInput {
	resp := make([]pulumi.StringInput, 0)
	for _, s := range iamGcpServiceAccounts {
		resp = append(resp, pulumi.Sprintf("serviceAccount:%s", s))
	}
	return resp
}

func GetManagedZoneNameOutputName(domainName string) string {
	return pulumigoogleprovider.PulumiOutputName(dns.ManagedZone{},
		commonsdnszone.GetZoneName(domainName), englishword.EnglishWord_name.String())
}

func GetManagedZoneNameserversOutputName(domainName string) string {
	return pulumigoogleprovider.PulumiOutputName(dns.ManagedZone{},
		commonsdnszone.GetZoneName(domainName), englishword.EnglishWord_nameservers.String())
}

func GetManagedZoneGcpProjectIdOutputName(domainName string) string {
	return pulumigoogleprovider.PulumiOutputName(dns.ManagedZone{},
		commonsdnszone.GetZoneName(domainName), englishword.EnglishWord_project.String(),
		englishword.EnglishWord_id.String())
}

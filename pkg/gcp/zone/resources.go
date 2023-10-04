package zone

import (
	pb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/dnszone/stack/gcp"
	wordpb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/english/rpc/enums"
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/network/dns/domain/rpc/enums"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	commonsdnsdomain "github.com/plantoncloud-inc/go-commons/network/dns/domain"
	commonsdnszone "github.com/plantoncloud-inc/go-commons/network/dns/zone"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/name/provider/cloud/gcp/output"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/dns"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	GcpProvider        *pulumigcp.Provider
	StackResourceInput *pb.DnsZoneGcpStackResourceInput
	Labels             map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) (*dns.ManagedZone, error) {
	addedManagedZone, err := addZone(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add domain")
	}
	if err := addIamPolicy(ctx, input.StackResourceInput, addedManagedZone); err != nil {
		return nil, errors.Wrapf(err, "failed to add iam policy for managed zone")
	}
	return addedManagedZone, nil
}

// addIamPolicy creates iam policy granting gcp service accounts permissions required for managing records in the zone.
func addIamPolicy(ctx *pulumi.Context, stackResourceInput *pb.DnsZoneGcpStackResourceInput,
	addedManagedZone *dns.ManagedZone) error {
	zoneName := commonsdnszone.GetZoneName(stackResourceInput.DnsZone.Metadata.Name)
	// todo: the correct resource to use is https://cloud.google.com/dns/docs/zones/iam-per-resource-zones#gcloud
	// but the resource is not yet available in the gcp provider.
	// as a temporary workaround, granting dns admin role to all the service accounts to the entire project.
	// this method grants much broader permissions which allow the service account to control all the zones in the project.
	_, err := projects.NewIAMBinding(ctx, zoneName, &projects.IAMBindingArgs{
		Members: pulumi.StringArray(getIamBindingMembers(stackResourceInput.DnsZone.Spec.Gcp.IamServiceAccounts)),
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

func addZone(ctx *pulumi.Context, input *Input) (*dns.ManagedZone, error) {
	zoneName := commonsdnszone.GetZoneName(input.StackResourceInput.DnsZone.Metadata.Name)
	z, err := dns.NewManagedZone(ctx, zoneName, &dns.ManagedZoneArgs{
		Name:                    pulumi.String(zoneName),
		Project:                 pulumi.String(input.StackResourceInput.DnsZone.Spec.Gcp.ProjectId),
		Description:             pulumi.String(fmt.Sprintf("env zone for %s", input.StackResourceInput.DnsZone.Metadata.Name)),
		DnsName:                 pulumi.String(commonsdnsdomain.SuffixDot(input.StackResourceInput.DnsZone.Metadata.Name)),
		Visibility:              pulumi.String(commonsdnsdomain.GeVisibility(enums.DnsDomainVisibility_DNS_DOMAIN_VISIBILITY_EXTERNAL)),
		PrivateVisibilityConfig: nil,
		Labels:                  pulumi.ToStringMap(input.Labels),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add zone for %s domain", input.StackResourceInput.DnsZone.Metadata.Name)
	}
	ctx.Export(GetManagedZoneNameOutputName(input.StackResourceInput.DnsZone.Metadata.Name), z.Name)
	ctx.Export(GetManagedZoneNameserversOutputName(input.StackResourceInput.DnsZone.Metadata.Name), z.NameServers)
	ctx.Export(GetManagedZoneGcpProjectIdOutputName(input.StackResourceInput.DnsZone.Metadata.Name), z.Project)
	return z, nil
}

// temporarily comment out this bits as internal zones are currently not supported
//func getPrivateVisibilityConfig(inputZoneConfig *pb.IngressDomainStackInput) dns.ManagedZonePrivateVisibilityConfigPtrInput {
//	if inputZoneConfig.Visibility == rpc.DnsDomainVisibility_pub {
//		return nil
//	}
//	return dns.ManagedZonePrivateVisibilityConfigPtrInput(
//		&dns.ManagedZonePrivateVisibilityConfigArgs{
//			Networks: dns.ManagedZonePrivateVisibilityConfigNetworkArray{
//				dns.ManagedZonePrivateVisibilityConfigNetworkArgs{
//					NetworkUrl: pulumi.String(inputZoneConfig.NetworkSelfLink),
//				},
//			},
//		})
//}

func GetManagedZoneNameOutputName(domainName string) string {
	return puluminameoutputgcp.Name(dns.ManagedZone{}, commonsdnszone.GetZoneName(domainName), wordpb.Word_name.String())
}

func GetManagedZoneNameserversOutputName(domainName string) string {
	return puluminameoutputgcp.Name(dns.ManagedZone{}, commonsdnszone.GetZoneName(domainName), wordpb.Word_nameservers.String())
}

func GetManagedZoneGcpProjectIdOutputName(domainName string) string {
	return puluminameoutputgcp.Name(dns.ManagedZone{}, commonsdnszone.GetZoneName(domainName), wordpb.Word_project.String(), wordpb.Word_id.String())
}
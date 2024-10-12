package main

import (
	gcpdnszonev1 "buf.build/gen/go/project-planton/apis/protocolbuffers/go/project/planton/provider/gcp/gcpdnszone/v1"
	"github.com/pkg/errors"
	"github.com/project-planton/gcp-dns-zone-pulumi-module/pkg"
	"github.com/project-planton/pulumi-module-golang-commons/pkg/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpdnszonev1.GcpDnsZoneStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return pkg.Resources(ctx, stackInput)
	})
}

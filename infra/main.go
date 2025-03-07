package main

import (
	"github.com/pulumi/pulumi-oci/sdk/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ociConfig := config.New(ctx, "oci")
		tenancyOcid := ociConfig.Require("tenancyOcid")
		compartment, err := identity.NewCompartment(ctx, "horus", &identity.CompartmentArgs{
			CompartmentId: pulumi.String(tenancyOcid),
			Name:          pulumi.String("horus"),
			Description:   pulumi.String("Horus Project"),
		})
		if err != nil {
			return err
		}

		ctx.Export("compartment", compartment.Name)

		return nil
	})
}

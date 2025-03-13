package main

import (
	"context"
	"fmt"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
	"os"
)

/*
TODO: For testing. Delete later.
*/

func main() {
	stackitClient := stackit.New(&options.ClientOptions{
		Region: "eu01",
	})

	err := stackitClient.Create(context.Background(), "058a4fe4-542f-447c-b06b-de4baed88601")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
func createOrStartServer(cmd *cobra.Command, args []string) error {
	options, err := options.FromEnv(false)
	if err != nil {
		return err
	}

	ctx := context.Background()

	stackit := stackit.New(options.ClientOptions)

	req, publicKey, privateKey, err := h.BuildServerOptions(ctx, options)
	if err != nil {
		return err
	}
	if publicKey == nil {
		return errors.New("no public key generated")
	}

	diskSize, err := strconv.Atoi(options.DiskSize)
	if err != nil {
		return errors.Wrap(err, "parse disk size")
	}

	return stackit.Create(ctx, req, diskSize, *publicKey, privateKey)
}

*/

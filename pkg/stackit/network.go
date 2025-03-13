package stackit

import (
	"context"

	"github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

func (s *Stackit) createNetwork(ctx context.Context, projectId string, networkName string) (*iaas.Network, error) {
	createNetworkPayload := iaas.CreateNetworkPayload{
		Name: utils.Ptr(networkName),
		AddressFamily: &iaas.CreateNetworkAddressFamily{
			Ipv4: &iaas.CreateNetworkIPv4Body{
				PrefixLength: utils.Ptr(int64(29)), // 29 is the largest possible prefix length
			},
		},
	}

	network, err := s.client.CreateNetwork(ctx, projectId).CreateNetworkPayload(createNetworkPayload).Execute()
	if err != nil {
		return nil, err
	}

	network, err = wait.CreateNetworkWaitHandler(ctx, s.client, projectId, *network.NetworkId).WaitWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return network, nil
}

func (s *Stackit) deleteNetwork(ctx context.Context, projectId string, networkId string) error {
	return s.client.DeleteNetworkExecute(ctx, projectId, networkId)
}

package seth_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goplugin/seth"
	"github.com/stretchr/testify/require"

	network_debug_contract "github.com/goplugin/seth/contracts/bind/debug"
	link_token "github.com/goplugin/seth/contracts/bind/link"
	network_debug_sub_contract "github.com/goplugin/seth/contracts/bind/sub"
)

// Shows how to deploy a contract with parameterless constructor and bind it to it's Geth wrapper
func TestDeploymentParameterlessConstructorExample(t *testing.T) {
	commonEnvVars(t)

	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initalise seth")
	contractData, err := c.DeployContractFromContractStore(c.NewTXOpts(), "NetworkDebugSubContract.abi")
	require.NoError(t, err, "failed to deploy sub-debug contract")

	contract, err := network_debug_sub_contract.NewNetworkDebugSubContract(contractData.Address, c.Client)
	require.NoError(t, err, "failed to create debug contract instance")

	_, err = c.Decode(contract.TraceOneInt(c.NewTXOpts(), big.NewInt(1)))
	require.NoError(t, err, "failed to decode transaction")
}

// Shows how to deploy a contract with constructor with parameters and bind it to it's Geth wrapper
func TestDeploymentConstructorWithParametersExample(t *testing.T) {
	commonEnvVars(t)

	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initalise seth")
	contractData, err := c.DeployContractFromContractStore(c.NewTXOpts(), "NetworkDebugSubContract.abi", common.Address{})
	require.NoError(t, err, "failed to deploy debug contract")

	contract, err := network_debug_contract.NewNetworkDebugContract(contractData.Address, c.Client)
	require.NoError(t, err, "failed to create debug contract instance")

	_, err = c.Decode(contract.ProcessUintArray(c.NewTXOpts(), []*big.Int{big.NewInt(1)}))
	require.NoError(t, err, "failed to decode transaction")
}

// Shows how to deploy a contract with parameterless constructor that takes ABI and BIN from Geth wrapper
// and bind it to that wrapper
func TestDeploymentFromGethWrapperExample(t *testing.T) {
	commonEnvVars(t)

	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initalise seth")
	abi, err := network_debug_contract.NetworkDebugContractMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")
	contractData, err := c.DeployContract(c.NewTXOpts(), "NetworkDebugSubContract", *abi, common.FromHex(network_debug_contract.NetworkDebugContractBin))
	require.NoError(t, err, "failed to deploy sub-debug contract from wrapper's ABI/BIN")

	contract, err := network_debug_sub_contract.NewNetworkDebugSubContract(contractData.Address, c.Client)
	require.NoError(t, err, "failed to create debug contract instance")

	_, err = c.Decode(contract.TraceOneInt(c.NewTXOpts(), big.NewInt(1)))
	require.NoError(t, err, "failed to decode transaction")
}

func TestDeploymentLinkTokenFromGethWrapperExample(t *testing.T) {
	commonEnvVars(t)

	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initalise seth")
	abi, err := link_token.LinkTokenMetaData.GetAbi()
	c.ContractStore.ABIs["LinkToken.abi"] = *abi
	require.NoError(t, err, "failed to get ABI")
	contractData, err := c.DeployContract(c.NewTXOpts(), "LinkToken", *abi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy link token contract from wrapper's ABI/BIN")

	contract, err := link_token.NewLinkToken(contractData.Address, c.Client)
	require.NoError(t, err, "failed to create debug contract instance")

	_, err = c.Decode(contract.Mint(c.NewTXOpts(), c.Addresses[0], big.NewInt(1)))
	require.Error(t, err, "did not fail to mint tokens")
}

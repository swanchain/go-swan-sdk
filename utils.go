package swan

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	OrchestratorPublicAddressTestnet = "0x29eD49c8E973696D07E7927f748F6E5Eacd5516D"
	OrchestratorPublicAddressMainnet = "0x4B98086A20f3C19530AF32D21F85Bc6399358e20"
)

func contractInfoVerified(contractInfo ContractInfo, signature string, orchestratorPublicAddress string) bool {
	// Convert contract info to JSON string
	messageJSON, err := json.Marshal(contractInfo)
	if err != nil {
		return false
	}

	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(string(messageJSON))) + string(messageJSON))
	hash := crypto.Keccak256Hash(hashedMessage)

	decodedMessage, err := hexutil.Decode(signature)
	if err != nil {
		return false
	}

	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublicKeyECDSA == nil {
		err = fmt.Errorf("could not get a public get from the message signature")
	}
	if err != nil {
		return false
	}

	return orchestratorPublicAddress == crypto.PubkeyToAddress(*sigPublicKeyECDSA).String()
}

package keys

import (
	"net/http"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func listKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys",
		Long: `Return a list of all public keys stored by this key manager
along with their associated name and address.`,
		RunE: runListCmd,
	}
	cmd.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

func runListCmd(cmd *cobra.Command, args []string) error {
	kb, err := NewKeyBaseFromHomeFlag()
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err == nil {
		printInfos(infos)
	}
	return err
}





// used for outputting keys.Info over REST
type KeyOutput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	// TODO add pubkey?
	// Pubkey  string `json:"pubkey"`
}

func QueryKeysRequestHandler(w http.ResponseWriter, r *http.Request) {

	kb, err := NewKeyBaseFromHomeFlag()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	infos, err := kb.List()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}



	// an empty list will be JSONized as null, but we want to keep the empty list
	if len(infos) == 0 {
		w.Write([]byte("[]"))
		return
	}
	keysOutput := make([]KeyOutput, len(infos))
	for i, info := range infos {
		keysOutput[i] = KeyOutput{Name: info.GetName(), Address: info.GetAddress().String()}
	}
	output, err := json.MarshalIndent(keysOutput, "", "  ")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(output)
}

package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	cfg "github.com/line/ostracon/config"
	tmos "github.com/line/ostracon/libs/os"
	tmrand "github.com/line/ostracon/libs/rand"
	"github.com/line/ostracon/p2p"
	"github.com/line/ostracon/privval"
	"github.com/line/ostracon/types"
	tmtime "github.com/line/ostracon/types/time"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Ostracon",
		RunE:  initFiles,
	}

	AddInitFlags(cmd)
	return cmd
}

func AddInitFlags(cmd *cobra.Command) {
	cmd.Flags().String("priv_key_type", config.PrivKeyType,
		"Specify validator's private key type (ed25519 | composite)")
}

func initFiles(cmd *cobra.Command, args []string) error {
	return initFilesWithConfig(config)
}

func initFilesWithConfig(config *cfg.Config) error {
	// private validator
	privValKeyFile := config.PrivValidatorKeyFile()
	privValStateFile := config.PrivValidatorStateFile()
	privKeyType := config.PrivValidatorKeyType()
	var pv *privval.FilePV
	if tmos.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
		logger.Info("Found private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	} else {
		var err error
		pv, err = privval.GenFilePV(privValKeyFile, privValStateFile, privKeyType)
		if err != nil {
			return err
		}
		if pv != nil {
			pv.Save()
		}
		logger.Info("Generated private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	}

	nodeKeyFile := config.NodeKeyFile()
	if tmos.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	// genesis file
	genFile := config.GenesisFile()
	if tmos.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genDoc := types.GenesisDoc{
			ChainID:         fmt.Sprintf("test-chain-%v", tmrand.Str(6)),
			GenesisTime:     tmtime.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
			VoterParams:     types.DefaultVoterParams(),
		}
		pubKey, err := pv.GetPubKey()
		if err != nil {
			return fmt.Errorf("can't get pubkey: %w", err)
		}
		genDoc.Validators = []types.GenesisValidator{{
			Address: pubKey.Address(),
			PubKey:  pubKey,
			Power:   10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}

	return nil
}

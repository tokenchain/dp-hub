package main

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/bech32"
	"github.com/tokenchain/ixo-blockchain/app"
	cli2 "github.com/tokenchain/ixo-blockchain/client/cli"
	rest2 "github.com/tokenchain/ixo-blockchain/x/dapmining/client/rest"
	"github.com/tomasen/realip"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	CoinDistrHanlderRes struct {
		Result string `json:"result" yaml:"result"`
	}
	ClaimReq struct {
		Address  string `json:"address" yaml:"address"`
		Response string `json:"response" yaml:"response"`
	}
)

const (
	FCCLI                  = "dpfaucet"
	flagRecaptchaSecretKey = "recaptcha_secret_key"
	flagAmountDap          = "amount_dap_faucet"
	flagAmountDollar       = "amount_dollar_faucet"
	flagKey                = "key"
	flagPass               = "pass"
	flagNode               = "node"
	flagUri                = "public_url"
	flagRewardGenInterval  = "reward_interval"
)

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, "=", value)
		return value
	} else {
		log.Fatal("Error loading environment variable: ", key)
		return ""
	}
}
func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	rootCmd := &cobra.Command{
		Use:     "dpfaucet",
		Aliases: []string{"dpf"},
		Short:   "dpFaucet Client",
	}
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().String(flagAmountDap, "0", "The amount of dap")
	rootCmd.PersistentFlags().String(flagAmountDollar, "0", "The amount of dollar")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		if err := cli2.InitConfigFaucet(rootCmd); err != nil {
			return err
		}
		if err := viper.BindPFlag(flags.FlagChainID, rootCmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagRecaptchaSecretKey, rootCmd.PersistentFlags().Lookup(flagRecaptchaSecretKey)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagAmountDap, rootCmd.PersistentFlags().Lookup(flagAmountDap)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagAmountDollar, rootCmd.PersistentFlags().Lookup(flagAmountDollar)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagKey, rootCmd.PersistentFlags().Lookup(flagKey)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagPass, rootCmd.PersistentFlags().Lookup(flagPass)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagNode, rootCmd.PersistentFlags().Lookup(flagNode)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagUri, rootCmd.PersistentFlags().Lookup(flagUri)); err != nil {
			return err
		}
		if err := viper.BindPFlag(flagRewardGenInterval, rootCmd.PersistentFlags().Lookup(flagRewardGenInterval)); err != nil {
			return err
		}

		return nil
	}
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		rest2.ServeCommand(cdc, registerRoutes),
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
	)
}

// registerRoutes registers the routes from the different modules for the LCD.
func registerRoutes(rs *rest2.RestServer) {
	rs.Mux.HandleFunc("/claim", CoinDistriHanlder(rs.CliCtx)).Methods("POST")
}

func CoinDistriHanlder(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chain := viper.GetString(flags.FlagChainID)
		//recaptchaSecretKey := viper.GetString(flagRecaptchaSecretKey)
		amountFaucet := viper.GetString(flagAmountDap)
		amountDollar := viper.GetString(flagAmountDollar)
		key := viper.GetString(flagKey)
		pass := viper.GetString(flagPass)
	    //	node := viper.GetString(flagNode)
		//  publicUrl := viper.GetString(flagUri)
	    //	rewardGenInterval := viper.GetString(flagRewardGenInterval)
	
		resp := CoinDistrHanlderRes{
			Result: "success",
		}
		var claim ClaimReq
		// decode JSON response from front end
		decoder := json.NewDecoder(r.Body)
		decoderErr := decoder.Decode(&claim)
		if decoderErr != nil {
			panic(decoderErr)
		}
		// make sure address is bech32
		readableAddress, decodedAddress, decodeErr := bech32.DecodeAndConvert(claim.Address)
		if decodeErr != nil {
			panic(decodeErr)
		}
		// re-encode the address in bech32
		encodedAddress, encodeErr := bech32.ConvertAndEncode(readableAddress, decodedAddress)
		if encodeErr != nil {
			panic(encodeErr)
		}
		// make sure captcha is valid
		clientIP := realip.FromRequest(r)
		captchaResponse := claim.Response
		captchaPassed, captchaErr := recaptcha.Confirm(clientIP, captchaResponse)
		if captchaErr != nil {
			panic(captchaErr)
		}
		// send the coins!
		if captchaPassed {
			sendFaucet := fmt.Sprintf(
				"%s send --to=%v --name=%v --chain-id=%v --amount=%v",
				FCCLI, encodedAddress, key, chain, amountFaucet)
			fmt.Println(time.Now().UTC().Format(time.RFC3339), encodedAddress, "[1]")
			executeCmd(sendFaucet, pass)

			time.Sleep(5 * time.Second)

			sendSteak := fmt.Sprintf(
				"%s send --to=%v --name=%v --chain-id=%v --amount=%v",
				FCCLI, encodedAddress, key, chain, amountDollar)
			fmt.Println(time.Now().UTC().Format(time.RFC3339), encodedAddress, "[2]")
			executeCmd(sendSteak, pass)
		}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func executeCmd(command string, writes ...string) {
	cmd, wc, _ := goExecute(command)
	for _, write := range writes {
		wc.Write([]byte(write + "\n"))
	}
	cmd.Wait()
}

func goExecute(command string) (cmd *exec.Cmd, pipeIn io.WriteCloser, pipeOut io.ReadCloser) {
	cmd = getCmd(command)
	pipeIn, _ = cmd.StdinPipe()
	pipeOut, _ = cmd.StdoutPipe()
	go cmd.Start()
	time.Sleep(time.Second)
	return cmd, pipeIn, pipeOut
}

func getCmd(command string) *exec.Cmd {
	// split command into command and args
	split := strings.Split(command, " ")
	var cmd *exec.Cmd
	if len(split) == 1 {
		cmd = exec.Command(split[0])
	} else {
		cmd = exec.Command(split[0], split[1:]...)
	}
	return cmd
}

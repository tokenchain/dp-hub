DID base DXP transaction self sign and board cast design
=========================================================

## Step to generation ED25519 Base Key


## Source code for the transaction builder

```go

func NewDidTxBuild(ctx context.CLIContext, msg sdk.Msg, ixoDid exported.IxoDid) SignTxPack {
	instance := SignTxPack{
		ctxCli: ctx,
		msg:    msg,
		did:    ixoDid,
	}
	instance.txBldr = auth.NewTxBuilderFromCLI(ctx.Input)
	return instance
}

func (tb SignTxPack) collectMsgs() []sdk.Msg {
	return []sdk.Msg{tb.msg}
}

func (tb SignTxPack) collectSignatures(sig IxoSignature) []IxoSignature {
	return []IxoSignature{sig}
}

// sign the message in here
func (tb SignTxPack) SignMsgForSignature(msg auth.StdSignMsg) IxoSignature {
	privateKey := tb.did.GetPriKeyByte()

	if l := len(privateKey); l != ed25519.PrivateKeySize {
		panic("ed25519: bad private key length: " + strconv.Itoa(l))
	}

	signatureBytes := ed25519.Sign(privateKey[:], msg.Bytes())
	fmt.Println("===> ⚠️ check signed message =====")
	fmt.Println(msg.Bytes())

	//return NewSignature(time.Now(), signatureBytes[:])
	return NewSignature(time.Now(), signatureBytes)
}

func (tb SignTxPack) CollectSignedMessage(standardMsg auth.StdSignMsg) IxoTx {
	signature := tb.SignMsgForSignature(standardMsg)
	//sign message
	signingSignature := tb.collectSignatures(signature)
	//collection of messages
	messages := tb.collectMsgs()
	return NewIxoTx(messages, standardMsg.Fee, signingSignature, standardMsg.Memo)
}

func (tb SignTxPack) printUnsignedStdTx(stdSignMsg auth.StdSignMsg) error {
	if tb.txBldr.SimulateAndExecute() {
		if err := tb.doSimulate(); err != nil {
			return err
		}
	}

	var json []byte
	var err error

	if viper.GetBool(flags.FlagIndentResponse) {
		json, err = tb.ctxCli.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
	} else {
		json, err = tb.ctxCli.Codec.MarshalJSON(stdSignMsg)
	}

	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(tb.ctxCli.Output, "%s\n", json)

	return nil
}
func (tb SignTxPack) doSimulate() error {
	if tb.ctxCli.Simulate {
		txBldr, err := utils.EnrichWithGas(tb.txBldr, tb.ctxCli, tb.collectMsgs())
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}
	return nil
}
func (tb SignTxPack) DebugTxDecode() error {

	txBldr, err := utils.PrepareTxBuilder(tb.txBldr, tb.ctxCli)
	if err != nil {
		return err
	}

	stdSignMsg, err := txBldr.BuildSignMsg([]sdk.Msg{tb.msg})
	if err != nil {
		return err
	}

	if tb.ctxCli.Simulate {
		if err := tb.printUnsignedStdTx(stdSignMsg); err != nil {
			return err
		}
		return nil
	}
	//will print the message  check signed message ==
	signTxMsg := tb.CollectSignedMessage(stdSignMsg)
	bz, err := tb.ctxCli.Codec.MarshalJSON(signTxMsg)
	if err != nil {
		return fmt.Errorf("Could not marshall tx to binary. Error: %s! ", err.Error())
	}

	tx, err := DefaultTxDecoder(tb.ctxCli.Codec)(bz)
	if err != nil {
		return fmt.Errorf("Could not marshall tx to binary. Error: %s! ", err.Error())
	}

	fmt.Println("=============== signed bytes==============")
	fmt.Println(signTxMsg.GetFirstSignature())
	fmt.Println("=============== decoded signed bytes==============")
	fmt.Println(tx)
	return nil
}
func (tb SignTxPack) CompleteAndBroadcastTxCLI() error {

	txBldr, err := utils.PrepareTxBuilder(tb.txBldr, tb.ctxCli)
	if err != nil {
		return err
	}

	stdSignMsg, err := txBldr.BuildSignMsg([]sdk.Msg{tb.msg})
	if err != nil {
		return err
	}

	if tb.ctxCli.Simulate {
		if err := tb.printUnsignedStdTx(stdSignMsg); err != nil {
			return err
		}
		return nil
	}
	/*
		fmt.Println("=============== public key ==============")
		fmt.Println(tb.did.GetPubKey())
		fmt.Println("=============== private key ==============")
		fmt.Println(tb.did.GetPriKeyByte())
	*/
	if !tb.ctxCli.SkipConfirm {

		var json []byte

		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = tb.ctxCli.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
		} else {
			json, err = tb.ctxCli.Codec.MarshalJSON(stdSignMsg)
		}

		if err != nil {
			panic(err)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)

		ok, err := input.GetConfirmation("Confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}
	//will print the message  check signed message ==
	signTxMsg := tb.CollectSignedMessage(stdSignMsg)
	/*
		fmt.Println("=============== pre-tx-signature ==============")
		fmt.Println(signTxMsg.GetFirstSignature())
	*/
	bz, err := tb.ctxCli.Codec.MarshalJSON(signTxMsg)
	if err != nil {
		return fmt.Errorf("Could not marshall tx to binary. Error: %s! ", err.Error())
	}

	res, err := tb.ctxCli.BroadcastTx(bz)
	if err != nil {
		return fmt.Errorf("Could not broadcast tx. Error: %s! ", err.Error())
	}

	fmt.Println(res.String())
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.TxHash)
	return nil
}

```
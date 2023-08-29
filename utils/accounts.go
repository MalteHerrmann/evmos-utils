package utils

// GetAccounts is a method to retrieve the binaries keys from the configured
// keyring backend and stores it in the Binary struct.
func (bin Binary) GetAccounts() error {
	out, err := ExecuteBinaryCmd(&bin, BinaryCmdArgs{
		Subcommand: []string{"keys", "list", "--output=json"},
	})
	if err != nil {
		return err
	}

	accounts, err := ParseAccountsFromOut(out)
	if err != nil {
		return err
	}

	bin.Accounts = accounts

	return nil
}

package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/greenix-studio/cli/internal/auth"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out of your Greenix account",
	Long: `Sign out of your Greenix account.

Revokes the CLI session on the server and removes the stored credentials
from this machine.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLogout()
	},
}

func runLogout() error {
	creds, err := auth.Load()
	if errors.Is(err, auth.ErrNotLoggedIn) {
		color.Yellow("⚠️  You are not signed in.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not read credentials: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	revokeErr := auth.Revoke(ctx, creds.AccessToken)

	if err := auth.Clear(); err != nil {
		return fmt.Errorf("could not remove credentials: %w", err)
	}

	color.Green("✅ Signed out %s", displayIdentity(creds.Email, creds.Name, creds.UserID))
	if revokeErr != nil {
		color.Yellow("⚠️  Local session removed, but the server could not be reached to revoke it: %s", revokeErr)
	}
	return nil
}

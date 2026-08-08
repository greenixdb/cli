package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/greenixdb/cli/internal/auth"
)

var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	Short:   "Show the currently signed-in Greenix account",
	Aliases: []string{"me"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWhoami()
	},
}

func runWhoami() error {
	creds, err := auth.Load()
	if errors.Is(err, auth.ErrNotLoggedIn) {
		color.Yellow("⚠️  Not signed in. Run `greenix login` to sign in.")
		return fmt.Errorf("not logged in")
	}
	if err != nil {
		return fmt.Errorf("could not read credentials: %w", err)
	}

	if creds.Expired() {
		color.Yellow("⚠️  Your session has expired. Run `greenix login` to sign in again.")
		return fmt.Errorf("session expired")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user, err := auth.Me(ctx, creds.AccessToken)
	if errors.Is(err, auth.ErrNotLoggedIn) {
		_ = auth.Clear()
		color.Yellow("⚠️  Your session is no longer valid. Run `greenix login` to sign in again.")
		return fmt.Errorf("session invalid")
	}
	if err != nil {
		color.Green("👤 %s", displayIdentity(creds.Email, creds.Name, creds.UserID))
		color.Yellow("⚠️  Showing cached details — could not reach the Greenix API: %s", err)
		return nil
	}

	color.Green("👤 %s", displayIdentity(user.Email, user.Name, user.ID))
	if user.ID != "" {
		fmt.Printf("   User ID: %s\n", user.ID)
	}
	if user.Plan != "" {
		fmt.Printf("   Plan:    %s\n", user.Plan)
	}
	if !creds.LoggedInAt.IsZero() {
		fmt.Printf("   Since:   %s\n", creds.LoggedInAt.Local().Format("2006-01-02 15:04"))
	}
	if !creds.ExpiresAt.IsZero() {
		fmt.Printf("   Expires: %s\n", creds.ExpiresAt.Local().Format("2006-01-02 15:04"))
	}
	return nil
}

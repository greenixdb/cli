package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/greenixdb/cli/internal/auth"
)

var loginNoBrowser bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in to your Greenix account",
	Long: `Sign in to your Greenix account.

Opens auth.greenixdb.com in your browser where you can continue with Google.
Once you have approved the request you can close the browser tab — the CLI
detects the result automatically and stores the session on this machine.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLogin()
	},
}

func init() {
	loginCmd.Flags().BoolVar(&loginNoBrowser, "no-browser", false, "print the sign-in URL instead of opening a browser")
}

func runLogin() error {
	if existing, err := auth.Load(); err == nil && !existing.Expired() {
		color.Yellow("⚠️  Already signed in as %s", displayIdentity(existing.Email, existing.Name, existing.UserID))
		fmt.Println("   Run `greenix logout` first to sign in with a different account.")
		return nil
	}

	session, err := auth.NewSession(Version)
	if err != nil {
		return fmt.Errorf("could not start a login session: %w", err)
	}

	color.Cyan("🔐 Signing in to Greenix...")
	fmt.Println()

	opened := false
	if !loginNoBrowser {
		if err := auth.OpenBrowser(session.VerificationURL); err == nil {
			opened = true
		}
	}
	if opened {
		fmt.Println("   Opened your browser to continue with Google.")
		fmt.Println("   If it did not open, visit this URL:")
	} else {
		fmt.Println("   Open this URL in your browser to continue:")
	}
	color.Blue("   %s", session.VerificationURL)
	fmt.Println()
	fmt.Println("   Waiting for you to finish signing in (Ctrl+C to cancel)...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	status, err := session.WaitForApproval(ctx, 2*time.Second, 5*time.Minute)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("login cancelled")
		}
		color.Red("❌ Login failed: %s", err)
		return fmt.Errorf("login failed")
	}

	switch status.Status {
	case "authorized":
		// continue below
	case "denied":
		color.Red("❌ Login was denied in the browser.")
		return fmt.Errorf("login denied")
	case "expired":
		color.Red("❌ The login request expired. Run `greenix login` again.")
		return fmt.Errorf("login expired")
	default:
		color.Red("❌ Login failed: %s", firstNonEmpty(status.Error, status.Status))
		return fmt.Errorf("login failed")
	}

	if status.AccessToken == "" {
		color.Red("❌ The auth server did not return a session token.")
		return fmt.Errorf("login failed")
	}

	creds := &auth.Credentials{
		AccessToken: status.AccessToken,
		TokenType:   firstNonEmpty(status.TokenType, "Bearer"),
		UserID:      status.User.ID,
		Email:       status.User.Email,
		Name:        status.User.Name,
		Plan:        status.User.Plan,
		LoggedInAt:  time.Now(),
	}
	if status.ExpiresIn > 0 {
		creds.ExpiresAt = time.Now().Add(time.Duration(status.ExpiresIn) * time.Second)
	}
	if err := auth.Save(creds); err != nil {
		return fmt.Errorf("could not save credentials: %w", err)
	}

	path, _ := auth.CredentialsPath()
	fmt.Println()
	color.Green("✅ Signed in as %s", displayIdentity(creds.Email, creds.Name, creds.UserID))
	fmt.Printf("   Session saved to %s\n", path)
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func displayIdentity(email, name, id string) string {
	switch {
	case email != "" && name != "":
		return fmt.Sprintf("%s <%s>", name, email)
	case email != "":
		return email
	case name != "":
		return name
	case id != "":
		return id
	default:
		return "your Greenix account"
	}
}

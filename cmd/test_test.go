package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func resetFlags() {
	flagTarget = ""
	flagSearch = ""
	flagFull = false
	flagGenTemplate = ""
	flagSwitch = false
}

func newCmd() *cobra.Command {
	resetFlags()
	cmd := &cobra.Command{
		Use:           "umafetch",
		Short:         "Umamusume - Fastfetch theme",
		Version:       "1.0.0",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetVersionTemplate(`v{{.Version}}`)
	cmd.AddCommand(umasCmd, templateCmd, installCmd)
	return cmd
}

func execute(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestRootHelp(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "umafetch") {
		t.Error("help should contain 'umafetch'")
	}
}

func TestRootVersion(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "--version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "v1.0.0") {
		t.Errorf("version should contain 'v1.0.0', got: %s", out)
	}
}

func TestUmasDefault(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Total umas:") {
		t.Errorf("should show total, got: %s", out)
	}
}

func TestUmasFull(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas", "--full")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Outfits") {
		t.Errorf("should show 'Outfits' in full list, got: %s", out)
	}
}

func TestUmasSearch(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas", "--search", "Special Week")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Special Week") {
		t.Errorf("should find Special Week, got: %s", out)
	}
}

func TestUmasSearchCaseInsensitive(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas", "--search", "special week")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Special Week") {
		t.Errorf("should find Special Week case-insensitive, got: %s", out)
	}
}

func TestUmasSearchPartial(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas", "--search", "spec")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Special Week") {
		t.Errorf("should find Special Week with partial match, got: %s", out)
	}
}

func TestUmasSearchNotFound(t *testing.T) {
	cmd := newCmd()
	_, err := execute(cmd, "umas", "--search", "NobodyExists")
	if err == nil {
		t.Fatal("expected error for non-existent uma")
	}
	if !strings.Contains(err.Error(), "no uma found") {
		t.Errorf("error should contain 'no uma found', got: %v", err)
	}
}

func TestUmasSearchFull(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas", "--search", "Special Week", "--full")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "#") {
		t.Errorf("full search should show hex colors, got: %s", out)
	}
}

func TestUmasTarget(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas", "--target", "Special Week, 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Selected: Special Week (outfit=1)") {
		t.Errorf("should show selected outfit, got: %s", out)
	}
}

func TestUmasTargetRandom(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "umas", "--target", "Special Week")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Selected: Special Week (random outfit)") {
		t.Errorf("should show random outfit, got: %s", out)
	}
}

func TestUmasTargetNotFound(t *testing.T) {
	cmd := newCmd()
	_, err := execute(cmd, "umas", "--target", "Nobody, 1")
	if err == nil {
		t.Fatal("expected error for non-existent uma")
	}
	if !strings.Contains(err.Error(), "no uma found") {
		t.Errorf("error should contain 'no uma found', got: %v", err)
	}
}

func TestUmasTargetInvalidOutfit(t *testing.T) {
	cmd := newCmd()
	_, err := execute(cmd, "umas", "--target", "Special Week, 99")
	if err == nil {
		t.Fatal("expected error for invalid outfit")
	}
	if !strings.Contains(err.Error(), "outfit 99 not found") {
		t.Errorf("error should contain 'outfit 99 not found', got: %v", err)
	}
}

func TestUmasTargetInvalidOrder(t *testing.T) {
	cmd := newCmd()
	_, err := execute(cmd, "umas", "--target", "Special Week, abc")
	if err == nil {
		t.Fatal("expected error for invalid order number")
	}
	if !strings.Contains(err.Error(), "invalid order number") {
		t.Errorf("error should contain 'invalid order number', got: %v", err)
	}
}

func TestTemplateSwitch(t *testing.T) {
	cmd := newCmd()
	out, err := execute(cmd, "template", "--switch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Template switched:") {
		t.Errorf("should show template switched, got: %s", out)
	}
}

func TestTemplateGenInvalid(t *testing.T) {
	cmd := newCmd()
	_, err := execute(cmd, "template", "--gen-template", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
	if !strings.Contains(err.Error(), "unknown template") {
		t.Errorf("error should contain 'unknown template', got: %v", err)
	}
}

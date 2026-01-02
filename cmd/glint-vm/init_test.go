package main

import (
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestInitCommand_Bash(t *testing.T) { //nolint:paralleltest // uses t.Setenv via setupTestEnv
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "auto-switch"},
				},
				Action: initCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"glint-vm", "init", "bash"})
	})

	if !strings.Contains(output, "export GLINT_VM_ROOT=") {
		t.Error("Output should contain GLINT_VM_ROOT export")
	}

	if !strings.Contains(output, "export PATH=") {
		t.Error("Output should contain PATH export")
	}

	if !strings.Contains(output, "glint-vm()") {
		t.Error("Output should contain glint-vm wrapper function")
	}

	if strings.Contains(output, "PROMPT_COMMAND") {
		t.Error("Output should not contain PROMPT_COMMAND without --auto-switch")
	}
}

func TestInitCommand_Bash_AutoSwitch(t *testing.T) { //nolint:paralleltest // uses t.Setenv via setupTestEnv
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "auto-switch"},
				},
				Action: initCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"glint-vm", "init", "--auto-switch", "bash"})
	})

	if !strings.Contains(output, "_glint_vm_auto_switch") {
		t.Error("Output should contain _glint_vm_auto_switch function")
	}

	if !strings.Contains(output, "PROMPT_COMMAND") {
		t.Error("Output should contain PROMPT_COMMAND for bash")
	}
}

func TestInitCommand_Zsh(t *testing.T) { //nolint:paralleltest // uses t.Setenv via setupTestEnv
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "auto-switch"},
				},
				Action: initCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"glint-vm", "init", "zsh"})
	})

	if !strings.Contains(output, "export GLINT_VM_ROOT=") {
		t.Error("Output should contain GLINT_VM_ROOT export")
	}

	if !strings.Contains(output, "glint-vm()") {
		t.Error("Output should contain glint-vm wrapper function")
	}
}

func TestInitCommand_Zsh_AutoSwitch(t *testing.T) { //nolint:paralleltest // uses t.Setenv via setupTestEnv
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "auto-switch"},
				},
				Action: initCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"glint-vm", "init", "--auto-switch", "zsh"})
	})

	if !strings.Contains(output, "_glint_vm_auto_switch") {
		t.Error("Output should contain _glint_vm_auto_switch function")
	}

	if !strings.Contains(output, "add-zsh-hook") {
		t.Error("Output should contain add-zsh-hook for zsh")
	}

	if !strings.Contains(output, "chpwd") {
		t.Error("Output should contain chpwd hook for zsh")
	}
}

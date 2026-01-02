package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewIntegrator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		shellName string
		wantType  string
		wantErr   bool
	}{
		{
			name:      "bash shell",
			shellName: bash,
			wantType:  bash,
			wantErr:   false,
		},
		{
			name:      "zsh shell",
			shellName: zsh,
			wantType:  zsh,
			wantErr:   false,
		},
		{
			name:      "unsupported shell",
			shellName: "fish",
			wantType:  "",
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			integrator, err := NewIntegrator(test.shellName)

			if test.wantErr {
				if err == nil {
					t.Errorf("NewIntegrator() expected error but got nil")
				}

				return
			}

			if err != nil {
				t.Errorf("NewIntegrator() unexpected error: %v", err)

				return
			}

			if integrator.Name() != test.wantType {
				t.Errorf("NewIntegrator() got shell type %s, want %s", integrator.Name(), test.wantType)
			}
		})
	}
}

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name     string
		shellEnv string
		want     string
	}{
		{
			name:     "bash from env",
			shellEnv: "/bin/bash",
			want:     bash,
		},
		{
			name:     "zsh from env",
			shellEnv: "/usr/bin/zsh",
			want:     zsh,
		},
		{
			name:     "no shell env",
			shellEnv: "",
			want:     bash, // default fallback
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			originalShell := os.Getenv("SHELL")

			t.Setenv("SHELL", originalShell)
			t.Setenv("SHELL", test.shellEnv)

			got := DetectShell()
			if got != test.want {
				t.Errorf("DetectShell() = %s, want %s", got, test.want)
			}
		})
	}
}

func TestGetGlintVMRoot(t *testing.T) {
	tests := []struct {
		name       string
		envRoot    string
		wantPrefix string
	}{
		{
			name:       "custom root from env",
			envRoot:    "/custom/path",
			wantPrefix: "/custom/path",
		},
		{
			name:       "default root",
			envRoot:    "",
			wantPrefix: "glint-vm", // Should end with glint-vm
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Save original env
			originalRoot := os.Getenv("GLINT_VM_ROOT")

			defer func() {
				if originalRoot != "" {
					t.Setenv("GLINT_VM_ROOT", originalRoot)
				}
			}()

			if test.envRoot != "" {
				t.Setenv("GLINT_VM_ROOT", test.envRoot)
			}

			got := GetGlintVMRoot()
			if !strings.Contains(got, test.wantPrefix) {
				t.Errorf("GetGlintVMRoot() = %s, want to contain %s", got, test.wantPrefix)
			}
		})
	}
}

func TestGetCurrentDir(t *testing.T) {
	t.Parallel()

	root := GetGlintVMRoot()
	expected := filepath.Join(root, "current")
	got := GetCurrentDir()

	if got != expected {
		t.Errorf("GetCurrentDir() = %s, want %s", got, expected)
	}
}

func TestGetCurrentBinaryPath(t *testing.T) {
	t.Parallel()

	root := GetGlintVMRoot()
	expected := filepath.Join(root, "current", "golangci-lint")
	got := GetCurrentBinaryPath()

	if got != expected {
		t.Errorf("GetCurrentBinaryPath() = %s, want %s", got, expected)
	}
}

func TestBashShell_Name(t *testing.T) {
	t.Parallel()

	shell := &BashShell{}
	if shell.Name() != bash {
		t.Errorf("BashShell.Name() = %s, want bash", shell.Name())
	}
}

func TestBashShell_GenerateInit(t *testing.T) {
	t.Parallel()

	shell := &BashShell{}

	tests := []struct {
		name         string
		opts         InitOptions
		wantContains []string
		wantExcludes []string
	}{
		{
			name: "basic init without auto-switch",
			opts: InitOptions{AutoSwitch: false},
			wantContains: []string{
				"export GLINT_VM_ROOT=",
				"export PATH=",
				"current:",
				"glint-vm()",
				"command glint-vm",
			},
			wantExcludes: []string{
				"_glint_vm_auto_switch",
				"PROMPT_COMMAND",
			},
		},
		{
			name: "init with auto-switch",
			opts: InitOptions{AutoSwitch: true},
			wantContains: []string{
				"export GLINT_VM_ROOT=",
				"export PATH=",
				"current:",
				"glint-vm()",
				"command glint-vm",
				"_glint_vm_auto_switch",
				"PROMPT_COMMAND",
				"GLINT_VM_ROOT}/versions/",
				"_GLINT_VM_NOTIFIED_VERSION",
			},
			wantExcludes: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			output := shell.GenerateInit(test.opts)

			for _, want := range test.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("BashShell.GenerateInit() output missing %q", want)
				}
			}

			for _, exclude := range test.wantExcludes {
				if strings.Contains(output, exclude) {
					t.Errorf("BashShell.GenerateInit() output should not contain %q", exclude)
				}
			}
		})
	}
}

func TestBashShell_GenerateUse(t *testing.T) {
	t.Parallel()

	shell := &BashShell{}
	version := "v1.55.2"

	output := shell.GenerateUse(version)

	wantContains := []string{
		"export PATH=",
		"current:",
		"export GLINT_VM_VERSION=\"v1.55.2\"",
		"echo \"Switched to golangci-lint v1.55.2\" >&2",
	}

	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("BashShell.GenerateUse() output missing %q", want)
		}
	}
}

func TestZshShell_Name(t *testing.T) {
	t.Parallel()

	shell := &ZshShell{}
	if shell.Name() != zsh {
		t.Errorf("ZshShell.Name() = %s, want zsh", shell.Name())
	}
}

func TestZshShell_GenerateInit(t *testing.T) {
	t.Parallel()

	shell := &ZshShell{}

	tests := []struct {
		name         string
		opts         InitOptions
		wantContains []string
		wantExcludes []string
	}{
		{
			name: "basic init without auto-switch",
			opts: InitOptions{AutoSwitch: false},
			wantContains: []string{
				"export GLINT_VM_ROOT=",
				"export PATH=",
				"current:",
				"glint-vm()",
				"command glint-vm",
			},
			wantExcludes: []string{
				"_glint_vm_auto_switch",
				"chpwd",
				"add-zsh-hook",
			},
		},
		{
			name: "init with auto-switch",
			opts: InitOptions{AutoSwitch: true},
			wantContains: []string{
				"export GLINT_VM_ROOT=",
				"export PATH=",
				"current:",
				"glint-vm()",
				"command glint-vm",
				"_glint_vm_auto_switch",
				"chpwd",
				"add-zsh-hook",
				"GLINT_VM_ROOT}/versions/",
				"_GLINT_VM_NOTIFIED_VERSION",
			},
			wantExcludes: []string{
				"PROMPT_COMMAND", // bash-specific
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			output := shell.GenerateInit(test.opts)

			for _, want := range test.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("ZshShell.GenerateInit() output missing %q", want)
				}
			}

			for _, exclude := range test.wantExcludes {
				if strings.Contains(output, exclude) {
					t.Errorf("ZshShell.GenerateInit() output should not contain %q", exclude)
				}
			}
		})
	}
}

func TestZshShell_GenerateUse(t *testing.T) {
	t.Parallel()

	shell := &ZshShell{}
	version := "v1.54.0"

	output := shell.GenerateUse(version)

	wantContains := []string{
		"export PATH=",
		"current:",
		"export GLINT_VM_VERSION=\"v1.54.0\"",
		"echo \"Switched to golangci-lint v1.54.0\" >&2",
	}

	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("ZshShell.GenerateUse() output missing %q", want)
		}
	}
}

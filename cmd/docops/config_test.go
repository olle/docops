package main

import "testing"

func TestParseConfigHTMLFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-html"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.html || cfg.markdown || cfg.txt || cfg.pdf || cfg.web || cfg.tui || cfg.debug > 0 {
		t.Errorf("config = %+v, want html only", cfg)
	}
}

func TestParseConfigMDFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-md"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.markdown || cfg.html || cfg.txt || cfg.pdf || cfg.web || cfg.tui || cfg.debug > 0 {
		t.Errorf("config = %+v, want markdown only", cfg)
	}
}

func TestParseConfigTXTFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-txt"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.txt || cfg.html || cfg.markdown || cfg.pdf || cfg.web || cfg.tui || cfg.debug > 0 {
		t.Errorf("config = %+v, want plain text only", cfg)
	}
}

func TestParseConfigPDFFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-pdf"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.pdf || cfg.html || cfg.markdown || cfg.txt || cfg.web || cfg.tui || cfg.debug > 0 {
		t.Errorf("config = %+v, want pdf only", cfg)
	}
}

func TestParseConfigTUIFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-tui"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.tui || cfg.html || cfg.markdown || cfg.txt || cfg.pdf || cfg.web || cfg.debug > 0 {
		t.Errorf("config = %+v, want tui only", cfg)
	}
}

func TestParseConfigDebugDepthFlags(t *testing.T) {
	for _, test := range []struct {
		flag  string
		depth int
	}{
		{flag: "-d", depth: 1},
		{flag: "-dd", depth: 2},
		{flag: "-ddd", depth: 3},
		{flag: "-debug", depth: 1},
		{flag: "--debug", depth: 1},
	} {
		t.Run(test.flag, func(t *testing.T) {
			cfg, err := parseConfig([]string{test.flag})
			if err != nil {
				t.Fatal(err)
			}
			if cfg.debug != test.depth {
				t.Errorf("debug depth = %d, want %d", cfg.debug, test.depth)
			}
			if cfg.html || cfg.markdown || cfg.txt || cfg.pdf || cfg.web || cfg.tui {
				t.Errorf("config = %+v, want debug only", cfg)
			}
		})
	}
}

func TestParseConfigDebugFalse(t *testing.T) {
	for _, flag := range []string{"-d=false", "-dd=false", "-ddd=false", "-debug=false"} {
		cfg, err := parseConfig([]string{flag})
		if err != nil {
			t.Fatal(err)
		}
		if cfg.debug != 0 {
			t.Errorf("%s: debug depth = %d, want 0", flag, cfg.debug)
		}
	}
}

func TestParseConfigRejectsInvalidDebugValue(t *testing.T) {
	for _, flag := range []string{"-d=maybe", "-debug=maybe"} {
		if _, err := parseConfig([]string{flag}); err == nil {
			t.Errorf("expected %s to fail", flag)
		}
	}
}

func TestParseConfigRejectsMultipleDebugFlags(t *testing.T) {
	for _, args := range [][]string{
		{"-d", "-dd"},
		{"-dd", "-ddd"},
		{"-ddd", "-debug"},
	} {
		if _, err := parseConfig(args); err == nil {
			t.Errorf("parseConfig(%q) = no error, want debug conflict", args)
		}
	}
}

func TestParseConfigWebUsesDefaultAddr(t *testing.T) {
	cfg, err := parseConfig([]string{"-web"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.web || cfg.addr != defaultAddr {
		t.Errorf("config = %+v, want web with %s", cfg, defaultAddr)
	}
}

func TestParseConfigWebWithEqualsAddr(t *testing.T) {
	cfg, err := parseConfig([]string{"-web=127.0.0.1:9090"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.web || cfg.addr != "127.0.0.1:9090" {
		t.Errorf("config = %+v, want web with 127.0.0.1:9090", cfg)
	}
}

func TestParseConfigWebFalse(t *testing.T) {
	cfg, err := parseConfig([]string{"-web=false"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.web {
		t.Errorf("config = %+v, want web disabled", cfg)
	}
}

func TestParseConfigRejectsModesTogether(t *testing.T) {
	for _, args := range [][]string{
		{"-html", "-md"},
		{"-html", "-txt"},
		{"-pdf", "-tui"},
		{"-debug", "-html"},
		{"-web", "-tui"},
		{"-d", "-web"},
	} {
		if _, err := parseConfig(args); err == nil {
			t.Errorf("parseConfig(%q) = no error, want mode conflict", args)
		}
	}
}

func TestParseConfigRejectsRemovedFlags(t *testing.T) {
	for _, args := range [][]string{
		{"-p"},
		{"-server"},
		{"-serve"},
		{"-name", "Ada"},
		{"-out", "out.md"},
		{"-addr", defaultAddr},
		{"-markdown"},
	} {
		if _, err := parseConfig(args); err == nil {
			t.Errorf("parseConfig(%q) = no error, want removed flag to fail", args)
		}
	}
}

func TestParseConfigRejectsWebWithSpaceAddr(t *testing.T) {
	if _, err := parseConfig([]string{"-web", "127.0.0.1:9090"}); err == nil {
		t.Error("expected -web with a separate address to fail")
	}
}

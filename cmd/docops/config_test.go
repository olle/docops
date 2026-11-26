package main

import "testing"

func TestParseConfigHTMLFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-html"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.html || cfg.markdown || cfg.txt || cfg.pdf || cfg.web || cfg.tui || cfg.debug {
		t.Errorf("config = %+v, want html only", cfg)
	}
}

func TestParseConfigMDFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-md"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.markdown || cfg.html || cfg.txt || cfg.pdf || cfg.web || cfg.tui || cfg.debug {
		t.Errorf("config = %+v, want markdown only", cfg)
	}
}

func TestParseConfigTXTFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-txt"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.txt || cfg.html || cfg.markdown || cfg.pdf || cfg.web || cfg.tui || cfg.debug {
		t.Errorf("config = %+v, want plain text only", cfg)
	}
}

func TestParseConfigPDFFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-pdf"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.pdf || cfg.html || cfg.markdown || cfg.txt || cfg.web || cfg.tui || cfg.debug {
		t.Errorf("config = %+v, want pdf only", cfg)
	}
}

func TestParseConfigTUIAndDebugFlags(t *testing.T) {
	for _, test := range []struct {
		flag  string
		check func(config) bool
	}{
		{flag: "-tui", check: func(cfg config) bool { return cfg.tui }},
		{flag: "-debug", check: func(cfg config) bool { return cfg.debug }},
		{flag: "--debug", check: func(cfg config) bool { return cfg.debug }},
	} {
		t.Run(test.flag, func(t *testing.T) {
			cfg, err := parseConfig([]string{test.flag})
			if err != nil {
				t.Fatal(err)
			}
			if !test.check(cfg) {
				t.Errorf("config = %+v, want flag %q enabled", cfg, test.flag)
			}
		})
	}
}

func TestParseConfigDebugFalse(t *testing.T) {
	cfg, err := parseConfig([]string{"-debug=false"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.debug {
		t.Errorf("config = %+v, want debug disabled", cfg)
	}
}

func TestParseConfigRejectsInvalidDebugValue(t *testing.T) {
	if _, err := parseConfig([]string{"-debug=maybe"}); err == nil {
		t.Error("expected invalid debug value to fail")
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
		{"-d"},
		{"--d"},
	} {
		if _, err := parseConfig(args); err == nil {
			t.Errorf("parseConfig(%q) = no error, want removed flag to fail", args)
		}
	}
}

func TestParseConfigRejectsServeWithSpaceAddr(t *testing.T) {
	if _, err := parseConfig([]string{"-web", "127.0.0.1:9090"}); err == nil {
		t.Error("expected -web with a separate address to fail")
	}
}

package main

import "testing"

func TestParseConfigHTMLFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-html"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.html || cfg.markdown {
		t.Errorf("config = %+v, want html", cfg)
	}
}

func TestParseConfigMDFlag(t *testing.T) {
	cfg, err := parseConfig([]string{"-md"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.markdown || cfg.html {
		t.Errorf("config = %+v, want markdown", cfg)
	}
}

func TestParseConfigDebugFlags(t *testing.T) {
	for _, name := range []string{"-debug", "--debug", "-d", "--d"} {
		t.Run(name, func(t *testing.T) {
			cfg, err := parseConfig([]string{name})
			if err != nil {
				t.Fatal(err)
			}
			if !cfg.debug || cfg.server || cfg.tui || cfg.html || cfg.markdown {
				t.Errorf("config = %+v, want debug only", cfg)
			}
		})
	}
}

func TestParseConfigDebugFalse(t *testing.T) {
	cfg, err := parseConfig([]string{"-d=false"})
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

func TestParseConfigRejectsDebugWithHTML(t *testing.T) {
	if _, err := parseConfig([]string{"-debug", "-html"}); err == nil {
		t.Error("expected -debug with -html to fail")
	}
}

func TestParseConfigRejectsMarkdownFlag(t *testing.T) {
	if _, err := parseConfig([]string{"-markdown"}); err == nil {
		t.Error("expected removed -markdown flag to fail")
	}
}

func TestParseConfigServerUsesDefaultAddr(t *testing.T) {
	cfg, err := parseConfig([]string{"-server"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.server || cfg.addr != defaultAddr {
		t.Errorf("config = %+v, want server with %s", cfg, defaultAddr)
	}
}

func TestParseConfigServerWithEqualsAddr(t *testing.T) {
	cfg, err := parseConfig([]string{"-server=127.0.0.1:9090"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.server || cfg.addr != "127.0.0.1:9090" {
		t.Errorf("config = %+v, want server with 127.0.0.1:9090", cfg)
	}
}

func TestParseConfigServerWithSpaceAddr(t *testing.T) {
	cfg, err := parseConfig([]string{"-server", "127.0.0.1:9090"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.server || cfg.addr != "127.0.0.1:9090" {
		t.Errorf("config = %+v, want server with 127.0.0.1:9090", cfg)
	}
}

func TestParseConfigServerFalse(t *testing.T) {
	cfg, err := parseConfig([]string{"-server=false"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.server {
		t.Errorf("config = %+v, want server disabled", cfg)
	}
}

func TestParseConfigRejectsHTMLWithMD(t *testing.T) {
	if _, err := parseConfig([]string{"-html", "-md"}); err == nil {
		t.Error("expected -html with -md to fail")
	}
}

func TestParseConfigRejectsHTMLWithServer(t *testing.T) {
	if _, err := parseConfig([]string{"-html", "-server"}); err == nil {
		t.Error("expected -html with -server to fail")
	}
}

func TestParseConfigRejectsNameFlag(t *testing.T) {
	if _, err := parseConfig([]string{"-name", "Ada"}); err == nil {
		t.Error("expected removed -name flag to fail")
	}
}

func TestParseConfigRejectsOutFlag(t *testing.T) {
	if _, err := parseConfig([]string{"-out", "out.md"}); err == nil {
		t.Error("expected removed -out flag to fail")
	}
}

func TestParseConfigRejectsAddrFlag(t *testing.T) {
	if _, err := parseConfig([]string{"-addr", defaultAddr}); err == nil {
		t.Error("expected removed -addr flag to fail")
	}
}

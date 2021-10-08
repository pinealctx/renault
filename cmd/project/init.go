package project

import (
	"fmt"
	"github.com/pinealctx/renault/pkg/paths"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var rootItem = Item{
	Name: "[ProjectName]",
	Children: []Item{
		{
			Name: "api",
			Children: []Item{
				{
					Name: "models",
					Children: []Item{
						{
							Name: "README.md",
							Template: `# Models

这个目录用于存储一些基础且需跨仓库通用的数据结构。`,
						},
					},
				},
				{
					Name: "pb",
					Children: []Item{
						{
							Name: "README.md",
							Template: `# ProtoBuffer

这个目录用于存储关于ProtoBuffer以及GRPC的定义。`,
						},
					},
				},
				{
					Name: "sql",
					Children: []Item{
						{
							Name: "README.md",
							Template: `# SQL

这个目录用于存储关于Mysql或其他数据库SQL脚本。`,
						},
					},
				},
				{
					Name: "swagger",
					Children: []Item{
						{
							Name: "README.md",
							Template: `# Swagger

这个目录用于存储关于Swagger(openapi)相关的接口定义。`,
						},
					},
				},
			},
		},
		{
			Name: "cmd",
			Children: []Item{
				{
					Name: "[ProjectName]",
					Children: []Item{
						{
							Name:     "main.go",
							Template: "package main\n\nfunc main() {\n\n}\n\n",
						},
					},
				},
			},
		},
		{
			Name: "configs",
			Children: []Item{
				{
					Name: "README.md",
					Template: `# Configs

这个目录用于存储一些关于服务相关的配置。`,
				},
			},
		},
		{
			Name: "internal",
			Children: []Item{
				{
					Name: "README.md",
					Template: `# Internal

这个目录用于存储仅本仓库内会使用到而不愿意暴露给外部使用的package。`,
				},
			},
		},
		{
			Name: "pkg",
			Children: []Item{
				{
					Name: "README.md",
					Template: `# Pkg

这个目录用于存储一些可供外部仓库公用的package。`,
				},
			},
		},
		{
			Name: "scripts",
			Children: []Item{
				{
					Name: "README.md",
					Template: `# Scripts

这个目录用于存储一些脚本文件等。`,
				},
				{
					Name: "lint.sh",
					Template: `#!/usr/bin/env bash

go fmt ./...
go vet ./...
go mod tidy
golangci-lint run ./...`,
				},
			},
		},
		{
			Name: "README.md",
			Template: `# [ProjectName]

This is an interesting service.`,
		},
		{
			Name: ".gitignore",
			Template: `
.idea
.vscode
.DS_Store
`,
		},
		{
			Name: ".golangci.yml",
			Template: `run:
  concurrency: 4
  timeout: 10m
  issues-exit-code: 1
  tests: false
  skip-dirs:
    - tools
  skip-dirs-use-default: true
  allow-parallel-runners: false

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true
  uniq-by-line: true
  sort-results: true

linters-settings:
  dogsled:
    max-blank-identifiers: 2
  dupl:
    threshold: 150
  errcheck:
    check-type-assertions: false
    check-blank: false
  exhaustive:
    check-generated: false
    default-signifies-exhaustive: true
  funlen:
    lines: 120
    statements: 60
  gocognit:
    min-complexity: 20
  nestif:
    min-complexity: 4
  goconst:
    min-len: 3
    min-occurrences: 3
  gocritic:
    disabled-checks:
      - commentFormatting
      - ifElseChain
  gocyclo:
    min-complexity: 20
  godot:
    scope: declarations
    capital: false
  godox:
    keywords:
      - NOTE
      - OPTIMIZE
      - HACK
  gofmt:
    simplify: true
  goimports:
    local-prefixes: github.com/org/project
  golint:
    min-confidence: 0.8
  gomnd:
    settings:
      mnd:
        checks:
          - argument
          - case
          - condition
          - return
        ignored-numbers: 0,1
  gomodguard:
    blocked:
      local_replace_directives: false
  govet:
    check-shadowing: true
    fieldalignment: true
    settings:
      printf:
        funcs:
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Infof
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Warnf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Errorf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Fatalf

    enable:
      - atomicalign
    enable-all: false
    disable:
      - shadow
    disable-all: false
  depguard:
    list-type: blacklist
    include-go-root: false
    packages:
      - github.com/sirupsen/logrus
    packages-with-error-message:
      - github.com/sirupsen/logrus: "logging is allowed only by logutils.Log"
  lll:
    line-length: 600
    tab-width: 1
  misspell:
    locale: US
  nakedret:
    max-func-lines: 30
  prealloc:
    simple: true
    range-loops: true
    for-loops: false
  predeclared:
    ignore: ""
    q: false
  nolintlint:
    allow-leading-space: true
    require-explanation: true
    require-specific: true
  rowserrcheck:
    packages:
      - github.com/jmoiron/sqlx
  testpackage:
    skip-regexp: (export|internal)_test\.go
  thelper:
    test:
      first: true
      name: true
      begin: true
    benchmark:
      first: true
      name: true
      begin: true
  unparam:
    check-exported: false
  unused:
    check-exported: false
  whitespace:
    multi-if: false
    multi-func: false
  wsl:
    strict-append: true
    allow-assign-and-call: true
    allow-multiline-assign: true
    allow-cuddle-declarations: false
    allow-trailing-comment: false
    force-case-trailing-whitespace: 0
    force-err-cuddling: false
    allow-separated-leading-comment: false
  gofumpt:
    extra-rules: false
  errorlint:
    errorf: true
  makezero:
    always: false
  forbidigo:
    forbid:
      - fmt.Errorf
      - fmt.Print.*
      - ginkgo\\.F.*

linters:
  disable-all: true
  enable:
    - deadcode
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unused
    - asciicheck
    - bodyclose
    - dupl
    - exportloopref
    - funlen
    - gocognit
    - goconst
    - gocritic
    - gocyclo
    - godox
    - gosec
    - lll
    - nestif
    - prealloc
    - predeclared
    - exportloopref
    - unparam
    - megacheck

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
        - lll
        - funlen
    - linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
        - lll
        - funlen
      source: "^//go:generate "
severity:
  default-severity: error
  case-sensitive: false
  rules:
    - linters:
        - dupl
      severity: info
`,
		},
		{
			Name: "Makefile",
			Template: `lint:
	@bash ./scripts/lint.sh
`,
		},
		{
			Name: "go.mod",
			Template: `module [Module]

go 1.16

`,
		},
	},
}

var initCommand = &cli.Command{
	Name:  "init",
	Usage: "Initialize the project structure.",
	Flags: []cli.Flag{&cli.StringFlag{
		Name:     "name",
		Usage:    "Specify the project module name, eg: github.com/pinealctx/renault.",
		Required: true,
	}},
	Action: initProject,
}

func initProject(c *cli.Context) error {
	var name = c.String("name")
	var list = strings.Split(name, "/")
	var root = list[len(list)-1]
	var pwd, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("InitProject: Getwd error: %+v", err)
	}
	var fullPath = path.Join(pwd, root)
	var exist bool
	if exist, err = paths.Exists(fullPath); err != nil {
		return fmt.Errorf("InitProject: Check path exist error: %+v", err)
	}
	if exist {
		return fmt.Errorf("InitProject: Root path [%s] exist", root)
	}
	if err = rootItem.Apply(pwd, map[string]string{
		"[ProjectName]": root,
		"[Module]":      name,
	}); err != nil {
		return fmt.Errorf("InitProject: Apply error: %+v", err)
	}
	return nil
}

type Item struct {
	Name     string `json:"name,omitempty"`
	Template string `json:"template,omitempty"`
	Children []Item `json:"children,omitempty"`
}

func (i Item) Apply(parent string, args map[string]string) error {
	var currentPath = path.Join(parent, replaceStr(i.Name, args))
	if i.Template != "" {
		return ioutil.WriteFile(currentPath, []byte(replaceStr(i.Template, args)), 0755)
	}
	var err = os.Mkdir(currentPath, 0755)
	if err != nil {
		return err
	}
	for _, c := range i.Children {
		if err = c.Apply(currentPath, args); err != nil {
			return err
		}
	}
	return nil
}

func replaceStr(s string, args map[string]string) string {
	for k, v := range args {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

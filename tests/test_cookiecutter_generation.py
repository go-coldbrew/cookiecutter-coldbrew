# -*- coding: utf-8 -*-
import os
import re

import pytest
from binaryornot.check import is_binary

PATTERN = r"{{(\s?cookiecutter)[.](.*?)}}"
RE_OBJ = re.compile(PATTERN)


def build_files_list(root_dir):
    """Build a list containing absolute paths to the generated files."""
    return [
        os.path.join(dirpath, file_path)
        for dirpath, subdirs, files in os.walk(root_dir)
        for file_path in files
    ]


def check_paths(paths):
    """Assert that no cookiecutter variables remain unrendered."""
    for path in paths:
        if is_binary(path):
            continue
        with open(path, "r", encoding="latin-1") as f:
            for line in f:
                match = RE_OBJ.search(line)
                msg = "cookiecutter variable not replaced in {}"
                assert match is None, msg.format(path)


# ---------------------------------------------------------------------------
# Structure tests
# ---------------------------------------------------------------------------


class TestProjectStructure:
    """Tests that verify the generated project has the correct file layout."""

    def test_project_directory_name(self, bake_project):
        project = bake_project()
        assert project.name == "testservice"
        assert project.is_dir()

    def test_custom_app_name_with_spaces(self, bake_project):
        # Don't set app_name — let the Jinja2 filter derive it from name
        project = bake_project(full_context={
            "source_path": "github.com/testorg",
            "name": "My Cool Service",
            "grpc_package": "com.github.testorg",
            "service_name": "TestSvc",
            "project_short_description": "A test service.",
            "docker_image": "alpine:latest",
            "docker_build_image": "golang",
            "docker_build_image_version": "1.26",
        })
        assert project.name == "my_cool_service"
        assert project.is_dir()

    def test_expected_files_exist(self, bake_project):
        project = bake_project()
        expected_files = [
            "main.go",
            "Makefile",
            "Dockerfile",
            "docker-compose.local.yml",
            "deploy/prometheus.yml",
            "go.mod",
            "README.md",
            "AGENTS.md",
            "CLAUDE.md",
            "local.env.example",
            ".dockerignore",
            ".editorconfig",
            ".gitignore",
            ".golangci.yml",
            ".mockery.yaml",
            "buf.yaml",
            "buf.gen.yaml",
            "config/config.go",
            "service/service.go",
            "service/healthcheck.go",
            "service/service_test.go",
            "service/healthcheck_test.go",
            "version/version.go",
            ".github/workflows/go.yml",
            ".gitlab-ci.yml",
            "third_party/OpenAPI/embed.go",
        ]
        for filepath in expected_files:
            assert (project / filepath).exists(), f"Missing expected file: {filepath}"

    def test_proto_file_exists(self, bake_project):
        project = bake_project()
        proto_files = list((project / "proto").glob("*.proto"))
        assert len(proto_files) == 1
        assert proto_files[0].name == "testservice.proto"

    def test_removed_files_absent(self, bake_project):
        project = bake_project()
        removed = [
            "tools/tools.go",
            "go.sum",
            "local.env",
        ]
        for filepath in removed:
            assert not (project / filepath).exists(), f"File should not exist: {filepath}"

    def test_no_unrendered_variables(self, bake_project):
        project = bake_project()
        paths = build_files_list(str(project))
        assert paths
        check_paths(paths)


# ---------------------------------------------------------------------------
# Content tests — Go files
# ---------------------------------------------------------------------------


class TestGoFileContent:
    """Tests that verify Go source file content is correct."""

    def test_go_mod_module_path(self, bake_project):
        project = bake_project()
        content = (project / "go.mod").read_text()
        assert "module github.com/testorg/testservice" in content

    def test_go_mod_tool_directive(self, bake_project):
        project = bake_project()
        content = (project / "go.mod").read_text()
        assert "tool (" in content
        assert "github.com/bufbuild/buf/cmd/buf" in content
        assert "github.com/golangci/golangci-lint/v2/cmd/golangci-lint" in content
        assert "github.com/vektra/mockery/v2" in content
        assert "golang.org/x/vuln/cmd/govulncheck" in content

    def test_go_mod_coldbrew_dependencies(self, bake_project):
        project = bake_project()
        content = (project / "go.mod").read_text()
        assert "github.com/go-coldbrew/core" in content
        assert "github.com/go-coldbrew/errors" in content
        assert "github.com/go-coldbrew/log" in content

    def test_go_mod_go_version(self, bake_project):
        project = bake_project()
        content = (project / "go.mod").read_text()
        assert "go 1.26" in content

    def test_main_go_imports(self, bake_project):
        project = bake_project()
        content = (project / "main.go").read_text()
        assert "github.com/testorg/testservice/config" in content
        assert "github.com/testorg/testservice/service" in content
        assert "github.com/go-coldbrew/core" in content

    def test_healthcheck_uses_caller_context(self, bake_project):
        project = bake_project()
        content = (project / "service/healthcheck.go").read_text()
        assert "hcServer.Check(ctx," in content
        assert "context.Background()" not in content

    def test_healthcheck_no_typo(self, bake_project):
        project = bake_project()
        content = (project / "service/healthcheck.go").read_text()
        assert "Not Ready to serve traffic" in content
        assert "Not Ready to server traffic" not in content

    def test_service_formatting(self, bake_project):
        project = bake_project()
        content = (project / "service/service.go").read_text()
        assert "func (s *svc) Stop()" in content
        assert "func (s*svc)" not in content

    def test_version_app_name(self, bake_project):
        project = bake_project()
        content = (project / "version/version.go").read_text()
        assert 'AppName = "testservice"' in content

    def test_config_embeds_coldbrew(self, bake_project):
        project = bake_project()
        content = (project / "config/config.go").read_text()
        assert "cbConfig.Config" in content
        assert "envconfig" in content


# ---------------------------------------------------------------------------
# Content tests — Proto
# ---------------------------------------------------------------------------


class TestProtoContent:
    def test_proto_package(self, bake_project):
        project = bake_project()
        content = (project / "proto/testservice.proto").read_text()
        assert "package com.github.testorg;" in content

    def test_proto_go_package(self, bake_project):
        project = bake_project()
        content = (project / "proto/testservice.proto").read_text()
        assert 'go_package = "github.com/testorg/testservice/proto;testservice"' in content

    def test_proto_service_name(self, bake_project):
        project = bake_project()
        content = (project / "proto/testservice.proto").read_text()
        assert "service TestSvc {" in content

    def test_custom_service_name(self, bake_project):
        project = bake_project({"service_name": "OrderService"})
        proto_file = list((project / "proto").glob("*.proto"))[0]
        content = proto_file.read_text()
        assert "service OrderService {" in content


# ---------------------------------------------------------------------------
# Content tests — Docker
# ---------------------------------------------------------------------------


class TestDockerContent:
    def test_dockerfile_uses_copy_not_add(self, bake_project):
        project = bake_project()
        content = (project / "Dockerfile").read_text()
        # COPY should be used for source code, not ADD
        assert "COPY . ." in content
        assert "\nADD " not in content

    def test_build_alpine_cgo_disabled(self, bake_project):
        """CGO_ENABLED=0 is set in Makefile's build-alpine target (called by Dockerfile)."""
        project = bake_project()
        content = (project / "Makefile").read_text()
        # Find the build-alpine target and verify CGO_ENABLED=0
        assert "CGO_ENABLED=0 go build" in content

    def test_dockerfile_ca_certificates(self, bake_project):
        project = bake_project()
        content = (project / "Dockerfile").read_text()
        assert "ca-certificates" in content

    def test_dockerfile_expose_ports(self, bake_project):
        project = bake_project()
        content = (project / "Dockerfile").read_text()
        assert "EXPOSE 9090 9091" in content

    def test_dockerfile_images(self, bake_project):
        project = bake_project()
        content = (project / "Dockerfile").read_text()
        assert "FROM golang:1.26 AS build-stage" in content
        assert "FROM alpine:latest" in content

    def test_dockerignore_excludes_secrets(self, bake_project):
        project = bake_project()
        content = (project / ".dockerignore").read_text()
        lines = {line.strip() for line in content.splitlines() if line.strip()}
        for entry in ["local.env", "vendor", ".github"]:
            assert entry in lines
        # .git must NOT be excluded — needed for version metadata during docker build
        assert ".git" not in lines
        assert ".git/" not in lines


# ---------------------------------------------------------------------------
# Content tests — Makefile
# ---------------------------------------------------------------------------


class TestMakefileContent:
    def test_git_branch_defined(self, bake_project):
        project = bake_project()
        content = (project / "Makefile").read_text()
        assert "GIT_BRANCH=" in content

    def test_go_tool_commands(self, bake_project):
        project = bake_project()
        content = (project / "Makefile").read_text()
        assert "go tool buf generate proto" in content
        assert "go tool golangci-lint run" in content
        assert "go tool mockery" in content
        assert "go tool govulncheck" in content

    def test_local_stack_targets(self, bake_project):
        project = bake_project()
        content = (project / "Makefile").read_text()
        assert "local-stack:" in content
        assert "local-stack-down:" in content
        assert "local-stack-logs:" in content
        assert "local-stack-reset:" in content
        assert "local-psql:" in content
        assert "docker-compose.local.yml" in content

    def test_bench_run_pattern(self, bake_project):
        project = bake_project()
        content = (project / "Makefile").read_text()
        assert "-run=^$$" in content
        assert "-run=^B" not in content

    def test_fmt_target(self, bake_project):
        project = bake_project()
        content = (project / "Makefile").read_text()
        assert "fmt:" in content
        assert "gofmt -w ." in content



# ---------------------------------------------------------------------------
# Content tests — CI
# ---------------------------------------------------------------------------


class TestCIContent:
    def test_github_actions_checkout_v4(self, bake_project):
        project = bake_project()
        content = (project / ".github/workflows/go.yml").read_text()
        assert "actions/checkout@v4" in content
        assert "actions/checkout@v3" not in content

    def test_github_actions_setup_go_v5(self, bake_project):
        project = bake_project()
        content = (project / ".github/workflows/go.yml").read_text()
        assert "actions/setup-go@v5" in content
        assert "actions/setup-go@v4" not in content

    def test_github_actions_cache_enabled(self, bake_project):
        project = bake_project()
        content = (project / ".github/workflows/go.yml").read_text()
        assert "cache: true" in content

    def test_gitlab_ci_no_go_mod_tidy_in_before_script(self, bake_project):
        project = bake_project()
        content = (project / ".gitlab-ci.yml").read_text()
        # before_script should not contain go mod tidy (conflicts with -mod=readonly)
        lines = content.split("\n")
        in_before_script = False
        for line in lines:
            if "before_script:" in line:
                in_before_script = True
            elif in_before_script and not line.startswith("    -") and line.strip():
                in_before_script = False
            if in_before_script:
                assert "go mod tidy" not in line

    def test_gitlab_ci_no_vet_job(self, bake_project):
        project = bake_project()
        content = (project / ".gitlab-ci.yml").read_text()
        assert "vet:" not in content
        assert "go vet" not in content

    def test_gitlab_ci_go_tool_cobertura(self, bake_project):
        project = bake_project()
        content = (project / ".gitlab-ci.yml").read_text()
        assert "go tool gocover-cobertura" in content
        assert "go install github.com/boumenot/gocover-cobertura" not in content


# ---------------------------------------------------------------------------
# Content tests — Config files
# ---------------------------------------------------------------------------


class TestDockerCompose:
    def test_compose_service_name(self, bake_project):
        project = bake_project()
        content = (project / "docker-compose.local.yml").read_text()
        assert "testservice:" in content
        assert "9090:9090" in content
        assert "9091:9091" in content

    def test_compose_profiles(self, bake_project):
        project = bake_project()
        content = (project / "docker-compose.local.yml").read_text()
        assert 'profiles: ["deps"]' in content
        assert 'profiles: ["obs"]' in content

    def test_compose_db_name(self, bake_project):
        project = bake_project()
        content = (project / "docker-compose.local.yml").read_text()
        assert "POSTGRES_DB: testservice_dev" in content

    def test_agents_md_local_stack(self, bake_project):
        project = bake_project()
        content = (project / "AGENTS.md").read_text()
        assert "make local-stack" in content
        assert "PROFILES=" in content
        assert "localhost:9091" in content


class TestConfigFiles:
    def test_gitignore_entries(self, bake_project):
        project = bake_project()
        content = (project / ".gitignore").read_text()
        for entry in ["local.env", "cover.html", "cover.out", "misc"]:
            assert entry in content

    def test_local_env_example_exists(self, bake_project):
        project = bake_project()
        content = (project / "local.env.example").read_text()
        assert "ENVIRONMENT" in content

    def test_editorconfig(self, bake_project):
        project = bake_project()
        content = (project / ".editorconfig").read_text()
        assert "root = true" in content
        assert "[*.go]" in content

    def test_agents_md(self, bake_project):
        project = bake_project()
        content = (project / "AGENTS.md").read_text()
        assert "make build" in content
        assert "make test" in content
        assert "make lint" in content
        assert "govulncheck" in content
        assert "make generate" in content
        assert "Never edit generated files" in content

    def test_claude_md_imports_agents(self, bake_project):
        project = bake_project()
        content = (project / "CLAUDE.md").read_text()
        assert "@AGENTS.md" in content


# ---------------------------------------------------------------------------
# Parameterized tests
# ---------------------------------------------------------------------------


@pytest.mark.parametrize("go_version", ["1.26", "1.25"])
def test_go_version_selection(bake_project, go_version):
    project = bake_project({"docker_build_image_version": go_version})
    go_mod = (project / "go.mod").read_text()
    assert f"go {go_version}" in go_mod

    dockerfile = (project / "Dockerfile").read_text()
    assert f"golang:{go_version}" in dockerfile

    ci = (project / ".github/workflows/go.yml").read_text()
    assert f'go-version: "{go_version}"' in ci

    gitlab = (project / ".gitlab-ci.yml").read_text()
    assert f'"golang:{go_version}"' in gitlab


def test_custom_source_path(bake_project):
    project = bake_project({"source_path": "gitlab.com/mycompany"})
    go_mod = (project / "go.mod").read_text()
    assert "module gitlab.com/mycompany/testservice" in go_mod

    main_go = (project / "main.go").read_text()
    assert "gitlab.com/mycompany/testservice/config" in main_go

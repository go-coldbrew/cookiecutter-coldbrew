# -*- coding: utf-8 -*-

import os
import re
import sh
import json # Added
from pathlib import Path # Added

import pytest
from binaryornot.check import is_binary
from io import open

PATTERN = '{{(\s?cookiecutter)[.](.*?)}}'
RE_OBJ = re.compile(PATTERN)

test_contexts = [
    {
        'name': 'My Test App',
        'source_path': 'github.com/ankurs',
        'project_short_description': 'A default description for My Test App.',
        'docker_image': 'alpine:latest', # Base image for final stage
        'docker_build_image': 'golang',   # Image for build stage
        'docker_build_image_version': '1.24',
        'service_name': 'MainService',
        'grpc_package': 'com.github.ankurs.main',
        # 'app_name' will be derived from 'name'
    },
    {
        'name': 'another app project',
        'source_path': 'gitlab.com/tester',
        'project_short_description': 'Description for another app project.',
        'docker_image': 'alpine:3.18',
        'docker_build_image': 'golang',
        'docker_build_image_version': '1.23',
        'service_name': 'SecondaryService',
        'grpc_package': 'com.gitlab.tester.secondary',
    },
    {
        'name': 'WebService1',
        'source_path': 'mygit.org/user',
        'project_short_description': 'WebService1: A microservice.',
        'docker_image': 'distroless/static:latest',
        'docker_build_image': 'golang',
        'docker_build_image_version': '1.24', # Intentionally using 1.24 again for variety
        'service_name': 'PrimaryEndpoint',
        'grpc_package': 'org.mycompany.webservice1.primary',
    },
]

def build_files_list(root_dir):
    """Build a list containing absolute paths to the generated files."""
    return [
        os.path.join(dirpath, file_path)
        for dirpath, subdirs, files in os.walk(root_dir)
        for file_path in files
    ]

def check_paths(paths):
    """Method to check all paths have correct substitutions,
    used by other tests cases
    """
    # Assert that no match is found in any of the files
    for path in paths:
        if is_binary(path):
            continue
        for line in open(path, 'r', encoding="latin-1"):
            match = RE_OBJ.search(line)
            msg = 'cookiecutter variable not replaced in {}'
            assert match is None, msg.format(path)

@pytest.mark.parametrize("test_case_context", test_contexts)
def test_default_configuration(cookies, test_case_context):
    # --- Test Setup ---
    # app_name is derived from 'name' in cookiecutter.json, simulate that for expected_app_name
    expected_app_name = test_case_context['name'].replace(' ', '_').lower()
    
    # Note: The template's cookiecutter.json uses {{cookiecutter.name}} to derive app_name.
    # We pass the full test_case_context to cookies.bake().

    # --- 1. Initial Project Bake & Directory Structure ---
    result = cookies.bake(extra_context=test_case_context)
    assert result.exit_code == 0, \
        f"Cookiecutter bake failed for context '{test_case_context['name']}': {result.exception}"
    assert result.exception is None, \
        f"Cookiecutter bake raised an exception for context '{test_case_context['name']}': {result.exception}"
    assert result.project.basename == expected_app_name, \
        f"Project directory name is '{result.project.basename}', expected '{expected_app_name}' for context '{test_case_context['name']}'"
    assert result.project.isdir(), \
        f"Project directory '{result.project}' not found or not a directory for context '{test_case_context['name']}'"

    project_path = Path(result.project)

    # --- 2. Check for Leftover Cookiecutter Variables ---
    # This ensures all template variables were correctly rendered across all generated files.
    paths = build_files_list(str(project_path))
    assert paths, f"No files found in generated project at {project_path} for context '{test_case_context['name']}'"
    check_paths(paths) # check_paths has its own assertion messages referencing specific files

    # --- 3. Rendered Variables in Key Project Files ---
    # Verify that context variables are correctly rendered in specific important files.

    # README.md
    readme_path = project_path / "README.md"
    assert readme_path.is_file(), \
        f"README.md not found at '{readme_path}' for context '{test_case_context['name']}'"
    readme_content = readme_path.read_text(encoding="utf-8")
    assert test_case_context['project_short_description'] in readme_content, \
        f"project_short_description missing in README for context '{test_case_context['name']}'"
    assert expected_app_name in readme_content, \
        f"expected_app_name missing in README for context '{test_case_context['name']}'"

    # go.mod
    go_mod_path = project_path / "go.mod"
    assert go_mod_path.is_file(), \
        f"go.mod not found at '{go_mod_path}' for context '{test_case_context['name']}'"
    go_mod_content = go_mod_path.read_text(encoding="utf-8")
    expected_go_mod_module_line = f"module {test_case_context['source_path']}/{expected_app_name}"
    assert expected_go_mod_module_line in go_mod_content, \
        f"Module line '{expected_go_mod_module_line}' incorrect/missing in go.mod for context '{test_case_context['name']}'"
    expected_go_version_line = f"go {test_case_context['docker_build_image_version'].split('-')[0]}"
    assert expected_go_version_line in go_mod_content, \
        f"Go version line starting with '{expected_go_version_line}' incorrect/missing in go.mod for context '{test_case_context['name']}'"

    # {{cookiecutter.app_name}}/main.go (path uses expected_app_name)
    main_go_path = project_path / expected_app_name / "main.go"
    assert main_go_path.is_file(), \
        f"main.go not found at '{main_go_path}' for context '{test_case_context['name']}'"
    main_go_content = main_go_path.read_text(encoding="utf-8")
    assert f"package main" in main_go_content, \
        f"'package main' not found in '{main_go_path}' for context '{test_case_context['name']}'"

    # Makefile
    makefile_path = project_path / "Makefile"
    assert makefile_path.is_file(), \
        f"Makefile not found at '{makefile_path}' for context '{test_case_context['name']}'"
    makefile_content = makefile_path.read_text(encoding="utf-8")
    assert f"APP_NAME := {expected_app_name}" in makefile_content, \
        f"APP_NAME incorrect/missing in Makefile for context '{test_case_context['name']}'"
    assert f"BINARY_NAME={expected_app_name}" in makefile_content, \
        f"BINARY_NAME incorrect/missing in Makefile for context '{test_case_context['name']}'"

    # --- 4. Dockerfile Content ---
    # Verify the Dockerfile uses the correct build image and version.
    dockerfile_path = project_path / "Dockerfile"
    assert dockerfile_path.is_file(), \
        f"Dockerfile not found at '{dockerfile_path}' for context '{test_case_context['name']}'"
    dockerfile_content = dockerfile_path.read_text(encoding="utf-8")
    expected_docker_build_line = f"FROM {test_case_context['docker_build_image']}:{test_case_context['docker_build_image_version']} AS build-stage"
    assert expected_docker_build_line in dockerfile_content, \
        f"Dockerfile build stage line incorrect/missing for context '{test_case_context['name']}'.\n" \
        f"Expected line: '{expected_docker_build_line}'\nActual content:\n{dockerfile_content}"

    # --- 5. Copied Files (Integrity Check for _copy_without_render mechanism) ---
    # These files are listed in cookiecutter.json's `_copy_without_render` setting.
    # This section verifies if they are copied verbatim. If they are (incorrectly) rendered,
    # their content would differ from the source.
    # Note: The `_copy_without_render` paths in `cookiecutter.json` (e.g., "third_party/OpenAPI/*")
    # are relative to the template root. If the actual source files are nested within a templated
    # directory (e.g., `{{cookiecutter.app_name}}/third_party/...`), they might not be matched
    # by these patterns and thus would be rendered. This test helps detect such discrepancies.
    template_dir = Path(cookies.template)
    # This is the literal name of the directory in the template source structure
    template_app_dir_name = "{{cookiecutter.app_name}}" 

    files_to_compare = [
        "third_party/OpenAPI/swagger-ui.css",
        "third_party/OpenAPI/swagger-ui-bundle.js",
    ]

    for file_rel_path_str in files_to_compare:
        file_rel_path = Path(file_rel_path_str)
        
        source_file_path = template_dir / template_app_dir_name / file_rel_path
        # In the generated project, these files are at `project_root/third_party/...`
        # because the `{{cookiecutter.app_name}}` part of the path is the root of the project.
        generated_file_path = project_path / file_rel_path 

        assert source_file_path.is_file(), \
            f"Source file for comparison not found: '{source_file_path}' (context: '{test_case_context['name']}')"
        assert generated_file_path.is_file(), \
            f"Generated file for comparison not found: '{generated_file_path}' (context: '{test_case_context['name']}')"

        source_content = source_file_path.read_bytes()
        generated_content = generated_file_path.read_bytes()

        assert source_content == generated_content, \
            f"Content of '{file_rel_path}' differs between source template and generated project " \
            f"for context '{test_case_context['name']}'. This might indicate it was (incorrectly) rendered " \
            f"instead of copied directly, or an unexpected modification occurred."

    # --- 6. Hook-Generated File Existence ---
    # These files are generated by post-generation hooks, typically involving `make` commands
    # like `make install generate mock`, which run tools like buf, mockery, and statik.
    
    # Proto-generated Go files
    proto_dir = project_path / "proto"
    assert (proto_dir / f"{expected_app_name}.pb.go").is_file(), \
        f"Proto file '{proto_dir / f'{expected_app_name}.pb.go'}' not found for context '{test_case_context['name']}'"
    assert (proto_dir / f"{expected_app_name}_grpc.pb.go").is_file(), \
        f"Proto GRPC file '{proto_dir / f'{expected_app_name}_grpc.pb.go'}' not found for context '{test_case_context['name']}'"
    assert (proto_dir / f"{expected_app_name}.pb.gw.go").is_file(), \
        f"Proto Gateway file '{proto_dir / f'{expected_app_name}.pb.gw.go'}' not found for context '{test_case_context['name']}'"
    assert (proto_dir / f"{expected_app_name}.pb.vt.go").is_file(), \
        f"Proto VT file '{proto_dir / f'{expected_app_name}.pb.vt.go'}' not found for context '{test_case_context['name']}'"

    # OpenAPI/Swagger file
    openapi_dir = project_path / "third_party" / "OpenAPI"
    assert (openapi_dir / f"{expected_app_name}.swagger.json").is_file(), \
        f"Swagger JSON file '{openapi_dir / f'{expected_app_name}.swagger.json'}' not found for context '{test_case_context['name']}'"

    # Statik file (for embedding OpenAPI spec)
    statik_dir = project_path / "statik"
    assert (statik_dir / "statik.go").is_file(), \
        f"Statik file '{statik_dir / 'statik.go'}' not found for context '{test_case_context['name']}'"

    # --- 7. Git Initialization ---
    # Verifies that the `git init` command in the post-generation hook was successful.
    git_dir = project_path / ".git"
    assert git_dir.is_dir(), \
        f".git directory not found in the generated project at '{git_dir}' for context '{test_case_context['name']}'."

    # --- 8. Hook Command Success (Implicit Check) ---
    # Successful execution of `result = cookies.bake(...)` with `result.exit_code == 0`
    # (checked in section 1) implicitly confirms that pre/post-generation hooks, including
    # `make` commands, completed without error.
    # If hooks failed, result.exit_code would be non-zero.

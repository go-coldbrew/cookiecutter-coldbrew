# -*- coding: utf-8 -*-

import os
import re
import sh

import pytest
from binaryornot.check import is_binary
from io import open

PATTERN = '{{(\s?cookiecutter)[.](.*?)}}'
RE_OBJ = re.compile(PATTERN)

@pytest.fixture
def context():
    return {
        'source_path': 'github.com/ankurs',
        'app_name': 'MyTestProject',
        'project_short_description': 'A short description of the project.',
        "docker_image": "lacion/docker-alpine:latest",
        "docker_build_image": "golang:gobuildimage",
}

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

def test_default_configuration(cookies, context):
    result = cookies.bake(extra_context=context)
    assert result.exit_code == 0
    assert result.exception is None
    assert result.project.basename == context['app_name']
    assert result.project.isdir()

    paths = build_files_list(str(result.project))
    assert paths
    check_paths(paths)

def test_disable_features(cookies, feature_context):
    result = cookies.bake(extra_context=feature_context)
    assert result.exit_code == 0
    assert result.exception is None
    assert result.project.basename == feature_context['app_name']
    assert result.project.isdir()

    paths = build_files_list(str(result.project))
    assert paths
    check_paths(paths)

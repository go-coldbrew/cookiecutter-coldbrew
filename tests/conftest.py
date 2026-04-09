# -*- coding: utf-8 -*-
from pathlib import Path

import pytest
from cookiecutter.main import cookiecutter


@pytest.fixture
def template_dir():
    """Path to the cookiecutter template root."""
    return str(Path(__file__).parent.parent.resolve())


@pytest.fixture
def default_context():
    """Default context for baking the template."""
    return {
        "source_path": "github.com/testorg",
        "name": "TestService",
        "app_name": "testservice",
        "grpc_package": "com.github.testorg",
        "service_name": "TestSvc",
        "project_short_description": "A test service.",
        "docker_image": "alpine:latest",
        "docker_build_image": "golang",
        "docker_build_image_version": "1.26",
    }


@pytest.fixture
def bake_project(tmp_path, template_dir, default_context):
    """Factory fixture that bakes a project.

    Returns a function that accepts optional context overrides and
    returns a pathlib.Path to the generated project directory.

    Pass with_hooks=True to enable post-generation hooks (skips proto init
    via COOKIECUTTER_SKIP_PROTO_INIT=1 to keep tests fast).
    """
    def _bake(extra_context=None, full_context=None, with_hooks=False):
        if full_context is not None:
            ctx = full_context
        else:
            ctx = {**default_context}
            if extra_context:
                ctx.update(extra_context)

        import os
        old_env = os.environ.get("COOKIECUTTER_SKIP_PROTO_INIT")
        if with_hooks:
            os.environ["COOKIECUTTER_SKIP_PROTO_INIT"] = "1"

        try:
            project_dir = cookiecutter(
                template_dir,
                output_dir=str(tmp_path),
                no_input=True,
                extra_context=ctx,
                accept_hooks=with_hooks,
            )
        finally:
            if with_hooks:
                if old_env is None:
                    os.environ.pop("COOKIECUTTER_SKIP_PROTO_INIT", None)
                else:
                    os.environ["COOKIECUTTER_SKIP_PROTO_INIT"] = old_env

        return Path(project_dir)

    return _bake

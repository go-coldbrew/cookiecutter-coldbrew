"""
Does the following:

1. Inits git if used
2. Deletes dockerfiles if not going to be used
3. Deletes config utils if not needed
"""
from __future__ import print_function
import os, sys
import shutil
from subprocess import Popen

# Get the root project directory
PROJECT_DIRECTORY = os.path.realpath(os.path.curdir)

def remove_file(filename):
    """
    generic remove file from project dir
    """
    fullpath = os.path.join(PROJECT_DIRECTORY, filename)
    if os.path.exists(fullpath):
        os.remove(fullpath)

def init_git():
    """
    Initialises git on the new project folder
    """
    GIT_COMMANDS = [
        ["go", "fmt"],
        ["git", "init"],
        ["git", "add", "."],
        ["git", "commit", "-a", "-m", "Initial Commit."]
    ]

    for command in GIT_COMMANDS:
        git = Popen(command, cwd=PROJECT_DIRECTORY)
        git.wait()

def init_proto():
    print("Starting proto initialization...")
    print("Step 1/5: Fetching Go modules (this might take a few minutes)...")
    code = Popen(["go", "mod", "download", "all"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        print("Error: Failed to fetch Go modules.")
        sys.exit(code)

    print("Step 2/5: Running 'make install'...")
    code = Popen(["make", "install"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        print("Error: 'make install' failed.")
        sys.exit(code)

    print("Step 3/5: Running 'make generate'...")
    code = Popen(["make", "generate"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        print("Error: 'make generate' failed.")
        sys.exit(code)

    print("Step 4/5: Tidying Go modules...")
    code = Popen(["go", "mod", "tidy"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        print("Error: 'go mod tidy' failed.")
        sys.exit(code)

    print("Step 5/5: Running 'make mock'...")
    code = Popen(["make", "mock"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        print("Error: 'make mock' failed.")
        sys.exit(code)

    print("Proto initialization completed successfully.")


def remove_docker_files():
    """
    Removes files needed for docker if it isn't going to be used
    """
    for filename in ["Dockerfile",]:
        os.remove(os.path.join(
            PROJECT_DIRECTORY, filename
        ))

init_proto()
init_git()

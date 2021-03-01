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
    code = Popen(["make","install"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        sys.exit(code)
    code = Popen(["make","generate"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        sys.exit(code)
    code = Popen(["go","mod", "tidy"], cwd=PROJECT_DIRECTORY).wait()
    if code > 0:
        sys.exit(code)


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

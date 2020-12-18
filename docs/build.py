# This file is part of arduino-lint.

# Copyright 2020 ARDUINO SA (http://www.arduino.cc/)

# This software is released under the GNU General Public License version 3,
# which covers the main part of arduino-lint.
# The terms of this license can be found at:
# https://www.gnu.org/licenses/gpl-3.0.en.html

# You can be released from the requirements of the above licenses by purchasing
# a commercial license. Buying such a license is mandatory if you want to
# modify or otherwise use the software for commercial activities involving the
# Arduino software without disclosing the source code of your own applications.
# To purchase a commercial license, send an email to license@arduino.cc.
import os
import sys
import re
import unittest
import subprocess

import click
from git import Repo

# In order to provide support for multiple arduino-lint releases, Documentation is versioned so that visitors can select
# which version of the documentation website should be displayed. Unfortunately this feature isn't provided by GitHub
# pages or MkDocs, so we had to implement it on top of the generation process.
#
# Before delving into the details of the generation process, here follow some requirements that were established to
# provide versioned documentation:
#
# - A special version of the documentation called `dev` is provided to reflect the status of the arduino-lint on the
#   `main` branch - this includes unreleased features and bugfixes.
# - Docs are versioned after the minor version of an arduino-lint release. For example, arduino-lint `0.99.1` and
#   `0.99.2` will be both covered by documentation version `0.99`.
# - The landing page of the documentation website will automatically redirect visitors to the documentation of the most
#   recently released version of arduino-lint.
#
# To implement the requirements above, the execution of MkDocs is wrapped using a CLI tool called Mike
# (https://github.com/jimporter/mike) that does a few things for us:
#
# - It runs MkDocs targeting subfolders named after the arduino-lint version, e.g. documentation for version `0.10.1`
#   can be found under the folder `0.10`.
# - It injects an HTML control into the documentation website that lets visitors choose which version of the docs to
#   browse from a dropdown list.
# - It provides a redirect to a version we decide when visitors hit the landing page of the documentation website.
# - It pushes generated contents to the `gh-pages` branch.
#
# In order to avoid unwanted changes to the public website hosting the arduino-lint documentation, only Mike is allowed
# to push changes to the `gh-pages` branch, and this only happens from within the CI, in the "Publish documentation"
# workflow: https://github.com/arduino/arduino-lint/blob/master/.github/workflows/publish-docs.yml
#
# The CI is responsible for guessing which version of arduino-lint we're building docs for, so that generated content
# will be stored in the appropriate section of the documentation website. Because this guessing might be fairly complex,
# the logic is implemented in this Python script. The script will determine the version of arduino-lint that was
# modified in the current commit (either `dev` or an official, numbered release) and whether the redirect to the latest
# version that happens on the landing page should be updated or not.


DEV_BRANCHES = ["main"]


class TestScript(unittest.TestCase):
    def test_get_docs_version(self):
        ver, alias = get_docs_version("main", [])
        self.assertEqual(ver, "dev")
        self.assertEqual(alias, "")

        release_names = ["1.4.x", "0.13.x"]
        ver, alias = get_docs_version("0.13.x", release_names)
        self.assertEqual(ver, "0.13")
        self.assertEqual(alias, "")
        ver, alias = get_docs_version("1.4.x", release_names)
        self.assertEqual(ver, "1.4")
        self.assertEqual(alias, "latest")

        ver, alias = get_docs_version("0.1.x", [])
        self.assertIsNone(ver)
        self.assertIsNone(alias)


def get_docs_version(ref_name, release_branches):
    if ref_name in DEV_BRANCHES:
        return "dev", ""

    if ref_name in release_branches:
        # if version is latest, add an alias
        alias = "latest" if ref_name == release_branches[0] else ""
        # strip `.x` suffix from the branch name to get the version: 0.3.x -> 0.3
        return ref_name[:-2], alias

    return None, None


def get_rel_branch_names(blist):
    """Get the names of the release branches, sorted from newest to older.

    Only process remote refs so we're sure to get all of them and clean up the
    name so that we have a list of strings like 0.6.x, 0.7.x, ...
    """
    pattern = re.compile(r"origin/(\d+\.\d+\.x)")
    names = []
    for b in blist:
        res = pattern.search(b.name)
        if res is not None:
            names.append(res.group(1))

    # Since sorting is stable, first sort by major...
    names = sorted(names, key=lambda x: int(x.split(".")[0]), reverse=True)
    # ...then by minor
    return sorted(names, key=lambda x: int(x.split(".")[1]), reverse=True)


@click.command()
@click.option("--test", is_flag=True)
@click.option("--dry", is_flag=True)
@click.option("--remote", default="origin", help="The git remote where to push.")
def main(test, dry, remote):
    # Run tests if requested
    if test:
        unittest.main(argv=[""], exit=False)
        sys.exit(0)

    # Detect repo root folder
    here = os.path.dirname(os.path.realpath(__file__))
    repo_dir = os.path.join(here, "..")

    # Get current repo
    repo = Repo(repo_dir)

    # Get the list of release branch names
    rel_br_names = get_rel_branch_names(repo.refs)

    # Deduce docs version from current branch. Use the 'latest' alias if
    # version is the most recent
    docs_version, alias = get_docs_version(repo.active_branch.name, rel_br_names)
    if docs_version is None:
        print(f"Can't get version from current branch '{repo.active_branch}', skip docs generation")
        return 0

    # Taskfile args aren't regular args so we put everything in one string
    cmd = (f"task docs:publish DOCS_REMOTE={remote} DOCS_VERSION={docs_version} DOCS_ALIAS={alias}",)

    if dry:
        print(cmd)
        return 0

    subprocess.run(cmd, shell=True, check=True, cwd=repo_dir)


# Usage:
#
#     To run the tests:
#         $python build.py test
#
#     To run the script (must be run from within the repo tree):
#         $python build.py
#
if __name__ == "__main__":
    sys.exit(main())

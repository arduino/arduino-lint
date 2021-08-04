# Source:
# https://github.com/arduino/tooling-project-assets/blob/main/workflow-templates/assets/deploy-mkdocs-versioned/siteversion/test/test_all.py

# Copyright 2020 ARDUINO SA (http://www.arduino.cc/)

# This software is released under the GNU General Public License version 3
# The terms of this license can be found at:
# https://www.gnu.org/licenses/gpl-3.0.en.html

# You can be released from the requirements of the above licenses by purchasing
# a commercial license. Buying such a license is mandatory if you want to
# modify or otherwise use the software for commercial activities involving the
# Arduino software without disclosing the source code of your own applications.
# To purchase a commercial license, send an email to license@arduino.cc.
import siteversion


def test_get_docs_version():
    data = siteversion.get_docs_version("main", [])
    assert data["version"] == "dev"
    assert data["alias"] == ""

    release_names = ["1.4.x", "0.13.x"]
    data = siteversion.get_docs_version("0.13.x", release_names)
    assert data["version"] == "0.13"
    assert data["alias"] == ""
    data = siteversion.get_docs_version("1.4.x", release_names)
    assert data["version"] == "1.4"
    assert data["alias"] == "latest"

    data = siteversion.get_docs_version("0.1.x", [])
    assert data["version"] is None
    assert data["alias"] is None

# This file is part of arduino-lint.
#
# Copyright 2020 ARDUINO SA(http: // www.arduino.cc/)
#
# This software is released under the GNU General Public License version 3,
# which covers the main part of arduino-lint.
# The terms of this license can be found at:
# https: // www.gnu.org/licenses/gpl-3.0.en.html
#
# You can be released from the requirements of the above licenses by purchasing
# a commercial license. Buying such a license is mandatory if you want to
# modify or otherwise use the software for commercial activities involving the
# Arduino software without disclosing the source code of your own applications.
# To purchase a commercial license, send an email to license@arduino.cc.
import json
import pathlib
import platform
import typing

import dateutil.parser
import invoke.context
import pytest
import semver

test_data_path = pathlib.Path(__file__).resolve().parent.joinpath("testdata")


def test_defaults(run_command):
    result = run_command(cmd=[], custom_working_dir=test_data_path.joinpath("recursive"))
    assert result.ok


@pytest.mark.parametrize(
    "project_folder, compliance_level",
    [("Strict", "strict"), ("Specification", "specification"), ("Permissive", "permissive"), ("Invalid", None)],
)
def test_compliance(run_command, project_folder, compliance_level):
    project_path = test_data_path.joinpath("compliance", project_folder)
    expected_ok = False
    for compliance_setting in ["strict", "specification", "permissive"]:
        if compliance_setting == compliance_level:
            expected_ok = True

        result = run_command(cmd=["--compliance", compliance_setting, project_path])
        assert result.ok == expected_ok


def test_compliance_invalid(run_command):
    result = run_command(cmd=["--compliance", "foo", test_data_path.joinpath("ValidSketch")])
    assert not result.ok


def test_format(run_command):
    project_path = test_data_path.joinpath("ValidSketch")
    result = run_command(cmd=["--format", "text", project_path])
    assert result.ok
    with pytest.raises(json.JSONDecodeError):
        json.loads(result.stdout)

    result = run_command(cmd=["--format", "json", project_path])
    assert result.ok
    json.loads(result.stdout)

    result = run_command(cmd=["--format", "foo", project_path])
    assert not result.ok


def test_help(run_command):
    result = run_command(cmd=["--help"])
    assert result.ok
    assert "Usage:" in result.stdout


@pytest.mark.parametrize(
    "project_folder, expected_exit_statuses",
    [
        ("Submit", {"submit": 0, "update": 1, "false": 0}),
        ("Update", {"submit": 1, "update": 0, "false": 0}),
        ("False", {"submit": 1, "update": 1, "false": 0}),
        ("Invalid", {"submit": 1, "update": 1, "false": 1}),
    ],
)
def test_library_manager(run_command, project_folder, expected_exit_statuses):
    project_path = test_data_path.joinpath("library-manager", project_folder)
    for library_manager_setting, expected_exit_status in expected_exit_statuses.items():
        result = run_command(cmd=["--library-manager", library_manager_setting, project_path])
        assert result.exited == expected_exit_status


def test_library_manager_invalid(run_command):
    result = run_command(cmd=["--library-manager", "foo", test_data_path.joinpath("ValidSketch")])
    assert not result.ok


@pytest.mark.parametrize(
    "project_folder, expected_exit_statuses",
    [
        ("Sketch", {"sketch": 0, "library": 1, "platform": 1, "package-index": 1, "all": 0}),
        ("Library", {"sketch": 1, "library": 0, "platform": 1, "package-index": 1, "all": 0}),
        ("Platform", {"sketch": 1, "library": 1, "platform": 0, "package-index": 1, "all": 0}),
        ("PackageIndex", {"sketch": 1, "library": 1, "platform": 1, "package-index": 0, "all": 0}),
    ],
)
def test_project_type(run_command, project_folder, expected_exit_statuses):
    project_path = test_data_path.joinpath("project-type", project_folder)
    for project_type, expected_exit_status in expected_exit_statuses.items():
        result = run_command(cmd=["--project-type", project_type, project_path])
        assert result.exited == expected_exit_status


def test_project_type_invalid(run_command):
    result = run_command(cmd=["--project-type", "foo", test_data_path.joinpath("ValidSketch")])
    assert not result.ok


def test_recursive(run_command):
    valid_projects_path = test_data_path.joinpath("recursive")
    result = run_command(cmd=["--recursive", "true", valid_projects_path])
    assert result.ok

    result = run_command(cmd=["--recursive", "false", valid_projects_path])
    assert not result.ok


def test_recursive_invalid(run_command):
    result = run_command(cmd=["--recursive", "foo", test_data_path.joinpath("ValidSketch")])
    assert not result.ok


def test_report_file(run_command, working_dir):
    project_path = test_data_path.joinpath("ValidSketch")
    report_file_name = "report.json"
    result = run_command(cmd=["--report-file", report_file_name, project_path])
    assert result.ok
    with pathlib.Path(working_dir, report_file_name).open() as report_file:
        report = json.load(report_file)

    assert pathlib.PurePath(report["configuration"]["paths"][0]) == project_path
    assert report["configuration"]["projectType"] == "all"
    assert report["configuration"]["recursive"]
    assert pathlib.PurePath(report["projects"][0]["path"]) == project_path
    assert report["projects"][0]["projectType"] == "sketch"
    assert report["projects"][0]["summary"]["pass"]
    assert report["projects"][0]["summary"]["errorCount"] == 0
    assert report["summary"]["pass"]
    assert report["summary"]["errorCount"] == 0


def test_verbose(run_command):
    project_path = test_data_path.joinpath("verbose", "HasWarnings")
    result = run_command(cmd=["--format", "text", project_path])
    assert result.ok
    assert "result: pass" not in result.stdout
    assert "result: fail" in result.stdout

    result = run_command(cmd=["--format", "text", "--verbose", project_path])
    assert result.ok
    assert "result: pass" in result.stdout

    result = run_command(cmd=["--format", "json", project_path])
    assert result.ok
    report = json.loads(result.stdout)
    assert True not in [check.get("result") == "pass" for check in report["projects"][0]["checks"]]
    assert True in [check.get("result") == "fail" for check in report["projects"][0]["checks"]]

    result = run_command(cmd=["--format", "json", "--verbose", project_path])
    assert result.ok
    report = json.loads(result.stdout)
    assert True in [check.get("result") == "pass" for check in report["projects"][0]["checks"]]
    assert True in [check.get("result") == "fail" for check in report["projects"][0]["checks"]]


def test_version(run_command):
    result = run_command(cmd=["--version"])
    assert result.ok
    output_list = result.stdout.strip().split(sep=" ")
    assert semver.VersionInfo.isvalid(version=output_list[0])
    dateutil.parser.isoparse(output_list[1])


@pytest.fixture(scope="function")
def run_command(pytestconfig, working_dir) -> typing.Callable[..., invoke.runners.Result]:
    """Provide a wrapper around invoke's `run` API so that every test will work in the same temporary folder.

    Useful reference:
        http://docs.pyinvoke.org/en/1.4/api/runners.html#invoke.runners.Result
    """

    arduino_lint_path = pathlib.Path(pytestconfig.rootdir).parent / "arduino-lint"

    def _run(
        cmd: list,
        custom_working_dir: typing.Optional[str] = None,
        custom_env: typing.Optional[dict] = None
    ) -> invoke.runners.Result:
        if cmd is None:
            cmd = []
        if not custom_working_dir:
            custom_working_dir = working_dir
        quoted_cmd = []
        for token in cmd:
            quoted_cmd.append(f'"{token}"')
        cli_full_line = '"{}" {}'.format(arduino_lint_path, " ".join(quoted_cmd))
        run_context = invoke.context.Context()
        # It might happen that we need to change directories between drives on Windows,
        # in that case the "/d" flag must be used otherwise directory wouldn't change
        cd_command = "cd"
        if platform.system() == "Windows":
            cd_command += " /d"
        # Context.cd() is not used since it doesn't work correctly on Windows.
        # It escapes spaces in the path using "\ " but it doesn't always work,
        # wrapping the path in quotation marks is the safest approach
        with run_context.prefix(f'{cd_command} "{custom_working_dir}"'):
            return run_context.run(
                command=cli_full_line, echo=False, hide=True, warn=True, env=custom_env, encoding="utf-8"
            )

    return _run


@pytest.fixture(scope="function")
def working_dir(tmpdir_factory) -> str:
    """Create a temporary folder for the test to run in. It will be created before running each test and deleted at the
    end. This way all the tests work in isolation.
    """
    work_dir = tmpdir_factory.mktemp(basename="ArduinoLintTestWork")
    yield str(work_dir)

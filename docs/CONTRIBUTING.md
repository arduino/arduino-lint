# How to contribute

First of all, thanks for contributing! This document provides some basic guidelines for contributing to this repository.

There are several ways you can get involved:

| Type of contribution                              | Contribution method                                     |
| ------------------------------------------------- | ------------------------------------------------------- |
| - Support request<br/>- Question<br/>- Discussion | Post on the [Arduino Forum][forum]                      |
| - Bug report<br/>- Feature request                | Issue report (read the [issue guidelines][issues])      |
| Beta testing                                      | Try out the [nightly build][nightly]                    |
| - Bug fix<br/>- Enhancement                       | Pull Request (read the [pull request guidelines][prs])  |
| Monetary                                          | - [Donate][donate]<br/>- [Buy official products][store] |

[forum]: https://forum.arduino.cc/c/using-arduino/project-guidance/19
[issues]: #issue-reports
[nightly]: installation.md#nightly-builds
[prs]: #pull-requests
[donate]: https://www.arduino.cc/en/donate/
[store]: https://store.arduino.cc

## Issue Reports

Do you need help or have a question about using Arduino Lint? Support requests should be made to the
[Arduino forum](https://forum.arduino.cc/c/using-arduino/project-guidance/19).

High quality bug reports and feature requests are valuable contributions to the Arduino Lint project.

### Before reporting an issue

- Give the [nightly build](installation.md#nightly-builds) a test drive to see if your issue was already resolved.
- Search [existing pull requests and issues](https://github.com/arduino/arduino-lint/issues?q=) to see if it was already
  reported. If you have additional information to provide about an existing issue, please comment there. You can use
  [GitHub's "Reactions" feature](https://github.blog/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) if
  you only want to express support.

### Qualities of an excellent report

- The issue title should be descriptive. Vague titles make it difficult to understand the purpose of the issue, which
  might cause your issue to be overlooked.
- Provide a full set of steps necessary to reproduce the issue. Demonstration code or commands should be complete and
  simplified to the minimum necessary to reproduce the issue.
- Be responsive. We may need you to provide additional information in order to investigate and resolve the issue.
- If you find a solution to your problem, please comment on your issue report with an explanation of how you were able
  to fix it and close the issue.

## Pull Requests

To propose improvements or fix a bug, feel free to submit a PR.

### Pull request checklist

In order to ease code reviews and have your contributions merged faster, here is a list of items you can check before
submitting a PR:

- Create small PRs that are narrowly focused on addressing a single concern.
- Write tests for the code you wrote.
- Maintain [**clean commit history**](https://www.freshconsulting.com/insights/blog/atomic-commits/) and use
  [**meaningful commit messages**](https://cbea.ms/git-commit/). PRs with messy commit history are difficult to review
  and require a lot of work to be merged.
- <a id="breaking"></a> If the PR contains a breaking change, please start the commit message and PR title with the
  string **[breaking]**. Don't forget to describe in the PR description what changes users might need to make in their
  workflow or application due to this PR. A breaking change is a change that forces users to change their command-line
  invocations or parsing of the JSON formatted output when upgrading from an older version of Arduino Lint.
- PR titles indirectly become part of the CHANGELOG so it's crucial to provide a good record of **what** change is being
  made in the title; **why** it was made will go in the PR description, along with
  [a link to a GitHub issue](https://docs.github.com/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue)
  if one exists.
- Open your PR against the `main` branch.
- Your PR must pass all [CI](https://wikipedia.org/wiki/Continuous_integration) tests before we will merge it. You can
  run the CI in your fork by clicking the "Actions" tab and then the "I understand my workflows..." button. If you're
  seeing an error and don't think it's your fault, it may not be! The reviewer will help you if there are test failures
  that seem not related to the change you are making.

### Development prerequisites

To build Arduino Lint from sources you need the following tools to be available in your local environment:

- [Go](https://golang.org/doc/install)
- [Taskfile](https://taskfile.dev/#/installation) to help you run the most common tasks from the command line

If you want to run integration tests or work on documentation, you will also need:

- A working [Python](https://www.python.org/downloads/) environment.
  - The **Python** version in use is defined in the `tool.poetry.dependencies` field of
    [`pyproject.toml`](https://github.com/arduino/arduino-lint/blob/main/pyproject.toml).
- [Poetry](https://python-poetry.org/docs/).
- [**Node.js** / **npm**](https://nodejs.org/en/download/) - Node.js dependencies management tool.
  - The **Node.js** version in use is defined in the `engines.node` field of
    [`package.json`](https://github.com/arduino/arduino-lint/blob/main/package.json).
  - **ⓘ** [**nvm**](https://github.com/nvm-sh/nvm#installing-and-updating) is recommended if you want to manage multiple
    installations of **Node.js** on your system.

### Building the source code

From the project folder root, just run:

```shell
task build
```

The project uses Go modules, so dependencies will be downloaded automatically. At the end of the build, you should find
the `arduino-lint` executable in the same folder.

<a id="running-the-tests"></a>

### Running the checks

There are several checks and test suites in place to ensure the code works as expected and is written in a way that's
consistent across the whole codebase. To avoid pushing changes that will cause the CI system to fail, you can run the
checks locally by running this command:

```
task check
```

#### Go unit tests

To run only the Go unit tests, run:

```
task go:test
```

By default, all tests for all Arduino Lint's Go packages are run. To run unit tests for only one or more specific
packages, you can set the `TARGETS` environment variable, e.g.:

```
TARGETS=./internal/rule task go:test
```

Alternatively, to run only some specific test(s), you can specify a regex to match against the test function name, e.g.:

```
TEST_REGEX='^TestLibraryProperties.*' task go:test
```

Both can be combined as well, typically to run only a specific test:

```
TEST_REGEX='^TestFindProjects$' TARGETS=./internal/project task go:test
```

#### Integration tests

Being a command line interface, Arduino Lint is heavily interactive and it has to stay consistent in accepting the user
input and providing the expected output and proper exit codes.

For these reasons, in addition to regular unit tests the project has a suite of integration tests that actually run
`arduino-lint` in a different process and assess the options are correctly understood and the output is what we expect.

##### Running tests

After the software requirements have been installed, you should be able to run the tests with:

```
task go:test-integration
```

This will automatically install the necessary dependencies, if not already installed, and run the integration tests
automatically.

To run specific tests, you must run `pytest` from the virtual environment created by Poetry.

```
poetry run pytest tests/test_all.py::test_report_file
```

You can avoid writing the `poetry run` prefix each time by creating a new shell inside the virtual environment:

```
poetry shell
pytest test_lib.py
pytest test_lib.py::test_list
```

### Dependency license metadata

Metadata about the license types of all dependencies is cached in the repository. To update this cache, run the
following command from the repository root folder:

```
task general:cache-dep-licenses
```

The necessary **Licensed** tool can be installed by following
[these instructions](https://github.com/github/licensed#as-an-executable).

Unfortunately, **Licensed** does not have Windows support.

An updated cache is also generated whenever the cache is found to be outdated by the by the "Check Go Dependencies" CI
workflow and made available for download via the `dep-licenses-cache`
[workflow artifact](https://docs.github.com/actions/managing-workflow-runs/downloading-workflow-artifacts).

<a id="linting-and-formatting"></a> <a id="configuration-files-formatting"></a> <a id="documentation-formatting"></a>

### Automated corrections

Tools are provided to automatically bring the project into compliance with some required checks. Run them all with this
command:

```
task fix
```

### Working on documentation

Documentation is provided to final users in form of static HTML content generated from a tool called
[MkDocs](https://www.mkdocs.org/) and hosted on [GitHub Pages](https://pages.github.com/):
https://arduino.github.io/arduino-lint/dev/

#### Local development

The documentation consists of static content written over several Markdown files under the `docs` subfolder of the
Arduino Lint repository, as well as the dynamically generated [command line reference](commands/arduino-lint.md).

When working on the documentation, it is useful to be able to see the effect the changes will have on the generated
documentation website. You can build the documentation website and serve it from your personal computer by running the
command:

```
task website:serve
```

The documentation will build. If you don't see any error, open `http://127.0.0.1:8000` in your browser to view the local
copy of the documentation.

#### Documentation publishing

The Arduino Lint git repository has a special branch called `gh-pages` that contains the generated HTML code for the
documentation website. Every time a change is pushed to this special branch, GitHub automatically triggers a deployment
to pull the change and publish a new version of the website. Do not open Pull Requests to push changes to the `gh-pages`
branch; that will be done exclusively from the CI.

For details on the documentation publishing system, see:
https://github.com/arduino/tooling-project-assets/blob/main/workflow-templates/deploy-cobra-mkdocs-versioned-poetry.md

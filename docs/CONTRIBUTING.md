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

[forum]: https://forum.arduino.cc/index.php?board=3.0
[issues]: #issue-reports
[nightly]: installation.md#nightly-builds
[prs]: #pull-requests
[donate]: https://www.arduino.cc/en/Main/Contribute
[store]: https://store.arduino.cc

## Issue Reports

Do you need help or have a question about using Arduino Lint? Support requests should be made to the
[Arduino forum](https://forum.arduino.cc/index.php?board=3.0).

High quality bug reports and feature requests are valuable contributions to the Arduino Lint project.

### Before reporting an issue

- Give the [nightly build](installation.md#nightly-builds) a test drive to see if your issue was already resolved.
- Search [existing pull requests and issues](https://github.com/arduino/arduino-lint/issues?q=) to see if it was already
  reported. If you have additional information to provide about an existing issue, please comment there. You can use
  [GitHub's "Reactions" feature](https://github.com/blog/2119-add-reactions-to-pull-requests-issues-and-comments) if you
  only want to express support.

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
- Maintain [**clean commit history**](http://www.freshconsulting.com/atomic-commits) and use
  [**meaningful commit messages**](http://chris.beams.io/posts/git-commit). PRs with messy commit history are difficult
  to review and require a lot of work to be merged.
- <a id="breaking"></a> If the PR contains a breaking change, please start the commit message and PR title with the
  string **[breaking]**. Don't forget to describe in the PR description what changes users might need to make in their
  workflow or application due to this PR. A breaking change is a change that forces users to change their command-line
  invocations or parsing of the JSON formatted output when upgrading from an older version of Arduino Lint.
- PR titles indirectly become part of the CHANGELOG so it's crucial to provide a good record of **what** change is being
  made in the title; **why** it was made will go in the PR description, along with
  [a link to a GitHub issue](https://docs.github.com/en/free-pro-team@latest/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue)
  if one exists.
- Open your PR against the `main` branch.
- Your PR must pass all [CI](https://en.wikipedia.org/wiki/Continuous_integration) tests before we will merge it. You
  can run the CI in your fork by clicking the "Actions" tab and then the "I understand my workflows..." button. If
  you're seeing an error and don't think it's your fault, it may not be! The reviewer will help you if there are test
  failures that seem not related to the change you are making.

### Development prerequisites

To build Arduino Lint from sources you need the following tools to be available in your local environment:

- [Go](https://golang.org/doc/install) version 1.16 or later
- [Taskfile](https://taskfile.dev/#/installation) to help you run the most common tasks from the command line

If you want to run integration tests or work on documentation, you will also need:

- A working [Python](https://www.python.org/downloads/) environment, version 3.9 or later.
- [Poetry](https://python-poetry.org/docs/).

### Building the source code

From the project folder root, just run: F

```shell
task build
```

The project uses Go modules, so dependencies will be downloaded automatically. At the end of the build, you should find
the `arduino-lint` executable in the same folder.

### Running the tests

There are several checks and test suites in place to ensure the code works as expected and is written in a way that's
consistent across the whole codebase. To avoid pushing changes that will cause the CI system to fail, you can run most
of the tests locally.

To ensure code style is consistent, run:

```
task check
```

To run all tests:

```
task test
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

#### Linting and formatting

When editing any Python file in the project, remember to run linting checks with:

```
task python:check
```

This will run [`flake8`](https://flake8.pycqa.org/) automatically and return any error in the code formatting.

In case of linting errors you should be able to solve most of them by automatically formatting with:

```shell
task python:format
```

#### Configuration files formatting

We use [Prettier](https://prettier.io/) to automatically format all YAML files in the project. Keeping and enforcing a
formatting standard helps everyone make small PRs and avoids the introduction of formatting changes made by unconfigured
editors.

There are several ways to run Prettier. If you're using Visual Studio Code you can easily use the
[`prettier-vscode` extension](https://github.com/prettier/prettier-vscode) to automatically format as you write.

Otherwise you can use the following tasks. To do so you'll need to install `npm` if not already installed. Check the
[official documentation](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm) to learn how to install
`npm` for your platform.

To check if the files are correctly formatted run:

```shell
task config:check
```

If the output tells you that some files are not formatted correctly run:

```shell
task config:format
```

Checks are automatically run on every pull request to verify that configuration files are correctly formatted.

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

The documentation will build. If you don't see any error, open `http://127.0.0.1:8000` in your browser to local the
local copy of the documentation.

#### Documentation formatting

We use [Prettier](https://prettier.io/) to automatically format all Markdown files in the project. Keeping and enforcing
a formatting standard helps everyone make small PRs and avoids the introduction of formatting changes made by
misconfigured editors.

There are several ways to run Prettier. If you're using Visual Studio Code you can easily use the
[`prettier-vscode` extension](https://github.com/prettier/prettier-vscode) to automatically format as you write.

Otherwise you can use the following tasks. To do so you'll need to install `npm` if not already installed. Check the
[official documentation](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm) to learn how to install
`npm` for your platform.

To check if the files are correctly formatted run:

```shell
task website:check
```

If the output tells you that some files are not formatted correctly run:

```shell
task docs:format
```

Checks are automatically run on every pull request to verify that documentation is correctly formatted.

#### Documentation publishing

The Arduino Lint git repository has a special branch called `gh-pages` that contains the generated HTML code for the
documentation website. Every time a change is pushed to this special branch, GitHub automatically triggers a deployment
to pull the change and publish a new version of the website. Do not open Pull Requests to push changes to the `gh-pages`
branch; that will be done exclusively from the CI.

For details on the documentation publishing system, see:
https://github.com/arduino/tooling-project-assets/blob/main/workflow-templates/deploy-cobra-mkdocs-versioned-poetry.md

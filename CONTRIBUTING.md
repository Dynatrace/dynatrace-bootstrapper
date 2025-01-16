# Contributing

- [Pull requests](#pull-requests)
- [Quick start](#quick-start)
- [Unit tests](#unit-tests)

## Pull requests

Make sure all the following are true when creating a pull-request:

- The PR has a meaningful title [guidelines](https://github.com/kubernetes/community/blob/master/contributors/guide/pull-requests.md#use-imperative-mood-in-your-commit-message-subject).
- The PR is labeled accordingly with a **single** label.
- Unit tests have been updated/added.

## Quick start

1. Fork the dynatrace-bootstrapper repository and get the source code:

    ```sh
    git clone https://github.com/<your_username>/dynatrace-bootstrapper
    cd dynatrace-bootstrapper
    ```

2. Install development prerequisites:

   ```sh
   make prerequisites
   ```

3. Create a new branch to work on:

    > Group your branch into a category using a prefix for your branch name, like `feature/`, `ci/`, `bugfix/`, `doc/`.

   ```sh
   git checkout -b feature/your-branch
   ```

4. Once the changes are finished, make sure there are no warnings in the code.

    > **NOTE:**
    > Unit tests can also be automatically run via pre-commit hook, installed by running `make prerequisites/setup-pre-commit`.
    > With the pre-commit hook can only commit code that passes all checks.

    ```sh
    make test
    ```

5. Create a pull request from the fork ([see guide](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)), with a proper title and fill out the description template. Once everything is ready, set the PR ready for review.

6. A maintainer will review the pull request and make comments. It's preferable to add additional commits over amending and force-pushing since it can be difficult to follow code reviews when the commit history changes. Commits will be squashed when they're merged.

## Unit tests

Run the go unit tests via make:

```sh
make test
```

### Mocking

For our mocking needs we trust in [testify](https://github.com/stretchr/testify) while using [mockery](https://github.com/vektra/mockery) to generate our mocks.
We check in our mocks to improve code readability especially when reading via GitHub and to remove the dependency on make scripts to run our tests.
Mockery only has to be run when adding new mocks or have to be updated to changed interfaces.

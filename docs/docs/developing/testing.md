---
sidebar_position: 3
---

# Testing

## Unit tests

To run all unit tests for all modules, you can use the following command:

```bash
$ just test-unit
```

## End-to-end tests

To run all end-to-end tests for all modules, you can use the following command:

```bash
$ just test-e2e
```

The recipe takes an optional argument for which image-versions (docker tag) to use (e.g. `just test-e2e latest`).
If no argument is provided, it will default to `local` and also build the docker images locally with the local tag.

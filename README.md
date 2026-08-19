# go-test-app

Minimal Go program used as a build target for the rollout-computer CI.
Its pipeline (`.ci/pipeline.yaml`) runs a build step, a parallel group
(unit tests + vet), and a summary step inside one microVM job.

Push this repository to a public remote, then point the CI's "New build"
dialog at it.

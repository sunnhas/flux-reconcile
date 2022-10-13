# flux-reconcile

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/sunnhas/flux-reconcile?label=latest%20release&logo=github&sort=semver)
![License](https://img.shields.io/github/license/sunnhas/flux-reconcile)

Is a CLI to trigger reconciliation in a [FluxCD]() cluster.

## Install

Binaries are released on [GitHub](https://github.com/sunnhas/flux-reconcile/releases) within releases to the right.
You can download the binary for your platform manually.

Put it in one of the system bin folders or somewhere on your `PATH` and make it executable `chmod +x flux-reconcile`.

## How to use?

When installed you can run `flux-reconcile`:

```shell
flux-reconcile {webhook} --key secret-key
```

Per default flux-reconcile will trigger a `GitRepository` reconciliation.
To specify which resources to target use `--resources` flag.

```shell
flux-reconcile {webhook} --key secret-key --resources GitRepository,ImageAutomation
```

You can put the signature key in the environment with `FR_KEY` instead of the `--key` flag.

```shell
export FR_KEY=secret-key
flux-reconcile {webhook}
```

## Flux cluster requirements

Before reconciliation can be triggered the cluster needs to have a webhook receiver of `generic-hmac`.
The configuration of this can be found at the [Flux documents](https://fluxcd.io/flux/components/notification/receiver/#generic-hmac-receiver).

When that is configured you will need the 3 things below:

- Endpoint for the receiver
- The secret key used for the signature
- The resources you want to trigger

[![CI](https://github.com/RedbackThomson/provider-metronome/actions/workflows/ci.yml/badge.svg)](https://github.com/RedbackThomson/provider-metronome/actions/workflows/ci.yml)
[![GitHub release](https://img.shields.io/github/release/redbackthomson/provider-metronome/all.svg?style=flat-square)](https://github.com/redbackthomson/provider-metronome/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/redbackthomson/provider-metronome)](https://goreportcard.com/report/github.com/redbackthomson/provider-metronome)

# provider-metronome

`provider-metronome` is a Crossplane Provider that enables deployment and management
of [Metronome](https://metronome.com) resources provisioned by Crossplane.

The provider currently supports the following resources:

- [BillableMetric](https://docs.metronome.com/api/#billable-metrics)
- [CustomFieldKey](https://docs.metronome.com/api/#custom-fields)
- [Product](https://docs.metronome.com/api/#products)
- [Rate](https://docs.metronome.com/api/#rate-cards)
- [RateCard](https://docs.metronome.com/api/#rate-cards)

## Install

If you would like to install `provider-metronome` without modifications, you may do
so using the Crossplane CLI in a Kubernetes cluster where Crossplane is
installed:

```console
crossplane xpkg install provider xpkg.upbound.io/redbackthomson/provider-metronome:v0.0.1
```

Then you will need to create a `ProviderConfig` that specifies the API key
connected to your Metronome account. This is commonly done by storing the API
key in a secret that the `ProviderConfig` references.

## Developing locally

**Pre-requisite:** A Kubernetes cluster with Crossplane installed

To run the `provider-metronome` controller against your existing local cluster,
simply run:

```console
make run
```

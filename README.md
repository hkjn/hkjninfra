hkjninfra
=======

Repo hkjninfra holds some infrastructure plans for hkjn.me.

## Dockerized Terraform

A dockerized Terraform alias `tf` can be added by:

```
source tf_dockerized.sh
```

After this, we can use `tf plan`, `tf apply` and other commands.

## Dockerized gcloud tools

An alias for working with a dockerized set of `gcloud` tools can be
added by the command:

```
source gcloud_dockerized.sh
```

After this, the `gcd` command enters an interactive container which
can work towards the GCE project.

## Ignition format of  `user-data`

The `user-data` field for CoreOS machines should be in Ignition format:

* https://coreos.com/os/docs/latest/booting-on-google-compute-engine.html

The more human-friendly cloud-config YAML format can be used, and
transpiled to Ignition JSON using the `ct` tool:

* https://github.com/coreos/container-linux-config-transpiler/releases

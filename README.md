hkjninfra
=======

Repo hkjninfra holds some infrastructure plans for hkjn.me.

## Dockerized Terraform

A dockerized Terraform alias `tf` can be added by:

```
source tf_dockerized.sh
```

## Remote state

The project uses GCP to store state remotely:

```
bash remote_config.sh
```
source tf_dockerized.sh
run_tf remote config \
       -backend=gcs \
       -backend-config="bucket=hkjn-terraform-state" \
       -backend-config="path=hkjninfra/prod/terraform.tfstate" \
       -backend-config="project=henrik-jonsson"

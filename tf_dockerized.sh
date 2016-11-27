run_tf() {
  docker run --rm -it -v $(pwd):/home/tfuser \
         -e GOOGLE_APPLICATION_CREDENTIALS=/home/tfuser/.gcp/tf-dns-editor.json \
         hkjn/terraform $@
}

alias tf='run_tf $@'

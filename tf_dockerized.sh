run_tf() {
	python ignite.py
	if [[ $? -ne 0 ]]; then
		echo "ignite.py failed, bailing." >&2
		return
	fi
	docker run --rm -it -v $(pwd):/home/tfuser \
	           -e GOOGLE_APPLICATION_CREDENTIALS=/home/tfuser/.gcp/tf-dns-editor.json \
		   hkjn/terraform $@ -var 'git_commit='$(last_commit)''
}

alias tf='run_tf $@'

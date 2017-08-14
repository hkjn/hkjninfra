run_tf() {
	for f in *.yml; do
		echo "Generating ${f}.json.."
		ct < ${f} > ${f}.json
		if [[ $? -ne 0 ]]; then
			echo "FATAL: Failed to generate .json using ct." >&2
			return
		fi
	done
	docker run --rm -it -v $(pwd):/home/tfuser \
	           -e GOOGLE_APPLICATION_CREDENTIALS=/home/tfuser/.gcp/tf-dns-editor.json \
	hkjn/terraform $@
}

alias tf='run_tf $@'

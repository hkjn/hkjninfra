run_tf() {
	for host in zg1 zg3; do
		echo "Generating ${host}.json.."
		python ignite.py ${host} > bootstrap_${host}.json
		if [[ $? -ne 0 ]]; then
			echo "ignite.py failed, bailing." >&2
			return
		fi
	done
	docker run --rm -it -v $(pwd):/home/tfuser \
	           -e GOOGLE_APPLICATION_CREDENTIALS=/home/tfuser/.gcp/tf-dns-editor.json \
	hkjn/terraform $@
}

alias tf='run_tf $@'

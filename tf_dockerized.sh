ignite_diff() {
	local target
	target=${1}
	if [[ ! -e bootstrap/bootstrap_${target}.json ]]; then
		echo "Unknown target ${target}." >&2
		return 1
	fi
	run_tf output hkjn_ignite_json_${target} > /tmp/${target}_output.json
	if [[ $? -ne 0 ]]; then
		echo "tf output command failed: $(cat /tmp/${target}_output.json)" >&2
	fi
	diff <(jq '.' < /tmp/${target}_output.json) <(jq '.' < bootstrap/bootstrap_${target}.json)
}

run_tf() {
	local action
	action=$1
	if [[ "${action}" = plan ]]; then
		python ignite.py
		if [[ $? -ne 0 ]]; then
			echo "ignite.py failed, bailing." >&2
			return
		fi
	fi
	# TODO: Below we take current VERSION file, but could run an older version
	# for some targets, as specified in ignition.py..
	echo "version = \"$(cat VERSION)\"" > terraform.tfvars
	# TODO: By doing 'tf plan -detailed-exitcode', we can check for status 2 (there was
	# a diff), and if so, run the ignite_diff command above for a not-horrible comparison
	# on where the metadata differs.
	docker run --rm -it -v $(pwd):/home/tfuser \
	           -e GOOGLE_APPLICATION_CREDENTIALS=/home/tfuser/.gcp/tf-dns-editor.json \
		   hkjn/terraform $@
}

alias tf='run_tf $@'

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
	# TODO: Could store user-data fields with 'tf output', then
	# diff the before / after and show better diffs:
	# 1. tf output hkjn_ignite_json_zg1 > zg1.json
	# 2. diff <(jq '.' < zg1.json) <(jq '.' < bootstrap/bootstrap_zg1.json)
	# By doing 'tf plan -detailed-exitcode', we can check for status 2 (there was
	# a diff), and if so, run the diff command above for a not-horrible comparison
	# on where the metadata differs
	docker run --rm -it -v $(pwd):/home/tfuser \
	           -e GOOGLE_APPLICATION_CREDENTIALS=/home/tfuser/.gcp/tf-dns-editor.json \
		   hkjn/terraform $@
}

alias tf='run_tf $@'

#!usr/bin/env python
#
# Create Ignition JSON for specified instance.
#
import json


INSTANCES = {
    'core': '1.4.3',
    'decenter_world': '1.4.3',
    #'builder': '1.2.5',
}


UPDATE_CONF_FILE = {
    'filesystem': 'root',
    'path': '/etc/coreos/update.conf',
    'contents': {
        'source': 'data:,GROUP%3Dbeta%0AREBOOT_STRATEGY%3D%22etcd-lock%22',
        'verification': {},
    },
    'mode': 420,
    'user': {},
    'group': {},
}


def new_file(filename, checksum, url):
    """Return Ignition config for specified file.
    
    Args:
        filename: str with filename under /opt/bin.
        checksum: str with expected checksum of file.
    Returns:
        Dict describing Ignitiion config for file.
    """
    return {
        'filesystem': 'root',
        'path': '/opt/bin/{}'.format(filename),
        'contents': {
            'source': url,    
            'verification': {'hash': 'sha512-{}'.format(checksum) },
        },
        'mode': 493,
        'user': {},
        'group': {},
    }


def new_unit(unit):
    """Return Ignition config for a systemd unit.

    Args:
        unit: A str like 'bitcoin.service', identifying a file under units/.
    Returns:
        Dict with Ignition config for the systemd unit.
    """

    unit_contents = ''
    with open('units/{}'.format(unit)) as unit_file:
        unit_contents = unit_file.read()
    return {
        'name': unit,
        'enable': True,
        'contents': unit_contents,
    }


def new_unit_dropin(unit, dropin_name):
    """Return Ignition config for a systemd dropin for a unit.

    Args:
        unit: A str like 'bitcoin.service', identifying the unit of the dropin.
        dropin_name: A str like '10_override_storage.conf', identifying the name of the dropin.
    Returns:
        Dict with Ignition config for the systemd unit dropin.
    """

    dropin_contents = ''
    with open('units/{}'.format(dropin_name)) as dropin_file:
        dropin_contents = dropin_file.read()
    return {
        'name': unit,
        'dropins': [{
            'name': dropin_name,
            'contents': dropin_contents,
        }]
    }


def get_checksums(version):
    result = {}
    with open('{}.sha512'.format(version)) as checksum_file:
        for line in checksum_file.readlines():
            parts = line.split()
            if len(parts) != 2:
                raise RuntimeError('Invalid line in checksum file: {}'.format(line))
            checksum, release_file = parts[0], parts[1]
            result[release_file] = checksum
    return result

    
def get_config(instance, version, checksums):
    """Returns Ignition config for the instance.
    
    Returns:

        Dict with Ignition config.
    """
    
    shared_files = [
        UPDATE_CONF_FILE,
        new_file('gather_facts', checksums['gather_facts'], 'https://github.com/hkjn/hkjninfra/releases/download/{}/gather_facts'.format(version)),
        new_file('tclient', checksums['tclient_x86_64'], 'https://github.com/hkjn/hkjninfra/releases/download/{}/tclient_x86_64'.format(version)),
    ]
    shared_units = [
        new_unit('tclient.service'),
        new_unit('tclient.timer'),
    ]
    files = []
    units = []
    filesystem = []
    decenter_version = '1.1.5' # TODO: Should come from fetch + checksums file.
    instance_configs = {
        'core': {
            'files': [],
            'units': [
                new_unit_dropin('docker.service', '10_override_storage.conf'),
                new_unit('bitcoin.service'),
                new_unit('containers.mount'),
            ],
        },
        'decenter_world': {
            'files': [
                new_file(
                    'decenter_world',
                    checksums['decenter_world_x86_64'],
                    'https://github.com/hkjn/decenter.world/releases/download/{}/decenter_world_x86_64'.format(decenter_version),
                ),
                new_file(
                    'decenter_redirector',
                    checksums['decenter_redirector_x86_64'],
                    'https://github.com/hkjn/decenter.world/releases/download/{}/decenter_redirector_x86_64'.format(decenter_version),
                ),
            ],
            'units': [
                new_unit('decenter.service'),
                new_unit('decenter_redirector.service'),
                new_unit('etc-secrets.mount'),
            ],
        },
    }
    if instance not in instance_configs:
        raise RuntimeError('Unknown instance {}'.format(instance))

    return {
        'ignition': {
            'version': '2.0.0',
            'config': {}
        },
        'storage': {
            'filesystem': filesystem,
            'files': shared_files + instance_configs[instance]['files'],
        },
        'systemd': {
            'units': shared_units + instance_configs[instance]['units'],
            'networkd': {},
            'passwd': {},
        },
    }


def run():
    """Generate Ignition JSON config files for all instances.
    """

    print('Generating Ignition JSON..')
    for instance in sorted(INSTANCES):
        version = INSTANCES[instance]
        checksums = {}
        print('Using checksums from version {} for {}..'.format(version, instance))
        try:
            checksums = get_checksums(version)
        except IOError as ioerr:
            raise RuntimeError('Checksums unavailable: {}'.format(version, ioerr))
        for release_file in sorted(checksums):
            print('Checksum for {} {}: {}'.format(release_file, version, checksums[release_file]))

        json_path = 'bootstrap/bootstrap_{}.json'.format(instance)
        print('Generating {}..'.format(json_path))
        with open(json_path, 'w') as json_file:
            json_file.write(json.dumps(get_config(instance, version, checksums)))


if __name__ == '__main__':
    run()

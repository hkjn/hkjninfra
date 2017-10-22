
import json
import os


INSTANCES = {
    'core': '1.4.3',
    'decenter_world': '1.5.0',
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


def new_file(path, checksum, url):
    """Return Ignition config for specified file.
    
    Args:
        path: str with full path to file on remote host.
        checksum: str with expected checksum of file.
    Returns:
        Dict describing Ignitiion config for file.
    """
    return {
        'filesystem': 'root',
        'path': path,
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

    
def get_config(instance, hkjninfra_version, decenter_version, checksums, sshash):
    """Returns Ignition config for the instance.
    
    Args:
        instance: str with the instance to create config for.
        hkjninfra_version: version of hkjninfra binaries to run on the instance.
        decenter_version: version of decenter.world binaries to run on the
        instance.
        checksums: dict of binary name to checksums for binaries.
        sshash: str with the secretservice hash to use.
    Returns:
        Dict with Ignition config.
    """
    
    shared_files = [
        UPDATE_CONF_FILE,
        new_file(
            '/opt/bin/gather_facts',
            checksums['gather_facts'],
            'https://github.com/hkjn/hkjninfra/releases/download/{}/gather_facts'.format(hkjninfra_version),
        ),
        new_file(
            '/etc/ssl/mon_ca.pem',
            'cf8032384a17fc591f83030a3a536f8b80e79c7d6e5e839e5003b587163bf371a51ab1b14dc047486cfd55fb74238d69253fcb27d49245ba249359584b169bb4',
            'https://admin1.hkjn.me/{0}/files/certs/mon_ca.pem'.format(sshash),
        ),
        new_file(
            '/opt/bin/tclient',
            checksums['tclient_x86_64'],
            'https://github.com/hkjn/hkjninfra/releases/download/{}/tclient_x86_64'.format(hkjninfra_version),
        ),
    ]
    shared_units = [
        new_unit('tclient.service'),
        new_unit('tclient.timer'),
    ]
    files = []
    units = []
    filesystem = []
    instance_configs = {
        'core': {
            'files': [
            ],
            'units': [
                new_unit_dropin('docker.service', '10_override_storage.conf'),
                new_unit('bitcoin.service'),
                new_unit('containers.mount'),
            ],
        },
        'decenter_world': {
            'files': [
                new_file(
                    '/opt/bin/decenter_world',
                    checksums['decenter_world_x86_64'],
                    'https://github.com/hkjn/decenter.world/releases/download/{}/decenter_world_x86_64'.format(decenter_version),
                ),
                new_file(
                    '/opt/bin/decenter_redirector',
                    checksums['decenter_redirector_x86_64'],
                    'https://github.com/hkjn/decenter.world/releases/download/{}/decenter_redirector_x86_64'.format(decenter_version),
                ),
                new_file(
                    '/etc/ssl/client.pem',
                    'e9c89e72d89d93fb7015c63413927913cbfca8eb5cd50523b105d0c37a89bf8ab8965a4b489c5ec54d64e9afb4a480aacac52b47be03d1044ad1d4ca9cd2844d',
                    'https://admin1.hkjn.me/{0}/files/certs/client.pem'.format(sshash),
                ),
                new_file(
                    '/etc/ssl/client-key.pem',
                    'e01f327e1089d652a521ddc402171dcdaeaf70f8f59b916a426efabe9a8e5ae78d3842cbb7e0caa11630b33b6b1ceb7e1f8871bedfde0f312a9b577697269d6b',
                    'https://admin1.hkjn.me/{0}/files/certs/client-key.pem'.format(sshash),
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

    decenter_version = '1.1.7' # TODO: Should come from fetch + checksums file.
    sshash = os.environ.get('SECRETSERVICE_HASH')
    if not sshash:
        raise RuntimeError('No SECRETSERVICE_HASH set in environment.')
    print('Generating Ignition JSON..')
    for instance in sorted(INSTANCES):
        hkjninfra_version = INSTANCES[instance]
        checksums = {}
        print('Using checksums from version {} for {}..'.format(hkjninfra_version, instance))
        try:
            checksums = get_checksums(hkjninfra_version)
        except IOError as ioerr:
            raise RuntimeError('Checksums unavailable: {}'.format(hkjninfra_version, ioerr))
        for release_file in sorted(checksums):
            print('Checksum for {} {}: {}'.format(release_file, hkjninfra_version, checksums[release_file]))

        json_path = 'bootstrap/bootstrap_{}.json'.format(instance)
        print('Generating {}..'.format(json_path))
        with open(json_path, 'w') as json_file:
            conf = get_config(instance, hkjninfra_version, decenter_version, checksums, sshash)
            json_file.write(json.dumps(conf))


if __name__ == '__main__':
    run()

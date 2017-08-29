#!usr/bin/env python
#
# Create Ignition JSON for specified instance.
#
import json


INSTANCES = {
    'zg1': '1.1.4',
    'zg3': '1.1.4',
}


def get_shared_files(version, checksums):
    """Return Ignition config for the shared files for all instances.
    
    Returns:
        List of dict of files.
    """
    return [
        {
            'filesystem': 'root',
            'path': '/etc/coreos/update.conf',
            'contents': {
                'source': 'data:,GROUP%3Dbeta%0AREBOOT_STRATEGY%3D%22etcd-lock%22',
                'verification': {},
            },
            'mode': 420,
            'user': {},
            'group': {},
        }, {
            'filesystem': 'root',
            'path': '/opt/bin/gather_facts',
            'contents': {
                'source': 'https://github.com/hkjn/hkjninfra/releases/download/{}/gather_facts'.format(version),
                'verification': {
                    'hash': 'sha512-{}'.format(checksums['gather_facts']),
                },
            },
            'mode': 493,
            'user': {},
            'group': {},
        }, {
            'filesystem': 'root',
            'path': '/opt/bin/report_client',
            'contents': {
                'source': 'https://github.com/hkjn/hkjninfra/releases/download/{}/tclient_x86_64'.format(version),
                'verification': {
                'hash': 'sha512-{}'.format(checksums['tclient_x86_64']),
                },
            },
            'mode': 493,
            'user': {},
            'group': {},
        },
    ]


def get_shared_units():
    """Return Ignition config for the shared systemd units for all instances.
    
    Returns:
        List of dict of systemd units.
    """
    rsc = ''
    with open('units/report_client.service') as rc:
        rcs = rc.read()
    return [
        {
            'name': 'report_client.service',
            'enable': True,
            'contents': rcs,
        }, {
            'name': 'report_client.timer',
            'enable': True,
            'contents': '[Unit]\nDescription=Timer that starts report_client.service\n\n[Timer]\n# Run every 5 min.\nOnCalendar=*:0/5\nPersistent=true\n\n[Install]\nWantedBy=multi-user.target\n',
        },
    ]


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
    shared_files = get_shared_files(version, checksums)
    shared_units = get_shared_units()
    files = []
    units = []
    filesystem = []
    if instance == 'zg1':
#        files = [
#            {
#                'filesystem': 'root',
#                'path': '/etc/systemd/system/docker.service.d/10-override-storage.conf',
#                'contents': '[Service]\nEnvironment=\"DOCKER_OPTS=-g /containers/docker -s overlay2\"',
#            },
#        ]
        units = [
            {
# TODO: Reenable the dropin to change storage for docker.
#                'name': 'docker.service',
#                'dropins': [
#                    {
#                        'name': '10_override_storage.conf',
#                        'contents': '[Service]\nEnvironment=\"DOCKER_OPTS=-g /containers/docker -s overlay2\"',
#                    },
#                ],
#            }, {
                'name': 'bitcoin.service',
                'enable': True,
                'contents': '[Unit]\nDescription=bitcoind\nAfter=network-online.target\n\n[Service]\nExecStartPre=-/bin/bash -c \"docker pull hkjn/bitcoin:$(uname -m)\"\nExecStartPre=-/usr/bin/docker stop bitcoin\nExecStartPre=-/usr/bin/docker rm bitcoin\nExecStart=/bin/bash -c \" \\\n  docker run --name bitcoin \\\n             -p 8333:8333 \\\n             --memory=1050m \\\n             --cpu-shares=128 \\\n             -v /containers/bitcoin:/home/bitcoin/.bitcoin \\\n             hkjn/bitcoin:$(uname -m)\"\nRestart=always\n\n[Install]\nWantedBy=multi-user.target\n',
            }, {
                'name': 'containers.mount',
                'enable': True,
                'contents': '[Mount]\nWhat=/dev/disk/by-id/scsi-0Google_PersistentDisk_persistent-disk-1-part1\nWhere=/containers\nType=xfs\n\n[Install]\nRequiredBy=local-fs.target\n',
            },
        ]
    elif instance == 'zg3':
        filesystems = [{
            'mount': {
                'device': '/dev/disk/by-id/scsi-0Google_PersistentDisk_persistent-disk-1',
                'format': 'ext4',
            },
        }]
        files = [
            {
                "filesystem": "root",
                "path": "/opt/bin/decenter_world",
                "contents": {
                    "source": "https://github.com/hkjn/decenter.world/releases/download/1.1.2/decenter_world_x86_64",
                    "verification": {
                        'hash': 'sha512-ed0fa9f29b504fb30ce7c33afc743e636bccffa6a9bd5630f9fd374cf6076725e6d44d8e2b11ed82f849de90cf009199bf2f19aa803ffd1830ddd75a837aeb06',
                    },
            },
                "mode": 493,
                "user": {},
                "group": {},
            }, {
                "filesystem": "root",
                "path": "/opt/bin/decenter_redirector",
                "contents": {
                    "source": "https://github.com/hkjn/decenter.world/releases/download/1.1.1/decenter_redirector_x86_64",
                    "verification": {
                        "hash": "sha512-8026412bcc856bb073e01a5b984e0a2161049b76e575cfc506bf733a03ca70ed2fffe0a83d269a578936a545066500603b3a86c9e6a47905376108b8af41837e",
                    },
                },
                "mode": 493,
                "user": {},
                "group": {},
            },
        ]
        units = [
            {
               "name": "decenter.service",
                "enable": True,
                "contents": "[Unit]\nDescription=decenter.world server\nAfter=network-online.target\n\n[Service]\nEnvironment=PATH=/usr/bin/:/opt/bin:/bin\nEnvironment=DECENTER_WORLD_ADDR=:443\nExecStart=/usr/bin/bash -c \"decenter_world\"\nRestart=always\n\n[Install]\nWantedBy=multi-user.target\n",
            }, {
                "name": "decenter_redirector.service",
                "enable": True,
                "contents": "[Unit]\n Description=decenter.world redirector server\n After=network-online.target\n\n [Service]\n Environment=PATH=/usr/bin/:/opt/bin:/bin\n ExecStart=/usr/bin/bash -c \"decenter_redirector\"\n Restart=always\n\n [Install]\n WantedBy=multi-user.target\n"
            }, {
                "name": "etc-secrets.mount",
                "enable": True,
                "contents": "[Mount]\nWhat=/dev/disk/by-id/scsi-0Google_PersistentDisk_persistent-disk-1\nWhere=/etc/secrets\nType=ext4\n\n[Install]\nRequiredBy=local-fs.target\n"
            },
        ]
    else:
        raise RuntimeError('Unknown instance {}'.format(instance))

    return {
        'ignition': {
            'version': '2.0.0',
                'config': {}
            },
            'storage': {
                'filesystem': filesystem,
                'files': shared_files + files,
            },
            'systemd': {
                'units': shared_units + units,
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

        json_path = 'bootstrap_{}.json'.format(instance)
        print('Generating {}..'.format(json_path))
        with open(json_path, 'w') as json_file:
            json_file.write(json.dumps(get_config(instance, version, checksums)))


if __name__ == '__main__':
    run()

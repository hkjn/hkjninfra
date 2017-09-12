#!/usr/bin/env python

import ignite

import unittest


class IgniteTest(unittest.TestCase):
    def test_get_config(self):
        self.maxDiff = None
        checksums = {
            'gather_facts': '55bb96874add4d200274cf1796c622da8e92244ad5b5fa15818bc516c5ed249e9cd98a736d44b66c7e03ca2b52e5aa898717fbd7d08ff13cd94de38ba2aef8c8',
            'tclient_x86_64': 'bf080645783c999f1a2bc8bc306660df8dbf496c6b7f98cf1d257d43c544050a7f4b6d1d9ba962c1d45fae8eb373061d3350e191edd73f1b44f01fc01448177f',
        }
        got = ignite.get_config('core', '1.1.0', checksums)
        want = {
                'ignition': {'config': {}, 'version': '2.0.0'},
                'storage': {
                    'files': [{
                        'contents': {
                            'source': 'data:,GROUP%3Dbeta%0AREBOOT_STRATEGY%3D%22etcd-lock%22',
                            'verification': {},
                        },
                        'filesystem': 'root',
                        'group': {},
                        'mode': 420,
                        'path': '/etc/coreos/update.conf',
                        'user': {},
                    }, {
                        'contents': {
                            'source': 'https://github.com/hkjn/hkjninfra/releases/download/1.1.0/gather_facts',
                            'verification': {
                                'hash': 'sha512-55bb96874add4d200274cf1796c622da8e92244ad5b5fa15818bc516c5ed249e9cd98a736d44b66c7e03ca2b52e5aa898717fbd7d08ff13cd94de38ba2aef8c8',
                            },
                        },
                        'filesystem': 'root',
                        'group': {},
                        'mode': 493,
                        'path': '/opt/bin/gather_facts',
                        'user': {},
                    },
                    {
                        'contents': {
                            'source': 'https://github.com/hkjn/hkjninfra/releases/download/1.1.0/tclient_x86_64',
                            'verification': {
                                'hash': 'sha512-bf080645783c999f1a2bc8bc306660df8dbf496c6b7f98cf1d257d43c544050a7f4b6d1d9ba962c1d45fae8eb373061d3350e191edd73f1b44f01fc01448177f',
                            },
                        },
                        'filesystem': 'root',
                        'group': {},
                        'mode': 493,
                        'path': '/opt/bin/report_client',
                        'user': {
                        },
                    }],
                    'filesystem': [],
                },
                'systemd': {
                    'networkd': {},
                    'passwd': {},
                    'units': [
                        {
                            'contents': '[Unit]\nDescription=report client\nAfter=network-online.target\n\n[Service]\nEnvironment=PATH=/usr/bin/:/opt/bin:/bin\nEnvironment=REPORT_ADDR=mon.hkjn.me:50051\nEnvironment=REPORT_FACTS_PATH=/etc/report_facts.json\nExecStartPre=-/bin/bash -c "gather_facts > /etc/report_facts.json"\nExecStart=/bin/bash -c report_client\n\n[Install]\nWantedBy=multi-user.target\n\n',
                            'enable': True,
                            'name': 'report_client.service',
                        }, {
                            'contents': '[Unit]\nDescription=Timer that starts report_client.service\n\n[Timer]\n# Run every 5 min.\nOnCalendar=*:0/5\nPersistent=true\n\n[Install]\nWantedBy=multi-user.target\n',
                            'enable': True,
                            'name': 'report_client.timer',
                        }, {
                            'dropins': [
                                {
                                    'contents': '[Service]\nEnvironment="DOCKER_OPTS=-g /containers/docker -s overlay2"\n',
                                    'name': '10_override_storage.conf',
                                },
                            ],
                            'name': 'docker.service',
                        }, {
                            'contents': '[Unit]\nDescription=bitcoind\nAfter=network-online.target\n\n[Service]\nExecStartPre=-/bin/bash -c "docker pull hkjn/bitcoin:$(uname -m)"\nExecStartPre=-/usr/bin/docker stop bitcoin\nExecStartPre=-/usr/bin/docker rm bitcoin\nExecStart=/bin/bash -c " \\\n  docker run --name bitcoin \\\n             -p 8333:8333 \\\n             --memory=1050m \\\n             --cpu-shares=128 \\\n             -v /containers/bitcoin:/home/bitcoin/.bitcoin \\\n             hkjn/bitcoin:$(uname -m) -dbcache=500 -printtoconsole"\nRestart=always\n\n[Install]\nWantedBy=multi-user.target\n',
                            'enable': True,
                            'name': 'bitcoin.service',
                        }, {
                            'contents': '[Mount]\nWhat=/dev/disk/by-id/scsi-0Google_PersistentDisk_persistent-disk-1-part1\nWhere=/containers\nType=xfs\n\n[Install]\nRequiredBy=local-fs.target\n',
                            'enable': True,
                            'name': 'containers.mount',
                        },
                    ],
                }
            }
        self.assertEqual(want, got)


if __name__ == '__main__':
    unittest.main()

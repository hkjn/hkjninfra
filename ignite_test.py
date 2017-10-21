#!/usr/bin/env python

import ignite

import unittest


class IgniteTest(unittest.TestCase):
    def test_get_config(self):
        self.maxDiff = None
        checksums = {
            'decenter_world_x86_64': 'ce1ebb0b0192b61a6d94eba19db94d8fc4dbb72e8188ac69cc1c88e209f68ae9c57d32b07287abbc4149439ae41ba1f519bd4360dc3b8c7add47f65a52d33dab',
            'decenter_redirector_x86_64': '735b7d08a457918eadfc049cc0063466756a7684bb5bd7cc660afcc4660a87d8685fca361e4b01678a0f01aea79f6831b9f1ca1c184bd44a7a2a63cb2b534d7c',
            'gather_facts': '55bb96874add4d200274cf1796c622da8e92244ad5b5fa15818bc516c5ed249e9cd98a736d44b66c7e03ca2b52e5aa898717fbd7d08ff13cd94de38ba2aef8c8',
            'tclient_x86_64': 'bf080645783c999f1a2bc8bc306660df8dbf496c6b7f98cf1d257d43c544050a7f4b6d1d9ba962c1d45fae8eb373061d3350e191edd73f1b44f01fc01448177f',
        }
        got = ignite.get_config('core', '1.1.0', checksums, 'fakehash')
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
                    }, {
                        'contents': {
                            'source': 'https://admin1.hkjn.me/fakehash/files/certs/mon_ca.pem',
                            'verification': {
                                'hash':
                                'sha512-cf8032384a17fc591f83030a3a536f8b80e79c7d6e5e839e5003b587163bf371a51ab1b14dc047486cfd55fb74238d69253fcb27d49245ba249359584b169bb4',
                                }
                            },
                        'filesystem': 'root',
                        'group': {},
                        'mode': 493,
                        'path': '/etc/ssl/mon_ca.pem',
                        'user': {},
                    }, {
                        'contents': {
                            'source': 'https://github.com/hkjn/hkjninfra/releases/download/1.1.0/tclient_x86_64',
                            'verification': {
                                'hash': 'sha512-bf080645783c999f1a2bc8bc306660df8dbf496c6b7f98cf1d257d43c544050a7f4b6d1d9ba962c1d45fae8eb373061d3350e191edd73f1b44f01fc01448177f',
                            },
                        },
                        'filesystem': 'root',
                        'group': {},
                        'mode': 493,
                        'path': '/opt/bin/tclient',
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
                            'contents': '\n'.join([
                                '[Unit]',
                                'Description=tclient',
                                'After=network-online.target',
                                '',
                                '[Service]',
                                'Environment=REPORT_ADDR=mon.hkjn.me:50051',
                                'Environment=REPORT_FACTS_PATH=/etc/report_facts.json',
                                'Environment=REPORT_TLS_CA_CERT=/etc/ssl/mon_ca.pem',
                                'Environment=REPORT_TLS_CERT=/etc/ssl/client.pem',
                                'Environment=REPORT_TLS_KEY=/etc/ssl/client-key.pem',
                                '',
                                'ExecStartPre=-/bin/bash -c "/opt/bin/gather_facts > /etc/report_facts.json"',
                                'ExecStart=/opt/bin/tclient',
                                '',
                                '[Install]',
                                'WantedBy=multi-user.target',
                                '',
                            ]),
                            'enable': True,
                            'name': 'tclient.service',
                        }, {
                            'contents': '[Unit]\nDescription=Timer that starts report_client.service\n\n[Timer]\n# Run every 5 min.\nOnCalendar=*:0/5\nPersistent=true\n\n[Install]\nWantedBy=multi-user.target\n',
                            'enable': True,
                            'name': 'tclient.timer',
                        }, {
                            'dropins': [
                                {
                                    'contents': '[Service]\nEnvironment="DOCKER_OPTS=-g /containers/docker -s overlay2"\n',
                                    'name': '10_override_storage.conf',
                                },
                            ],
                            'name': 'docker.service',
                        }, {
                            'contents': '[Unit]\nDescription=bitcoind\nAfter=network-online.target\n\n[Service]\nExecStartPre=-/bin/bash -c "docker pull hkjn/bitcoin:$(uname -m)"\nExecStartPre=-/usr/bin/docker stop bitcoin\nExecStartPre=-/usr/bin/docker rm bitcoin\nExecStart=/bin/bash -c " \\\n  docker run --name bitcoin \\\n             -p 8333:8333 \\\n             --memory=1050m \\\n             --cpu-shares=128 \\\n             -v /containers/bitcoin:/home/bitcoin/.bitcoin \\\n             hkjn/bitcoin:$(uname -m) -dbcache=800 -onlynet=ipv4 -printtoconsole"\nRestart=always\n\n[Install]\nWantedBy=multi-user.target\n',
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

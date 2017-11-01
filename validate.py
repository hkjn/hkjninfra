import json
import sys
import urllib2

data=json.loads(open("bootstrap/{}.json".format(sys.argv[1])).read())
for f in data['storage']['files']:
    if 'source' not in f['contents'] or 'http' not in f['contents']['source']:
        continue
    url = f['contents']['source']
    digest = f['contents']['verification']['hash'].lstrip('sha512-')
    print('{} {}'.format(url, digest))
    print('Fetching {}..'.format(url))
    response = urllib2.urlopen(url)
    html = response.read()
    with open('/tmp/{}'.format(digest), 'w+') as tmpfile:
        tmpfile.write(html)
        print('Wrote /tmp/{}'.format(digest))


# if 'source', fetch and compare with 'verification'['hash']



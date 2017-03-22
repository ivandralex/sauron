import sys
import string
import pandas

stop_words = ['iPad', 'iPhone', 'Android', 'Phone', 'Touch', 'Macintosh']

data_path = sys.argv[1]
output_path = sys.argv[2]
label = sys.argv[3]
data = pandas.read_csv(data_path)

labeled = data.query('label == ' + label)

ips = []
user_agents = []

with open(output_path, 'w+') as f:
    for i,r in enumerate(labeled.index):
        user_agent = labeled.ix[r]['user_agent']
        ip = labeled.ix[r]['ip']
        skip = False
        if user_agent == '-':
            skip = True
        for word in stop_words:
            if word in user_agent:
                skip = True
                break
        if skip:
            continue

        if ip in ips:
            continue
        if not user_agent in user_agents:
            print user_agent
            user_agents.append(user_agent)
        #ips.append(ip)
        #line = '"' + ip + '","' + user_agent + '"'
        line = ip
        f.write(line + '\n')
    f.close()

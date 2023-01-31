import json
import sys
import getopt
import os
import socket


def usage():
    print """usage:
python jxwaf_init.py --api_key=a2dde899-96a7-40s2-88ba-31f1f75f1552 --api_password=653cbbde-1cac-11ea-978f-2e728ce88125 --waf_server=http://192.168.1.1
"""


file_path = "/opt/jxwaf/nginx/conf/jxwaf/jxwaf_config.json"


def main(argv):
    try:
        opts, args = getopt.getopt(sys.argv[1:], 'h', ['help', 'api_key=', 'api_password=', 'waf_server='])
    except getopt.GetoptError:
        usage()
        sys.exit()
    if os.path.exists(file_path) == False:
        print "Error: /opt/jxwaf/nginx/conf/jxwaf/jxwaf_config.json is not exist"
        sys.exit()
    for opt, arg in opts:
        if opt in ['-h', '--help']:
            usage()
            sys.exit()
        elif opt in ['--api_key']:
            api_key = arg
            json_data = ""
            f = open(file_path, 'r')
            json_data = json.loads(f.read())
            json_data['waf_api_key'] = api_key
            f.close()
            ff = open(file_path, 'w')
            ff.write(json.dumps(json_data))
            ff.close()
        elif opt in ['--api_password']:
            api_password = arg
            json_data = ""
            f = open(file_path, 'r')
            json_data = json.loads(f.read())
            json_data['waf_api_password'] = api_password
            f.close()
            ff = open(file_path, 'w')
            ff.write(json.dumps(json_data))
            ff.close()
        elif opt in ['--waf_server']:
            waf_server = arg
            if str(waf_server)[-1] == "/":
                print "Error: waf_server is contain uri path "
                sys.exit()
            waf_server_update = waf_server + "/waf_update"
            waf_server_monitor = waf_server + "/waf_monitor"
            waf_name_list_item_update_website = waf_server + "/waf_name_list_item_update"
            waf_add_name_list_item_website = waf_server + "/api/add_name_list_item"
            f = open(file_path, 'r')
            json_data = json.loads(f.read())
            json_data['waf_update_website'] = waf_server_update
            json_data['waf_monitor_website'] = waf_server_monitor
            json_data['waf_name_list_item_update_website'] = waf_name_list_item_update_website
            json_data['waf_add_name_list_item_website'] = waf_add_name_list_item_website
            json_data['waf_node_hostname'] = socket.gethostname()
            f.close()
            ff = open(file_path, 'w')
            ff.write(json.dumps(json_data))
            ff.close()
        else:
            print "Error: invalid parameters"
            usage()
            sys.exit()
    f = open(file_path, 'r')
    json_data = json.loads(f.read())
    result_api_key = json_data['waf_api_key']
    result_api_password = json_data['waf_api_password']
    waf_update_website = json_data['waf_update_website']
    f.close()
    print "config file:  " + file_path
    print "config result:"
    print "api_key is %s,access_secret is %s " % (result_api_key, result_api_password)
    print "waf_update_website is %s " % (waf_update_website)
    # print  json.dumps(json_data)
    data = {"api_key": result_api_key, "api_password": result_api_password}



if __name__ == '__main__':
    if len(sys.argv) == 1:
        usage()
        sys.exit()
    main(sys.argv)

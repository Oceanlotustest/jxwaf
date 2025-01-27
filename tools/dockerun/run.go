package main

import (
	"os"
	"syscall"
	"log"
	"github.com/google/uuid"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}


var (
	Http_prot = getEnv("HTTP_PORT","80")
	Https_port = getEnv("HTTPs_PORT","443")
	Nginx_config = fmt.Sprintf(`#user  nobody;
worker_processes  auto;

#error_log  logs/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;

#pid        logs/nginx.pid;


worker_rlimit_nofile 102400;
events {
    #multi_accept on;
    worker_connections  10240;
    #use epoll;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    #log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
    #                  '$status $body_bytes_sent "$http_referer" '
    #                  '"$http_user_agent" "$http_x_forwarded_for"';

    #access_log  logs/access.log  main;
    client_body_buffer_size  10m;
    client_max_body_size 100m;
    sendfile        on;
    #tcp_nopush     on;
	resolver  114.114.114.114;
  resolver_timeout 5s;
    #keepalive_timeout  0;
    keepalive_timeout  65;
    lua_ssl_trusted_certificate  /etc/pki/tls/certs/ca-bundle.crt;
    lua_ssl_verify_depth 3;
lua_shared_dict waf_conf_data 100m;
lua_shared_dict jxwaf_sys 100m;
lua_shared_dict jxwaf_limit_req 100m;
lua_shared_dict jxwaf_limit_count 100m;
lua_shared_dict jxwaf_limit_domain 100m;
lua_shared_dict jxwaf_limit_ip_count 100m;
lua_shared_dict jxwaf_limit_ip 100m;
lua_shared_dict jxwaf_limit_bot 100m;
lua_shared_dict jxwaf_flow_block_ip 100m;
lua_shared_dict jxwaf_public 100m;
init_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/init.lua;
init_worker_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/init_worker.lua;
rewrite_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/rewrite.lua;
access_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/access.lua;
#header_filter_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/header_filter.lua;
#body_filter_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/body_filter.lua;
log_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/log.lua;
rewrite_by_lua_no_postpone on;
    #gzip  on;
	upstream jxwaf {
	server www.jxwaf.com;
  balancer_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/balancer.lua;
}
lua_code_cache on;
    server {
        listen       %s;
        server_name  localhost;

        #charset koi8-r;

        #access_log  logs/host.access.log  main;
        set $proxy_pass_https_flag "false";
        location / {
            #root   html;
           # index  index.html index.htm;
          if ($proxy_pass_https_flag = "true"){
            proxy_pass https://jxwaf;
          }
          if ($proxy_pass_https_flag = "false"){
            proxy_pass http://jxwaf;
          }

           proxy_set_header Host  $http_host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }

        #error_page  404              /404.html;

        # redirect server error pages to the static page /50x.html
        #
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
	 #proxy_pass http://www.jxwaf.com;
        }

        # proxy the PHP scripts to Apache listening on 127.0.0.1:80
        #
        #location ~ \.php$ {
        #    proxy_pass   http://127.0.0.1;
        #}

        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
        #
        #location ~ \.php$ {
        #    root           html;
        #    fastcgi_pass   127.0.0.1:9000;
        #    fastcgi_index  index.php;
        #    fastcgi_param  SCRIPT_FILENAME  /scripts$fastcgi_script_name;
        #    include        fastcgi_params;
        #}

        # deny access to .htaccess files, if Apache's document root
        # concurs with nginx's one
        #
        #location ~ /\.ht {
        #    deny  all;
        #}
    }


    # another virtual host using mix of IP-, name-, and port-based configuration
    #
    #server {
    #    listen       8000;
    #    listen       somename:8080;
    #    server_name  somename  alias  another.alias;

    #    location / {
    #        root   html;
    #        index  index.html index.htm;
    #    }
    #}


    # HTTPS server
    #
    #server {
    #    listen       443 ssl;
    #    server_name  localhost;

    #    ssl_certificate      cert.pem;
    #    ssl_certificate_key  cert.key;

    #    ssl_session_cache    shared:SSL:1m;
    #    ssl_session_timeout  5m;

    #    ssl_ciphers  HIGH:!aNULL:!MD5;
    #    ssl_prefer_server_ciphers  on;

    #    location / {
    #        root   html;
    #        index  index.html index.htm;
    #    }
    #}
    server {
        listen       %s ssl;
        server_name  localhost;

        ssl_certificate      full_chain.pem;
        ssl_certificate_key  private.key;

        ssl_session_cache    shared:SSL:1m;
        ssl_session_timeout  5m;
        ssl_session_tickets off;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
        ssl_ciphers "EECDH+AESGCM:EDH+AESGCM:AES256+EECDH:AES256+EDH:ECDHE-RSA-AES128-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA128:DHE-RSA-AES128-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES128-GCM-SHA128:ECDHE-RSA-AES128-SHA384:ECDHE-RSA-AES128-SHA128:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES128-SHA:DHE-RSA-AES128-SHA128:DHE-RSA-AES128-SHA128:DHE-RSA-AES128-SHA:DHE-RSA-AES128-SHA:ECDHE-RSA-DES-CBC3-SHA:EDH-RSA-DES-CBC3-SHA:AES128-GCM-SHA384:AES128-GCM-SHA128:AES128-SHA128:AES128-SHA128:AES128-SHA:AES128-SHA:DES-CBC3-SHA:HIGH:!aNULL:!eNULL:!EXPORT:!DES:!MD5:!PSK:!RC4";
        ssl_prefer_server_ciphers  on;
        ssl_certificate_by_lua_file /opt/jxwaf/lualib/resty/jxwaf/ssl.lua;
        set $proxy_pass_https_flag "false";
        location / {
            root   html;
            index  index.html index.htm;
          if ($proxy_pass_https_flag = "true"){
            proxy_pass https://jxwaf;
          }
          if ($proxy_pass_https_flag = "false"){
            proxy_pass http://jxwaf;
          }
	    proxy_ssl_server_name on;
	    proxy_ssl_name $http_host;
	    proxy_ssl_session_reuse off;
            proxy_set_header Host  $http_host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
    }


}`,Http_prot,Https_port)

)




var (
	file_path = "/opt/jxwaf/nginx/conf/jxwaf/jxwaf_config.json"
	server = getEnv("JXWAF_SERVER","")
	wafapikey = getEnv("WAF_API_KEY","")
	wafapipassword = getEnv("WAF_API_PASSWORD","")

	cmd = "/opt/jxwaf/nginx/sbin/nginx"
	args = []string{
		"jxwaf",
		"-g",
		"daemon off;",
	}
)


func nginx_config() error{
	err := os.WriteFile("/opt/jxwaf/nginx/conf/nginx.conf", []byte(Nginx_config),0644)
	if err != nil {
		log.Fatal(err)
	}
	return err 
}


func evndata()  map[string]string{
	evndata := map[string]string{
		"waf_update_website": server + "/waf_update",
		"waf_monitor_website": server + "/waf_monitor",
		"waf_name_list_item_update_website": server + "/waf_name_list_item_update",
		"waf_add_name_list_item_website": server + "/api/add_name_list_item",
		"waf_node_hostname":"docker_" + uuid.New().String(),
		"waf_api_key": wafapikey,
		"waf_api_password": wafapipassword,

	}
	return evndata 
}

func waf_init(){
	// confmap := evndata()
	err := nginx_config()
	if err != nil {
		log.Fatal(err)
	} 

	var confmap map[string]string
	file, err := os.Open(file_path)
    if err != nil {
       log.Panic("config open err: ",err)
    }
    defer file.Close()
    content, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(content, &confmap)
	if err != nil {
		log.Panic("json marshal err: ",err)
	}
	if confmap["waf_api_key"] == "" {
		confmap := evndata()
		confjson, err := json.Marshal(&confmap)
		if err != nil {
			log.Print(err)
		}
		ioutil.WriteFile(file_path, confjson, 0644)
		log.Print(string(confjson))
	}else{
		hostname := confmap["waf_node_hostname"]
		node_uuid := confmap["waf_node_uuid"]
		confmap := evndata()
		confmap["waf_node_hostname"] = hostname
		confmap["waf_node_uuid"] = node_uuid
		confjson, err := json.Marshal(&confmap)
		if err != nil {
			log.Print(err)
		}
		ioutil.WriteFile(file_path, confjson, 0644)
		log.Print(string(confjson))
	}

	if err := syscall.Exec(cmd, args,os.Environ()); err != nil {
		log.Fatal(err)
	}

}


func main() {
	
	waf_init()
}

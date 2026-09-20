multipass shell app1

sudo apt update
sudo apt install nginx -y


sudo tee /var/www/html/index.html > /dev/null <<'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>App Server 1</title>
</head>
<body>
    <h1>Hello from APP SERVER 1</h1>
    <h2>Server: app1</h2>
</body>
</html>
EOF

sudo systemctl restart nginx

curl localhost

exit

multipass shell app2

sudo apt update
sudo apt install nginx -y


sudo tee /var/www/html/index.html > /dev/null <<'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>App Server 2</title>
</head>
<body>
    <h1>Hello from APP SERVER 2</h1>
    <h2>Server: app2</h2>
</body>
</html>
EOF


sudo systemctl restart nginx

curl localhost

exit

multipass shell lb

sudo apt update
sudo apt install nginx -y

sudo cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup

sudo vi /etc/nginx/conf.d/loadbalancer.conf

upstream backend_servers {

    server 192.168.252.8:80;
    server 192.168.252.9:80;

}

server {

    listen 80;

    location / {

        proxy_pass http://backend_servers;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    }
}

# remove the default site
sudo rm -f /etc/nginx/sites-enabled/default

sudo nginx -t

sudo systemctl restart nginx

sudo systemctl status nginx

exit


for i in {1..10}
do
    echo "Request $i"
    curl -s http://192.168.252.7 | grep Server
done

multipass exec app1 -- sudo systemctl stop nginx

for i in {1..10}
do
    echo "Request $i"
    curl -s http://192.168.252.7 | grep Server
done

# Show weighted load balancing

multipass shell lb

sudo vi /etc/nginx/conf.d/loadbalancer.conf


upstream backend_servers {

    server 192.168.252.8:80 weight=3;

    server 192.168.252.9:80 weight=1;

}

sudo nginx -t
sudo systemctl restart nginx

exit

-- test it

# Show least-connections load balancing


sudo vi /etc/nginx/conf.d/loadbalancer.conf

upstream backend_servers {

    least_conn;

    server 192.168.64.11:80;
    server 192.168.64.12:80;

}


sudo nginx -t
sudo systemctl restart nginx

exit

-- test it

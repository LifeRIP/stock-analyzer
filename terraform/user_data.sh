#!/bin/bash
echo "user_data started!" >> /var/log/user-data.log
yum update -y
yum install -y docker git
systemctl start docker
systemctl enable docker
usermod -aG docker ec2-user

# Docker Compose v2
curl -SL https://github.com/docker/compose/releases/download/v2.34.0/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

# Clona el repo de la aplicación
cd /home/ec2-user
git clone https://github.com/LifeRIP/stock-analyzer app
cd app

# Ejecuta en modo producción
# docker-compose --profile prod up -d

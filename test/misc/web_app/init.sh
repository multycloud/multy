#!/bin/bash -xe
sudo su
yum update -y
curl --silent --location https://rpm.nodesource.com/setup_14.x | bash
yum -y install git nodejs mysql
git clone https://github.com/FaztTech/nodejs-mysql-links.git
cd nodejs-mysql-links
export DATABASE_HOST='${db_host}'
export DATABASE_USER='${db_username}'
export DATABASE_PASSWORD='${db_password}'
# both aws and az will try to run this command but only one will succeed
mysql -h $DATABASE_HOST -P 3306 -u $DATABASE_USER --password=$DATABASE_PASSWORD -e 'source database/db.sql' || true
npm i
npm run build
npm start
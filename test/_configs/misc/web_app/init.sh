#!/bin/bash -xe

{
export DATABASE_HOST='${db_host}'
export DATABASE_USER='${db_username}'
export DATABASE_PASSWORD='${db_password}'
sudo yum update -y
curl --silent --location https://rpm.nodesource.com/setup_14.x | sudo bash
sudo yum -y install git nodejs mysql
sudo chmod a+rwx .
git clone https://github.com/FaztTech/nodejs-mysql-links.git
cd nodejs-mysql-links

# both aws and az will try to run this command but only one will succeed
mysql -h $DATABASE_HOST -P 3306 -u $DATABASE_USER --password=$DATABASE_PASSWORD -e 'source database/db.sql' || true

npm i
npm run build
npm start
} |& tee -a logs.txt
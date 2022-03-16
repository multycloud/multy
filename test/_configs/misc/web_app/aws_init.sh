#!/bin/bash -xe

{
date
sudo yum update -y

sudo yum install epel-release -y
sudo yum install jq -y

region=$(curl -s http://169.254.169.254/latest/meta-data/placement/region)

curl --silent --location https://rpm.nodesource.com/setup_14.x | sudo bash
sudo yum -y install git
sudo yum -y install nodejs
sudo yum -y install mysql
sudo chmod a+rwx .
git clone https://github.com/FaztTech/nodejs-mysql-links.git
cd nodejs-mysql-links

export DATABASE_HOST=$(aws ssm get-parameter --with-decryption --name "/${vault_name}/${db_host_secret_name}" --region "$region" | jq ".Parameter.Value" -r)
export DATABASE_USER=$(aws ssm get-parameter --with-decryption --name "/${vault_name}/${db_username_secret_name}" --region "$region" | jq ".Parameter.Value" -r)
export DATABASE_PASSWORD=$(aws ssm get-parameter --with-decryption --name "/${vault_name}/${db_password_secret_name}" --region "$region" | jq ".Parameter.Value" -r)

# both aws and az will try to run this command but only one will succeed
mysql -h $DATABASE_HOST -P 3306 -u $DATABASE_USER --password=$DATABASE_PASSWORD -e 'source database/db.sql' || true

npm i
npm run build
date
npm start

} |& tee -a logs.txt
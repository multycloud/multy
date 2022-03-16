#!/bin/bash -xe

{
date
sudo yum update -y

sudo yum install epel-release -y
sudo yum install jq -y

# az cli
sudo rpm --import https://packages.microsoft.com/keys/microsoft.asc
sudo sh -c 'echo -e "[azure-cli]
name=Azure CLI
baseurl=https://packages.microsoft.com/yumrepos/azure-cli
enabled=1
gpgcheck=1
gpgkey=https://packages.microsoft.com/keys/microsoft.asc" > /etc/yum.repos.d/azure-cli.repo'
sudo yum install azure-cli -y
az login --identity --allow-no-subscriptions

curl --silent --location https://rpm.nodesource.com/setup_14.x | sudo bash
sudo yum -y install git
sudo yum -y install nodejs
sudo yum -y install mysql
sudo chmod a+rwx .
git clone https://github.com/FaztTech/nodejs-mysql-links.git
cd nodejs-mysql-links

export DATABASE_HOST=$(az keyvault secret show --vault-name '${vault_name}' -n '${db_host_secret_name}' | jq ".value" -r)
export DATABASE_USER=$(az keyvault secret show --vault-name '${vault_name}' -n '${db_username_secret_name}' | jq ".value" -r)
export DATABASE_PASSWORD=$(az keyvault secret show --vault-name '${vault_name}' -n '${db_password_secret_name}' | jq ".value" -r)

# both aws and az will try to run this command but only one will succeed
mysql -h $DATABASE_HOST -P 3306 -u $DATABASE_USER --password=$DATABASE_PASSWORD -e 'source database/db.sql' || true

npm i
npm run build
date
npm start

} |& tee -a logs.txt
systemctl stop primeprice.com
rm -rf /var/www/aws-scraping/backend/primeprice.com
rm -rf /var/www/aws-scraping/backend/text.log
cp ./primeprice.com /var/www/aws-scraping/backend/
systemctl start primeprice.com
systemctl restart mongodb
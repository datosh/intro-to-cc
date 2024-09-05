#!/bin/bash

sudo apt  update
sudo apt -y install apache2

echo '<html><body><h1>Hello Krakow!</h1></body></html>' | sudo tee /var/www/html/index.html

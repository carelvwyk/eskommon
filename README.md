Eskom Mon
=

Eskom is currently under siege by unions demanding higher pay rises (https://www.enca.com/south-africa/eskom-not-moving-from-its-no-wage-increase-stance). This is causing unecessary power blackouts and I would like to know whenever my electricy supply goes down. 
I have a Mecer 2000VU UPS powering a raspberry pi and my wifi router. The pi monitors the UPS and pushes battery percentage to AWS CloudWatch metrics. I configured an alarm on CW to SMS me whenever UPS battery percentage drops below 100%, indicating a power failure.

This project is Go based because I love Go.
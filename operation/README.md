# Operation

The `operation` directory contains all configuration templates and scripts to install and run the festvials-gateway.

* `install.sh` script to install festvials-gateway on a VM
* `service_template.service` festivals gateway unit file for `systemctl`
* `ufw_app_profile` firewall app profile file for `ufw`
* `update.sh` script to update the festivals-gateway

## Deployment

Follow the [**deployment guide**](DEPLOYMENT.md) for deploying the festivals-gateway inside a virtual machine or the [**local deployment guide**](./local/README.md) for running it on your macOS developer machine.

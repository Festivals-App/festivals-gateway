# Development Deployment

This deployment guide explains how to deploy the FestivalsApp Gateway using certificates intended for development purposes.

## Prerequisites

This guide assumes you have already created a Virtual Machine (VM) following the [VM deployment guide](https://github.com/Festivals-App/festivals-documentation/tree/main/deployment/vm-deployment).

Before starting the installation, ensure you have:

- Created and configured your VM
- SSH access secured and logged in as the new admin user
- Your server's IP address (use `ip a` to check)
- A server name matching the Common Name (CN) for your server certificate (e.g., `gateway.festivalsapp.home` for a hostname `gateway`).

I use the development wildcard server certificate (`CN=*festivalsapp.home`) for this guide.

  > **DON'T USE THIS IN PRODUCTION, SEE [festivals-pki](https://github.com/Festivals-App/festivals-pki) FOR SECURITY BEST PRACTICES FOR PRODUCTION**

## 1. Installing the FestivalsApp Gateway

Run the following commands to download and install the FestivalsApp Gateway:

```bash
curl -o install.sh https://raw.githubusercontent.com/Festivals-App/festivals-gateway/master/operation/install.sh
chmod +x install.sh
sudo ./install.sh
```

The config file is located at:

  > `/etc/festivals-gateway.conf`.

You also need to provide certificates in the right format and location:

  > Root CA certificate     `/usr/local/festivals-gateway/ca.crt`  
  > Server certificate is   `/usr/local/festivals-gateway/server.crt`  
  > Server key is           `/usr/local/festivals-gateway/server.key`  

Where the root CA certificate is required to validate incoming requests and the server certificate and key is requires to make outgoing connections.
For instructions on how to manage and create the certificates see the [festivals-pki](https://github.com/Festivals-App/festivals-pki) repository.

## 2. Copying mTLS Certificates to the VM

Copy the server mTLS certificates from your development machine to the VM:

```bash
scp /opt/homebrew/etc/pki/ca.crt <user>@<ip-address>:.
scp /opt/homebrew/etc/pki/issued/server.crt <user>@<ip-address>:.
scp /opt/homebrew/etc/pki/private/server.key <user>@<ip-address>:.
```

Once copied, SSH into the VM and move them to the correct location:

```bash
sudo mv ca.crt /usr/local/festivals-gateway/ca.crt
sudo mv server.crt /usr/local/festivals-gateway/server.crt
sudo mv server.key /usr/local/festivals-gateway/server.key
```

Set the correct permissions:

```bash
# Change owner to web user
sudo chown www-data:www-data /usr/local/festivals-gateway/ca.crt
sudo chown www-data:www-data /usr/local/festivals-gateway/server.crt
sudo chown www-data:www-data /usr/local/festivals-gateway/server.key
# Set secure permissions
sudo chmod 640 /usr/local/festivals-gateway/ca.crt
sudo chmod 640 /usr/local/festivals-gateway/server.crt
sudo chmod 600 /usr/local/festivals-gateway/server.key
```

## 3. Configuring the FestivalsApp Gateway

Open the configuration file:

```bash
sudo nano /etc/festivals-gateway.conf
```

Set the server name, heartbeat endpoint and authentication endpoint:

```ini
[service]
bin-host = "<server name>"
# For example:
# bind-address = "festivalsapp.home"

[heartbeat]
endpoint = "<discovery endpoint>"
#For example: endpoint = "https://discovery.festivalsapp.home/loversear"

[authentication]
endpoint = "<authentication endpoint>"
# endpoint = "https://identity-0.festivalsapp.home:22580"
```

And now let's start the service:

```bash
sudo systemctl start festivals-gateway
```

## **🚀 The gateway should now be running successfully. 🚀**

### Optional: Setting Up DNS Resolution  

For the services in the FestivalsApp backend to function correctly, proper DNS resolution is required.
This is because mTLS is configured to validate the client’s certificate identity based on its DNS hostname.  

If you don’t have a DNS server to manage DNS for your development VMs, you can manually configure DNS resolution
by adding the necessary entries to each server’s `/etc/hosts` file:  

```bash
sudo nano /etc/hosts
```

Add the following entries:  

```ini
<ip address> <server name>
<identity ip address> <heartbeat endpoint name>
<identity ip address> <auth endpoint name>
...

# Example:  
# 192.168.8.186 festivalsapp.home
# 192.168.8.186 discovery.festivalsapp.home
# 192.168.8.185 identity-0.festivalsapp.home
# ...
```

**Keep in mind that you will need to update each machine’s `hosts` file whenever you add a new VM or if any IP addresses change.**

## Testing

Lets login as the default admin user and get the server info:

```bash
curl -H "Api-Key: TEST_API_KEY_001" -u "admin@email.com:we4711" --cert /opt/homebrew/etc/pki/issued/api-client.crt --key /opt/homebrew/etc/pki/private/api-client.key --cacert /opt/homebrew/etc/pki/ca.crt https://identity-0.festivalsapp.home:22580/users/login
```

This should return a JWT Token `<Header.<Payload>.<Signatur>`

  > eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.
  > eyJVc2VySUQiOiIxIiwiVXNlclJvbGUiOjQyLCJVc2VyRmVzdGl2YWxzIjpbXSwiVXNlckFydGlzdHMiOltdLCJVc2VyTG9jYXRpb25zIjpbXSwiVXNlckV2ZW50cyI6W10sIlVzZXJMaW5rcyI6W10sIlVzZXJQbGFjZXMiOltdLCJVc2VySW1hZ2VzIjpbXSwiVXNlclRhZ3MiOltdLCJpc3MiOiJpZGVudGl0eS0wLmZlc3RpdmFsc2FwcC5ob21lIiwiZXhwIjoxNzQwMjMxMTQ4fQ.
  > geBq1pxEvqwjnKA5YTHQ8IjJc9mwkpsQIRy1kGc63oNXzyAhPrPJsepICXxr2yVmB0E8oDECXLn4Cy5V_p4UAduWXnc0r8S05ijV8NCfmsEcJg-RRO8POkGykiC2mrn-XR8Nf8OF0fLp7Mhsb0_aqBoTOLdtB9V7IV49-JjWwX5gHl3HuXGOOhe4n_epumc8w8yDxYakWeaBFtEtaRmhFXK_yttexYOLP6Z1BBTL005hBGhO58qVW0cfgf_t5VWBpUnz3zqdC-GFegItqJQbKZ2pmfmXNz_AoJf2JmPtCzpJ4lG6QeSslvdFuwaCdYpDQPOvnMSIORwrAq_FL2m7qw

Use this to make authorized calls to the Gateway:

```bash
curl -H "Api-Key: TEST_API_KEY_001" -H "Authorization: Bearer <JWT>" --cert /opt/homebrew/etc/pki/issued/api-client.crt --key /opt/homebrew/etc/pki/private/api-client.key --cacert /opt/homebrew/etc/pki/ca.crt https://gateway.festivalsapp.home/info
```

# Running the gateway server locally on you mac

## Prerequisites

- You also need to add the dev root ca certifiate
  
```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/new-root-certificate.crt
```

```bash
make install
```

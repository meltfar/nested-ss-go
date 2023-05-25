# A modified ss to support nested ss protocol
## the text below is just used to remind myself

### What is "Nested"?
- It means that this version should be used combine with other vanilla ss.
- This ss needs you to hold a private ss.

### How
- After pointed to a vanilla ss, this nested ss would provide another layer of encryption: 
1. It would encrypt data traffic with selected encryption method 
2. Then relay the encrypted traffic to the vanilla ss through socks5 protocol
3. traffic would be double-encrypted by the vanilla ss, then relay to the first ss server
4. traffic decrypted by dangerous public ss server, but actually it would still be encrypted because of our nested ss.
5. relay to the real ss backend you provided, decrypt, relay to its real destination
6. enjoy

#!/usr/bin/env python3

import hashlib
from base64 import b64decode as decode

def verify_password(challenge_password, password:str, salt_len:int=4):
    challenge_bytes = decode(challenge_password[6:])
    
    digest = challenge_bytes[:-salt_len]
    salt = challenge_bytes[-salt_len:]

    hr = hashlib.sha1(password.encode('utf-8'))
    hr.update(salt)
    
    return digest == hr.digest()


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(
        formatter_class=argparse.RawDescriptionHelpFormatter,
        description='''Verify hashed SSHA password

passwords from ldapsearch command maybe base64 encoded (attribute name, "userPassword::")
you may need to decode these passwords to see the hash with `base64 -d` command:

> HASHED_PASSWORD=$(echo BASE64_PASSWORD | base64 -d)''')

    parser.add_argument("hashed_secret", type=str, help="Hashed password challange; should start with {SSHA}...")
    parser.add_argument("password", type=str, help="Plain-text password")
    parser.add_argument("--salt", type=int, default=4, help="Salt length for SSHA; default is 4 (SSHA default for `slappasswd` command)")

    args = parser.parse_args()

    challange_psswd = args.hashed_secret
    passwd = args.password
    salt_len = args.salt

    if verify_password(challange_psswd, passwd, salt_len=salt_len):
        print("Password is correct.")
    else:
        print("Password is incorrect.")

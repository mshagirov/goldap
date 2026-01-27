#!/usr/bin/env python3

import hashlib
from base64 import b64decode as decode

def verify_ssha_password(challenge_password, password:str):
    challenge_bytes = decode(challenge_password[6:])
    
    digest = challenge_bytes[:20]
    salt = challenge_bytes[20:]

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

    args = parser.parse_args()

    challange_psswd = args.hashed_secret
    passwd = args.password

    if verify_ssha_password(challange_psswd, passwd):
        print("Password is correct.")
    else:
        print("Password is incorrect.")

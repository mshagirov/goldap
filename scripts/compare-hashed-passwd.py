#!/usr/bin/env python3
import hashlib
from base64 import b64decode as decode

def verify_password(challenge_password, password):
    challenge_bytes = decode(challenge_password[6:])
    
    digest = challenge_bytes[:-4]
    salt = challenge_bytes[-4:]

    hr = hashlib.sha1(password.encode('utf-8'))
    hr.update(salt)
    
    return digest == hr.digest()


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="Verify hashed SSHA password")
    parser.add_argument("hashed_secret")
    parser.add_argument("password")

    args = parser.parse_args()

    challange_psswd = args.hashed_secret
    passwd = args.password
    if verify_password(challange_psswd, passwd):
        print("Password is correct.")
    else:
        print("Password is incorrect.")

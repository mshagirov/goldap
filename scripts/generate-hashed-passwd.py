#!/usr/bin/env python3

import hashlib
from base64 import b64encode as encode
import os

def make_secret(password : str):
    salt = os.urandom(4)

    h = hashlib.sha1((password).encode('utf-8'))
    h.update(salt) # concat salt

    return "{SSHA}" + encode(h.digest() + salt).decode('ascii')


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="Generate hashed SSHA password")
    parser.add_argument("password")

    args = parser.parse_args()

    passwd = args.password

    print(make_secret(passwd))

#!/usr/bin/env python3

import hashlib
from base64 import b64encode as encode
import os

def make_secret(password : str, salt_len : int = 4):
    salt = os.urandom(salt_len)

    h = hashlib.sha1((password).encode('utf-8'))
    h.update(salt) # concat salt

    return "{SSHA}" + encode(h.digest() + salt).decode('ascii')


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="Generate hashed SSHA password")
    parser.add_argument("password", type=str, help="Plain-text password")
    parser.add_argument("--salt", type=int, default=4, help="Salt length for SSHA; default is 4 (SSHA default for `slappasswd` command)")
    args = parser.parse_args()

    passwd = args.password
    salt_len = args.salt

    print(make_secret(passwd, salt_len))

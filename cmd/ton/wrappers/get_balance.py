#!/usr/bin/env python3.7
import sys
import subprocess

def get_nano_grams(workdir,address):
       stdoutdata = subprocess.getoutput(f"{workdir}/wrappers/getaccount {workdir} {address}")

       begin = stdoutdata.find("balance:")

       if begin == -1:
           return False

       begin += 80

       f = stdoutdata[begin:]

       return stdoutdata[begin:f.find(")") + begin]

balance = get_nano_grams(sys.argv[1], sys.argv[2])

if balance == False:
    print("error")
else:
    print(balance)
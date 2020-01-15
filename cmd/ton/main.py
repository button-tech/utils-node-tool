#!/usr/bin/env python3
from fastapi import FastAPI, HTTPException
import subprocess
from pydantic import BaseModel
import uvicorn
import uuid
from starlette.middleware.cors import CORSMiddleware
import asyncio
import os

fift = "./liteclient-build/crypto/fift -I ./ton/crypto/fift/lib/"
lite_client = "./liteclient-build/lite-client/lite-client"

class RunMethod(BaseModel):
    value: str

class ValidatorPubKey(BaseModel):
    value: str

class Boc(BaseModel):
    data: str

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/ton/send")
async def send_boc(boc: Boc):
    result = await create_and_send_boc(boc.data)
    if result == "error":
        raise HTTPException(status_code=500, detail="err")
    return result

@app.get("/seqno/{address}")
async def get_seqno(address: str):
    stdout = await cli_call("getaccount " + address)

    begin = stdout.find("data:")
    if begin == -1:
        raise HTTPException(status_code=500, detail="err")

    begin += 70

    seqno = stdout[begin:begin + 8]

    return {"seqno": seqno}

@app.get("/getBalance/{address}")
async def get_balance(address: str):
    stdout = await cli_call("getaccount " + address)

    balance = parse_stdout(stdout, "account balance is")

    if balance == "error":
        raise HTTPException(status_code=500, detail="err")

    balance = replace_multiple(balance, ["ng", " ", "\n"], "")

    return {"balance": int(balance)}

@app.get("/getAccount/{address}")
async def get_account(address: str):
    if len(address) != 48 and len(address) != 67 and len(address) != 66:
      raise HTTPException(status_code=500, detail="bad address")

    result = await get_address_info(address)
    if result == "error":
      raise HTTPException(status_code=500, detail="err")

    return result

@app.post("/runmethod")
async def runmethod(m: RunMethod):
    stdout = await run_cli_method(m.value)

    if stdout == "error":
        raise HTTPException(status_code=500, detail="bad params")

    return {"result": stdout}


@app.get("/lastTx/{address}")
async def last_tx(address: str):
    result = await get_last_tx(address)
    if result == "error":
        raise HTTPException(status_code=500, detail="bad exec")
    return result


@app.get("/activateCheck/{address}")
async def activate_check(address:str):
    result = await check_contract_code(address)
    return {"result": result}


async def check_contract_code(address: str)->bool:
    stdout = await cli_call("getaccount " + address)

    if stdout.find("code:") == -1:
        return False

    return True

async def get_address_info(address: str) -> dict:
    stdout = await cli_call("getaccount " + address)

    result = parse_stdout(stdout, "got account state for ", "with respect to blocks")
    if result == "error":
        return "error"

    result = result.replace(" ", "").split(":")
    if len(result) != 2:
        return "error"

    bounce = True
    if stdout.find("balance:") == -1:
        bounce = False

    addresses = subprocess.getoutput(f'{fift} -s ./fift-scripts/addresses.fif {result[0]} 0x{result[1]}').split("\n")

    if len(addresses) != 3:
        return "error"

    return {"workchainId": result[0], "bounce": bounce, "fullAddress": result[1], "nonBounceableAddress": addresses[1], "shortAddress": addresses[2]}


async def run_cli_method(params: str) -> list:
    stdout = await cli_call("runmethod " + params)

    result = parse_stdout(stdout, "result:  [ ")

    if result == "error":
       return "error"

    result = replace_multiple(result, ["[", "]", "(", ")"], "")

    return result.split()

async def cli_call(cmd):
    proc = await asyncio.create_subprocess_exec(
        lite_client, f'-c {cmd}',
        stderr=asyncio.subprocess.PIPE)

    data = await proc.stderr.read()

    data = data.decode('ascii').rstrip()

    await proc.wait()

    return data

async def get_last_tx(address: str) -> dict:
    stdout = await cli_call("getaccount " + address)

    result = parse_stdout(stdout, "last transaction ", "account balance is")
    if result == "error":
        raise HTTPException(status_code=500, detail="bad exec")

    result = result.replace("=", "").replace("lt", "").replace("hash", "").split()

    if len(result) != 2:
        return "error"
    if len(result[1]) != 64:
        return "error"

    return {"lt": result[0], "hash": result[1]}

def replace_multiple(mainString, toBeReplaces, newString):
    for elem in toBeReplaces:
        if elem in mainString:
            mainString = mainString.replace(elem, newString)

    return mainString


def parse_stdout(stdout: str, start_phrase: str, end_phrase: str="")-> str:
    begin = stdout.find(start_phrase)
    if begin == -1:
        return "error"
    if end_phrase == "":
        result = stdout[begin + len(start_phrase):]
    else:
        end = stdout.find(end_phrase)
        result = stdout[begin + len(start_phrase):end]

    return result

async def create_and_send_boc(hexData):
    fileName = str(uuid.uuid4().hex)

    text = '''
     B{''' + hexData + '''}
     "''' + fileName + '''.boc"
     tuck
     B>file
     ."(Saved to file " type .")" cr
     '''

    try:
        with open(f'{fileName}.fif', "w") as f:
            f.write(text)
    except:
        return "err"

    os.system(f'{fift} {fileName}.fif')
    await cli_call("sendfile " + fileName + ".boc")
    os.remove(f'./{fileName}.boc')
    os.remove(f'./{fileName}.fif')

    return {"result": "ok"}

if __name__ == "__main__":
   uvicorn.run("main:app", host="0.0.0.0", port=3000, workers=8, loop="asyncio")

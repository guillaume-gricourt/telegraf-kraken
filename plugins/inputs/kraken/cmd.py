data = {"pair": ",".join(pairs)})

url = "https://api.kraken.com/0/public/Ticker"
curl "https://api.kraken.com/0/public/Ticker?pair=XBTUSD"

{"error":[],"result":{"XXBTZUSD":{"a":["16730.30000","11","11.000"],"b":["16730.20000","1","1.000"],"c":["16730.20000","0.00994822"],"v":["719.44918285","993.33983790"],"p":["16671.03757","16650.04408"],"t":[13040,18346],"l":["16550.00000","16550.00000"],"h":["16759.80000","16759.80000"],"o":"16615.00000"},"XRPUSDT":{"a":["0.35315000","800","800.000"],"b":["0.35290000","800","800.000"],"c":["0.35350000","56.48122000"],"v":["906378.81058203","944769.43803412"],"p":["0.32727408","0.32775350"],"t":[871,918],"l":["0.30058000","0.30058000"],"h":["0.35500000","0.35500000"],"o":"0.33760000"}}

session = requests.Session()
self.response = session.post(
            url,
            data=data,
            headers={},
            timeout=None,

        if self.response.status_code not in (200, 201, 202):
            self.response.raise_for_status()

curl "https://api.kraken.com/0/public/Ticker?pair=XBTUSD"

{"error":[],"result":{"XXBTZUSD":
    {"a":["16511.70000","1","1.000"],
     "b":["16511.60000","2","2.000"],
     "c":["16511.60000","0.00040000"],
     "v":["663.22810627","2172.02438751"],
     "p":["16518.01733","16548.75222"],
     "t":[5562,19311],
     "l":["16490.40000","16464.10000"],
     "h":["16542.40000","16629.70000"],
     "o":"16528.70000"}}}

error : array of string
type: str

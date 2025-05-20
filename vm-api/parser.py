import datetime as dt

filePath="/root/VirtualizationServer/database/tuna_ports"
f = open(filePath)
now = dt.datetime.now()
hhmmss = list(map(int, str(now).split()[1].split(".")[0].split(":")))
data = f.readlines()
for string in data:
    parts = string.split()
    tm = list(map(int, parts[0].split("[")[-1][:-1:].split(":")))
    if len(tm) != 3: continue
    if tm[0] == hhmmss[0] and (tm[1] - hhmmss[1] < 2):
        if parts[1] == "Forwarding":
            port = int(parts[2].split(":")[-1])
            print(port)
            break
f.close()
with open(filePath, "r+") as f:
    f.truncate(0)
f.close()

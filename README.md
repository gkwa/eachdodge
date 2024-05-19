# eachdodge

Purpose:


## example usage

```bash

root@ip-172-31-29-172:~/eachdodge# out=$(mktemp /tmp/eachdodge-XXXX.json)
root@ip-172-31-29-172:~/eachdodge# eachdodge ips2 --outfile=$out
root@ip-172-31-29-172:~/eachdodge# jq -r '.[] | select(.ipVersion == "IPv4" and .interface != "lo" and .isInterface == true) | .ip' $out | jq -R . | jq -s 'join(",")' | tr -d '"'
172.31.29.172,10.0.3.1,172.17.0.1
root@ip-172-31-29-172:~/eachdodge#

```

## install eachdodge


on macos/linux:
```bash

brew install taylormonacelli/homebrew-tools/eachdodge

```


on windows:

```powershell

TBD

```

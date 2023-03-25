#!/bin/bash

set -e

# Change this according your portfolio
pairs=(BTC ADA ETH SOL BNB DOT AVAX LINK FTM MATIC ROSE MANA SAND NEAR AUDIO)
timeframe=1d
days=751
start=2022-04-14
end=2022-11-10

for pair in ${pairs[@]}; do
  ninjabot download --pair ${pair}BUSD --timeframe $timeframe --output testdata/${pair}BUSD-${timeframe}.csv --start $start --end $end #--days $days
done

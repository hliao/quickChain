#!/usr/bin/env bash
set -euo pipefail

# Config
N="${1:-100}"                               # number of chunks (default 100)
RPC="${RPC:-http://127.0.0.1:26657}"        # CometBFT RPC endpoint
OUT="${OUT:-.}"                             # output dir (default: current dir)

# Utils
b64dec() { (echo -n "$1" | base64 -D 2>/dev/null) || (echo -n "$1" | base64 -d 2>/dev/null); }
b64out() { jq -r '.result.tx_result.data // .result.deliver_tx.data'; }
codeout() { jq -r '.result.tx_result.code // .result.deliver_tx.code'; }
now_human() { date '+%Y-%m-%d %H:%M:%S%z'; }
now_epoch() { date +%s; }

SCRIPT_T0=$(now_epoch); SCRIPT_T0_H=$(now_human)

S1_T0=$(now_epoch); S1_T0_H=$(now_human)
echo "Step 1/5: Health check -> $RPC/status [start=$S1_T0_H]"
if curl -s "$RPC/status" >/dev/null; then
  echo "  OK"
else
  echo "  WARN: $RPC/status not reachable; continuing..."
fi
S1_T1=$(now_epoch); S1_T1_H=$(now_human); echo "Step 1/5: done [end=$S1_T1_H duration=$((S1_T1-S1_T0))s]"

S2_T0=$(now_epoch); S2_T0_H=$(now_human)
echo "Step 2/5: Prepare output dir [start=$S2_T0_H]"
echo "  Dir: $OUT"
mkdir -p "$OUT"
S2_T1=$(now_epoch); S2_T1_H=$(now_human); echo "Step 2/5: done [end=$S2_T1_H duration=$((S2_T1-S2_T0))s]"

S3_T0=$(now_epoch); S3_T0_H=$(now_human)
echo "Step 3/5: Generate ${N} chunks (1024 bytes each) [start=$S3_T0_H]"
for ((i=1;i<=N;i++)); do
  FILE="$OUT/chunk_${i}.bin"
  dd if=/dev/urandom of="$FILE" bs=1024 count=1 status=none
  if (( i % 10 == 0 || i == N )); then
    NOW_H=$(now_human); NOW_E=$(now_epoch)
    printf "  Generated %d/%d (%.0f%%) [t=%s elapsed=%ss]\n" "$i" "$N" "$(awk -v a=$i -v b=$N 'BEGIN{print (a*100.0)/b}')" "$NOW_H" "$((NOW_E-S3_T0))"
  fi
done
S3_T1=$(now_epoch); S3_T1_H=$(now_human); echo "Step 3/5: done [end=$S3_T1_H duration=$((S3_T1-S3_T0))s]"

S4_T0=$(now_epoch); S4_T0_H=$(now_human)
echo "Step 4/5: Broadcast ${N} transactions [start=$S4_T0_H]"
> "$OUT/keys.txt"
sent_ok=0; sent_fail=0
for ((i=1;i<=N;i++)); do
  FILE="$OUT/chunk_${i}.bin"
  HEX=$(xxd -p -c 1024 "$FILE" | tr -d '\n')
  RESP="$OUT/resp_${i}.json"
  curl -s "$RPC/broadcast_tx_commit?tx=0x$HEX" -o "$RESP"

  CODE=$(codeout < "$RESP")
  if [[ "$CODE" != "0" ]]; then
    sent_fail=$((sent_fail+1))
  else
    B64=$(b64out < "$RESP")
    KEY=$(b64dec "$B64")
    echo "$KEY" >> "$OUT/keys.txt"
    sent_ok=$((sent_ok+1))
  fi

  if (( i % 5 == 0 || i == N )); then
    NOW_H=$(now_human); NOW_E=$(now_epoch)
    printf "  Broadcast %d/%d | ok=%d fail=%d [t=%s elapsed=%ss]\n" "$i" "$N" "$sent_ok" "$sent_fail" "$NOW_H" "$((NOW_E-S4_T0))"
  fi
done
S4_T1=$(now_epoch); S4_T1_H=$(now_human); echo "Step 4/5: done [end=$S4_T1_H duration=$((S4_T1-S4_T0))s]"

S5_T0=$(now_epoch); S5_T0_H=$(now_human)
echo "Step 5/5: Query and verify ${sent_ok} keys [start=$S5_T0_H]"
q_ok=0; q_fail=0
lineno=0
while IFS= read -r KEY; do
  lineno=$((lineno+1))
  FILE="$OUT/chunk_${lineno}.bin"               # same order as sent
  KEY_HEX=$(printf "%s" "$KEY" | xxd -p -c 256)
  RET="$OUT/retrieved_${lineno}.bin"

  if ! curl -s "$RPC/abci_query?path=%22/get%22&data=0x${KEY_HEX}" \
    | jq -r '.result.response.value' \
    | (base64 -D 2>/dev/null > "$RET" || base64 -d > "$RET"); then
    q_fail=$((q_fail+1))
  else
    if cmp -s "$FILE" "$RET"; then
      q_ok=$((q_ok+1))
    else
      q_fail=$((q_fail+1))
    fi
  fi

  if (( lineno % 5 == 0 || lineno == sent_ok )); then
    NOW_H=$(now_human); NOW_E=$(now_epoch)
    printf "  Query %d/%d | ok=%d fail=%d [t=%s elapsed=%ss]\n" "$lineno" "$sent_ok" "$q_ok" "$q_fail" "$NOW_H" "$((NOW_E-S5_T0))"
  fi
done < "$OUT/keys.txt"
S5_T1=$(now_epoch); S5_T1_H=$(now_human); echo "Step 5/5: done [end=$S5_T1_H duration=$((S5_T1-S5_T0))s]"

SCRIPT_T1=$(now_epoch); SCRIPT_T1_H=$(now_human)
echo "Done. [start=$SCRIPT_T0_H end=$SCRIPT_T1_H duration=$((SCRIPT_T1-SCRIPT_T0))s]"
echo "Summary: sent_ok=$sent_ok sent_fail=$sent_fail | query_ok=$q_ok query_fail=$q_fail"
echo "Artifacts: $OUT"




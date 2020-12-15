const axios = require("axios");

async function printRatio(producer, votes) {
  const url = `https://www.iostabc.com/api/account/${producer}/actions?page=1&size=100`;
  const result = await axios.get(url);
  let first = null;
  let second = null;
  let last = null;
  let idx = 0;
  while (idx < result.data.actions.length) {
    const entry = result.data.actions[idx];
    if (
      entry.contract == "token.iost" &&
      entry.action_name == "destroy" &&
      entry.data.startsWith('["contribute","')
    ) {
      if (first === null) {
        first = entry;
      } else if (second == null) {
        second = entry;
      }
    }
    idx++;
  }

  if (first === null || second === null || first.tx_hash === second.tx_hash) {
    return;
  }
  const j = JSON.parse(first.data);
  const vol = Number(j[2]);
  const yearReward =
    vol *
    ((365 * 24 * 3600 * 2) / (Number(first.block) - Number(second.block)));
  console.log(
    `${producer},${Math.floor(votes)},${((yearReward / votes) * 100).toFixed(
      2
    )}%`
  );
}

async function main() {
  try {
    var idx = 1;
    while (idx < 6) {
      const url =
        "https://www.iostabc.com/api/producers?page=" +
        idx +
        "&size=100&sort_by=votes&order=desc&search=";
      const result = await axios.get(url);
      for (idx in result.data.producers) {
        const item = result.data.producers[idx];
        if (
          item.block_count > 0 &&
          item.isProducer &&
          item.online &&
          item.votes > 10000000
        ) {
          await printRatio(item.account, item.votes);
        }
      }
      idx += 1;
      break; // only one page
    }
  } catch (e) {
    console.log("CATCH", e);
  }
}

main();

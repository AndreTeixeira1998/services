const { CryptoService } = require("micro-js-client/crypto");

// Get news related to a currency
async function getCryptocurrencyNews() {
  let cryptoService = new CryptoService(process.env.MICRO_API_TOKEN);
  let rsp = await cryptoService.news({
    symbol: "BTCUSD",
  });
  console.log(rsp);
}

getCryptocurrencyNews();

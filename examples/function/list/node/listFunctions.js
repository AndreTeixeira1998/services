const { FunctionService } = require("micro-js-client/function");

// List all the deployed functions
async function listFunctions() {
  let functionService = new FunctionService(process.env.MICRO_API_TOKEN);
  let rsp = await functionService.list({});
  console.log(rsp);
}

listFunctions();

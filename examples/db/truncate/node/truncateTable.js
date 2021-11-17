const { DbService } = require("micro-js-client/db");

// Truncate the records in a table
async function truncateTable() {
  let dbService = new DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.truncate({
    table: "users",
  });
  console.log(rsp);
}

truncateTable();
